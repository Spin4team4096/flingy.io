package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

func startWebsocketServer() {
	http.HandleFunc("/ws", handleConnection)
	go http.ListenAndServe("127.0.0.1:4096", nil)
	log.Println("Started websocket server")
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("connection panic:", rec)
		}
	}()
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"flingy.io", "beta.flingy.io", "localhost:*"},
	})
	if err != nil {
		log.Println("websocket upgrade failed:", err)
		return
	}
	defer conn.CloseNow()

	conn.SetReadLimit(1024)

	ctx := context.Background()

	c := &clientData{
		conn: conn,
		ctx:  ctx,
		slot: -1,
		queue: queue{
			control: make(chan []byte, 8),
			event:   make(chan []byte, 64),
			state:   make(chan []byte, 1),
		},
		done: make(chan struct{}),
	}
	go c.sendLoop()
	defer close(c.done)
	for {
		_, data, err := conn.Read(ctx) // the _ is msgType, i assume we don't care bout that
		if err != nil {
			if c.slot >= 0 {
				lobbiesMu.Lock()
				if lobby := lobbies[c.lobbyID]; lobby != nil {
					despawnPlayer(lobby, int(c.slot))
					delete(lobby.Clients, uint8(c.slot))
				}
				lobbiesMu.Unlock()
			}
			return // client disconnected
		}
		if err := decodePacket(c, data); err != nil {
			log.Println("decode error:", err)
			continue
		}
	}
}

// PACKET DOCS

// 1: Welcome: Sends PlayerSlot when they join a lobby
// 2: Disconnect: Tells that player to fuck off
//		1 = nonexistent lobby
//		2 = full
//		3 = duplicate
// 3: Lobby Data: Sends data about the game
// 4: One Username: When a player joins, send that username to all clients
// 5: All usernames: When game state changes (0>1, 1>0, 1>2), send all usernames

// -1: Join Request: Adds player to requested lobby
// -2: Client Data: Client inputs go here

func encodePacket(packetType int, data any) []byte {
	version := 1 // REMEMBER TO KEEP SYNCED WITH CLIENT
	switch packetType {
	case 1:
		w := newPacketWriter(2 + 2)
		w.u8(byte(version))
		w.u8(byte(packetType))
		w.u16(uint16(data.(uint16)))
		return w.done()
	case 2:
		w := newPacketWriter(2 + 1)
		w.u8(byte(version))
		w.u8(byte(packetType))
		w.u8(uint8(data.(uint8)))
		return w.done()
	case 3:
		lobby := data.(*Lobby)

		playerCount := 0
		playersAlive := 0
		for i := range lobby.Player {
			if lobby.Player[i].Active {
				playerCount++
			}
			if lobby.Player[i].HP > 0 {
				playersAlive++
			}
		}

		shortenLeaderboard := false
		lowLimit := false
		count := 0
		if playerCount > 64 {
			// fucking shit we gotta compress this leaderboard lol
			shortenLeaderboard = true
			count = playersAlive
		} else {
			count = playerCount
		}
		if lobby.Config.PlayerLimit <= 255 {
			lowLimit = true
		}

		size := 2 + 12 + (29 * count)
		if !lowLimit {
			size += 2 * count
		}
		w := newPacketWriter(size)
		w.u8(byte(version))
		w.u8(byte(packetType))

		// HEADER
		w.u32(uint32(lobby.Tick))
		border := float32(lobby.BorderRadius)
		w.f32(border)
		borderLimit := (float64(border) + 50)
		borderScale := math.Floor(32767 / borderLimit)
		mouseLimit := (float64(border) + 512)
		mouseScale := math.Floor(32767 / mouseLimit)
		w.u8(byte(clamp(30*8, 0, 255))) // add tps counter here

		flags := byte(0)
		if lobby.Closing {
			flags |= 1 << 0
		}
		if shortenLeaderboard {
			flags |= 1 << 1
		}
		if lowLimit {
			flags |= 1 << 2
		}
		flags |= (byte(lobby.State) & 0b11) << 3
		w.u8(byte(flags))

		// PLAYERS
		w.u16(uint16(count))
		for i := range lobby.Player {
			if !shortenLeaderboard && !lobby.Player[i].Active {
				continue
			}
			if shortenLeaderboard && lobby.Player[i].HP <= 0 {
				continue
			}

			p := lobby.Player[i]

			if lowLimit {
				w.u8(byte(i))
				w.u8(byte(p.Rank))
			} else {
				w.u16(uint16(i))
				w.u16(uint16(p.Rank))
			}

			w.i16(int16(math.Round(clamp(p.X, -borderLimit, borderLimit) * borderScale)))
			w.i16(int16(math.Round(clamp(p.Y, -borderLimit, borderLimit) * borderScale)))
			w.i16(int16(math.Round(p.SpeedX * 800))) // range +-40
			w.i16(int16(math.Round(p.SpeedY * 800)))
			w.i16(int16(math.Round(clamp(p.FlailX, -borderLimit, borderLimit) * borderScale)))
			w.i16(int16(math.Round(clamp(p.FlailY, -borderLimit, borderLimit) * borderScale)))
			w.i16(int16(math.Round(p.FlailSpeedX * 400))) // range +-80
			w.i16(int16(math.Round(p.FlailSpeedY * 400)))
			w.i16(int16(math.Round(math.Remainder(p.Direction, 360) * 182)))
			w.i16(int16(math.Round(p.AngularVelocity * 256))) // range +-128
			w.i16(int16(math.Round(clamp(p.WorldMouseX, -mouseLimit, mouseLimit) * mouseScale)))
			w.i16(int16(math.Round(clamp(p.WorldMouseY, -mouseLimit, mouseLimit) * mouseScale)))
			w.u8(byte(math.Round(clamp(p.HP, 0, 31.875) * 8)))
			w.u8(byte(math.Round(clamp(p.FlailRadius, 0, 51) * 5)))

			flags := byte(0)
			if p.Alive {
				flags |= 1 << 0
			}
			if p.FlailConnected {
				flags |= 1 << 1
			}
			if p.Immune {
				flags |= 1 << 2
			}
			flags |= byte(p.ColorIndex&0b11111) << 3
			w.u8(flags)
		}

		return w.done()
	case 4:
		d := data.(usernamePacket)

		w := newPacketWriter(2 + 2 + 1 + len(d.username))
		w.u8(byte(version))
		w.u8(byte(packetType))
		w.u16(uint16(d.slot))
		w.str(string(d.username))
		return w.done()
	case 5:
		lobby := data.(*Lobby)

		size := 2 + 2
		for i := range lobby.Player {
			size += 1 + len(lobby.Player[i].Username)
		}
		w := newPacketWriter(size)
		w.u8(byte(version))
		w.u8(byte(packetType))

		w.u16(uint16(lobby.Config.PlayerLimit))
		for i := range lobby.Player {
			w.str(string(lobby.Player[i].Username))
		}

		return w.done()

	}
	return nil
}

func decodePacket(c *clientData, data []byte) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panicked decoding packet: %v", rec)
		}
	}()
	if len(data) < 2 {
		return fmt.Errorf("packet too short to even read")
	}
	r := newPacketReader(data)
	if version := r.u8(); version == 1 { // REMEMBER TO KEEP SYNCED WITH CLIENT
		packetType := r.i8()
		switch packetType {
		case -1:
			lobbyID := r.str()
			username := r.str()
			runes := []rune(username)
			if len(runes) > 20 {
				runes = runes[:20]
			}
			username = string(runes)

			if c.slot >= 0 {
				c.send(1, encodePacket(2, uint8(3)))
				return nil
			}

			lobbiesMu.Lock()
			lobby := lobbies[lobbyID]
			slot := -1
			playerCount := 0
			if lobby != nil {
				for i := range lobby.Player {
					if lobby.Player[i].Active {
						playerCount++
						continue
					}
					if slot == -1 {
						spawnPlayer(lobby, i)
						lobby.Player[i].Username = username
						lobby.Player[i].ColorIndex = 0
						slot = i
						lobby.Clients[uint8(slot)] = c
						playerCount++
					}
				}
			} else {
				lobbiesMu.Unlock()
				c.send(1, encodePacket(2, uint8(1)))
				return nil
			}

			if slot < 0 {
				lobbiesMu.Unlock()
				c.send(1, encodePacket(2, uint8(2)))
				return nil
			}

			c.slot = int16(slot)
			c.lobbyID = lobbyID
			c.send(1, encodePacket(5, lobby))

			if playerCount < 64 {
				for _, other := range lobby.Clients {
					if other != c {
						other.send(2, encodePacket(4, usernamePacket{slot: uint16(slot), username: username}))
					}
				}
			}

			c.send(1, encodePacket(1, uint16(slot)))
			lobbiesMu.Unlock()
			return nil
		case -2:
			flags := r.u8()
			down := (flags)&1 != 0
			color := (flags >> 1) & 0b11111

			x := float64(r.f32())
			y := float64(r.f32())

			if !isFinite(x) || !isFinite(y) {
				return nil
			}

			lobbiesMu.Lock()
			lobby := lobbies[c.lobbyID]
			slot := c.slot

			lobby.Player[slot].LeftMouseDown = down
			lobby.Player[slot].ColorIndex = int(color)
			lobby.Player[slot].WorldMouseX = x
			lobby.Player[slot].WorldMouseY = y
			lobbiesMu.Unlock()

			return nil
		default:
			return fmt.Errorf("bad packet type: %d", packetType)
		}
	} else {
		return fmt.Errorf("bad version: %d", version)
	}
}

func (c *clientData) send(priority int, pkt []byte) {
	switch priority {
	case 1:
		select {
		case c.queue.control <- pkt:
		default:
			log.Println("control queue fail")
		}
	case 2:
		select {
		case c.queue.event <- pkt:
		default:
			select {
			case <-c.queue.event:
			default:
			}
			select {
			case c.queue.event <- pkt:
			default:
			}
		}
	case 3:
		select {
		case c.queue.state <- pkt:
		default:
			select {
			case <-c.queue.state:
			default:
			}
			select {
			case c.queue.state <- pkt:
			default:
			}
		}
	}
}

func (c *clientData) write(pkt []byte) bool {
	ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
	defer cancel()
	if err := c.conn.Write(ctx, websocket.MessageBinary, pkt); err != nil {
		c.conn.CloseNow()
		return false
	}
	return true
}

func (c *clientData) sendLoop() {
	for {
		select {
		case pkt := <-c.queue.control:
			if !c.write(pkt) {
				return
			}
		default:
		}
		select {
		case pkt := <-c.queue.event:
			if !c.write(pkt) {
				return
			}
		default:
		}
		select {
		case pkt := <-c.queue.state:
			if !c.write(pkt) {
				return
			}
		default:
		}
		select {
		case pkt := <-c.queue.control:
			if !c.write(pkt) {
				return
			}
		case pkt := <-c.queue.event:
			if !c.write(pkt) {
				return
			}
		case pkt := <-c.queue.state:
			if !c.write(pkt) {
				return
			}
		case <-c.done:
			return
		}

	}
}

func firebaseLoop() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		type pub struct {
			id          string
			playerCount int
			status      int
		}
		updates := make([]pub, 0, len(lobbies))
		lobbiesMu.Lock()
		for id, lobby := range lobbies {
			playerCount := 0
			playersAlive := 0
			for i := range lobby.Player {
				if lobby.Player[i].Active {
					playerCount++
				}
				if lobby.Player[i].HP > 0 {
					playersAlive++
				}
			}

			status := 0
			switch lobby.State {
			case 0:
				status = 0
			case 1:
				status = 0
			case 2:
				if playerCount == playersAlive {
					status = 0
				} else {
					status = 1
				}
			case 3:
				status = 1
			}
			if status != lobby.PublishedStatus || time.Since(lobby.LastUpdate) >= 2*time.Minute {
				lobby.PublishedStatus = status
				lobby.LastUpdate = time.Now()
				updates = append(updates, pub{id, playerCount, status})
			}
		}
		lobbiesMu.Unlock()

		for _, u := range updates {
			client.NewRef("flingy2/lobbies/"+u.id).Update(ctx, map[string]any{
				"playerCount": u.playerCount,
				"status":      u.status,
				"lastSeen":    time.Now().UnixMilli(),
			})
		}
	}
}

// Structs

type clientData struct {
	conn    *websocket.Conn
	ctx     context.Context
	lobbyID string
	slot    int16
	queue   queue
	done    chan struct{}
}

type queue struct {
	control chan []byte
	event   chan []byte
	state   chan []byte
}

type notifParams struct {
	notifID int
	variant int
	p1      uint8
	p2      uint8
	text    string
	color   uint8
}

type usernamePacket struct {
	slot     uint16
	username string
}

// Specific Functions

func notify(l *Lobby, slot int, p notifParams) {
	// packet := encodePacket(2, p)

	// then actually send data to all players

	// but for now, stub:
	log.Printf("would notify: all=%v slot=%d notifID=%d variant=%d params=%+v", l == nil, slot, p.notifID, p.variant, p)
}

// Byte Helpers

type packetWriter struct {
	buffer []byte
	offset int
}

func newPacketWriter(size int) *packetWriter {
	return &packetWriter{buffer: make([]byte, size)}
}

func (w *packetWriter) u8(v byte) {
	w.buffer[w.offset] = v
	w.offset++
}

func (w *packetWriter) i8(v int8) {
	w.buffer[w.offset] = byte(v)
	w.offset++
}

func (w *packetWriter) u16(v uint16) {
	binary.LittleEndian.PutUint16(w.buffer[w.offset:], v)
	w.offset += 2
}

func (w *packetWriter) i16(v int16) {
	binary.LittleEndian.PutUint16(w.buffer[w.offset:], uint16(v))
	w.offset += 2
}

func (w *packetWriter) u32(v uint32) {
	binary.LittleEndian.PutUint32(w.buffer[w.offset:], v)
	w.offset += 4
}

func (w *packetWriter) i32(v int32) {
	binary.LittleEndian.PutUint32(w.buffer[w.offset:], uint32(v))
	w.offset += 4
}

func (w *packetWriter) f32(v float32) {
	binary.LittleEndian.PutUint32(w.buffer[w.offset:], math.Float32bits(v))
	w.offset += 4
}

func (w *packetWriter) str(s string) {
	w.buffer[w.offset] = byte(len(s))
	w.offset++
	copy(w.buffer[w.offset:], s)
	w.offset += len(s)
}

func (w *packetWriter) done() []byte {
	return w.buffer[:w.offset]
}

type packetReader struct {
	data   []byte
	offset int
}

func newPacketReader(data []byte) *packetReader {
	return &packetReader{data: data}
}

func (r *packetReader) u8() byte {
	v := r.data[r.offset]
	r.offset++
	return v
}

func (r *packetReader) i8() int8 {
	v := int8(r.data[r.offset])
	r.offset++
	return v
}

func (r *packetReader) u16() uint16 {
	v := binary.LittleEndian.Uint16(r.data[r.offset:])
	r.offset += 2
	return v
}

func (r *packetReader) i16() int16 {
	v := int16(binary.LittleEndian.Uint16(r.data[r.offset:]))
	r.offset += 2
	return v
}

func (r *packetReader) u32() uint32 {
	v := binary.LittleEndian.Uint32(r.data[r.offset:])
	r.offset += 4
	return v
}

func (r *packetReader) i32() int32 {
	v := int32(binary.LittleEndian.Uint32(r.data[r.offset:]))
	r.offset += 4
	return v
}

func (r *packetReader) f32() float32 {
	v := math.Float32frombits(binary.LittleEndian.Uint32(r.data[r.offset:]))
	r.offset += 4
	return v
}

func (r *packetReader) str() string {
	n := int(r.data[r.offset])
	r.offset++
	s := string(r.data[r.offset : r.offset+n])
	r.offset += n
	return s
}

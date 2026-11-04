package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Structs
type LobbyConfig struct {
	ID            string
	Name          string
	PlayerLimit   int
	BorderSize    int
	MatchTime     int
	CountdownTime int
	BotSpinSpeed  int
	BotAccuracy   int
	BotPercentage float32
	AddMoney      bool
}
type Player struct {
	Active                   bool
	Alive                    bool
	HP                       float32
	Username                 string
	X, Y                     float32
	SpeedX, SpeedY           float32
	Size                     float32
	FlailX, FlailY           float32
	FlailSpeedX, FlailSpeedY float32
	FlailConnected           bool
	FlailSize                float32
	Direction                float32
	AngularVelocity          float32
	WorldMouseX, WorldMouseY float32
	Rank                     int
	Immune                   bool
	Color                    string
}
type Lobby struct {
	Config          LobbyConfig
	Tick            int
	Running         bool
	Closing         bool
	State           int // 0 = waiting, 1 = countdown, 2 = playing, 3 = ending
	StateTick       int
	PublishedStatus int // 0 = open, 1 = busy
	LastUpdate      time.Time
	Player          []Player
	Clients         map[uint8]*clientData
}

type ServerEntry struct {
	URL         string `json:"url"`
	Status      int    `json:"status"`
	PlayerCount int    `json:"playerCount"`
	CreatedAt   int64  `json:"createdAt"`
	LastSeen    int64  `json:"lastSeen"`
}

var (
	lobbies   = map[string]*Lobby{}
	lobbiesMu sync.Mutex
)

// Specific Functions
func createLobby(config LobbyConfig) {
	if err := registerLobby(config); err != nil {
		log.Println("Could not register lobby:", err)
		return
	}

	lobbiesMu.Lock()
	if lobbies[config.ID] == nil {
		lobbies[config.ID] = &Lobby{
			Config:     config,
			Running:    true,
			LastUpdate: time.Now(),
			Player:     make([]Player, config.PlayerLimit),
			Clients:    make(map[uint8]*clientData),
		}
	} else {
		log.Printf("Lobby %s already exists in memory", config.ID)
	}
	lobbiesMu.Unlock()
}

func endLobby(id string) {
	lobbiesMu.Lock()
	if lobbies[id] == nil {
		lobbiesMu.Unlock()
		return
	}
	lobbies[id].Closing = true
	lobbiesMu.Unlock()
	// broadcast("Server will close after this match", "#ff6666", -1)
}

func killLobby(id string) {
	lobbiesMu.Lock()
	if lobbies[id] == nil {
		lobbiesMu.Unlock()
		return
	}
	delete(lobbies, id)
	lobbiesMu.Unlock()
	unregisterLobby(id)
}

func unregisterLobby(id string) {
	serversRef := client.NewRef("flingy2/lobbies")
	if err := serversRef.Child(id).Delete(ctx); err != nil {
		log.Println("Error removing lobby from firebase:", err)
	}
}

func registerLobby(config LobbyConfig) error {
	serversRef := client.NewRef("flingy2/lobbies")
	var servers map[string]ServerEntry
	if err := serversRef.Get(ctx, &servers); err != nil {
		return fmt.Errorf("error reading servers: %w", err)
	}
	if _, ok := servers[config.ID]; ok {
		return fmt.Errorf("lobby %s already registered", config.ID)
	} else {
		entry := ServerEntry{
			URL:         os.Getenv("WS_URL"),
			Status:      0,
			PlayerCount: 0,
			CreatedAt:   time.Now().UnixMilli(),
			LastSeen:    time.Now().UnixMilli(),
		}
		if err := serversRef.Child(config.ID).Set(ctx, entry); err != nil {
			return fmt.Errorf("error adding server: %w", err)
		}
	}
	return nil
}

func unregisterAll() {
	serversRef := client.NewRef("flingy2/lobbies")
	if err := serversRef.Delete(ctx); err != nil {
		log.Println("Error resetting firebase:", err)
	}
}

type statusEntry struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

func setStatus(text string, color string) {
	serversRef := client.NewRef("flingy2/status")
	entry := statusEntry{
		Text:  text,
		Color: color,
	}
	if err := serversRef.Set(ctx, entry); err != nil {
		log.Println("failed to set status:", err)
	} else {
		log.Printf("set status to %s (%s)", text, color)
	}
}

func spawnPlayer(l *Lobby, slot int) {
	p := &l.Player[slot]

	p.Active = true
	p.Alive = true
	p.HP = 3
	p.X = float32(randomFloat(float64(-l.Config.BorderSize), float64(l.Config.BorderSize)))
	p.Y = float32(randomFloat(float64(-l.Config.BorderSize), float64(l.Config.BorderSize)))
	p.SpeedX = 0
	p.SpeedY = 0
	p.Size = 10
	p.FlailX = p.X
	p.FlailY = p.Y
	p.FlailSpeedX = 0
	p.FlailSpeedY = 0
	p.FlailConnected = true
	p.FlailSize = 5
	p.Direction = 0
	p.AngularVelocity = 0
	p.WorldMouseX = 0
	p.WorldMouseY = 0
	p.Rank = 0
	p.Immune = true
}

func despawnPlayer(l *Lobby, slot int) {
	l.Player[slot] = Player{}
}

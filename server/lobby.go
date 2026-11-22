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
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	PlayerLimit   int     `json:"playerLimit"`
	BorderRadius  int     `json:"borderRadius"`
	MatchTime     int     `json:"matchTime"`
	CountdownTime int     `json:"countdownTime"`
	BotSpinSpeed  int     `json:"botSpinSpeed"`
	BotAccuracy   int     `json:"botAccuracy"`
	BotPercentage float64 `json:"botPercentage"`
	AddMoney      bool    `json:"addMoney"`
}
type Player struct {
	Active                   bool
	Alive                    bool
	HP                       float64
	Username                 string
	X, Y                     float64
	SpeedX, SpeedY           float64
	Radius                   float64
	FlailX, FlailY           float64
	FlailSpeedX, FlailSpeedY float64
	FlailConnected           bool
	FlailRadius              float64
	Direction                float64
	AngularVelocity          float64
	WorldMouseX, WorldMouseY float64
	Rank                     int
	Immune                   bool
	ColorIndex               int

	LeftMouseDown    bool
	FlingUsed        bool
	TimeLeftImmune   float64
	TimeUntilRespawn float64
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
	BorderRadius    float64
	Player          []Player
	Clients         map[uint8]*clientData
}

type ServerEntry struct {
	URL         string      `json:"url"`
	Status      int         `json:"status"`
	PlayerCount int         `json:"playerCount"`
	CreatedAt   int64       `json:"createdAt"`
	LastSeen    int64       `json:"lastSeen"`
	Config      LobbyConfig `json:"config"`
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
			Config:       config,
			Running:      true,
			LastUpdate:   time.Now(),
			Player:       make([]Player, config.PlayerLimit),
			Clients:      make(map[uint8]*clientData),
			BorderRadius: float64(config.BorderRadius),
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
			Config:      config,
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
	p.X = randomFloat(float64(-l.BorderRadius), float64(l.BorderRadius))
	p.Y = randomFloat(float64(-l.BorderRadius), float64(l.BorderRadius))
	p.SpeedX = 0
	p.SpeedY = 0
	p.Radius = 5
	p.FlailX = p.X
	p.FlailY = p.Y
	p.FlailSpeedX = 0
	p.FlailSpeedY = 0
	p.FlailConnected = true
	p.FlailRadius = 2.5
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

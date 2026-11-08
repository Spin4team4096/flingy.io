package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// Main Loop
func main() {
	log.Println("Starting Server...")
	godotenv.Load()

	// initialise FireBase
	ctx = context.Background()
	conf := &firebase.Config{
		DatabaseURL: os.Getenv("FIREBASE_DATABASE_URL"),
	}
	// Fetch the service account key JSON file contents
	opt := option.WithAuthCredentialsFile(option.ServiceAccount, os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err = app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	log.Println("Initialised firebase")

	// Create lobbies
	data, err := os.ReadFile("lobbies.json")
	if err != nil {
		log.Fatalln("Error reading lobbies.json:", err)
	}
	var configs []LobbyConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		log.Fatalln("Error parsing lobbies.json:", err)
	}
	serversRef := client.NewRef("flingy2/lobbies")
	var servers map[string]ServerEntry
	if err := serversRef.Get(ctx, &servers); err != nil {
		log.Println("Error reading servers:", err)
	}
	if servers == nil {
		servers = map[string]ServerEntry{}
	}
	for i := range configs {
		n := 1
		for {
			candidate := fmt.Sprintf("%s%04d", configs[i].ID, n)
			if _, ok := servers[candidate]; !ok {
				break
			}
			n++
		}
		configs[i].ID = fmt.Sprintf("%s%04d", configs[i].ID, n)
		servers[configs[i].ID] = ServerEntry{}
		createLobby(configs[i])
	}
	log.Println("Created lobbies")
	startWebsocketServer()
	startConsole()
	go firebaseLoop()
	log.Println("Started server")

	// Main tick loop
	ticker := time.NewTicker(time.Second / 30)
	last := time.Now()
	for range ticker.C {
		now := time.Now()
		dt := float32(now.Sub(last).Seconds() * 30)
		last = now
		if dt > 1.5 {
			dt = 1.5
		}
		lobbiesMu.Lock()
		for _, lobby := range lobbies {
			lobby.Tick++
			lobby.StateTick++

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

			switch lobby.State {
			case 0:
				if playerCount > 1 {
					lobby.State = 1
					lobby.StateTick = 0
					log.Printf("Starting Match")
				}
			case 1:
				if playerCount <= 1 {
					lobby.State = 0
					lobby.StateTick = 0
					log.Printf("Not enough players anymore")
					// reset lobby data too, idk how
				} else if lobby.StateTick >= lobby.Config.CountdownTime*30 {
					lobby.State = 2
					lobby.StateTick = 0
					// shrink players
					// broadcast game started
					log.Printf("Match Started")
					// spawn game started sound (distanceless)
				} else if playersAlive > 0 {
					n := 0
					for i := range lobby.Config.PlayerLimit {
						if lobby.Player[i].HP > 0 {
							n++
							angle := 2 * math.Pi * (float64(n) / float64(playersAlive))
							tx := float32(math.Cos(angle) * float64(lobby.Config.BorderSize) * 0.8)
							ty := float32(math.Sin(angle) * float64(lobby.Config.BorderSize) * 0.8)

							lobby.Player[i].X += (tx - lobby.Player[i].X) * 0.05
							lobby.Player[i].Y += (ty - lobby.Player[i].Y) * 0.05
						}
					}
				}
			case 2:
				if playerCount == 0 {
					lobby.State = 0
					lobby.StateTick = 0
					log.Printf("Everyone left, BOOO")
					// reset lobby data too, idk how
				} else if playersAlive <= 1 {
					lobby.State = 3
					lobby.StateTick = 0
					// broadcast winner
					log.Printf("Somebody won, idk")
					// spawn game ended sound (distanceless)
				}
			case 3:
				if lobby.StateTick >= 30*5 || playerCount == 0 {
					lobby.State = 0
					lobby.StateTick = 0
					log.Printf("Match Ended")

					for _, c := range lobby.Clients {
						c.conn.CloseNow()
					}
					lobby.Clients = make(map[uint8]*clientData)

					for i := range lobby.Player {
						despawnPlayer(lobby, i)
					}
				}
			}

			// player updates
			for i := range lobby.Player {
				p := &lobby.Player[i]

				if p.Alive {
					if p.Immune {
						// fuck, i forgot to add immunity timer in struct. will do later
					}

				}
			}

		}
		lobbiesMu.Unlock()
	}
}

var (
	ctx    context.Context
	client *db.Client
)

// Helper Functions
func randomInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func randomFloat(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

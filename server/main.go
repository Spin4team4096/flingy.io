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
			candidate := fmt.Sprintf("%s_%04d", configs[i].ID, n)
			if _, ok := servers[candidate]; !ok {
				break
			}
			n++
		}
		configs[i].ID = fmt.Sprintf("%s_%04d", configs[i].ID, n)
		servers[configs[i].ID] = ServerEntry{}
		createLobby(configs[i])
	}
	log.Println("Created lobbies")
	startWebsocketServer()
	startConsole()
	go firebaseLoop()
	log.Println("Started server")

	// Main tick loop
	targetTPS := 30
	ticker := time.NewTicker(time.Second / time.Duration(targetTPS))
	last := time.Now()
	for range ticker.C {
		now := time.Now()
		dt := now.Sub(last).Seconds() * 30
		last = now
		if dt > 1.5 {
			dt = 1.5
		}
		lobbiesMu.Lock()
		for id, lobby := range lobbies {

			lobby.Tick++
			lobby.StateTick++

			playerCount := 0
			playersAlive := 0
			playersUnharmed := 0
			teamCount := 0
			teamList := map[int]int{}
			teamColors := [4]int{0, 8, 2, 4}
			for i := range lobby.Player {
				p := &lobby.Player[i]
				if p.Active {
					if lobby.Config.Gamemode == "TDM" {
						if lobby.State <= 1 {
							p.Team = playerCount % lobby.Config.TeamCount
						}
						p.ColorIndex = teamColors[p.Team]
					}

					playerCount++
				}
				if p.HP > 0 {
					group := p.Team
					if group == -1 {
						group = i
					}
					if _, ok := teamList[group]; !ok {
						teamCount++
						teamList[group] = teamCount
					}

					playersAlive++
				}
				if p.HP >= 3 && p.Active {
					playersUnharmed++
				}
			}

			switch lobby.State {
			case 0:
				if lobby.Closing {
					for _, c := range lobby.Clients {
						c.conn.CloseNow()
					}

					delete(lobbies, id)
					go unregisterLobby(id)
					continue
				}
				lobby.Open = true
				if playerCount >= lobby.Config.MinPlayers {
					lobby.State = 1
					lobby.StateTick = 0
					log.Printf("[%s] Preparing match...", id)

					packet := encodePacket(5, lobby) // send all usernames
					for _, c := range lobby.Clients {
						c.send(1, packet)
					}
				}
				lobby.BorderRadius = float64(lobby.Config.BorderRadius)
			case 1:
				lobby.Open = playerCount < lobby.Config.PlayerLimit
				if playerCount < lobby.Config.MinPlayers {
					lobby.State = 0
					lobby.StateTick = 0
					log.Printf("[%s] Not enough players anymore", id)

					packet := encodePacket(5, lobby)
					for _, c := range lobby.Clients {
						c.send(1, packet)
					}

					// reset lobby data too, idk how
				} else if lobby.StateTick >= lobby.Config.CountdownTime*targetTPS {
					lobby.State = 2
					lobby.StateTick = 0
					for i := range lobby.Config.PlayerLimit {
						p := &lobby.Player[i]
						p.FlailRadius = 2.5
						p.Immune = true
						p.TimeLeftImmune = 60
					}
					// broadcast game started
					log.Printf("[%s] Match started!", id)

					packet := encodePacket(5, lobby)
					for _, c := range lobby.Clients {
						c.send(1, packet)
					}

					// spawn game started sound (distanceless)
				} else if playersAlive > 0 {
					// countdown pull
					for i := range lobby.Config.PlayerLimit {
						p := &lobby.Player[i]
						if lobby.Player[i].HP > 0 {
							group := p.Team
							if group == -1 {
								group = i
							}
							angle := 2 * math.Pi * float64(teamList[group]) / float64(teamCount)
							tx := math.Cos(angle) * float64(lobby.BorderRadius) * 0.8
							ty := math.Sin(angle) * float64(lobby.BorderRadius) * 0.8
							k := 1 - math.Pow(0.93, dt)
							p.X += (tx - p.X) * k
							p.Y += (ty - p.Y) * k
						}
					}
				}
			case 2:
				if playerCount == playersUnharmed && playerCount < lobby.Config.PlayerLimit && lobby.Config.Gamemode == "FFA" {
					lobby.Open = true
				} else {
					lobby.Open = false
				}
				if playerCount == 0 {
					lobby.State = 0
					lobby.StateTick = 0
					log.Printf("[%s] Everyone left early, BOOO", id)
					// reset lobby data too, idk how
				} else if teamCount <= 1 {
					lobby.State = 3
					lobby.StateTick = 0
					// broadcast winner
					log.Printf("[%s] Somebody won, idk who", id)
					// spawn game ended sound (distanceless)
				}
				lobby.BorderRadius = clamp(lobby.BorderRadius-(float64(lobby.Config.BorderRadius)-float64(lobby.Config.MinBorderRadius))/(float64(lobby.Config.MatchTime)*30)*dt, float64(lobby.Config.MinBorderRadius), float64(lobby.Config.BorderRadius))
			case 3:
				lobby.Open = false
				if lobby.StateTick >= targetTPS*5 || playerCount == 0 {
					lobby.State = 0
					lobby.StateTick = 0
					log.Printf("[%s] Match ended", id)

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
						p.TimeLeftImmune -= dt
						if p.TimeLeftImmune <= 0 {
							p.Immune = false
						}
					}

					if lobby.BorderRadius <= float64(lobby.Config.MinBorderRadius) {
						p.HP -= (dt / 30) / (10 + p.FlailRadius/5)
						if p.HP <= 0 {
							p.HP = 0
							p.Alive = false

							playersAlive := 0
							for i := range lobby.Player {
								if lobby.Player[i].HP > 0 {
									playersAlive++
								}
							}
							p.Rank = playersAlive
						}
					}

					// direction

					raw := directionTo(p.WorldMouseX, p.WorldMouseY, p.X, p.Y)
					p.SpinTarget += math.Remainder(raw-p.LastTarget, 360)
					p.LastTarget = raw

					if p.FlailConnected {
						diff := p.SpinTarget - p.Direction
						maxSpeed := float64(128)
						if diff > 360+maxSpeed {
							p.SpinTarget -= 360 * math.Ceil((diff-(360+maxSpeed))/360)
						} else if diff < -360-maxSpeed {
							p.SpinTarget += 360 * math.Ceil((-diff-(360+maxSpeed))/360)
						}
						diff = p.SpinTarget - p.Direction

						// change 2.0 to make it snappier
						gain := 0.03 * (1 - math.Pow(0.9, 2.0)) / (1 - 0.9)
						damping := math.Pow(0.9, 2.0*dt)
						brake := 1 + math.Sqrt(p.FlailRadius*2)/10
						accel := clamp(diff, -maxSpeed, maxSpeed) * gain / brake

						p.AngularVelocity = p.AngularVelocity*damping + accel*dt
						p.Direction += p.AngularVelocity * dt
					} else {
						p.AngularVelocity = 0
						p.Direction = directionTo(p.WorldMouseX, p.WorldMouseY, p.X, p.Y)
						p.SpinTarget, p.LastTarget = raw, raw
					}

					sin, cos := math.Sincos(degToRad(p.Direction))

					// fling
					if p.FlailConnected && p.LeftMouseDown && !p.FlingUsed && (lobby.State == 2 || lobby.State == 0) {
						p.FlailConnected = false
						p.SpeedX += -cos * (7.5 + p.FlailRadius*0.15)
						p.SpeedY += -sin * (7.5 + p.FlailRadius*0.15)
						p.FlingUsed = true
						p.Immune = false
					}
					if !p.LeftMouseDown {
						p.FlingUsed = false
					}

					if p.FlailConnected {
						p.FlailRadius = clamp(p.FlailRadius+((math.Abs(p.AngularVelocity)/100-0.05)*dt), 2.5, 50)
					} else {
						if p.FlailRadius > 2.5 {
							p.FlailRadius = clamp(p.FlailRadius-(0.15*dt), 0, 50)
						} else {
							p.FlailRadius = clamp(p.FlailRadius-(0.05*dt), 0, 50)
						}

						if p.FlailRadius <= 0 || p.FlailSpeedX*p.FlailSpeedX+p.FlailSpeedY*p.FlailSpeedY < 0.25 {
							p.FlailConnected = true
							p.FlingUsed = true
						}
					}

					// movement
					if (p.WorldMouseX-p.X)*(p.WorldMouseX-p.X)+(p.WorldMouseY-p.Y)*(p.WorldMouseY-p.Y) > 25*25 {
						p.SpeedX += cos * 0.3 * dt
						p.SpeedY += sin * 0.3 * dt
					}

					p.SpeedX *= math.Pow(0.92, dt)
					p.SpeedY *= math.Pow(0.92, dt)
					p.X += p.SpeedX * dt
					p.Y += p.SpeedY * dt

					r := math.Hypot(p.X, p.Y)
					if max := float64(lobby.BorderRadius) - p.Radius; r > max {
						p.X *= max / r
						p.Y *= max / r
					}

					// flail
					if p.FlailConnected {
						oldX, oldY := p.FlailX, p.FlailY
						p.FlailX = p.X + cos*(16+p.Radius+p.FlailRadius)
						p.FlailY = p.Y + sin*(16+p.Radius+p.FlailRadius)
						p.FlailSpeedX = (p.FlailX - oldX) / dt
						p.FlailSpeedY = (p.FlailY - oldY) / dt
					} else {
						p.FlailSpeedX *= math.Pow(0.96, dt)
						p.FlailSpeedY *= math.Pow(0.96, dt)
						p.FlailX += p.FlailSpeedX * dt
						p.FlailY += p.FlailSpeedY * dt
					}

					// collision
					for n := range lobby.Player {
						c := &lobby.Player[n]
						if (p.FlailX-c.X)*(p.FlailX-c.X)+(p.FlailY-c.Y)*(p.FlailY-c.Y) <= (p.FlailRadius+c.Radius)*(p.FlailRadius+c.Radius) && c.Alive {
							if p == c {
								p.FlailConnected = true
								p.FlingUsed = true
							} else if lobby.State == 2 && !p.Immune && !c.Immune && (p.Team != c.Team || p.Team == -1 || c.Team == -1) {
								c.HP -= 1
								c.Alive = false
								c.TimeUntilRespawn = 90
								// spawn death sound
								if killer := lobby.Clients[uint8(i)]; killer != nil && lobby.Config.AddMoney {
									killer.send(2, encodePacket(6, uint8(1)))
								}
								if c.HP <= 0 {
									// broadcast elimination to killer
									// spawn lose sound for losing player only
									c.HP = 0
									playersAlive := 0
									for i := range lobby.Player {
										if lobby.Player[i].HP > 0 {
											playersAlive++
										}
									}
									c.Rank = playersAlive
								} else {
									// broadcast kill to killer
									c.X = randomFloat(-lobby.BorderRadius, lobby.BorderRadius)
									c.Y = randomFloat(-lobby.BorderRadius, lobby.BorderRadius)
									c.SpeedX = 0
									c.SpeedY = 0
									c.FlailX = c.X
									c.FlailY = c.Y
								}
							}
						}
					}
				} else if p.HP > 0 {
					p.TimeUntilRespawn -= dt
					if p.TimeUntilRespawn <= 0 {
						p.Alive = true
						p.FlailRadius *= 0.64
						p.FlailConnected = true
						p.FlingUsed = true
						p.Immune = true
						p.TimeLeftImmune = 60
					}
				} else {
					// out for good, can't do anything
				}
			}

			if len(lobby.Clients) > 0 {
				packet := encodePacket(3, lobby)
				for _, c := range lobby.Clients {
					c.send(3, packet)
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

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func radToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func directionTo(nx, ny, ox, oy float64) float64 {
	return radToDeg(math.Atan2(ny-oy, nx-ox))
}

func clamp(v, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, v))
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

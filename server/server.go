package server

import (
	"encoding/json"
	"log"
	"net"
	"playground/raylib-go/models"
	"time"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450
const SCORE_LIMIT int32 = 3

// var mu sync.Mutex
type InputMessage struct {
	Addr string
	Cmd  string // e.g., "R_J", "R_K"
}

type Coordinate struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Ball struct {
	Position Coordinate `json:"position"`
	Speed    Coordinate `json:"speed"`
	Radius   float32    `json:"radius"`
	IsActive bool       `json:"is_active"`
}

type Player struct {
	Position Coordinate `json:"position"`
	Size     Coordinate `json:"size"`
	Score    int32      `json:"score"`
}

type Game struct {
	Red    Player              `json:"red"`
	Blue   Player              `json:"blue"`
	Ball   Ball                `json:"ball"`
	Conn   map[string]net.Addr `json:"conn"`
	Client *models.Client      `json:"client"`
}

var g *Game = &Game{
	Red: Player{
		Position: Coordinate{
			X: 0,
			Y: 200,
		},
		Size: Coordinate{
			X: 10,
			Y: 100,
		},
	},

	Blue: Player{
		Position: Coordinate{
			X: float32(SCREEN_WIDTH) - 10,
			Y: 200,
		},
		Size: Coordinate{
			X: 10,
			Y: 100,
		},
	},
	Ball: Ball{
		Position: Coordinate{
			X: float32(SCREEN_WIDTH) / 2,
			Y: float32(SCREEN_HEIGHT) / 2,
		},
		Speed: Coordinate{
			X: 3.0,
			Y: 0.0,
		},
		IsActive: false,
	},
	Conn: make(map[string]net.Addr),
}

func StartServer() {
	conn, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		log.Fatalf("unable to start the server => Error:%v ", err.Error())
	}
	inputChan := make(chan InputMessage, 100)
	// defer conn.Close()
	log.Printf("listening connection on port :8000")
	// server listener routine that gonna listen from client
	go func() {
		for {
			var buf [10]byte

			n, addr, err := conn.ReadFrom(buf[:])
			if err != nil {
				log.Printf("error when setup: %v", err.Error())
			}
			cmd := string(buf[:n])
			// log.Printf("i receive and transfer: %v\n", cmd)
			inputChan <- InputMessage{
				Cmd:  cmd,
				Addr: addr.String(),
			}

		}

	}()

	// TODO: ball routine

	go func() {
		ticker := time.NewTicker(16 * time.Millisecond)
		defer ticker.Stop()
		for {

			select {
			case msg := <-inputChan:
				{

					log.Printf("received: %v \n\n", string(msg.Cmd))
					if _, exists := g.Conn[msg.Addr]; !exists && len(g.Conn) < 2 {
						udpAddr, err := net.ResolveUDPAddr("udp", msg.Addr)
						if err == nil {
							g.Conn[msg.Addr] = udpAddr
						}
					}
					if err != nil {
						log.Printf("unable to read: Error: %v\n", err.Error())
					}

					// TODO: only read sent value from the valid clients

					// log.Printf("received: %v \n\n", string(msg.Cmd))

					switch msg.Cmd {
					case "R_J":
						{

							if g.Red.Position.Y < float32(SCREEN_HEIGHT)-g.Red.Size.Y/2 {
								g.Red.Position.Y = g.Red.Position.Y + 2
							}

						}
					case "R_K":
						{
							if g.Red.Position.Y > g.Red.Size.Y/2 {
								g.Red.Position.Y = g.Red.Position.Y - 2
							}

						}
					case "R_H":
						{
							if g.Red.Position.X > 0 {
								g.Red.Position.X = g.Red.Position.X - 2
							}

						}
					case "R_L":
						{

							if g.Red.Position.X < float32(SCREEN_WIDTH)-g.Red.Size.X {
								g.Red.Position.X = g.Red.Position.X + 2
							}
						}
					case "B_H":
						{

							if g.Blue.Position.X > 0 {
								g.Blue.Position.X = g.Blue.Position.X - 2
							}
						}
					case "B_L":
						{
							if g.Blue.Position.X < float32(SCREEN_WIDTH)-g.Blue.Size.X {
								g.Blue.Position.X = g.Blue.Position.X + 2
							}
						}
					case "B_J":
						{
							if g.Blue.Position.Y < float32(SCREEN_HEIGHT)-g.Blue.Size.Y/2 {
								g.Blue.Position.Y = g.Blue.Position.Y + 2
							}
						}
					case "B_K":
						{
							if g.Blue.Position.Y > g.Blue.Size.Y/2 {
								g.Blue.Position.Y = g.Blue.Position.Y - 2
							}
						}
					// collisions
					case "R_B":
						{

							g.Ball.Speed.X = 3.0
							g.Ball.Speed.Y = (g.Ball.Position.Y - g.Red.Position.Y) / (g.Red.Size.Y / 2) * 5

						}
					case "B_B":
						{
							g.Ball.Speed.X = -3.0
							g.Ball.Speed.Y = (g.Ball.Position.Y - g.Blue.Position.Y) / (g.Blue.Size.Y / 2) * 5

						}
					case "START":
						{

							g.Ball.IsActive = true
						}

					default:

					}

				}
			case <-ticker.C:
				{
					if g.Ball.IsActive {
						g.Ball.Position.X = g.Ball.Position.X + float32(g.Ball.Speed.X)
						g.Ball.Position.Y = g.Ball.Position.Y + float32(g.Ball.Speed.Y)
					} else {
						g.Ball.Position = Coordinate{X: g.Red.Position.X + 2*g.Red.Size.X + g.Ball.Radius, Y: g.Red.Position.Y}
					}
					// ball.Position.Y = ball.Position.Y + float32(ball.Speed.Y)
					if g.Ball.Position.X <= 0 {
						g.Blue.Score = g.Blue.Score + 1
						// reset
						g.Reset()
					}
					if g.Ball.Position.X >= float32(SCREEN_WIDTH) {
						g.Red.Score = g.Red.Score + 1
						// reset
						g.Reset()
					}

					if g.Ball.Position.Y <= 0 {
						g.Ball.Speed.Y = 3.0
					}
					if g.Ball.Position.Y >= float32(SCREEN_HEIGHT) {
						g.Ball.Speed.Y = -3.0
					}

					msg, err := json.Marshal(g)
					if err != nil {
						log.Println("unable to marshal the message")
						log.Printf("error: %v \n", err.Error())
					}
					log.Printf("new state |> %v \n", string(msg))

					for _, c := range g.Conn {
						if c != nil {
							log.Printf("number ofcurrent client: %v\n", len(g.Conn))
							log.Printf("client: %v", c)
							conn.WriteTo(msg, c)
						}
					}

				}
			}
		}

	}()
}

func sendResponse(conn net.PacketConn, msg []byte) {
	for _, c := range g.Conn {
		if c != nil {
			log.Printf("number ofcurrent client: %v\n", len(g.Conn))
			conn.WriteTo(msg, c)
		}
	}
}

func (g *Game) Reset() {
	if g.Red.Score >= SCORE_LIMIT {
		// rl.BeginDrawing()
		// rl.ClearBackground(rl.Black)
		// rl.DrawText("RED WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Red)
		// rl.EndDrawing()
		// rl.WaitTime(2)
		// os.Exit(0)
		// TODO: send red wins message
	}

	if g.Blue.Score >= SCORE_LIMIT {
		// rl.BeginDrawing()
		// rl.ClearBackground(rl.Black)
		// rl.DrawText("BLUE WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Blue)
		// rl.EndDrawing()
		// rl.WaitTime(2)
		// os.Exit(0)
		// TODO: send blue wins message

	}

	// TODO: put this on the client

	// rl.BeginDrawing()
	// rl.ClearBackground(rl.Black)
	// rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", g.Red.Score, g.Blue.Score), int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 24, rl.White)
	// rl.EndDrawing()
	// rl.WaitTime(2.5)

	// INFO: reset blue
	g.Blue.Position.X = float32(SCREEN_WIDTH) - 10
	g.Blue.Position.Y = 200
	g.Blue.Size.X = 10
	g.Blue.Size.Y = 100

	// INFO: reset red
	g.Red.Position.X = 0
	g.Red.Position.Y = 200
	g.Red.Size.X = 10
	g.Red.Size.Y = 100

	// INFO: reset ball
	g.Ball.Position = Coordinate{X: g.Red.Position.X + 2*g.Red.Size.X + g.Ball.Radius, Y: g.Red.Position.Y}
	g.Ball.Speed.Y = 0.0
	g.Ball.Speed.X = 3.0
	g.Ball.IsActive = false
}

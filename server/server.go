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

var g *models.ServerGame = &models.ServerGame{
	Red: models.ServerPlayer{
		Position: models.Coordinate{
			X: 0,
			Y: 200,
		},
		Size: models.Coordinate{
			X: 10,
			Y: 100,
		},
	},

	Blue: models.ServerPlayer{
		Position: models.Coordinate{
			X: float32(SCREEN_WIDTH) - 10,
			Y: 200,
		},
		Size: models.Coordinate{
			X: 10,
			Y: 100,
		},
	},
	Ball: models.ServerBall{
		Position: models.Coordinate{
			X: float32(SCREEN_WIDTH) / 2,
			Y: float32(SCREEN_HEIGHT) / 2,
		},
		Speed: models.Coordinate{
			X: 3.0,
			Y: 0.0,
		},
		IsActive: false,
	},
	Conn:         make(map[string]net.Addr),
	Winner:       make(map[string]bool),
	ScoreUpdated: true,
}

func StartServer() {

	// INFO: trick to get local ip by making udp request for google dns
	connection, err := net.Dial("udp", "8.8.8.8:80")
	defer connection.Close()
	localAddr := connection.LocalAddr().(*net.UDPAddr).IP
	conn, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		log.Fatalf("unable to start the server => Error:%v ", err.Error())
	}
	inputChan := make(chan models.InputMessage, 100)

	log.Printf("listening connection on port %v:8000", localAddr)
	// INFO: server listener routine that gonna listen from client
	go func() {
		for {
			var buf [10]byte

			n, addr, err := conn.ReadFrom(buf[:])
			if err != nil {
				log.Printf("error when setup: %v", err.Error())
			}
			cmd := string(buf[:n])
			// log.Printf("i receive and transfer: %v\n", cmd)
			inputChan <- models.InputMessage{
				Cmd:  cmd,
				Addr: addr.String(),
			}

		}

	}()

	// INFO: game state modifer and sender routine

	go func() {
		ticker := time.NewTicker(16 * time.Millisecond)
		defer ticker.Stop()
		for {

			select {
			case msg := <-inputChan:
				{

					if _, exists := g.Conn[msg.Addr]; !exists && len(g.Conn) < 2 {
						udpAddr, err := net.ResolveUDPAddr("udp", msg.Addr)
						if err == nil {
							g.Conn[msg.Addr] = udpAddr
						}
					}
					if err != nil {
						log.Printf("unable to read: Error: %v\n", err.Error())
					}

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
					case "START_R":
						{
							if !g.Ball.IsActive {
								g.Ball.IsActive = true
								g.Ball.Speed.X = 3.0

							}

						}
					case "START_B":
						{
							if !g.Ball.IsActive {
								g.Ball.Speed.X = -3.0
								g.Ball.IsActive = true

							}
						}
					default:

					}

				}
			case <-ticker.C:
				{
					g.ScoreUpdated = false
					if g.Ball.IsActive {
						g.Ball.Position.X = g.Ball.Position.X + float32(g.Ball.Speed.X)
						g.Ball.Position.Y = g.Ball.Position.Y + float32(g.Ball.Speed.Y)
					} else {
						g.Ball.Position = models.Coordinate{X: g.Red.Position.X + 2*g.Red.Size.X + g.Ball.Radius, Y: g.Red.Position.Y}
					}
					// ball.Position.Y = ball.Position.Y + float32(ball.Speed.Y)
					if g.Ball.Position.X <= 0 {
						g.Blue.Score = g.Blue.Score + 1
						// TODO: reset
						reset()
					}
					if g.Ball.Position.X >= float32(SCREEN_WIDTH) {
						g.Red.Score = g.Red.Score + 1
						// TODO: reset
						reset()
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
					// log.Printf("new state |> %v \n", string(msg))

					for _, addr := range g.Conn {
						// log.Printf("client: %v", addr)
						sendResponse(conn, addr, msg)
					}

				}
			}
		}

	}()
}

func sendResponse(conn net.PacketConn, addr net.Addr, msg []byte) {
	if addr != nil {
		conn.WriteTo(msg, addr)
	}
}

func reset() {
	if g.Red.Score >= SCORE_LIMIT {
		g.Winner["red"] = true
	}
	if g.Blue.Score >= SCORE_LIMIT {
		g.Winner["blue"] = true
	}

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
	g.Ball.Position = models.Coordinate{X: g.Red.Position.X + 2*g.Red.Size.X + g.Ball.Radius, Y: g.Red.Position.Y}
	g.Ball.Speed.Y = 0.0
	g.Ball.Speed.X = 3.0
	g.Ball.IsActive = false

	// INFO: update score
	g.ScoreUpdated = true
}

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"playground/raylib-go/client"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450
const SCORE_LIMIT int32 = 3

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
	Client *client.Client      `json:"client"`
}

var GAME *Game = &Game{
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

func NewGame() *Game {
	return GAME
}

func StartServer() {
	conn, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		log.Fatalf("unable to start the server => Error:%v ", err.Error())
	}
	log.Printf("listening connection on port :8000")
	go func() {
		for {
			var buf [10]byte

			_, addr, err := conn.ReadFrom(buf[0:])

			if len(GAME.Conn) < 2 {
				GAME.Conn[addr.String()] = addr
			}

			if err != nil {
				log.Printf("unable to read: Error: %v\n", err.Error())
			}

			// TODO: only read sent value from the valid clients

			var msg []byte
			switch cmd := string(buf[0:6]); cmd {
			case "RKEY_J":
				{

					if GAME.Red.Position.Y < float32(SCREEN_HEIGHT)-GAME.Red.Size.Y/2 {
						GAME.Red.Position.Y = GAME.Red.Position.Y + 2
					}
					log.Printf("position: %v\n", GAME.Red.Position.Y)

					log.Printf("message: |> %v \n", string(msg))

					msg, err = json.Marshal(GAME.Red.Position)
					if err != nil {
						log.Println("unable to marshal the message")
					}
				}
			default:

			}
			// if string(buf[0:6]) == "RKEY_J" {
			// 	// sendResponse(conn, msg)
			// }

			for _, c := range GAME.Conn {
				if c != nil {
					log.Printf("number ofcurrent client: %v\n", len(GAME.Conn))
					conn.WriteTo(msg, c)
				}
			}

		}

	}()
}

func sendResponse(conn net.PacketConn, msg []byte) {
	for _, c := range GAME.Conn {
		if c != nil {
			log.Printf("number ofcurrent client: %v\n", len(GAME.Conn))
			conn.WriteTo(msg, c)
		}
	}
}

func (g *Game) Reset() {
	if g.Red.Score >= SCORE_LIMIT {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("RED WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Red)
		rl.EndDrawing()
		rl.WaitTime(2)
		os.Exit(0)
	}

	if g.Blue.Score >= SCORE_LIMIT {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("BLUE WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Blue)
		rl.EndDrawing()
		rl.WaitTime(2)
		os.Exit(0)
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", g.Red.Score, g.Blue.Score), int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 24, rl.White)
	rl.EndDrawing()
	rl.WaitTime(2.5)

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

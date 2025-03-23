package server

import (
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

type Ball struct {
	Position rl.Vector2
	Speed    rl.Vector2
	Radius   float32
	IsActive bool
}

type Player struct {
	Position rl.Vector2
	Size     rl.Vector2
	Score    int32
}

type Game struct {
	Red    Player
	Blue   Player
	Ball   Ball
	Conn   map[string]net.Addr
	Client *client.Client
}

var GAME *Game = &Game{
	Red: Player{
		Position: rl.Vector2{
			X: 0,
			Y: 200,
		},
		Size: rl.Vector2{
			X: 10,
			Y: 100,
		},
	},

	Blue: Player{
		Position: rl.Vector2{
			X: float32(SCREEN_WIDTH) - 10,
			Y: 200,
		},
		Size: rl.Vector2{
			X: 10,
			Y: 100,
		},
	},
	Ball: Ball{
		Position: rl.Vector2{
			X: float32(SCREEN_WIDTH) / 2,
			Y: float32(SCREEN_HEIGHT) / 2,
		},
		Speed: rl.Vector2{
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
				log.Printf("unable to read: Error: %v", err.Error())
				return
			}

			log.Printf("|> %v", string(buf[0:]))
			if string(buf[0:]) == "RKEY_J" {
				if GAME.Red.Position.Y < float32(SCREEN_HEIGHT)-GAME.Red.Size.Y/2 {
					GAME.Red.Position.Y = GAME.Red.Position.Y + 2
				}

			}
			for _, c := range GAME.Conn {
				if c != nil {
					log.Printf("number ofcurrent client: %v\n", len(GAME.Conn))
					conn.WriteTo([]byte("dont blame me i am sending something"), c)
				}
			}

		}

	}()
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
	g.Ball.Position = rl.Vector2{X: g.Red.Position.X + 2*g.Red.Size.X + g.Ball.Radius, Y: g.Red.Position.Y}
	g.Ball.Speed.Y = 0.0
	g.Ball.Speed.X = 3.0
	g.Ball.IsActive = false
}

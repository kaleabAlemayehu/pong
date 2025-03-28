package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"playground/raylib-go/client"
	"playground/raylib-go/server"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450
const SCORE_LIMIT int32 = 3

type Ball struct {
	Position rl.Vector2 `json:"position"`
	Speed    rl.Vector2 `json:"speed"`
	Radius   float32    `json:"radius"`
	IsActive bool       `json:"is_active"`
}

type Player struct {
	Position rl.Vector2 `json:"position"`
	Size     rl.Vector2 `json:"size"`
	Score    int32      `json:"score"`
}

type Game struct {
	Red    Player         `json:"red"`
	Blue   Player         `json:"blue"`
	Ball   Ball           `json:"ball"`
	Client *client.Client `json:"client"`
}

var g *Game = &Game{
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
}

func main() {
	host := flag.Bool("host", false, "to host server of the game")
	flag.Parse()
	message := make(chan []byte, 10000)
	input := make(chan string, 10000)

	if *host == true {
		server.StartServer()
	}

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "pong")
	defer rl.CloseWindow()

	client.ListeningClient(input, message)

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {

		select {
		case msg, ok := <-message:
			if !ok {
				log.Printf("there is something wrong with the recieved data\n")

			}
			err := json.Unmarshal(msg, g)
			if err != nil {
				log.Printf("error happend when recieved data from the server unmarshaled\n")
			}

		default:
			handleMovement(input)
			if rl.CheckCollisionCircleRec(rl.Vector2(g.Ball.Position), g.Ball.Radius, rl.Rectangle{X: g.Red.Position.X, Y: g.Red.Position.Y - g.Red.Size.Y/2, Width: g.Red.Size.X, Height: g.Red.Size.Y}) {
				g.Ball.Speed.X = 3.0
				g.Ball.Speed.Y = (g.Ball.Position.Y - g.Red.Position.Y) / (g.Red.Size.Y / 2) * 5
				// send red collition with ball

			}
			if rl.CheckCollisionCircleRec(rl.Vector2(g.Ball.Position), 10.0, rl.Rectangle{X: g.Blue.Position.X, Y: g.Blue.Position.Y - g.Blue.Size.Y/2, Width: g.Blue.Size.X, Height: g.Blue.Size.Y}) {
				g.Ball.Speed.X = -3.0
				g.Ball.Speed.Y = (g.Ball.Position.Y - g.Blue.Position.Y) / (g.Blue.Size.Y / 2) * 5
				// send blue collition with ball
			}

			// INFO: drawing start
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)
			rl.DrawRectangleRec(rl.Rectangle{X: g.Red.Position.X, Y: g.Red.Position.Y - (g.Red.Size.Y / 2), Width: g.Red.Size.X, Height: g.Red.Size.Y}, rl.Red)
			rl.DrawRectangleRec(rl.Rectangle{X: g.Blue.Position.X, Y: g.Blue.Position.Y - (g.Blue.Size.Y / 2), Width: g.Blue.Size.X, Height: g.Blue.Size.Y}, rl.Blue)
			rl.DrawCircleV(rl.Vector2(g.Ball.Position), 10.0, rl.White)
			rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", g.Red.Score, g.Blue.Score), int32(rl.GetScreenWidth()/2)-115, 0, 20, rl.White)
			rl.EndDrawing()
		}
	}
}

func handleMovement(input chan string) {
	if rl.IsKeyDown(rl.KeyJ) {
		input <- "R_J"
		// _, err := g.Client.Conn.Write([]byte("R_J"))
		// if err != nil {
		// 	log.Println("pressing j is not sending data")
		// }

		// var res []byte
		//
		// n, _, err := g.Client.Conn.ReadFromUDP(res[:])
		// if err != nil {
		// 	log.Println("unable to get response for the server")
		// }
		// log.Printf("response from the server RKEY_J: %v", string(res[0:n]))
		//
		// log.Printf("value: %v\n", string(res[:n]))
		// var resval rl.Vector2
		// err = json.Unmarshal(res[:n], &resval)
		// if err != nil {
		// 	log.Println("error form client:", err.Error())
		// }
		//
		// g.Red.Position = resval
	}
	if rl.IsKeyDown(rl.KeyK) {
		if g.Red.Position.Y > g.Red.Size.Y/2 {
			g.Red.Position.Y = g.Red.Position.Y - 2
		}
	}
	if rl.IsKeyDown(rl.KeyH) {
		if g.Red.Position.X > 0 {
			g.Red.Position.X = g.Red.Position.X - 2
		}
	}
	if rl.IsKeyDown(rl.KeyL) {
		if g.Red.Position.X < float32(SCREEN_WIDTH)-g.Red.Size.X {
			g.Red.Position.X = g.Red.Position.X + 2
		}
	}

	if rl.IsKeyDown(rl.KeyA) {
		// B_H
		if g.Blue.Position.X > 0 {
			g.Blue.Position.X = g.Blue.Position.X - 2
		}
	}

	if rl.IsKeyDown(rl.KeyF) {
		// B_L
		if g.Blue.Position.X < float32(SCREEN_WIDTH)-g.Blue.Size.X {
			g.Blue.Position.X = g.Blue.Position.X + 2
		}
	}

	if rl.IsKeyDown(rl.KeyS) {
		if g.Blue.Position.Y < float32(SCREEN_HEIGHT)-g.Blue.Size.Y/2 {
			g.Blue.Position.Y = g.Blue.Position.Y + 2
		}
	}
	if rl.IsKeyDown(rl.KeyD) {
		if g.Blue.Position.Y > g.Blue.Size.Y/2 {
			g.Blue.Position.Y = g.Blue.Position.Y - 2
		}
	}
	if rl.IsKeyDown(rl.KeySpace) {
		g.Ball.IsActive = true
	}
}

func Reset() {
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

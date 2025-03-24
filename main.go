package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"playground/raylib-go/client"
	"playground/raylib-go/server"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450
const SCORE_LIMIT int32 = 3

func main() {
	host := flag.Bool("host", false, "to host server of the game")
	flag.Parse()

	if *host == true {
		server.StartServer()

	}

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "pong")
	defer rl.CloseWindow()

	g := server.NewGame()
	g.Client = client.NewClient()

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		handleMovement(g)
		if g.Ball.IsActive {
			g.Ball.Position.X = g.Ball.Position.X + float32(g.Ball.Speed.X)
			g.Ball.Position.Y = g.Ball.Position.Y + float32(g.Ball.Speed.Y)
		} else {
			g.Ball.Position = server.Coordinate{X: g.Red.Position.X + 2*g.Red.Size.X + g.Ball.Radius, Y: g.Red.Position.Y}
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

		if rl.CheckCollisionCircleRec(rl.Vector2(g.Ball.Position), g.Ball.Radius, rl.Rectangle{X: g.Red.Position.X, Y: g.Red.Position.Y - g.Red.Size.Y/2, Width: g.Red.Size.X, Height: g.Red.Size.Y}) {
			g.Ball.Speed.X = 3.0
			g.Ball.Speed.Y = (g.Ball.Position.Y - g.Red.Position.Y) / (g.Red.Size.Y / 2) * 5

		}
		if rl.CheckCollisionCircleRec(rl.Vector2(g.Ball.Position), 10.0, rl.Rectangle{X: g.Blue.Position.X, Y: g.Blue.Position.Y - g.Blue.Size.Y/2, Width: g.Blue.Size.X, Height: g.Blue.Size.Y}) {
			g.Ball.Speed.X = -3.0
			g.Ball.Speed.Y = (g.Ball.Position.Y - g.Blue.Position.Y) / (g.Blue.Size.Y / 2) * 5
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

func handleMovement(g *server.Game) {
	if rl.IsKeyDown(rl.KeyJ) {
		_, err := g.Client.Conn.Write([]byte("RKEY_J\n"))
		if err != nil {
			log.Println("pressing j is not sending data")
		}
		var res [256]byte

		g.Client.Conn.ReadFromUDP(res[0:])
		if err != nil {
			log.Println("unable to get response for the server")
		}
		log.Printf("response from the server RKEY_J: %v", string(res[0:]))
		var resVal server.Coordinate = g.Red.Position
		err = json.Unmarshal(res[:], &resVal)
		if err != nil {
			log.Println(err.Error())
		}
		g.Red.Position = resVal
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
		if g.Blue.Position.X > 0 {
			g.Blue.Position.X = g.Blue.Position.X - 2
		}
	}

	if rl.IsKeyDown(rl.KeyF) {
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

// func Reset() {
// 	if g.Red.Score >= SCORE_LIMIT {
// 		rl.BeginDrawing()
// 		rl.ClearBackground(rl.Black)
// 		rl.DrawText("RED WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Red)
// 		rl.EndDrawing()
// 		rl.WaitTime(2)
// 		os.Exit(0)
// 	}
//
// 	if g.Blue.Score >= SCORE_LIMIT {
// 		rl.BeginDrawing()
// 		rl.ClearBackground(rl.Black)
// 		rl.DrawText("BLUE WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Blue)
// 		rl.EndDrawing()
// 		rl.WaitTime(2)
// 		os.Exit(0)
// 	}
//
// 	rl.BeginDrawing()
// 	rl.ClearBackground(rl.Black)
// 	rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", g.Red.Score, g.Blue.Score), int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 24, rl.White)
// 	rl.EndDrawing()
// 	rl.WaitTime(2.5)
//
// 	// INFO: reset blue
// 	g.Blue.Position.X = float32(SCREEN_WIDTH) - 10
// 	g.Blue.Position.Y = 200
// 	g.Blue.Size.X = 10
// 	g.Blue.Size.Y = 100
//
// 	// INFO: reset red
// 	g.Red.Position.X = 0
// 	g.Red.Position.Y = 200
// 	g.Red.Size.X = 10
// 	g.Red.Size.Y = 100
//
// 	// INFO: reset ball
// 	g.Ball.Position = rl.Vector2{X: g.Red.Position.X + 2*g.Red.Size.X + g.Ball.Radius, Y: g.Red.Position.Y}
// 	g.Ball.Speed.Y = 0.0
// 	g.Ball.Speed.X = 3.0
// 	g.Ball.IsActive = false
// }

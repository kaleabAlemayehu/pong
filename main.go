package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"playground/raylib-go/client"
	model "playground/raylib-go/models"
	"playground/raylib-go/server"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450
const SCORE_LIMIT int32 = 3

var IS_RED bool

var g *model.Game = &model.Game{
	Red: model.Player{
		Position: rl.Vector2{
			X: 0,
			Y: 200,
		},
		Size: rl.Vector2{
			X: 10,
			Y: 100,
		},
	},

	Blue: model.Player{
		Position: rl.Vector2{
			X: float32(SCREEN_WIDTH) - 10,
			Y: 200,
		},
		Size: rl.Vector2{
			X: 10,
			Y: 100,
		},
	},
	Ball: model.Ball{
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
	serverIp := flag.String("join", "", "ip address of the server")
	flag.Parse()
	message := make(chan *model.Game, 10000)
	input := make(chan string, 10000)

	if *host == true {
		IS_RED = true
		server.StartServer()
	}

	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "pong")
	defer rl.CloseWindow()

	client.ListeningClient(input, message, *serverIp)

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {

		select {
		case msg, ok := <-message:
			if !ok {
				log.Printf("there is something wrong with the recieved data\n")
			}
			g = msg

		default:
			handleMovement(input)
			if rl.CheckCollisionCircleRec(rl.Vector2(g.Ball.Position), g.Ball.Radius, rl.Rectangle{X: g.Red.Position.X, Y: g.Red.Position.Y - g.Red.Size.Y/2, Width: g.Red.Size.X, Height: g.Red.Size.Y}) {
				// send red collition with ball
				input <- "R_B"
			}
			if rl.CheckCollisionCircleRec(rl.Vector2(g.Ball.Position), 10.0, rl.Rectangle{X: g.Blue.Position.X, Y: g.Blue.Position.Y - g.Blue.Size.Y/2, Width: g.Blue.Size.X, Height: g.Blue.Size.Y}) {
				// send blue collition with ball
				input <- "B_B"
			}

			for key, value := range g.Winner {
				if key == "red" && value {
					redWin()
				}
				if key == "blue" && value {
					blueWin()
				}
			}
			if g.ScoreUpdated {
				reset()
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
	if IS_RED {

		if rl.IsKeyDown(rl.KeyJ) {
			input <- "R_J"
		}
		if rl.IsKeyDown(rl.KeyK) {
			input <- "R_K"
		}
		if rl.IsKeyDown(rl.KeyH) {
			input <- "R_H"
		}
		if rl.IsKeyDown(rl.KeyL) {
			input <- "R_L"
		}
		if rl.IsKeyDown(rl.KeySpace) {
			input <- "START_R"
		}
	} else {
		if rl.IsKeyDown(rl.KeyH) {
			input <- "B_H"
		}

		if rl.IsKeyDown(rl.KeyL) {
			input <- "B_L"
		}

		if rl.IsKeyDown(rl.KeyJ) {
			input <- "B_J"
		}
		if rl.IsKeyDown(rl.KeyK) {
			input <- "B_K"
		}
		if rl.IsKeyDown(rl.KeySpace) {
			input <- "START_B"
		}
	}
}

func redWin() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText("RED WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Red)
	rl.EndDrawing()
	rl.WaitTime(2)
	os.Exit(0)
}

func blueWin() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText("BLUE WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Blue)
	rl.EndDrawing()
	rl.WaitTime(2)
	os.Exit(0)
}

func reset() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", g.Red.Score, g.Blue.Score), int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 24, rl.White)
	rl.EndDrawing()
	rl.WaitTime(2.5)
}

package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450
const SCORE_LIMIT int32 = 3

type score struct {
	red  int32
	blue int32
}

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

var RED Player = Player{
	Position: rl.Vector2{
		X: 0,
		Y: 200,
	},
	Size: rl.Vector2{
		X: 10,
		Y: 100,
	},
}

var BLUE Player = Player{
	Position: rl.Vector2{
		X: float32(SCREEN_WIDTH) - 10,
		Y: 200,
	},
	Size: rl.Vector2{
		X: 10,
		Y: 100,
	},
}

var ball Ball = Ball{
	Position: rl.Vector2{
		X: float32(SCREEN_WIDTH) / 2,
		Y: float32(SCREEN_HEIGHT) / 2,
	},
	Speed: rl.Vector2{
		X: 3.0,
		Y: 0.0,
	},
	IsActive: false,
}

var SCORE score

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "pong")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		handleMovement()
		if ball.IsActive {
			ball.Position.X = ball.Position.X + float32(ball.Speed.X)
			ball.Position.Y = ball.Position.Y + float32(ball.Speed.Y)
		} else {
			ball.Position = rl.Vector2{X: RED.Position.X + 2*RED.Size.X + ball.Radius, Y: RED.Position.Y}
		}
		// ball.Position.Y = ball.Position.Y + float32(ball.Speed.Y)
		if ball.Position.X <= 0 {
			BLUE.Score = BLUE.Score + 1
			// reset
			ball.reset()
		}
		if ball.Position.X >= float32(SCREEN_WIDTH) {
			RED.Score = RED.Score + 1
			// reset
			ball.reset()
		}

		if ball.Position.Y <= 0 {
			ball.Speed.Y = 3.0
		}
		if ball.Position.Y >= float32(SCREEN_HEIGHT) {
			ball.Speed.Y = -3.0
		}

		if rl.CheckCollisionCircleRec(ball.Position, ball.Radius, rl.Rectangle{X: RED.Position.X, Y: RED.Position.Y - RED.Size.Y/2, Width: RED.Size.X, Height: RED.Size.Y}) {
			ball.Speed.X = 3.0
			ball.Speed.Y = (ball.Position.Y - RED.Position.Y) / (RED.Size.Y / 2) * 5

		}
		if rl.CheckCollisionCircleRec(ball.Position, 10.0, rl.Rectangle{X: BLUE.Position.X, Y: BLUE.Position.Y - BLUE.Size.Y/2, Width: BLUE.Size.X, Height: BLUE.Size.Y}) {
			ball.Speed.X = -3.0
			ball.Speed.Y = (ball.Position.Y - BLUE.Position.Y) / (BLUE.Size.Y / 2) * 5
			// INFO: ball.speed.x = (ball.position.x - player.position.x)/(player.size.x/2)*5;
		}

		// INFO: drawing start
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawRectangleRec(rl.Rectangle{X: RED.Position.X, Y: RED.Position.Y - (RED.Size.Y / 2), Width: RED.Size.X, Height: RED.Size.Y}, rl.Red)
		rl.DrawRectangleRec(rl.Rectangle{X: BLUE.Position.X, Y: BLUE.Position.Y - (BLUE.Size.Y / 2), Width: BLUE.Size.X, Height: BLUE.Size.Y}, rl.Blue)
		rl.DrawCircleV(ball.Position, 10.0, rl.White)
		rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", RED.Score, BLUE.Score), int32(rl.GetScreenWidth()/2)-115, 0, 20, rl.White)
		rl.EndDrawing()
	}
}

func handleMovement() {
	if rl.IsKeyDown(rl.KeyJ) {
		if RED.Position.Y < float32(SCREEN_HEIGHT)-RED.Size.Y/2 {
			RED.Position.Y = RED.Position.Y + 2
		}
	}
	if rl.IsKeyDown(rl.KeyK) {
		if RED.Position.Y > RED.Size.Y/2 {
			RED.Position.Y = RED.Position.Y - 2
		}
	}
	if rl.IsKeyDown(rl.KeyH) {
		if RED.Position.X > 0 {
			RED.Position.X = RED.Position.X - 2
		}
	}
	if rl.IsKeyDown(rl.KeyL) {
		if RED.Position.X < float32(SCREEN_WIDTH)-RED.Size.X {
			RED.Position.X = RED.Position.X + 2
		}
	}

	if rl.IsKeyDown(rl.KeyA) {
		if BLUE.Position.X > 0 {
			BLUE.Position.X = BLUE.Position.X - 2
		}
	}

	if rl.IsKeyDown(rl.KeyF) {
		if BLUE.Position.X < float32(SCREEN_WIDTH)-BLUE.Size.X {
			BLUE.Position.X = BLUE.Position.X + 2
		}
	}

	if rl.IsKeyDown(rl.KeyS) {
		if BLUE.Position.Y < float32(SCREEN_HEIGHT)-BLUE.Size.Y/2 {
			BLUE.Position.Y = BLUE.Position.Y + 2
		}
	}
	if rl.IsKeyDown(rl.KeyD) {
		if BLUE.Position.Y > BLUE.Size.Y/2 {
			BLUE.Position.Y = BLUE.Position.Y - 2
		}
	}
	if rl.IsKeyDown(rl.KeySpace) {
		ball.IsActive = true
	}
}

func (b *Ball) reset() {
	if RED.Score >= SCORE_LIMIT {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("RED WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Red)
		rl.EndDrawing()
		rl.WaitTime(2)
		os.Exit(0)
	}

	if BLUE.Score >= SCORE_LIMIT {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("BLUE WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Blue)
		rl.EndDrawing()
		rl.WaitTime(2)
		os.Exit(0)
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", RED.Score, BLUE.Score), int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 24, rl.White)
	rl.EndDrawing()
	rl.WaitTime(2.5)

	// INFO: reset blue
	BLUE.Position.X = float32(SCREEN_WIDTH) - 10
	BLUE.Position.Y = 200
	BLUE.Size.X = 10
	BLUE.Size.Y = 100

	// INFO: reset red
	RED.Position.X = 0
	RED.Position.Y = 200
	RED.Size.X = 10
	RED.Size.Y = 100

	// INFO: reset ball
	ball.Position = rl.Vector2{X: RED.Position.X + 2*RED.Size.X + b.Radius, Y: RED.Position.Y}
	b.Speed.Y = 0.0
	b.Speed.X = 3.0
	b.IsActive = false
}

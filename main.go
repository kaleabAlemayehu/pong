package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450
const SCORE_LIMIT int32 = 5

type score struct {
	red  int32
	blue int32
}

type ball struct {
	Position   rl.Vector2
	DirectionX float32
	DirectionY float32
}

var SCORE score

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "pong")
	defer rl.CloseWindow()

	red := rl.Rectangle{
		X:      0,
		Y:      200,
		Width:  10,
		Height: 100,
	}
	blue := rl.Rectangle{
		X:      float32(SCREEN_WIDTH) - 10,
		Y:      200,
		Width:  10,
		Height: 100,
	}

	ball := ball{
		Position: rl.Vector2{
			X: float32(SCREEN_WIDTH) / 2,
			Y: float32(SCREEN_HEIGHT) / 2,
		},
		DirectionX: -3.0,
		DirectionY: -3.0,
	}
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		handleMovement(&red, &blue)
		// INFO:  make the ball occilate
		ball.Position.X = ball.Position.X + float32(ball.DirectionX)
		ball.Position.Y = ball.Position.Y + float32(ball.DirectionY)
		if ball.Position.X <= 0 {
			SCORE.blue = SCORE.blue + 1
			// reset
			ball.reset()
		}
		if ball.Position.X >= float32(SCREEN_WIDTH) {
			SCORE.red = SCORE.red + 1
			// reset
			ball.reset()
		}

		if ball.Position.Y <= 0 {
			ball.DirectionY = 3.0
		}
		if ball.Position.Y >= float32(SCREEN_HEIGHT) {
			ball.DirectionY = -3.0
		}

		if rl.CheckCollisionCircleRec(ball.Position, 10.0, red) {
			ball.DirectionX = 3.0
		}
		if rl.CheckCollisionCircleRec(ball.Position, 10.0, blue) {
			ball.DirectionX = -3.0
		}

		// INFO: drawing start
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawRectangleRec(red, rl.Red)
		rl.DrawRectangleRec(blue, rl.Blue)
		rl.DrawCircleV(ball.Position, 10.0, rl.White)
		rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", SCORE.red, SCORE.blue), int32(rl.GetScreenWidth()/2)-115, 0, 20, rl.White)
		rl.EndDrawing()
	}
}

func handleMovement(red *rl.Rectangle, blue *rl.Rectangle) {
	if rl.IsKeyDown(rl.KeyJ) {
		if red.Y < float32(SCREEN_HEIGHT)-red.Height {
			red.Y = red.Y + 2
		}
	}
	if rl.IsKeyDown(rl.KeyK) {
		if red.Y > 0 {
			red.Y = red.Y - 2
		}
	}
	if rl.IsKeyDown(rl.KeyH) {
		if red.X > 0 {
			red.X = red.X - 2
		}
	}
	if rl.IsKeyDown(rl.KeyL) {
		if red.X < float32(SCREEN_WIDTH)-red.Width {
			red.X = red.X + 2
		}
	}

	if rl.IsKeyDown(rl.KeyA) {
		if blue.X > 0 {
			blue.X = blue.X - 2
		}
	}

	if rl.IsKeyDown(rl.KeyF) {
		if blue.X < float32(SCREEN_WIDTH)-blue.Width {
			blue.X = blue.X + 2
		}
	}

	if rl.IsKeyDown(rl.KeyS) {
		if blue.Y < float32(SCREEN_HEIGHT)-blue.Height {
			blue.Y = blue.Y + 2
		}
	}
	if rl.IsKeyDown(rl.KeyD) {
		if blue.Y > 0 {
			blue.Y = blue.Y - 2
		}
	}
}

func (b *ball) reset() {
	if SCORE.red >= SCORE_LIMIT {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("RED WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Red)
		rl.EndDrawing()
		rl.WaitTime(2)
		os.Exit(0)
	}

	if SCORE.blue >= SCORE_LIMIT {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("BLUE WINS!!", int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 32, rl.Blue)
		rl.EndDrawing()
		rl.WaitTime(2)
		os.Exit(0)
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText(fmt.Sprintf("RED: %v |<=>| BLUE: %v", SCORE.red, SCORE.blue), int32(rl.GetScreenWidth()/2)-115, int32(rl.GetScreenHeight()/2), 24, rl.White)
	rl.EndDrawing()
	rl.WaitTime(2.5)

	b.Position = rl.Vector2{
		X: float32(SCREEN_WIDTH) / 2,
		Y: float32(SCREEN_HEIGHT) / 2,
	}
	b.DirectionX = -3.0
	b.DirectionY = -3.0
}

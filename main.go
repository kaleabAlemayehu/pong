package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450

type ball struct {
	Position   rl.Vector2
	DirectionX float32
	DirectionY float32
}

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
			ball.DirectionX = 3.0
		}
		if ball.Position.X >= float32(SCREEN_WIDTH) {
			ball.DirectionX = -3.0
		}

		if ball.Position.Y <= 0 {
			ball.DirectionY = 3.0
		}
		if ball.Position.Y >= float32(SCREEN_HEIGHT) {
			ball.DirectionY = -3.0
		}

		log.Printf("ball positionX %v\n", ball.Position.X)
		log.Printf("ball positionY %v\n", ball.Position.Y)

		if rl.CheckCollisionCircleRec(ball.Position, 10.0, red) {
			log.Printf("red hit me\n")
			ball.DirectionX = 3.0
		}
		if rl.CheckCollisionCircleRec(ball.Position, 10.0, blue) {
			log.Printf("blue hit me")
			ball.DirectionX = -3.0
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawRectangleRec(red, rl.Red)
		rl.DrawRectangleRec(blue, rl.Blue)
		rl.DrawCircleV(ball.Position, 10.0, rl.White)
		rl.EndDrawing()
	}
}

func handleMovement(red *rl.Rectangle, blue *rl.Rectangle) {
	if rl.IsKeyDown(rl.KeyJ) {
		red.Y = red.Y + 2
	}
	if rl.IsKeyDown(rl.KeyK) {
		red.Y = red.Y - 2
	}
	if rl.IsKeyDown(rl.KeyH) {
		red.X = red.X - 2
	}
	if rl.IsKeyDown(rl.KeyL) {
		red.X = red.X + 2
	}

	if rl.IsKeyDown(rl.KeyA) {
		blue.X = blue.X - 2
	}
	if rl.IsKeyDown(rl.KeyF) {
		blue.X = blue.X + 2
	}
	if rl.IsKeyDown(rl.KeyS) {
		blue.Y = blue.Y + 2
	}
	if rl.IsKeyDown(rl.KeyD) {
		blue.Y = blue.Y - 2
	}
}

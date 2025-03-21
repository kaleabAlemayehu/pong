package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450

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
	ball := rl.Vector2{
		X: float32(SCREEN_WIDTH) / 2,
		Y: float32(SCREEN_HEIGHT) / 2,
	}
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		handleMovement(red, blue)
		if rl.CheckCollisionCircleRec(ball, 10.0, red) {
			log.Printf("red hit me")
		}

		if rl.CheckCollisionCircleRec(ball, 10.0, blue) {
			log.Printf("blue hit me")
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawRectangleRec(red, rl.Red)
		rl.DrawRectangleRec(blue, rl.Blue)
		rl.DrawCircleV(ball, 10.0, rl.White)
		rl.EndDrawing()
	}
}

func handleMovement(red rl.Rectangle, blue rl.Rectangle) {
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

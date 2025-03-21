package main

import rl "github.com/gen2brain/raylib-go/raylib"

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
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

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

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawRectangleRec(red, rl.Red)
		rl.EndDrawing()
	}
}

package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"net"
)

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

type Client struct {
	Conn *net.UDPConn
}

type Game struct {
	Red    Player `json:"red"`
	Blue   Player `json:"blue"`
	Ball   Ball   `json:"ball"`
	Client Client `json:"client"`
}

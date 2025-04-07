package models

import (
	"net"
)

type InputMessage struct {
	Addr string
	Cmd  string // e.g., "R_J", "R_K"
}

type Coordinate struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type ServerBall struct {
	Position Coordinate `json:"position"`
	Speed    Coordinate `json:"speed"`
	Radius   float32    `json:"radius"`
	IsActive bool       `json:"is_active"`
}

type ServerPlayer struct {
	Position Coordinate `json:"position"`
	Size     Coordinate `json:"size"`
	Score    int32      `json:"score"`
}

type ServerGame struct {
	Red    ServerPlayer        `json:"red"`
	Blue   ServerPlayer        `json:"blue"`
	Ball   ServerBall          `json:"ball"`
	Conn   map[string]net.Addr `json:"conn"`
	Client Client              `json:"client"`
	Winner map[string]bool     `json:"winner_player"`
}

package client

import (
	"log"
	"net"
)

type Client struct {
	Conn *net.UDPConn
}

func NewClient() *Client {
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8000")
	if err != nil {
		log.Fatalf("client unable to resovle upd address Error: %v ", err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("client unable to connect upd server Error: %v ", err.Error())
	}
	log.Println("connected succesfully to the server")
	return &Client{Conn: conn}
}

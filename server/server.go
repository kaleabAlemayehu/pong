package server

import (
	"log"
	"net"
)

func StartServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8000")
	if err != nil {
		log.Fatalf("unable to setup the address => Error:%v ", err.Error())
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("unable to start the server => Error:%v ", err.Error())
	}
	log.Printf("listening connection on port :8000")
	go func() {
		for {
			var buf [512]byte

			_, addr, err := conn.ReadFromUDP(buf[0:])
			if err != nil {
				log.Printf("unable to read: Error: %v", err.Error())
				return
			}
			log.Printf("|> %v", string(buf[0:]))
			conn.WriteToUDP(buf[0:], addr)
		}

	}()
}

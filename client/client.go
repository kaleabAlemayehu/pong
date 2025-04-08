package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	model "playground/raylib-go/models"
	"time"
)

func ListeningClient(input chan string, msg chan *model.Game, localAddr string) {
	// INFO: sending to broadcast network so that i can find the server
	// udpAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:8000")
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:8000", localAddr))
	// udpAddr, err := net.ResolveUDPAddr("udp", "192.168.1.255:8000")
	if err != nil {
		log.Printf("client unable to resovle upd address Error: %v ", err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Printf("client unable to connect upd server Error: %v ", err.Error())
	}
	// INFO: listening routine from server
	go func() {

		res := make([]byte, 10240)
		for {
			var g model.Game
			conn.SetReadDeadline(time.Now().Add(3 * time.Second))
			n, _, err := conn.ReadFromUDP(res[:])
			if err != nil {
				log.Printf("unable to start the server => Error:%v ", err.Error())
			}

			// log.Printf("i receive: %v", string(res[:n]))

			err = json.Unmarshal(res[:n], &g)
			if err != nil {
				log.Printf("error happend when recieved data from the server unmarshaled\n")
				log.Printf("client routine the error: %v", err.Error())
			}
			select {
			case msg <- &g: // Non-blocking send
			default:
				log.Println("Warning: Message dropped (buffer full)")
			}
		}
	}()
	// INFO: sending to server routine
	go func() {
		for {
			select {
			case i := <-input:
				{
					log.Printf("i sent %v\n", i)
					_, err := conn.Write([]byte(i))
					if err != nil {
						log.Printf("error while sending: %v", err.Error())
					}
				}
			}
		}
	}()
}

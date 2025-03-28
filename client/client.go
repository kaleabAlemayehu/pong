package client

import (
	"log"
	"net"
)

type Client struct {
	Conn *net.UDPConn
}

func ListeningClient(input chan string, msg chan []byte) {
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8000")
	if err != nil {
		log.Printf("client unable to resovle upd address Error: %v ", err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Printf("client unable to connect upd server Error: %v ", err.Error())
	}
	log.Println("this does print out")
	// listening routine from server
	go func() {

		res := make([]byte, 10240)
		for {
			n, _, err := conn.ReadFrom(res[:])
			_ = n
			if err != nil {
				log.Printf("unable to start the server => Error:%v ", err.Error())
			}
			select {
			case msg <- res[:n]: // Non-blocking send
			default:
				log.Println("Warning: Message dropped (buffer full)")
			}
		}
	}()
	// sending to server routine
	go func() {
		for {
			select {
			case i := <-input:
				{
					n, err := conn.Write([]byte(i))
					if err != nil {
						log.Printf("error while sending: %v", err.Error())
					} else {
						log.Printf("%d bytes sent to server\n", n)
					}
				}
			}
		}
	}()
}

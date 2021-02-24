package cmd

import (
	"log"
	"net"
	"os"
	"time"

	socks "github.com/fangdingjun/socks-go"
)

func openSocks(socksAddress string) {
	// Listen host:port for socks server if required envs provided
	if os.Getenv("SOCKS_ENABLED") == "true" {

		conn, err := net.Listen("tcp", socksAddress)
		if err != nil {
			log.Fatal(err)
		}

		for {
			c, err := conn.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("connected from %s", c.RemoteAddr())
			// Default keepalive - 15s
			d := net.Dialer{Timeout: 10 * time.Second}
			s := socks.Conn{Conn: c, Dial: d.Dial}
			go s.Serve()
		}
	}
}

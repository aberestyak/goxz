package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type endpoint struct {
	Host string
	Port int
}

func (endpoint *endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Printf("error while copy remote->local: %s\n", err)
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Printf("error while copy local->remote: %s\n", err)
		}
		chDone <- true
	}()

	<-chDone
}

func main() {
	if err := validateEnvs(); err != nil {
		log.Fatalf("Error while envs validation: %s\n", err)
	}
	// Create socks server
	go openSocks(localEndpoint.String())

	sshConfig, err := makeSSHConfig(os.Getenv("REMOTE_SSH_USER"))
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
	if err != nil {
		log.Fatalf("Dial INTO remote server error: %s\n", err)
	}
	log.Println("Connected to " + serverEndpoint.String())

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		log.Fatalf("Listen open port ON remote server error: %s\n", err)
	}
	defer listener.Close()

	// handle incoming connections on reverse forwarded port
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", localEndpoint.String())
		if err != nil {
			log.Fatalf("Dial INTO local service error: %s\n", err)
		}
		// Remote connection
		client, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		handleClient(client, local)
	}

}

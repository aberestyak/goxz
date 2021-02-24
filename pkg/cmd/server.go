package cmd

import (
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func startServer() {
	if err := validateEnvs(); err != nil {
		log.Fatalf("Error while envs validation: %s\n", err)
	}
	// Create socks server
	go openSocks(localEndpoint.String())
	log.Printf("Started SOCKS server on %s\n", localEndpoint.String())

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
	defer serverConn.Close()

	// Listen on remote server port
	// If remote inerface no fully configured (https://github.com/golang/go/issues/18806)
	// port will be opened on 0.0.0.0
	listener, err := serverConn.Listen("tcp4", remoteEndpoint.String())
	if err != nil {
		log.Fatalf("Listen open port ON remote server error: %s\n", err)
	}
	defer listener.Close()

	// handle incoming connections on reverse forwarded port
	for {
		// Remote connection
		client, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		// Process each connection in separate goroutine
		go func() {
			// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
			local, err := net.Dial("tcp", localEndpoint.String())
			if err != nil {
				log.Fatalf("Dial INTO local service error: %s\n", err)
			}
			transferData(client, local)
			local.Close()
			return
		}()
	}

}

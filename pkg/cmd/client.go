package cmd

import (
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func startClient() {
	if err := validateEnvs(); err != nil {
		log.Fatalf("Error while envs validation: %s\n", err)
	}

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

	// Listen on remote local port
	local, err := net.Listen("tcp", localEndpoint.String())
	if err != nil {
		log.Fatalf("Listen open port ON local server error: %s\n", err)
	}
	defer local.Close()
	log.Println("Start linstening " + localEndpoint.String())

	// handle incoming connections on reverse forwarded port
	for {
		client, err := local.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		// Establish new connection with remote server for each incoming connection
		go func() {
			remote, err := serverConn.Dial("tcp4", remoteEndpoint.String())
			if err != nil {
				log.Fatalln(err)
			}
			transferData(client, remote)
			remote.Close()
		}()
	}
}

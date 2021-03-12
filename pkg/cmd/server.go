package cmd

import (
	"errors"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func serveTCP() error {
	sshConfig, err := makeSSHConfig(os.Getenv("BASTION_SSH_USER"))
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", bastionEndpoint.String(), sshConfig)
	if err != nil {
		return errors.New("Dial INTO remote server error: " + err.Error())
	}
	log.Println("Connected to " + bastionEndpoint.String())
	defer serverConn.Close()

	// Listen on remote server port
	// If remote inerface no fully configured (https://github.com/golang/go/issues/18806)
	// port will be opened on 0.0.0.0
	listener, err := serverConn.Listen("tcp4", bastionTCPEndpoint.String())
	if err != nil {
		return errors.New("Listen open port ON remote server error: " + err.Error())
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
			// Open a (local) connection to localTCPEndpoint whose content will be forwarded so serverEndpoint
			local, err := net.Dial("tcp", localTCPEndpoint.String())
			if err != nil {
				log.Fatalf("Dial INTO local service error: %s\n", err)
			}
			transferData(client, local)
			local.Close()
		}()
	}
}

func serveUDPOverTCP() error {
	// Open local TCP port
	localTCPConn, err := net.Listen("tcp", serverTCPtoUDPEndpoint.String())
	if err != nil {
		return errors.New("Listen open port ON local server error: " + err.Error())
	}
	defer localTCPConn.Close()
	log.Println("Start linstening TCP to UDP: " + serverTCPtoUDPEndpoint.String())

	for {
		udpBuf := make([]byte, 1500)
		tcpBuf := make([]byte, 1500)

		clientTCP, err := localTCPConn.Accept()
		if err != nil {
			return err
		}
		defer clientTCP.Close()

		go func() {
			remoteUDPConn, err := net.DialUDP("udp", nil, remoteUDPEndpoint)
			if err != nil {
				log.Printf("Cant dial remote TCP to UDP endpoint: %v", err)
			}
			defer remoteUDPConn.Close()

			tcpBufLen, _ := clientTCP.Read(tcpBuf)

			// Send UDP to remote (target host(DNS))
			if _, err := remoteUDPConn.Write(tcpBuf[:tcpBufLen]); err != nil {
				log.Printf("Could not send TCP to UDP response: %v", err)
			}
			for {
				// Read UDP from remote (target host(DNS))
				udpBufLen, _, err := remoteUDPConn.ReadFromUDP(udpBuf)
				if err != nil {
					log.Printf("Could not read UDP response: %v", err)
				}
				// Finally send UDP response via TCP
				if _, err := clientTCP.Write(udpBuf[:udpBufLen]); err != nil {
					log.Printf("Could not send UDP response via TCP: %v", err)
				}
			}
		}()
	}
}

func startServer() {

	// Create socks server
	go openSocks(localTCPEndpoint.String())
	log.Printf("Started SOCKS server on %s\n", localTCPEndpoint.String())

	go serveUDPOverTCP()

	if err := serveTCP(); err != nil {
		log.Fatal(err.Error())
	}
}

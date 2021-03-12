package cmd

import (
	"errors"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func forwardTCP() error {
	sshConfig, err := makeSSHConfig(os.Getenv("BASTION_SSH_USER"))
	if err != nil {
		return err
	}
	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", bastionEndpoint.String(), sshConfig)
	if err != nil {
		return errors.New("Dial INTO remote server error: " + err.Error())
	}
	log.Println("Connected to " + bastionEndpoint.String())
	defer serverConn.Close()

	// Listen on remote local port
	local, err := net.Listen("tcp", localTCPEndpoint.String())
	if err != nil {
		return errors.New("Listen open port ON local server error: " + err.Error())
	}
	defer local.Close()
	log.Println("Start linstening " + localTCPEndpoint.String())

	// handle incoming connections on reverse forwarded port
	for {
		client, err := local.Accept()
		if err != nil {
			return err
		}
		// Establish new connection with remote server for each incoming connection
		go func() {
			remote, err := serverConn.Dial("tcp4", bastionTCPEndpoint.String())
			if err != nil {
				log.Fatalln(err)
			}
			transferData(client, remote)
			remote.Close()
		}()
	}
}

func forwardUDPOverTCP() error {
	// Open local UDP port
	localTCPConn, err := net.ListenUDP("udp", clientUDPEndpoint)
	if err != nil {
		return errors.New("Listen open port ON local server error: " + err.Error())
	}
	defer localTCPConn.Close()
	log.Println("Start linstening " + clientUDPEndpoint.String())

	for {
		udpBuf := make([]byte, 1500)
		udpBufLen, clientAddress, err := localTCPConn.ReadFromUDP(udpBuf)
		if err != nil {
			return err
		}
		go func() {
			remoteTCPConn, err := net.Dial("tcp4", serverTCPtoUDPEndpoint.String())
			if err != nil {
				log.Fatalf("Cant dial remote TCP to UDP endpoint: %v", err)
			}
			TCPtoUDP(udpBuf, udpBufLen, clientAddress, localTCPConn, remoteTCPConn)
		}()
	}
}

func startClient() {

	go forwardUDPOverTCP()

	if err := forwardTCP(); err != nil {
		log.Fatal(err.Error())
	}
}

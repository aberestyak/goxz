package cmd

import (
	"io"
	"log"
	"net"
	"strings"
)

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
func transferData(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		if _, err := io.Copy(client, remote); err != nil {
			// To prevent "errors" whith closed connection
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("error while copy remote->local: %s\n", err)
			}
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		if _, err := io.Copy(remote, client); err != nil {
			// To prevent "errors" whith closed connection
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("error while copy local->remote: %s\n", err)
			}
		}
		chDone <- true
	}()

	<-chDone
}

func TCPtoUDP(udpBuf []byte, udpBufLen int, sourceAddress *net.UDPAddr, localConn *net.UDPConn, remoteTCPConn net.Conn) {
	tcpBuf := make([]byte, 1500)
	defer remoteTCPConn.Close()

	if _, err := remoteTCPConn.Write(udpBuf[:udpBufLen]); err != nil {
		log.Fatalf("Could not write packet to tcp connection: %v", err)
	}

	tcpBufLen, _ := remoteTCPConn.Read(tcpBuf)
	if _, err := localConn.WriteToUDP(tcpBuf[:tcpBufLen], sourceAddress); err != nil {
		log.Fatalf("Could not send TCP to UDP response: %v", err)
	}
}

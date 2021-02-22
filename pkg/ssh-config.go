package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Get private key for ssh authentication
func parsePrivateKey() (ssh.Signer, error) {
	// Check SSH_KEY_PATH env, and file existance
	if os.Getenv("SSH_KEY_PATH") != "" {
		if _, err := os.Stat(os.Getenv("SSH_KEY_PATH")); err != nil || os.IsExist(err) {
			return nil, errors.New("File " + os.Getenv("SSH_KEY_PATH") + " doesn't exist!")
		}
		buff, _ := ioutil.ReadFile(os.Getenv("SSH_KEY_PATH"))
		// read $HOME/.ssh/id_rsa, check if SSH_KEY_PASSPHRASE provided
		if os.Getenv("SSH_KEY_PASSPHRASE") != "" {
			signer, err := ssh.ParsePrivateKeyWithPassphrase(buff, []byte(os.Getenv("SSH_KEY_PASSPHRASE")))
			return signer, err
		}
		return ssh.ParsePrivateKey(buff)
	}
	return nil, errors.New("SSH_KEY_PATH not provided")
}

func openSSHAgent() (agent.ExtendedAgent, error) {
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		socket := os.Getenv("SSH_AUTH_SOCK")

		conn, err := net.Dial("unix", socket)
		if err != nil {
			return nil, err
		}
		return agent.NewClient(conn), nil
	}
	return nil, errors.New("No SSH_AUTH_SOCK provided")
}

// SSH config will use 3 strategies: ssh-agent, keys and password
func makeSSHConfig(user string) (*ssh.ClientConfig, error) {
	// Create basic client config and add methods if they are valid
	config := ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	// Create private key auth, if SSH_KEY_PATH is set
	if key, err := parsePrivateKey(); err == nil {
		config.Auth = append(config.Auth, ssh.PublicKeys(key))
	} else {
		log.Printf("Skipping private key authentication: %v\n", err)
	}

	// Create ssh-aget auth, if SSH_AUTH_SOCK is set
	if agentAuth, err := openSSHAgent(); err == nil {
		config.Auth = append(config.Auth, ssh.PublicKeysCallback(agentAuth.Signers))
	} else {
		log.Printf("Skipping ssh-agent authentication: %v\n", err)
	}

	// Create user password auth, if SSH_USER_PASSWORD is set
	if os.Getenv("SSH_USER_PASSWORD") != "" {
		config.Auth = append(config.Auth, ssh.Password(os.Getenv("SSH_USER_PASSWORD")))
	}

	return &config, nil
}

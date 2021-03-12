package cmd // local service to be forwarded
import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

type endpoint struct {
	Host string
	Port int
}

var mandatoryServerEnvs = []string{"BASTION_HOST_SSH", "BASTION_PORT_SSH", "BASTION_SSH_USER", "LOCAL_HOST_TCP", "LOCAL_PORT_TCP", "SERVER_HOST_TCP-UDP", "SERVER_PORT_TCP-UDP", "REMOTE_HOST_UDP", "REMOTE_PORT_UDP"}
var mandatoryClientEnvs = []string{"BASTION_HOST_SSH", "BASTION_PORT_SSH", "BASTION_SSH_USER", "LOCAL_HOST_TCP", "LOCAL_PORT_TCP", "CLIENT_HOST_UDP", "CLIENT_PORT_UDP", "SERVER_HOST_TCP-UDP", "SERVER_PORT_TCP-UDP"}
var authEnvs = []string{"SSH_KEY_PATH", "SSH_AUTH_SOCK", "SSH_USER_PASSWORD"}

func (endpoint *endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// Check that all mandotory and, at least one authEnv, envs are set.
func validateEnvs(mode string) error {
	validate := func(mandatoryEnvs []string) error {
		for _, env := range mandatoryEnvs {
			if os.Getenv(env) == "" {
				return errors.New("mandatory env " + env + " is not set")
			}
		}
		authEnvsFound := false
		for _, env := range authEnvs {
			if os.Getenv(env) != "" {
				authEnvsFound = true
				break
			}
		}
		if authEnvsFound {
			return nil
		} else {
			return errors.New("none of auth envs (SSH_KEY_PATH, SSH_AUTH_SOCK, SSH_USER_PASSWORD) is set")
		}
	}

	switch mode {
	case "server":
		if err := validate(mandatoryServerEnvs); err != nil {
			return err
		}
	case "client":
		if err := validate(mandatoryClientEnvs); err != nil {
			return err
		}
	default:
		return ErrIncorrectArg
	}
	return nil
}

var localTCPEndpoint = endpoint{
	Host: os.Getenv("LOCAL_HOST_TCP"),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("LOCAL_PORT_TCP"))
		return port
	}(),
}

var clientUDPEndpoint = &net.UDPAddr{
	IP: net.ParseIP(os.Getenv("CLIENT_HOST_UDP")),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("CLIENT_PORT_UDP"))
		return port
	}(),
}

// remote SSH server
var bastionEndpoint = endpoint{
	Host: os.Getenv("BASTION_HOST_SSH"),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("BASTION_PORT_SSH"))
		return port
	}(),
}

// remote forwarding port (on remote SSH server network)
var bastionTCPEndpoint = endpoint{
	Host: os.Getenv("BASTION_HOST_TCP"),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("BASTION_PORT_TCP"))
		return port
	}(),
}

// remote for to forward UDP traffic. Bastion here
var serverTCPtoUDPEndpoint = endpoint{
	Host: os.Getenv("SERVER_HOST_TCP-UDP"),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("SERVER_PORT_TCP-UDP"))
		return port
	}(),
}

var remoteUDPEndpoint = &net.UDPAddr{
	IP: net.ParseIP(os.Getenv("REMOTE_HOST_UDP")),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("REMOTE_PORT_UDP"))
		return port
	}(),
}

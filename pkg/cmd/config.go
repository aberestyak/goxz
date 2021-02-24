package cmd // local service to be forwarded
import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type endpoint struct {
	Host string
	Port int
}

var mandatoryEnvs = []string{"LOCAL_HOST", "LOCAL_PORT", "REMOTE_SSH_HOST", "REMOTE_SSH_PORT", "REMOTE_SSH_USER", "REMOTE_FORWARD_HOST", "REMOTE_FORWARD_PORT"}
var authEnvs = []string{"SSH_KEY_PATH", "SSH_AUTH_SOCK", "SSH_USER_PASSWORD"}

func (endpoint *endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// Check that all mandotory and, at least one authEnv, envs are set.
func validateEnvs() error {
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
	}

	return errors.New("none of auth envs (SSH_KEY_PATH, SSH_AUTH_SOCK, SSH_USER_PASSWORD) is set")
}

var localEndpoint = endpoint{
	Host: os.Getenv("LOCAL_HOST"),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("LOCAL_PORT"))
		return port
	}(),
}

// remote SSH server
var serverEndpoint = endpoint{
	Host: os.Getenv("REMOTE_SSH_HOST"),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("REMOTE_SSH_PORT"))
		return port
	}(),
}

// remote forwarding port (on remote SSH server network)
var remoteEndpoint = endpoint{
	Host: os.Getenv("REMOTE_FORWARD_HOST"),
	Port: func() int {
		port, _ := strconv.Atoi(os.Getenv("REMOTE_FORWARD_PORT"))
		return port
	}(),
}

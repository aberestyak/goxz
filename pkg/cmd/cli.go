package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	example = `
	# start goxz with local SOCKS server
	%[1]s server
	# start goxz client for connection to SOCKS server
	%[1]s client
`
)

// ErrInsufficientArgs is thrown if arg len <1 or >2
var ErrInsufficientArgs = fmt.Errorf("\nincorrect number or arguments, see --help for usage instructions")

// ErrInsufficientArgs is thrown if arg len <1 or >2
var ErrIncorrectArg = fmt.Errorf("\nincorrect argument, see --help for usage instructions")

// NewCmdGoXz creates the cobra command to be executed
func NewCmdGoXz() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "goxz <server/client>",
		Short:        "Start goxz server or client",
		Example:      fmt.Sprintf(example, "goxz"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateArgs(args); err != nil {
				return err
			}
			if err := start(args); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

// Validate ensures proper command usage
func validateArgs(args []string) error {
	argLen := len(args)
	if argLen != 1 {
		return ErrInsufficientArgs
	}
	// if len = 1, then it's always 0 index
	if args[0] != "server" && args[0] != "client" {
		return ErrIncorrectArg
	}
	return nil
}

// Start - start client or server
func start(args []string) error {
	if args[0] == "server" {
		startServer()
	} else {
		startClient()
	}

	return nil
}

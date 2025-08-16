package client

import (
	"example.com/db/internal/shutil"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"example.com/db/internal/types"
)

type Redis struct {
	*Client
}

func (c *Redis) SetClient(client *Client) {
	c.Client = client
}

func (c *Redis) Ping() error {
	args := []string{"-u", c.URL}
	args = append(args, "PING")
	cmd := exec.Command(c.getDefaultCommand(), args...)
	err := shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(io.Discard),
		shutil.WithStderr(io.Discard),
	)
	if err == nil {
		fmt.Println(shutil.ColorGreen("pong"))
	}
	return err
}

func (c *Redis) Connect() error {
	args := []string{"-u", c.URL}
	prog := c.Client.GetInteractiveREPL()
	if prog == "" {
		prog = c.getDefaultCommand()
	}
	cmd := exec.Command(prog, args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Redis) RunQuery(query string, args ...string) error {
	defaultArgs := []string{}
	switch c.GetFormat() {
	case types.JSON:
		defaultArgs = append(defaultArgs, "--json")
	case types.CSV:
		defaultArgs = append(defaultArgs, "--csv")
	case "":
		// nothing to-do
	default:
		return fmt.Errorf("%w: driver %s", UnsupportedFormat, c.Driver)
	}

	args = append(args, defaultArgs...)
	args = append(args, "-u", c.URL)
	args = append(args, strings.Fields(query)...)
	// fmt.Print(args)
	cmd := exec.Command(c.getDefaultCommand(), args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Redis) ListTables() error {
	return c.RunQuery("KEYS *")
}

func (c *Redis) String() string {
	return c.URL
}

func (c *Redis) getDefaultCommand() string {
	if shutil.IsCommandInstalled("valkey-cli") {
		return "valkey-cli"
	}
	return "redis-cli"
}

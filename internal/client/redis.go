package client

import (
	"example.com/db/internal/shutil"
	"fmt"
	"io"
	"os/exec"
	_ "strings"
)

type Redis struct {
	*Client
}

func (c *Redis) SetClient(client *Client) {
	c.Client = client
}

func (c *Redis) Ping() error {
	args := []string{}
	args = append(args, "PING")
	cmd := exec.Command(c.getDefaultCommand(), args...)
	err := shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
	if err == nil {
		fmt.Println(shutil.ColorGreen("pong"))
	}
	return err
}

func (c *Redis) Connect() error {
	args := []string{}
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
	args = append(args, query)
	cmd := exec.Command(c.getDefaultCommand(), args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(io.Discard),
		shutil.WithStderr(io.Discard),
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

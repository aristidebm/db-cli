package client

import (
	_ "fmt"
	"os/exec"
	_ "strings"

	"example.com/db/internal/shutil"
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
	cmd := exec.Command("redis-cli", args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Redis) Connect() error {
	args := []string{}
	cmd := exec.Command("redis-cli", args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Redis) RunQuery(query string) error {
	args := []string{}
	// parts := strings.Fields(query)
	args = append(args, query)
	cmd := exec.Command("redis-cli", args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Redis) ListTables() error {
	return c.RunQuery("KEYS *")
}

func (c *Redis) ListDatabases() error {
	return c.RunQuery("CONFIG GET databases")
}

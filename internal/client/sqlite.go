package client

import (
	"example.com/db/internal/shutil"
	"fmt"
	"os/exec"
)

type SQLite struct {
	client *Client
}

func (c *SQLite) SetClient(client *Client) {
	c.client = client
}

func (c *SQLite) Ping() error {
	args := []string{}
	args = append(args, c.client.Path, "SELECT 1;")
	cmd := exec.Command("sqlite3", args...)
	return shutil.Run(cmd, shutil.WithStdin(c.client.Stdin), shutil.WithStdout(c.client.Stdout), shutil.WithStderr(c.client.Stderr))
}

func (c *SQLite) Connect() error {
	args := []string{}
	args = append(args, c.client.Path)
	cmd := exec.Command("sqlite3", args...)
	return shutil.RunInteractive(cmd, shutil.WithStdin(c.client.Stdin), shutil.WithStdout(c.client.Stdout), shutil.WithStderr(c.client.Stderr))
}

func (c *SQLite) RunQuery(query string) error {
	if query == "" {
		return c.Connect()
	}
	args := []string{c.client.Path}
	args = append(args, query)
	cmd := exec.Command("sqlite3", args...)
	return shutil.Run(cmd, shutil.WithStdin(c.client.Stdin), shutil.WithStdout(c.client.Stdout), shutil.WithStderr(c.client.Stderr))
}

func (c *SQLite) ListTables() error {
	return c.RunQuery(".tables")
}

func (c *SQLite) ListDatabases() error {
	return fmt.Errorf("%w driver %s", UnsupportedCommand, c.client.Driver)
}

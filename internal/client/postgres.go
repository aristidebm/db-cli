package client

import (
	"example.com/db/internal/shutil"
	"os/exec"
)

type Postgres struct {
	client *Client
}

func (c *Postgres) SetClient(client *Client) {
	c.client = client
}

func (c *Postgres) Ping() error {
	args := []string{}
	args = append(args, "-c", "SELECT 1;")
	cmd := exec.Command("psql", args...)
	return shutil.Run(cmd, shutil.WithStdin(c.client.Stdin), shutil.WithStdout(c.client.Stdout), shutil.WithStderr(c.client.Stderr))
}

func (c *Postgres) Connect() error {
	args := []string{}
	cmd := exec.Command("psql", args...)
	return shutil.RunInteractive(cmd, shutil.WithStdin(c.client.Stdin), shutil.WithStdout(c.client.Stdout), shutil.WithStderr(c.client.Stderr))
}

func (c *Postgres) RunQuery(query string) error {
	args := []string{}
	args = append(args, "-c", query)
	cmd := exec.Command("psql", args...)
	return shutil.Run(cmd, shutil.WithStdin(c.client.Stdin), shutil.WithStdout(c.client.Stdout), shutil.WithStderr(c.client.Stderr))
}

func (c *Postgres) ListTables() error {
	return c.RunQuery("\\dt")
}

func (c *Postgres) ListDatabases() error {
	return c.RunQuery("\\l")
}

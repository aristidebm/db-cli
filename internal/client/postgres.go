package client

import (
	"fmt"
	"os/exec"

	"example.com/db/internal/shutil"
)

type Postgres struct {
	*Client
}

func (c *Postgres) SetClient(client *Client) {
	c.Client = client
}

func (c *Postgres) Ping() error {
	args := []string{}
	args = append(args, "-c", "SELECT 1;")
	cmd := exec.Command("psql", args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Postgres) Connect() error {
	args := []string{}
	cmd := exec.Command("psql", args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Postgres) RunQuery(query string) error {
	args := []string{}
	args = append(args, "-c", query)

	switch c.Format {
	case CSV:
		args = append(args, "--pset format csv")
	case HTML:
		args = append(args, "--pset format html")
	case LATEX:
		args = append(args, "--pset format latex")
	case ASCIIDOC:
		args = append(args, "--pset format asciidoc")
	default:
		return fmt.Errorf("%w: driver %s", UnsupportedFormat, c.Driver)
	}

	cmd := exec.Command("psql", args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Postgres) ListTables() error {
	return c.RunQuery("\\dt")
}

func (c *Postgres) ListDatabases() error {
	return c.RunQuery("\\l")
}

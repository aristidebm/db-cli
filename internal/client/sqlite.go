package client

import (
	"fmt"
	"os/exec"

	"example.com/db/internal/shutil"
)

type SQLite struct {
	*Client
}

func (c *SQLite) SetClient(client *Client) {
	c.Client = client
}

func (c *SQLite) Ping() error {
	args := []string{}
	args = append(args, c.Path, "SELECT 1;")
	cmd := exec.Command("sqlite3", args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *SQLite) Connect() error {
	args := []string{}
	args = append(args, c.Path)
	cmd := exec.Command("sqlite3", args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *SQLite) RunQuery(query string) error {
	if query == "" {
		return c.Connect()
	}

	args := []string{}
	switch c.Format {
	case JSON:
		args = append(args, "--json")
	case CSV:
		args = append(args, "--csv")
	case MARKDOWN:
		args = append(args, "--markdown")
	case HTML:
		args = append(args, "--html")
	default:
		return fmt.Errorf("%w: driver %s", UnsupportedFormat, c.Driver)
	}

	args = append(args, c.Path, query)
	cmd := exec.Command("sqlite3", args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *SQLite) ListTables() error {
	return c.RunQuery(".tables")
}

func (c *SQLite) ListDatabases() error {
	return fmt.Errorf("%w driver %s", UnsupportedCommand, c.Driver)
}

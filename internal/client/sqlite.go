package client

import (
	"fmt"
	"io"
	"os/exec"

	"example.com/db/internal/shutil"
	"example.com/db/internal/types"
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

func (c *SQLite) Connect() error {
	args := []string{}
	args = append(args, c.Path)

	prog := c.Client.GetInteractiveREPL()
	if prog == "" {
		prog = "sqlite3"
	}
	cmd := exec.Command(prog, args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *SQLite) RunQuery(query string, args ...string) error {
	switch c.GetFormat() {
	case types.JSON:
		args = append(args, "--json")
	case types.CSV:
		args = append(args, "--csv")
	case types.MARKDOWN, "markdown":
		args = append(args, "--markdown")
	case types.HTML:
		args = append(args, "--html")
	case "":
		// nothing to-do
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

func (c *SQLite) String() string {
	return c.URL
}

package client

import (
	"fmt"
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
	cmd := exec.Command(c.getPingCommand(), args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *SQLite) Connect() error {
	args := []string{}
	args = append(args, c.Path)
	cmd := exec.Command(c.getConnectCommand(), args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *SQLite) RunQuery(query string) error {
	args := []string{}
	switch c.Format {
	case types.JSON:
		args = append(args, "--json")
	case types.CSV:
		args = append(args, "--csv")
	case types.MARKDOWN:
		args = append(args, "--markdown")
	case types.HTML:
		args = append(args, "--html")
	default:
		return fmt.Errorf("%w: driver %s", UnsupportedFormat, c.Driver)
	}

	args = append(args, c.Path, query)
	cmd := exec.Command(c.getQueryCommand(), args...)
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

func (c *SQLite) getPingCommand() string {
	if c.Client.SourceConfig.Ping != "" {
		return c.Client.SourceConfig.Ping
	}
	if c.Client.DriverConfig.Ping != "" {
		return c.Client.DriverConfig.Ping
	}
	return "sqlite3"
}

func (c *SQLite) getConnectCommand() string {
	if c.Client.SourceConfig.Connect != "" {
		return c.Client.SourceConfig.Connect
	}
	if c.Client.DriverConfig.Connect != "" {
		return c.Client.DriverConfig.Connect
	}
	return "sqlite3"
}

func (c *SQLite) getQueryCommand() string {
	if c.Client.SourceConfig.Query != "" {
		return c.Client.SourceConfig.Query
	}
	if c.Client.DriverConfig.Query != "" {
		return c.Client.DriverConfig.Query
	}
	return "sqlite3"
}

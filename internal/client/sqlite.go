package client

import (
	"example.com/db/internal/shutil"
	"fmt"
)

type SQLite struct {
	client *Client
}

func (c *SQLite) SetClient(client *Client) {
	c.client = client
}

func (c *SQLite) Ping() error {
	return shutil.RunCommand("sqlite3", c.client.Path, "SELECT 1;")
}

func (c *SQLite) Connect() error {
	return shutil.RunInteractiveCommand("sqlite3", c.client.Path)
}

func (c *SQLite) RunQuery(query string) error {
	if query == "" {
		return c.Connect()
	}
	args := []string{c.client.Path}
	// if separator != "" {
	// 	args = append(args, "-separator", separator)
	// }
	args = append(args, query)
	return shutil.RunCommand("sqlite3", args...)
}

func (c *SQLite) ListTables() error {
	return c.RunQuery(".tables")
}

func (c *SQLite) ListDatabases() error {
	return fmt.Errorf("%w driver %s", UnsupportedCommand, c.client.Driver)
}

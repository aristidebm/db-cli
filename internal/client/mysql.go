package client

import (
	"os/exec"

	"example.com/db/internal/shutil"
)

type MySQL struct {
	*Client
}

func (c *MySQL) SetClient(client *Client) {
	c.Client = client
}

func (c *MySQL) Ping() error {
	args := []string{}
	args = append(args, "-e", "SELECT 1;")
	cmd := exec.Command("mysql", args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *MySQL) Connect() error {
	args := []string{}
	cmd := exec.Command("mysql", args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *MySQL) RunQuery(query string) error {
	args := []string{}
	args = append(args, "-e", query)

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

	cmd := exec.Command("mysql", args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *MySQL) ListTables() error {
	return c.RunQuery("SHOW TABLES;")
}

func (c *MySQL) ListDatabases() error {
	return c.RunQuery("SHOW DATABASES;")
}

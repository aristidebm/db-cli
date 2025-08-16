package client

import (
	"fmt"
	"io"
	"os/exec"

	"example.com/db/internal/shutil"
	"example.com/db/internal/types"
)

type Postgres struct {
	*Client
}

func (c *Postgres) SetClient(client *Client) {
	c.Client = client
}

func (c *Postgres) Ping() error {
	args := []string{c.URL}
	args = append(args, "-c", "SELECT 1;")
	cmd := exec.Command("psql", args...)
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

func (c *Postgres) Connect() error {
	args := []string{c.URL}
	prog := c.Client.GetInteractiveREPL()
	if prog == "" {
		prog = "psql"
	}
	cmd := exec.Command(prog, args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Postgres) RunQuery(query string, args ...string) error {
	// Disable the pager by default, the user can enable it if needed
	args = append(args, c.URL, "-c", query,
		"--pset", "pager=off",
		"--pset", "footer=off",
	)
	switch c.GetFormat() {
	case types.CSV:
		args = append(args, "--pset", "format=csv")
	case types.HTML:
		args = append(args, "--pset", "format=html")
	case types.LATEX:
		args = append(args, "--pset", "format=latex")
	case types.ASCIIDOC:
		args = append(args, "--pset", "format=asciidoc")
	case types.UNALIGNED:
		args = append(args, "--pset", "format=unaligned")
	case "":
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
	// \dt add some annoying headers
	// return c.RunQuery("\\dt")
	return c.RunQuery("SELECT schemaname AS \"Schema\", tablename AS \"Name\", 'table' AS \"Type\", tableowner AS \"Owner\" FROM pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema') ORDER BY schemaname, tablename;")
}

func (c *Postgres) String() string {
	return c.URL
}

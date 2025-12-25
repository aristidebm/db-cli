package client

import (
	"bytes"
	"example.com/db/internal/shutil"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"example.com/db/internal/types"
)

var commonTypeOperations = map[string][]string{
	"hash":        {"HGETALL", "HSET"},
	"set":         {"SMEMBERS", "SADD"},
	"zset":        {"ZRANGE", "ZADD"},
	"list":        {"LRANGE", "LPUSH"},
	"string":      {"GET", "SET"},
	"stream":      {"XRANGE", "XADD"},
	"hyperloglog": {"PFCOUNT", "PFADD"},
}

type Redis struct {
	*Client
}

func (c *Redis) SetClient(client *Client) {
	c.Client = client
}

func (c *Redis) Ping() error {
	args := []string{"-u", c.URL}
	args = append(args, "PING")
	cmd := exec.Command(c.getDefaultCommand(), args...)
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

func (c *Redis) Connect() error {
	args := []string{"-u", c.URL}
	prog := c.Client.GetInteractiveREPL()
	if prog == "" {
		prog = c.getDefaultCommand()
	}
	cmd := exec.Command(prog, args...)
	return shutil.RunInteractive(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Redis) RunQuery(query string, args ...string) error {
	defaultArgs := []string{}
	switch c.GetFormat() {
	case types.JSON:
		defaultArgs = append(defaultArgs, "--json")
	case types.CSV:
		defaultArgs = append(defaultArgs, "--csv")
	case "":
		// nothing to-do
	default:
		return fmt.Errorf("%w: scheme %s", UnsupportedFormat, c.Scheme)
	}

	args = append(args, defaultArgs...)
	args = append(args, "-u", c.URL)
	args = append(args, strings.Fields(query)...)
	cmd := exec.Command(c.getDefaultCommand(), args...)
	return shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(c.Stdout),
		shutil.WithStderr(c.Stderr),
	)
}

func (c *Redis) ListTables() error {
	args := []string{"-u", c.URL, "--scan", "--pattern", "*"}
	var buf bytes.Buffer

	cmd := exec.Command(c.getDefaultCommand(), args...)
	err := shutil.Run(cmd,
		shutil.WithStdin(c.Stdin),
		shutil.WithStdout(&buf),
		shutil.WithStderr(&buf),
	)

	if err != nil {
		return err
	}

	keys, err := io.ReadAll(&buf)
	if err != nil {
		return err
	}

	var typeBuffer bytes.Buffer
	var entryNumber = 0
	for key := range bytes.Lines(keys) {
		args = []string{"-u", c.URL, "TYPE", strings.TrimSpace(string(key))}
		cmd := exec.Command(c.getDefaultCommand(), args...)
		// fmt.Println(cmd)
		err := shutil.Run(cmd,
			shutil.WithStdin(c.Stdin),
			shutil.WithStdout(&typeBuffer),
			shutil.WithStderr(&typeBuffer),
		)
		if err != nil {
			return err
		}

		value, err := io.ReadAll(&typeBuffer)
		if err != nil {
			return err
		}
		typeStr := strings.TrimSpace(string(value))
		output := []string{strings.ToUpper(typeStr)}
		if v, ok := commonTypeOperations[typeStr]; ok {
			output = append(output, v...)
		} else {
			output = append(output, "", "")
		}
		output = append(output, string(key))

		// Decide the std-out to use
		stdOut := c.Stdout
		if stdOut == nil {
			stdOut = os.Stdout
		}

		// print the header
		if entryNumber == 0 {
			fmt.Fprintln(stdOut, strings.Join([]string{"TYPE", "ACCESS", "MUTATE", "KEYS"}, "\t\t"))
		}

		fmt.Fprint(stdOut, strings.Join(output, "\t\t"))

		entryNumber += 1
	}
	return nil
}

func (c *Redis) String() string {
	return c.URL
}

func (c *Redis) getDefaultCommand() string {
	if shutil.IsCommandInstalled("valkey-cli") {
		return "valkey-cli"
	}
	return "redis-cli"
}

package client

import (
	"io"
	"strings"

	"example.com/db/internal/shutil"
	"example.com/db/internal/types"
)

type DataSource interface {
	Ping() error
	Connect() error
	RunQuery(string) error
	ListTables() error
	ListDatabases() error
	SetClient(client *Client)
	String() string
}

type Client struct {
	URL          string
	Driver       string
	Host         string
	Port         string
	User         string
	Pass         string
	DBName       string
	Path         string
	Stdout       io.Writer
	Stdin        io.Reader
	Stderr       io.Writer
	Format       types.Format
	SourceConfig types.Source
	DriverConfig types.Driver
	DataSource
}

func NewClient(connectionURL string) (*Client, error) {
	client := &Client{URL: connectionURL}
	if err := client.parseURL(); err != nil {
		return nil, err
	}
	client.Format = types.DEFAULT
	return client, nil
}

func (c *Client) parseURL() error {
	u, err := shutil.ParseURL(c.URL)
	if err != nil {
		return err
	}
	c.Driver = u.Scheme
	c.Host = u.Hostname()
	c.Port = u.Port()
	if u.User != nil {
		c.User = u.User.Username()
		if pass, ok := u.User.Password(); ok {
			c.Pass = pass
		}
	}

	c.DBName = strings.TrimPrefix(u.Path, "/")

	// Handle special cases
	switch c.Driver {
	case "sqlite3":
		c.Path = u.Path
		c.DataSource = &SQLite{}
	case "redis":
		// Redis doesn't use traditional database names
		if c.Port == "" {
			c.Port = "6379"
		}
		c.DataSource = &Redis{}
	case "postgres", "postgresql":
		if c.Port == "" {
			c.Port = "5432"
		}
		c.DataSource = &Postgres{}
	case "mysql":
		if c.Port == "" {
			c.Port = "3306"
		}
		c.DataSource = &MySQL{}
	}

	if c.DataSource == nil {
		return shutil.URLParseError
	}

	c.DataSource.SetClient(c)
	return nil
}

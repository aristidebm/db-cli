package client

import (
	_ "fmt"
	"io"
	"net/url"
	"strings"
)

type DataSource interface {
	Ping() error
	Connect() error
	RunQuery(string) error
	ListTables() error
	ListDatabases() error
	SetClient(client *Client)
}

type Client struct {
	URL    string
	Driver string
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
	Path   string
	Stdout io.Writer
	Stdin  io.Reader
	Stderr io.Writer
	DataSource
}

func NewClient(connectionURL string) *Client {
	client := &Client{URL: connectionURL}
	client.parseURL()
	return client
}

func (c *Client) parseURL() error {
	u, err := url.Parse(c.URL)
	if err != nil {
		// Handle non-standard URLs like sqlite3:////path
		if strings.HasPrefix(c.URL, "sqlite3:") {
			c.Driver = "sqlite3"
			c.Path = strings.TrimPrefix(c.URL, "sqlite3:")
			// Remove leading slashes for absolute paths
			c.Path = strings.TrimLeft(c.Path, "/")
			if !strings.HasPrefix(c.Path, "/") {
				c.Path = "/" + c.Path
			}
			return URLParseError
		}
		return URLParseError
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
		return URLParseError
	}

	c.DataSource.SetClient(c)
	return nil
}

// func (c *DatabaseClient) Ping() error {
// 	switch c.Driver {
// 	case "postgres", "postgresql":
// 		return c.pingPostgreSQL()
// 	case "mysql":
// 		return c.pingMySQL()
// 	case "sqlite3":
// 		return c.pingSQLite()
// 	case "redis":
// 		return c.pingRedis()
// 	default:
// 		return fmt.Errorf("unsupported driver: %s", c.Driver)
// 	}
// }
//
// func (c *DatabaseClient) Connect() error {
// 	switch c.Driver {
// 	case "postgres", "postgresql":
// 		return c.connectPostgreSQL()
// 	case "mysql":
// 		return c.connectMySQL()
// 	case "sqlite3":
// 		return c.connectSQLite()
// 	case "redis":
// 		return c.connectRedis()
// 	default:
// 		return fmt.Errorf("unsupported driver: %s", c.Driver)
// 	}
// }
//
// func (c *DatabaseClient) RunQuery(query, separator string) error {
// 	if query == "" {
// 		return c.Connect()
// 	}
//
// 	switch c.Driver {
// 	case "postgres", "postgresql":
// 		return c.runPostgreSQLQuery(query, separator)
// 	case "mysql":
// 		return c.runMySQLQuery(query, separator)
// 	case "sqlite3":
// 		return c.runSQLiteQuery(query, separator)
// 	case "redis":
// 		return c.runRedisCommand(query)
// 	default:
// 		return fmt.Errorf("unsupported driver: %s", c.Driver)
// 	}
// }
//
// func (c *DatabaseClient) ListTables() error {
// 	switch c.Driver {
// 	case "postgres", "postgresql":
// 		return c.runPostgreSQLQuery("\\dt", "")
// 	case "mysql":
// 		return c.runMySQLQuery("SHOW TABLES;", "")
// 	case "sqlite3":
// 		return c.runSQLiteQuery(".tables", "")
// 	case "redis":
// 		return c.runRedisCommand("KEYS *")
// 	default:
// 		return fmt.Errorf("unsupported driver: %s", c.Driver)
// 	}
// }
//
// func (c *DatabaseClient) ListDatabases() error {
// 	switch c.Driver {
// 	case "postgres", "postgresql":
// 		return c.runPostgreSQLQuery("\\l", "")
// 	case "mysql":
// 		return c.runMySQLQuery("SHOW DATABASES;", "")
// 	case "sqlite3":
// 		return fmt.Errorf("sqlite3 doesn't support multiple databases")
// 	case "redis":
// 		return c.runRedisCommand("CONFIG GET databases")
// 	default:
// 		return fmt.Errorf("unsupported driver: %s", c.Driver)
// 	}
// }
//
// // PostgreSQL methods
// func (c *DatabaseClient) pingPostgreSQL() error {
// 	args := c.buildPSQLArgs()
// 	args = append(args, "-c", "SELECT 1;")
// 	cmd := exec.Command("psql", args...)
// 	return cmd.Run()
// }
//
// func (c *DatabaseClient) connectPostgreSQL() error {
// 	args := c.buildPSQLArgs()
// 	cmd := exec.Command("psql", args...)
// 	cmd.Stdin = nil
// 	cmd.Stdout = nil
// 	cmd.Stderr = nil
// 	return runInteractiveCommand("psql", args...)
// }
//
// func (c *DatabaseClient) runPostgreSQLQuery(query, separator string) error {
// 	args := c.buildPSQLArgs()
// 	if separator != "" {
// 		args = append(args, "-F", separator)
// 	}
// 	args = append(args, "-c", query)
// 	return runCommand("psql", args...)
// }
//
// func (c *DatabaseClient) buildPSQLArgs() []string {
// 	var args []string
//
// 	if c.Host != "" {
// 		args = append(args, "-h", c.Host)
// 	}
// 	if c.Port != "" {
// 		args = append(args, "-p", c.Port)
// 	}
// 	if c.User != "" {
// 		args = append(args, "-U", c.User)
// 	}
// 	if c.DBName != "" {
// 		args = append(args, "-d", c.DBName)
// 	}
//
// 	return args
// }
//
// // MySQL methods
// func (c *DatabaseClient) pingMySQL() error {
// 	args := c.buildMySQLArgs()
// 	args = append(args, "-e", "SELECT 1;")
// 	cmd := exec.Command("mysql", args...)
// 	return cmd.Run()
// }
//
// func (c *DatabaseClient) connectMySQL() error {
// 	args := c.buildMySQLArgs()
// 	return runInteractiveCommand("mysql", args...)
// }
//
// func (c *DatabaseClient) runMySQLQuery(query, separator string) error {
// 	args := c.buildMySQLArgs()
// 	if separator != "" {
// 		// MySQL doesn't have a direct separator flag, but we can format output
// 		args = append(args, "--batch", "--raw")
// 	}
// 	args = append(args, "-e", query)
// 	return runCommand("mysql", args...)
// }
//
// func (c *DatabaseClient) buildMySQLArgs() []string {
// 	var args []string
//
// 	if c.Host != "" {
// 		args = append(args, "-h", c.Host)
// 	}
// 	if c.Port != "" {
// 		args = append(args, "-P", c.Port)
// 	}
// 	if c.User != "" {
// 		args = append(args, "-u", c.User)
// 	}
// 	if c.Pass != "" {
// 		args = append(args, "-p"+c.Pass)
// 	}
// 	if c.DBName != "" {
// 		args = append(args, c.DBName)
// 	}
//
// 	return args
// }
//
// // SQLite methods
// func (c *DatabaseClient) pingSQLite() error {
// 	cmd := exec.Command("sqlite3", c.Path, "SELECT 1;")
// 	return cmd.Run()
// }
//
// func (c *DatabaseClient) connectSQLite() error {
// 	return runInteractiveCommand("sqlite3", c.Path)
// }
//
// func (c *DatabaseClient) runSQLiteQuery(query, separator string) error {
// 	args := []string{c.Path}
// 	if separator != "" {
// 		args = append(args, "-separator", separator)
// 	}
// 	args = append(args, query)
// 	return runCommand("sqlite3", args...)
// }
//
// // Redis methods
// func (c *DatabaseClient) pingRedis() error {
// 	args := c.buildRedisArgs()
// 	args = append(args, "PING")
// 	cmd := exec.Command("redis-cli", args...)
// 	return cmd.Run()
// }
//
// func (c *DatabaseClient) connectRedis() error {
// 	args := c.buildRedisArgs()
// 	return runInteractiveCommand("redis-cli", args...)
// }
//
// func (c *DatabaseClient) runRedisCommand(command string) error {
// 	args := c.buildRedisArgs()
// 	// Split command into parts
// 	parts := strings.Fields(command)
// 	args = append(args, parts...)
// 	return runCommand("redis-cli", args...)
// }
//
// func (c *DatabaseClient) buildRedisArgs() []string {
// 	var args []string
//
// 	if c.Host != "" {
// 		args = append(args, "-h", c.Host)
// 	}
// 	if c.Port != "" {
// 		args = append(args, "-p", c.Port)
// 	}
// 	if c.Pass != "" {
// 		args = append(args, "-a", c.Pass)
// 	}
//
// 	return args
// }

package client

type Redis struct {
	client *Client
}

func (c *Redis) Ping() error                 {}
func (c *Redis) Connect() error              {}
func (c *Redis) RunQuery(query string) error {}
func (c *Redis) ListTables() error
func (c *Redis) ListDatabases() error     {}
func (c *Redis) SetClient(client *Client) {}

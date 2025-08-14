package client

type Postgres struct {
	client *Client
}

func (c *Postgres) Ping() error                 {}
func (c *Postgres) Connect() error              {}
func (c *Postgres) RunQuery(query string) error {}
func (c *Postgres) ListTables() error
func (c *Postgres) ListDatabases() error     {}
func (c *Postgres) SetClient(client *Client) {}

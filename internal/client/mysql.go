package client

type MySQL struct {
	client *Client
}

func (c *MySQL) Ping() error           {}
func (c *MySQL) Connect() error        {}
func (c *MySQL) RunQuery(string) error {}
func (c *MySQL) ListTables() error
func (c *MySQL) ListDatabases() error     {}
func (c *MySQL) SetClient(client *Client) {}

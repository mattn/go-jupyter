package jupyter

type Config struct {
	Token  string
	Origin string
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{config: config}
}

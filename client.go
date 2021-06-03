package jupyter

import (
	"github.com/google/uuid"
)

type Config struct {
	Token   string
	Origin  string
	Session string
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	if config.Session == "" {
		config.Session = uuid.NewString()
	}
	return &Client{config: config}
}

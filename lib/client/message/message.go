package message

import (
	"github.com/snowmerak/rmq/gen/proto/message"
	"google.golang.org/grpc"
)

type Client struct {
	target string
	message.MessageQueueClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		target:             conn.Target(),
		MessageQueueClient: message.NewMessageQueueClient(conn),
	}
}

func (c *Client) Target() string {
	return c.target
}

package message

import (
	"context"
	"github.com/snowmerak/rmq/gen/bean"
	"github.com/snowmerak/rmq/gen/proto/message"
)

type Server struct {
	message.UnimplementedMessageQueueServer
	bean *bean.Bean
}

func New(bean *bean.Bean) *Server {
	return &Server{
		bean: bean,
	}
}

func (s Server) Send(ctx context.Context, msg *message.SendMsg) (*message.ReplyMsg, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) mustEmbedUnimplementedMessageQueueServer() {
}

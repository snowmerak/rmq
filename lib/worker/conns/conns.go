package conns

import (
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

//go:bean
type GrpcConns struct {
	connections map[string]*grpc.ClientConn
	rwLock      sync.RWMutex
	logger      *zerolog.Logger
}

func New(logger *zerolog.Logger) *GrpcConns {
	return &GrpcConns{
		connections: map[string]*grpc.ClientConn{},
		logger:      logger,
	}
}

func (c *GrpcConns) Add(conn *grpc.ClientConn) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	c.connections[conn.Target()] = conn
}

func (c *GrpcConns) Remove(target ...string) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	for _, t := range target {
		delete(c.connections, t)
	}
}

func (c *GrpcConns) Get(target ...string) []*grpc.ClientConn {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	if len(target) == 0 {
		conns := make([]*grpc.ClientConn, 0, len(c.connections))
		for _, conn := range c.connections {
			conns = append(conns, conn)
		}
		return conns
	}

	conns := make([]*grpc.ClientConn, 0, len(target))
	for _, t := range target {
		conn, ok := c.connections[t]
		switch ok {
		case true:
			conns = append(conns, conn)
		case false:
			conn, err := grpc.Dial(t, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				c.logger.Warn().Err(err).Str("target", t).Msg("dial failed")
				continue
			}
			conns = append(conns, conn)
			c.connections[t] = conn
		}
	}
	return conns
}

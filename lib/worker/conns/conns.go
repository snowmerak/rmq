package conns

import (
	"google.golang.org/grpc"
	"sync"
)

//go:bean
type GrpcConns struct {
	connections map[string]*grpc.ClientConn
	rwLock      sync.RWMutex
}

func New() *GrpcConns {
	return &GrpcConns{
		connections: map[string]*grpc.ClientConn{},
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
		if conn, ok := c.connections[t]; ok {
			conns = append(conns, conn)
		}
	}
	return conns
}

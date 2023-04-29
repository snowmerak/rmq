package station

import (
	"context"
	"github.com/snowmerak/rmq/gen/bean"
	"github.com/snowmerak/rmq/gen/proto/message"
	client "github.com/snowmerak/rmq/lib/client/message"
	"github.com/snowmerak/rmq/lib/worker/bridge"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
)

type Station struct {
	in            *bridge.Bridge[[]byte]
	clients       []*atomic.Pointer[client.Client]
	clientsMap    map[string]int
	nilIndices    []int
	clientsRwLock sync.RWMutex

	cancel context.CancelFunc

	bean *bean.Bean
}

func New(queueSize uint64, bean *bean.Bean, conn ...*grpc.ClientConn) *Station {
	in := bridge.New[[]byte](queueSize)
	clients := make([]*atomic.Pointer[client.Client], len(conn))
	for i, c := range conn {
		clients[i].Store(client.New(c))
	}

	ctx, cancel := context.WithCancel(context.Background())

	st := &Station{
		in:      in,
		clients: clients,
		cancel:  cancel,

		clientsMap: map[string]int{},
		nilIndices: []int{},

		bean: bean,
	}

	go func() {
		for {
			data, err := in.Pop()
			if err != nil {
				bean.Logger().Error().Err(err).Msg("station pop from 'in' bridge error")
				break
			}

			if err := bean.Pool().Submit(func() {
				for i := range clients {
					cli := clients[i].Load()
					if cli == nil {
						continue
					}
					reply, err := cli.Send(ctx, &message.SendMsg{
						Data: data,
					})
					if err != nil {
						bean.Logger().Error().Err(err).Str("target", cli.Target()).Msg("send to server error")
						return
					}

					if !reply.Success {
						bean.Logger().Error().Str("target", cli.Target()).Bool("reply", reply.Success).Msg("failed to send to server")
					}
				}
			}); err != nil {
				bean.Logger().Error().Err(err).Msg("submit to pool error")
				return
			}
		}
	}()

	return st
}

func (s *Station) Push(data []byte) error {
	return s.in.Push(data)
}

func (s *Station) Close() {
	s.cancel()
}

func (s *Station) AddClient(c ...*grpc.ClientConn) {
	s.clientsRwLock.Lock()
	defer s.clientsRwLock.Unlock()

	for _, conn := range c {
		if _, ok := s.clientsMap[conn.Target()]; ok {
			s.bean.Logger().Warn().Str("target", conn.Target()).Msg("client already exists")
			continue
		}
		if len(s.nilIndices) > 0 {
			index := s.nilIndices[len(s.nilIndices)-1]
			s.nilIndices = s.nilIndices[:len(s.nilIndices)-1]
			s.clientsMap[conn.Target()] = index
			s.clients[index].Store(client.New(conn))
			continue
		}
		s.clientsMap[conn.Target()] = len(s.clients)
		s.clients = append(s.clients, &atomic.Pointer[client.Client]{})
		s.clients[len(s.clients)-1].Store(client.New(conn))
	}
}

func (s *Station) RemoveClient(t ...string) {
	for _, target := range t {
		if _, ok := s.clientsMap[target]; !ok {
			s.bean.Logger().Warn().Str("target", target).Msg("client not exists")
			continue
		}
		index := s.clientsMap[target]
		delete(s.clientsMap, target)
		s.clients[index].Store(nil)
		s.nilIndices = append(s.nilIndices, index)
	}
}

const (
	CheckTypeAdd = iota
	CheckTypeRem
)

func (s *Station) Check(t ...string) (int, []string) {
	s.clientsRwLock.RLock()
	defer s.clientsRwLock.RUnlock()

	result := []string(nil)
	typ := CheckTypeAdd
	if len(t) < len(s.clientsMap) {
		typ = CheckTypeRem
	}

	switch typ {
	case CheckTypeAdd:
		for _, target := range t {
			if _, ok := s.clientsMap[target]; ok {
				continue
			}
			result = append(result, target)
		}
	case CheckTypeRem:
		tempMap := make(map[string]struct{}, len(t))
		for _, target := range t {
			tempMap[target] = struct{}{}
		}
		for target := range s.clientsMap {
			if _, ok := tempMap[target]; ok {
				continue
			}
			result = append(result, target)
		}
	}

	return typ, result
}

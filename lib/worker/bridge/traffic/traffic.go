package traffic

import (
	"context"
	"github.com/snowmerak/rmq/gen/bean"
	"github.com/snowmerak/rmq/lib/worker/bridge"
	"github.com/snowmerak/rmq/lib/worker/bridge/traffic/station"
)

//go:bean
type Traffic struct {
	in  *bridge.Bridge[Signal]
	out map[string]*station.Station
}

func New(queueSize uint64, bean *bean.Bean) *Traffic {
	in := bridge.New[Signal](queueSize)
	out := map[string]*station.Station{}

	go func() {
		for {
			signal, err := in.Pop()
			if err != nil {
				bean.Logger().Error().Err(err).Msg("pop from 'in' bridge error")
				break
			}

			list, isCacheHit, err := bean.Subset().GetMembers(context.Background(), signal.Direction)
			if err != nil {
				bean.Logger().Error().Err(err).Str("direction", signal.Direction).Msg("cannot get targets")
				break
			}

			st, ok := out[signal.Direction]
			if !ok {
				st = station.New(queueSize, bean)
				out[signal.Direction] = st
			}

			if !isCacheHit {
				typ, targets := st.Check(list...)
				switch typ {
				case station.CheckTypeAdd:
					st.AddClient(bean.GrpcConns().Get(targets...)...)
				case station.CheckTypeRem:
					st.RemoveClient(targets...)
				}
			}

			if err := st.Push(signal.Msg); err != nil {
				bean.Logger().Err(err).Msg("cannot send message to station")
			}
		}
	}()

	return &Traffic{
		in:  in,
		out: out,
	}
}

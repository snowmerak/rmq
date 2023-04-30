package bean

import (
	"github.com/snowmerak/rmq/lib/client/subset"
	"github.com/snowmerak/rmq/lib/worker/conns"
	"github.com/snowmerak/rmq/lib/worker/logger"
	"github.com/snowmerak/rmq/lib/worker/pool"
)

type Bean struct {
	subset    *subset.Subset
	grpcconns *conns.GrpcConns
	logger    *logger.Logger
	pool      *pool.Pool
}

type Builder struct {
	bean *Bean
}

func New() *Builder {
	return &Builder{new(Bean)}
}

func (b *Builder) Build() *Bean {
	return b.bean
}

func (b *Builder) AddSubset(subset *subset.Subset) *Builder {
	b.bean.subset = subset
	return b
}

func (b *Builder) AddGrpcConns(grpcconns *conns.GrpcConns) *Builder {
	b.bean.grpcconns = grpcconns
	return b
}

func (b *Builder) AddLogger(logger *logger.Logger) *Builder {
	b.bean.logger = logger
	return b
}

func (b *Builder) AddPool(pool *pool.Pool) *Builder {
	b.bean.pool = pool
	return b
}

func (b *Bean) Subset() *subset.Subset {
	return b.subset
}

func (b *Bean) GrpcConns() *conns.GrpcConns {
	return b.grpcconns
}

func (b *Bean) Logger() *logger.Logger {
	return b.logger
}

func (b *Bean) Pool() *pool.Pool {
	return b.pool
}

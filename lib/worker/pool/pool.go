package pool

import "github.com/panjf2000/ants"

//go:bean
type Pool struct {
	pool *ants.Pool
}

func New(size int, opts ...ants.Option) (*Pool, error) {
	p, err := ants.NewPool(size, opts...)
	if err != nil {
		return nil, err
	}

	return &Pool{
		pool: p,
	}, nil
}

func (p *Pool) Submit(task func()) error {
	return p.pool.Submit(task)
}

func (p *Pool) Release() {
	p.pool.Release()
}

package main

type Resetable interface {
	Reset()
}

type Pool[T Resetable] struct {
	pool chan T
}

func NewPool[T Resetable]() *Pool[T] {
	return &Pool[T]{
		pool: make(chan T),
	}
}

func (p *Pool[T]) Get() T {
	return <-p.pool
}

func (p *Pool[T]) Put(obj T) {
	p.pool <- obj
}

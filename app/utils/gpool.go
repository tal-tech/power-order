package utils

import "sync"

// 协程池
type Gpool struct {
	queue chan int
	wg *sync.WaitGroup
}

func NewGPool(size int) *Gpool {
	if size < 1 {
		size = 10
	}
	pool := new(Gpool)
	pool.queue = make(chan int , size)
	pool.wg = &sync.WaitGroup{}
	return pool
}

func (p *Gpool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

func (p *Gpool) Done() {
	<-p.queue
	p.wg.Done()
}

func (p *Gpool) Wait() {
	p.wg.Wait()
}

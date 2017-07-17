package main

import (
	"sync"
)

type Iterator interface {
	Next() (interface{}, bool)
}

type Observable interface {
	Iterator
	Subscribe(Observer) (<-chan error, <-chan struct{})
}

type Observer func(interface{})

type RxQueue chan interface{}

func NewQueue(emittable ...func() (interface{})) RxQueue {
	rxq := make(RxQueue)
	var wg sync.WaitGroup
	for _, em := range emittable {
		wg.Add(1)
		go func(em func() interface{}) {
			rxq <- em()
			wg.Done()
		}(em)
	}
	go func() {
		wg.Wait()
		close(rxq)
	}()
	return rxq
}

func (rxq RxQueue) Next() (interface{}, bool) {
	if next, ok := <-rxq; ok {
		return next, false
	}
	return nil, true
}
func (rxq RxQueue) Subscribe(observe Observer) (<-chan error, <-chan struct{}) {
	errc := make(chan error)
	donec := make(chan struct{})

	go func() {
		defer close(donec)
		defer close(errc)
		for {
			next, done := rxq.Next()
			if done {
				donec <- struct{}{}
				return
			}
			if err, ok := next.(error); ok {
				errc <- err
				return
			}
			observe(next)
		}
	}()
	return errc, donec
}


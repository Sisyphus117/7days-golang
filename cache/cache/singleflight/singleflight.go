package singleflight

import "sync"

type call struct {
	wg    sync.WaitGroup
	value any
	err   error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func NewGroup() *Group {
	return &Group{m: make(map[string]*call)}
}

func (g *Group) Call(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()

	if c, has := g.m[key]; has {
		c.wg.Wait()
		return c.value, c.err
	}
	c := NewCall()
	c.wg.Add(1)
	c.value, c.err = fn()
	c.wg.Done()

	g.mu.Unlock()
	return c.value, c.err

}

func NewCall() *call {
	return &call{}
}

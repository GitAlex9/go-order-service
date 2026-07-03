package id

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type CounterGenerator struct {
	mutex    sync.RWMutex
	counters map[EntityPrefix]*atomic.Uint64
}

func NewCounterGenerator() *CounterGenerator {
	return &CounterGenerator{
		counters: make(map[EntityPrefix]*atomic.Uint64),
	}
}

func (g *CounterGenerator) Generate(prefix EntityPrefix) string {
	counter := g.getCounter(prefix)
	nextID := counter.Add(1)
	return fmt.Sprintf("%s-%03d", prefix, nextID)
}

func (g *CounterGenerator) getCounter(prefix EntityPrefix) *atomic.Uint64 {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	counter, exists := g.counters[prefix]

	if !exists {
		counter = &atomic.Uint64{}
		g.counters[prefix] = counter
	}
	return counter
}

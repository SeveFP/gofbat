package gofbat

import (
	"sync"
	"time"
)

type Batcher interface {
	Batch(key string, f func())
}

type element struct {
	f      func()
	nCalls int
	timer  *time.Timer
}

type batcher struct {
	elements        map[string]*element
	triggerCallAt   time.Duration
	triggerAtNCalls int
	timedOnly       bool
	sync.Mutex
}

func NewBatcher(triggerCallAt time.Duration, triggerAtNcalls int, timedOnly bool) Batcher {
	return &batcher{
		triggerCallAt:   triggerCallAt,
		triggerAtNCalls: triggerAtNcalls,
		timedOnly:       timedOnly,
		elements:        make(map[string]*element),
	}
}

func (b *batcher) Batch(key string, f func()) {
	b.Lock()
	defer b.Unlock()

	e, ok := b.elements[key]
	if !ok {
		e = &element{
			f:      f,
			nCalls: 1,
			timer:  b.newTimer(key),
		}

		b.elements[key] = e

		return
	}

	e.nCalls++
	e.timer.Stop()
	e.timer = b.newTimer(key)
	b.elements[key] = e

	if !b.timedOnly && e.nCalls >= b.triggerAtNCalls {
		go b.callAndDelete(key)
	}
}

func (b *batcher) callAndDelete(key string) {
	b.Lock()
	defer b.Unlock()

	e, ok := b.elements[key]
	if !ok {
		return
	}

	e.f()

	delete(b.elements, key)
}

func (b *batcher) newTimer(key string) *time.Timer {
	return time.AfterFunc(b.triggerCallAt, func() {
		b.callAndDelete(key)
	})
}

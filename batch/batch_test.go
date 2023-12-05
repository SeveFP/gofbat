package gofbat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatcher(t *testing.T) {
	t.Run("triggerAtNCalls", func(t *testing.T) {
		b := NewBatcher(time.Hour, 3, false)

		var called int
		done := make(chan struct{})
		f := func() {
			called++
			done <- struct{}{}
		}

		b.Batch("key", f)
		b.Batch("key", f)
		b.Batch("key", f)

		<-done
		assert.Equal(t, 1, called)
	})

	t.Run("triggerCallAt", func(t *testing.T) {
		b := NewBatcher(time.Second, 3, false)

		var called int
		done := make(chan struct{})
		f := func() {
			called++
			done <- struct{}{}
		}

		b.Batch("key", f)

		<-done
		assert.Equal(t, 1, called)
	})

	t.Run("onlyTimed", func(t *testing.T) {
		b := NewBatcher(time.Second*5, 1, true)

		var called int
		done := make(chan struct{})
		f := func() {
			called++
			done <- struct{}{}
		}

		b.Batch("key", f)
		b.Batch("key", f)
		b.Batch("key", f)
		b.Batch("key", f)
		b.Batch("key", f)
		b.Batch("key", f)

		<-done
		assert.Equal(t, 1, called)
	})

}

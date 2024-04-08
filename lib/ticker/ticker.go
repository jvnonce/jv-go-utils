package ticker

import "time"

type FinishFunc[T comparable] func(id T)

type ticker[T comparable] struct {
	id       T
	finishAt time.Time
	onFinish FinishFunc[T]
	ticker   *time.Ticker
}

// Ticker interface
type Ticker[T comparable] interface {
	// Starts ticker
	Start()
	// Reset ticker for new time to finish
	Reset(finishAt time.Time)
}

// Ticker constructor
func New[T comparable](id T, finishAt time.Time, onFinish FinishFunc[T]) Ticker[T] {
	return &ticker[T]{
		id:       id,
		finishAt: finishAt,
		onFinish: onFinish,
	}
}

func (t *ticker[T]) Start() {
	go func() {
		now := time.Now()
		if t.finishAt.Before(now) {
			go t.onFinish(t.id)
			return
		}
		d := t.finishAt.Sub(now)
		t.ticker = time.NewTicker(d)
		defer func() {
			t.ticker.Stop()
		}()
		for {
			<-t.ticker.C
			t.onFinish(t.id)
			return
		}
	}()
}

func (t *ticker[T]) Reset(finishAt time.Time) {
	now := time.Now()
	t.finishAt = finishAt
	if t.finishAt.Before(now) {
		go t.onFinish(t.id)
		return
	}
	if t.ticker != nil {
		d := t.finishAt.Sub(now)
		t.ticker.Reset(d)
	} else {
		t.Start()
	}
}

package sigtrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/x-mod/event"
)

type Handler func()

type Capture struct {
	notify  chan os.Signal
	traps   map[string]Handler
	close   chan struct{}
	serving *event.Event
	stopped *event.Event
}

type CaptureOpt func(*Capture)

func Trap(sig os.Signal, handler Handler) CaptureOpt {
	return func(ca *Capture) {
		if handler != nil {
			ca.traps[sig.String()] = handler
		}
	}
}

func New(opts ...CaptureOpt) *Capture {
	ca := &Capture{
		traps:   make(map[string]Handler),
		close:   make(chan struct{}),
		serving: event.New(),
		stopped: event.New(),
	}
	for _, o := range opts {
		o(ca)
	}
	ca.notify = make(chan os.Signal, 1)
	return ca
}

func (cap *Capture) Serve(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("context required")
	}
	if len(cap.traps) == 0 {
		return nil
	}

	signal.Notify(cap.notify)
	defer cap.stopped.Fire()
	cap.serving.Fire()

	for {
		select {
		case <-cap.close:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-cap.notify:
			cap.fire(sig)
		}
	}
}

func (cap *Capture) fire(sig os.Signal) {
	if handler, ok := cap.traps[sig.String()]; ok {
		handler()
	}
}
func (cap *Capture) Serving() <-chan struct{} {
	return cap.serving.Done()
}
func (cap *Capture) Close() <-chan struct{} {
	if cap.serving.HasFired() {
		close(cap.close)
		return cap.stopped.Done()
	}
	return event.Done()
}

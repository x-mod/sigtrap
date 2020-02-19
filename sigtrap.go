package sigtrap

import (
	"context"
	"os"
	"os/signal"
)

type Handler func()

type Capture struct {
	notify chan os.Signal
	traps  map[string]Handler
	stop   chan struct{}
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
		traps: make(map[string]Handler),
	}
	for _, o := range opts {
		o(ca)
	}
	ca.notify = make(chan os.Signal, 1)
	return ca
}

func (ca *Capture) Serve(ctx context.Context) error {
	if len(ca.traps) == 0 {
		return nil
	}
	ca.stop = make(chan struct{})
	signal.Notify(ca.notify)
	for {
		select {
		case <-ca.stop:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-ca.notify:
			ca.fire(sig)
		}
	}
}

func (ca *Capture) fire(sig os.Signal) {
	if handler, ok := ca.traps[sig.String()]; ok {
		handler()
	}
}

func (ca *Capture) Close() {
	if ca.stop != nil {
		close(ca.stop)
	}
}

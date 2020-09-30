package main

import (
	"context"
	"log"
	"syscall"

	"github.com/x-mod/sigtrap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := sigtrap.New(
		sigtrap.Trap(syscall.SIGINT, sigtrap.Handler(func() {
			log.Println("sig INT catched")
			cancel()
		})),
		sigtrap.Trap(syscall.SIGTERM, sigtrap.Handler(cancel)),
	)
	defer c.Close()

	go c.Serve(ctx)
	<-c.Serving()
	log.Println("sigtrap is serv ...")
	<-ctx.Done()
	log.Println("sigtrap done ... ", ctx.Err())
}

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
		sigtrap.Trap(syscall.SIGINT, sigtrap.Handler(cancel)),
		sigtrap.Trap(syscall.SIGTERM, sigtrap.Handler(cancel)),
	)
	defer c.Close()
	log.Println("sigtrap: waiting ...")
	log.Println("sigtrap:", c.Serve(ctx))
}

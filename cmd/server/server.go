package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-tire/pkg/config"
	"github.com/golang-tire/pkg/log"
)

var (
	debugMode = flag.Bool("debug", false, "run in debug mode")
)

func cliContext() context.Context {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT}
	var sig = make(chan os.Signal, len(signals))
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(sig, signals...)
	go func() {
		<-sig
		cancel()
	}()
	return ctx
}

func main() {
	flag.Parse()

	ctx := cliContext()
	err := log.Init(ctx, *debugMode)
	if err != nil {
		panic(err)
	}

	err = config.Init("config", "yaml", "")
	if err != nil {
		panic(err)
	}

	err = setupModules(ctx)
	if err != nil {
		panic(err)
	}
}

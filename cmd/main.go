package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arimatakao/encrypted-messages/cmd/config"
	"github.com/arimatakao/encrypted-messages/internal/server"
)

var cfgpath *string
var l *log.Logger

func init() {
	cfgpath = flag.String("config", "./config.yaml", "config path")
	flag.Parse()

	l = log.New(os.Stdout, "INFO ", log.LstdFlags)
	l.Printf("inited logger")
}

func main() {
	if err := config.LoadConfig(*cfgpath); err != nil {
		l.Fatalf("error while load configuration: %s", err.Error())
	}

	app, err := server.NewServer(l)
	if err != nil {
		l.Fatalf("server cannot be inited: %s", err.Error())
	}

	go func() {
		if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatalf("get fatal internal error: %s", err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	l.Printf("starting graceful shutdown")
	if err := app.Shutdown(context.Background()); err != nil {
		l.Fatalf("shutdown with error: %s", err.Error())
	}
	l.Printf("shutdown complete successful")
	os.Exit(0)
}

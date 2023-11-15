package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/arimatakao/encrypted-messages/cmd/config"
	"github.com/arimatakao/encrypted-messages/internal/server"
)

var cfgpath *string

func init() {
	cfgpath = flag.String("config", "./config.yaml", "config path")
	flag.Parse()
}

func main() {
	if err := config.LoadConfig(*cfgpath); err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}
	fmt.Printf("Port: %s\nBaseUrl: %s\n", config.Conf.App.Port, config.Conf.App.BaseUrl)
	app := server.NewServer(config.Conf)
	log.Fatal(app.Run())
}

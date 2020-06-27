package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	flagConfigFile = flag.String("config", "abserver.json", "config file path")
)

const (
	ABSERVER_DB_NAME = "abserver"
	ROLE_ADMIN       = 1
	ROLE_USER        = 2
)

func main() {
	flag.Parse()

	var err error
	config, err := NewConfig(*flagConfigFile)
	if err != nil {
		panic(err)
	}

	StartWebServer(config)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

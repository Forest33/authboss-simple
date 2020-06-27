package main

import (
	"flag"
	"github.com/aarondl/tpl"
	"github.com/volatiletech/authboss-clientstate"
	_ "github.com/volatiletech/authboss/auth"
	_ "github.com/volatiletech/authboss/logout"
	_ "github.com/volatiletech/authboss/recover"
	_ "github.com/volatiletech/authboss/register"
	"os"
	"os/signal"
)

var (
	storer *Storer
	ws     *WebServer

	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer

	templates tpl.Templates
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

	storer = NewStorer(config)
	ws = StartWebServer(config)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

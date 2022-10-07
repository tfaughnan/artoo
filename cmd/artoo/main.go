package main

import (
	"log"

	"github.com/tfaughnan/artoo/client"
	"github.com/tfaughnan/artoo/config"
	"github.com/tfaughnan/artoo/plugin/echo"
	"github.com/tfaughnan/artoo/plugin/help"
	"github.com/tfaughnan/artoo/plugin/openai"
	"github.com/tfaughnan/artoo/plugin/raw"
	"github.com/tfaughnan/artoo/plugin/tmdb"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	c := client.NewClient(cfg)

	c.RegisterLineHandler(`^:(?P<server>\S+) 001 (?P<nick>\S+) :(?P<body>.*)$`, c.Handle001)
	c.RegisterLineHandler(`^PING (?P<token>\S+)$`, c.HandlePing)
	c.RegisterLineHandler(`^:(?P<nick>\S+)!(?P<user>\S+)@(?P<host>\S+) PRIVMSG (?P<target>\S+) :(?P<body>.*)$`, c.HandlePrivmsg)

	c.RegisterPlugin(echo.Plugin)
	c.RegisterPlugin(openai.Plugin)
	c.RegisterPlugin(raw.Plugin)
	c.RegisterPlugin(tmdb.Plugin)
	c.RegisterPlugin(help.Plugin)

	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}

	log.Fatal(c.MainLoop())
}

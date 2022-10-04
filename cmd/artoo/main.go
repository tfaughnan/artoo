package main

import (
	"log"

	"github.com/tfaughnan/artoo/client"
	"github.com/tfaughnan/artoo/config"
	"github.com/tfaughnan/artoo/plugin/echo"
	"github.com/tfaughnan/artoo/plugin/openai"
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
	c.RegisterPluginHandler(echo.EchoPattern, echo.EchoHandler)
	c.RegisterPluginHandler(openai.OpenaiPattern, openai.OpenaiHandler)

	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}

	log.Fatal(c.MainLoop())
}

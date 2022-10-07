package bonk

import (
	"regexp"
	"strings"

	"github.com/tfaughnan/artoo/client"
)

var pattern = regexp.MustCompile(`^\.bonk\s+(?P<nick>.+)$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "bonk",
	Desc:    "Bonk a user",
	Usage:   ".bonk <nick>",
	Example: ".bonk deckard",
}

func handler(c *client.Client, lgroups, bgroups map[string]string) {
	c.PrintfPrivmsg(lgroups["target"], "\U0001F6A8 BONK! \U0001F528 %s, go to horny jail!", strings.TrimSpace(bgroups["nick"]))
}

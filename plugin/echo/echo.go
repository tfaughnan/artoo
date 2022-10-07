package echo

import (
	"regexp"

	"github.com/tfaughnan/artoo/client"
)

var pattern = regexp.MustCompile(`^\.echo\s+(?P<text>.+)$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "echo",
	Desc:    "Echo text back.",
	Usage:   ".echo <text>",
}

func handler(c *client.Client, lgroups, bgroups map[string]string) {
	c.PrintfPrivmsg(lgroups["target"], bgroups["text"])
}

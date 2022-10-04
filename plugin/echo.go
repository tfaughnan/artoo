package plugin

import "github.com/tfaughnan/artoo/client"

var EchoPattern = `^\.echo\s+(?P<text>.+)$`

func EchoHandler(c *client.Client, lgroups, bgroups map[string]string) {
	c.PrintfPrivmsg(lgroups["target"], bgroups["text"])
}

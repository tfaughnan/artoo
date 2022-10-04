package echo

import "github.com/tfaughnan/artoo/client"

var Pattern = `^\.echo\s+(?P<text>.+)$`

func Handler(c *client.Client, lgroups, bgroups map[string]string) {
	c.PrintfPrivmsg(lgroups["target"], bgroups["text"])
}

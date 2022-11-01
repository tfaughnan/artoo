package raw

import (
	"log"
	"regexp"

	"github.com/tfaughnan/artoo/client"
)

var pattern = regexp.MustCompile(`^\.raw\s+(?P<cmd>.+)$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "raw",
	Desc:    "Send a raw command to the IRC connection socket (owner-only).",
	Usage:   ".raw <command>",
	Example: ".raw JOIN #bots",
}

func handler(c *client.Client, lgroups, bgroups map[string]string) {
	user := lgroups["user"]
	cmd := bgroups["cmd"]
	if user != c.Cfg.Owner {
		log.Printf("Unauthorized raw command from <%s>: %s\n", user, cmd)
		c.PrintfPrivmsg(lgroups["target"], "Unauthorized")
	} else {
		c.PrintRaw(cmd)
	}
}

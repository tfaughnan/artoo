package raw

import (
	"log"

	"github.com/tfaughnan/artoo/client"
)

var Pattern = `^\.raw\s+(?P<cmd>.+)$`

func Handler(c *client.Client, lgroups, bgroups map[string]string) {
	nick := lgroups["nick"]
	cmd := bgroups["cmd"]
	if nick != c.Cfg.Owner {
		log.Printf("Unauthorized raw command: <%s> %s\n", nick, cmd)
		c.PrintfPrivmsg(lgroups["target"], "Unauthorized")
	} else {
		c.PrintRaw(cmd)
	}
}

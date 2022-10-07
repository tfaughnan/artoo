package help

import (
	"regexp"
	"strings"

	"github.com/tfaughnan/artoo/client"
)

// TODO: make <plugin> optional in regex
var pattern = regexp.MustCompile(`^\.help(\s+(?P<plugin>.+))?$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "help",
	Desc:    "Print artoo's help text.",
	Usage:   ".help || .help <plugin>",
	Example: ".help openai",
}

func handler(c *client.Client, lgroups, bgroups map[string]string) {
	target := lgroups["target"]
	plugin := bgroups["plugin"]

	if strings.TrimSpace(plugin) != "" {
		for _, p := range c.Plugins {
			if p.Name == plugin {
				c.PrintfPrivmsg(target, "%-12s : %s", "Name", p.Name)
				c.PrintfPrivmsg(target, "%-12s : %s", "Description", p.Desc)
				c.PrintfPrivmsg(target, "%-12s : %s", "Usage", p.Usage)
				c.PrintfPrivmsg(target, "%-12s : %s", "Example", p.Example)
				return
			}
		}
	}

	plist := ""
	for i, p := range c.Plugins {
		if i > 0 {
			plist += ", "
		}
		plist += p.Name
	}
	c.PrintfPrivmsg(target, "Available plugins: %s\n", plist)
	c.PrintfPrivmsg(target, "Use \".help <plugin>\" for more info.")
}

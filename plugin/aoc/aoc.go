package aoc

import (
	"log"
	"regexp"
	"sort"
	"strconv"

	"github.com/tfaughnan/artoo/client"
	"github.com/tfaughnan/artoo/style"
)

var pattern = regexp.MustCompile(`^\.aoc(\s+(?P<n>.+))?$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "aoc",
	Desc:    "List top n users on an Advent of Code private leaderboard (default n=5)",
	Usage:   ".aoc [n]",
	Example: ".aoc 20",
}

func handler(c *client.Client, lgroups, bgroups map[string]string) {
	target := lgroups["target"]
	n, err := strconv.Atoi(bgroups["n"])
	if err != nil {
		n = 5
	}

	r, err := fetchLeaderboard(c.Cfg.Aoc, c.Cfg.HttpTimeout)
	if err != nil {
		log.Println(err)
		c.PrintfPrivmsg(target, "%s%s%v%s", style.Color, style.Red, err, style.Reset)
		return
	}

	members := make([]member, 0, len(r.Members))
	for _, m := range r.Members {
		members = append(members, m)
	}
	sort.Slice(members, func(i, j int) bool {
		return members[i].LocalScore > members[j].LocalScore
	})

	if n > len(members) {
		n = len(members)
	}

	c.PrintfPrivmsg(target, "%sAdvent of Code %s \U0001F384 Private Leaderboard 181542%s", style.Bold, r.Event, style.Reset)
	c.PrintfPrivmsg(target, "%s%-6s %-6s %-25s %s%s", style.Bold, "Rank", "Score", "Stars", "Name", style.Reset)
	for i := range members[:n] {
		starString := ""
		for d := 1; d <= 25; d++ {
			if parts, ok := members[i].Days[strconv.Itoa(d)]; ok {
				if len(parts) == 2 {
					starString += "*"
				} else {
					starString += "-"
				}
			} else {
				starString += " "
			}
		}
		c.PrintfPrivmsg(target, "%-6d %-6d %s%s%-25s%s %s", i+1, members[i].LocalScore, style.Color, style.Yellow, starString, style.Reset, members[i].Name)
	}
}

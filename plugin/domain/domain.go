package domain

import (
	"log"
	"regexp"

	"github.com/tfaughnan/artoo/client"
	"github.com/tfaughnan/artoo/style"
)

var pattern = regexp.MustCompile(`^\.domain\s+(?P<query>.+)$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "domain",
	Desc:    "Check domain availability via Gandi",
	Usage:   ".domain <query>",
	Example: ".domain tyrellcorp.biz",
}

func handler(c *client.Client, lgroups, bgroups map[string]string) {
	query := bgroups["query"]
	target := lgroups["target"]

	dr, err := fetchDomain(c.Cfg.Domain, c.Cfg.HttpTimeout, query)
	if err != nil {
		log.Println(err)
		c.PrintfPrivmsg(target, "%s%s%v%s", style.Color, style.Red, err, style.Reset)
		return
	}

	if len(dr.Products) == 0 {
		c.PrintfPrivmsg(target, "%s: %s%sinvalid%s (no results)", query, style.Color, style.Red, style.Reset)
	} else if dr.Products[0].Status != "available" {
		c.PrintfPrivmsg(target, "%s: %s%sunavailable%s", query, style.Color, style.Red, style.Reset)
	} else {
		c.PrintfPrivmsg(target, "%s: available for %s%s%s %.2f%s", query, style.Color, style.Green, c.Cfg.Domain.Currency, dr.Products[0].Prices[0].PriceBeforeTaxes, style.Reset)
	}
}

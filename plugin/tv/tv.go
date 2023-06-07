package tv

import (
	"fmt"
	"log"
	"regexp"

	"github.com/tfaughnan/artoo/client"
	"github.com/tfaughnan/artoo/config"
	"github.com/tfaughnan/artoo/style"
)

var pattern = regexp.MustCompile(`^\.tv\s+(?P<query>.+)$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "tv",
	Desc:    "Search for TV shows on TMDB",
	Usage:   ".tv <query>",
	Example: ".tv godfather part ii",
}

type movie struct {
	ID          int
	Title       string
	Year        string
	Director    string
	Rating      float32
	RatingColor string
	URL         string
	Overview    string
}

func handler(c *client.Client, lgroups, bgroups map[string]string) {
	query := bgroups["query"]
	target := lgroups["target"]

	m, err := fetchMovie(c.Cfg.Tmdb, c.Cfg.HttpTimeout, query)
	if err != nil {
		log.Println(err)
		c.PrintfPrivmsg(target, "%v", err)
		return
	} else if m == (movie{}) {
		c.PrintfPrivmsg(target, "No results")
		return
	}

	c.PrintfPrivmsg(target, "%s%s%s%s (%s)%s  %s%s[%.1f/10]%s @ %s",
		style.Color, style.Blue, style.Bold, m.Title, m.Year, style.Reset, style.Color, m.RatingColor, m.Rating, style.Reset,
		m.URL)
	c.PrintfPrivmsg(target, "%sOverview:%s %s", style.Bold, style.Reset, m.Overview)
}

func fetchMovie(cfg config.TmdbConfig, timeout int, query string) (movie, error) {
	sresp, err := fetchSearch(cfg, timeout, query)
	if err != nil {
		return movie{}, err
	} else if sresp.TotalResults == 0 {
		return movie{}, nil
	}

	res := sresp.Results[0]
	fmt.Printf("%#v", res)
	m := movie{}
	m.ID = res.ID
	m.Title = res.OrigTitle
	m.Year = res.ReleaseDate[0:4]
	m.Rating = res.VoteAverage
	m.URL = fmt.Sprintf("https://www.themoviedb.org/tv/%d", res.ID)

	switch {
	case res.VoteAverage < 5:
		m.RatingColor = style.Red
	case res.VoteAverage < 7.5:
		m.RatingColor = style.Orange
	default:
		m.RatingColor = style.Green
	}

	cresp, err := fetchCredits(cfg, timeout, m.ID)
	if err != nil {
		return movie{}, err
	}

	for _, person := range cresp.Crew {
		if person.Job == "Director" {
			m.Director = person.Name
			break
		}
	}

	m.Overview = res.Overview

	return m, nil
}

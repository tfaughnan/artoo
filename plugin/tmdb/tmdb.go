package tmdb

import (
	"fmt"
	"log"
	"regexp"

	"github.com/tfaughnan/artoo/client"
	"github.com/tfaughnan/artoo/config"
	"github.com/tfaughnan/artoo/style"
)

var pattern = regexp.MustCompile(`^\.kino\s+(?P<query>.+)$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "tmdb",
	Desc:    "Search for movies on TMDB",
	Usage:   ".kino <query>",
	Example: ".kino godfather part ii",
}

type movie struct {
	ID          int
	Title       string
	Year        string
	Director    string
	Rating      float32
	RatingColor string
	URL         string
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

	var ratingText string
	if m.Rating > 0 {
		ratingText = fmt.Sprintf("%.1f", m.Rating)
	} else {
		ratingText = "-"
	}

	c.PrintfPrivmsg(target, "%s%s%s%s (%s)%s dir. %s %s%s[%s/10]%s @ %s",
		style.Color, style.Blue, style.Bold, m.Title, m.Year, style.Reset,
		m.Director, style.Color, m.RatingColor, ratingText, style.Reset,
		m.URL)
}

func fetchMovie(cfg config.TmdbConfig, timeout int, query string) (movie, error) {
	sresp, err := fetchSearch(cfg, timeout, query)
	if err != nil {
		return movie{}, err
	} else if sresp.TotalResults == 0 {
		return movie{}, nil
	}

	res := sresp.Results[0]
	m := movie{}
	m.ID = res.ID
	m.Title = res.Title
	m.Year = res.ReleaseDate[0:4]
	m.Rating = res.VoteAverage
	m.URL = fmt.Sprintf("https://www.themoviedb.org/movie/%d", res.ID)

	switch {
	case res.VoteCount == 0:
		m.RatingColor = style.Grey
		m.Rating = -1
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

	return m, nil
}

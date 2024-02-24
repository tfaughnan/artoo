package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strings"

	"github.com/tfaughnan/artoo/config"
)

var IrcMaxBytes = 400

type LineHandler func(lgroups map[string]string)

type Client struct {
	cfgPath      string // XXX: filepath type instead?
	Cfg          config.Config
	conn         net.Conn
	r            *textproto.Reader
	w            *textproto.Writer
	lineHandlers map[*regexp.Regexp]LineHandler // XXX
	Plugins      []Plugin
}

func NewClient(cfg config.Config) *Client {
	c := Client{}
	c.Cfg = cfg
	c.lineHandlers = make(map[*regexp.Regexp]LineHandler)
	c.Plugins = make([]Plugin, 0)
	return &c
}

func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.Cfg.Host, c.Cfg.Port)
	conn, err := dial(addr, c.Cfg.SSL)
	if err != nil {
		return err
	}

	c.r = textproto.NewReader(bufio.NewReader(conn))
	c.w = textproto.NewWriter(bufio.NewWriter(conn))

	if c.Cfg.Pass != "" {
		c.w.PrintfLine("PASS %s:%s", c.Cfg.Nick, c.Cfg.Pass)
	}
	c.w.PrintfLine("NICK %s", c.Cfg.Nick)
	c.w.PrintfLine("USER %s 0 * %s", c.Cfg.User, c.Cfg.Real)

	return nil
}

func (c *Client) MainLoop() error {
	for {
		line, err := c.r.ReadLine()
		if err != nil {
			return err
		}

		if c.Cfg.Verbose {
			log.Printf("> %s\n", line)
		}

		for re := range c.lineHandlers {
			if groups := matchGroups(re, line); groups != nil {
				c.lineHandlers[re](groups)
			}
		}
	}
}

func (c *Client) RegisterLineHandler(pattern string, fn LineHandler) {
	re := regexp.MustCompile(pattern)
	c.lineHandlers[re] = fn
}

func (c *Client) RegisterPlugin(p Plugin) {
	c.Plugins = append(c.Plugins, p)
}

func (c *Client) Handle001(lgroups map[string]string) {
	for _, channel := range c.Cfg.Channels {
		c.w.PrintfLine("JOIN %s", channel)
	}
}

func (c *Client) HandlePing(lgroups map[string]string) {
	c.w.PrintfLine("PONG %s", lgroups["token"])
}

func (c *Client) HandlePrivmsg(lgroups map[string]string) {
	if lgroups["target"] == c.Cfg.Nick {
		lgroups["target"] = lgroups["nick"] // for direct messages
	}

	for _, p := range c.Plugins {
		if bgroups := matchGroups(p.Pattern, lgroups["body"]); bgroups != nil {
			p.Handler(c, lgroups, bgroups)
		}
	}
}

func (c *Client) PrintfPrivmsg(target string, format string, args ...any) {
	maxBytes := IrcMaxBytes - len("PRIVMSG "+target+" :"+"\r\n")
	for _, line := range strings.Split(fmt.Sprintf(format, args...), "\n") {
		msg := ""
		length := len(line)
		for i, r := range line {
			msg += string(r)
			if (i+1)%maxBytes == 0 || i == length-1 {
				c.w.PrintfLine("PRIVMSG %s :%s", target, msg)
				msg = ""
			}
		}
	}
}

func (c *Client) PrintRaw(cmd string) {
	c.w.PrintfLine(cmd)
}

func dial(addr string, ssl bool) (net.Conn, error) {
	if ssl {
		return tls.Dial("tcp", addr, nil)
	} else {
		return net.Dial("tcp", addr)
	}
}

func matchGroups(re *regexp.Regexp, s string) map[string]string {
	groups := make(map[string]string)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			groups[name] = match[i]
		}
	}

	return groups
}

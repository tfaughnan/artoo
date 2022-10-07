package client

import "regexp"

type PluginHandler func(c *Client, lgroups, bgroups map[string]string)
type Plugin struct {
	Pattern *regexp.Regexp
	Handler PluginHandler
	Name    string
	Desc    string
	Usage   string
	Example string
}

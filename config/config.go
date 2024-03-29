package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Host        string   `toml:"host"`
	Port        int      `toml:"port"`
	SSL         bool     `toml:"ssl"`
	Nick        string   `toml:"nick"`
	User        string   `toml:"user"`
	Real        string   `toml:"real"`
	Pass        string   `toml:"pass"`
	Channels    []string `toml:"channels"`
	Owner       string   `toml:"owner"`
	Verbose     bool     `toml:"verbose"`
	HttpTimeout int      `toml:"http_timeout"`

	// plugin-specific configuration
	Openai OpenaiConfig `toml:"openai"`
	Tmdb   TmdbConfig   `toml:"tmdb"`
	Domain DomainConfig `toml:"domain"`
}

type OpenaiConfig struct {
	ApiURL      string  `toml:"api_url"`
	ApiKey      string  `toml:"api_key"`
	Model       string  `toml:"model"`
	MaxTokens   int     `toml:"max_tokens"`
	Temperature float32 `toml:"temperature"`
}

type TmdbConfig struct {
	ApiURL string `toml:"api_url"`
	ApiKey string `toml:"api_key"`
}

type DomainConfig struct {
	ApiURL   string `toml:"api_url"`
	ApiKey   string `toml:"api_key"`
	Currency string `toml:"currency"`
}

func LoadConfig() (Config, error) {
	path := getConfigPath()
	log.Printf("Loading configuration from %s\n", path)
	return loadConfigFromFile(path)
}

func loadConfigFromFile(path string) (Config, error) {
	var cfg Config
	md, err := toml.DecodeFile(path, &cfg)
	setDefaults(&cfg, md)
	return cfg, err
}

func getConfigPath() string {
	var pathArg string
	flag.StringVar(&pathArg, "c", "", "config file")
	flag.Parse()
	if pathArg != "" {
		return pathArg
	}

	pathSys := "/etc/artoo.toml"
	home, err := os.UserHomeDir()
	if err != nil {
		return pathSys
	}
	pathHome := filepath.Join(home, ".config/artoo.toml")
	if _, err := os.Stat(pathHome); errors.Is(err, os.ErrNotExist) {
		return pathSys
	}
	return pathHome
}

func setDefaults(cfg *Config, md toml.MetaData) {
	if !md.IsDefined("host") {
		cfg.Host = "127.0.0.1"
	}
	if !md.IsDefined("port") {
		cfg.Port = 6667
	}
	if !md.IsDefined("ssl") {
		cfg.SSL = false
	}
	if !md.IsDefined("nick") {
		cfg.Nick = "artoo"
	}
	if !md.IsDefined("user") {
		cfg.User = "artoo"
	}
	if !md.IsDefined("real") {
		cfg.Real = "https://github.com/tfaughnan/artoo"
	}
	if !md.IsDefined("pass") {
		cfg.Pass = ""
	}
	if !md.IsDefined("channels") {
		cfg.Channels = []string{}
	}
	// TODO: default owner?
	if !md.IsDefined("verbose") {
		cfg.Verbose = false
	}
	if !md.IsDefined("http_timeout") {
		cfg.HttpTimeout = 60
	}
}

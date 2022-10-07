package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tfaughnan/artoo/client"
	"github.com/tfaughnan/artoo/config"
	"github.com/tfaughnan/artoo/style"
)

var pattern = regexp.MustCompile(`^\.prompt\s+(?P<prompt>.+)$`)
var Plugin = client.Plugin{
	Pattern: pattern,
	Handler: handler,
	Name:    "openai",
	Desc:    "Fetch a completion for the prompt from OpenAI.",
	Usage:   ".prompt <prompt>",
	Example: ".prompt Write a haiku about IRC",
}

func handler(c *client.Client, lgroups, bgroups map[string]string) {
	prompt := bgroups["prompt"]
	target := lgroups["target"]

	comp, err := fetchCompletion(c.Cfg.Openai, c.Cfg.HttpTimeout, prompt)
	if err != nil {
		log.Println(err)
		c.PrintfPrivmsg(target, "Request failed: %v", err)
		return
	}

	sep := " "
	if strings.Contains(comp, "\n") {
		// if multiline completion, print prompt on own line for aesthetics
		sep = "\n"
	}
	promptReminder := fmt.Sprintf("%s%s[%s%s%s]:%s", style.Color, style.Blue,
		style.Italics, prompt, style.Italics, style.Reset)
	c.PrintfPrivmsg(target, "%s%s%s", promptReminder, sep, comp)
}

type reqPayload struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
}

type respPayload struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created uint64   `json:"created"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage"`
}

type choice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
	// LogProbs     LogProbs    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type usage struct{}

func fetchCompletion(cfg config.OpenaiConfig, timeout int, prompt string) (string, error) {
	payload := reqPayload{prompt, cfg.Model, cfg.MaxTokens, cfg.Temperature}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	url := cfg.ApiURL + "/completions"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payloadJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.ApiKey))

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("Received status \"%s\"", resp.Status))
		return "", err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body) // TODO: change name
	if err != nil {
		return "", err
	}

	var results respPayload // TODO: change names
	if err := json.Unmarshal(bodyBytes, &results); err != nil {
		return "", err
	}

	choice := results.Choices[0]
	comp := strings.TrimSpace(choice.Text)
	if choice.FinishReason == "length" {
		comp += "..."
	}

	return comp, nil
}

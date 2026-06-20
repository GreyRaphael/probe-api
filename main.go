package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var version = "0.3.0"

type probeResult struct {
	label string
	code  int
}

type probeTarget struct {
	path    string
	body    string
	label   string
	headers map[string]string
}

func probe(url, body string, headers map[string]string) int {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}

func buildOpenAITargets(model string) []probeTarget {
	rBody, _ := json.Marshal(map[string]any{"model": model, "input": "hi"})
	cBody, _ := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{{"role": "user", "content": "hi"}},
	})

	authPrefix := "Bearer "

	return []probeTarget{
		{
			path:  "/responses",
			body:  string(rBody),
			label: "Responses API",
			headers: map[string]string{
				"Authorization": authPrefix,
			},
		},
		{
			path:  "/chat/completions",
			body:  string(cBody),
			label: "Chat Completions API",
			headers: map[string]string{
				"Authorization": authPrefix,
			},
		},
	}
}

func buildAnthropicTargets(model string) []probeTarget {
	aBody, _ := json.Marshal(map[string]any{
		"model":      model,
		"max_tokens": 1,
		"messages":   []map[string]string{{"role": "user", "content": "hi"}},
	})

	return []probeTarget{
		{
			path:  "/messages",
			body:  string(aBody),
			label: "Anthropic Messages API",
			headers: map[string]string{
				"x-api-key":       "",
				"anthropic-version": "2023-06-01",
			},
		},
	}
}

func printResult(label string, code int) {
	switch {
	case code == 404 || code == 405 || code == 0:
		fmt.Printf("  %-30s -> HTTP %-3d  [X] not supported\n", label, code)
	default:
		fmt.Printf("  %-30s -> HTTP %-3d  [OK] endpoint exists\n", label, code)
	}
}

func main() {
	var (
		showVersion bool
		showHelp    bool
		model       string
		anthropic   bool
	)

	flag.BoolVar(&showVersion, "version", false, "Print version and exit")
	flag.BoolVar(&showVersion, "v", false, "Print version and exit (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Print usage and exit")
	flag.BoolVar(&showHelp, "h", false, "Print usage and exit (shorthand)")
	flag.StringVar(&model, "model", "test", "Model name to use in probe requests")
	flag.StringVar(&model, "m", "test", "Model name (shorthand)")
	flag.BoolVar(&anthropic, "anthropic", false, "Probe Anthropic Messages API instead of OpenAI")
	flag.BoolVar(&anthropic, "a", false, "Probe Anthropic Messages API (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "probe-api \u2014 Detect API protocol support\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  probe-api [flags] <base_url> <api_key>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  base_url    API base URL including path prefix\n")
		fmt.Fprintf(os.Stderr, "  api_key     API key / token for authentication\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # OpenAI-compatible (default)\n")
		fmt.Fprintf(os.Stderr, "  probe-api https://openrouter.ai/api/v1 sk-or-xxx\n")
		fmt.Fprintf(os.Stderr, "  probe-api -m mimo-v2.5-pro https://example.com/v1 tp-xxx\n\n")
		fmt.Fprintf(os.Stderr, "  # Anthropic-compatible\n")
		fmt.Fprintf(os.Stderr, "  probe-api -a https://api.deepseek.com/anthropic sk-xxx\n")
		fmt.Fprintf(os.Stderr, "  probe-api -a -m claude-sonnet-4-20250514 https://api.anthropic.com sk-xxx\n\n")
		fmt.Fprintf(os.Stderr, "  # Misc\n")
		fmt.Fprintf(os.Stderr, "  probe-api -v\n\n")
		fmt.Fprintf(os.Stderr, "Exit Codes:\n")
		fmt.Fprintf(os.Stderr, "  0   All probes completed\n")
		fmt.Fprintf(os.Stderr, "  1   Usage error or missing arguments\n")
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("probe-api %s\n", version)
		os.Exit(0)
	}

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	baseURL := args[0]
	apiKey := args[1]

	if model == "" {
		model = "test"
	}

	fmt.Printf("Probing: %s\n", baseURL)
	fmt.Println("---")

	var targets []probeTarget
	if anthropic {
		targets = buildAnthropicTargets(model)
	} else {
		targets = buildOpenAITargets(model)
	}

	for _, t := range targets {
		url := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(t.path, "/")
		// Inject API key into headers
		for k := range t.headers {
			if t.headers[k] == "" {
				t.headers[k] = apiKey
			} else {
				t.headers[k] = t.headers[k] + apiKey
			}
		}
		code := probe(url, t.body, t.headers)
		printResult(t.label, code)
	}

	fmt.Println("---")
	fmt.Println("404/405 = not supported | 200/400/422 = endpoint exists | 0 = connection failed")
}

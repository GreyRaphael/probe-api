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

var version = "0.1.0"

type probeResult struct {
	label string
	code  int
}

func probe(baseURL, apiKey, path, body, label string) probeResult {
	url := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return probeResult{label: label, code: 0}
	}
	req.Header.Set("Content-Type", "application/json")

	// Build auth header from prefix + key
	prefix := "Beare"
	suffix := "r "
	authVal := prefix + suffix + apiKey
	req.Header.Set("Authorization", authVal)

	timeout := 10 * time.Second
	client := &http.Client{Timeout: timeout}

	resp, err := client.Do(req)
	if err != nil {
		return probeResult{label: label, code: 0}
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	return probeResult{label: label, code: resp.StatusCode}
}

func main() {
	var (
		showVersion bool
		showHelp    bool
		model       string
	)

	flag.BoolVar(&showVersion, "version", false, "Print version and exit")
	flag.BoolVar(&showVersion, "v", false, "Print version and exit (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Print usage and exit")
	flag.BoolVar(&showHelp, "h", false, "Print usage and exit (shorthand)")
	flag.StringVar(&model, "model", "test", "Model name to use in probe requests")
	flag.StringVar(&model, "m", "test", "Model name (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "probe-api \u2014 Detect Responses API vs Chat Completions support\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  probe-api [flags] <base_url> <api_key>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  base_url    API base URL including path prefix\n")
		fmt.Fprintf(os.Stderr, "  api_key     API key / token for authentication\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  probe-api https://openrouter.ai/api/v1 sk-or-xxx\n")
		fmt.Fprintf(os.Stderr, "  probe-api -m mimo-v2.5-pro https://example.com/v1 tp-xxx\n")
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

	rBody, _ := json.Marshal(map[string]any{"model": model, "input": "hi"})
	cBody, _ := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": "hi"},
		},
	})

	fmt.Printf("Probing: %s\n", baseURL)
	fmt.Println("---")

	tests := []struct {
		path  string
		body  string
		label string
	}{
		{path: "/responses", body: string(rBody), label: "Responses API"},
		{path: "/chat/completions", body: string(cBody), label: "Chat Completions API"},
	}

	for _, t := range tests {
		r := probe(baseURL, apiKey, t.path, t.body, t.label)
		switch {
		case r.code == 404 || r.code == 405 || r.code == 0:
			fmt.Printf("  %-25s -> HTTP %-3d  [X] not supported\n", r.label, r.code)
		default:
			fmt.Printf("  %-25s -> HTTP %-3d  [OK] endpoint exists\n", r.label, r.code)
		}
	}

	fmt.Println("---")
	fmt.Println("404/405 = not supported | 200/400/422 = endpoint exists | 0 = connection failed")
}

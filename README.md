# probe-api

[中文文档](README-CN.md)

Detect whether an OpenAI-compatible API provider supports **Responses API** or **Chat Completions API** (or both).

Useful for configuring tools like [Codex CLI](https://github.com/openai/codex) that require a specific wire protocol — set `wire_api = "chat"` when the provider only supports Chat Completions.

## Install

```bash
# Download from GitHub Releases (recommended)
# https://github.com/GreyRaphael/probe-api/releases

# From source
go install github.com/GreyRaphael/probe-api@latest

# Or build locally
git clone https://github.com/GreyRaphael/probe-api.git
cd probe-api
make build
```

## Usage

```bash
probe-api [flags] <base_url> <api_key>
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--version` | `-v` | | Print version and exit |
| `--help` | `-h` | | Print usage and exit |
| `--model` | `-m` | `test` | Model name for probe requests |

### Examples

```bash
# OpenRouter (supports both APIs)
probe-api https://openrouter.ai/api/v1 sk-or-xxx

# Xiaomi MiMo (Chat Completions only)
probe-api -m mimo-v2.5-pro https://token-plan-cn.xiaomimimo.com/v1 tp-xxx

# OpenAI (supports both APIs)
probe-api https://api.openai.com/v1 sk-xxx

# Local Ollama (Chat Completions only, no API key needed)
probe-api http://localhost:11434/v1 ollama
```

### Output

```
Probing: https://openrouter.ai/api/v1
---
  Responses API             -> HTTP 400  [OK] endpoint exists
  Chat Completions API      -> HTTP 400  [OK] endpoint exists
---
404/405 = not supported | 200/400/422 = endpoint exists | 0 = connection failed
```

### HTTP Code Reference

| Code | Meaning |
|------|---------|
| 200 / 400 / 422 | Endpoint exists (400 = bad request / model not found) |
| 401 / 403 | Endpoint exists but auth failed (key invalid or missing) |
| 404 / 405 | Endpoint does not exist |
| 0 | Connection failed (timeout / DNS / network error) |

## Codex CLI Integration

If the provider only supports Chat Completions, add this to `~/.codex/config.toml`:

```toml
[model_providers.myprovider]
name = "My Provider"
base_url = "https://example.com/v1"
env_key = "MY_API_KEY"
wire_api = "chat"
requires_openai_auth = false
```

## Release

Push an annotated tag to trigger GitHub Actions — builds `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64` (tar.gz) and `windows/amd64` (zip). Each archive contains a single `probe-api` (or `probe-api.exe`) binary.

```bash
git tag -a v0.1.0 -m "probe-api v0.1.0: detect Responses API vs Chat Completions support"
git push origin v0.1.0
```

## Build

```bash
make build          # local binary
make package        # cross-compile + archive into dist/
make clean          # remove build artifacts
```

## License

MIT

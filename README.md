# probe-api

Detect whether an OpenAI-compatible API provider supports **Responses API** or **Chat Completions API** (or both).

Useful for configuring tools like Codex CLI that need to know which wire protocol a provider speaks.

## Install

```bash
# From source
go install github.com/GreyRaphael/probe-api@latest

# Or download from GitHub Releases
# https://github.com/GreyRaphael/probe-api/releases

# Or build locally
git clone https://github.com/GreyRaphael/probe-api.git
cd probe-api
make build
```

## Release

Push a tag to trigger a GitHub Actions release with pre-built binaries:

```bash
git tag -a v0.1.0 -m "probe-api v0.1.0: detect Responses API vs Chat Completions support"
git push origin v0.1.0
```

Binaries: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64` (tar.gz), `windows/amd64` (zip). Each archive contains a single `probe-api` (or `probe-api.exe`) binary.

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
# OpenRouter (supports both)
probe-api https://openrouter.ai/api/v1 sk-or-xxx

# Xiaomi MiMo (Chat Completions only)
probe-api -m mimo-v2.5-pro https://token-plan-cn.xiaomimimo.com/v1 tp-xxx

# OpenAI
probe-api https://api.openai.com/v1 sk-xxx
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

| HTTP Code | Meaning |
|-----------|---------|
| 200 / 400 / 422 | Endpoint exists (400 = bad request, model not found) |
| 401 / 403 | Endpoint exists but auth failed |
| 404 / 405 | Endpoint does not exist |
| 0 | Connection failed (timeout / DNS / network) |

## Cross-compile

```bash
make build-all   # outputs to dist/
```

## License

MIT

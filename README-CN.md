# probe-api

[English](README.md)

检测 OpenAI 兼容 API 是否支持 **Responses API** 或 **Chat Completions API**。

配置 [Codex CLI](https://github.com/openai/codex) 等工具时非常有用——如果 provider 只支持 Chat Completions，需要设置 `wire_api = "chat"`。

## 安装

```bash
# 从 GitHub Releases 下载（推荐）
# https://github.com/GreyRaphael/probe-api/releases

# 从源码安装
go install github.com/GreyRaphael/probe-api@latest

# 或者本地编译
git clone https://github.com/GreyRaphael/probe-api.git
cd probe-api
make build
```

## 用法

```bash
probe-api [flags] <base_url> <api_key>
```

### 参数

| 参数 | 缩写 | 默认值 | 说明 |
|------|------|--------|------|
| `--version` | `-v` | | 打印版本号 |
| `--help` | `-h` | | 打印帮助信息 |
| `--model` | `-m` | `test` | 探测请求使用的模型名 |

### 示例

```bash
# OpenRouter（两个 API 都支持）
probe-api https://openrouter.ai/api/v1 sk-or-xxx

# 小米 MiMo（仅支持 Chat Completions）
probe-api -m mimo-v2.5-pro https://token-plan-cn.xiaomimimo.com/v1 tp-xxx

# OpenAI（两个 API 都支持）
probe-api https://api.openai.com/v1 sk-xxx

# 本地 Ollama（仅 Chat Completions，不需要 API key）
probe-api http://localhost:11434/v1 ollama
```

### 输出示例

```
Probing: https://openrouter.ai/api/v1
---
  Responses API             -> HTTP 400  [OK] endpoint exists
  Chat Completions API      -> HTTP 400  [OK] endpoint exists
---
404/405 = not supported | 200/400/422 = endpoint exists | 0 = connection failed
```

### HTTP 状态码说明

| 状态码 | 含义 |
|--------|------|
| 200 / 400 / 422 | 端点存在（400 = 请求格式错误 / 模型名不对） |
| 401 / 403 | 端点存在但认证失败（key 无效或缺失） |
| 404 / 405 | 端点不存在 |
| 0 | 连接失败（超时 / DNS / 网络错误） |

## Codex CLI 集成

如果 provider 只支持 Chat Completions，在 `~/.codex/config.toml` 中添加：

```toml
[model_providers.myprovider]
name = "My Provider"
base_url = "https://example.com/v1"
env_key = "MY_API_KEY"
wire_api = "chat"
requires_openai_auth = false
```

## 发版

推送 annotated tag 触发 GitHub Actions 自动构建，产物包括 `linux/amd64`、`linux/arm64`、`darwin/amd64`、`darwin/arm64`（tar.gz）和 `windows/amd64`（zip）。每个压缩包内只有一个 `probe-api`（或 `probe-api.exe`）二进制文件。

```bash
git tag -a v0.1.0 -m "probe-api v0.1.0: detect Responses API vs Chat Completions support"
git push origin v0.1.0
```

## 编译

```bash
make build          # 本地编译
make package        # 交叉编译 + 打包到 dist/
make clean          # 清理构建产物
```

## 许可证

MIT

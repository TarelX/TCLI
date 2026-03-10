<div align="center">

```
 ████████╗  ██████╗ ██╗      ██╗
 ╚══██╔══╝ ██╔════╝ ██║      ██║
    ██║    ██║      ██║      ██║
    ██║    ██║      ██║      ██║
    ██║    ╚██████╗ ███████╗ ██║
    ╚═╝     ╚═════╝ ╚══════╝ ╚═╝
```

**Terminal Code Intelligence CLI**

面向中文开发者的 AI 代码辅助工具

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/TarelX/TCLI?style=flat-square)](https://github.com/TarelX/TCLI/releases)

</div>

---

## 功能特性

- **`tcli ask`** — 单次 AI 问答，支持管道输入和文件上下文
- **`tcli commit`** — AI 自动生成 git commit message，交互式确认
- **`tcli chat`** — 全屏交互式对话界面（TUI），流式响应，`@file` 文件引用
- **`tcli review`** — AI 代码审查，支持 staged/diff/文件
- **`tcli fix`** — 接收错误信息，AI 给出修复建议
- **`tcli explain`** — AI 解释代码文件或指定函数
- **多 Provider 支持** — OpenAI / Anthropic / 兼容接口（Deepseek、Qwen、Ollama 等）
- **中文优先** — 默认中文提示词和输出，专为中文开发者优化

## 安装

### 从 Release 下载（推荐）

前往 [Releases](https://github.com/TarelX/TCLI/releases) 下载对应平台的二进制文件。

### 从源码构建

```bash
git clone https://github.com/TarelX/TCLI.git
cd TCLI
go build -o tcli .
```

## 快速上手

### 1. 配置 API Key

```bash
# OpenAI 兼容接口（推荐，支持多种模型）
tcli config set default.provider openai
tcli config set default.model gpt-4o
tcli config set providers.openai.api_key sk-xxx
tcli config set providers.openai.base_url https://api.openai.com/v1

# 或 Anthropic
tcli config set default.provider anthropic
tcli config set default.model claude-3-7-sonnet
tcli config set providers.anthropic.api_key sk-ant-xxx
```

### 2. 开始使用

```bash
# 单次问答
tcli ask "用 Go 写一个快速排序"

# 带文件上下文
tcli ask "这个文件做什么" --file main.go

# 管道输入
cat error.log | tcli ask "这个报错什么原因"

# AI 生成 commit message
git add .
tcli commit

# 交互式对话
tcli chat

# 代码审查
tcli review main.go
tcli review --staged

# 错误修复
go build 2>&1 | tcli fix

# 代码解释
tcli explain main.go
tcli explain main.go --func ParseConfig
```

## 命令速查

| 命令 | 说明 | 示例 |
|------|------|------|
| `tcli ask [问题]` | 单次 AI 问答 | `tcli ask "什么是 goroutine"` |
| `tcli ask -f <file>` | 带文件上下文问答 | `tcli ask "解释一下" -f main.go` |
| `tcli commit` | AI 生成 commit message | `tcli commit --convention` |
| `tcli commit --push` | 生成并自动提交推送 | `tcli commit --push` |
| `tcli chat` | 交互式对话界面 | `tcli chat` |
| `tcli review <file>` | 代码审查 | `tcli review ./src/` |
| `tcli review --staged` | 审查 staged 文件 | `tcli review --staged` |
| `tcli fix` | 修复错误（管道输入） | `npm run build 2>&1 \| tcli fix` |
| `tcli fix --error "msg"` | 直接传入错误信息 | `tcli fix --error "undefined: foo" -f main.go` |
| `tcli explain <file>` | 解释代码 | `tcli explain internal/ai/client.go` |
| `tcli config set` | 设置配置项 | `tcli config set default.model gpt-4o` |
| `tcli config list` | 查看所有配置 | `tcli config list` |
| `tcli version` | 显示版本信息 | `tcli version` |

## Chat 交互模式

`tcli chat` 提供全屏交互式对话界面：

- **`@file`** — 输入 `@` 触发文件补全，将文件内容注入对话上下文
- **`/clear`** — 清空对话历史
- **`/copy`** — 复制最后一条 AI 回复到剪贴板
- **`/help`** — 查看帮助
- **Tab** — 补全候选项
- **↑↓** — 切换候选
- **Ctrl+C** — 退出

## 全局参数

| 参数 | 说明 |
|------|------|
| `--config <path>` | 指定配置文件路径 |
| `-p, --provider <name>` | 覆盖默认 provider |
| `-m, --model <name>` | 覆盖默认模型 |
| `--no-context` | 不注入代码上下文 |
| `--raw` | 纯文本输出（适合管道） |

## 配置文件

配置文件位于 `~/.tcli/config.yaml`（Windows: `%APPDATA%\tcli\config.yaml`）：

```yaml
version: 1
default:
  provider: openai
  model: gpt-4o
  max_tokens: 8192
  temperature: 0.7
providers:
  openai:
    api_key: sk-xxx
    base_url: https://api.openai.com/v1
  anthropic:
    api_key: sk-ant-xxx
  compatible:
    - name: deepseek
      api_key: sk-xxx
      base_url: https://api.deepseek.com/v1
      model: deepseek-coder
context:
  max_tokens: 32768
  include_git_info: true
  ignore_patterns:
    - "*.lock"
    - node_modules/
    - vendor/
ui:
  theme: auto
  language: zh
  show_token_count: true
```

## 技术栈

- **语言**: Go 1.25+
- **CLI 框架**: [cobra](https://github.com/spf13/cobra) + [viper](https://github.com/spf13/viper)
- **TUI**: [bubbletea v2](https://charm.land/bubbletea) + [bubbles v2](https://charm.land/bubbles) + [lipgloss v2](https://charm.land/lipgloss)
- **AI SDK**: [go-openai](https://github.com/sashabaranov/go-openai) + [anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go)
- **Markdown 渲染**: [glamour](https://github.com/charmbracelet/glamour)
- **发布**: [goreleaser](https://goreleaser.com/)

## 开发

```bash
# 克隆
git clone https://github.com/TarelX/TCLI.git
cd TCLI

# 安装依赖
go mod tidy

# 开发构建
go build -o tcli .

# 带版本信息构建
go build -ldflags "-X main.version=0.1.0 -X main.buildDate=$(date -u +%Y-%m-%d)" -o tcli .
```

## 开发者

**Ti**

## License

[MIT](LICENSE)

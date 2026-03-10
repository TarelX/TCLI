# TCli 任务完成清单

> 最后更新：2026-03-10

---

## 项目初始化

- [x] 创建工作目录 `f:\桌面\TCil\`
- [x] 编写 `开发方案.md`（技术栈、架构、功能设计、开发计划）
- [x] 创建 `go.mod`（模块名：`github.com/your-username/tcli`，Go 1.23）
- [x] 创建完整目录结构
  - [x] `cmd/`
  - [x] `internal/ai/`
  - [x] `internal/tui/splash/`
  - [x] `internal/tui/chat/`
  - [x] `internal/tui/components/`
  - [x] `internal/tui/styles/`
  - [x] `internal/context/`
  - [x] `internal/token/`
  - [x] `internal/prompt/`
  - [x] `internal/config/`
- [x] 创建 `.gitignore`
- [x] 创建 `.goreleaser.yaml`（多平台发布配置）

---

## Phase 1 — 骨架与核心命令（第 1-2 周）

### 骨架文件（已创建）

- [x] `main.go` — 程序入口，注入版本信息
- [x] `cmd/root.go` — Cobra 根命令，全局 flag，Viper 初始化
- [x] `cmd/version.go` — `tcli version` 命令
- [x] `cmd/chat.go` — `tcli chat` 命令骨架（Phase 2 实现）
- [x] `cmd/ask.go` — `tcli ask` 命令骨架
- [x] `cmd/review.go` — `tcli review` 命令骨架（Phase 3 实现）
- [x] `cmd/commit.go` — `tcli commit` 命令骨架
- [x] `cmd/fix.go` — `tcli fix` 命令骨架（Phase 3 实现）
- [x] `cmd/explain.go` — `tcli explain` 命令骨架（Phase 3 实现）
- [x] `cmd/config.go` — `tcli config set/list/providers` 命令
- [x] `internal/config/config.go` — 配置结构体定义
- [x] `internal/config/path.go` — 跨平台配置路径解析
- [x] `internal/ai/client.go` — AI 统一接口（Client interface）
- [x] `internal/ai/openai.go` — OpenAI 实现（流式 + 非流式）
- [x] `internal/ai/anthropic.go` — Anthropic 实现（流式 + 非流式）
- [x] `internal/ai/compatible.go` — OpenAI 兼容接口（Deepseek/Qwen/Ollama）
- [x] `internal/ai/factory.go` — 根据配置自动创建 AI 客户端
- [x] `internal/token/counter.go` — Token 数量估算
- [x] `internal/token/budget.go` — Token 预算管理 + 自动截断
- [x] `internal/context/gitinfo.go` — Git 信息收集
- [x] `internal/context/collector.go` — 代码上下文收集
- [x] `internal/context/fileread.go` — 文件读取 + 语言识别
- [x] `internal/prompt/templates.go` — 内置中文提示词模板
- [x] `internal/prompt/loader.go` — 自定义提示词加载
- [x] `internal/tui/styles/theme.go` — Lipgloss 色彩主题
- [x] `internal/tui/splash/model.go` — 启动动画（四阶段）
- [x] `internal/tui/chat/model.go` — Chat TUI Model 结构
- [x] `internal/tui/chat/update.go` — Chat Update 逻辑骨架
- [x] `internal/tui/chat/view.go` — Chat View 渲染骨架
- [x] `internal/tui/root.go` — RootModel（splash → chat 切换）

### 待实现（Phase 1 核心逻辑）

- [ ] 安装 Go（≥ 1.23）并执行 `go mod tidy` 拉取依赖
- [ ] `cmd/ask.go` — 完整实现
  - [ ] 读取管道输入（`os.Stdin`）
  - [ ] 读取 `--file` 参数文件内容
  - [ ] 调用 `ai.NewClientFromConfig()` 创建客户端
  - [ ] 流式输出到终端（`--raw` 模式）
  - [ ] Glamour Markdown 渲染（默认模式）
- [ ] `cmd/commit.go` — 完整实现
  - [ ] 调用 `context.GitStagedDiff()` 获取 diff
  - [ ] staged 为空时提示并退出
  - [ ] 调用 AI 生成 commit message
  - [ ] 交互式确认（y/n/edit）
  - [ ] `--push` flag：确认后执行 `git commit -m "xxx"`
- [ ] `cmd/config.go` — `providers` 子命令实现
- [ ] `internal/config/path.go` — `EnsureConfigDir()` 并写入默认配置

---

## Phase 2 — TUI 交互界面（第 3-4 周）

- [ ] 完善 `internal/tui/chat/update.go`
  - [ ] 流式响应通过 channel 注入 Bubbletea 消息循环
  - [ ] `/clear` 命令清空历史
  - [ ] `/copy` 命令复制最后回复到剪贴板
  - [ ] `/model` 命令切换模型
  - [ ] `/save` 命令保存对话到文件
- [ ] 完善 `internal/tui/chat/view.go`
  - [ ] 用 Glamour 渲染消息中的 Markdown / 代码块
  - [ ] Token 超过 80% 时状态栏变色警告
- [ ] `internal/tui/components/statusbar.go` — 独立状态栏组件
- [ ] `internal/tui/components/spinner.go` — 等待动画组件
- [ ] `internal/tui/components/input.go` — 多行输入框封装
- [ ] 完善 `cmd/chat.go` — 调用 `tui.StartChat()`
- [ ] Token 预算超限时自动截断历史，显示提示

---

## Phase 3 — 代码感知功能（第 5-6 周）

- [ ] `cmd/review.go` — 完整实现
  - [ ] `--staged` 模式
  - [ ] `--diff <ref>` 模式
  - [ ] 单文件 / 目录模式
  - [ ] 结构化审查报告输出
- [ ] `cmd/fix.go` — 完整实现
  - [ ] 读取管道错误信息
  - [ ] 结合 `--file` 上下文
  - [ ] 输出修复建议
- [ ] `cmd/explain.go` — 完整实现
  - [ ] `--func` 过滤指定函数
  - [ ] 自动识别语言类型
- [ ] 兼容 provider 完整支持
  - [ ] Deepseek Coder v3
  - [ ] Qwen 2.5 Coder
  - [ ] Ollama 本地模型
- [ ] 项目级 `.tcli/config.yaml` 合并加载
- [ ] 自定义提示词模板（`.tcli/prompts/`）

---

## Phase 4 — 发布与完善（第 7-8 周）

- [ ] 配置 GitHub Actions CI（Push Tag 自动触发 Release）
  - [ ] 创建 `.github/workflows/release.yml`
- [ ] Shell 补全
  - [ ] `tcli completion bash`
  - [ ] `tcli completion zsh`
  - [ ] `tcli completion fish`
  - [ ] `tcli completion powershell`
- [ ] `tcli update` 命令（自动检查并下载最新版本）
- [ ] 完整中文 README.md
  - [ ] 安装说明（Homebrew / 直接下载 / 从源码构建）
  - [ ] 快速上手
  - [ ] 配置参考
  - [ ] 命令速查表
- [ ] Windows 兼容性专项测试
  - [ ] PowerShell 颜色渲染
  - [ ] CMD 基础功能
  - [ ] Windows Terminal 完整体验
- [ ] 发布 v0.1.0

---

## 后期规划（v1.0+）

- [ ] MCP 协议客户端支持
- [ ] tree-sitter 本地 AST 代码分析
- [ ] `tcli index` — 构建项目代码索引（类 Aider Repo Map）
- [ ] 多轮对话历史持久化（`~/.tcli/history/`）
- [ ] 插件系统（`.tcli/plugins/`）

---

## 备注

> 目前本机未安装 Go，需先完成以下步骤才能运行项目：
>
> 1. 从 [go.dev/dl](https://go.dev/dl/) 下载安装 Go ≥ 1.23
> 2. 在 `f:\桌面\TCil\` 目录执行：
>    ```powershell
>    go mod tidy
>    go build -o tcli.exe .
>    .\tcli.exe version
>    ```
> 3. 配置 API Key：
>    ```powershell
>    .\tcli.exe config set providers.anthropic.api_key sk-ant-xxx
>    .\tcli.exe config set default.provider anthropic
>    ```

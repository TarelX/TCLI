package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TarelX/TCLI/internal/ai"
	cctx "github.com/TarelX/TCLI/internal/context"
	"github.com/TarelX/TCLI/internal/prompt"
	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var askCmd = &cobra.Command{
	Use:   "ask [问题]",
	Short: "单次 AI 问答（支持管道输入）",
	Long: `向 AI 提问，流式输出结果。支持通过管道传入上下文：

  tcli ask "这个函数做什么？" --file main.go
  cat error.log | tcli ask "这个报错什么原因"
  git diff HEAD~1 | tcli ask "帮我写 commit message"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAsk,
}

func init() {
	askCmd.Flags().StringP("file", "f", "", "附加指定文件内容作为上下文")
	askCmd.Flags().Bool("raw", false, "纯文本输出，不渲染 Markdown")
}

func runAsk(cmd *cobra.Command, args []string) error {
	// 1. 获取问题文本
	question := ""
	if len(args) > 0 {
		question = args[0]
	}

	// 2. 读取管道输入（如果有）
	pipeInput := ""
	if !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		data, err := io.ReadAll(os.Stdin)
		if err == nil && len(data) > 0 {
			pipeInput = string(data)
		}
	}

	// 没有问题也没有管道输入，提示用户
	if question == "" && pipeInput == "" {
		return fmt.Errorf("请提供问题，例如：tcli ask \"这段代码做什么？\"")
	}

	// 如果只有管道输入没有问题，默认问题为"请分析以下内容"
	if question == "" {
		question = "请分析以下内容"
	}

	// 3. 读取 --file 指定的文件内容
	fileContent := ""
	filePath, _ := cmd.Flags().GetString("file")
	if filePath != "" {
		content, err := cctx.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("读取文件失败：%w", err)
		}
		fileContent = content
	}

	// 4. 收集代码上下文（除非 --no-context）
	contextInfo := ""
	noContext, _ := cmd.Root().PersistentFlags().GetBool("no-context")
	if !noContext {
		ctx := cctx.Collect(viper.GetInt("context.max_tokens"))
		contextInfo = ctx.Summary
	}

	// 5. 组装 prompt
	userContent := prompt.AskWithContext(question, contextInfo, fileContent)
	if pipeInput != "" {
		userContent = "## 管道输入\n\n```\n" + pipeInput + "\n```\n\n" + userContent
	}

	messages := []ai.Message{
		ai.SystemMessage(prompt.SystemBase),
		ai.UserMessage(userContent),
	}

	// 6. 创建 AI 客户端
	client, err := ai.NewClientFromConfig()
	if err != nil {
		return err
	}

	// 7. 判断输出模式
	rawMode, _ := cmd.Flags().GetBool("raw")
	// 如果输出到管道（非终端），自动切换 raw 模式
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		rawMode = true
	}

	req := ai.Request{
		Messages:    messages,
		MaxTokens:   viper.GetInt("default.max_tokens"),
		Temperature: float32(viper.GetFloat64("default.temperature")),
		Stream:      true,
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 8192
	}
	if req.Temperature == 0 {
		req.Temperature = 0.7
	}

	// 8. 流式输出
	var fullResponse strings.Builder
	err = client.Stream(context.Background(), req, func(delta string, done bool, streamErr error) {
		if streamErr != nil {
			fmt.Fprintf(os.Stderr, "\n错误：%v\n", streamErr)
			return
		}
		if !done {
			fullResponse.WriteString(delta)
			if rawMode {
				fmt.Print(delta)
			}
		}
	})
	if err != nil {
		return fmt.Errorf("AI 请求失败：%w", err)
	}

	// 9. 渲染最终输出
	if rawMode {
		// 流式已输出，补个换行
		fmt.Println()
	} else {
		// Glamour Markdown 渲染
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(100),
		)
		if err != nil {
			// 渲染器创建失败，回退到原始输出
			fmt.Println(fullResponse.String())
			return nil
		}
		rendered, err := renderer.Render(fullResponse.String())
		if err != nil {
			fmt.Println(fullResponse.String())
			return nil
		}
		fmt.Print(rendered)
	}

	return nil
}

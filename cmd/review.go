package cmd

import (
	"context"
	"fmt"
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

var reviewCmd = &cobra.Command{
	Use:   "review [文件或目录]",
	Short: "AI 代码审查",
	Long: `对指定文件或目录进行 AI 代码审查，输出结构化报告：
  - 潜在 Bug
  - 性能问题
  - 安全问题
  - 可读性建议

  tcli review main.go
  tcli review ./src/
  tcli review --staged        # 审查 git staged 的文件
  tcli review --diff HEAD~1   # 审查最近一次提交的改动`,
	RunE: runReview,
}

func init() {
	reviewCmd.Flags().Bool("staged", false, "审查 git staged 的文件")
	reviewCmd.Flags().String("diff", "", "审查指定 git ref 的改动（如 HEAD~1）")
}

func runReview(cmd *cobra.Command, args []string) error {
	// 1. 收集要审查的代码内容
	var codeContent string
	var err error

	staged, _ := cmd.Flags().GetBool("staged")
	diffRef, _ := cmd.Flags().GetString("diff")

	switch {
	case staged:
		codeContent, err = cctx.GitStagedDiff()
		if err != nil {
			return fmt.Errorf("获取 staged diff 失败：%w", err)
		}
		if strings.TrimSpace(codeContent) == "" {
			return fmt.Errorf("没有 staged 的改动")
		}
	case diffRef != "":
		codeContent, err = cctx.GitDiff(diffRef)
		if err != nil {
			return fmt.Errorf("获取 diff 失败：%w", err)
		}
		if strings.TrimSpace(codeContent) == "" {
			return fmt.Errorf("指定的 ref 没有改动")
		}
	case len(args) > 0:
		codeContent, err = cctx.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("读取文件失败：%w", err)
		}
	default:
		return fmt.Errorf("请指定要审查的文件，或使用 --staged / --diff 参数")
	}

	// 2. 调用 AI 审查
	messages := []ai.Message{
		ai.SystemMessage(prompt.SystemReview),
		ai.UserMessage("请对以下代码进行审查：\n\n" + codeContent),
	}

	return runAIAndRender(messages, "🔍 正在审查代码...")
}

// runAIAndRender 是 review/fix/explain 共用的 AI 调用 + 渲染逻辑
func runAIAndRender(messages []ai.Message, hint string) error {
	client, err := ai.NewClientFromConfig()
	if err != nil {
		return err
	}

	fmt.Println(hint)

	req := ai.Request{
		Messages:    messages,
		MaxTokens:   viper.GetInt("default.max_tokens"),
		Temperature: float32(viper.GetFloat64("default.temperature")),
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 8192
	}
	if req.Temperature == 0 {
		req.Temperature = 0.5
	}

	result, err := client.Complete(context.Background(), req)
	if err != nil {
		return fmt.Errorf("AI 请求失败：%w", err)
	}

	// 渲染输出
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println(result)
		return nil
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		fmt.Println(result)
		return nil
	}
	rendered, err := renderer.Render(result)
	if err != nil {
		fmt.Println(result)
		return nil
	}
	fmt.Print(rendered)
	return nil
}

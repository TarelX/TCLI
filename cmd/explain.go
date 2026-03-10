package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/TarelX/TCLI/internal/ai"
	cctx "github.com/TarelX/TCLI/internal/context"
	"github.com/TarelX/TCLI/internal/prompt"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain [文件]",
	Short: "解释代码文件或函数的作用",
	Long: `让 AI 解释指定文件或函数的功能、逻辑和设计意图：

  tcli explain main.go
  tcli explain main.go --func ParseConfig
  cat main.go | tcli explain`,
	RunE: runExplain,
}

func init() {
	explainCmd.Flags().String("func", "", "只解释指定函数名")
}

func runExplain(cmd *cobra.Command, args []string) error {
	// 1. 获取代码内容（文件参数 或 管道输入）
	var codeContent string

	if len(args) > 0 {
		content, err := cctx.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("读取文件失败：%w", err)
		}
		codeContent = content
	} else if !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		data, err := io.ReadAll(os.Stdin)
		if err == nil && len(data) > 0 {
			codeContent = "```\n" + string(data) + "\n```"
		}
	}

	if codeContent == "" {
		return fmt.Errorf("请指定要解释的文件，例如：tcli explain main.go")
	}

	// 2. 组装消息
	funcName, _ := cmd.Flags().GetString("func")
	userContent := "请解释以下代码：\n\n" + codeContent
	if funcName != "" {
		userContent = fmt.Sprintf("请重点解释其中的 `%s` 函数：\n\n%s", funcName, codeContent)
	}

	messages := []ai.Message{
		ai.SystemMessage(prompt.SystemExplain),
		ai.UserMessage(userContent),
	}

	return runAIAndRender(messages, "📖 正在分析代码...")
}

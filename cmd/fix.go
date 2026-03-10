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

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "修复代码错误（接收管道输入）",
	Long: `接收编译器或运行时的错误信息，给出修复建议：

  go build 2>&1 | tcli fix
  npm run build 2>&1 | tcli fix
  tcli fix --error "undefined: foo" --file main.go`,
	RunE: runFix,
}

func init() {
	fixCmd.Flags().String("error", "", "直接传入错误信息文本")
	fixCmd.Flags().StringP("file", "f", "", "关联的源文件")
}

func runFix(cmd *cobra.Command, args []string) error {
	// 1. 获取错误信息（管道输入 或 --error 参数）
	errorText, _ := cmd.Flags().GetString("error")

	if errorText == "" && !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		data, err := io.ReadAll(os.Stdin)
		if err == nil && len(data) > 0 {
			errorText = string(data)
		}
	}

	if errorText == "" {
		return fmt.Errorf("请提供错误信息，例如：go build 2>&1 | tcli fix")
	}

	// 2. 读取关联源文件（可选）
	fileContent := ""
	filePath, _ := cmd.Flags().GetString("file")
	if filePath != "" {
		content, err := cctx.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("读取文件失败：%w", err)
		}
		fileContent = content
	}

	// 3. 组装消息
	userContent := "## 错误信息\n\n```\n" + errorText + "\n```"
	if fileContent != "" {
		userContent += "\n\n## 相关代码\n\n" + fileContent
	}

	messages := []ai.Message{
		ai.SystemMessage(prompt.SystemFix),
		ai.UserMessage(userContent),
	}

	return runAIAndRender(messages, "🔧 正在分析错误...")
}

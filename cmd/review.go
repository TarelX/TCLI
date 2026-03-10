package cmd

import (
	"github.com/spf13/cobra"
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
	// TODO Phase 3
	cmd.Println("review 命令正在开发中（Phase 3）")
	return nil
}

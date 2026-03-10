package cmd

import (
	"github.com/spf13/cobra"
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
	// TODO Phase 1
	cmd.Println("ask 命令正在开发中（Phase 1）")
	return nil
}

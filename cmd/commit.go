package cmd

import (
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "AI 生成 git commit message",
	Long: `分析 git staged 的改动，自动生成 commit message。

  tcli commit                  # 生成后手动确认再提交
  tcli commit --push           # 生成后自动 git commit && git push
  tcli commit --convention     # 强制 Conventional Commits 格式
  tcli commit --lang en        # 英文 commit（默认中文）`,
	RunE: runCommit,
}

func init() {
	commitCmd.Flags().Bool("push", false, "生成并确认后自动执行 git commit && git push")
	commitCmd.Flags().Bool("convention", false, "强制使用 Conventional Commits 格式")
	commitCmd.Flags().String("lang", "zh", "commit message 语言（zh | en）")
}

func runCommit(cmd *cobra.Command, args []string) error {
	// TODO Phase 1
	cmd.Println("commit 命令正在开发中（Phase 1）")
	return nil
}

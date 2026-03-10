package cmd

import (
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "启动交互式 AI 对话界面（TUI）",
	Long: `启动基于 Bubbletea 的全交互式对话界面。
支持多行输入、流式响应、Markdown 渲染、Token 用量实时显示。

快捷键：
  Shift+Enter   换行
  Enter         发送消息
  ↑ / ↓        翻历史输入
  Ctrl+C        退出

内置命令（输入 / 开头）：
  /clear        清空对话历史
  /copy         复制最后一条回复到剪贴板
  /model        切换模型
  /save         保存对话到文件`,
	RunE: runChat,
}

func init() {
	chatCmd.Flags().Bool("no-splash", false, "跳过启动动画")
}

func runChat(cmd *cobra.Command, args []string) error {
	// TODO Phase 2: 启动 TUI
	// noSplash, _ := cmd.Flags().GetBool("no-splash")
	// return tui.StartChat(noSplash)
	cmd.Println("chat 命令正在开发中（Phase 2）")
	return nil
}

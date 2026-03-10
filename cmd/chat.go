package cmd

import (
	"github.com/TarelX/TCLI/internal/ai"
	"github.com/TarelX/TCLI/internal/tui"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "启动交互式 AI 对话界面（TUI）",
	Long: `启动基于 Bubbletea 的全交互式对话界面。
支持多行输入、流式响应、Markdown 渲染、Token 用量实时显示。

快捷键：
  Enter         发送消息
  Ctrl+C        退出

内置命令（输入 / 开头）：
  /clear        清空对话历史
  /copy         复制最后一条回复到剪贴板
  /help         查看帮助`,
	RunE: runChat,
}

func init() {
	chatCmd.Flags().Bool("no-splash", false, "跳过启动动画")
}

func runChat(cmd *cobra.Command, args []string) error {
	// 创建 AI 客户端
	client, err := ai.NewClientFromConfig()
	if err != nil {
		return err
	}

	return tui.StartChat(client, appVersion)
}

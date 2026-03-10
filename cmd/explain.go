package cmd

import (
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
	// TODO Phase 3
	cmd.Println("explain 命令正在开发中（Phase 3）")
	return nil
}

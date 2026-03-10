package cmd

import (
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
	// TODO Phase 3
	cmd.Println("fix 命令正在开发中（Phase 3）")
	return nil
}

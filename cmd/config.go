package cmd

import (
	"fmt"

	"github.com/TarelX/TCLI/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理 TCli 配置",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "设置配置项",
	Example: `  tcli config set default.provider anthropic
  tcli config set default.model claude-3-7-sonnet
  tcli config set providers.openai.api_key sk-xxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 确保配置目录存在
		cfgDir, err := config.EnsureConfigDir()
		if err != nil {
			return fmt.Errorf("创建配置目录失败：%w", err)
		}
		viper.Set(args[0], args[1])
		if err := viper.WriteConfig(); err != nil {
			// 配置文件不存在则创建
			viper.SetConfigFile(cfgDir + "/config.yaml")
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("写入配置失败：%w", err)
			}
		}
		fmt.Printf("✓ 已设置 %s = %s\n", args[0], args[1])
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "显示所有配置项",
	Run: func(cmd *cobra.Command, args []string) {
		settings := viper.AllSettings()
		if len(settings) == 0 {
			fmt.Println("暂无配置，请先运行 'tcli config set' 进行配置。")
			return
		}
		// TODO: 格式化输出，隐藏 api_key 敏感值
		fmt.Printf("配置文件：%s\n\n", viper.ConfigFileUsed())
		for k, v := range settings {
			fmt.Printf("  %-40s %v\n", k, v)
		}
	},
}

var configProvidersCmd = &cobra.Command{
	Use:   "providers",
	Short: "显示已配置的 provider 列表",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Phase 1
		fmt.Println("providers 子命令正在开发中（Phase 1）")
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configProvidersCmd)
}

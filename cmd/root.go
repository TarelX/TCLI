package cmd

import (
	"fmt"
	"os"

	"github.com/TarelX/TCLI/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appVersion string
	buildDate  string
	cfgFile    string
)

var rootCmd = &cobra.Command{
	Use:   "tcli",
	Short: "TCli — AI 驱动的代码 CLI 工具",
	Long: `TCli (Terminal Code Intelligence CLI)
面向中文开发者的 AI 代码辅助工具，支持 OpenAI / Anthropic / 本地模型。

使用 'tcli --help' 查看所有命令。`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 是 main.go 的入口，注入版本信息
func Execute(version, date string) {
	appVersion = version
	buildDate = date
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "错误："+err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径（默认 ~/.tcli/config.yaml）")
	rootCmd.PersistentFlags().StringP("provider", "p", "", "覆盖默认 provider（openai | anthropic | compatible）")
	rootCmd.PersistentFlags().StringP("model", "m", "", "覆盖默认模型名")
	rootCmd.PersistentFlags().Bool("no-context", false, "不注入代码上下文")
	rootCmd.PersistentFlags().Bool("raw", false, "纯文本输出，不渲染 Markdown（适合管道）")

	_ = viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider"))
	_ = viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model"))

	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(askCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(fixCmd)
	rootCmd.AddCommand(explainCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cfgDir, err := config.DefaultConfigDir()
		if err == nil {
			viper.AddConfigPath(cfgDir)
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("TCLI")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时静默，首次使用正常
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "读取配置文件出错："+err.Error())
		}
	}
}

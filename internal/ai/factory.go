package ai

import (
	"errors"
	"os"

	"github.com/TarelX/TCLI/internal/config"
	"github.com/spf13/viper"
)

// NewClientFromConfig 根据当前 viper 配置创建合适的 AI 客户端
func NewClientFromConfig() (Client, error) {
	provider := viper.GetString("default.provider")
	model := viper.GetString("default.model")

	// 命令行 flag 优先级最高
	if p := viper.GetString("provider"); p != "" {
		provider = p
	}
	if m := viper.GetString("model"); m != "" {
		model = m
	}

	switch Provider(provider) {
	case ProviderOpenAI:
		apiKey := viper.GetString("providers.openai.api_key")
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		if apiKey == "" {
			return nil, errors.New("未找到 OpenAI API Key，请运行：tcli config set providers.openai.api_key sk-xxx")
		}
		baseURL := viper.GetString("providers.openai.base_url")
		return NewOpenAIClient(apiKey, model, baseURL), nil

	case ProviderAnthropic:
		apiKey := viper.GetString("providers.anthropic.api_key")
		if apiKey == "" {
			apiKey = os.Getenv("ANTHROPIC_API_KEY")
		}
		if apiKey == "" {
			return nil, errors.New("未找到 Anthropic API Key，请运行：tcli config set providers.anthropic.api_key sk-ant-xxx")
		}
		baseURL := viper.GetString("providers.anthropic.base_url")
		return NewAnthropicClient(apiKey, model, baseURL), nil

	case ProviderCompatible:
		// 从 compatible 列表中找到名称匹配的条目
		var entries []config.CompatibleEntry
		if err := viper.UnmarshalKey("providers.compatible", &entries); err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.Model == model || e.Name == provider {
				apiKey := e.APIKey
				if apiKey == "" {
					apiKey = os.Getenv("TCLI_COMPATIBLE_API_KEY")
				}
				return NewCompatibleClient(e.Name, apiKey, e.Model, e.BaseURL), nil
			}
		}
		return nil, errors.New("在 compatible 配置中未找到匹配的 provider，请检查 ~/.tcli/config.yaml")

	default:
		return nil, errors.New("未知的 provider：" + provider + "，可选值：openai | anthropic | compatible")
	}
}

package config

// Config 是全局配置的完整结构，对应 ~/.tcli/config.yaml
type Config struct {
	Version   int             `mapstructure:"version"`
	Default   DefaultConfig   `mapstructure:"default"`
	Providers ProvidersConfig `mapstructure:"providers"`
	Context   ContextConfig   `mapstructure:"context"`
	UI        UIConfig        `mapstructure:"ui"`
}

type DefaultConfig struct {
	Provider    string  `mapstructure:"provider"`
	Model       string  `mapstructure:"model"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float32 `mapstructure:"temperature"`
}

type ProvidersConfig struct {
	OpenAI     OpenAIConfig      `mapstructure:"openai"`
	Anthropic  AnthropicConfig   `mapstructure:"anthropic"`
	Compatible []CompatibleEntry `mapstructure:"compatible"`
}

type OpenAIConfig struct {
	APIKey  string   `mapstructure:"api_key"`
	BaseURL string   `mapstructure:"base_url"`
	Models  []string `mapstructure:"models"`
}

type AnthropicConfig struct {
	APIKey string   `mapstructure:"api_key"`
	Models []string `mapstructure:"models"`
}

// CompatibleEntry 代表一个 OpenAI 兼容接口（Deepseek、Qwen、Ollama 等）
type CompatibleEntry struct {
	Name    string `mapstructure:"name"`
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

type ContextConfig struct {
	MaxTokens      int      `mapstructure:"max_tokens"`
	IncludeGitInfo bool     `mapstructure:"include_git_info"`
	IgnorePatterns []string `mapstructure:"ignore_patterns"`
	ExtraFiles     []string `mapstructure:"extra_files"`
}

type UIConfig struct {
	Theme          string `mapstructure:"theme"`
	Language       string `mapstructure:"language"`
	ShowTokenCount bool   `mapstructure:"show_token_count"`
}

// DefaultConfig 返回开箱即用的默认配置
func DefaultValues() Config {
	return Config{
		Version: 1,
		Default: DefaultConfig{
			Provider:    "anthropic",
			Model:       "claude-3-7-sonnet",
			MaxTokens:   8192,
			Temperature: 0.7,
		},
		Context: ContextConfig{
			MaxTokens:      32768,
			IncludeGitInfo: true,
			IgnorePatterns: []string{"*.lock", "node_modules/", "vendor/", ".git/"},
		},
		UI: UIConfig{
			Theme:          "auto",
			Language:       "zh",
			ShowTokenCount: true,
		},
	}
}

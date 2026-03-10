package ai

// CompatibleClient 对接 OpenAI 兼容接口（Deepseek / Qwen / Ollama 等）
// 直接复用 OpenAIClient，仅更改 BaseURL 和 Provider 标识
type CompatibleClient struct {
	*OpenAIClient
	name string
}

// NewCompatibleClient 创建兼容接口客户端
// name: 用于显示的名称（如 "deepseek"、"ollama"）
func NewCompatibleClient(name, apiKey, model, baseURL string) *CompatibleClient {
	return &CompatibleClient{
		OpenAIClient: NewOpenAIClient(apiKey, model, baseURL),
		name:         name,
	}
}

func (c *CompatibleClient) Provider() Provider { return ProviderCompatible }

func (c *CompatibleClient) ModelName() string {
	return c.name + "/" + c.OpenAIClient.modelName
}

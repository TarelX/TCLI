package prompt

import (
	"os"
	"path/filepath"
	"strings"
)

// LoadCustom 从 .tcli/prompts/ 目录加载自定义提示词模板
// 文件名（不含扩展名）作为模板名，如 review.md → "review"
func LoadCustom(promptsDir string) (map[string]string, error) {
	result := make(map[string]string)

	entries, err := os.ReadDir(promptsDir)
	if os.IsNotExist(err) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		data, err := os.ReadFile(filepath.Join(promptsDir, e.Name()))
		if err != nil {
			continue
		}
		result[name] = string(data)
	}
	return result, nil
}

package context

import (
	"os"
	"path/filepath"
)

// ReadFile 读取文件内容，返回带语言标注的 Markdown 代码块
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lang := detectLang(path)
	return "```" + lang + "\n" + string(data) + "\n```", nil
}

// detectLang 根据文件扩展名返回代码块语言标识
func detectLang(path string) string {
	ext := filepath.Ext(path)
	langs := map[string]string{
		".go":   "go",
		".ts":   "typescript",
		".tsx":  "tsx",
		".js":   "javascript",
		".jsx":  "jsx",
		".rs":   "rust",
		".py":   "python",
		".java": "java",
		".c":    "c",
		".cpp":  "cpp",
		".cs":   "csharp",
		".rb":   "ruby",
		".sh":   "bash",
		".yaml": "yaml",
		".yml":  "yaml",
		".json": "json",
		".toml": "toml",
		".md":   "markdown",
		".sql":  "sql",
	}
	if lang, ok := langs[ext]; ok {
		return lang
	}
	return ""
}

// fileExists 检查文件是否存在
func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

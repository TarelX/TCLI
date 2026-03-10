package chat

import (
	"os"
	"path/filepath"
	"strings"
)

// 可用的 slash 命令列表
var slashCommands = []string{"/clear", "/copy", "/help"}

// updateCompletions 根据当前输入内容更新候选列表
// 只在输入内容实际变化时才重新计算，避免每次 tick 都重置 completionIdx
func (m *Model) updateCompletions() {
	text := m.input.Value()

	// 输入没变化，不重算（保护 completionIdx 不被重置）
	if text == m.lastInput {
		return
	}
	m.lastInput = text

	// 检测 / 命令补全（仅当输入以 / 开头且只有一个词时）
	if strings.HasPrefix(text, "/") && !strings.Contains(text, " ") {
		prefix := text
		var matches []string
		for _, cmd := range slashCommands {
			if strings.HasPrefix(cmd, prefix) {
				matches = append(matches, cmd)
			}
		}
		if len(matches) > 0 && text != matches[0] {
			m.completions = matches
			m.completionMode = "cmd"
			m.completionIdx = 0
			return
		}
	}

	// 检测 @ 文件补全：找到最后一个 @ 后的部分文件名
	atIdx := strings.LastIndex(text, "@")
	if atIdx >= 0 {
		afterAt := text[atIdx+1:]
		// 只在 @ 后没有空格时触发（说明用户还在输入文件名）
		if !strings.Contains(afterAt, " ") {
			matches := matchFiles(afterAt, 8)
			if len(matches) > 0 {
				m.completions = matches
				m.completionMode = "file"
				m.completionIdx = 0
				return
			}
		}
	}

	// 没有匹配，清除补全状态
	m.completions = nil
	m.completionMode = ""
	m.completionIdx = 0
}

// applyCompletion 将当前选中的候选项插入输入框
func (m *Model) applyCompletion() {
	if len(m.completions) == 0 {
		return
	}
	selected := m.completions[m.completionIdx]
	text := m.input.Value()

	switch m.completionMode {
	case "cmd":
		// 替换整个输入为选中的命令
		m.input.Reset()
		m.input.InsertString(selected + " ")
	case "file":
		// 找到最后一个 @，替换 @ 后面的部分
		atIdx := strings.LastIndex(text, "@")
		if atIdx >= 0 {
			before := text[:atIdx+1]
			m.input.Reset()
			m.input.InsertString(before + selected + " ")
		}
	}

	m.completions = nil
	m.completionMode = ""
}

// matchFiles 扫描当前目录，返回匹配 prefix 的文件路径（最多 maxResults 个）
func matchFiles(prefix string, maxResults int) []string {
	var results []string

	// 分离目录和文件名前缀
	dir := "."
	filePrefix := prefix
	if idx := strings.LastIndexAny(prefix, "/\\"); idx >= 0 {
		dir = prefix[:idx]
		if dir == "" {
			dir = "."
		}
		filePrefix = prefix[idx+1:]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		name := entry.Name()
		// 跳过隐藏文件和常见无用目录
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "__pycache__" {
			continue
		}

		fullPath := name
		if dir != "." {
			fullPath = filepath.Join(dir, name)
		}

		// 目录加 / 后缀方便继续补全
		if entry.IsDir() {
			fullPath += "/"
		}

		if filePrefix == "" || strings.HasPrefix(strings.ToLower(name), strings.ToLower(filePrefix)) {
			results = append(results, fullPath)
			if len(results) >= maxResults {
				break
			}
		}
	}

	return results
}

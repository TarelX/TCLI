package context

import (
	"fmt"
	"strings"
)

// CodeContext 收集到的代码上下文
type CodeContext struct {
	Git        GitInfo
	ProjectType string
	Summary    string // 最终拼装的上下文字符串，注入到 system prompt
}

// Collect 收集当前工作目录的代码上下文
func Collect(maxTokens int) CodeContext {
	ctx := CodeContext{}

	ctx.Git = CollectGitInfo()
	ctx.ProjectType = detectProjectType()

	ctx.Summary = buildSummary(ctx, maxTokens)
	return ctx
}

// detectProjectType 识别项目类型
func detectProjectType() string {
	checks := []struct {
		file    string
		name    string
	}{
		{"go.mod", "Go"},
		{"package.json", "Node.js"},
		{"Cargo.toml", "Rust"},
		{"requirements.txt", "Python"},
		{"pyproject.toml", "Python"},
		{"pom.xml", "Java (Maven)"},
		{"build.gradle", "Java (Gradle)"},
		{"Gemfile", "Ruby"},
	}
	for _, c := range checks {
		if fileExists(c.file) {
			return c.name
		}
	}
	return "未知"
}

func buildSummary(ctx CodeContext, _ int) string {
	var sb strings.Builder

	sb.WriteString("## 项目上下文\n\n")
	sb.WriteString(fmt.Sprintf("- 项目类型：%s\n", ctx.ProjectType))

	if ctx.Git.IsRepo {
		sb.WriteString(fmt.Sprintf("- 当前分支：%s\n", ctx.Git.Branch))
		if len(ctx.Git.RecentCommits) > 0 {
			sb.WriteString("- 最近提交：\n")
			for _, c := range ctx.Git.RecentCommits {
				sb.WriteString(fmt.Sprintf("  - %s\n", c))
			}
		}
		if ctx.Git.StatusSummary != "" {
			sb.WriteString("- 未提交变更：\n```\n")
			sb.WriteString(ctx.Git.StatusSummary)
			sb.WriteString("\n```\n")
		}
	}

	return sb.String()
}

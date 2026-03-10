package context

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitInfo 包含当前仓库的基本信息
type GitInfo struct {
	Branch        string
	RecentCommits []string // 最近 3 条 commit message
	StatusSummary string   // git status --short 的输出
	IsRepo        bool
}

// CollectGitInfo 收集当前目录的 git 信息
func CollectGitInfo() GitInfo {
	if !isGitRepo() {
		return GitInfo{IsRepo: false}
	}
	return GitInfo{
		IsRepo:        true,
		Branch:        gitBranch(),
		RecentCommits: gitRecentCommits(3),
		StatusSummary: gitStatus(),
	}
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

func gitBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func gitRecentCommits(n int) []string {
	out, err := exec.Command("git", "log", "--oneline", "-n", fmt.Sprintf("%d", n)).Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	result := make([]string, 0, n)
	for _, l := range lines {
		if l != "" {
			result = append(result, l)
		}
	}
	return result
}

func gitStatus() string {
	out, err := exec.Command("git", "status", "--short").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GitStagedDiff 返回 staged 文件的 diff 内容
func GitStagedDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--staged").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// GitDiff 返回指定 ref 的 diff 内容
func GitDiff(ref string) (string, error) {
	out, err := exec.Command("git", "diff", ref).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

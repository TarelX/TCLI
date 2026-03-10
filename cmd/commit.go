package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TarelX/TCLI/internal/ai"
	cctx "github.com/TarelX/TCLI/internal/context"
	"github.com/TarelX/TCLI/internal/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "AI 生成 git commit message",
	Long: `分析 git staged 的改动，自动生成 commit message。

  tcli commit                  # 生成后手动确认再提交
  tcli commit --push           # 生成后自动 git commit && git push
  tcli commit --convention     # 强制 Conventional Commits 格式
  tcli commit --lang en        # 英文 commit（默认中文）`,
	RunE: runCommit,
}

func init() {
	commitCmd.Flags().Bool("push", false, "生成并确认后自动执行 git commit && git push")
	commitCmd.Flags().Bool("convention", false, "强制使用 Conventional Commits 格式")
	commitCmd.Flags().String("lang", "zh", "commit message 语言（zh | en）")
}

func runCommit(cmd *cobra.Command, args []string) error {
	// 1. 获取 staged diff
	diff, err := cctx.GitStagedDiff()
	if err != nil {
		return fmt.Errorf("获取 git diff 失败：%w", err)
	}
	diff = strings.TrimSpace(diff)
	if diff == "" {
		return fmt.Errorf("没有 staged 的改动。请先使用 git add 添加文件")
	}

	// 2. 构建 system prompt
	convention, _ := cmd.Flags().GetBool("convention")
	lang, _ := cmd.Flags().GetString("lang")

	systemPrompt := prompt.SystemCommit
	if convention {
		systemPrompt += "\n\n必须使用 Conventional Commits 格式（feat/fix/docs/style/refactor/test/chore）。"
	}
	if lang == "en" {
		systemPrompt += "\n\n请使用英文撰写 commit message。"
	}

	// 3. 组装消息
	userContent := "以下是 git staged 的改动，请生成 commit message：\n\n```diff\n" + diff + "\n```"

	messages := []ai.Message{
		ai.SystemMessage(systemPrompt),
		ai.UserMessage(userContent),
	}

	// 4. 创建 AI 客户端
	client, err := ai.NewClientFromConfig()
	if err != nil {
		return err
	}

	req := ai.Request{
		Messages:    messages,
		MaxTokens:   viper.GetInt("default.max_tokens"),
		Temperature: float32(viper.GetFloat64("default.temperature")),
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 2048
	}
	if req.Temperature == 0 {
		req.Temperature = 0.3
	}

	// 5. 调用 AI 生成（非流式，一次性获取完整结果）
	fmt.Println("🤖 正在分析改动，生成 commit message...")
	commitMsg, err := client.Complete(context.Background(), req)
	if err != nil {
		return fmt.Errorf("AI 生成失败：%w", err)
	}
	commitMsg = strings.TrimSpace(commitMsg)

	// 6. 显示生成的 commit message
	fmt.Println()
	fmt.Println("──────────────────────────────────────")
	fmt.Println(commitMsg)
	fmt.Println("──────────────────────────────────────")
	fmt.Println()

	// 7. 交互式确认
	autoPush, _ := cmd.Flags().GetBool("push")
	action := promptUserAction(autoPush)

	switch action {
	case "y":
		// 执行 git commit
		if err := gitCommit(commitMsg); err != nil {
			return fmt.Errorf("git commit 失败：%w", err)
		}
		fmt.Println("✓ 已提交")

		if autoPush {
			fmt.Println("📤 正在推送...")
			if err := gitPush(); err != nil {
				return fmt.Errorf("git push 失败：%w", err)
			}
			fmt.Println("✓ 已推送")
		}

	case "e":
		// 编辑模式：将消息写入临时文件，用 git commit -e
		if err := gitCommitEdit(commitMsg); err != nil {
			return fmt.Errorf("git commit 失败：%w", err)
		}
		fmt.Println("✓ 已提交（编辑模式）")

	case "n":
		fmt.Println("已取消")

	case "c":
		// 复制到剪贴板
		if err := copyToClipboard(commitMsg); err != nil {
			fmt.Printf("复制失败：%v\n", err)
		} else {
			fmt.Println("✓ 已复制到剪贴板")
		}
	}

	return nil
}

// promptUserAction 显示交互式选项，返回用户选择
func promptUserAction(showPush bool) string {
	pushHint := ""
	if showPush {
		pushHint = " + push"
	}
	fmt.Printf("[y] 确认提交%s  [e] 编辑后提交  [c] 复制  [n] 取消: ", pushHint)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if input == "" {
			return "y"
		}
		return input
	}
	return "n"
}

// gitCommit 执行 git commit
func gitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// gitCommitEdit 使用编辑器模式提交（将AI生成的消息作为初始内容）
func gitCommitEdit(message string) error {
	cmd := exec.Command("git", "commit", "-e", "-m", message)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// gitPush 执行 git push
func gitPush() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// copyToClipboard 复制文本到剪贴板
func copyToClipboard(text string) error {
	// Windows: clip, macOS: pbcopy, Linux: xclip
	var cmd *exec.Cmd
	switch {
	case isCommandAvailable("clip"):
		cmd = exec.Command("clip")
	case isCommandAvailable("pbcopy"):
		cmd = exec.Command("pbcopy")
	case isCommandAvailable("xclip"):
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return fmt.Errorf("未找到剪贴板工具（clip/pbcopy/xclip）")
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// isCommandAvailable 检查命令是否存在
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

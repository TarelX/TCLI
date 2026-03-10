package config

import (
	"os"
	"path/filepath"
)

// DefaultConfigDir 返回 TCli 全局配置目录
// Windows: %APPDATA%\tcli
// macOS/Linux: ~/.tcli
func DefaultConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		// 回退到 home 目录
		home, err2 := os.UserHomeDir()
		if err2 != nil {
			return "", err2
		}
		return filepath.Join(home, ".tcli"), nil
	}
	return filepath.Join(base, "tcli"), nil
}

// EnsureConfigDir 确保配置目录存在
func EnsureConfigDir() (string, error) {
	dir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

// ProjectConfigDir 返回项目级配置目录（.tcli/）
// 从当前目录向上查找，直到根目录
func ProjectConfigDir() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for {
		candidate := filepath.Join(dir, ".tcli")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 创建路径
func CreateAllPath(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return os.Chmod(path, os.ModePerm)
}

func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// FormatPath 格式化地址格式
func FormatPath(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

// 判断文件或路径是否存在
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}

	return filepath.Dir(path), nil
}

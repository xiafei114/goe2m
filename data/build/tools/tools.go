package tools

import (
	"os"
	"os/exec"
	"path/filepath"
)

// GetModelPath 获取目录地址
func GetModelPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path := filepath.Dir(file)
	path, _ = filepath.Abs(path)

	return path
}

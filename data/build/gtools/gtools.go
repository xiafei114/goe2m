package gtools

import (
	"fmt"
	"goe2m/data/config"
)

// Execute 执行
func Execute() {
	fmt.Println(config.GetInFilePath())
}

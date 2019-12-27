package config

// Config custom config struct
type Config struct {
	InFilePath string `yaml:"in_file_path"`
	OutDir     string `yaml:"out_dir"`
}

// SetOutDir Setting Output Directory.设置输出目录
func SetOutDir(outDir string) {
	_map.OutDir = outDir
}

// GetOutDir Get Output Directory.获取输出目录
func GetOutDir() string {
	return _map.OutDir
}

// SetInFilePath Setting Output Directory.设置输出目录
func SetInFilePath(inFilePath string) {
	_map.InFilePath = inFilePath
}

// GetInFilePath Get Output Directory.获取输出目录
func GetInFilePath() string {
	return _map.InFilePath
}

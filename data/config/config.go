package config

// Config custom config struct
type Config struct {
	InFilePath  string `yaml:"in_file_path"`
	OutDir      string `yaml:"out_dir"`
	ProjectName string `yaml:"project_name"`
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

// SetProjectName 项目名称
func SetProjectName(projectName string) {
	_map.ProjectName = projectName
}

// GetProjectName Get 项目名称
func GetProjectName() string {
	return _map.ProjectName
}

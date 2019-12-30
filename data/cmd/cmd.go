package cmd

import (
	"goe2m/data/build/gen"
	"goe2m/data/config"
	"os"

	"github.com/spf13/cobra"
)

var inFilePath string
var outDir string
var projectName string

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "gorm mysql reflect tools",
	Long:  `base on gorm tools for mysql database to golang struct`,
	Run: func(cmd *cobra.Command, args []string) {
		gen.Execute()
		// Start doing things.开始做事情
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&inFilePath, "inFilePath", "i", "", "Excel文件路径")
	rootCmd.MarkFlagRequired("inFilePath")

	rootCmd.PersistentFlags().StringVarP(&outDir, "outdir", "o", "", "输出目录")
	rootCmd.MarkFlagRequired("outdir")

	rootCmd.PersistentFlags().StringVarP(&projectName, "projectName", "p", "", "项目名称")
	rootCmd.MarkFlagRequired("projectName")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(outDir) > 0 {
		config.SetOutDir(outDir)
	}

	if len(inFilePath) > 0 {
		config.SetInFilePath(inFilePath)
	}

	if len(projectName) > 0 {
		config.SetProjectName(projectName)
	}
}

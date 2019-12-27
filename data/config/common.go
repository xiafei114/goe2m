package config

import (
	"fmt"
	"goe2m/data/build/tools"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var _map = Config{}

func init() {
	onInit()
}

func onInit() {
	path := tools.GetModelPath()
	err := InitFile(path + "/config.yml")
	if err != nil {
		fmt.Println("InitFile: ", err.Error())
		return
	}
}

// InitFile default value from file .
func InitFile(filename string) error {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(bs, &_map); err != nil {
		fmt.Println("read toml error: ", err.Error())
		return err
	}

	return nil
}

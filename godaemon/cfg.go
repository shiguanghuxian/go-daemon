package godaemon

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

/* 处理配置文件 */

type CFG struct {
	Debug        bool           `yaml:"debug" json:"debug"`
	Interval     string         `yaml:"interval" json:"interval"`
	FileNotExist bool           `yaml:"file_not_exist" json:"file_not_exist"`
	CmdOut       bool           `yaml:"cmd_out" json:"cmd_out"`
	Apps         []*Application `yaml:"application" json:"application"`
}

type Application struct {
	Path string `yaml:"path" json:"path"`
	Cmd  string `yaml:"cmd" json:"cmd"`
	Exec string `yaml:"exec" json:"exec"`
	Args string `yaml:"args" json:"args"`
}

// NewServerConfig 初始化一个配置文件对象
func NewCFG(path string) (config *CFG, err error) {
	if path == "" {
		path = GetRootDir() + "/config/cfg.yaml"
	}
	cnfBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	config = &CFG{}
	err = yaml.Unmarshal(cnfBytes, config)
	if err != nil {
		return
	}
	return
}

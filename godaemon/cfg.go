package godaemon

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

/* 处理配置文件 */

// CFG 根配置
type CFG struct {
	Debug        bool           `yaml:"debug"`
	Interval     string         `yaml:"interval"`
	FileNotExist bool           `yaml:"file_not_exist"`
	CmdOut       bool           `yaml:"cmd_out"`
	Apps         []*Application `yaml:"application"`
}

// Application 每个应用的配置
type Application struct {
	PidPath  string `yaml:"pid_path"` // pid 文件路径
	WorkPath string `yaml:"workpath"` // 工作目录
	// Cmd      string `yaml:"cmd"`      // 要执行的命令
	// Args     string `yaml:"args"`     // 命令参数
}

// NewCFG 初始化一个配置文件对象
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

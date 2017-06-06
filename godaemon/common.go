package godaemon

/* 公共处理函数 */

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// 获取程序跟目录
func GetRootDir() string {
	// 文件不存在获取执行路径
	file, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		file = "/usr/bin"
	}
	return file
}

// 写pid到文件
func WritePidToFile(filename string) error {
	return ioutil.WriteFile(GetRootDir()+"/var/"+filename+".pid", []byte(strconv.Itoa(os.Getpid())), 0666)
}

// 删除pid文件
func RemovePidFile(filename string) error {
	return os.Remove(GetRootDir() + "/var/" + filename + ".pid")
}

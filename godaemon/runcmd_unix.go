// +build linux freebsd darwin

package godaemon

import (
	"errors"
	"fmt"
	"log"

	sh "github.com/codeskyblue/go-sh"
	"github.com/shirou/gopsutil/process"
)

/* 非windows执行指令 */

// 根据应用程序路径执行指定指令
func (gd *GoDaemon) runCmd(appPath string, cmd ...interface{}) error {
	// logger.Info("执行命令")
	log.Println("执行命令")
	// log.Println(cmd)
	if len(cmd) == 0 {
		return errors.New("指令不能为空")
	}
	// 命令参数列表
	args := make([]interface{}, 0)
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	// 执行的指令
	command := "./" + fmt.Sprint(cmd[0])
	// 创建一个ssh session
	session := sh.NewSession()
	session.ShowCMD = false
	session.SetDir(appPath)
	// session.SetEnv("PATH", appPath)
	err := session.Command(command, args...).Command("exit").Run()
	return err
}

// 杀死指定的pid进程
func (gd *GoDaemon) killPid(pid int32) error {
	// 创建一个进程操作对象
	appProcess, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	err = appProcess.Kill()
	if err != nil {
		// 尝试其它方式杀死进程

		return err
	}
	return nil
}

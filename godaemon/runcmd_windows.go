// +build windows

package godaemon

import (
	"errors"
	"fmt"
	"log"

	sh "github.com/codeskyblue/go-sh"
)

/* 执行windows指令 */

// 根据应用程序路径执行指定指令
func (d *GoDaemon) runCmd(appPath string, cmd ...interface{}) error {
	// logger.Info("执行命令")
	log.Println("执行命令")
	// log.Println(cmd)
	if len(cmd) < 2 {
		return errors.New("指令不能为空")
	}
	// 创建一个ssh session
	session := sh.NewSession()
	session.ShowCMD = false
	// **
	// session.SetDir(appPath)
	// // "cmd", "/c",
	// err := session.Command("cmd", "/c", "start", cmd[0], cmd[1]).Command("exit").Run()
	// **

	session.SetDir(appPath)
	err := session.Command("C:\\Windows\\System32\\cmd.exe", "/c", "start", cmd[0], cmd[1]).Command("exit").Run()
	return err
}

// 杀死指定的pid进程 taskkill /PID pid
func (d *GoDaemon) killPid(pid int32) error {
	if pid < 2 {
		return errors.New("pid错误")
	}
	session := sh.NewSession()
	session.ShowCMD = false
	err := session.Command("taskkill", "/PID", fmt.Sprint(pid)).Command("exit").Run()
	if err != nil {
		return err
	}
	return nil
}

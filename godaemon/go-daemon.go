package godaemon

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/donnie4w/go-logger/logger"
	ps "github.com/mitchellh/go-ps"
	"github.com/robfig/cron"
)

/* 守护服务核心逻辑 */

type GoDaemon struct {
	cfg     *CFG
	crontab *cron.Cron
}

// NewGoDaemon 实例化守护进程
func NewGoDaemon() *GoDaemon {
	return &GoDaemon{
		crontab: cron.New(), // 定时任务对象
	}
}

// 设置配置文件
func (gd *GoDaemon) SetCfgPath(path string) (err error) {
	gd.cfg, err = NewCFG(path)
	return err
}

// Run 启动程序
func (gd *GoDaemon) Run() {
	// 没有配置实例情况使用默认
	var err error
	if gd.cfg == nil {
		gd.cfg, err = NewCFG("")
		if err != nil {
			log.Panicln(err)
		}
	}
	// 是否控制台打印日志
	logger.SetConsole(true)
	// 判断pid路径是否为空
	if len(gd.cfg.Apps) == 0 {
		logger.Warn("要守护的程序列表不存在")
	}
	// 设置任务
	gd.task()
	// 开启定时任务
	gd.crontab.Start()
}

// Stop 结束
func (gd *GoDaemon) Stop() {
	gd.crontab.Stop()
}

// 定时执行检察所有pid是否存活
func (gd *GoDaemon) task() {
	for _, v := range gd.cfg.Apps {
		v := v
		gd.crontab.AddFunc(fmt.Sprintf("@every %s", gd.cfg.Interval), func() {
			// 捕获异常
			defer func() {
				if r := recover(); r != nil {
					logger.Error(fmt.Sprintf("%T", r))
					for i := 1; i < 20; i++ {
						_, file, line, _ := runtime.Caller(i)
						if line == 0 {
							break
						}
						logger.Error(fmt.Sprintf("文件:%s,代码行:%d", file, line))
					}
				}
			}()
			// 读取pid文件
			pid, err := gd.readPid(v.Path)
			if err != nil {
				logger.Error(fmt.Sprintf("pid 文件读取错误[%s]", err.Error()))
				return
			}
			// 读取结果为空直接返回-可能是文件不存在或真的没有pid内容
			if pid == -1 {
				return
			}
			if pid == 0 && gd.cfg.FileNotExist == true {
				goto START
			} else {
				// 检查pid是否存在
				process, err := ps.FindProcess(pid)
				if err != nil {
					logger.Error(fmt.Sprintf("根据pid获取进程信息错误[%s]", err.Error()))
					return
				}
				if process != nil {
					log.Println(fmt.Sprintf("进程存在pid:%d，进程名:%s", pid, process.Executable())) // 只打印到控制台
					return
				}
			}
		START:
			/* 处理重启服务 */
			// 判断是否是执行脚本文件形式
			if v.Exec != "" {
				// 判断操作系统
				if runtime.GOOS == "windows" {
					gd.runCmd("cmd", "/c", v.Exec)
				} else {
					gd.runCmd("/bin/bash", v.Exec)
				}
			} else {
				// 命令也为空时
				if v.Cmd == "" {
					return
				}
				if v.Args == "" {
					gd.runCmd(v.Cmd)
				} else {
					gd.runCmd(v.Cmd, strings.Split(v.Args, " ")...)
				}
			}
		})
	}
}

// 执行命令
func (gd *GoDaemon) runCmd(cmdStr string, args ...string) {
	// 执行程序
	cmd := exec.Command(cmdStr, args...)
	// cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("执行命令错误[%s]", err.Error()))
		return
	} else {
		log.Println("执行成功:", cmdStr)
	}
	// 不执行此函数，直接启动命令形式无法正常启动
	err = cmd.Start()
	if err != nil {
		return
	}
	cmd.Wait()
}

// 从文件中读取pid并检查文件是否存在
func (gd *GoDaemon) readPid(path string) (int, error) {
	// 判断文件是否存在
	if gd.pathExists(path) == false {
		return 0, nil
	}
	// 读取文件
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return -1, err
	}
	// 去空格转string
	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return -1, err
	}
	return pid, nil
}

// 判断文件是否存在
func (gd *GoDaemon) pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

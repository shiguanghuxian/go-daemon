package godaemon

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/donnie4w/go-logger/logger"
	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/process"
)

// 守护进程服务对象

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

// SetCfgPath 设置配置文件
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
	logger.SetConsole(gd.cfg.Debug)
	logPath := fmt.Sprintf("%s/logs", GetRootDir())
	logger.SetRollingFile(logPath, "godaemon.log", 10, 10, logger.MB)
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
		go gd.runApp(v)

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
			// 检查程序是否存活，不存在则运行
			gd.runApp(v)
		})
	}
}

// 检查程序是否存活，不存在则运行
func (gd *GoDaemon) runApp(app *Application) {
	log.Println("检查程序，工作目录：", app.WorkPath)
	// logger.Info("检查程序，工作目录：", app.WorkPath)
	// 获取pid文件内容
	pid, err := gd.getPid(app.PidPath)
	if err != nil {
		logger.Error("读取pid文件错误：", err)
		return
	}
	// 如果pid为0，则pid文件不存在，如果配置文件不存在启动程序，则启动
	if pid == 0 {
		if gd.cfg.FileNotExist == true {
			// ## run 启动程序
			err = gd.runCmd(app.WorkPath, "control", "start")
			if err != nil {
				logger.Error(err)
			}
		}
		return
	}
	/* pid不为0，则判进程是否存在，不存在则先终止一下，再启动 */
	// 检查进程是否存在
	isExists, err := process.PidExists(int32(pid))
	if err != nil {
		logger.Error("检查pid进程是否存在错误：", err)
		return
	}
	if isExists == false {
		// ## run 启动程序
		err = gd.runCmd(app.WorkPath, "control", "restart")
		if err != nil {
			logger.Error(err)
		}
		return
	}

}

// 获取pid信息
func (gd *GoDaemon) getPid(path string) (int32, error) {
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
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		return -1, err
	}
	return int32(pid), nil
}

// 判断文件是否存在
func (gd *GoDaemon) pathExists(path string) bool {
	_, err := os.Stat(path)
	// log.Println(path)
	// log.Println(info)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

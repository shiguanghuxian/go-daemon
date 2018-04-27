package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"git.53it.net/go-daemon/godaemon"
	"github.com/kardianos/service"
)

// 程序版本
var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
	GIT_HASH   string
)

// 程序全局变量
var (
	configPath *string
	logger     service.Logger
	goDaemon   *godaemon.GoDaemon
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// 开启异步任务，开启服务
	go p.run()
	return nil
}

func (p *program) run() {
	// 存储pid
	err := godaemon.WritePidToFile("go-daemon")
	if err != nil {
		log.Println("写pid文件错误")
	}
	// 实例化web服务
	goDaemon = godaemon.NewGoDaemon()
	goDaemon.Run()
}

func (p *program) Stop(s service.Service) error {
	// 删除pid文件
	godaemon.RemovePidFile("go-daemon")
	// 关闭连接
	goDaemon.Stop()
	// 停止任务，3秒后返回
	<-time.After(time.Second * 1)
	return nil
}

func main() {
	// 全部核心运行程序
	runtime.GOMAXPROCS(runtime.NumCPU())

	svcConfig := &service.Config{
		Name:        "go-daemon.53it.net",
		DisplayName: "go-daemon.53it.net",
		Description: "go-daemon 使用golang开发的守护其它进程的服务",
	}
	// 实例化
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	// 接收一个参数，install | uninstall | start | stop | restart
	if len(os.Args) > 1 {
		if os.Args[1] == "-v" || os.Args[1] == "-version" {
			ver := fmt.Sprintf("Version: %s\nBuilt: %s\nGo version: %s\nGit commit: %s", VERSION, BUILD_TIME, GO_VERSION, GIT_HASH)
			fmt.Println(ver)
			return
		}
		// 判断是否传入配置文件路径
		if os.Args[1] == "-config" {
			if len(os.Args) < 3 {
				log.Println("参数错误,需要传入配置文件地址")
				return
			}
			err = goDaemon.SetCfgPath(os.Args[2])
			if err != nil {
				log.Println("配置文件读取错误", err)
				return
			}
		} else {
			err = service.Control(s, os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		// err = s.Run()
		// if err != nil {
		// 	logger.Error(err)
		// }
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

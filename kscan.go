package main

import (
	"embed"
	"fmt"
	"kscan/app"
	"kscan/core/slog"
	"kscan/core/tips"
	"kscan/lib/color"
	"kscan/run"
	"os"
	"runtime"
	"time"

	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/pool"
	"github.com/lcvvvv/stdio"
)

// logo信息
// const logo = `
//      _   __
//     /#| /#/
//     |#|/#/  _____  _____     *     _   _
//     |#.#/  /Edge/ /Forum\   /#\   /#\ /#\
//     |##|  |#|____ |#|      /Kv2\  |##\|#|
//     |#.#\  \r0cky\|#|     /#/_\#\ |#.#.#|
//     |#|\#\/\___|#||#|____/#/Rui\#\|#|\##|
//     \#| \#\lcvvvv/ \aels/#/ v1.87#\#/ \#/

// `
const logo = ``

// 帮助信息
const help = `
optional arguments:
  -h , --help     show this help message and exit
  -f , --fofa     从fofa获取检测对象，需提前配置环境变量:FOFA_EMAIL、FOFA_KEY
  -t , --target   指定探测对象
`

const usage = "usage: kscan [-h,--help] (-t,--target) [options]\n\n"

func main() {
	startTime := time.Now()

	//环境初始化
	Init()
	//kscan模块启动
	if len(app.Setting.Target) > 0 {
		//扫描模块初始化
		InitKscan()
		//开始扫描
		run.Start()
	}
	//计算程序运行时间
	elapsed := time.Since(startTime)
	slog.Printf(slog.INFO, "程序执行总时长为：[%s]", elapsed.String())
}

func Init() {
	app.Args.SetLogo(logo)
	app.Args.SetUsage(usage)
	app.Args.SetHelp(help)
	//参数初始化
	app.Args.Parse()
	//基础输出初始化
	stdio.SetEncoding(app.Args.Encoding)
	//参数合法性校验
	app.Args.CheckArgs()
	//日志初始化
	if app.Args.Debug {
		slog.SetLevel(slog.DEBUG)
	} else {
		slog.SetLevel(slog.INFO)
	}
	//color包初始化
	if os.Getenv("KSCAN_COLOR") == "1" {
		color.Enabled()
	}
	if app.Args.CloseColor == true {
		color.Disabled()
	}
	//pool包初始化
	pool.SetLogger(slog.Debug())
	//配置文件初始化
	app.ConfigInit()
	//Output初始化
	if app.Setting.Output != nil {
		slog.SetOutput(app.Setting.Output)
	}
	fmt.Println("Tips:", tips.GetTips())
	slog.Println(slog.INFO, "当前环境为：", runtime.GOOS, ", 输出编码为：", app.Setting.Encoding)
	if runtime.GOOS == "windows" && app.Setting.CloseColor == true {
		slog.Println(slog.INFO, "在Windows系统下，默认不会开启颜色展示，可以通过添加环境变量开启哦：KSCAN_COLOR=TRUE")
	}
}

//go:embed static/fingerprint.txt
var fingerprintEmbed embed.FS

const (
	qqwryPath       = "qqwry.dat"
	fingerprintPath = "static/fingerprint.txt"
)

func InitKscan() {
	//HTTP指纹库初始化
	fs, _ := fingerprintEmbed.Open(fingerprintPath)
	if n, err := appfinger.InitDatabaseFS(fs); err != nil {
		slog.Println(slog.ERROR, "指纹库加载失败，请检查【fingerprint.txt】文件", err)
	} else {
		slog.Printf(slog.INFO, "成功加载HTTP指纹:[%d]条", n)
	}
	//超时及日志配置
	gonmap.SetLogger(slog.Debug())
	slog.Printf(slog.INFO, "成功加载NMAP探针:[%d]个,指纹[%d]条", gonmap.UsedProbesCount, gonmap.UsedMatchCount)
}

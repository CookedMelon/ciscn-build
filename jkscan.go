package main

import (
	"embed"
	"jkscan/app"
	"jkscan/core/slog"
	"jkscan/lib/color"
	"jkscan/run"
	"os"
	"time"

	"github.com/lcvvvv/pool"

	"github.com/lcvvvv/appfinger"
)

// logo信息
const logo = `
    8  8   8  8""""8 8""""8 8""""8 8"""8 
    8  8   8  8      8    " 8    8 8   8 
    8e 8eee8e 8eeeee 8e     8eeee8 8e  8 
    88 88   8     88 88     88   8 88  8 
e   88 88   8 e   88 88   e 88   8 88  8 
8eee88 88   8 8eee88 88eee8 88   8 88  8 
                                         
`

// 帮助信息
const help = `
optional arguments:
  -h , --help     show this help message and exit
  -t , --target   指定探测对象
`

const usage = "usage: jkscan [-h,--help] (-t,--target) [options]\n\n"

func main() {
	startTime := time.Now()

	//环境初始化
	Init()
	//jkscan模块启动
	if len(app.Setting.Target) > 0 {
		//扫描模块初始化
		Initjkscan()
		run.Timer()
		//开始扫描
		run.Start()
	}
	run.FlushBuffer()
	//计算程序运行时间
	elapsed := time.Since(startTime)
	slog.Printf(slog.INFO, "程序执行总时长为：[%s]", elapsed.String())
}

func Init() {
	app.Args.SetLogo(logo)
	app.Args.SetUsage(usage)
	app.Args.SetHelp(help)
	//参数初始化
	app.Args.Init()
	//参数合法性校验
	app.Args.CheckAvailable()
	//日志初始化
	if app.Args.Debug {
		slog.SetLevel(slog.DEBUG)
	} else {
		slog.SetLevel(slog.INFO)
	}
	//color包初始化
	if os.Getenv("jkscan_COLOR") == "1" {
		color.Enabled()
	}
	//pool包初始化
	pool.SetLogger(slog.Debug())
	//配置文件初始化
	app.ConfigInit()
	//Output初始化
	if app.Setting.Output != nil {
		slog.SetOutput(app.Setting.Output)
	}
}

//go:embed static/fingerprint.txt
var fingerprintEmbed embed.FS

const fingerprintPath = "static/fingerprint.txt"

func Initjkscan() {
	//HTTP指纹库初始化
	fs, _ := fingerprintEmbed.Open(fingerprintPath)
	if n, err := appfinger.InitDatabaseFS(fs); err != nil {
		slog.Println(slog.ERROR, "指纹库加载失败，static/fingerprint.txt文件有误", err)
	} else {
		slog.Printf(slog.INFO, "成功加载HTTP指纹:[%d]条", n)
	}
}

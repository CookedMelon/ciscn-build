package app

import (
	"fmt"
	"jkscan/lib/sflag"
	"os"
)

type args struct {
	USAGE, HELP, LOGO, SYNTAX string

	Help, Debug, Check, NotPing          bool
	ScanVersion, DownloadQQwry, CloseCDN bool
	Output, Proxy, Encoding              string
	Port, ExcludedPort                   []int
	Path, Host, Target                   []string
	OutputJson, OutputCSV                string
	Spy, Touch                           string
	Top, Threads                         int
	//输出修饰
	Match, NotMatch string
}

var Args = args{}

// Init 初始化参数
func (o *args) Init() {
	//自定义Usage
	sflag.SetUsage(o.LOGO)
	//定义参数
	o.define()
	//实例化参数值
	sflag.InitArgs()
	//输出LOGO
	o.printLOGO()
}

// 定义参数
func (o *args) define() {
	sflag.BoolVar(&o.Help, "h", false)
	sflag.BoolVar(&o.Help, "help", false)
	sflag.BoolVar(&o.Debug, "debug", false)
	sflag.BoolVar(&o.Debug, "d", false)
	sflag.IntVar(&o.Top, "top", 400)
	sflag.IntVar(&o.Threads, "threads", 800)
	//spy模块
	sflag.AutoVarString(&o.Spy, "spy", "None")
	//jkscan模块
	sflag.StringSpliceVar(&o.Target, "target")
	sflag.StringSpliceVar(&o.Target, "t")
	sflag.IntSpliceVar(&o.Port, "port")
	sflag.IntSpliceVar(&o.Port, "p")
	//输出模块
	sflag.StringVar(&o.OutputJson, "oJ", "")
}

func (o *args) SetLogo(logo string) {
	o.LOGO = logo
}

func (o *args) SetUsage(usage string) {
	o.USAGE = usage
}

func (o *args) SetSyntax(syntax string) {
	o.SYNTAX = syntax
}

func (o *args) SetHelp(help string) {
	o.HELP = help
}

// 校验参数有效性
func (o *args) CheckAvailable() {
	if len(o.Port) > 0 && o.Top != 400 {
		fmt.Print("--port、--top参数不能同时使用")
		os.Exit(0)
	}
	//判断内容
	if o.Top != 0 && (o.Top > 1000 || o.Top < 1) {
		fmt.Print("TOP参数输入错误,TOP参数应为1-1000之间的整数。")
		os.Exit(0)
	}
	if o.Threads != 0 && o.Threads > 2048 {
		fmt.Print("--threads不得大于2048")
		os.Exit(0)
	}
}

// 输出LOGO
func (arg *args) printLOGO() {
	if len(os.Args) == 1 {
		fmt.Print(arg.LOGO)
		fmt.Print(arg.USAGE)
		os.Exit(0)
	}
	if arg.Help {
		fmt.Print(arg.LOGO)
		fmt.Print(arg.USAGE)
		fmt.Print(arg.HELP)
		os.Exit(0)
	}
	//打印logo
	fmt.Print(arg.LOGO)
}

package run

import (
	"fmt"
	"jkscan/app"
	"jkscan/core/scanner"
	"jkscan/core/slog"
	"jkscan/lib/misc"
	"jkscan/lib/uri"
	"net"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lcvvvv/gonmap"

	"github.com/lcvvvv/appfinger"

	"github.com/atotto/clipboard"
)

func Start() {
	//启用定时输出负载情况
	go PrintStatus()
	//下发扫描任务
	var wg = &sync.WaitGroup{}
	wg.Add(3) // 表示有三个扫描器
	Scanner1 = getScanner1(wg)
	Scanner2 = getScanner2(wg)
	Scanner3 = getScanner3(wg)
	//扫描器进入监听状态
	start()
	//开始分发扫描任务
	for _, expr := range app.Setting.Target {
		GetTask(expr)
	}
	//根据扫描情况，关闭scanner
	go stop()
	wg.Wait()
}

func GetTask(expr string) {
	if expr == "" {
		return
	}
	if expr == "paste" || expr == "clipboard" {
		if clipboard.Unsupported {
			slog.Println(slog.ERROR, runtime.GOOS, "clipboard unsupported")
		}
		clipboardStr, _ := clipboard.ReadAll()
		for _, line := range strings.Split(clipboardStr, "\n") {
			line = strings.ReplaceAll(line, "\r", "")
			GetTask(line)
		}
		return
	}

	if uri.IsIPRanger(expr) {
		for _, ip := range uri.RangerToIP(expr) {
			GetTask(ip.String())
		}
		return
	}
	if uri.IsHostPath(expr) {
		if !app.Setting.Check {
			GetTask(uri.GetNetlocWithHostPath(expr))
		}
		return
	}
	if uri.IsIPv4(expr) {
		Scanner1.Push(net.ParseIP(expr))
		return
	}

	if uri.IsNetWithPort(expr) {
		netloc, port := uri.SplitWithNetlocPort(expr)
		if uri.IsIPv4(netloc) {
			Scanner2.Push(net.ParseIP(netloc), port)
		}
		if !app.Setting.Check {
			GetTask(netloc)
		}
		return
	}
	slog.Println(slog.WARN, "无法识别的Target:", expr)
}

var (
	Scanner1 *scanner.IPClient
	Scanner2 *scanner.PortClient
	Scanner3 *scanner.URLClient
)

func start() {
	go Scanner1.Start()
	go Scanner2.Start()
	time.Sleep(time.Second * 1)
	slog.Println(slog.INFO, "准备就绪")
}

func stop() {
	for {
		time.Sleep(time.Second)
		if Scanner1.RunningThreads() == 0 && !Scanner1.IsDone() {
			Scanner1.Stop()
		}
		if !Scanner1.IsDone() {
			continue
		}
		if Scanner2.RunningThreads() == 0 && !Scanner2.IsDone() {
			Scanner2.Stop()
		}
		if !Scanner2.IsDone() {
			continue
		}
		if Scanner3.RunningThreads() == 0 && !Scanner3.IsDone() {
			Scanner3.Stop()
		}
	}
}

func getScanner1(wg *sync.WaitGroup) *scanner.IPClient {
	IPConfig := scanner.DefaultConfig()
	IPConfig.Threads = 200
	IPConfig.Timeout = 200 * time.Millisecond
	IPConfig.HostDiscoverClosed = app.Setting.ClosePing
	client := scanner.NewIPScanner(IPConfig)
	client.HandlerDie = func(addr net.IP) {
		slog.Println(slog.DEBUG, addr.String(), " is die")
	}
	client.HandlerAlive = func(addr net.IP) {
		//启用端口存活性探测任务下发器
		slog.Println(slog.DEBUG, addr.String(), " is alive")
		for _, port := range app.Setting.Port {
			Scanner2.Push(addr, port)
		}
	}
	client.HandlerError = func(addr net.IP, err error) {
		slog.Println(slog.DEBUG, "IPScanner Error: ", addr.String(), err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}

func setTimeout(i int) time.Duration {
	switch {
	case i > 10000:
		return time.Millisecond * 200
	case i > 5000:
		return time.Millisecond * 300
	case i > 1000:
		return time.Millisecond * 400
	default:
		return time.Millisecond * 500
	}
}

func getScanner2(wg *sync.WaitGroup) *scanner.PortClient {
	PortConfig := scanner.DefaultConfig()
	PortConfig.Threads = app.Setting.Threads
	PortConfig.Timeout = setTimeout(len(app.Setting.Port))
	if app.Setting.ScanVersion {
		PortConfig.DeepInspection = true
	}
	client := scanner.NewPortScanner(PortConfig)
	client.HandlerClosed = func(addr net.IP, port int) {
		//nothing
	}
	client.HandlerOpen = func(addr net.IP, port int) {
		printOpenResult(addr, port)
	}
	client.HandlerNotMatched = func(addr net.IP, port int, response string) {
		printUnknownResult(addr, port, response)
	}
	client.HandlerMatched = func(addr net.IP, port int, response *gonmap.Response) {
		URLRaw := fmt.Sprintf("%s://%s:%d", response.FingerPrint.Service, addr.String(), port)
		URL, _ := url.Parse(URLRaw)
		if appfinger.SupportCheck(URL.Scheme) {
			return
		}
		printNmapResult(URL, response)
	}

	client.HandlerError = func(addr net.IP, port int, err error) {
		slog.Println(slog.DEBUG, "PortScanner Error: ", fmt.Sprintf("%s:%d", addr.String(), port), err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}
func getScanner3(wg *sync.WaitGroup) *scanner.URLClient {
	URLConfig := scanner.DefaultConfig()
	URLConfig.Threads = app.Setting.Threads/2 + 1

	client := scanner.NewURLScanner(URLConfig)
	client.HandlerMatched = func(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
		printAppResult(URL, banner, finger)
	}
	client.HandlerError = func(url *url.URL, err error) {
		slog.Println(slog.DEBUG, "URLScanner Error: ", url.String(), err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}
func printNmapResult(URL *url.URL, resp *gonmap.Response) {
	if responseFilter(resp.Raw) {
		return
	}
	finger := resp.FingerPrint
	m := misc.TurnMap(finger)
	m["Response"] = resp.Raw
	m["IP"] = URL.Hostname()
	m["Port"] = URL.Port()
	fmt.Println("normal1")
	fmt.Println(URL)
	fmt.Println(finger.Service)
	// misc.PrintMap(m)
	fmt.Println("------------------------------")
	m["URL"] = URL.String() //nmap扫出的url没有URL
	tmap := misc.GetService(m)
	fmt.Println(tmap)
	tmpService := Service{
		tmap.Port,
		tmap.Protocol,
		tmap.ServiceApp,
	}
	tmpData := Data{
		[]Service{tmpService},
		tmap.DeviceInfo,
		tmap.Honeypot,
	}
	Add(tmap.Ip, tmpData)
	misc.Printinfo(tmap)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++")
	// outputHandler(URL, finger.Service, m)
}

func printAppResult(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
	if responseFilter(banner.Response, banner.Cert) {
		return
	}
	m := misc.TurnMap(finger)
	m["Service"] = URL.Scheme
	m["FoundIP"] = banner.FoundIP
	m["Response"] = banner.Response
	m["Cert"] = banner.Cert
	m["Header"] = banner.Header
	m["Body"] = banner.Body
	m["ICP"] = banner.ICP
	m["FingerPrint"] = m["ProductName"]
	delete(m, "ProductName")
	m["Port"] = uri.GetURLPort(URL)
	if m["Port"] == "" {
		slog.Println(slog.WARN, "无法获取端口号：", URL)
	}
	if hostname := URL.Hostname(); uri.IsIPv4(hostname) {
		m["IP"] = hostname
	}
	fmt.Println("normal2")
	fmt.Println(URL)
	fmt.Println(banner.Title)
	fmt.Println("------------------------------")
	m["URL"] = URL.String() //appfinger扫出的url没有URL
	tmap := misc.GetService(m)
	fmt.Println(tmap)
	tmpService := Service{
		tmap.Port,
		tmap.Protocol,
		tmap.ServiceApp,
	}
	tmpData := Data{
		[]Service{tmpService},
		tmap.DeviceInfo,
		tmap.Honeypot,
	}
	Add(tmap.Ip, tmpData)
	misc.Printinfo(tmap)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++")
}

func printUnknownResult(addr net.IP, port int, response string) {
	if responseFilter(response) {
		return
	}
	//输出结果
	fmt.Println("Unknown")
	target := fmt.Sprintf("unknown://%s:%d", addr.String(), port)
	URL, _ := url.Parse(target)
	fmt.Println(URL)
	fmt.Println("无法识别的协议")
	misc.PrintMap(map[string]string{
		"Response": response,
		"IP":       URL.Hostname(),
		"Port":     strconv.Itoa(port),
	})
	// outputHandler(URL, "无法识别该协议", map[string]string{
	// 	"Response": response,
	// 	"IP":       URL.Hostname(),
	// 	"Port":     strconv.Itoa(port),
	// })
	fmt.Println("------------------------------")
	m := map[string]string{"IP": URL.Hostname(), "Port": strconv.Itoa(port), "Response": response, "URL": URL.String()}
	// m := map[string]interface{}{"IP": URL.Hostname(), "Port": strconv.Itoa(port), "Response": response, "URL": URL.String()}
	tmap := misc.GetService(m)
	fmt.Println(tmap)
	tmpService := Service{
		tmap.Port,
		tmap.Protocol,
		tmap.ServiceApp,
	}
	tmpData := Data{
		[]Service{tmpService},
		tmap.DeviceInfo,
		tmap.Honeypot,
	}
	Add(tmap.Ip, tmpData)
	misc.Printinfo(tmap)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++")
}

func printOpenResult(addr net.IP, port int) {
	// //输出结果
	fmt.Println("empty")
	protocol := gonmap.GuessProtocol(port) //获取协议
	target := fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port)
	URL, _ := url.Parse(target)
	fmt.Println(URL)
	fmt.Println("response is empty2")
	misc.PrintMap(
		map[string]string{
			"IP":   URL.Hostname(),
			"Port": strconv.Itoa(port),
		})
	// outputHandler(URL, "response is empty2", map[string]string{
	// 	"IP":   URL.Hostname(),
	// 	"Port": strconv.Itoa(port),
	// })
	fmt.Println("------------------------------")
	m := map[string]string{"IP": URL.Hostname(), "Port": strconv.Itoa(port), "URL": URL.String()}
	// m := map[string]interface{}{"IP": URL.Hostname(), "Port": strconv.Itoa(port), "URL": URL.String()}
	tmap := misc.GetService(m)
	fmt.Println(tmap)
	tmpService := Service{
		tmap.Port,
		tmap.Protocol,
		tmap.ServiceApp,
	}
	tmpData := Data{
		[]Service{tmpService},
		tmap.DeviceInfo,
		tmap.Honeypot,
	}
	Add(tmap.Ip, tmpData)
	misc.Printinfo(tmap)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++")
}
func responseFilter(strArgs ...string) bool {
	var match = app.Setting.Match
	var notMatch = app.Setting.NotMatch

	if match != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.Match) {
				return false
			}
		}
	}

	if notMatch != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.NotMatch) {
				return true
			}
		}
	}
	return false
}

func PrintStatus() {
	for {
		time.Sleep(time.Second * 1)
		var (
			IPnum   = Scanner1.RunningThreads()
			Portnum = Scanner2.RunningThreads()
		)
		if time.Now().Unix()%180 == 0 {
			warn := fmt.Sprintf("当前存活协程数：IP：%d 个，Port：%d 个", IPnum, Portnum)
			slog.Println(slog.WARN, warn)
		}
	}
}

package run

import (
	"fmt"
	"jkscan/app"
	"jkscan/core/scanner"
	"jkscan/core/slog"
	"jkscan/lib/color"
	"jkscan/lib/misc"
	"jkscan/lib/uri"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/simplehttp"
	"github.com/lcvvvv/stdio/chinese"
)

const file_path = "output.json"

func Start() {
	//启用看门狗函数定时输出负载情况
	go watchDog()
	//下发扫描任务
	var wg = &sync.WaitGroup{}
	wg.Add(3) // 表示有三个扫描器
	IPScanner = generateIPScanner(wg)
	PortScanner = generatePortScanner(wg)
	URLScanner = generateURLScanner(wg)
	//扫描器进入监听状态
	start()
	//开始分发扫描任务
	for _, expr := range app.Setting.Target {
		pushTarget(expr)
	}
	slog.Println(slog.INFO, "所有扫描任务已下发完毕")
	//根据扫描情况，关闭scanner
	go stop()
	wg.Wait()
}

func pushTarget(expr string) {
	if expr == "" {
		return
	}
	if expr == "paste" || expr == "clipboard" {
		if clipboard.Unsupported == true {
			slog.Println(slog.ERROR, runtime.GOOS, "clipboard unsupported")
		}
		clipboardStr, _ := clipboard.ReadAll()
		for _, line := range strings.Split(clipboardStr, "\n") {
			line = strings.ReplaceAll(line, "\r", "")
			pushTarget(line)
		}
		return
	}
	if uri.IsIPv4(expr) {
		IPScanner.Push(net.ParseIP(expr))
		if app.Setting.Check == true {
			pushURLTarget(uri.URLParse("http://"+expr), nil)
			pushURLTarget(uri.URLParse("https://"+expr), nil)
		}
		return
	}
	if uri.IsIPv6(expr) {
		slog.Println(slog.WARN, "暂时不支持IPv6的扫描对象：", expr)
		return
	}
	if uri.IsCIDR(expr) {
		for _, ip := range uri.CIDRToIP(expr) {
			pushTarget(ip.String())
		}
		return
	}
	if uri.IsIPRanger(expr) {
		for _, ip := range uri.RangerToIP(expr) {
			pushTarget(ip.String())
		}
		return
	}
	if uri.IsHostPath(expr) {
		pushURLTarget(uri.URLParse("http://"+expr), nil)
		pushURLTarget(uri.URLParse("https://"+expr), nil)
		if app.Setting.Check == false {
			pushTarget(uri.GetNetlocWithHostPath(expr))
		}
		return
	}
	if uri.IsNetlocPort(expr) {
		netloc, port := uri.SplitWithNetlocPort(expr)
		if uri.IsIPv4(netloc) {
			PortScanner.Push(net.ParseIP(netloc), port)
		}
		if app.Setting.Check == false {
			pushTarget(netloc)
		}
		return
	}
	if uri.IsURL(expr) {
		pushURLTarget(uri.URLParse(expr), nil)
		if app.Setting.Check == false {
			pushTarget(uri.GetNetlocWithURL(expr))
		}
		return
	}
	slog.Println(slog.WARN, "无法识别的Target字符串:", expr)
}

func pushURLTarget(URL *url.URL, response *gonmap.Response) {
	var cli *http.Client
	//判断是否初始化client
	if app.Setting.Proxy != "" || app.Setting.Timeout != 3*time.Second {
		cli = simplehttp.NewClient()
	}
	//判断是否需要设置代理
	if app.Setting.Proxy != "" {
		simplehttp.SetProxy(cli, app.Setting.Proxy)
	}
	//判断是否需要设置超时参数
	if app.Setting.Timeout != 3*time.Second {
		simplehttp.SetTimeout(cli, app.Setting.Timeout)
	}

	//判断是否存在请求修饰性参数
	if len(app.Setting.Host) == 0 && len(app.Setting.Path) == 0 {
		URLScanner.Push(URL, response, nil, cli)
		return
	}

	//如果存在，则逐一建立请求下发队列
	var reqs []*http.Request
	for _, host := range app.Setting.Host {
		req, _ := simplehttp.NewRequest(http.MethodGet, URL.String(), nil)
		req.Host = host
		reqs = append(reqs, req)
	}
	for _, path := range app.Setting.Path {
		req, _ := simplehttp.NewRequest(http.MethodGet, URL.String()+path, nil)
		reqs = append(reqs, req)
	}
	for _, req := range reqs {
		URLScanner.Push(req.URL, response, req, cli)
	}
}

var (
	IPScanner   *scanner.IPClient
	PortScanner *scanner.PortClient
	URLScanner  *scanner.URLClient
)

func start() {
	go IPScanner.Start()
	go PortScanner.Start()
	go URLScanner.Start()
	time.Sleep(time.Second * 1)
	slog.Println(slog.INFO, "准备就绪")
}

func stop() {
	for {
		time.Sleep(time.Second)
		if IPScanner.RunningThreads() == 0 && IPScanner.IsDone() == false {
			IPScanner.Stop()
			slog.Println(slog.DEBUG, "检测到所有IP检测任务已完成，IP扫描引擎已停止")
		}
		if IPScanner.IsDone() == false {
			continue
		}
		if PortScanner.RunningThreads() == 0 && PortScanner.IsDone() == false {
			PortScanner.Stop()
			slog.Println(slog.DEBUG, "检测到所有Port检测任务已完成，Port扫描引擎已停止")
		}
		if PortScanner.IsDone() == false {
			continue
		}
		if URLScanner.RunningThreads() == 0 && URLScanner.IsDone() == false {
			URLScanner.Stop()
			slog.Println(slog.DEBUG, "检测到所有URL检测任务已完成，URL扫描引擎已停止")
		}
	}
}

func generateIPScanner(wg *sync.WaitGroup) *scanner.IPClient {
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
			PortScanner.Push(addr, port)
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

func getTimeout(i int) time.Duration {
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

func generatePortScanner(wg *sync.WaitGroup) *scanner.PortClient {
	PortConfig := scanner.DefaultConfig()
	PortConfig.Threads = app.Setting.Threads
	PortConfig.Timeout = getTimeout(len(app.Setting.Port))
	if app.Setting.ScanVersion == true {
		PortConfig.DeepInspection = true
	}
	client := scanner.NewPortScanner(PortConfig)
	client.HandlerClosed = func(addr net.IP, port int) {
		//nothing
	}
	client.HandlerOpen = func(addr net.IP, port int) {
		outputOpenResponse(addr, port)
	}
	client.HandlerNotMatched = func(addr net.IP, port int, response string) {
		outputUnknownResponse(addr, port, response)
	}
	client.HandlerMatched = func(addr net.IP, port int, response *gonmap.Response) {
		URLRaw := fmt.Sprintf("%s://%s:%d", response.FingerPrint.Service, addr.String(), port)
		URL, _ := url.Parse(URLRaw)
		if appfinger.SupportCheck(URL.Scheme) == true {
			pushURLTarget(URL, response)
			return
		}
		outputNmapFinger(URL, response)
	}

	client.HandlerError = func(addr net.IP, port int, err error) {
		slog.Println(slog.DEBUG, "PortScanner Error: ", fmt.Sprintf("%s:%d", addr.String(), port), err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}

func generateURLScanner(wg *sync.WaitGroup) *scanner.URLClient {
	URLConfig := scanner.DefaultConfig()
	URLConfig.Threads = app.Setting.Threads/2 + 1

	client := scanner.NewURLScanner(URLConfig)
	client.HandlerMatched = func(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
		outputAppFinger(URL, banner, finger)
	}
	client.HandlerError = func(url *url.URL, err error) {
		slog.Println(slog.DEBUG, "URLScanner Error: ", url.String(), err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}

//	func outputNmapFinger(URL *url.URL, resp *gonmap.Response) {
//		if responseFilter(resp.Raw) == true {
//			return
//		}
//		finger := resp.FingerPrint
//		m := misc.ToMap(finger)
//		m["Response"] = resp.Raw
//		m["IP"] = URL.Hostname()
//		m["Port"] = URL.Port()
//		fmt.Println("normal1")
//		outputHandler(URL, finger.Service, m)
//	}
func outputNmapFinger(URL *url.URL, resp *gonmap.Response) {
	if responseFilter(resp.Raw) == true {
		return
	}
	finger := resp.FingerPrint
	m := misc.ToMap(finger)
	m["Response"] = resp.Raw
	m["IP"] = URL.Hostname()
	m["Port"] = URL.Port()
	fmt.Println("normal1")
	fmt.Println(URL)
	fmt.Println(finger.Service)
	misc.PrintMap(m)
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

//	func outputAppFinger(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
//		if responseFilter(banner.Response, banner.Cert) == true {
//			return
//		}
//		m := misc.ToMap(finger)
//		m["Service"] = URL.Scheme
//		m["FoundIP"] = banner.FoundIP
//		m["Response"] = banner.Response
//		m["Cert"] = banner.Cert
//		m["Header"] = banner.Header
//		m["Body"] = banner.Body
//		m["ICP"] = banner.ICP
//		m["FingerPrint"] = m["ProductName"]
//		delete(m, "ProductName")
//		m["Port"] = uri.GetURLPort(URL)
//		if m["Port"] == "" {
//			slog.Println(slog.WARN, "无法获取端口号：", URL)
//		}
//		if hostname := URL.Hostname(); uri.IsIPv4(hostname) {
//			m["IP"] = hostname
//		}
//		fmt.Println("normal2")
//		outputHandler(URL, banner.Title, m)
//	}
func outputAppFinger(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
	if responseFilter(banner.Response, banner.Cert) == true {
		return
	}
	m := misc.ToMap(finger)
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
	misc.PrintMap(m)
	// outputHandler(URL, banner.Title, m)
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

//	func outputUnknownResponse(addr net.IP, port int, response string) {
//		if responseFilter(response) == true {
//			return
//		}
//		//输出结果
//		fmt.Println("Unknown")
//		target := fmt.Sprintf("unknown://%s:%d", addr.String(), port)
//		URL, _ := url.Parse(target)
//		outputHandler(URL, "无法识别该协议", map[string]string{
//			"Response": response,
//			"IP":       URL.Hostname(),
//			"Port":     strconv.Itoa(port),
//		})
//	}
func outputUnknownResponse(addr net.IP, port int, response string) {
	if responseFilter(response) == true {
		return
	}
	//输出结果
	fmt.Println("Unknown")
	target := fmt.Sprintf("unknown://%s:%d", addr.String(), port)
	URL, _ := url.Parse(target)
	fmt.Println(URL)
	fmt.Println("无法识别该协议")
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

//	func outputOpenResponse(addr net.IP, port int) {
//		// //输出结果
//		fmt.Println("empty")
//		protocol := gonmap.GuessProtocol(port) //获取协议
//		target := fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port)
//		URL, _ := url.Parse(target)
//		outputHandler(URL, "response is empty2", map[string]string{
//			"IP":   URL.Hostname(),
//			"Port": strconv.Itoa(port),
//		})
//	}
func outputOpenResponse(addr net.IP, port int) {
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
			if strings.Contains(str, app.Setting.Match) == true {
				return false
			}
		}
	}

	if notMatch != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.NotMatch) == true {
				return true
			}
		}
	}
	return false
}

var (
	disableKey       = []string{"MatchRegexString", "Service", "ProbeName", "Response", "Cert", "Header", "Body", "IP"}
	importantKey     = []string{"ProductName", "DeviceType"}
	varyImportantKey = []string{"Hostname", "FingerPrint", "ICP"}
)

func getHTTPDigest(s string) string {
	var length = 24
	var digestBuf []rune
	_, body := simplehttp.SplitHeaderAndBody(s)
	body = chinese.ToUTF8(body)
	for _, r := range []rune(body) {
		buf := []byte(string(r))
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if len(buf) > 1 {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(body, length-len(digestBuf))
}

func getRawDigest(s string) string {
	var length = 24
	if len(s) < length {
		return s
	}
	var digestBuf []rune
	for _, r := range []rune(s) {
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if 0x20 <= r && r <= 0x7E {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(s, length-len(digestBuf))
}

func outputHandler(URL *url.URL, keyword string, m map[string]string) {
	// fmt.Println(m)
	m = misc.FixMap(m)
	if respRaw := m["Response"]; respRaw != "" {
		if m["Service"] == "http" || m["Service"] == "https" {
			m["Digest"] = strconv.Quote(getHTTPDigest(respRaw))
		} else {
			m["Digest"] = strconv.Quote(getRawDigest(respRaw))
		}
	}
	m["Length"] = strconv.Itoa(len(m["Response"]))
	sourceMap := misc.CloneMap(m)
	for _, keyword := range disableKey {
		delete(m, keyword)
	}
	for key, value := range m {
		if key == "FingerPrint" {
			continue
		}
		m[key] = misc.StrRandomCut(value, 24)
	}
	fingerPrint := color.StrMapRandomColor(m, true, importantKey, varyImportantKey)
	fingerPrint = misc.FixLine(fingerPrint)
	format := "%-30v %-" + strconv.Itoa(misc.AutoWidth(color.Clear(keyword), 26+color.Count(keyword))) + "v %s"
	printStr := fmt.Sprintf(format, URL.String(), keyword, fingerPrint)
	slog.Println(slog.DATA, printStr)
	// 输出
	sourceMap["URL"] = URL.String()
	sourceMap["Keyword"] = keyword
	misc.PrintMap(m)

	// fmt.Println(tmap)
	if jw := app.Setting.OutputJson; jw != nil {
		sourceMap["URL"] = URL.String()
		sourceMap["Keyword"] = keyword
		jw.Push(misc.TidyMap(sourceMap))
	}
	if cw := app.Setting.OutputCSV; cw != nil {
		sourceMap["URL"] = URL.String()
		sourceMap["Keyword"] = keyword
		delete(sourceMap, "Header")
		delete(sourceMap, "Cert")
		delete(sourceMap, "Response")
		delete(sourceMap, "Body")
		sourceMap["Digest"] = strconv.Quote(sourceMap["Digest"])
		for key, value := range sourceMap {
			sourceMap[key] = chinese.ToUTF8(value)
		}
		cw.Push(sourceMap)
	}
}

func watchDog() {
	for {
		time.Sleep(time.Second * 1)
		var (
			nIP   = IPScanner.RunningThreads()
			nPort = PortScanner.RunningThreads()
			nURL  = URLScanner.RunningThreads()
		)
		if time.Now().Unix()%180 == 0 {
			warn := fmt.Sprintf("当前存活协程数：IP：%d 个，Port：%d 个，URL：%d 个", nIP, nPort, nURL)
			slog.Println(slog.WARN, warn)
		}
	}
}

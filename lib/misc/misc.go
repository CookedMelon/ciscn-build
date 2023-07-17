package misc

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

var protocol_list = []string{
	"ssh",
	"http",
	"https",
	"rtsp",
	"ftp",
	"telnet",
	"amqp",
	"mongodb",
	"redis",
	"mysql",
}

type Info struct {
	Ip         string
	Port       int
	Protocol   string
	ServiceApp []string
	DeviceInfo []string
	Honeypot   []string
}

func Printinfo(info Info) {
	fmt.Println("ip:", info.Ip)
	fmt.Println("port:", info.Port)
	fmt.Println("protocol:", info.Protocol)
	fmt.Println("service_app:", info.ServiceApp)
	fmt.Println("deviceinfo:", info.DeviceInfo)
	fmt.Println("honeypot:", info.Honeypot)
}
func IsDuplicate[T any](slice []T, val T) bool {
	for _, item := range slice {
		if fmt.Sprint(item) == fmt.Sprint(val) {
			return true
		}
	}
	return false
}

func RemoveDuplicateElement[T any](slice []T, elems ...T) []T {
	slice = append(slice, elems...)
	set := make(map[string]struct{}, len(slice))
	j := 0
	for _, v := range slice {
		_, ok := set[fmt.Sprint(v)]
		if ok {
			continue
		}
		set[fmt.Sprint(v)] = struct{}{}
		slice[j] = v
		j++
	}
	return slice[:j]
}

func FixLine(line string) string {
	line = strings.ReplaceAll(line, "\t", "")
	line = strings.ReplaceAll(line, "\r", "")
	line = strings.ReplaceAll(line, "\n", "")
	line = strings.ReplaceAll(line, "\xc2\xa0", "")
	line = strings.ReplaceAll(line, " ", "")
	return line
}

func Xrange(args ...int) []int {
	var start, stop int
	var step = 1
	var r []int
	switch len(args) {
	case 1:
		stop = args[0]
		start = 0
	case 2:
		start, stop = args[0], args[1]
	case 3:
		start, stop, step = args[0], args[1], args[2]
	default:
		return nil
	}
	if start > stop {
		return nil
	}
	if step < 0 {
		return nil
	}

	for i := start; i <= stop; i += step {
		r = append(r, i)
	}
	return r
}

func MustLength(s string, i int) string {
	if len(s) > i {
		return s[:i]
	}
	return s
}

func Percent(int1 int, int2 int) string {
	float1 := float64(int1)
	float2 := float64(int2)
	f := 1 - float1/float2
	f = f * 100
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func StrRandomCut(s string, length int) string {
	sRune := []rune(s)
	if len(sRune) > length {
		i := rand.Intn(len(sRune) - length)
		return string(sRune[i : i+length])
	} else {
		return s
	}
}
func removeEmptyStrings(slice []string) []string {
	newSlice := make([]string, 0)
	for _, str := range slice {
		if str != "" {
			newSlice = append(newSlice, str)
		}
	}
	return newSlice
}
func Base64Encode(keyword string) string {
	input := []byte(keyword)
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

func Base64Decode(encodeString string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	return string(decodeBytes), err
}

func CloneStrMap(strMap map[string]string) map[string]string {
	newStrMap := make(map[string]string)
	for k, v := range strMap {
		newStrMap[k] = v
	}
	return newStrMap
}

func CloneIntMap(intMap map[int]string) map[int]string {
	newIntMap := make(map[int]string)
	for k, v := range intMap {
		newIntMap[k] = v
	}
	return newIntMap
}

func RandomString(i ...int) string {
	var length int
	var str string
	if len(i) != 1 {
		length = 32
	} else {
		length = i[0]
	}
	Char := "01234567890abcdef"
	for range Xrange(length) {
		j := rand.Intn(len(Char) - 1)
		str += Char[j : j+1]
	}
	return str
}

func Intersection(a []string, b []string) (inter []string) {
	for _, s1 := range a {
		for _, s2 := range b {
			if s1 == s2 {
				inter = append(inter, s1)
			}
		}
	}
	return inter
}

func FixMap(m map[string]string) map[string]string {
	var arr []string
	rm := make(map[string]string)
	for key, value := range m {
		if value == "" {
			continue
		}
		if IsDuplicate(arr, value) {
			if key != "Username" && key != "Password" {
				continue
			}
		}
		arr = append(arr, value)
		rm[key] = value
	}
	return rm
}

func CloneMap(m map[string]string) map[string]string {
	var nm = make(map[string]string)
	for key, value := range m {
		nm[key] = value
	}
	return nm
}

func AutoWidth(s string, length int) int {
	length1 := len(s)
	length2 := len([]rune(s))

	if length1 == length2 {
		return length
	}
	return length - (length1-length2)/2
}

func ToMap(param interface{}) map[string]string {
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	m := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		// 通过interface方法来获取key所对应的值
		var cell string
		switch s := v.Field(i).Interface().(type) {
		case string:
			cell = s
		case []string:
			cell = strings.Join(s, "; ")
		case int:
			cell = strconv.Itoa(s)
		case Stringer:
			cell = s.String()
		default:
			continue
		}
		m[t.Field(i).Name] = cell
	}
	return m
}

type Stringer interface {
	String() string
}

func CopySlice[T any](slice []T) []T {
	v := make([]T, len(slice))
	copy(v, slice)
	return v
}

func TidyMap(m map[string]string) map[string]string {
	var nm = make(map[string]string)
	for key, value := range m {
		if value == "" || key == "MatchRegexString" || key == "Response" || key == "Body" || key == "Cert" || key == "Header" {
			continue
		}
		nm[key] = value
	}
	return nm
}
func PrintMap(m map[string]string) {
	fmt.Println("map:")
	for key, value := range m {
		fmt.Println(key, ":", value)
	}
}
func GetProtocol(m map[string]string) string {
	var protocol string
	// 获取m的“URL”字段的://前的字符内容
	if strings.Contains(m["URL"], "://") {
		protocol = strings.Split(m["URL"], "://")[0]
		//如果protocol在protocol_list中，则返回protocol
		if IsDuplicate(protocol_list, protocol) {
			return protocol
		} else {
			return ""
		}
	}
	return ""
}

// 返回值{"ip","port","protocol","service_app","deviceinfo","honeypot"}
func GetService(m map[string]string) Info {
	var info Info
	var service_app = make([]string, 0)
	//将X-Powered-By，Response，Server，FingerPrint合并为一个字符串
	var detail_all_in_one = strings.ToLower(m["X-Powered-By"] + m["Response"] + m["Server"] + m["FingerPrint"])
	//识别openssh
	//将response中的SSH-2.0-OpenSSH_8.0转为openssh/8.0，如果版本号匹配失败(没匹配到SSH-2.0-OpenSSH_)则为openssh/N
	//ssh之后得改下，不能用split这种
	if strings.Contains(detail_all_in_one, "openssh") {
		service_app = append(service_app, GetOpenSSH(m))
	}
	// 识别wordpress，从Body中提取content="WordPress 6.0.5"字样为workpress/6.0.5，若未匹配到则为workpress/N
	if strings.Contains(detail_all_in_one, "wordpress") {
		service_app = append(service_app, GetWordPress(m))
	}

	//识别windows,FingerPrint中出现即可，出现Windows CE 6.00等记为windows/6.00，否则记为windows/N
	if strings.Contains(detail_all_in_one, "windows") {
		service_app = append(service_app, GetWindows(m))
	}

	//识别nginx
	//目前没有看到有信息能指出nginx的版本，有些时候可能能从body里找到
	if strings.Contains(detail_all_in_one, "nginx") {
		service_app = append(service_app, GetNginx(m))
	}

	//识别Jetty
	//对Jetty(9.4.11.v20180605)记为Jetty/9.4.11.v20180605，对Jetty记为Jetty/N
	if strings.Contains(detail_all_in_one, "jetty") {
		service_app = append(service_app, GetJetty(m))
	}

	//识别debian
	//若server中有Debian则记为debian/N，若使用正则识别出response中出现Debian-5则记为debian/5
	if strings.Contains(detail_all_in_one, "debian") {
		service_app = append(service_app, GetDebian(m))
	}
	//识别Grafana v9.1.2
	//若FingerPrint中有grafana则记为grafana/N，若使用正则识别出Body中出现grafana v9.1.2等字样则记为grafana/v9.1.2
	if strings.Contains(detail_all_in_one, "grafana") {
		service_app = append(service_app, GetGrafana(m))
	}

	//玛德写完才发现sql不在考察范围内
	// //识别mysql
	// //将Version中的5.5.68-MariaDB记为5.5.68
	// if strings.Contains(m["Service"], "mysql") {
	// 	re := regexp.MustCompile(`^[\d\.]+`)
	// 	match := re.FindString(m["Version"])
	// 	if match != "" {
	// 		answer = append(answer, "mysql/"+match)
	// 	} else {
	// 		answer = append(answer, "mysql/N")
	// 	}
	// }

	//识别node.js
	//若FingerPrint中有Node.js则记为node.js/N，若使用正则识别出Body中出现node.js v9.1.2等字样则记为node.js/9.1.2
	if strings.Contains(detail_all_in_one, "node.js") {
		service_app = append(service_app, GetNodeJS(m))
	}
	//识别express
	//若FingerPrint中有Express则记为express/N
	if strings.Contains(detail_all_in_one, "express") {
		service_app = append(service_app, "express/N")
	}
	// if strings.Contains(m["FingerPrint"], "Express") {
	// 	service_app = append(service_app, "express/N")
	// }
	//识别asp.net
	//若X-Powered-By中有ASP.NET则记为asp.net/N
	if strings.Contains(detail_all_in_one, "asp.net") {
		service_app = append(service_app, "asp.net/N")
	}
	// if strings.Contains(m["X-Powered-By"], "ASP.NET") {
	// 	service_app = append(service_app, "asp.net/N")
	// }
	//识别php
	//若Server中有PHP/5.4.16，记为php/5.4.16
	if strings.Contains(detail_all_in_one, "php") {
		service_app = append(service_app, GetPHP(m))
	}

	//识别Microsoft-HTTPAPI
	//若Server中有Microsoft-HTTPAPI/2.0则记为microsoft-httpapi/2.0
	if strings.Contains(detail_all_in_one, strings.ToLower("Microsoft-HTTPAPI")) {
		service_app = append(service_app, GetHttpAPI(m))
	}

	//识别apache
	//若Server中有Apache/2.4.29 (Ubuntu)则记为apache/2.4.29
	if strings.Contains(detail_all_in_one, strings.ToLower("apache")) {
		service_app = append(service_app, GetApache(m))
	}

	//识别OpenResty
	//若Server中有OpenResty
	if strings.Contains(detail_all_in_one, "openresty") {
		service_app = append(service_app, "openresty/N")
	}
	// if strings.Contains(m["Server"], "OpenResty") {
	// 	service_app = append(service_app, "openresty/N")
	// }
	//识别IIS
	//若Server中有IIS/10.0(Microsoft-IIS/10.0)则记为iis/10.0
	if strings.Contains(detail_all_in_one, strings.ToLower("Microsoft-IIS")) {
		service_app = append(service_app, GetIIS(m))
	}

	//识别OpenSSL
	//若Server中有OpenSSL/1.1.1，则记为openssl/1.1.1
	if strings.Contains(detail_all_in_one, strings.ToLower("openssl")) {
		service_app = append(service_app, GetOpenSSL(m))
	}
	//识别elasticsearch
	//若Body中出现Elasticsearch，则记为elasticsearch/N
	if strings.Contains(detail_all_in_one, "elasticsearch") {
		service_app = append(service_app, "elasticsearch/N")
	}
	//识别LiteSpeed
	if strings.Contains(detail_all_in_one, "litespeed") {
		service_app = append(service_app, "litespeed/N")
	}
	//识别rabbitmq
	if strings.Contains(detail_all_in_one, "rabbitmq") {
		service_app = append(service_app, "rabbitmq/N")
	}
	//识别micro_httpd
	if strings.Contains(detail_all_in_one, "micro_httpd") {
		service_app = append(service_app, "micro_httpd/N")
	}
	//识别grafana
	if strings.Contains(detail_all_in_one, "grafana") {
		service_app = append(service_app, "grafana/N")
	}
	//识别Weblogic
	if strings.Contains(detail_all_in_one, "weblogic") {
		service_app = append(service_app, "weblogic/N")
	}
	//识别java
	if strings.Contains(detail_all_in_one, "java") && strings.Contains(detail_all_in_one, "javarmi") && strings.Contains(detail_all_in_one, "javadoc") && strings.Contains(detail_all_in_one, "javaex") && strings.Contains(detail_all_in_one, "javascript") {
		service_app = append(service_app, "java/N")
	}
	// 识别ubuntu
	if strings.Contains(detail_all_in_one, "ubuntu") {
		service_app = append(service_app, GetUbuntu(m))
	}
	// 识别centos
	if strings.Contains(detail_all_in_one, "centos") {
		service_app = append(service_app, "centos/N")
	}
	info.ServiceApp = removeEmptyStrings(service_app)
	// 识别蜜罐
	honeypot := make([]string, 0)
	honeypot_list := []string{"glastopf", "kippo", "hfish"}
	//循环遍历，若匹配成功添加端口/蜜罐名
	for _, honeypot_name := range honeypot_list {
		if strings.Contains(detail_all_in_one, honeypot_name) {
			honeypot = append(honeypot, m["Port"])
			honeypot = append(honeypot, honeypot_name)
		}
	}
	info.Honeypot = honeypot
	//识别设备信息
	deviceinfo := make([]string, 0)
	if strings.Contains(detail_all_in_one, "pfsense") {
		deviceinfo = append(deviceinfo, "firewall/pfsense")
	} else if strings.Contains(detail_all_in_one, "hikvision-cameras") {
		deviceinfo = append(deviceinfo, "webcam/hikvision")
	} else if strings.Contains(detail_all_in_one, "dahua-cameras") {
		deviceinfo = append(deviceinfo, "webcam/dahua")
	} else if strings.Contains(detail_all_in_one, "cisco-switch") {
		deviceinfo = append(deviceinfo, "switch/cisco")
	} else if strings.Contains(detail_all_in_one, "synology-sAS") {
		deviceinfo = append(deviceinfo, "nas/synology")
	}
	info.DeviceInfo = deviceinfo

	// 识别IP地址
	info.Ip = m["IP"]
	// 识别端口
	info.Port, _ = strconv.Atoi(m["Port"])
	// 识别协议
	info.Protocol = GetProtocol(m)

	return info
}

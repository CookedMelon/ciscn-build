package misc

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
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
func GetService(m map[string]string) []string {
	var answer = make([]string, 0)
	//识别openssh
	//将response中的SSH-2.0-OpenSSH_8.0转为openssh/8.0，如果版本号匹配失败(没匹配到SSH-2.0-OpenSSH_)则为openssh/N
	//ssh之后得改下，不能用split这种
	re := regexp.MustCompile(`SSH-2\.0-OpenSSH_([\d\.]+)`)
	match := re.FindStringSubmatch(m["Response"])
	if match != nil {
		fmt.Println("find openssh:", m["Response"])
		if len(match) < 2 {
			if regexp.MustCompile(`OpenSSH`).FindString(m["Response"]) != "" {
				answer = append(answer, "openssh/N")
			}
			answer = append(answer, "openssh/N")
		}
		// 如果匹配到 "SSH-2.0-OpenSSH_(数字小数点组合)"，则返回 "openssh/(数字小数点组合)"
		answer = append(answer, "openssh/"+match[1])
	}
	if strings.Contains(m["Response"], "OpenSSH") && !strings.Contains(m["Response"], "SSH-2.0-OpenSSH_") {
		fmt.Println("find openssh但是无版本:", m["Response"])
		answer = append(answer, "openssh/N")
	}

	// fmt.Println("find openssh:", answer)
	// 识别wordpress，从Body中提取content="WordPress 6.0.5"字样为workpress/6.0.5，若未匹配到则为workpress/N
	if strings.Contains(m["FingerPrint"], "WordPress") {
		fmt.Println("find wordpress:", m["FingerPrint"])
		re := regexp.MustCompile(`WordPress (\d+\.\d+(\.\d+){0,2})`)
		matches := re.FindAllString(m["Body"], -1)

		if len(matches) > 0 {
			answer = append(answer, strings.ToLower(matches[0]))
		} else {
			answer = append(answer, "wordpress/N")
		}
	}
	//识别windows,FingerPrint中出现即可，出现Windows CE 6.00等记为windows/6.00，否则记为windows/N
	if strings.Contains(m["FingerPrint"], "Windows") {
		fmt.Println("find windows:", m["FingerPrint"])
		re := regexp.MustCompile(`Windows CE (\d+\.\d+)`)
		matches := re.FindAllString(m["FingerPrint"], -1)
		if len(matches) > 0 {
			var temp = strings.ToLower(matches[0])
			answer = append(answer, strings.Replace(temp, " ce ", "/", 1))
		} else {
			answer = append(answer, "windows/N")
		}
	}

	//识别nginx
	//目前没有看到有信息能指出nginx的版本，有些时候可能能从body里找到
	if strings.Contains(m["FingerPrint"], "nginx") {
		answer = append(answer, "nginx/N")
	} else if strings.Contains(m["Server"], "nginx") {
		//如Server: nginx/1.18.0 (Ubuntu)提取为nginx/1.18.0，如Server: nginx/1.28.0提取为nginx/1.28.0
		re := regexp.MustCompile(`nginx/(\d+\.\d+(\.\d+){0,2})`)
		matches := re.FindAllString(m["Server"], -1)
		if len(matches) > 0 {
			answer = append(answer, strings.ToLower(matches[0]))
		} else {
			answer = append(answer, "nginx/N")
		}
	}
	//识别Jetty
	//对Jetty(9.4.11.v20180605)记为Jetty/9.4.11.v20180605，对Jetty记为Jetty/N
	if strings.Contains(m["Server"], "Jetty") {
		re := regexp.MustCompile(`Jetty(\((.*?)\))?`)
		matches := re.FindStringSubmatch(m["Server"])
		if matches[2] == "" { // 如果括号内无任何文本，则替换为 "Jetty/N"
			answer = append(answer, "jetty/N")
		} else {
			answer = append(answer, "jetty/"+matches[2])
		}
	}
	//识别debian
	//若server中有Debian则记为debian/N，若使用正则识别出response中出现Debian-5则记为debian/5
	re = regexp.MustCompile(`Debian-(\d+)`)
	match = re.FindStringSubmatch(m["Response"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Debian-数字"
		answer = append(answer, "debian/"+match[1])
	} else if strings.Contains(m["Server"], "Debian") {
		answer = append(answer, "debian/N")
	} else {
		match = re.FindStringSubmatch(m["Response"])
		if len(match) > 0 && match[1] != "" { // 如果匹配到 "Debian-数字"
			answer = append(answer, "debian/"+match[1])
		} else if strings.Contains(m["Response"], "Debian") || strings.Contains(m["Response"], "debian") {
			answer = append(answer, "debian/N")
		}
	}
	//识别Grafana v9.1.2
	//若FingerPrint中有grafana则记为grafana/N，若使用正则识别出Body中出现grafana v9.1.2等字样则记为grafana/v9.1.2
	if strings.Contains(m["FingerPrint"], "Grafana") {
		re := regexp.MustCompile(`Grafana (\d+\.\d+(\.\d+){0,2})`)
		matches := re.FindAllString(m["Body"], -1)
		if len(matches) > 0 {
			answer = append(answer, strings.ToLower(matches[0]))
		} else {
			answer = append(answer, "grafana/N")
		}
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
	if strings.Contains(m["FingerPrint"], "Node.js") {
		re := regexp.MustCompile(`Node.js (\d+\.\d+(\.\d+){0,2})`)
		matches := re.FindAllString(m["Body"], -1)
		if len(matches) > 0 {
			answer = append(answer, strings.ToLower(matches[0]))
		} else {
			answer = append(answer, "node.js/N")
		}
	}
	//识别express
	//若FingerPrint中有Express则记为express/N
	if strings.Contains(m["FingerPrint"], "Express") {
		answer = append(answer, "express/N")
	}
	//识别asp.net
	//若X-Powered-By中有ASP.NET则记为asp.net/N
	if strings.Contains(m["X-Powered-By"], "ASP.NET") {
		answer = append(answer, "asp.net/N")
	}
	//识别php
	//若Server中有PHP/5.4.16，记为php/5.4.16
	re = regexp.MustCompile(`(?i)PHP/([\d\.]+)`)
	match = re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "PHP/数字"
		answer = append(answer, "php/"+match[1])
	} else if strings.Contains(m["Server"], "PHP") {
		answer = append(answer, "php/N")
	}
	//识别Microsoft-HTTPAPI
	//若Server中有Microsoft-HTTPAPI/2.0则记为microsoft-httpapi/2.0
	re = regexp.MustCompile(`(?i)Microsoft-HTTPAPI/([\d\.]+)`)
	match = re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Microsoft-HTTPAPI/数字"
		answer = append(answer, "microsoft-httpapi/"+match[1])
	} else if strings.Contains(m["Server"], "Microsoft-HTTPAPI") {
		answer = append(answer, "microsoft-httpapi/N")
	}
	//识别apache
	//若Server中有Apache/2.4.29 (Ubuntu)则记为apache/2.4.29
	re = regexp.MustCompile(`(?i)Apache/([\d\.]+)`)
	match = re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Apache/数字"
		answer = append(answer, "apache/"+match[1])
	} else if strings.Contains(m["Server"], "Apache") && !strings.Contains(m["Server"], "Apache-Coyote") {
		answer = append(answer, "apache/N")
	}
	//识别OpenResty
	//若Server中有OpenResty
	if strings.Contains(m["Server"], "OpenResty") {
		answer = append(answer, "openresty/N")
	}
	//识别IIS
	//若Server中有IIS/10.0(Microsoft-IIS/10.0)则记为iis/10.0
	re = regexp.MustCompile(`(?i)IIS/([\d\.]+)`)
	match = re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "iis/数字"
		answer = append(answer, "iis/"+match[1])
	} else if strings.Contains(m["Server"], "iis") {
		answer = append(answer, "iis/N")
	}
	//识别OpenSSL
	//若Server中有OpenSSL/1.1.1，则记为openssl/1.1.1
	re = regexp.MustCompile(`(?i)OpenSSL/([\d\.]+)`)
	match = re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "openssl/数字"
		answer = append(answer, "openssl/"+match[1])
	} else if strings.Contains(m["Server"], "openssl") {
		answer = append(answer, "openssl/N")
	}
	//识别elasticsearch
	//若Body中出现Elasticsearch，则记为elasticsearch/N
	if strings.Contains(m["Body"], "Elasticsearch") {
		answer = append(answer, "elasticsearch/N")
	}
	return answer
}

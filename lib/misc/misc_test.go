package misc

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestFixLine(t *testing.T) {

	var s = "1 1  1          1"
	fmt.Println(s)
	fmt.Println(FixLine(s))

}
func TestWordPress(t *testing.T) {
	var m = map[string]string{}
	m["FingerPrint"] = "WordPress\t; WordPress-PluginWooCommerce\t; PHP\t; Openfire\t; Apache\t; PoweredBy\t; MetaGenerator\t; JQuery\t; Open-Graph-Protocol\t; Frame\t; wordpress\t; Apache httpd/2.4.29; Apache httpd;  v"
	m["Body"] = `Some random text with content="WordPress 6.0.5.2" and other versions like content="WordPress 5.4" or content="WordPress 7.0.1.2".`
	var answer = make([]string, 0)
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
	fmt.Println(answer)
}
func TestWindows(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	m["FingerPrint"] = "Windows CE 6.00"
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
	fmt.Println(answer)
}
func TestNginx(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	// m["Server"] = "nginx/1.18.0 (Ubuntu)"
	m["Server"] = "nginx/1.0"
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
	fmt.Println(answer)
}
func TestJetty(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	// m["Server"] = "nginx/1.18.0 (Ubuntu)"
	// m["Server"] = "Jetty"
	m["Server"] = "Jetty(9.4.11.v20180"
	if strings.Contains(m["Server"], "Jetty") {
		re := regexp.MustCompile(`Jetty(\((.*?)\))?`)
		matches := re.FindStringSubmatch(m["Server"])
		if matches[2] == "" { // 如果括号内无任何文本，则替换为 "Jetty/N"
			answer = append(answer, "jetty/N")
		} else {
			answer = append(answer, "jetty/"+matches[2])
		}
	}
	fmt.Println(answer)
}
func TestDebian(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	// m["Response"] = "Debian-10"
	m["Server"] = "(Debian)"
	re := regexp.MustCompile(`Debian-(\d+)`)
	match := re.FindStringSubmatch(m["Response"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Debian-数字"
		answer = append(answer, "debian/"+match[1])
	} else {
		if strings.Contains(m["Server"], "Debian") {
			answer = append(answer, "debian/N")
		}
	}
	fmt.Println(answer)
}
func TestMysql(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	// m["Response"] = "Debian-10"
	m["Service"] = "mysql;ssh"
	m["Version"] = "5.5.68-MariaDB"
	if strings.Contains(m["Service"], "mysql") {
		re := regexp.MustCompile(`^[\d\.]+`)
		match := re.FindString(m["Version"])
		if match != "" {
			answer = append(answer, "mysql/"+match)
		} else {
			answer = append(answer, "mysql/N")
		}
	}
	fmt.Println(answer)
}

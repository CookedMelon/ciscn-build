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
func TestSSH(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	m["Response"] = "SSH-2.0-OpenSSH_8.4p1 Debian-5+deb11u1"
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
	fmt.Println(answer)
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
	m["Server"] = "nginx （Ubuntu）"
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
	m["Server"] = "Jetty(10.0.13)"
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
	// m["Server"] = "(Debian)"
	m["Response"] = "220 ProFTPD Server (Debian)"
	re := regexp.MustCompile(`Debian-(\d+)`)
	match := re.FindStringSubmatch(m["Response"])
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
func TestNodeJs(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	m["FingerPrint"] = "Node.js 12.22.1"
	if strings.Contains(m["FingerPrint"], "Node.js") {
		re := regexp.MustCompile(`Node.js (\d+\.\d+(\.\d+){0,2})`)
		matches := re.FindAllString(m["Body"], -1)
		if len(matches) > 0 {
			answer = append(answer, strings.ToLower(matches[0]))
		} else {
			answer = append(answer, "node.js/N")
		}
	}
	fmt.Println(answer)
}
func TestPHP(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	m["Server"] = "Apache/2.4.29 (Ubuntu) PHP/7.2 lll"
	re := regexp.MustCompile(`(?i)PHP/([\d\.]+)`)
	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "PHP/数字"
		answer = append(answer, "php/"+match[1])
	} else if strings.Contains(m["Server"], "PHP") {
		answer = append(answer, "php/N")
	}
	fmt.Println(answer)
}
func TestMicrosoftHTTPAPI(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	re := regexp.MustCompile(`(?i)Microsoft-HTTPAPI/([\d\.]+)`)
	m["Server"] = "Microsoft-HTTPAPI"
	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Microsoft-HTTPAPI/数字"
		answer = append(answer, "microsoft-httpapi/"+match[1])
	} else if strings.Contains(m["Server"], "Microsoft-HTTPAPI") {
		answer = append(answer, "microsoft-httpapi/N")
	}
	fmt.Println(answer)
}
func TestApache(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	m["Server"] = "Apache"
	// m["Server"] = "Apache/2.4.29 (Ubuntu) PHP/7.2 lll"
	re := regexp.MustCompile(`(?i)Apache/([\d\.]+)`)
	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Apache/数字"
		answer = append(answer, "apache/"+match[1])
	} else if strings.Contains(m["Server"], "Apache") && !strings.Contains(m["Server"], "Apache-Coyote") {
		answer = append(answer, "apache/N")
	}
	fmt.Println(answer)
}
func TestSSL(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	m["Server"] = "Apache/2.4.56 (Win64) OpenSSL/1.1.1t PHP/8.2.4"
	re := regexp.MustCompile(`(?i)OpenSSL/([\d\.]+)`)
	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "openssl/数字"
		answer = append(answer, "openssl/"+match[1])
	} else if strings.Contains(m["Server"], "openssl") {
		answer = append(answer, "openssl/N")
	}
	fmt.Println(answer)
}
func TestUbuntu(t *testing.T) {
	var m = map[string]string{}
	var answer = make([]string, 0)
	// m["Response"] = "Debian-10"
	// m["Server"] = "(Debian)"
	m["Response"] = "SSH-2.0-OpenSSH_6.6.1p1 2020Ubuntu-2ubuntu2"
	re := regexp.MustCompile(`Ubuntu-(\d+)`)
	match := re.FindStringSubmatch(m["Response"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Debian-数字"
		answer = append(answer, "ubuntu/"+match[1])
	} else if strings.Contains(m["Server"], "Ubuntu") {
		answer = append(answer, "ubuntu/N")
	} else {
		match = re.FindStringSubmatch(m["Response"])
		if len(match) > 0 && match[1] != "" { // 如果匹配到 "Debian-数字"
			answer = append(answer, "ubuntu/"+match[1])
		} else if strings.Contains(m["Response"], "Ubuntu") || strings.Contains(m["Response"], "ubuntu") {
			answer = append(answer, "ubuntu/N")
		}
	}
	fmt.Println(answer)
}
func TestWorkPress(t *testing.T) {
	m := map[string]string{}
	m["FingerPrint"] = "WordPress"
	m["Body"] = `string<meta name="generator" content="WordPress 6.2.2">string`
	fmt.Println(m["Body"])
	answer := GetWordPress(m)
	fmt.Println(answer)
}
func TestNgioinx(t *testing.T) {
	m := map[string]string{}
	// m["FingerPrint"] = "WordPress"
	m["Server"] = `nginx/1.22.0`
	fmt.Println(m["Body"])
	answer := GetNginx(m)
	fmt.Println(answer)
}
func TestUbuntu2(t *testing.T) {
	m := map[string]string{}
	// m["FingerPrint"] = "WordPress"
	// m["Response"] = `SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.1 `
	m["Server"] = `(Ubuntu) `
	answer := GetUbuntu(m)
	fmt.Println(answer)
}

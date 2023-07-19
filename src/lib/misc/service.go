package misc

import (
	"fmt"
	"regexp"
	"strings"
)

func GetOpenSSH(m map[string]string) string {

	re := regexp.MustCompile(`SSH-2\.0-OpenSSH_([\d\.]+)`)
	match := re.FindStringSubmatch(m["Response"])
	if match != nil {
		if len(match) < 2 {
			if regexp.MustCompile(`OpenSSH`).FindString(m["Response"]) != "" {
				return "openssh/N"
			}
			return "openssh/N"
		}
		// 如果匹配到 "SSH-2.0-OpenSSH_(数字小数点组合)"，则返回 "openssh/(数字小数点组合)"
		return "openssh/" + match[1]
	}
	if strings.Contains(m["Response"], "OpenSSH") && !strings.Contains(m["Response"], "SSH-2.0-OpenSSH_") {
		return "openssh/N"
	}
	return ""
}

func GetWordPress(m map[string]string) string {
	if strings.Contains(m["FingerPrint"], "WordPress") {
		fmt.Println("find wordpress:", m["FingerPrint"])
		re := regexp.MustCompile(`content="WordPress ([\d\.]+)"`)
		matches := re.FindStringSubmatch(m["Body"])
		if len(matches) >= 2 {
			return "wordpress/" + matches[1]
		} else {
			return "wordpress/N"
		}
	}
	return ""
}

func GetWindows(m map[string]string) string {

	if strings.Contains(m["FingerPrint"], "Windows") {
		fmt.Println("find windows:", m["FingerPrint"])
		re := regexp.MustCompile(`Windows CE (\d+\.\d+)`)
		matches := re.FindAllString(m["FingerPrint"], -1)
		if len(matches) > 0 {
			var temp = strings.ToLower(matches[0])
			return strings.Replace(temp, " ce ", "/", 1)
		} else {
			return "windows/N"
		}
	}
	return ""
}
func GetNginx(m map[string]string) string {

	if strings.Contains(m["Server"], "nginx") {
		//如Server: nginx/1.18.0 (Ubuntu)提取为nginx/1.18.0，如Server: nginx/1.28.0提取为nginx/1.28.0
		// re := regexp.MustCompile(`nginx/(\d+\.\d+(\.\d+){0,2})`)
		re := regexp.MustCompile(`nginx/([\d\.]+)`)
		matches := re.FindStringSubmatch(m["Server"])
		if len(matches) >= 2 {
			return "nginx/" + matches[1]
		} else {
			return "nginx/N"
		}
	} else if strings.Contains(m["FingerPrint"], "nginx") {
		return "nginx/N"
	}
	return ""
}
func GetJetty(m map[string]string) string {

	if strings.Contains(m["Server"], "Jetty") {
		re := regexp.MustCompile(`Jetty(\((.*?)\))?`)
		matches := re.FindStringSubmatch(m["Server"])
		if matches[2] == "" { // 如果括号内无任何文本，则替换为 "Jetty/N"
			return "jetty/N"
		} else {
			return "jetty/" + matches[2]
		}
	}
	return ""
}
func GetDebian(m map[string]string) string {

	re := regexp.MustCompile(`Debian-(\d+)`)
	match := re.FindStringSubmatch(m["Response"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Debian-数字"
		return "debian/" + match[1]
	} else if strings.Contains(m["Server"], "Debian") {
		return "debian/N"
	} else {
		match = re.FindStringSubmatch(m["Response"])
		if len(match) > 0 && match[1] != "" { // 如果匹配到 "Debian-数字"
			return "debian/" + match[1]
		} else if strings.Contains(m["Response"], "Debian") || strings.Contains(m["Response"], "debian") {
			return "debian/N"
		}
	}
	return ""
}
func GetGrafana(m map[string]string) string {

	if strings.Contains(m["FingerPrint"], "Grafana") {
		re := regexp.MustCompile(`Grafana (\d+\.\d+(\.\d+){0,2})`)
		matches := re.FindAllString(m["Body"], -1)
		if len(matches) > 0 {
			return strings.ToLower(matches[0])
		} else {
			return "grafana/N"
		}
	}
	return ""
}
func GetNodeJS(m map[string]string) string {

	if strings.Contains(m["FingerPrint"], "Node.js") {
		re := regexp.MustCompile(`Node.js (\d+\.\d+(\.\d+){0,2})`)
		matches := re.FindAllString(m["Body"], -1)
		if len(matches) > 0 {
			return strings.ToLower(matches[0])
		} else {
			return "node.js/N"
		}
	}
	return ""
}
func GetPHP(m map[string]string) string {

	re := regexp.MustCompile(`(?i)PHP/([\d\.]+)`)
	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "PHP/数字"
		return "php/" + match[1]
	} else if strings.Contains(m["Server"], "PHP") {
		return "php/N"
	} else if strings.Contains(m["X-Powered-By"], "PHP") {
		re = regexp.MustCompile(`(?i)PHP/([\d\.]+)`)
		match = re.FindStringSubmatch(m["X-Powered-By"])
		if len(match) > 0 && match[1] != "" { // 如果匹配到 "PHP/数字"
			return "php/" + match[1]
		} else if strings.Contains(m["X-Powered-By"], "PHP") {
			return "php/N"
		}
	}
	return ""
}
func GetHttpAPI(m map[string]string) string {

	re := regexp.MustCompile(`(?i)Microsoft-HTTPAPI/([\d\.]+)`)
	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Microsoft-HTTPAPI/数字"
		return "microsoft-httpapi/" + match[1]
	} else if strings.Contains(m["Server"], "Microsoft-HTTPAPI") {
		return "microsoft-httpapi/N"
	}
	return ""
}
func GetApache(m map[string]string) string {
	re := regexp.MustCompile(`(?i)Apache/([\d\.]+)`)

	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Apache/数字"
		return "apache/" + match[1]
	} else if strings.Contains(m["Server"], "Apache") && !strings.Contains(m["Server"], "Apache-Coyote") {
		return "apache/N"
	}
	return ""
}
func GetIIS(m map[string]string) string {

	re := regexp.MustCompile(`(?i)IIS/([\d\.]+)`)
	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "iis/数字"
		return "iis/" + match[1]
	} else if strings.Contains(m["Server"], "iis") {
		return "iis/N"
	}
	return ""
}
func GetOpenSSL(m map[string]string) string {

	re := regexp.MustCompile(`(?i)OpenSSL/([\d\.]+)`)
	match := re.FindStringSubmatch(m["Server"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "openssl/数字"
		return "openssl/" + match[1]
	} else if strings.Contains(m["Server"], "openssl") {
		return "openssl/N"
	}
	return ""
}
func GetUbuntu(m map[string]string) string {

	re := regexp.MustCompile(`Ubuntu-(\d+ubuntu[\d\.]+)`)
	match := re.FindStringSubmatch(m["Response"])
	if len(match) > 0 && match[1] != "" { // 如果匹配到 "Ubuntu-数字"
		return "ubuntu/" + match[1]
	} else if strings.Contains(m["Server"], "Ubuntu") {
		return "ubuntu/N"
	} else {
		match = re.FindStringSubmatch(m["Response"])
		if len(match) > 0 && match[1] != "" { // 如果匹配到 "ubuntu-数字"
			return "ubuntu/" + match[1]
		} else if strings.Contains(m["Response"], "Ubuntu") || strings.Contains(m["Response"], "ubuntu") {
			return "ubuntu/N"
		}
	}
	return ""
}

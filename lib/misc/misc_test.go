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

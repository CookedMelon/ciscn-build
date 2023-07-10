package run

import (
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {
	// InitJKjson()
}

func TestAdd1(t *testing.T) {
	jsonData := GetJson()
	fmt.Println(jsonData)
	Addip("192.168.1.0", jsonData)
	Addip("192.168.1.1", jsonData)
}

func TestAdd2(t *testing.T) {
	// Addip("192168.1.0")
}

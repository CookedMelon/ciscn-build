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
	Addip("192.168.1.3", jsonData)
	Addip("192.168.1.1", jsonData)
}

func TestAdd2(t *testing.T) {
	jsonService1 := Service{
		Protocol: "http",
		Port:     80,
		ServiceApp: []string{
			"Apache",
			"nginx",
		},
	}
	jsonService2 := Service{
		Protocol: "https",
		Port:     443,
		ServiceApp: []string{
			"Apache",
			"nginx",
		},
	}
	jsonService3 := Service{
		Protocol:   "https",
		Port:       4,
		ServiceApp: nil,
	}
	jsonService4 := Service{
		Protocol: "https",
		Port:     55,
		ServiceApp: []string{
			"abc",
			"def",
		},
	}
	jsonData1 := Data{
		Services:   []Service{jsonService1, jsonService2, jsonService3},
		DeviceInfo: []string{"Windows 10", "Windows 7"},
		Honeypot:   []string{"honeypot1", "honeypot2"},
	}
	jsonData2 := Data{
		Services:   []Service{jsonService1, jsonService2},
		DeviceInfo: []string{"Windows 10", "Windows 7"},
		Honeypot:   []string{"honeypot1"},
	}
	jsonData3 := Data{
		Services:   []Service{jsonService1, jsonService4},
		DeviceInfo: []string{"Windows 10", "Windows 7"},
		Honeypot:   []string{"honeypot1"},
	}
	Add("192.168.0.3", jsonData1)
	Add("192.168.0.3", jsonData2)
	Add("192.168.0.3", jsonData3)

}

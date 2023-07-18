package run

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	ap "jkscan/app"
)

type Data struct {
	Services   []Service `json:"services"`
	DeviceInfo []string  `json:"deviceinfo"`
	Honeypot   []string  `json:"honeypot"`
}
type Service struct {
	Port       int      `json:"port"`
	Protocol   string   `json:"protocol"`
	ServiceApp []string `json:"service_app"`
}

var globalData = make(map[string]Data)
var globalCount = 0
var firstRead = true

func GetJson() map[string]Data {
	// 读取原始的JSON数据
	if firstRead {
		file, err := ioutil.ReadFile(ap.Args.OutputJson)
		if err != nil {
			// 文件不存在，创建新的数据结构
			jsonData := make(map[string]Data)

			// 将数据结构转换为JSON格式
			updatedJSON, err := json.MarshalIndent(jsonData, "", "    ")
			if err != nil {
				fmt.Println("无法转换为JSON：", err)
				return jsonData
			}

			// 创建新的JSON文件并写入数据
			err = ioutil.WriteFile(ap.Args.OutputJson, updatedJSON, 0644)
			if err != nil {
				fmt.Println("无法写入文件：", err)
				return jsonData
			}

			fmt.Println("已创建新的JSON文件并初始化数据。")
			return jsonData
		}

		// 解析原始JSON数据
		var jsonData map[string]Data
		err = json.Unmarshal(file, &jsonData)
		if err != nil {
			fmt.Println("无法解析JSON：", err)
			return jsonData
		}

		fmt.Println("已加载现有的JSON数据。")
		firstRead = false
		globalData = jsonData
		return jsonData
	} else {
		return globalData
	}
}
func Addip(ip string, jsonData map[string]Data) {
	// 检查IP地址是否存在于JSON数据中
	fmt.Println(jsonData)
	_, exists := jsonData[ip]
	if !exists {
		// IP地址不存在，创建新的条目
		newData := Data{
			Services:   []Service{},
			DeviceInfo: []string{"new"},
			Honeypot:   []string{},
		}
		fmt.Println(newData)
		jsonData[ip] = newData
		fmt.Println(jsonData)

		// 将更新后的JSON数据写入文本文件
		updatedJSON, err := json.MarshalIndent(jsonData, "", "    ")
		if err != nil {
			fmt.Println("无法转换为JSON：", err)
			return
		}

		err = ioutil.WriteFile(ap.Args.OutputJson, updatedJSON, 0644)
		if err != nil {
			fmt.Println("无法写入文件：", err)
			return
		}

		fmt.Println("已向JSON数据添加新的条目。")
	} else {
		// IP地址已存在
		fmt.Println("JSON数据中已存在该IP地址。")
	}
}

/*
@author: xiongsp

@description: 添加数据到data.json

@usage: Add(ip, addData)

- ip: IP地址

- addData: 要添加的数据，类型为Data，差量更新
*/
func Add(ip string, addData Data) {
	// 检查IP地址是否存在于JSON数据中
	jsonData := GetJson()
	_, exists := jsonData[ip]
	if !exists {
		// IP地址不存在，创建新的条目
		jsonData[ip] = addData
		fmt.Println("已向JSON数据添加新的条目。")
	} else {
		// IP地址已存在
		// 进行差量更新
		oldData := jsonData[ip]
		resuldData := Data{}
		// 设备信息
		if addData.DeviceInfo != nil {
			resuldData.DeviceInfo = addData.DeviceInfo
		} else {
			resuldData.DeviceInfo = oldData.DeviceInfo
		}
		// 蜜罐信息，二者合并
		tmpHoneypot := make(map[string]bool)
		for _, v := range oldData.Honeypot {
			tmpHoneypot[v] = true
		}
		for _, v := range addData.Honeypot {
			tmpHoneypot[v] = true
		}
		for k := range tmpHoneypot {
			resuldData.Honeypot = append(resuldData.Honeypot, k)
		}
		// 服务信息，二者合并
		tmpService := make(map[string]Service)
		for _, v := range oldData.Services {
			tmpService[fmt.Sprintf("%s%d", v.Protocol, v.Port)] = v
		}
		for _, v := range addData.Services {
			tmpService[fmt.Sprintf("%s%d", v.Protocol, v.Port)] = v
		}
		for _, v := range tmpService {
			resuldData.Services = append(resuldData.Services, v)
		}
		jsonData[ip] = resuldData
		fmt.Println("已更新JSON数据。")
	}

	// 将更新后的JSON数据写入文本文件
	updatedJSON, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		fmt.Println("无法转换为JSON：", err)
		return
	}
	globalData = jsonData
	globalCount++
	fmt.Println(globalCount)
	// 每20次写入一次文件
	if globalCount >= 20 {
		globalCount = 0
		err = ioutil.WriteFile(ap.Args.OutputJson, updatedJSON, 0644)
		if err != nil {
			fmt.Println("无法写入文件：", err)
			return
		}
	}
}

func FlushBuffer() {
	jsonData := GetJson()
	// 将更新后的JSON数据写入文本文件
	updatedJSON, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		fmt.Println("无法转换为JSON：", err)
		return
	}
	err = ioutil.WriteFile(ap.Args.OutputJson, updatedJSON, 0644)
	if err != nil {
		fmt.Println("无法写入文件：", err)
		return
	}
}

package run

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	ap "jkscan/app"
	"os"
	"time"
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

// 记录最后一次更新时间
var lastUpdateTime = time.Now()

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
	lastUpdateTime = time.Now()
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

	// // 将更新后的JSON数据写入文本文件
	// updatedJSON, err := json.MarshalIndent(jsonData, "", "    ")
	// if err != nil {
	// 	fmt.Println("无法转换为JSON：", err)
	// 	return
	// }
	globalData = jsonData
	globalCount++
	fmt.Println(globalCount)
	// 每20次写入一次文件
	// if globalCount >= 20 {
	// 	globalCount = 0
	// 	err = ioutil.WriteFile(ap.Args.OutputJson, updatedJSON, 0644)
	// 	if err != nil {
	// 		fmt.Println("无法写入文件：", err)
	// 		return
	// 	}
	// }
}

// 读取globalData，写入文件
func WriteFile() {
	// 将更JSON数据写入文本文件
	updatedJSON, err := json.MarshalIndent(globalData, "", "    ")
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

// 定义定时函数，每60秒执行一次，调用WriteFile
func timer() {
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {
		//进行一次文件更新
		fmt.Println("定时器触发，进行一次文件更新。")
		WriteFile()
		//如果距离最近一次更新时间超过10分钟，则直接停止程序
		if time.Since(lastUpdateTime) > 20*time.Minute {
			fmt.Println("已超过20分钟未更新，程序即将退出。")
			os.Exit(-1)
			return
		}
	}
}

// 计时器函数
func Timer() {
	// 启动定时器
	fmt.Println("启动定时器。")
	go timer()
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

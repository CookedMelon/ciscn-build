package run

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Data struct {
	Services   []Service `json:"services"`
	DeviceInfo string    `json:"deviceinfo"`
	Honeypot   []string  `json:"honeypot"`
}
type Service struct {
	Port       int      `json:"port"`
	Protocol   string   `json:"protocol"`
	ServiceApp []string `json:"service_app"`
}

func GetJson() map[string]Data {
	// 读取原始的JSON数据
	file, err := ioutil.ReadFile("data.json")
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
		err = ioutil.WriteFile("data.json", updatedJSON, 0644)
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
	return jsonData
}
func Addip(ip string, jsonData map[string]Data) {
	// 检查IP地址是否存在于JSON数据中
	fmt.Println(jsonData)
	_, exists := jsonData[ip]
	if !exists {
		// IP地址不存在，创建新的条目
		newData := Data{
			Services:   []Service{},
			DeviceInfo: "new device",
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

		err = ioutil.WriteFile("data.json", updatedJSON, 0644)
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
func AddService(port int, ip, protocol, app string) {
	// 读取原始的JSON数据
	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println("无法读取文件：", err)
		return
	}

	// 解析原始JSON数据
	jsonData := make(map[string]Data)
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		fmt.Println("无法解析JSON：", err)
		return
	}

	// 检查IP地址是否存在于JSON数据中
	data, exists := jsonData[ip]
	if !exists {
		// IP地址不存在，创建新的条目
		data = Data{
			Services: []Service{},
		}
		jsonData[ip] = data
	}

	// 检查端口是否存在于服务中
	// service, exists := data.Services[port]
	// if exists {
	// 	// 端口已存在
	// 	if service.Port != port {
	// 		// 如果app不同，将app添加到ServiceApp字段中
	// 		service.ServiceApp = append(service.ServiceApp, app)
	// 		data.Services[port] = service
	// 	} else {
	// 		// 如果app相同，则不进行任何操作
	// 		fmt.Println("该IP地址和端口的服务已存在。")
	// 		return
	// 	}
	// } else {
	// 端口不存在，创建新的服务
	// newService := Service{
	// 	Port:       port,
	// 	Protocol:   protocol,
	// 	ServiceApp: []string{app},
	// }
	// data.Services[port] = newService
	// }

	// 将更新后的JSON数据写入文本文件
	updatedJSON, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		fmt.Println("无法转换为JSON：", err)
		return
	}

	err = ioutil.WriteFile("data.json", updatedJSON, 0644)
	if err != nil {
		fmt.Println("无法写入文件：", err)
		return
	}

	fmt.Println("已向JSON数据添加新的服务信息。")
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
		if addData.DeviceInfo != "" {
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

	err = ioutil.WriteFile("data.json", updatedJSON, 0644)
	if err != nil {
		fmt.Println("无法写入文件：", err)
		return
	}
}

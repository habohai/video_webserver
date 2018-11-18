package config

import (
	"encoding/json"
	"os"
)

// Configuration 配置结构体
type Configuration struct {
	LBAddr  string `json:"lb_addr"`
	OssAddr string `json:"oss_addr"`
}

var configuration *Configuration

func init() {
	file, _ := os.Open("./conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration = &Configuration{}

	if err := decoder.Decode(configuration); err != nil {
		panic(err)
	}
}

// GetLbAddr 获取Lb的地址
func GetLbAddr() string {
	return configuration.LBAddr
}

// GetOssAddr 获取Oss的地址
func GetOssAddr() string {
	return configuration.OssAddr
}

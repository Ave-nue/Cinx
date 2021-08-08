package utils

import (
	"cinx/ciface"
	"encoding/json"
	"io/ioutil"
	"os"
)

//用于存储框架全局变量
type GlobalConfig struct {
	//server相关
	//服务名称
	Name string
	//当前server对象
	TcpServer ciface.IServer
	//当前主机监听IP
	Host string
	//端口
	Port int

	//Cinx相关
	//版本号
	Version string
	//最大连接数
	MaxConn int
	//最大数据包包体大小
	MaxPackageSize uint32
	//业务worker池大小
	WorkerPoolSize uint32
	//消息队列缓冲区长度
	RequestQueneLength uint32
}

var GlobalCfg *GlobalConfig

//加载配置文件
func (cfg *GlobalConfig) LoadConfig() {
	//没有文件就不需要load了
	if _, err := os.Stat("config/Config_cinx.json"); err != nil {
		return
	}

	data, err := ioutil.ReadFile("config/Config_cinx.json")
	if err != nil {
		panic(err) //停止当前gorotine
	}
	//解析json
	err = json.Unmarshal(data, &GlobalCfg)
	if err != nil {
		panic(err) //停止当前gorotine
	}
}

//当调用此包时会自动调用其中的init方法，用于初始化变量等
func init() {
	//设置默认值
	GlobalCfg = &GlobalConfig{
		Name:               "CinxServer",
		Version:            "2021.8.7",
		Host:               "0.0.0.0",
		Port:               6608,
		MaxConn:            16384,
		MaxPackageSize:     4096,
		WorkerPoolSize:     8,
		RequestQueneLength: 1024,
	}

	//尝试从config/Config_cinx.json加载配置数据
	GlobalCfg.LoadConfig()
}

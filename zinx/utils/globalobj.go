package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

// 定义一些全局的参数，方便修改
type GlobalObj struct {
	TcpServer      ziface.IServer //当前zinx的全局Server对象
	Host           string         //当前服务器主机IP
	TcpPort        int            //当前服务器主机监听端口
	Name           string         //当前服务器名称
	Version        string         //当前zinx版本号
	MaxPacketSize  uint32         //读取数据包的最大值
	MaxConn        int            //当前服务器允许连接的最大连接数
	WorkerPoolSize uint32         //工作池worker数量
	MaxWorkTaskLen uint32         //任务消息队列最大request存储数量
	MaxMsgChanLen  uint32         //缓冲通道容量

	ConfilePath string //?

}

var GlobalObject *GlobalObj //
// 加载本地配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("zinx.json") //读取json文件内容
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject) //将数据解析到结构体GlobalObject中
	if err != nil {
		panic(err)
	}

}

// 提供初始化方法，调用Reload()加载json文件参数设置
func init() {
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        7777,
		Host:           "0.0.0.0",
		MaxConn:        12000,
		MaxPacketSize:  4096,
		WorkerPoolSize: 10,
		MaxWorkTaskLen: 1024,
		MaxMsgChanLen:  1024,
		ConfilePath:    "maintest/zinx.json",
	}
	GlobalObject.Reload()
}

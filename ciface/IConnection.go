package ciface

import "net"

type IConnection interface {
	//启动连接
	Start()
	//断开连接
	Stop()
	//获取连接句柄
	GetConnection() *net.TCPConn
	//获取连接ID
	GetID() uint32
	//获取IP:Port
	RemoteAddr() net.Addr
	//发送消息
	SendMessage(ID uint32, data []byte) error
	//连接属性操作
	SetProperty(string, interface{})
	GetProperty(string) (interface{}, error)
	RemoveProperty(string)
}

//定义连接处理业务的方法
type ConnectionFunc func(*net.TCPConn, []byte, int) error

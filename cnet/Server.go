package cnet

import (
	"cinx/ciface"
	"cinx/utils"
	"fmt"
	"net"
)

type Server struct {
	Name   string
	Socket string
	IP     string
	Port   int

	//当前的api
	MsgHandler ciface.IMessageHandler
	//连接管理器
	ConnMgr ciface.IConnectionManager
	//创建连接之后的hook
	OnConnectionStart func(conn ciface.IConnection)
	//断开连接之前的hook
	OnConnectionStop func(conn ciface.IConnection)
}

func (server *Server) Start() {
	//初始化消息队列及工作池
	server.MsgHandler.InitWorkerPool()

	//1 获取TCP的addr
	addr, err := net.ResolveTCPAddr(server.Socket, fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Printf("[Cnet]Cannot Resolve TCP addr netType = %s addr = %s:%d\n%s\n", server.Socket, server.IP, server.Port, err)
		return
	}

	//2 监听
	listener, err := net.ListenTCP(server.Socket, addr)
	if err != nil {
		fmt.Printf("[Cnet]Server start failed because can not listen with network = %s addr = %s\n%s\n", server.Socket, addr, err)
		return
	}
	fmt.Printf("[Cnet]Server \"%s\" start successed!", server.Name)

	//3 阻塞并处理客户端请求
	var connID uint32 = 0
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("[Cnet]Server connect build faild!\n%s\n", err)
			continue
		}

		//判断是否已经超出最大连接数
		if server.ConnMgr.Len() >= utils.GlobalCfg.MaxConn {
			fmt.Printf("[Cnet]Server connect is full!\n%s\n", err)
			//todo 给客户端响应一个错误码（超出最大连接数）
			conn.Close()
			continue
		}

		newConn := NewConnection(conn, connID, server)
		connID++

		go newConn.Start()
	}
}

func (server *Server) Stop() {
	//释放服务器资源，以及监听中的连接的停止和回收
	server.ConnMgr.Clear()
	fmt.Printf("[Cnet]Server \"%s\" stop successed!", server.Name)
}

func (server *Server) Serve() {
	//启动监听
	go server.Start()

	//todo 执行启动服务器后的其他逻辑

	//阻塞
	select {}
}

func (server *Server) AddRouter(msgID uint32, router ciface.IRouter) {
	server.MsgHandler.AddRouter(msgID, router)
	fmt.Println("[Cnet]Add router for server successed!")
}

func (server *Server) GetConnMgr() ciface.IConnectionManager {
	return server.ConnMgr
}

func (server *Server) GetMsgHandler() ciface.IMessageHandler {
	return server.MsgHandler
}

func NewServer() *Server {
	return &Server{
		Name:       utils.GlobalCfg.Name,
		Socket:     "tcp",
		IP:         utils.GlobalCfg.Host,
		Port:       utils.GlobalCfg.Port,
		MsgHandler: NewMessageHandler(),
		ConnMgr:    NewConnectionManager(),
	}
}

func (server *Server) SetOnConnectionStart(hookFunc func(ciface.IConnection)) {
	server.OnConnectionStart = hookFunc
}
func (server *Server) SetOnConnectionStop(hookFunc func(ciface.IConnection)) {
	server.OnConnectionStop = hookFunc
}
func (server *Server) CallOnConnectionStart(conn ciface.IConnection) {
	if server.OnConnectionStart == nil {
		fmt.Println("[Cnet]ConnectionStart function is empty!")
		return
	}

	server.OnConnectionStart(conn)
}
func (server *Server) CallOnConnectionStop(conn ciface.IConnection) {
	if server.OnConnectionStop == nil {
		fmt.Println("[Cnet]ConnectionStop function is empty!")
		return
	}

	server.OnConnectionStop(conn)
}

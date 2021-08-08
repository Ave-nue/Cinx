package cnet

import (
	"cinx/ciface"
	"cinx/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前连接
	Conn *net.TCPConn
	//连接ID
	ID uint32
	//当前请求累计数量
	requestCount uint32
	//连接状态
	isClosed bool
	//是否已退出连接 channel
	ExitChan chan bool
	//用于读写通道之间的通信
	msgChan chan []byte
	//server
	server ciface.IServer
	//连接属性集合
	property map[string]interface{} //此为万能指针写法
	//集合属性修改锁
	propertyLock sync.RWMutex
}

func NewConnection(conn *net.TCPConn, ID uint32, server ciface.IServer) *Connection {
	c := &Connection{
		Conn:         conn,
		ID:           ID,
		requestCount: 0,
		isClosed:     false,
		ExitChan:     make(chan bool),
		msgChan:      make(chan []byte),
		server:       server,
		property:     make(map[string]interface{}),
	}

	c.server.GetConnMgr().Add(c)

	return c
}

//读业务
func (conn *Connection) StartReader() {
	fmt.Println("[Cnet]Start Reader for Connection with ID:", conn.ID)
	defer fmt.Println("[Cnet]Read is finished for Connection with ID:", conn.ID)
	defer conn.Stop()

	//创建拆包工具
	dp := NewDataPack()
	for {
		//先读head
		headData := make([]byte, dp.GetHeadLength())
		_, err := io.ReadFull(conn.Conn, headData)
		if err != nil {
			fmt.Println("[Cnet]Read message head error\n", err)
			return
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("[Cnet]Head unpack error\n", err)
			return
		}
		if msg.GetLength() > 0 {
			//当数据长度>0时根据读出的head内容再读data
			msg.SetData(make([]byte, msg.GetLength()))
			_, err = io.ReadFull(conn.Conn, msg.GetData())
			if err != nil {
				fmt.Println("[Cnet]Data unpack error\n", err)
				return
			}
		}

		//绑定connection和data，封装为request
		req := Request{
			ID:      conn.requestCount,
			conn:    conn,
			message: msg,
		}
		conn.requestCount++

		//执行路由对应的数据处理逻辑
		if utils.GlobalCfg.WorkerPoolSize > 0 { //开启了工作池机制
			conn.server.GetMsgHandler().AddTask(&req)
		} else {
			conn.server.GetMsgHandler().DoMessage(&req)
		}
	}
}

//写业务
func (conn *Connection) StartWriter() {
	fmt.Println("[Cnet]Start Writer for Connection with ID:", conn.ID)
	defer fmt.Println("[Cnet]Write is finished for Connection with ID:", conn.ID)
	defer conn.Stop()

	for { //不断检查是否有消息需要写
		select {
		case data := <-conn.msgChan: //只要channal里有数据待写，就立马写出去
			if _, err := conn.Conn.Write(data); err != nil {
				fmt.Println("[Cnet]Write message head error\n", err)
				return
			}
		case <-conn.ExitChan:
			return
		}
	}
}

//启动连接
func (conn *Connection) Start() {
	fmt.Println("[Cnet]Start Connection with ID:", conn.ID)
	//启动读数据的线程
	go conn.StartReader()
	//启动写数据的线程
	go conn.StartWriter()

	//调用hook
	conn.server.CallOnConnectionStart(conn)
}

//断开连接
func (conn *Connection) Stop() {
	fmt.Println("[Cnet]Stop Connection with ID:", conn.ID)

	if conn.isClosed {
		return
	}

	conn.isClosed = true

	//调用hook
	conn.server.CallOnConnectionStop(conn)

	//关闭连接
	conn.Conn.Close()

	//告知writter关闭（似乎不需要
	//conn.ExitChan <- true

	//从连接管理中去除
	conn.server.GetConnMgr().Remove(conn)

	//回收channel
	close(conn.ExitChan)
	close(conn.msgChan)
}

//获取连接句柄
func (conn *Connection) GetConnection() *net.TCPConn {
	return conn.Conn
}

//获取连接ID
func (conn *Connection) GetID() uint32 {
	return conn.ID
}

//获取IP:Port
func (conn *Connection) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

//发送消息
func (conn *Connection) SendMessage(ID uint32, data []byte) error {
	if conn.isClosed {
		return errors.New("[Cnet]Connection has closed")
	}

	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(ID, data))
	if err != nil {
		return err
	}

	conn.msgChan <- binaryMsg

	return nil
}

func (conn *Connection) SetProperty(key string, value interface{}) {
	conn.propertyLock.Lock()
	defer conn.propertyLock.Unlock()

	conn.property[key] = value
}
func (conn *Connection) GetProperty(key string) (interface{}, error) {
	conn.propertyLock.RLock()
	defer conn.propertyLock.RUnlock()

	if value, ok := conn.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("[Cnet]Connection property not found")
	}
}
func (conn *Connection) RemoveProperty(key string) {
	conn.propertyLock.Lock()
	defer conn.propertyLock.Unlock()

	delete(conn.property, key)
}

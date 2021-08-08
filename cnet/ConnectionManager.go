package cnet

import (
	"cinx/ciface"
	"errors"
	"fmt"
	"sync"
)

type ConnectionManager struct {
	//管理的连接集合
	connections map[uint32]ciface.IConnection
	//读写互斥锁
	connLock sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]ciface.IConnection),
	}
}

//添加连接
func (connMgr *ConnectionManager) Add(conn ciface.IConnection) {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//加入
	connMgr.connections[conn.GetID()] = conn
	fmt.Printf("[Cnet]Add Connection with id = %d, now connection num = %d\n", conn.GetID(), connMgr.Len())
}

//删除
func (connMgr *ConnectionManager) Remove(conn ciface.IConnection) {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除
	delete(connMgr.connections, conn.GetID())
	fmt.Printf("[Cnet]Delete Connection with id = %d, now connection num = %d\n", conn.GetID(), connMgr.Len())
}

//根据ID获取连接
func (connMgr *ConnectionManager) Get(connID uint32) (ciface.IConnection, error) {
	//加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取连接总数
func (connMgr *ConnectionManager) Len() int {
	return len(connMgr.connections)
}

//清除并终止所有连接
func (connMgr *ConnectionManager) Clear() {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除并停止所有coonection
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}

	fmt.Printf("[Cnet]Clear all Connections , now connection num = %d\n", connMgr.Len())
}

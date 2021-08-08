package cnet

import "cinx/ciface"

type Request struct {
	//ID
	ID uint32
	//连接
	conn ciface.IConnection
	//消息
	message ciface.IMessage
}

func (request *Request) GetID() uint32 {
	return request.ID
}

func (request *Request) GetConnection() ciface.IConnection {
	return request.conn
}

func (request *Request) GetMessageID() uint32 {
	return request.message.GetID()
}

func (request *Request) GetData() []byte {
	return request.message.GetData()
}

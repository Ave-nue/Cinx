package ciface

type IRequest interface {
	//获取ID
	GetID() uint32
	//获取当前连接
	GetConnection() IConnection
	//消息ID
	GetMessageID() uint32
	//请求数据
	GetData() []byte
}

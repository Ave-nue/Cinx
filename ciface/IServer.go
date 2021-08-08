package ciface

type IServer interface {
	Start()
	Stop()
	Serve()

	//给当前的服务新增路由
	AddRouter(msgID uint32, ruter IRouter)
	//获取连接管理类
	GetConnMgr() IConnectionManager
	//获取消息管理类
	GetMsgHandler() IMessageHandler

	//hook函数相关
	SetOnConnectionStart(func(IConnection))
	SetOnConnectionStop(func(IConnection))
	CallOnConnectionStart(IConnection)
	CallOnConnectionStop(IConnection)
}

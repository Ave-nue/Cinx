package ciface

type IMessageHandler interface {
	//执行对应的router逻辑
	DoMessage(IRequest)
	//添加router
	AddRouter(uint32, IRouter)
	//初始化工作池
	InitWorkerPool()
	//压入新任务进入消息队列
	AddTask(IRequest)
}

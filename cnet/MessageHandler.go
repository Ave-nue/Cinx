package cnet

import (
	"cinx/ciface"
	"cinx/utils"
	"fmt"
)

type MessageHandler struct {
	//所有id到处理逻辑的对应关系表
	Apis map[uint32]ciface.IRouter
	//消息队列
	TaskQueue []chan ciface.IRequest
	//业务worker池size
	WorkerPoolSize uint32
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		Apis:           make(map[uint32]ciface.IRouter),
		TaskQueue:      make([]chan ciface.IRequest, utils.GlobalCfg.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalCfg.WorkerPoolSize,
	}
}

//执行对应的router逻辑
func (msgHandler *MessageHandler) DoMessage(request ciface.IRequest) {
	router, ok := msgHandler.Apis[request.GetMessageID()]
	if !ok { //是否已注册
		fmt.Println("[Cnet]No router for message id ", request.GetMessageID())
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

//添加router
func (msgHandler *MessageHandler) AddRouter(msgID uint32, router ciface.IRouter) {
	if _, ok := msgHandler.Apis[msgID]; ok { //是否已注册
		fmt.Println("[Cnet]Repeat message router id ", msgID)
		return
	}
	msgHandler.Apis[msgID] = router
	fmt.Println("[Cnet]Add router for handler successed!")
}

//初始化工作池
func (msgHandler *MessageHandler) InitWorkerPool() {
	for i := 0; i < int(msgHandler.WorkerPoolSize); i++ {
		msgHandler.TaskQueue[i] = make(chan ciface.IRequest, utils.GlobalCfg.RequestQueneLength)
		go msgHandler.StartWorker(i)
	}
}

//开始一项新工作
func (msgHandler *MessageHandler) StartWorker(workerID int) {
	fmt.Printf("[Cnet]New worker start with id = %d\n", workerID)
	taskQueue := msgHandler.TaskQueue[workerID]

	//阻塞等待队列中的消息
	for {
		select {
		case request := <-taskQueue:
			msgHandler.DoMessage(request)
		}
	}
}

//压入一个新的任务进入消息队列
func (msgHandler *MessageHandler) AddTask(request ciface.IRequest) {
	workerID := request.GetID() % msgHandler.WorkerPoolSize
	msgHandler.TaskQueue[workerID] <- request
	fmt.Printf("[Cnet]Add Task id = %d to worker id = %d\n", request.GetID(), workerID)
}

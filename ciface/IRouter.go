package ciface

type IRouter interface {
	//业务前方法
	PreHandle(request IRequest)
	//业务主方法
	Handle(request IRequest)
	//业务后方法
	PostHandle(request IRequest)
}

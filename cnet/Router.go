package cnet

import "cinx/ciface"

//用于自定义router的继承
type BaseRouter struct {
}

func (router *BaseRouter) PreHandle(request ciface.IRequest) {}

func (router *BaseRouter) Handle(request ciface.IRequest) {}

func (router *BaseRouter) PostHandle(request ciface.IRequest) {}

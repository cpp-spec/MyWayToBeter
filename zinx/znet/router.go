package znet

import "zinx/ziface"

// 实现router先嵌入此基类，然后根据需要对此基类的方法重写
type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(req ziface.IRequest)  {}
func (br *BaseRouter) Handle(req ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}

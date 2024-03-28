package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

/*
	func (this *PingRouter) PreHandle(request ziface.IRequest) {
		fmt.Println("call router PreHandle...")
		_, err := request.GetConnection().GetTCPConnection().Write([]byte("befor Ping ...\n"))
		if err != nil {
			fmt.Println("call back ping ping ping error")
		}

}
*/
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call router Handle...")
	fmt.Println("recv from client: msgId=", request.GetMsgId(), ",data=", string(request.GetDATa()))
	err := request.GetConnection().SendMsg(0, []byte("Pingpingping..."))
	if err != nil {
		fmt.Println(err)
	}
}

/*
	func (this *PingRouter) PostHandle(request ziface.IRequest) {
		fmt.Println("call router PostHandle...")
		_, err := request.GetConnection().GetTCPConnection().Write([]byte("after Ping ...\n"))
		if err != nil {
			fmt.Println("call back ping ping ping error")
		}
	}
*/
type helloZinxRouter struct {
	znet.BaseRouter
}

func (this *helloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("call helloZinxRouter Handle...")
	fmt.Println("recv from client: msgId=", request.GetMsgId(), ",data=", string(request.GetDATa()))
	err := request.GetConnection().SendBuffMsg(1, []byte("hello zinx router V0.8"))
	if err != nil {
		fmt.Println(err)
	}
}
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnection is called...")
	err := conn.SendMsg(2, []byte("DoConnection begin..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is called...")
}
func main() {
	s := znet.NewServer()

	//注册回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStart(DoConnectionLost)

	//s.AddRouter(&PingRouter{}) //配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &helloZinxRouter{})
	s.Serve()
}

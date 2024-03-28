package znet

import "zinx/ziface"

// 定义的Request接口结构体也很符合Irequest抽象层的两个接口
type Request struct {
	conn ziface.IConnection
	//data []byte
	msg ziface.IMessage //客户端的请求数据
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetDATa() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

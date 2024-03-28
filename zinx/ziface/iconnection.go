package ziface

import "net"

// 本质是一个接口
type IConnection interface {
	Start()                                      //启动链接
	Stop()                                       //关闭连接
	GetTCPConnection() *net.TCPConn              //获取原始的socket TCPConn
	GetConnID() uint32                           //获取当前连接id
	RemoteAddr() net.Addr                        //获取远程客户端地址信息
	SendMsg(msgId uint32, data []byte) error     //想要返回给客户端TLV格式的消息
	SendBuffMsg(msgId uint32, data []byte) error //想要返回给客户端TLV格式的消息 缓冲通道
}

// 统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error

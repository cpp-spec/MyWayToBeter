package ziface

type IMessage interface {
	GetDateLen() uint32
	GetMsgId() uint32
	GetData() []byte //[]byte切片类型，处理二进制数据。用于读写文件网络通信等

	SetDateLen(uint32)
	SetMsgId(uint32)
	SetData([]byte) //关于接口是否需要有返回值的问题，总结一下，set没有返回值，get有；而且这两种方法是互相对应的

}

package znet

// 此结构体组成部分是与IMessage接口相呼应的
type Message struct {
	Id      uint32
	Data    []byte
	DataLen uint32
}

// 创建初始化一个Message结构体
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

// 经验总结，设置接口，其实接口是抽象的，一个抽象的接口包含多个抽象的接口方法；下一步是具体实现抽象的接口对象，即创建一个对应结构体，这个结构体是方法们操作的对象，他应该
// 包含什么？我认为在代码的目的是实现功能，所以功能和业务需要什么他就应该有什么。比如此处，我们需要获取客户端发送数据的id编号，数据内容Data，数据长度DataLen。
// 然后提供创建初始化对象的方法New(),接着就是对抽象方法进行实现了。1.绑定方法作用的对象。2.关注输入参数和返回值。3.数据类型值得注意。
func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) GetDateLen() uint32 {
	return msg.DataLen
}

func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

func (msg *Message) SetDateLen(len uint32) {
	msg.DataLen = len
}

func (msg *Message) SetMsgId(id uint32) {
	msg.Id = id
}

// 定义interface接口，func函数实现。大大增强了灵活性，有点类似虚函数。
package ziface

// 抽象层的两个方法，客户端连接和客户端传来的数据
// 返回IConnection是因为要满足设计模式中抽象层依赖抽象层
type IRequest interface {
	GetConnection() IConnection //获取请求的连接信息
	GetDATa() []byte            //获取请求信息的数据
	GetMsgId() uint32
}

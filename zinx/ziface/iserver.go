package ziface

type IServer interface {
	Start()

	Stop()

	Serve()

	//AddRouter(router IRouter)//此前只有一种消息处理方法，不管来啥都一样的处理函数

	AddRouter(msgId uint32, router IRouter) //现在对不同的消息可以有不同的处理方法

	GetConnMgr() IConnManager //得到链接管理

	SetOnConnStart(func(IConnection)) //设置该server链接创建时的hook函数

	SetOnConnStop(func(IConnection)) //连接断开的hook函数

	CallOnConnStart(conn IConnection) //调用链接创建后需要回调的钩子方法

	CallOnConnStop(conn IConnection) //调用链接停止前的回调的业务方法
}

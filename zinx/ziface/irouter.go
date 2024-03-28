// 路由：在数据通信网络中，具体是指：多台设备进行通信时，彼此之间发送具有IP特征的数据包，当数据包经过具备路由功能的设备时，设备进行解包并查看IP报文的目的网络地址，并于自身维护的路由表条目进行匹配，符合则进行转发，否则丢弃报文，回应网络不可达。
// 网页开发中的网络钩子（Webhook）是一种通过自定义回调函数来增加或更改网页表现的方法。这些回调可被可能与原始网站或应用相关的第三方用户及开发者保存、修改与管理。术语“网络钩子”由杰夫·林德塞（Jeff Lindsay）于2007年通过给计算机编程术语“钩子”（Hook）加上前缀得来
package ziface

type IRouter interface {
	PreHandle(request IRequest)  //处理conn业务之前的钩子方法
	Handle(request IRequest)     //处理conn业务的方法
	PostHandle(request IRequest) //处理业务之后的钩子方法
}

/*
type IServer interface{
	Start()

	Stop()

	Serve()

	AddRouter(router IRouter)
}
*/

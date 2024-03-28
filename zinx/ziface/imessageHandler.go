package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)          //以非阻塞的方式处理消息
	AddRouter(msgId uint32, router IRouter) //传入消息id和对应的路由处理逻辑，就是为消息添加具体的处理逻辑
	StartWorkerPool()                       //启动worker工作池
	SendMsgToTaskQueue(request IRequest)    //将消息交给task queue，等待worker处理
}

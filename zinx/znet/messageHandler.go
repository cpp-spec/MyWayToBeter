package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter //存放每个msgid对应的处理方法的map属性
	WorkerPoolSize uint32                    //工作池worker总数量
	TaskQueue      []chan ziface.IRequest    //工作池的消息队列 注意为什么是存放的request，也说明了worker的作用
}

// 构造函数
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize), //这里要重新理解了，似乎workpoolsize是消息队列的数量
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()] //根据消息id确定对应的处理方法
	if !ok {
		fmt.Println("api msgId=", request.GetMsgId(), "is not found!")
		return
	}

	//执行对应的处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//判断当前msg是否已经绑定了API处理方法
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api,msgId=" + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Add api magId=", msgId)
}

// 启动一个worker，分配id和任务队列；阻塞等待消息队列中请求
func (mh *MsgHandle) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("workerId=", workerId, "is start")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 遍历启动所有worker
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//为任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkTaskLen)
		//创建worker协程
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 将消息发送给taskQueue;注意的是如何决定将消息发给哪个taskQueue
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//用Conn决定该有哪个worker来处理，取余是为了控制取值范围
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("add ConnId=", request.GetConnection().GetConnID(), "request msgId=", request.GetMsgId(), "to workerId=", workerId)
	mh.TaskQueue[workerId] <- request
}

package ziface

// 连接管理器
type IConnManager interface {
	Add(conn IConnection)                   //添加
	Remove(conn IConnection)                //删除
	Get(connId uint32) (IConnection, error) //根据connID获取链接
	Len() int                               //?注意关键字和系统，怪不得要大写首字母，len()是已经有的系统函数
	ClearConn()                             //删除并停止全部链接
}

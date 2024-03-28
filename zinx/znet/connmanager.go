package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //	读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection), //初始化返回map的引用
	}
}

func (ConnMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源，map加锁
	ConnMgr.connLock.Lock()
	defer ConnMgr.connLock.Unlock()
	ConnMgr.connections[conn.GetConnID()] = conn //键连接添加存储 指定id
	fmt.Println("connectin add to Connmanager succ,conn num=", ConnMgr.Len())

}

func (ConnMgr *ConnManager) Remove(conn ziface.IConnection) {
	ConnMgr.connLock.Lock()
	defer ConnMgr.connLock.Unlock()

	delete(ConnMgr.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnId=", conn.GetConnID(), "succ:conn num=", ConnMgr.Len())
}

func (ConnMgr *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	ConnMgr.connLock.Lock()
	defer ConnMgr.connLock.Unlock()
	if conn, ok := ConnMgr.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (ConnMgr *ConnManager) Len() int {
	return len(ConnMgr.connections) //间接访问，一种封装手段；获取当前连接数量
}

func (ConnMgr *ConnManager) ClearConn() {
	ConnMgr.connLock.Lock()
	defer ConnMgr.connLock.Unlock()
	for connId, conn := range ConnMgr.connections {
		conn.Stop()
		delete(ConnMgr.connections, connId)
	}
	fmt.Println("clear allconnections successfully:conn num=", ConnMgr.Len())
}

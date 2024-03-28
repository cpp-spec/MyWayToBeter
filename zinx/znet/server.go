package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// 定义Sever类。是对iserver接口的实现
type Server struct {
	Name      string //服务器名称
	IPversion string //TCP4 or other
	Ip        string //服务器绑定的IP地址
	Port      int    //服务绑定的端口
	//Router    ziface.IRouter //当前Server由用户绑定回调router
	msgHandler  ziface.IMsgHandle             //消息管理模块，用来绑定msgId和对应处理方法
	ConnMgr     ziface.IConnManager           //当前server的连接管理器
	OnConnStart func(conn ziface.IConnection) //函数作为成员变量，传入conn作为参数
	OnConnStop  func(conn ziface.IConnection) //连接断开时的回调函数
}

// 定义当前客户端连结的Handle API
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// 实现启动功能的接口
func (s *Server) Start() {
	fmt.Printf("[START] Server Server name:%s,listenner at IP:%s,Port:%d,is starting\n", s.Name, s.Ip, s.Port)
	//fmt.Printf("[START] Server name:%s,")
	fmt.Printf("zinx version:%s,MaxConn:%d,MaxPacketSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	//开启一个go做服务端listenner业务
	go func() {
		s.msgHandler.StartWorkerPool() //启动工作池
		//获取tcp addr
		addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		//开始监听
		listenner, err := net.ListenTCP(s.IPversion, addr)
		if err != nil {
			fmt.Println("listen", s.IPversion, "err", err)
			return
		}
		//监听启动成功，等待请求
		fmt.Println("start Zinx server", s.Name, "success,now listenning...")

		//自动生成id
		var cid uint32
		cid = 0

		//循环处理连接请求，Accept阻塞主城，收到消息后返回并创建协程异步处理
		for {
			conn, err := listenner.AcceptTCP()
			//conn,err:=listenner.Accept()
			if err != nil {
				fmt.Println("accept err", err)
				continue
			}

			//创建协程异步处理io
			/*	go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}
					//回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()	*/
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close() //超过最大连接数就关闭此链接
				continue
			}
			//创建connection,启动
			dealconn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			go dealconn.Start()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server,name", s.Name)
	s.ConnMgr.ClearConn() //清除所有链接，如果连接中的任务没有完成怎么办呢？
}

func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("Add router succ!")
}

/*
// 创建服务器句柄 V0.4以前没用全局参数模块

	func NewServer(name string) ziface.IServer {
		s := &Server{
			Name:      name,
			IPversion: "tcp4",
			Ip:        "0.0.0.0",
			Port:      7777,
			Router:    nil,
		}
		return s
	}
*/
//获取连接管理的方法
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 设置回调函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("--->CallOnConnStart...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("--->CallOnConnStop...")
		s.OnConnStop(conn)
	}
}
func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPversion:  "tcp4",
		Ip:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(), //创建connmanager
	}
	return s
}

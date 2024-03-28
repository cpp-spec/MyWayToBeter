package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
	//"golang.org/x/tools/go/analysis/passes/errorsas"
)

type Connection struct {
	TcpServer ziface.IServer //当前连接所属的server
	Conn      *net.TCPConn   //当前连接的tcp套接字
	ConnID    uint32         //当前连接的id，全局唯一
	isClosed  bool           //当前连接的关闭状态
	/*handleAPI    ziface.HandFunc //?	*/
	//Router       ziface.IRouter //该链接绑定的处理方法
	ExitBuffChan chan bool //定义一个通道，用来传输bool类型值，告知该链接已经退出停止的channel
	MsgHandler   ziface.IMsgHandle
	msgChan      chan []byte //读，写goroutine之间的通信通道
	msgBuffChan  chan []byte //读，写goroutine之间的通信通道 带缓冲
}

// 创建Connection的函数,类似于构造函数
/*func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection */
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server, //？
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		/*handleAPI:    callback_api,*/
		//Router:       router,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}
	c.TcpServer.GetConnMgr().Add(c) //将新创建的conn添加到管理，Conn新增了server属性，链接而属于server，由server管理
	return c
}
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	fmt.Println("正在向客户端发送信息...")
	if c.isClosed == true {
		return errors.New("connection closed when sendmsg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id:", msgId)
		return errors.New("pack error msg")
	}
	fmt.Println("测试sendmsg() msg大小:", len(msg)) //此处测试
	/*
		if _, err := c.Conn.Write(msg); err != nil {
			fmt.Println("write msg id", msgId, "error")
			c.ExitBuffChan <- true
			return errors.New("conn write error")
		}
	*/
	c.msgChan <- msg //通过channel发送给专门的写模块，意思是任务外包给专门的模块处理，任务交给你了，没完成是你的问题
	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	fmt.Println("正在向客户端发送信息...")
	if c.isClosed == true {
		return errors.New("connection closed when sendmsg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id:", msgId)
		return errors.New("pack error msg")
	}
	//fmt.Println("测试sendmsg() msg大小:", len(msg)) //此处测试
	/*
		if _, err := c.Conn.Write(msg); err != nil {
			fmt.Println("write msg id", msgId, "error")
			c.ExitBuffChan <- true
			return errors.New("conn write error")
		}
	*/
	c.msgBuffChan <- msg //通过channel发送给专门的写模块，意思是任务外包给专门的模块处理，任务交给你了，没完成是你的问题
	return nil
}

// 永久循环，阻塞等待服务器消息；如果收到数据，写入本地内存，数据读取完整后，交给开发者注册的handleAPI处理业务；异常处理是跳过本次循环
// 第一个括号里是方法的接收器，意思是此方法可以由一个指向Connection的指针调用，用来修改Connection结构体
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running...")
	defer fmt.Println(c.Conn.RemoteAddr().String(), "conn reader exit...") //defer函数用于在它包含的函数执行完毕后清理内存
	defer c.Stop()
	for { //将数据读到buf
		/*
			buf := make([]byte, 512)
			_, err := c.Conn.Read(buf)
			if err != nil {
				fmt.Println("recv buf err", err)
				c.ExitBuffChan <- true
				continue
			}
			//执行当前conn绑定的handle方法，可以看出，handle方法实在定义connection时确定的，连接和业务处理是绑定的
			/*if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
				fmt.Println("connID", c.ConnID, "handle is err")
				c.ExitBuffChan <- true
				return
			}*/
		/*
			//得到当前客户端请求的request数据
			req := Request{
				conn: c,
				data: buf,
			}*/
		//读取客户端msg head
		dp := NewDataPack() //创建封包拆包对象
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err", err)
			c.ExitBuffChan <- true
			continue
		}
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err", err)
			c.ExitBuffChan <- true
			continue
		}
		//读取datalen到data
		var data []byte
		if msg.GetDateLen() > 0 {
			data = make([]byte, msg.GetDateLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)
		//终于得到客户端请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		/*
			go func(request ziface.IRequest) {
				c.Router.PreHandle(request)
				c.Router.Handle(request)
				c.Router.PostHandle(request)
			}(&req)
		*/
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req) //执行对应的处理方法
		}

	}

}

// 写消息goroutine,将服务端数据发送至客户端
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")
	for {
		//阻塞等待写任务
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err", err, "writer exit!")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send data err", err, "writer exit!")
					return
				}
			} else {
				fmt.Println("msgBuffChan is closed")
				break
			}

		case <-c.ExitBuffChan: //conn关闭则退出
			return
		}
	}
}
func (c *Connection) Start() {
	//启动协程，异步读取客户端数据和业务
	go c.StartReader()
	//写会客户端协程
	go c.StartWriter()
	//根据用户传进来的创建链接的需处理业务，执行hook方法
	c.TcpServer.CallOnConnStart(c)

	for {
		//阻塞，若是收到退出指令，不再阻塞
		select {
		case <-c.ExitBuffChan:
			return
		}
	}

}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.TcpServer.CallOnConnStop(c) //调用server注册的hook函数
	c.Conn.Close()                //关闭socket链接
	c.ExitBuffChan <- true        //丛此处看出管道通信的意义，实现了与协程任务之间的通信，告知链接已经关闭；
	//将连接从管理器删除
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitBuffChan) //关闭通道
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

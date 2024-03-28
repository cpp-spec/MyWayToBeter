// 测试消息封装模块
package main

import (
	"fmt"
	"io"
	"net"
	"zinx/znet"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listening err:", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)

		}
		go func(conn net.Conn) {
			dp := znet.NewDataPack() //创建封包解包的对象
			for {
				headdata := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headdata)
				if err != nil {
					fmt.Println("read head err:", err)
					break
				}
				msgHead, err := dp.Unpack(headdata)
				if err != nil {
					fmt.Println("server unpack err:", err)
					return
				}
				if msgHead.GetDateLen() > 0 {
					msg := msgHead.(*znet.Message)
					msg.Data = make([]byte, msg.GetDateLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err:", err)
						return
					}
					fmt.Println("==>Recv Msg: Id=", msg.Id, ",len=", msg.DataLen, ",data=", string(msg.Data))
				}
			}

		}(conn)
	}
}

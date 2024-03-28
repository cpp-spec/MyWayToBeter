// 单元测试
package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func client() {
	fmt.Println("client start test...")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}
	for {
		_, err := conn.Write([]byte("hello ZINX"))
		if err != nil {
			fmt.Println("write error err", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err", err)
			return
		}
		fmt.Printf("server call back:%s,cnt=%d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}

func TestServer(t *testing.T) {
	//s := NewServer("[zinx V0.1]")
	s := NewServer("[zinx V0.2]")
	go client()
	s.Serve()
}

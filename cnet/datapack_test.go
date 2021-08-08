package cnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//用于封包和拆包的单元测试
func TestDataPack(t *testing.T) {
	// 服务端
	//创建Tcp连接
	listener, err := net.Listen("tcp", "127.0.0.1:6608")
	if err != nil {
		fmt.Println("server listen error")
		return
	}

	go func() {
		//拆包
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("listener accept error")
				continue
			}

			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					//先读head
					headData := make([]byte, dp.GetHeadLength())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error\n", err)
						return
					}
					if msgHead.GetLength() == 0 { //没有数据
						continue
					}

					//根据读出的head内容再读data
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetLength())
					_, err = io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack error\n", err)
						return
					}

					fmt.Println("Recive Message ID:", msg.ID, " Length:", msg.Length, " Data:", msg.Data)
				}
			}(conn)
		}
	}()

	// 客户端
	conn, err := net.Dial("tcp", "127.0.0.1:6608")
	if err != nil {
		fmt.Println("Client dial error\n", err)
		return
	}
	dp := NewDataPack()
	//连发两个message试试
	//第一个
	msg1 := &Message{
		ID:     1,
		Length: 6,
		Data:   []byte("你妈的"),
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("msg1 data pack error\n", err)
		return
	}
	//第二个
	msg2 := &Message{
		ID:     2,
		Length: 8,
		Data:   []byte("你妈码的"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("msg2 data pack error\n", err)
		return
	}
	//粘起来一起发
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	//阻塞
	select {}
}

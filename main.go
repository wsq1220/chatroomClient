package main

import (
	"fmt"
	"net"
	"github.com/wsq1220/chatroomClient/proto"
)

var userId int
var password string
var msgChan chan proto.UserRecvMsgReq

func init() {
	msgChan = make(chan proto.UserRecvMsgReq, 1000)
}

func main() {
	// 终端输入用户名和密码
	fmt.Println("please input yourid:")
	fmt.Scanf("%d", &userId)
	fmt.Println("please input your password:")
	fmt.Scanf("%s", &password)

	conn, err := net.Dial("tcp", "localhost:10000")
	if err != nil {
		fmt.Printf("client dial failed, err: %v\n", err)
		return
	}

	err = login(conn, userId, password)
	if err != nil {
		fmt.Printf("login failed, err: %v\n", err)
		return
	}

	go processServerMsg(conn)
	for {
		logic(conn)
	}
}
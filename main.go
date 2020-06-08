package main

import (
	"fmt"
	"net"

	"github.com/astaxie/beego/logs"
	"github.com/wsq1220/chatroomClient/proto"
)

var userId int
var password string
var msgChan chan proto.UserRecvMsgReq

func init() {
	msgChan = make(chan proto.UserRecvMsgReq, 1000)
}

func main() {
	if err := initLogger(); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		panic(err)
	}
	fmt.Println("init logger succ!")
	// also can use flag

	// 终端输入用户名和密码
	fmt.Println("please input your id and password, the format as id@password:")
	// 注意换行符  会读取换行符
	fmt.Scanf("%d@%s\n", &userId, &password)
	// fmt.Println("please input your password:")
	// fmt.Scanf("%s", &password)
	logs.Debug("id[%v], password[%v]", userId, password)

	conn, err := net.Dial("tcp", "localhost:10000")
	if err != nil {
		fmt.Printf("client dial failed, err: %v\n", err)
		return
	}

	err = login(conn, userId, password)
	if err != nil {
		fmt.Printf("login failed, err: %v\n", err)
		logs.Error("login failed, err: %v\n", err)
		return
	}

	go processServerMsg(conn)
	for {
		logic(conn)
	}
}

package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/wsq1220/chatroomClient/proto"
)

func listUnreadMsg(conn net.Conn) {
	select {
	case msg := <-msgChan:
		fmt.Println(msg.UserId, ":", msg.Data)
	default:
		fmt.Println("no more")
		return
	}
}

func enterTalk(conn net.Conn) {
	var msg string
	fmt.Println("enter text you want to send:")
	fmt.Scanf("%s", &msg)
	sendTextMsg(conn, msg)
}

func sendTextMsg(conn net.Conn, text string) {
	var msg proto.Message
	msg.Cmd = proto.UserSendMessageCmd

	var sendReq proto.SendMsgReq
	sendReq.UserId = userId
	sendReq.Data = text

	sendData, err := json.Marshal(sendReq)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	msg.Data = string(sendData)

	msgData, err := json.Marshal(msg)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	var buf [8192]byte
	packLen := uint32(len(msgData))
	binary.BigEndian.PutUint32(buf[0:4], packLen)

	_, err = conn.Write(buf[0:4])
	if err != nil {
		logs.Error("writr head data failed, err: %v", err)
		return
	}

	n, err := conn.Write(msgData)
	if err != nil {
		logs.Error("write body data failed, err: %v", err)
		return
	}

	if n != int(packLen) {
		errMsg := fmt.Sprintf("send data not finished! now:%v/%v", n, int(packLen))
		fmt.Println(errMsg)
		logs.Error(errMsg)
	}

	return
}

func enterMenu(conn net.Conn) {
	fmt.Println("------Menu Page--------")
	fmt.Println("1. list all online user")
	fmt.Println("2. list all unread message")
	fmt.Println("3. go to chat")
	fmt.Println("4. exit")
	fmt.Println("5. more")

	var op int
	fmt.Scanf("%d\n", &op)
	switch op {
	case 1:
		listUseronline()
	case 2:
		listUnreadMsg(conn)
	case 3:
		enterTalk(conn)
	case 4:
		os.Exit(0)
	case 5:
		fmt.Println("sorry, no more")
		os.Exit(0)
	}
}

func logic(conn net.Conn) {
	enterMenu(conn)
}

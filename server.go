package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/fatih/color"
	"github.com/wsq1220/chatroomClient/proto"
)

func processServerMsg(conn net.Conn) {
	logs.Debug("entered the processServerMsg")
	defer func() {
		if err := recover(); err != nil {
			logs.Error("Error when processing server msg in goroutine, err: %v", err)
			return
		}
	}()

	for {
		msg, err := readPackage(conn)
		if err != nil {
			logs.Error("read poackage failed, err: %v", err)
			os.Exit(0)
		}

		var userStatus proto.UserStatusNotify
		err = json.Unmarshal([]byte(msg.Data), &userStatus)
		if err != nil {
			logs.Error("json unmarshal failed, err: %v", err)
			return
		}

		switch msg.Cmd {
		case proto.UserStatusNotifyCmd:
			updateUserStatus(userStatus)
		case proto.UserRecvMessageCmd:
			recvMsgFromServer(msg)
		default:
			errMsg := fmt.Sprintf("[%v] is not supported!", msg.Cmd)
			err = fmt.Errorf(errMsg)
			// 需要显示在终端让用户知道
			color.Set(color.FgRed, color.Bold)
			fmt.Println(errMsg)
			color.Unset()
			logs.Error(errMsg)
			return
		}
	}
}

func recvMsgFromServer(msg proto.Message) {
	var recvMsgReq proto.UserRecvMsgReq
	err := json.Unmarshal([]byte(msg.Data), &recvMsgReq)
	if err != nil {
		logs.Error("json unmarshal failed, err: %v", err)
		return
	}
	msgStr := fmt.Sprintf("[%d]: %s", recvMsgReq.UserId, recvMsgReq.Data)
	// 显示用户消息在在终端
	color.Set(color.FgCyan)
	fmt.Println(msgStr)
	color.Unset()

	// 发送消息到管道
	msgChan <- recvMsgReq
}

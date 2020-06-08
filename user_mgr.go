package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/fatih/color"
	"github.com/wsq1220/chatroomClient/proto"
)

var onlineUserMap map[int]*proto.User = make(map[int]*proto.User, 100)

func listUseronline() {
	color.Set(color.FgCyan, color.Bold)
	fmt.Println("now online user:")
	color.Unset()
	// 除本人
	for id, _ := range onlineUserMap {
		if id == userId {
			continue
		}
		color.Set(color.FgHiBlue, color.Bold)
		fmt.Println(id)
		color.Unset()
	}
}

func updateUserStatus(userStatus proto.UserStatusNotify) {
	user, ok := onlineUserMap[userStatus.UserId]
	if !ok {
		logs.Warn("the user not online or not exist, create one new")
		user = &proto.User{}
		user.UserId = userStatus.UserId
	}

	user.Status = userStatus.Status
	onlineUserMap[userStatus.UserId] = user

	listUseronline()
}

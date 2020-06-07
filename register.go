package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"github.com/astaxie/beego/logs"
	"github.com/wsq1220/chatroomClient/proto"
)

func register(conn net.Conn, userId int, password string) (err error) {
	var registerMsg proto.Message
	registerMsg.Cmd = proto.UserRegisterCmd

	var register proto.Register
	register.User.UserId = userId
	// TODO use md5
	register.User.Password = password
	register.User.Alias = fmt.Sprintf("v-User%d", userId)
	register.User.Gender = "Unknown"
	register.User.Avatar = fmt.Sprintf("./img/%d.png", userId)
	logs.Debug("the user will register: %v", register.User)

	data, err := json.Marshal(register)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}
	registerMsg.Data = string(data)

	registerMsgData, err := json.Marshal(registerMsg)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	// TODO封装该方法
	var buf [4]byte
	packLen := uint32(len(registerMsgData))
	logs.Info("len of package: %v", packLen)

	binary.BigEndian.PutUint32(buf[0:4], packLen)

	n, err := conn.Write(buf[:])
	if err != nil || n != 4 {
		logs.Error("write data failed when registering, err: %v", err)
		return
	}

	n, err = conn.Write([]byte(data))
	if err != nil {
		logs.Error("write body data failed when registering, err: %v", err)
		return
	}

	// msg, err := readPackage(conn)
	// if err != nil {
	// 	logs.Error("read package failed, err: %v", err)
	// 	return
	// }
	// fmt.Println(msg)

	return
}

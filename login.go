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

// TODO
func login(conn net.Conn,userId int, password string) (err error) {
	var loginMsg proto.Message
	loginMsg.Cmd = proto.UserLoginCmd

	var login proto.Login
	login.Id = userId
	login.Password = password

	loginData, err := json.Marshal(login)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	loginMsg.Data = string(loginData)

	data, err := json.Marshal(loginMsg)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	var buf [8192]byte
	packLen := uint32(len(data))
	binary.BigEndian.PutUint32(buf[0:4], packLen)

	n, err := conn.Write([]byte(buf[0:4]))
	if err != nil || n != 4 {
		logs.Error("write head data failed, err: %v", err)
		return
	}

	n, err  = conn.Write(data)
	if err != nil {
		logs.Error("write body data failed, err: %v", err)
		return
	}

	if n != int(packLen) {
		logs.Error("send data not finished!")
		return
	}

	loginMsg, err = readPackage(conn)
	if err != nil {
		logs.Error("read package failed, err: %v", err)
		return
	}

	var loginResp proto.LoginResp
	err = json.Unmarshal([]byte(loginMsg.Data), &loginResp)
	if err != nil {
		logs.Error("json unmarshal failed, err: %v", err)
		return
	}

	if loginResp.StatusCode == 500 {
		fmt.Printf("user %v not register, start registering...\n", userId)
		logs.Error("user %v not register, start registering...", userId)
		if err = register(conn, userId, password); err != nil {
			fmt.Printf("register failed, err: %v\n", err)
			return
		}
		os.Exit(0)
	}

	for _, v := range loginResp.User {
		if v == userId {
			continue
		}
		fmt.Println("user logined: %v", v)

		// 添加到map中
		user := &proto.User{
			UserId: v,
		}
		onlineUserMap[v] = user
	}

	return
}
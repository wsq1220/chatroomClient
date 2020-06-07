package main

import (
	"net"
	"encoding/binary"
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/wsq1220/chatroomClient/proto"
)


func readPackage(conn net.Conn) (msg proto.Message, err error) {
	var buf [8192]byte
	n, err := conn.Read(buf[0:4])
	logs.Debug("the val of n: %v", n)
	if err != nil{
		logs.Error("client read head failed, err: %v", err)
		return
	}

	var packLen uint32
	packLen = binary.BigEndian.Uint32(buf[0:4])

	n, err = conn.Read(buf[0:packLen])
	if err != nil {
		logs.Error("client read body data failed, err: %v", err)
		return
	}

	err = json.Unmarshal([]byte(buf[0:packLen]), &msg)
	if err != nil {
		logs.Error("json unmarshal failed, err: %v", err)
		return
	}

	return
}
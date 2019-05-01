package main

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/gotcp"
)

// User : 登录玩家类
type User struct {
	gotcp.Session
	account string
}

// OnRecv : 接收到网络数据包，被触发
func (user *User) OnRecv(data []byte, flag byte) {
	cmd := gotcp.GetCmd(data)
	if user.IsVerified() == false && user.doVerify(protocol.CMD_LOGIN_ENUM(cmd), data, flag) == false {
		return
	}
	switch protocol.CMD_LOGIN_ENUM(cmd) {
	case protocol.CMD_LOGIN_LOGIN:
		// No need to do anything
	default:
		Ctx.Errorln("Unknown message number, message number is", cmd)
	}
}

// OnClose : 断开连接，被触发
func (user *User) OnClose() {
	if user.account != "" {
		Ctx.Infoln("Account connection disconnected, account:", user.account)
	}
}

func (user *User) doVerify(cmd protocol.CMD_LOGIN_ENUM, data []byte, flag byte) bool {
	if cmd != protocol.CMD_LOGIN_LOGIN {
		Ctx.Errorln("The expected message number is `protocol.CMD_LOGIN_LOGIN`(", int(protocol.CMD_LOGIN_LOGIN), "), but", cmd)
		user.Close()
		return false
	}
	msg := &protocol.MSG_LOGIN{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_LOGIN_LOGIN`(", int(protocol.CMD_LOGIN_LOGIN), ")")
		user.Close()
		return false
	}
	user.doLogin(msg.GetAccount(), msg.GetPassword(), msg.GetMode(), msg.GetUserdata())
	user.Verify()
	user.CloseAfterSending()
	return true
}

func (user *User) doLogin(account, passwd string, mode protocol.ENUM_LOGIN_MODE_ENUM, userdata []byte) {
	Ctx.Infoln("account =", account, "password =", passwd, "mode =", mode)
	user.account = account
	token, addrs, ports, nodeTypes, errCode := Ctx.Login(account, passwd, mode == protocol.ENUM_LOGIN_MODE_DEFAULT, userdata)
	if errCode == common.LoginSuccess {
		Ctx.Infoln("account =", account, "token =", token, "addr =", addrs, "port =", ports, "nodeTypes =", nodeTypes)
		msg := &protocol.MSG_LOGIN_RESULT{}
		msg.Err = protocol.ENUM_LOGIN_ERROR_OK
		msg.Token = token
		msg.Address = append(msg.Address, addrs...)
		msg.Port = append(msg.Port, ports...)
		for _, v := range nodeTypes {
			msg.NodeTyps = append(msg.NodeTyps, int32(v))
		}
		user.SendMsg(uint64(protocol.CMD_LOGIN_LOGIN), msg)
	} else {
		Ctx.Errorln("Login failed. error =", errCode, "account =", account)
		msg := &protocol.MSG_LOGIN_RESULT{}
		msg.Err = protocol.ENUM_LOGIN_ERROR_ENUM(errCode)
		user.SendMsg(uint64(protocol.CMD_LOGIN_LOGIN), msg)
	}
}

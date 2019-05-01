package main

import (
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
	if flag != 0 {
		Ctx.Errorln("Packet flag field error. account:", user.account)
		user.Close()
		return
	}
	cmd := gotcp.GetCmd(data)
	if user.IsVerified() == false && user.doVerify(protocol.CMD_GATEWAY_ENUM(cmd), data, flag) == false {
		return
	}
	switch protocol.CMD_GATEWAY_ENUM(cmd) {
	case protocol.CMD_GATEWAY_VERIFY_TOKEN:
		// No need to do anything
	default:
		data = append(data, flag)
		if Ctx.OnRecvFromClient(user.account, uint32(cmd), data) == false {
			Ctx.Errorln("Unknown message number, message number is", cmd, "account:", user.account)
			user.Close()
		}
	}
}

// OnClose : 断开连接，被触发
func (user *User) OnClose() {
	gateway.DelUser(user.account)
	Ctx.Infoln("Account connection disconnected, account:", user.account)
}

func (user *User) doVerify(cmd protocol.CMD_GATEWAY_ENUM, data []byte, flag byte) bool {
	if cmd != protocol.CMD_GATEWAY_VERIFY_TOKEN {
		Ctx.Errorln("The expected message number is `protocol.CMD_GATEWAY_VERIFY_TOKEN`(", int(protocol.CMD_GATEWAY_VERIFY_TOKEN), "), but", cmd)
		user.Close()
		return false
	}
	msg := &protocol.MSG_GATEWAY_VERIFY_TOKEN{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_GATEWAY_VERIFY_TOKEN`(", int(protocol.CMD_GATEWAY_VERIFY_TOKEN), ")")
		user.Close()
		return false
	}
	errcode := Ctx.VerifyToken(msg.GetAccount(), msg.GetToken(), user)
	repmsg := &protocol.MSG_GATEWAY_VERIFY_TOKEN_RESULT{}
	repmsg.Err = protocol.ENUM_GATEWAY_VERIFY_TOKEN_ERROR_ENUM(errcode)
	if errcode != 0 {
		Ctx.Errorln("Token verification failed, account:", msg.GetAccount(), "errcode:", errcode)
		user.SendMsg(uint64(protocol.CMD_GATEWAY_VERIFY_TOKEN), repmsg)
		user.Close()
		return false
	}
	user.Verify()
	if user.SendMsg(uint64(protocol.CMD_GATEWAY_VERIFY_TOKEN), repmsg) == false {
		Ctx.Errorln("Failed to send data, account:", msg.GetAccount())
		user.Close()
		return false
	}
	user.account = msg.GetAccount()
	kickOldUser := gateway.AddUser(user.account, user)
	if kickOldUser {
		Ctx.Infoln("Delete old player object, account:", user.account)
	}
	Ctx.Infoln("Token verification succeeded, account:", msg.GetAccount())
	return true
}

// GetAccount : 获取账号名
func (user *User) GetAccount() string {
	return user.account
}

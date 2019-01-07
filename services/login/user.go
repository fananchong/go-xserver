package main

import (
	"github.com/fananchong/gotcp"

	"github.com/fananchong/go-xserver/common"
	proto_login "github.com/fananchong/go-xserver/services/internal/protocol"
)

// User : 登录玩家类
type User struct {
	gotcp.Session
}

// OnRecv : 接收到网络数据包，被触发
func (user *User) OnRecv(data []byte, flag byte) {
	cmd := gotcp.GetCmd(data)
	if user.IsVerified() == false && user.doVerify(proto_login.CMD_LOGIN_ENUM(cmd), data, flag) == false {
		return
	}
	switch proto_login.CMD_LOGIN_ENUM(cmd) {
	case proto_login.CMD_LOGIN_LOGIN:
		// do nothing
	default:
		common.XLOG.Errorln("unknow cmd, cmd =", cmd)
	}
}

// OnClose : 断开连接，被触发
func (user *User) OnClose() {

}

func (user *User) doVerify(cmd proto_login.CMD_LOGIN_ENUM, data []byte, flag byte) bool {
	if cmd != proto_login.CMD_LOGIN_LOGIN {
		return false
	}
	msg := &proto_login.MSG_LOGIN{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		user.Close()
		return false
	}
	user.doLogin(msg.GetAccount(), msg.GetPassword(), msg.GetMode(), msg.GetUserdata())
	user.Verify()
	user.Close()
	return true
}

func (user *User) doLogin(account, passwd string, mode proto_login.ENUM_LOGIN_MODE_ENUM, userdata []byte) {
	common.XLOG.Infoln("account =", account, "password =", passwd, "mode =", mode)
	token, addr, port, errCode := common.XLOGIN.Login(account, passwd, mode == proto_login.ENUM_LOGIN_MODE_DEFAULT, userdata)
	if errCode == common.LoginSuccess {
		common.XLOG.Infoln("token =", token, "addr =", addr, "port =", port)
		msg := &proto_login.MSG_LOGIN_RESULT{}
		msg.Err = proto_login.ENUM_LOGIN_ERROR_OK
		msg.Token = token
		msg.Address = addr
		msg.Port = port
		user.SendMsg(uint64(proto_login.CMD_LOGIN_LOGIN), msg)
	} else {
		common.XLOG.Errorln("login fail. error =", errCode)
		msg := &proto_login.MSG_LOGIN_RESULT{}
		msg.Err = proto_login.ENUM_LOGIN_ERROR_ENUM(errCode)
		user.SendMsg(uint64(proto_login.CMD_LOGIN_LOGIN), msg)
	}
}

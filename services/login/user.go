package main

import (
	"github.com/fananchong/go-xserver/common"
	proto_login "github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/gotcp"
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
		Ctx.Log.Errorln("unknow cmd, cmd =", cmd)
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
	user.CloseAfterSending()
	return true
}

func (user *User) doLogin(account, passwd string, mode proto_login.ENUM_LOGIN_MODE_ENUM, userdata []byte) {
	Ctx.Log.Infoln("account =", account, "password =", passwd, "mode =", mode)
	token, addrs, ports, nodeTypes, errCode := Ctx.Login.Login(account, passwd, mode == proto_login.ENUM_LOGIN_MODE_DEFAULT, userdata)
	if errCode == common.LoginSuccess {
		Ctx.Log.Infoln("account =", account, "token =", token, "addr =", addrs, "port =", ports, "nodeTypes =", nodeTypes)
		msg := &proto_login.MSG_LOGIN_RESULT{}
		msg.Err = proto_login.ENUM_LOGIN_ERROR_OK
		msg.Token = token
		msg.Address = append(msg.Address, addrs...)
		msg.Port = append(msg.Port, ports...)
		for _, v := range nodeTypes {
			msg.NodeTyps = append(msg.NodeTyps, int32(v))
		}
		user.SendMsg(uint64(proto_login.CMD_LOGIN_LOGIN), msg)
	} else {
		Ctx.Log.Errorln("login fail. error =", errCode, "account =", account)
		msg := &proto_login.MSG_LOGIN_RESULT{}
		msg.Err = proto_login.ENUM_LOGIN_ERROR_ENUM(errCode)
		user.SendMsg(uint64(proto_login.CMD_LOGIN_LOGIN), msg)
	}
}

package main

import (
	"net/http"

	"github.com/fananchong/go-xserver/common"
	proto_login "github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/gogo/protobuf/proto"
)

func (login *Login) msgLogin(w http.ResponseWriter, r *http.Request) {
	data, ok := check(w, r)
	if !ok {
		w.Write(getErrRepString(proto_login.ENUM_LOGIN_ERROR_VERIFY_FAIL))
		return
	}
	msg := &proto_login.MSG_LOGIN{}
	if decodeMsg(data, msg) == nil {
		w.Write(getErrRepString(proto_login.ENUM_LOGIN_ERROR_VERIFY_FAIL))
		return
	}
	common.XLOG.Infoln("account =", msg.GetAccount(), "password =", msg.GetPassword(), "mode =", msg.GetMode())
	token, addr, port, errCode := common.XLOGIN.Login(msg.GetAccount(),
		msg.GetPassword(),
		msg.GetMode() == proto_login.ENUM_LOGIN_MODE_DEFAULT,
		msg.GetUserdata())
	if errCode != common.LoginSuccess {
		common.XLOG.Errorln("login fail. error =", errCode)
		w.Write(getErrRepString(proto_login.ENUM_LOGIN_ERROR_ENUM(errCode)))
		return
	}
	common.XLOG.Infoln("token =", token, "addr =", addr, "port =", port)
	rep := &proto_login.MSG_LOGIN_RESULT{}
	rep.Err = proto_login.ENUM_LOGIN_ERROR_OK
	rep.Token = token
	rep.Address = addr
	rep.Port = port
	succmsg, _ := proto.Marshal(rep)
	w.Write(succmsg)
}

package main

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/fananchong/go-xserver/common"
	proto_login "github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
)

func check(w http.ResponseWriter, req *http.Request) (data string, ok bool) {
	req.ParseForm()
	paramt, ok1 := req.Form["t"]
	paramd, ok2 := req.Form["d"]
	params, ok3 := req.Form["s"]
	if !ok1 || !ok2 || !ok3 {
		common.XLOG.Errorln("http request param error!")
		return
	}
	s1 := []byte(common.XCONFIG.Login.Sign1 + paramd[0] +
		common.XCONFIG.Login.Sign2 + paramt[0] +
		common.XCONFIG.Login.Sign3 + common.XCONFIG.Common.Version)
	s2 := md5.Sum(s1)
	s3 := fmt.Sprintf("%x", s2)
	if s3 != params[0] {
		common.XLOG.Errorln("sign error!",
			"client sign =", params[0],
			"server sign =", s3,
			"sign1 =", common.XCONFIG.Login.Sign1,
			"sign2 =", common.XCONFIG.Login.Sign2,
			"sign3 =", common.XCONFIG.Login.Sign3,
			"t =", paramt[0],
			"version =", common.XCONFIG.Common.Version)
		return
	}
	data = paramd[0]
	ok = true
	return
}

func getErrRepString(err proto_login.ENUM_LOGIN_ERROR_ENUM) []byte {
	rep := &proto_login.MSG_LOGIN_RESULT{}
	rep.Err = err
	data, _ := proto.Marshal(rep)
	return data
}

func decodeMsg(data string, msg proto.Message) proto.Message {
	return gotcp.DecodeCmdEx([]byte(data), 0, msg, 0)
}

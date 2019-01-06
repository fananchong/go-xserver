package components

import "github.com/fananchong/go-xserver/common"

// Login : 登陆模块
type Login struct {
}

// Start : 启动
func (login *Login) Start() bool {
	if getPluginType() == common.Login {
		common.XLOGIN = login
	}
	return true
}

// RegisterCustomAccountVerification : 注册回调
func (login *Login) RegisterCustomAccountVerification(f common.FuncTypeAccountVerification) {
}

// Login : 登陆处理
func (login *Login) Login(account, password string, defaultMode bool, userdata []byte) (token, address string, port int32, errcode common.LoginErrCode) {
	return
}

// Close : 关闭
func (login *Login) Close() {

}

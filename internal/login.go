package internal

import "github.com/fananchong/go-xserver/common"

// Login : 登陆模块
type Login struct {
}

// NewLogin : 登陆模块构造函数
func NewLogin() *Login {
	return &Login{}
}

// Start : 启动
func (login *Login) Start() bool {
	return false
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

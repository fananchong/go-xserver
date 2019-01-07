package main

import (
	"github.com/fananchong/go-xserver/common"
)

// Login : 登陆服务器
type Login struct {
}

// NewLogin : 构造函数
func NewLogin() *Login {
	return &Login{}
}

// Init : 初始化
func (login *Login) Init() {
	common.XLOGIN.RegisterCustomAccountVerification(login.customVerify)
	common.XTCPSVRFORCLIENT.RegisterSessType(User{})
}

// Start : 启动
func (login *Login) Start() bool {
	return true
}

// Close : 关闭
func (login *Login) Close() {

}

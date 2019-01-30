package main

import "github.com/fananchong/go-xserver/common"

// Login : 登陆服务器
type Login struct {
}

// NewLogin : 构造函数
func NewLogin() *Login {
	return &Login{}
}

// Start : 启动
func (login *Login) Start() bool {
	Ctx.Login.RegisterCustomAccountVerification(login.customVerify)
	Ctx.Login.RegisterAllocationNodeType([]common.NodeType{common.Gateway, common.Lobby})
	Ctx.ServerForClient.RegisterSessType(User{})
	return true
}

// Close : 关闭
func (login *Login) Close() {

}

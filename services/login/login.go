package main

import (
	"github.com/fananchong/go-xserver/config"
)

// Login : 登陆服务器
type Login struct {
}

// NewLogin : 构造函数
func NewLogin() *Login {
	return &Login{}
}

// Start : 启动
func (login *Login) Start() bool {
	Ctx.RegisterCustomAccountVerification(login.customVerify)
	//Ctx.RegisterAllocationNodeType([]config.NodeType{config.Gateway}) // Gateway 会随机中继 Lobby
	Ctx.RegisterAllocationNodeType([]config.NodeType{config.Gateway, config.Lobby}) // Gateway 会状态中继 Lobby
	Ctx.RegisterSessType(User{})
	return true
}

// Close : 关闭
func (login *Login) Close() {

}

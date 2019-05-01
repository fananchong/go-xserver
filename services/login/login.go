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
	Ctx.RegisterCustomAccountVerification(login.customVerify)
	//Ctx.RegisterAllocationNodeType([]common.NodeType{common.Gateway}) // Gateway 会随机中继 Lobby
	Ctx.RegisterAllocationNodeType([]common.NodeType{common.Gateway, common.Lobby}) // Gateway 会状态中继 Lobby
	Ctx.RegisterSessType(User{})
	return true
}

// Close : 关闭
func (login *Login) Close() {

}

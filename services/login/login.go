package main

import (
	"github.com/fananchong/go-xserver/common/config"
)

// LoginConfig : 登录服务器配置
type LoginConfig struct {
	AllocServers []config.NodeType `default:"[3]" desc:"给账号分配哪些类型服务器"`
}

// Login : 登陆服务器
type Login struct {
	cfg *LoginConfig
}

// NewLogin : 构造函数
func NewLogin() *Login {
	return &Login{
		cfg: &LoginConfig{},
	}
}

// Start : 启动
func (login *Login) Start() bool {
	if ok := Ctx.LoadConfig("login.toml", login.cfg); !ok {
		Ctx.Errorln("Load config fail, config: login.toml")
		return false
	}
	Ctx.RegisterCustomAccountVerification(login.customVerify)
	//Ctx.RegisterAllocationNodeType([]config.NodeType{config.Gateway}) // Gateway 会随机中继 Lobby
	//Ctx.RegisterAllocationNodeType([]config.NodeType{config.Gateway, services.Lobby}) // Gateway 会状态中继 Lobby
	Ctx.RegisterAllocationNodeType(login.cfg.AllocServers)
	Ctx.RegisterSessType(User{})
	return true
}

// Close : 关闭
func (login *Login) Close() {

}

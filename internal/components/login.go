package components

import (
	"github.com/fananchong/go-xserver/common"
)

// Login : 登陆模块
type Login struct {
	ctx              *common.Context
	verificationFunc common.FuncTypeAccountVerification
	allocType        []common.NodeType
}

// NewLogin : 实例化登陆模块
func NewLogin(ctx *common.Context) *Login {
	login := &Login{
		ctx: ctx,
	}
	login.ctx.Login = login
	return login
}

// Start : 启动
func (login *Login) Start() bool {
	if getPluginType(login.ctx) == common.Login {

	}
	return true
}

// RegisterCustomAccountVerification : 注册回调
func (login *Login) RegisterCustomAccountVerification(f common.FuncTypeAccountVerification) {
	login.verificationFunc = f
}

// RegisterAllocationNodeType : 注册要分配的服务器资源类型
func (login *Login) RegisterAllocationNodeType(types []common.NodeType) {
	login.allocType = append(login.allocType, types...)
}

// Login : 登陆处理
func (login *Login) Login(account, password string, defaultMode bool, userdata []byte) (token, address string, port int32, errcode common.LoginErrCode) {
	return
}

// Close : 关闭
func (login *Login) Close() {

}

package components

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/gotcp"
)

// 前置初始化好某些相互依赖的组件
func init() {
	common.XTCPSVRFORCLIENT = &gotcp.Server{}
	common.XTCPSVRFORINTRANET = &gotcp.Server{}
	common.XLOGIN = &Login{}
}

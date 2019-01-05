package internal

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/gotcp"
)

func initTCPServer() {
	server1 := &gotcp.Server{}
	server1.SetAddress(utils.GetIPOuter(), utils.GetDefaultServicePort())
	server1.SetUnfixedPort(true)
	common.XTCPSVRFORCLIENT = server1
	server2 := &gotcp.Server{}
	server2.SetAddress(utils.GetIPInner(), utils.GetIntranetListenPort())
	common.XTCPSVRFORINTRANET = server2
}

package components

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/gotcp"
)

// TCPServer : TCP Server 组件
type TCPServer struct {
}

// Start : 实例化组件
func (*TCPServer) Start() bool {
	if common.XTCPSVRFORCLIENT = startTCPServer(utils.GetIPOuter(), utils.GetDefaultServicePort()); common.XTCPSVRFORCLIENT == nil {
		return false
	}
	if common.XTCPSVRFORINTRANET = startTCPServer(utils.GetIPInner(), utils.GetIntranetListenPort()); common.XTCPSVRFORINTRANET == nil {
		return false
	}
	common.XCONFIG.Network.Port[0] = common.XTCPSVRFORCLIENT.(*gotcp.Server).GetRealPort()
	common.XCONFIG.Network.Port[1] = common.XTCPSVRFORINTRANET.(*gotcp.Server).GetRealPort()
	return true
}

// Close : 关闭组件
func (*TCPServer) Close() {
	if common.XTCPSVRFORCLIENT != nil {
		common.XTCPSVRFORCLIENT.(*gotcp.Server).Close()
		common.XTCPSVRFORCLIENT = nil
	}
	if common.XTCPSVRFORINTRANET != nil {
		common.XTCPSVRFORINTRANET.(*gotcp.Server).Close()
		common.XTCPSVRFORCLIENT = nil
	}
}

func startTCPServer(addr string, port int32) common.ITCPServer {
	s := &gotcp.Server{}
	s.SetAddress(addr, port)
	s.SetUnfixedPort(true)
	if s.Start() {
		return s
	}
	return nil
}

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
	loadPlugin()
	s := common.XTCPSVRFORCLIENT.(*gotcp.Server)
	if s.GetSessionType() != nil {
		if !startTCPServer(s, utils.GetIPOuter(), utils.GetDefaultServicePort()) {
			return false
		}
		common.XCONFIG.Network.Port[0] = s.GetRealPort()
	}
	s = common.XTCPSVRFORINTRANET.(*gotcp.Server)
	if s.GetSessionType() != nil {
		if !startTCPServer(s, utils.GetIPInner(), utils.GetIntranetListenPort()) {
			return false
		}
		common.XCONFIG.Network.Port[1] = s.GetRealPort()
	}
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

func startTCPServer(s *gotcp.Server, addr string, port int32) bool {
	s.SetAddress(addr, port)
	s.SetUnfixedPort(true)
	return s.Start()
}

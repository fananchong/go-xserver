package components

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/gotcp"
)

// TCPServer : TCP Server 组件
type TCPServer struct {
	ctx *common.Context
}

// NewTCPServer : 实例化
func NewTCPServer(ctx *common.Context) *TCPServer {
	server := &TCPServer{ctx: ctx}
	server.ctx.ServerForClient = &gotcp.Server{}
	server.ctx.ServerForIntranet = &gotcp.Server{}
	return server
}

// Start : 实例化组件
func (server *TCPServer) Start() bool {
	loadPlugin(server.ctx)
	s := server.ctx.ServerForClient.(*gotcp.Server)
	if s.GetSessionType() != nil {
		if !startTCPServer(s, utils.GetIPOuter(server.ctx), utils.GetDefaultServicePort(server.ctx)) {
			return false
		}
		server.ctx.Config.Network.Port[common.PORTFORCLIENT] = s.GetRealPort()
	}
	s = server.ctx.ServerForIntranet.(*gotcp.Server)
	if s.GetSessionType() != nil {
		if !startTCPServer(s, utils.GetIPInner(server.ctx), utils.GetIntranetListenPort(server.ctx)) {
			return false
		}
		server.ctx.Config.Network.Port[common.PORTFORINTRANET] = s.GetRealPort()
	}
	return true
}

// Close : 关闭组件
func (server *TCPServer) Close() {
	if server.ctx.ServerForClient != nil {
		server.ctx.ServerForClient.(*gotcp.Server).Close()
		server.ctx.ServerForClient = nil
	}
	if server.ctx.ServerForIntranet != nil {
		server.ctx.ServerForIntranet.(*gotcp.Server).Close()
		server.ctx.ServerForIntranet = nil
	}
}

func startTCPServer(s *gotcp.Server, addr string, port int32) bool {
	s.SetAddress(addr, port)
	s.SetUnfixedPort(true)
	return s.Start()
}

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
	server.ctx.ITCPServer = &gotcp.Server{}
	return server
}

// Start : 实例化组件
func (server *TCPServer) Start() bool {
	s := server.ctx.ITCPServer.(*gotcp.Server)
	if s.GetSessionType() != nil {
		// 填写 `0.0.0.0` ， 而不是 `utils.GetIPOuter(server.ctx)` 具体外网 IP ，是为了能支持阿里云 ECS 服务器
		if !startTCPServer(s, "0.0.0.0", utils.GetDefaultServicePort(server.ctx)) {
			return false
		}
		server.ctx.Config.Network.Port[common.PORTFORCLIENT] = s.GetRealPort()
	}
	return true
}

// Close : 关闭组件
func (server *TCPServer) Close() {
	if server.ctx.ITCPServer != nil {
		server.ctx.ITCPServer.(*gotcp.Server).Close()
		server.ctx.ITCPServer = nil
	}
}

func startTCPServer(s *gotcp.Server, addr string, port int32) bool {
	s.SetAddress(addr, port)
	s.SetUnfixedPort(true)
	return s.Start()
}

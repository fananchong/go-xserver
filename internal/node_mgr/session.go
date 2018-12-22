package nodemgr

import (
	"github.com/fananchong/gotcp"
)

// Session : 网络会话类
type Session struct {
	gotcp.Session
}

// OnRecv : 接收到网络数据包，被触发
func (sess *Session) OnRecv(data []byte, flag byte) {
	if sess.IsVerified() == false {
		sess.Verify()
	}
}

// OnClose : 断开连接，被触发
func (sess *Session) OnClose() {
}

package nodenormal

import (
	"fmt"
	"time"

	"github.com/fananchong/gotcp"
)

// Session : 网络会话类
type Session struct {
	*gotcp.Session
}

// Start : 开始连接 Mgr Server
func (sess *Session) Start() bool {
	sess.connectMgrServer()
	return true
}

func (sess *Session) connectMgrServer() {
	sess.Session = &gotcp.Session{}
TRY_AGAIN:
	addr, port := getMgrInfoByBlock()
	if sess.Connect(fmt.Sprintf("%s:%d", addr, port), sess) == false {
		time.Sleep(1 * time.Second)
		goto TRY_AGAIN
	}
}

// OnRecv : 接收到网络数据包，被触发
func (sess *Session) OnRecv(data []byte, flag byte) {
	if sess.IsVerified() == false {
		sess.Verify()
	}
}

// OnClose : 断开连接，被触发
func (sess *Session) OnClose() {
	go sess.connectMgrServer()
}

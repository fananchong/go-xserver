package components

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fananchong/go-xserver/common"
)

// Signal : 信号处理组件
type Signal struct {
	ctx *common.Context
	sig chan os.Signal
}

// NewSignal : 实例化
func NewSignal(ctx *common.Context) *Signal {
	return &Signal{ctx: ctx}
}

// Start : 实例化组件
func (s *Signal) Start() bool {
	OneComponentOK(s.ctx.Ctx)
	s.sig = make(chan os.Signal)
	signal.Notify(s.sig,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGTERM,
		syscall.SIGPIPE)

	for {
		select {
		case sig := <-s.sig:
			switch sig {
			case syscall.SIGPIPE:
			default:
				s.ctx.Log.Infoln("[app] recive signal:", sig)
				return true
			}
		}
	}
}

// Close : 关闭组件
func (s *Signal) Close() {
	// do nothing
}

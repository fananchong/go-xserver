package components

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fananchong/go-xserver/common"
)

// Signal : 信号处理组件
type Signal struct {
	sig chan os.Signal
}

// Start : 实例化组件
func (s *Signal) Start() bool {
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
				common.XLOG.Infoln("[app] recive signal:", sig)
				return true
			}
		}
	}
}

// Close : 关闭组件
func (s *Signal) Close() {
	// do nothing
}

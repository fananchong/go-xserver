package components

import (
	"net/http"
	_ "net/http/pprof" // 只使用 pprof 包的 init 函数

	"github.com/fananchong/go-xserver/common"
)

// Pprof : http pprof组件
type Pprof struct {
	web *http.Server
}

// Start : 实例化组件
func (pprof *Pprof) Start() bool {
	addr := common.XCONFIG.Common.Pprof
	if addr != "" {
		go func() {
			common.XLOG.Println("pprof listen :", addr)
			pprof.web = &http.Server{Addr: addr, Handler: nil}
			pprof.web.ListenAndServe()
		}()
	}
	return true
}

// Close : 关闭组件
func (pprof *Pprof) Close() {
	if pprof.web != nil {
		pprof.web.Close()
		pprof.web = nil
	}
}

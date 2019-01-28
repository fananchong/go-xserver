package components

import (
	"net/http"
	_ "net/http/pprof" // 只使用 pprof 包的 init 函数

	"github.com/fananchong/go-xserver/common"
)

// Pprof : http pprof组件
type Pprof struct {
	ctx *common.Context
	web *http.Server
}

// NewPprof : 实例化
func NewPprof(ctx *common.Context) *Pprof {
	return &Pprof{ctx: ctx}
}

// Start : 实例化组件
func (pprof *Pprof) Start() bool {
	addr := pprof.ctx.Config.Common.Pprof
	if addr != "" {
		go func() {
			pprof.ctx.Log.Println("pprof listen :", addr)
			pprof.web = &http.Server{Addr: addr, Handler: nil}
			pprof.web.ListenAndServe()
		}()
	}
	OneComponentOK(pprof.ctx.Ctx)
	return true
}

// Close : 关闭组件
func (pprof *Pprof) Close() {
	if pprof.web != nil {
		pprof.web.Close()
		pprof.web = nil
	}
}

package utils

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/fananchong/go-xserver/common"
)

// Ticker : 定时器帮助类
type Ticker struct {
	ctx       *common.Context
	ctxSub    context.Context
	ctxCancel context.CancelFunc
	interval  time.Duration
	f         func()
	name      string
}

// NewTickerHelper : 定时器帮助类的构造函数
func NewTickerHelper(name string, ctx *common.Context, interval time.Duration, f func()) *Ticker {
	return &Ticker{
		ctx:      ctx,
		interval: interval,
		f:        f,
		name:     name,
	}
}

// Start : 开始
func (helper *Ticker) Start() bool {
	helper.ctxSub, helper.ctxCancel = context.WithCancel(helper.ctx)
	go helper.loop()
	return true
}

func (helper *Ticker) loop() {
	timer := time.NewTicker(helper.interval)
	defer timer.Stop()
	for {
		select {
		case <-helper.ctxSub.Done():
			helper.ctx.Infoln("Ticker[", helper.name, "], off")
			return
		case <-timer.C:
			func() {
				defer func() {
					if err := recover(); err != nil {
						helper.ctx.Errorln("Ticker[", helper.name, "] except:", err, "\n", string(debug.Stack()))
					}
				}()
				helper.f()
			}()
		}
	}
}

// Close : 结束
func (helper *Ticker) Close() {
	if helper.ctxCancel != nil {
		helper.ctxCancel()
		helper.ctxCancel = nil
	}
}

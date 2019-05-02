package components

import (
	"time"

	"github.com/fananchong/go-xserver/common"
)

// Time : 随机函数组件
type Time struct {
	ctx   *common.Context
	delta int64
}

// NewTime : 实例化
func NewTime(ctx *common.Context) *Time {
	t := &Time{ctx: ctx}
	t.init()
	return t
}

// Start : 实例化组件
func (t *Time) Start() bool {
	return true
}

func (t *Time) init() {
	t.ctx.ITime = t
}

// Close : 关闭组件
func (*Time) Close() {
	// No need to do anything
}

// GetTickCount : 毫秒数
func (t *Time) GetTickCount() int64 {
	now := time.Now().UnixNano() / 1e6
	return now + t.delta
}

// SetDelta : 设置时间差（单位毫秒），用来快进或后退当前时间（调试时间相关功能时，会用到）
func (t *Time) SetDelta(delta int64) {
	t.delta = delta
}

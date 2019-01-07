package components

import (
	"math/rand"
	"time"

	"github.com/fananchong/go-xserver/common"
)

// Rand : 随机函数组件
type Rand struct {
	ctx *common.Context
}

// NewRand : 实例化
func NewRand(ctx *common.Context) *Rand {
	return &Rand{ctx: ctx}
}

// Start : 实例化组件
func (r *Rand) Start() bool {
	r.ctx.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return true
}

// Close : 关闭组件
func (*Rand) Close() {
	// do nothing
}

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
	r := &Rand{ctx: ctx}
	r.init()
	return r
}

// Start : 实例化组件
func (r *Rand) Start() bool {
	// do nothing
	return true
}

func (r *Rand) init() {
	r.ctx.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Close : 关闭组件
func (*Rand) Close() {
	// do nothing
}

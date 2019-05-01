package components

import (
	"math/rand"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/utils"
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
	return true
}

func (r *Rand) init() {
	r.ctx.IRand = rand.New(rand.NewSource(utils.GetMillisecondTimestamp()))
}

// Close : 关闭组件
func (*Rand) Close() {
	// No need to do anything
}

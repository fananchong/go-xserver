package components

import (
	"math/rand"
	"time"

	"github.com/fananchong/go-xserver/common"
)

// Rand : 随机函数组件
type Rand struct {
}

// Start : 实例化组件
func (*Rand) Start() bool {
	common.XRAND = rand.New(rand.NewSource(time.Now().UnixNano()))
	return true
}

// Close : 关闭组件
func (*Rand) Close() {
	// do nothing
}

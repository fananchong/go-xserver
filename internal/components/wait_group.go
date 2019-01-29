package components

import (
	"context"
	"sync"
)

// ContextValueType : context value 的 key 类型
type ContextValueType int

const (
	// WAITGROUP : context value 的 key 类型
	WAITGROUP ContextValueType = iota
)

// CreateContext : 获取 context.Context 对象
func CreateContext() context.Context {
	return context.WithValue(context.Background(), WAITGROUP, &sync.WaitGroup{})
}

// SetComponentCount : 设置组件数量
func SetComponentCount(ctx context.Context, count int) {
	ctx.Value(WAITGROUP).(*sync.WaitGroup).Add(count)
}

// OneComponentOK : 某组件初始化完毕
func OneComponentOK(ctx context.Context) {
	ctx.Value(WAITGROUP).(*sync.WaitGroup).Done()
}

// WaitComponent : 等待所有组件初始化完毕
func WaitComponent(ctx context.Context) {
	ctx.Value(WAITGROUP).(*sync.WaitGroup).Wait()
}

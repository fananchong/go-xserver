package misc

import (
	"context"
	"sync"

	"github.com/fananchong/go-xserver/common/config"
)

// 框架层的一些全局变量，可以按以下方式 Set/Get

// ContextValueType : context value 的 key 类型
type ContextValueType int

const (
	_contextValueKey ContextValueType = iota
	_waitGroup
	_pluginType
)

// CreateContext : 获取 context.Context 对象
func CreateContext() context.Context {
	contextValue := make(map[ContextValueType]interface{})
	contextValue[_waitGroup] = &sync.WaitGroup{}
	contextValue[_pluginType] = 0
	return context.WithValue(context.Background(), _contextValueKey, contextValue)
}

// SetComponentCount : 设置组件数量
func SetComponentCount(ctx context.Context, count int) {
	values := ctx.Value(_contextValueKey).(map[ContextValueType]interface{})
	values[_waitGroup].(*sync.WaitGroup).Add(count)
}

// OneComponentOK : 某组件初始化完毕
func OneComponentOK(ctx context.Context) {
	values := ctx.Value(_contextValueKey).(map[ContextValueType]interface{})
	values[_waitGroup].(*sync.WaitGroup).Done()
}

// WaitComponent : 等待所有组件初始化完毕
func WaitComponent(ctx context.Context) {
	values := ctx.Value(_contextValueKey).(map[ContextValueType]interface{})
	values[_waitGroup].(*sync.WaitGroup).Wait()
}

// SetPluginType : 设置插件类型
func SetPluginType(ctx context.Context, t config.NodeType) {
	values := ctx.Value(_contextValueKey).(map[ContextValueType]interface{})
	values[_pluginType] = t
}

// GetPluginType : 获取插件类型
func GetPluginType(ctx context.Context) config.NodeType {
	values := ctx.Value(_contextValueKey).(map[ContextValueType]interface{})
	return values[_pluginType].(config.NodeType)
}

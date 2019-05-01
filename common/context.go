package common

import (
	gocontext "context"

	"github.com/fananchong/go-xserver/common/context"
	"github.com/fananchong/go-xserver/config"
)

// Context : 应用程序上下文
type Context struct {
	gocontext.Context     // golang context ，可以用于控制并发，传递全局变量等
	context.IRand         // 随机数生成器
	*config.Config        // 配置对象  // TODO, 待优化，支持逻辑层配置接入
	context.ILogger       // 日志对象
	context.INode         // 提供消息中继等功能
	context.IRole2Account // 提供`根据角色名查找账号`的功能；角色名重名检查也可以用该接口
	context.ITCPServer    // 提供对外的 TCP 服务。 注册该字段相应接口，才会开启
	context.ILogin        // 节点类型为 Login，才会开启
	context.IGateway      // 节点类型为 Gateway ，才会开启
}

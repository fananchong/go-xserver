package common

import (
	"context"
	"math/rand"
)

// Context : 应用程序上下文
type Context struct {
	context.Context // golang context ，可以用于控制并发，传递全局变量等
	*rand.Rand      // golang rand ，用来随机数
	*Config         // 配置对象
	ILogger         // 日志对象
	INode           // 提供消息中继等功能
	IRole2Account   // 提供`根据角色名查找账号`的功能；角色名重名检查也可以用该接口
	ITCPServer      // 提供对外的 TCP 服务。 注册该字段相应接口，才会开启
	ILogin          // 节点类型为 Login，才会开启
	IGateway        // 节点类型为 Gateway ，才会开启
}

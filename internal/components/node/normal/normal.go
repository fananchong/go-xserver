package nodenormal

import (
	"os"
	"sync"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/internal/components/misc"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	nodegateway "github.com/fananchong/go-xserver/internal/components/node/gateway"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utils"
)

// 通过该类接入服务器组，该类主要处理与 Mgr Server 的交互

// Normal : 普通节点
type Normal struct {
	*Session
	ctx        *common.Context
	components []utils.IComponent
	mtx        sync.Mutex
}

// NewNormal : 普通节点实现类的构造函数
func NewNormal(ctx *common.Context) *Normal {
	normal := &Normal{
		ctx:     ctx,
		Session: NewSession(ctx),
	}
	pluginType := misc.GetPluginType(ctx)
	if pluginType != config.Mgr {
		normal.Info = &protocol.SERVER_INFO{}
		if pluginType != config.Gateway {
			normal.Info.Id = nodecommon.NodeID2ServerID(nodecommon.NewNID(ctx, pluginType))
		}
		normal.Info.Type = uint32(pluginType)
		normal.Info.Addrs = []string{utils.GetIPInner(ctx), utils.GetIPOuter(ctx)}
		normal.Info.Ports = ctx.Config().Network.Port
		// TODO: 后续支持
		// normal.Info.Overload
		// normal.Info.Version
		normal.init()
		normal.ctx.INode = normal
		if pluginType != config.Gateway {
			ctx.Infoln("NODE ID:", normal.GetID(), ", NODE TYPE:", pluginType)
			misc.SetPluginID(normal.ctx, uint32(normal.GetID()))
		}
	}
	return normal
}

func (normal *Normal) init() bool {
	// ping ticker
	pingTicker := utils.NewTickerHelper("PING", normal.Ctx, 5*time.Second, normal.Ping)

	// bind components
	normal.components = []utils.IComponent{
		normal.Session,
		pingTicker,
	}
	return true
}

// Start : 节点开始工作
func (normal *Normal) Start() bool {
	pluginType := misc.GetPluginType(normal.ctx)
	if pluginType != config.Mgr {
		go func() {
			misc.WaitComponent(normal.ctx)
			if pluginType == config.Gateway {
				normal.Info.Id = nodecommon.NodeID2ServerID(normal.ctx.IGateway.(*nodegateway.Gateway).GetID())
				normal.ctx.Infoln("NODE ID:", normal.GetID(), ", NODE TYPE:", pluginType)
				misc.SetPluginID(normal.ctx, uint32(normal.GetID()))
			}
			normal.ctx.Infoln("Service node start ...")
			if normal.start() == false {
				normal.ctx.Errorln("Service node failed to start")
				os.Exit(1)
			}
		}()
	}
	return true
}

func (normal *Normal) start() bool {
	normal.mtx.Lock()
	defer normal.mtx.Unlock()
	for _, v := range normal.components {
		if v != nil && v.Start() == false {
			panic("")
		}
	}
	return true
}

// Close : 关闭节点
func (normal *Normal) Close() {
	normal.mtx.Lock()
	defer normal.mtx.Unlock()
	for _, v := range normal.components {
		v.Close()
	}
	normal.Session.Shutdown()
}

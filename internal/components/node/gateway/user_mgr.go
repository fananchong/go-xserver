package nodegateway

import (
	"errors"
	"sync"
	"time"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/db"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
)

// User : 表示 1 个客户端对象
type User struct {
	Account         string
	Servers         map[uint32]common.NodeID
	ActiveTimestamp int64
}

// NewUser : 客户端对象构造函数
func NewUser(ctx *common.Context, account string) *User {
	user := &User{
		Account:         account,
		Servers:         make(map[uint32]common.NodeID),
		ActiveTimestamp: utils.GetMillisecondTimestamp(),
	}
	return user
}

// UserMgr : 客户端对象管理类
type UserMgr struct {
	ctx               *common.Context
	users             map[string]*User
	mutex             sync.RWMutex
	ServerRedisCli    go_redis_orm.IClient
	checkActiveTicker *utils.Ticker
}

// NewUserMgr : 客户端对象管理类构造函数
func NewUserMgr(ctx *common.Context) *UserMgr {
	userMgr := &UserMgr{
		ctx:   ctx,
		users: make(map[string]*User),
	}
	userMgr.checkActiveTicker = utils.NewTickerHelper("CHECK_ACTIVE", ctx, 1*time.Second, userMgr.checkActive)
	return userMgr
}

// AddUser : 加入一个玩家
func (userMgr *UserMgr) AddUser(account string, servers map[uint32]*protocol.SERVER_ID) error {
	userMgr.mutex.Lock()
	defer userMgr.mutex.Unlock()
	user := NewUser(userMgr.ctx, account)
	for nodeType, serverID := range servers {
		user.Servers[nodeType] = utility.ServerID2NodeID(serverID)
		key := db.GetKeyAllocServer(nodeType, account)
		if _, err := userMgr.ServerRedisCli.Do("EXPIRE", key, 365*86400); err != nil { // 设置账号分配的服务器资源信息，过期时间 1 年
			return err
		}
	}
	userMgr.users[account] = user
	return nil
}

// GetServerAndActive : 获取玩家对应服务器类型的服务器资源信息
func (userMgr *UserMgr) GetServerAndActive(account string, nodeType common.NodeType) (*common.NodeID, error) {
	userMgr.mutex.RLock()
	defer userMgr.mutex.RUnlock()
	if user, ok := userMgr.users[account]; ok {
		user.ActiveTimestamp = utils.GetMillisecondTimestamp()
		if id, ok := user.Servers[uint32(nodeType)]; ok {
			return &id, nil
		}
		return nil, nil
	}
	return nil, errors.New("No server information corresponding to the account was found")
}

func (userMgr *UserMgr) checkActive() {
	userMgr.mutex.Lock()
	defer userMgr.mutex.Unlock()

	// 检查闲置连接
	now := utils.GetMillisecondTimestamp()
	var dels []*User
	for _, user := range userMgr.users {
		if now-user.ActiveTimestamp >= userMgr.ctx.Config.Role.IdleTime*1000 {
			dels = append(dels, user)
		}
	}

	// TODO: 删除操作现在是 1 条 1 条执行，会很慢，极端情况下，是卡玩家登录的。
	//       待优化为 REDIS PIPELINING 模式
	//       参考 ： https://godoc.org/github.com/gomodule/redigo/redis#hdr-Pipelining ，类似：
	//       	c.Send("SET", "foo", "bar")
	//			c.Send("GET", "foo")
	//			c.Flush()
	//			c.Receive() // reply from SET
	//			v, err = c.Receive() // reply from GET

	// 删除闲置连接
	for _, user := range dels {
		for nodeType, serverID := range user.Servers {
			_ = serverID
			key := db.GetKeyAllocServer(nodeType, user.Account)
			if _, err := userMgr.ServerRedisCli.Do("EXPIRE", key, 300); err != nil { // 设置账号分配的服务器资源信息，过期时间 5 分钟
				userMgr.ctx.Log.Errorln(err, "account:", user.Account)
			}
		}
		delete(userMgr.users, user.Account)
		userMgr.ctx.Log.Infoln("Delete user information, account:", user.Account)
	}
}

// Start : 开始
func (userMgr *UserMgr) Start() {
	userMgr.checkActiveTicker.Start()
}

// Close : 结束
func (userMgr *UserMgr) Close() {
	userMgr.checkActiveTicker.Close()
}

package nodegateway

import (
	"errors"
	"fmt"
	"sync"
	"time"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
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
		ActiveTimestamp: time.Now().UnixNano() / 1e6,
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
		key := fmt.Sprintf("srv%d_%s", nodeType, account)
		if _, err := userMgr.ServerRedisCli.Do("EXPIRE", key, 365*86400); err != nil {
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
		user.ActiveTimestamp = time.Now().UnixNano() / 1e6
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
	// TODO: 待写，检查是否是闲置连接。
	// for _, user := range userMgr.users {

	// }
}

// Start : 开始
func (userMgr *UserMgr) Start() {
	userMgr.checkActiveTicker.Start()
}

// Close : 结束
func (userMgr *UserMgr) Close() {
	userMgr.checkActiveTicker.Close()
}

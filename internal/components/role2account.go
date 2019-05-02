package components

import (
	"sync"
	"time"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/internal/db"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/utils"
)

// AccountInfo :
type AccountInfo struct {
	Role            string
	Account         string
	ActiveTimestamp int64
}

// NewAccountInfo :
func NewAccountInfo(role, account string, t int64) *AccountInfo {
	info := &AccountInfo{
		Role:            role,
		Account:         account,
		ActiveTimestamp: t,
	}
	return info
}

// Role2Account : 本帮助类，提供`根据角色名查找账号`的功能
type Role2Account struct {
	ctx               *common.Context
	cache             map[string]*AccountInfo
	mutex             sync.RWMutex
	checkActiveTicker *utils.Ticker
}

// NewRole2Account : 构造函数
func NewRole2Account(ctx *common.Context) *Role2Account {
	role2account := &Role2Account{
		ctx:   ctx,
		cache: make(map[string]*AccountInfo),
	}
	role2account.checkActiveTicker = utils.NewTickerHelper("CHECK_ACTIVE", ctx, 1*time.Second, role2account.checkActive)
	ctx.IRole2Account = role2account
	return role2account
}

// Add : 加入本地缓存
func (role2account *Role2Account) Add(role, account string) {
	role2account.mutex.Lock()
	defer role2account.mutex.Unlock()
	role2account.cache[role] = NewAccountInfo(role, account, role2account.ctx.GetTickCount())
}

// AddAndInsertDB : 加入本地缓存并保存数据库
func (role2account *Role2Account) AddAndInsertDB(role, account string) bool {
	dbObj := db.NewRoleName(role2account.ctx.Config.DbRoleName.Name, role)
	dbObj.SetAccount(account)
	if err := dbObj.Save(); err != nil {
		role2account.ctx.Errorln(err, "role:", role, "account:", account)
		return false
	}
	role2account.mutex.Lock()
	defer role2account.mutex.Unlock()
	role2account.cache[role] = NewAccountInfo(role, account, role2account.ctx.GetTickCount())
	return true
}

// GetAndActive : 根据角色名，查找账号。先从本地缓存中找，没有则数据库中找
func (role2account *Role2Account) GetAndActive(role string) string {
	role2account.mutex.RLock()
	if info, ok := role2account.cache[role]; ok {
		info.ActiveTimestamp = role2account.ctx.GetTickCount()
		role2account.mutex.RUnlock()
		return info.Account
	}
	role2account.mutex.RUnlock()
	dbObj := db.NewRoleName(role2account.ctx.Config.DbRoleName.Name, role)
	if err := dbObj.Load(); err != nil {
		if err != go_redis_orm.ERR_ISNOT_EXIST_KEY {
			role2account.ctx.Errorln(err, "role:", role)
		}
		return ""
	}
	role2account.Add(role, dbObj.GetAccount())
	return dbObj.GetAccount()
}

func (role2account *Role2Account) checkActive() {
	// TODO: 需要新增最小堆排序容器字段，并维护。使该函数计算复杂度最高为 O(logn)
	role2account.mutex.Lock()
	defer role2account.mutex.Unlock()
	now := role2account.ctx.GetTickCount()
	var dels []*AccountInfo
	for _, info := range role2account.cache {
		if now-info.ActiveTimestamp >= 24*60*60*1000 { // 1 天
			dels = append(dels, info)
		}
	}
	for _, info := range dels {
		delete(role2account.cache, info.Role)
	}
}

// Start : 开始
func (role2account *Role2Account) Start() bool {
	role2account.checkActiveTicker.Start()
	return true
}

// Close : 结束
func (role2account *Role2Account) Close() {
	role2account.checkActiveTicker.Close()
}

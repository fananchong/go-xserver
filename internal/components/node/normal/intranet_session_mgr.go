package nodenormal

import (
	"sync"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
)

// User :
type User struct {
	Account         string
	Sess            *IntranetSession
	ActiveTimestamp int64
}

// NewUser :
func NewUser(account string, sess *IntranetSession) *User {
	user := &User{
		Account:         account,
		Sess:            sess,
		ActiveTimestamp: utils.GetMillisecondTimestamp(),
	}
	return user
}

// IntranetSessionMgr : IntranetSession 对象管理类
type IntranetSessionMgr struct {
	ctx               *common.Context
	users             map[string]*User
	mutex             sync.RWMutex
	checkActiveTicker *utils.Ticker
}

// NewIntranetSessionMgr :
func NewIntranetSessionMgr(ctx *common.Context) *IntranetSessionMgr {
	mgr := &IntranetSessionMgr{
		ctx:   ctx,
		users: make(map[string]*User),
	}
	mgr.checkActiveTicker = utils.NewTickerHelper("CHECK_ACTIVE", ctx, 1*time.Second, mgr.checkActive)
	return mgr
}

// AddUser :
func (mgr *IntranetSessionMgr) AddUser(account string, sess *IntranetSession) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.users[account] = NewUser(account, sess)
}

// DelUser :
func (mgr *IntranetSessionMgr) DelUser(account string) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	delete(mgr.users, account)
}

// GetAndActive :
func (mgr *IntranetSessionMgr) GetAndActive(account string) *IntranetSession {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	if user, ok := mgr.users[account]; ok {
		user.ActiveTimestamp = utils.GetMillisecondTimestamp()
		return user.Sess
	}
	return nil
}

func (mgr *IntranetSessionMgr) checkActive() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	now := utils.GetMillisecondTimestamp()
	var dels []*User
	for _, user := range mgr.users {
		if now-user.ActiveTimestamp >= mgr.ctx.Config.Role.IdleTime*2*1000 {
			dels = append(dels, user)
		}
	}
	for _, user := range dels {
		delete(mgr.users, user.Account)
	}
}

// Foreach : 遍历
func (mgr *IntranetSessionMgr) Foreach(f func(user *User) bool) {
	mgr.mutex.RLock()
	var users []*User
	for _, user := range mgr.users {
		users = append(users, user)
	}
	mgr.mutex.RUnlock()
	for _, user := range users {
		if !f(user) {
			break
		}
	}
}

// Start : 开始
func (mgr *IntranetSessionMgr) Start() {
	mgr.checkActiveTicker.Start()
}

// Close : 结束
func (mgr *IntranetSessionMgr) Close() {
	mgr.checkActiveTicker.Close()
}

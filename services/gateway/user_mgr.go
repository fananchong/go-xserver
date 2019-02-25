package main

import "sync"

// UserMgr : 玩家管理类
type UserMgr struct {
	users map[string]*User
	mutex sync.RWMutex
}

// NewUserMgr : 玩家管理类构造函数
func NewUserMgr() *UserMgr {
	userMgr := &UserMgr{}
	userMgr.users = make(map[string]*User)
	return userMgr
}

// AddUser : 加入一个玩家
func (userMgr *UserMgr) AddUser(account string, user *User) {
	// 如果存在旧玩家对象，关闭之
	if old := userMgr.GetUser(account); old != nil {
		Ctx.Log.Infoln("Delete old player object, account:", account)
		old.Close()
	}
	userMgr.mutex.Lock()
	defer userMgr.mutex.Unlock()
	userMgr.users[account] = user
}

// GetUser : 获取一个玩家
func (userMgr *UserMgr) GetUser(account string) *User {
	userMgr.mutex.RLock()
	defer userMgr.mutex.RUnlock()
	if user, ok := userMgr.users[account]; ok {
		return user
	}
	return nil
}

// DelUser : 删除一个玩家
func (userMgr *UserMgr) DelUser(account string) {
	userMgr.mutex.Lock()
	defer userMgr.mutex.Unlock()
	if _, ok := userMgr.users[account]; ok {
		delete(userMgr.users, account)
	}
}

package nodecommon

import (
	"sync"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
)

var sessionMgrObj = newSessionMgr()

// SessionMgr : 网络会话对象管理类
type SessionMgr struct {
	m       sync.RWMutex
	ss      map[common.NodeType][]*SessionBase
	ssByID  map[common.NodeID]*SessionBase
	counter uint32
}

// GetSessionMgr : 获取网络会话对象管理类实例
func GetSessionMgr() *SessionMgr {
	return sessionMgrObj
}

func newSessionMgr() *SessionMgr {
	return &SessionMgr{
		ss:     make(map[common.NodeType][]*SessionBase),
		ssByID: make(map[common.NodeID]*SessionBase),
	}
}

// Register : 有网络会话节点加入
func (sessmgr *SessionMgr) Register(sess *SessionBase) {
	t := sess.GetType()
	sid := sess.GetSID()
	sessmgr.m.Lock()
	defer sessmgr.m.Unlock()
	for sessmgr.deleteSessInSS(sid, t) {
		// do nothing
	}
	sessmgr.ss[t] = append(sessmgr.ss[t], sess)
	sessmgr.ssByID[sess.GetID()] = sess
}

// Lose1 : 丢失网络会话节点
func (sessmgr *SessionMgr) Lose1(sess *SessionBase) {
	t := sess.GetType()
	sid := sess.GetSID()
	sessmgr.Lose2(sid, t)
}

// Lose2 : 丢失网络会话节点
func (sessmgr *SessionMgr) Lose2(sid *protocol.SERVER_ID, t common.NodeType) {
	sessmgr.m.Lock()
	defer sessmgr.m.Unlock()
	for sessmgr.deleteSessInSS(sid, t) {
		// do nothing
	}
	delete(sessmgr.ssByID, utility.ServerID2NodeID(sid))
}

func (sessmgr *SessionMgr) deleteSessInSS(sid *protocol.SERVER_ID, t common.NodeType) bool {
	// 不用加锁，调用它的函数会加锁
	if lst, ok := sessmgr.ss[t]; ok {
		findindex := -1
		for i, v := range lst {
			if utility.EqualSID(v.GetSID(), sid) {
				findindex = i
				break
			}
		}
		if findindex >= 0 {
			sessmgr.ss[t] = append(lst[:findindex], lst[findindex+1:]...)
			return true
		}
	}
	return false
}

// GetByID : 根据 NID 获取网络会话节点
func (sessmgr *SessionMgr) GetByID(nid common.NodeID) *SessionBase {
	sessmgr.m.RLock()
	defer sessmgr.m.RUnlock()
	if v, ok := sessmgr.ssByID[nid]; ok {
		return v
	}
	return nil
}

// GetByType : 根据节点类型，获取某类网络会话节点列表
func (sessmgr *SessionMgr) GetByType(t common.NodeType) []*SessionBase {
	sessmgr.m.RLock()
	defer sessmgr.m.RUnlock()
	var ret []*SessionBase
	if lst, ok := sessmgr.ss[t]; ok {
		ret = append(ret, lst...)
	}
	return ret
}

// GetAll : 获取所有网络会话节点列表
func (sessmgr *SessionMgr) GetAll() []*SessionBase {
	sessmgr.m.RLock()
	defer sessmgr.m.RUnlock()
	var ret []*SessionBase
	for _, lst := range sessmgr.ss {
		ret = append(ret, lst...)
	}
	return ret
}

// SelectOne : 选择 1 个某类型的网络会话节点
func (sessmgr *SessionMgr) SelectOne(t common.NodeType) *SessionBase {
	lst := GetSessionMgr().GetByType(t)
	n := len(lst)
	if n > 0 {
		index := int32(sessmgr.counter % uint32(n))
		sessmgr.counter++
		return lst[index]
	}
	return nil
}

// ForByType : 根据类型遍历网络会话节点
func (sessmgr *SessionMgr) ForByType(t common.NodeType, f func(*SessionBase)) {
	lst := sessmgr.GetByType(t)
	for _, v := range lst {
		f(v)
	}
}

// ForAll : 遍历所有网络会话节点
func (sessmgr *SessionMgr) ForAll(f func(*SessionBase)) {
	lst := sessmgr.GetAll()
	for _, v := range lst {
		f(v)
	}
}

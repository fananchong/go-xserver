package nodemgr

import (
	"sync"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
)

var (
	xsessionmgr = newSessionMgr()
)

// SessionMgr : 网络会话对象管理类
type SessionMgr struct {
	m       sync.RWMutex
	ss      map[common.NodeType][]*Session
	ssByID  map[common.NodeID]*Session
	counter uint32
}

func newSessionMgr() *SessionMgr {
	return &SessionMgr{
		ss:     make(map[common.NodeType][]*Session),
		ssByID: make(map[common.NodeID]*Session),
	}
}

func (sessmgr *SessionMgr) register(sess *Session) {
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

func (sessmgr *SessionMgr) lose(sess *Session) {
	t := sess.GetType()
	sid := sess.GetSID()
	sessmgr.m.Lock()
	defer sessmgr.m.Unlock()
	for sessmgr.deleteSessInSS(sid, t) {
		// do nothing
	}
	delete(sessmgr.ssByID, sess.GetID())
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

func (sessmgr *SessionMgr) getByID(nid common.NodeID) *Session {
	sessmgr.m.RLock()
	defer sessmgr.m.RUnlock()
	if v, ok := sessmgr.ssByID[nid]; ok {
		return v
	}
	return nil
}

func (sessmgr *SessionMgr) getByType(t common.NodeType) []*Session {
	sessmgr.m.RLock()
	defer sessmgr.m.RUnlock()
	var ret []*Session
	if lst, ok := sessmgr.ss[t]; ok {
		ret = append(ret, lst...)
	}
	return ret
}

func (sessmgr *SessionMgr) getAll() []*Session {
	sessmgr.m.RLock()
	defer sessmgr.m.RUnlock()
	var ret []*Session
	for _, lst := range sessmgr.ss {
		ret = append(ret, lst...)
	}
	return ret
}

func (sessmgr *SessionMgr) selectOne(t common.NodeType) *Session {
	lst := xsessionmgr.getByType(t)
	n := len(lst)
	if n > 0 {
		index := int32(sessmgr.counter % uint32(n))
		sessmgr.counter++
		return lst[index]
	}
	return nil
}

func (sessmgr *SessionMgr) forByType(t common.NodeType, f func(*Session)) {
	lst := sessmgr.getByType(t)
	for _, v := range lst {
		f(v)
	}
}

func (sessmgr *SessionMgr) forAll(f func(*Session)) {
	lst := sessmgr.getAll()
	for _, v := range lst {
		f(v)
	}
}

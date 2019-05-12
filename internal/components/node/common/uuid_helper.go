package nodecommon

import (
	"fmt"
	"os"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/common/context"
	"github.com/fananchong/go-xserver/internal/protocol"
)

// NewNID : 生成一个NID
func NewNID(ctx *common.Context, t config.NodeType) context.NodeID {
	key := fmt.Sprintf("NODEID%d", t)
	id, err := ctx.GetUID(key)
	if err != nil {
		ctx.Errorln(err)
		os.Exit(1)
	}
	return context.NodeID(uint32(t)*10000 + uint32(id)%10000)
}

// NodeID2ServerID : context.NodeID 转化为 protocol.SERVER_ID
func NodeID2ServerID(nid context.NodeID) *protocol.SERVER_ID {
	sid := &protocol.SERVER_ID{
		ID: uint32(nid),
	}
	return sid
}

// ServerID2NodeID : protocol.SERVER_ID 转化为 context.NodeID
func ServerID2NodeID(sid *protocol.SERVER_ID) context.NodeID {
	return context.NodeID(sid.ID)
}

// EqualSID : 2 个 SID 是否相等
func EqualSID(sid1 *protocol.SERVER_ID, sid2 *protocol.SERVER_ID) bool {
	return sid1.GetID() == sid2.GetID()
}

// EqualNID : 2 个 NID 是否相等
func EqualNID(nid1 context.NodeID, nid2 context.NodeID) bool {
	return uint32(nid1) == uint32(nid2)
}

package nodecommon

import (
	"bytes"
	"encoding/binary"

	"github.com/fananchong/go-xserver/internal/protocol"
	uuid "github.com/satori/go.uuid"
)

// NewNID : 生成一个NID
func NewNID() NodeID {
	return UUID2NodeID(uuid.NewV1())
}

// UUID2NodeID : uuid.UUID 转化为 NodeID
func UUID2NodeID(uid uuid.UUID) NodeID {
	nid := NodeID{}
	copy(nid[:], uid[:])
	return nid
}

// NodeID2UUID : NodeID 转化为 uuid.UUID
func NodeID2UUID(nid NodeID) uuid.UUID {
	uid := uuid.UUID{}
	copy(uid[:], nid[:])
	return uid
}

// NodeID2ServerID : NodeID 转化为 protocol.SERVER_ID
func NodeID2ServerID(nid NodeID) *protocol.SERVER_ID {
	sid := &protocol.SERVER_ID{}
	sid.Low = binary.LittleEndian.Uint64(nid[:8])
	sid.High = binary.LittleEndian.Uint64(nid[8:])
	return sid
}

// ServerID2NodeID : protocol.SERVER_ID 转化为 NodeID
func ServerID2NodeID(sid *protocol.SERVER_ID) NodeID {
	nid := NodeID{}
	binary.LittleEndian.PutUint64(nid[:8], sid.GetLow())
	binary.LittleEndian.PutUint64(nid[8:], sid.GetHigh())
	return nid
}

// ServerID2UUID : protocol.SERVER_ID 转化为 uuid.UUID
func ServerID2UUID(sid *protocol.SERVER_ID) uuid.UUID {
	uid := uuid.UUID{}
	binary.LittleEndian.PutUint64(uid[:8], sid.GetLow())
	binary.LittleEndian.PutUint64(uid[8:], sid.GetHigh())
	return uid
}

// EqualSID : 2 个 SID 是否相等
func EqualSID(sid1 *protocol.SERVER_ID, sid2 *protocol.SERVER_ID) bool {
	return sid1.GetLow() == sid2.GetLow() && sid1.GetHigh() == sid2.GetHigh()
}

// EqualNID : 2 个 NID 是否相等
func EqualNID(nid1 NodeID, nid2 NodeID) bool {
	return bytes.Equal(nid1[:], nid2[:])
}

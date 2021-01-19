// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lobby_custom.proto

/*
	Package protocol is a generated protocol buffer package.

	It is generated from these files:
		lobby_custom.proto
		lobby.proto
		match.proto

	It has these top-level messages:
		ROLE_BASE_INFO
		ROLE_DETAIL_INFO
		CMD_LOBBY
		ENUM_LOBBY_COMMON_ERROR
		MSG_LOBBY_LOGIN
		MSG_LOBBY_LOGIN_RESULT
		MSG_LOBBY_CREATE_ROLE
		MSG_LOBBY_CREATE_ROLE_RESULT
		MSG_LOBBY_ENTER_GAME
		MSG_LOBBY_ENTER_GAME_RESULT
		MSG_LOBBY_CHAT
		MSG_LOBBY_MATCH
		MSG_LOBBY_MATCH_RESULT
		CMD_MATCH
		ENUM_MATCH_COMMON_ERROR
		MSG_MATCH_MATCH
		MSG_MATCH_MATCH_RESULT
*/
package protocol

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// 角色基本信息
type ROLE_BASE_INFO struct {
	RoleID   uint64 `protobuf:"varint,1,opt,name=RoleID,proto3" json:"RoleID,omitempty"`
	RoleName string `protobuf:"bytes,2,opt,name=RoleName,proto3" json:"RoleName,omitempty"`
}

func (m *ROLE_BASE_INFO) Reset()                    { *m = ROLE_BASE_INFO{} }
func (m *ROLE_BASE_INFO) String() string            { return proto.CompactTextString(m) }
func (*ROLE_BASE_INFO) ProtoMessage()               {}
func (*ROLE_BASE_INFO) Descriptor() ([]byte, []int) { return fileDescriptorLobbyCustom, []int{0} }

func (m *ROLE_BASE_INFO) GetRoleID() uint64 {
	if m != nil {
		return m.RoleID
	}
	return 0
}

func (m *ROLE_BASE_INFO) GetRoleName() string {
	if m != nil {
		return m.RoleName
	}
	return ""
}

// 角色详细信息（不是角色所有数据！一般主界面上的角色数据就放这里，比如游戏币、体力等等。）
type ROLE_DETAIL_INFO struct {
	BaseInfo *ROLE_BASE_INFO `protobuf:"bytes,1,opt,name=BaseInfo" json:"BaseInfo,omitempty"`
}

func (m *ROLE_DETAIL_INFO) Reset()                    { *m = ROLE_DETAIL_INFO{} }
func (m *ROLE_DETAIL_INFO) String() string            { return proto.CompactTextString(m) }
func (*ROLE_DETAIL_INFO) ProtoMessage()               {}
func (*ROLE_DETAIL_INFO) Descriptor() ([]byte, []int) { return fileDescriptorLobbyCustom, []int{1} }

func (m *ROLE_DETAIL_INFO) GetBaseInfo() *ROLE_BASE_INFO {
	if m != nil {
		return m.BaseInfo
	}
	return nil
}

func init() {
	proto.RegisterType((*ROLE_BASE_INFO)(nil), "protocol.ROLE_BASE_INFO")
	proto.RegisterType((*ROLE_DETAIL_INFO)(nil), "protocol.ROLE_DETAIL_INFO")
}
func (m *ROLE_BASE_INFO) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ROLE_BASE_INFO) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RoleID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintLobbyCustom(dAtA, i, uint64(m.RoleID))
	}
	if len(m.RoleName) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintLobbyCustom(dAtA, i, uint64(len(m.RoleName)))
		i += copy(dAtA[i:], m.RoleName)
	}
	return i, nil
}

func (m *ROLE_DETAIL_INFO) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ROLE_DETAIL_INFO) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.BaseInfo != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintLobbyCustom(dAtA, i, uint64(m.BaseInfo.Size()))
		n1, err := m.BaseInfo.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func encodeFixed64LobbyCustom(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32LobbyCustom(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintLobbyCustom(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ROLE_BASE_INFO) Size() (n int) {
	var l int
	_ = l
	if m.RoleID != 0 {
		n += 1 + sovLobbyCustom(uint64(m.RoleID))
	}
	l = len(m.RoleName)
	if l > 0 {
		n += 1 + l + sovLobbyCustom(uint64(l))
	}
	return n
}

func (m *ROLE_DETAIL_INFO) Size() (n int) {
	var l int
	_ = l
	if m.BaseInfo != nil {
		l = m.BaseInfo.Size()
		n += 1 + l + sovLobbyCustom(uint64(l))
	}
	return n
}

func sovLobbyCustom(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozLobbyCustom(x uint64) (n int) {
	return sovLobbyCustom(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ROLE_BASE_INFO) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLobbyCustom
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ROLE_BASE_INFO: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ROLE_BASE_INFO: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RoleID", wireType)
			}
			m.RoleID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLobbyCustom
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RoleID |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RoleName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLobbyCustom
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthLobbyCustom
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RoleName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLobbyCustom(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthLobbyCustom
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ROLE_DETAIL_INFO) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLobbyCustom
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ROLE_DETAIL_INFO: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ROLE_DETAIL_INFO: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLobbyCustom
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLobbyCustom
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BaseInfo == nil {
				m.BaseInfo = &ROLE_BASE_INFO{}
			}
			if err := m.BaseInfo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLobbyCustom(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthLobbyCustom
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipLobbyCustom(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLobbyCustom
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLobbyCustom
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLobbyCustom
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthLobbyCustom
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowLobbyCustom
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipLobbyCustom(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthLobbyCustom = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLobbyCustom   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("lobby_custom.proto", fileDescriptorLobbyCustom) }

var fileDescriptorLobbyCustom = []byte{
	// 170 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xca, 0xc9, 0x4f, 0x4a,
	0xaa, 0x8c, 0x4f, 0x2e, 0x2d, 0x2e, 0xc9, 0xcf, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2,
	0x00, 0x53, 0xc9, 0xf9, 0x39, 0x4a, 0x2e, 0x5c, 0x7c, 0x41, 0xfe, 0x3e, 0xae, 0xf1, 0x4e, 0x8e,
	0xc1, 0xae, 0xf1, 0x9e, 0x7e, 0x6e, 0xfe, 0x42, 0x62, 0x5c, 0x6c, 0x41, 0xf9, 0x39, 0xa9, 0x9e,
	0x2e, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x50, 0x9e, 0x90, 0x14, 0x17, 0x07, 0x88, 0xe5,
	0x97, 0x98, 0x9b, 0x2a, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0xe7, 0x2b, 0x79, 0x70, 0x09,
	0x80, 0x4d, 0x71, 0x71, 0x0d, 0x71, 0xf4, 0xf4, 0x81, 0x98, 0x63, 0xc2, 0xc5, 0xe1, 0x94, 0x58,
	0x9c, 0xea, 0x99, 0x97, 0x96, 0x0f, 0x36, 0x89, 0xdb, 0x48, 0x42, 0x0f, 0x66, 0xad, 0x1e, 0xaa,
	0x9d, 0x41, 0x70, 0x95, 0x4e, 0x02, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0,
	0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x49, 0x6c, 0x60, 0x4d, 0xc6, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xc7, 0x53, 0xf8, 0x4e, 0xc8, 0x00, 0x00, 0x00,
}

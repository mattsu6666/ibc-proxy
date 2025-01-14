// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ibc/lightclients/proxy/v1/proxy.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/ibc-go/modules/core/02-client/types"
	types1 "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ClientState struct {
	// client state corresponding to proxy
	// the type must implements ClientState interface
	ProxyClientState *types.Any `protobuf:"bytes,1,opt,name=proxy_client_state,json=proxyClientState,proto3" json:"proxy_client_state,omitempty"`
	// client id corresponding to upstream on proxy
	UpstreamClientId string `protobuf:"bytes,2,opt,name=upstream_client_id,json=upstreamClientId,proto3" json:"upstream_client_id,omitempty"`
	// the proxy commitment prefix of the proxy chain.
	ProxyPrefix *types1.MerklePrefix `protobuf:"bytes,3,opt,name=proxy_prefix,json=proxyPrefix,proto3" json:"proxy_prefix,omitempty"`
	// the ibc commitment prefix of the proxy chain
	IbcPrefix *types1.MerklePrefix `protobuf:"bytes,4,opt,name=ibc_prefix,json=ibcPrefix,proto3" json:"ibc_prefix,omitempty"`
	// a boolean value indicating whether or not to trust to create ProxyClient by the relayer
	// if true, trust that the values in each field of the client state are valid
	// if false, the client state must be setup through a channel between downstream and proxy chain. please see ./modules/proxy/keeper/relay.go
	TrustedSetup bool `protobuf:"varint,5,opt,name=trusted_setup,json=trustedSetup,proto3" json:"trusted_setup,omitempty"`
}

func (m *ClientState) Reset()         { *m = ClientState{} }
func (m *ClientState) String() string { return proto.CompactTextString(m) }
func (*ClientState) ProtoMessage()    {}
func (*ClientState) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b548f5864814422, []int{0}
}
func (m *ClientState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClientState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClientState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClientState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientState.Merge(m, src)
}
func (m *ClientState) XXX_Size() int {
	return m.Size()
}
func (m *ClientState) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientState.DiscardUnknown(m)
}

var xxx_messageInfo_ClientState proto.InternalMessageInfo

type ConsensusState struct {
	// consensus state corresponding to proxy
	// the type must implements ConsensusState interface
	ProxyConsensusState *types.Any `protobuf:"bytes,1,opt,name=proxy_consensus_state,json=proxyConsensusState,proto3" json:"proxy_consensus_state,omitempty"`
}

func (m *ConsensusState) Reset()         { *m = ConsensusState{} }
func (m *ConsensusState) String() string { return proto.CompactTextString(m) }
func (*ConsensusState) ProtoMessage()    {}
func (*ConsensusState) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b548f5864814422, []int{1}
}
func (m *ConsensusState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConsensusState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ConsensusState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ConsensusState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusState.Merge(m, src)
}
func (m *ConsensusState) XXX_Size() int {
	return m.Size()
}
func (m *ConsensusState) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusState.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusState proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ClientState)(nil), "ibc.lightclients.proxy.v1.ClientState")
	proto.RegisterType((*ConsensusState)(nil), "ibc.lightclients.proxy.v1.ConsensusState")
}

func init() {
	proto.RegisterFile("ibc/lightclients/proxy/v1/proxy.proto", fileDescriptor_7b548f5864814422)
}

var fileDescriptor_7b548f5864814422 = []byte{
	// 415 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xcd, 0xce, 0xd2, 0x40,
	0x14, 0x6d, 0x3f, 0x3f, 0x8d, 0xdf, 0x80, 0x86, 0x54, 0x4c, 0x0a, 0x8b, 0x42, 0x50, 0x23, 0x0b,
	0x99, 0x09, 0xba, 0x73, 0x27, 0x24, 0xfe, 0x2c, 0x4c, 0x4c, 0xd9, 0xb9, 0x81, 0xce, 0x74, 0x28,
	0x13, 0xdb, 0x4e, 0xd3, 0x99, 0x92, 0xf2, 0x06, 0x2e, 0x7d, 0x04, 0x1f, 0xc3, 0x47, 0x70, 0xc9,
	0xd2, 0xa5, 0x81, 0x17, 0x31, 0x33, 0xb7, 0x55, 0xdc, 0xe9, 0xee, 0xfe, 0x9c, 0x7b, 0xce, 0xb9,
	0x37, 0x17, 0x3d, 0x11, 0x94, 0x91, 0x54, 0x24, 0x3b, 0xcd, 0x52, 0xc1, 0x73, 0xad, 0x48, 0x51,
	0xca, 0xfa, 0x40, 0xf6, 0x73, 0x08, 0x70, 0x51, 0x4a, 0x2d, 0xbd, 0x81, 0xa0, 0x0c, 0x5f, 0xc2,
	0x30, 0x74, 0xf7, 0xf3, 0x61, 0x3f, 0x91, 0x89, 0xb4, 0x28, 0x62, 0x22, 0x18, 0x18, 0x0e, 0x12,
	0x29, 0x93, 0x94, 0x13, 0x9b, 0xd1, 0x6a, 0x4b, 0xa2, 0xbc, 0xe1, 0x1a, 0x8e, 0x8c, 0x24, 0x93,
	0x25, 0x27, 0xc0, 0x65, 0xb4, 0x20, 0x6a, 0x00, 0x4f, 0xff, 0x00, 0x64, 0x96, 0x09, 0x9d, 0xb5,
	0xa0, 0xdf, 0x19, 0x00, 0x27, 0xdf, 0xae, 0x50, 0x67, 0x69, 0x27, 0x57, 0x3a, 0xd2, 0xdc, 0x5b,
	0x20, 0xcf, 0xda, 0x5a, 0x03, 0xdd, 0x5a, 0x99, 0xaa, 0xef, 0x8e, 0xdd, 0x69, 0xe7, 0x79, 0x1f,
	0x83, 0x23, 0xdc, 0x3a, 0xc2, 0xaf, 0xf2, 0x43, 0xd8, 0xb3, 0xf8, 0x4b, 0x8e, 0x67, 0xc8, 0xab,
	0x0a, 0xa5, 0x4b, 0x1e, 0x65, 0x2d, 0x8d, 0x88, 0xfd, 0xab, 0xb1, 0x3b, 0xbd, 0x09, 0x7b, 0x6d,
	0x07, 0x06, 0xde, 0xc5, 0xde, 0x1b, 0xd4, 0x05, 0xc5, 0xa2, 0xe4, 0x5b, 0x51, 0xfb, 0xb7, 0xac,
	0xd6, 0x63, 0x6c, 0xce, 0x65, 0x36, 0xc0, 0x17, 0x9e, 0xf7, 0x73, 0xfc, 0x9e, 0x97, 0x9f, 0x52,
	0xfe, 0xc1, 0x62, 0xc3, 0x8e, 0x9d, 0x84, 0xc4, 0x5b, 0x22, 0x24, 0x28, 0x6b, 0x69, 0xae, 0xff,
	0x83, 0xe6, 0x46, 0x50, 0xd6, 0x90, 0x3c, 0x42, 0xf7, 0x74, 0x59, 0x29, 0xcd, 0xe3, 0xb5, 0xe2,
	0xba, 0x2a, 0xfc, 0xdb, 0x63, 0x77, 0x7a, 0x37, 0xec, 0x36, 0xc5, 0x95, 0xa9, 0xbd, 0xbc, 0xfe,
	0xfc, 0x75, 0xe4, 0x4c, 0x36, 0xe8, 0xfe, 0x52, 0xe6, 0x8a, 0xe7, 0xaa, 0x52, 0xb0, 0xf8, 0x5b,
	0xf4, 0xb0, 0x39, 0x5e, 0x5b, 0xff, 0x87, 0xfb, 0x3d, 0x80, 0xfb, 0xfd, 0xc5, 0x04, 0x0a, 0x8b,
	0xcd, 0xf7, 0x53, 0xe0, 0x1e, 0x4f, 0x81, 0xfb, 0xf3, 0x14, 0xb8, 0x5f, 0xce, 0x81, 0x73, 0x3c,
	0x07, 0xce, 0x8f, 0x73, 0xe0, 0x7c, 0x7c, 0x9d, 0x08, 0xbd, 0xab, 0xa8, 0x59, 0x8a, 0xc4, 0x91,
	0x8e, 0xd8, 0x2e, 0x12, 0x79, 0x1a, 0x51, 0x22, 0x28, 0x9b, 0xc1, 0xfb, 0x65, 0x32, 0xae, 0x52,
	0xae, 0xe0, 0x33, 0x67, 0xed, 0x6b, 0xd6, 0x75, 0xd3, 0xd6, 0x87, 0x82, 0x2b, 0x7a, 0xc7, 0x5a,
	0x79, 0xf1, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xc2, 0xc7, 0x41, 0x71, 0xc4, 0x02, 0x00, 0x00,
}

func (m *ClientState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClientState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClientState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TrustedSetup {
		i--
		if m.TrustedSetup {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if m.IbcPrefix != nil {
		{
			size, err := m.IbcPrefix.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProxy(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.ProxyPrefix != nil {
		{
			size, err := m.ProxyPrefix.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProxy(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.UpstreamClientId) > 0 {
		i -= len(m.UpstreamClientId)
		copy(dAtA[i:], m.UpstreamClientId)
		i = encodeVarintProxy(dAtA, i, uint64(len(m.UpstreamClientId)))
		i--
		dAtA[i] = 0x12
	}
	if m.ProxyClientState != nil {
		{
			size, err := m.ProxyClientState.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProxy(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ConsensusState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConsensusState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ConsensusState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ProxyConsensusState != nil {
		{
			size, err := m.ProxyConsensusState.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProxy(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProxy(dAtA []byte, offset int, v uint64) int {
	offset -= sovProxy(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ClientState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ProxyClientState != nil {
		l = m.ProxyClientState.Size()
		n += 1 + l + sovProxy(uint64(l))
	}
	l = len(m.UpstreamClientId)
	if l > 0 {
		n += 1 + l + sovProxy(uint64(l))
	}
	if m.ProxyPrefix != nil {
		l = m.ProxyPrefix.Size()
		n += 1 + l + sovProxy(uint64(l))
	}
	if m.IbcPrefix != nil {
		l = m.IbcPrefix.Size()
		n += 1 + l + sovProxy(uint64(l))
	}
	if m.TrustedSetup {
		n += 2
	}
	return n
}

func (m *ConsensusState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ProxyConsensusState != nil {
		l = m.ProxyConsensusState.Size()
		n += 1 + l + sovProxy(uint64(l))
	}
	return n
}

func sovProxy(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProxy(x uint64) (n int) {
	return sovProxy(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClientState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProxy
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClientState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClientState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProxyClientState", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProxy
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProxy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ProxyClientState == nil {
				m.ProxyClientState = &types.Any{}
			}
			if err := m.ProxyClientState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpstreamClientId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProxy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProxy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UpstreamClientId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProxyPrefix", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProxy
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProxy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ProxyPrefix == nil {
				m.ProxyPrefix = &types1.MerklePrefix{}
			}
			if err := m.ProxyPrefix.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcPrefix", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProxy
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProxy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.IbcPrefix == nil {
				m.IbcPrefix = &types1.MerklePrefix{}
			}
			if err := m.IbcPrefix.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TrustedSetup", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.TrustedSetup = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipProxy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProxy
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
func (m *ConsensusState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProxy
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ConsensusState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ConsensusState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProxyConsensusState", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProxy
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProxy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ProxyConsensusState == nil {
				m.ProxyConsensusState = &types.Any{}
			}
			if err := m.ProxyConsensusState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProxy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProxy
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
func skipProxy(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProxy
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
					return 0, ErrIntOverflowProxy
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProxy
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
			if length < 0 {
				return 0, ErrInvalidLengthProxy
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProxy
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProxy
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProxy        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProxy          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProxy = fmt.Errorf("proto: unexpected end of group")
)

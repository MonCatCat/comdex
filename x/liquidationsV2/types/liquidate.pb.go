// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: comdex/liquidationsV2/v1beta1/liquidate.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type WhiteListing struct {
	AppId       uint64 `protobuf:"varint,1,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty" yaml:"app_id"`
	AuctionType uint64 `protobuf:"varint,2,opt,name=auction_type,json=auctionType,proto3" json:"auction_type,omitempty" yaml:"auction_type"`
}

func (m *WhiteListing) Reset()         { *m = WhiteListing{} }
func (m *WhiteListing) String() string { return proto.CompactTextString(m) }
func (*WhiteListing) ProtoMessage()    {}
func (*WhiteListing) Descriptor() ([]byte, []int) {
	return fileDescriptor_631048b9d11253bf, []int{0}
}
func (m *WhiteListing) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WhiteListing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WhiteListing.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WhiteListing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WhiteListing.Merge(m, src)
}
func (m *WhiteListing) XXX_Size() int {
	return m.Size()
}
func (m *WhiteListing) XXX_DiscardUnknown() {
	xxx_messageInfo_WhiteListing.DiscardUnknown(m)
}

var xxx_messageInfo_WhiteListing proto.InternalMessageInfo

type LockedVault struct {
	LockedVaultId                uint64                                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" yaml:"id"`
	AppId                        uint64                                 `protobuf:"varint,2,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty" yaml:"id"`
	OriginalVaultId              uint64                                 `protobuf:"varint,3,opt,name=original_vault_id,json=originalVaultId,proto3" json:"original_vault_id,omitempty" yaml:"id"`
	ExtendedPairId               uint64                                 `protobuf:"varint,4,opt,name=extended_pair_vault_id,json=extendedPairVaultId,proto3" json:"extended_pair_vault_id,omitempty" yaml:"extended_pair_vault_id"`
	Owner                        string                                 `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty" yaml:"owner"`
	AmountIn                     github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=amount_in,json=amountIn,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount_in" yaml:"amount_in"`
	AmountOut                    github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,7,opt,name=amount_out,json=amountOut,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount_out" yaml:"amount_out"`
	CurrentCollaterlisationRatio github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=current_collateralisation_ratio,json=currentCollateralisationRatio,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"current_collateralisation_ratio" yaml:"current_collateralisation_ratio"`
	CollateralToBeAuctioned      github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,9,opt,name=collateral_to_be_auctioned,json=collateralToBeAuctioned,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"collateral_to_be_auctioned" yaml:"collateral_to_be_auctioned"`
	LiquidationTimestamp         time.Time                              `protobuf:"bytes,10,opt,name=liquidation_timestamp,json=liquidationTimestamp,proto3,stdtime" json:"liquidation_timestamp" yaml:"liquidation_timestamp"`
	IsInternalKeeper             bool                                   `protobuf:"varint,11,opt,name=is_internal_keeper,json=isInternalKeeper,proto3" json:"is_internal_keeper,omitempty" yaml:"is_intenal_keeper"`
	InternalKeeperAddress        string                                 `protobuf:"bytes,12,opt,name=internal_keeper_address,json=internalKeeperAddress,proto3" json:"internal_keeper_address,omitempty" yaml:"internal_keeper_address"`
	IsExternalKeeper             string                                 `protobuf:"bytes,13,opt,name=is_external_keeper,json=isExternalKeeper,proto3" json:"is_external_keeper,omitempty" yaml:"is_external_keeper"`
	ExternalKeeperAddress        string                                 `protobuf:"bytes,14,opt,name=external_keeper_address,json=externalKeeperAddress,proto3" json:"external_keeper_address,omitempty" yaml:"external_keeper_address"`
}

func (m *LockedVault) Reset()         { *m = LockedVault{} }
func (m *LockedVault) String() string { return proto.CompactTextString(m) }
func (*LockedVault) ProtoMessage()    {}
func (*LockedVault) Descriptor() ([]byte, []int) {
	return fileDescriptor_631048b9d11253bf, []int{1}
}
func (m *LockedVault) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LockedVault) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LockedVault.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LockedVault) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LockedVault.Merge(m, src)
}
func (m *LockedVault) XXX_Size() int {
	return m.Size()
}
func (m *LockedVault) XXX_DiscardUnknown() {
	xxx_messageInfo_LockedVault.DiscardUnknown(m)
}

var xxx_messageInfo_LockedVault proto.InternalMessageInfo

func init() {
	proto.RegisterType((*WhiteListing)(nil), "comdex.newliq.v1beta1.WhiteListing")
	proto.RegisterType((*LockedVault)(nil), "comdex.newliq.v1beta1.LockedVault")
}

func init() {
	proto.RegisterFile("comdex/liquidationsV2/v1beta1/liquidate.proto", fileDescriptor_631048b9d11253bf)
}

var fileDescriptor_631048b9d11253bf = []byte{
	// 856 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xcd, 0x6e, 0xdb, 0x36,
	0x1c, 0xb7, 0xb2, 0x26, 0x8b, 0xe9, 0xa4, 0x4d, 0xd5, 0xb8, 0x51, 0x8d, 0x56, 0xf4, 0x74, 0x08,
	0x72, 0x89, 0xb4, 0xb4, 0x97, 0x61, 0xbb, 0xcc, 0xee, 0x7a, 0x30, 0x16, 0x2c, 0x85, 0x10, 0x74,
	0x40, 0x0f, 0x13, 0x68, 0x89, 0x71, 0x88, 0xc8, 0xa2, 0x2a, 0x51, 0x6d, 0xf2, 0x0e, 0x3b, 0x74,
	0x6f, 0x31, 0x60, 0x2f, 0xb1, 0x63, 0x8e, 0x3d, 0x0e, 0x3b, 0x70, 0x9b, 0xf2, 0x06, 0x7a, 0x82,
	0x81, 0x1f, 0xfe, 0x90, 0xab, 0xa2, 0xc8, 0x25, 0x8e, 0xf8, 0xff, 0x7d, 0x42, 0x24, 0x05, 0x0e,
	0x43, 0x3a, 0x8d, 0xf0, 0xa5, 0x17, 0x93, 0x37, 0x05, 0x89, 0x10, 0x23, 0x34, 0xc9, 0x5f, 0x3d,
	0xf5, 0xde, 0x1e, 0x8d, 0x31, 0x43, 0x47, 0xf3, 0x65, 0xec, 0xa6, 0x19, 0x65, 0xd4, 0xec, 0x2a,
	0xb8, 0x9b, 0xe0, 0x77, 0x31, 0x79, 0xe3, 0x6a, 0x58, 0x6f, 0x77, 0x42, 0x27, 0x54, 0x22, 0x3c,
	0xf1, 0x9f, 0x02, 0xf7, 0xe0, 0x84, 0xd2, 0x49, 0x8c, 0x3d, 0xf9, 0x34, 0x2e, 0xce, 0x3c, 0x46,
	0xa6, 0x38, 0x67, 0x68, 0x9a, 0x6a, 0x80, 0x1d, 0xd2, 0x7c, 0x4a, 0x73, 0x6f, 0x8c, 0x72, 0x3c,
	0xb7, 0x0c, 0x29, 0x49, 0xd4, 0xdc, 0xf9, 0xd5, 0x00, 0x5b, 0x3f, 0x9f, 0x13, 0x86, 0x8f, 0x49,
	0xce, 0x48, 0x32, 0x31, 0x8f, 0xc0, 0x06, 0x4a, 0xd3, 0x80, 0x44, 0x96, 0xd1, 0x37, 0x0e, 0xee,
	0x0c, 0x7b, 0x25, 0x87, 0xeb, 0x83, 0x34, 0x1d, 0x45, 0x15, 0x87, 0xdb, 0x57, 0x68, 0x1a, 0x7f,
	0xeb, 0x28, 0x80, 0xe3, 0xaf, 0x23, 0xb1, 0x6e, 0x8e, 0xc0, 0x16, 0x2a, 0x42, 0xd1, 0x2b, 0x60,
	0x57, 0x29, 0xb6, 0xd6, 0x24, 0x71, 0xbf, 0xe4, 0xb0, 0x33, 0x50, 0xeb, 0xa7, 0x57, 0x29, 0xae,
	0x38, 0x7c, 0xa0, 0xe9, 0x4b, 0x60, 0xc7, 0xef, 0xa0, 0x05, 0xc6, 0xb9, 0xe9, 0x80, 0xce, 0x31,
	0x0d, 0x2f, 0x70, 0xf4, 0x0a, 0x15, 0x31, 0x33, 0x5d, 0xb0, 0x36, 0x4f, 0x62, 0x97, 0x1c, 0x6e,
	0x2f, 0x0d, 0x65, 0xa2, 0xb6, 0x92, 0x14, 0x69, 0xd6, 0x48, 0x64, 0x1e, 0xce, 0xd3, 0xab, 0x10,
	0x0f, 0x97, 0xd3, 0x2f, 0x61, 0x75, 0xf2, 0x63, 0x70, 0x9f, 0x66, 0x64, 0x42, 0x12, 0x14, 0x07,
	0x6f, 0x85, 0xa6, 0x60, 0x7e, 0x21, 0x99, 0xfd, 0x92, 0xc3, 0x7b, 0x27, 0x7a, 0xd8, 0xe8, 0x77,
	0x8f, 0xd6, 0xa7, 0xe6, 0x39, 0x78, 0x88, 0x2f, 0x19, 0x4e, 0x22, 0x1c, 0x05, 0x29, 0x22, 0xd9,
	0x42, 0xf2, 0x8e, 0x94, 0x7c, 0x56, 0x72, 0x78, 0xf7, 0x85, 0x46, 0xbc, 0x44, 0x24, 0x93, 0x8a,
	0x4f, 0x94, 0x62, 0x33, 0xd3, 0xf1, 0x1f, 0xe0, 0x25, 0xc2, 0xcc, 0xc9, 0x03, 0xeb, 0xf4, 0x5d,
	0x82, 0x33, 0x6b, 0xbd, 0x6f, 0x1c, 0xb4, 0x87, 0x8f, 0x44, 0xcb, 0x13, 0xb1, 0x50, 0x71, 0xb8,
	0xa5, 0xf4, 0xe4, 0xdc, 0xf1, 0x15, 0xce, 0xbc, 0x00, 0x6d, 0x34, 0xa5, 0x45, 0xc2, 0x02, 0x92,
	0x58, 0x1b, 0x92, 0xf4, 0xd3, 0x35, 0x87, 0xad, 0xbf, 0x39, 0xdc, 0x9f, 0x10, 0x76, 0x5e, 0x8c,
	0xdd, 0x90, 0x4e, 0x3d, 0xbd, 0x59, 0xd4, 0xcf, 0x61, 0x1e, 0x5d, 0x78, 0xe2, 0x1d, 0xe5, 0xee,
	0x28, 0x61, 0x25, 0x87, 0x9b, 0x03, 0x29, 0x31, 0x4a, 0x2a, 0x0e, 0x77, 0xf4, 0xab, 0x9c, 0x89,
	0x3a, 0xfe, 0x26, 0xd2, 0x53, 0x93, 0x02, 0xa0, 0xd7, 0x69, 0xc1, 0xac, 0x2f, 0xa5, 0xdb, 0xcb,
	0x5b, 0xbb, 0xb5, 0x95, 0xdb, 0x49, 0xc1, 0x2a, 0x0e, 0xef, 0xd7, 0xec, 0x68, 0xc1, 0x1c, 0x5f,
	0x17, 0x3a, 0x29, 0x98, 0xf9, 0xa7, 0x01, 0x60, 0x58, 0x64, 0x19, 0x4e, 0x58, 0x10, 0xd2, 0x38,
	0x46, 0x0c, 0x67, 0x28, 0x26, 0xb9, 0x3c, 0x6b, 0x41, 0x26, 0x7e, 0xac, 0x4d, 0x19, 0xe3, 0xf2,
	0x16, 0x31, 0x7e, 0xc0, 0x61, 0xc9, 0xe1, 0xe3, 0xe7, 0x4a, 0xf8, 0xb9, 0xd6, 0x9d, 0xc9, 0xfa,
	0xe2, 0x6f, 0xc5, 0xe1, 0xbe, 0x4a, 0xf6, 0x19, 0x7b, 0xc7, 0x7f, 0x12, 0xd6, 0x75, 0x50, 0x4d,
	0xc8, 0xfc, 0xc3, 0x00, 0xbd, 0x05, 0x37, 0x60, 0x34, 0x18, 0xe3, 0x40, 0x9f, 0x0c, 0x1c, 0x59,
	0x6d, 0x99, 0x3e, 0xb9, 0x75, 0xfa, 0xbd, 0x85, 0xdd, 0x29, 0x1d, 0xe2, 0xc1, 0x4c, 0xb0, 0xe2,
	0xf0, 0x2b, 0x1d, 0xfc, 0x93, 0xa6, 0x8e, 0xbf, 0x17, 0x36, 0xb3, 0xcd, 0xdf, 0x0c, 0xd0, 0x5d,
	0xba, 0xce, 0x82, 0xf9, 0xad, 0x63, 0x81, 0xbe, 0x71, 0xd0, 0x79, 0xda, 0x73, 0xd5, 0xbd, 0xe4,
	0xce, 0xee, 0x25, 0xf7, 0x74, 0x86, 0x18, 0x7e, 0x2f, 0x4a, 0x94, 0x1c, 0xee, 0x1e, 0x2f, 0x04,
	0xe6, 0xd3, 0x8a, 0xc3, 0xc7, 0x2a, 0x57, 0xa3, 0xbc, 0xf3, 0xfe, 0x1f, 0x68, 0xf8, 0xbb, 0x71,
	0x03, 0xd3, 0xfc, 0x05, 0x98, 0x24, 0x0f, 0x48, 0xc2, 0x70, 0x26, 0x8e, 0xf3, 0x05, 0xc6, 0x29,
	0xce, 0xac, 0x4e, 0xdf, 0x38, 0xd8, 0x1c, 0x7e, 0x5d, 0x72, 0xb8, 0x33, 0xca, 0x47, 0x7a, 0xf8,
	0xa3, 0x9c, 0x55, 0x1c, 0x5a, 0xfa, 0x34, 0x2b, 0xde, 0x82, 0xe6, 0xf8, 0x3b, 0x64, 0x05, 0x6d,
	0xe6, 0x60, 0x6f, 0x45, 0x3c, 0x40, 0x51, 0x94, 0xe1, 0x3c, 0xb7, 0xb6, 0xe4, 0xdb, 0xf9, 0xae,
	0xe4, 0xb0, 0x5b, 0x27, 0x0d, 0x14, 0xa0, 0xe2, 0xd0, 0xd6, 0x4e, 0xcd, 0x0a, 0x8e, 0xdf, 0x25,
	0x4d, 0x44, 0x33, 0x90, 0xa5, 0xc4, 0x15, 0xb0, 0x5c, 0x6a, 0x5b, 0xfa, 0x1d, 0xa9, 0x52, 0x2f,
	0x2e, 0x57, 0x4a, 0x3d, 0x9a, 0x97, 0x5a, 0xe1, 0xc9, 0x56, 0x75, 0xb8, 0x68, 0xb5, 0x82, 0x9a,
	0xb7, 0xba, 0xbb, 0x68, 0x55, 0x27, 0x7d, 0xd4, 0xea, 0x13, 0x0a, 0x8e, 0xdf, 0xc5, 0x4d, 0xc4,
	0xe1, 0xeb, 0xeb, 0xff, 0xec, 0xd6, 0xef, 0xa5, 0xdd, 0xba, 0x2e, 0x6d, 0xe3, 0x43, 0x69, 0x1b,
	0xff, 0x96, 0xb6, 0xf1, 0xfe, 0xc6, 0x6e, 0x7d, 0xb8, 0xb1, 0x5b, 0x7f, 0xdd, 0xd8, 0xad, 0xd7,
	0xdf, 0xd4, 0x76, 0xb8, 0xf8, 0x1e, 0x1e, 0xd2, 0xb3, 0x33, 0x12, 0x12, 0x14, 0xeb, 0x67, 0xef,
	0xa3, 0x0f, 0xaa, 0xdc, 0xf7, 0xe3, 0x0d, 0xb9, 0xe5, 0x9e, 0xfd, 0x1f, 0x00, 0x00, 0xff, 0xff,
	0x99, 0x33, 0xd7, 0x9e, 0x76, 0x07, 0x00, 0x00,
}

func (m *WhiteListing) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WhiteListing) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WhiteListing) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AuctionType != 0 {
		i = encodeVarintLiquidate(dAtA, i, uint64(m.AuctionType))
		i--
		dAtA[i] = 0x10
	}
	if m.AppId != 0 {
		i = encodeVarintLiquidate(dAtA, i, uint64(m.AppId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LockedVault) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LockedVault) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LockedVault) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ExternalKeeperAddress) > 0 {
		i -= len(m.ExternalKeeperAddress)
		copy(dAtA[i:], m.ExternalKeeperAddress)
		i = encodeVarintLiquidate(dAtA, i, uint64(len(m.ExternalKeeperAddress)))
		i--
		dAtA[i] = 0x72
	}
	if len(m.IsExternalKeeper) > 0 {
		i -= len(m.IsExternalKeeper)
		copy(dAtA[i:], m.IsExternalKeeper)
		i = encodeVarintLiquidate(dAtA, i, uint64(len(m.IsExternalKeeper)))
		i--
		dAtA[i] = 0x6a
	}
	if len(m.InternalKeeperAddress) > 0 {
		i -= len(m.InternalKeeperAddress)
		copy(dAtA[i:], m.InternalKeeperAddress)
		i = encodeVarintLiquidate(dAtA, i, uint64(len(m.InternalKeeperAddress)))
		i--
		dAtA[i] = 0x62
	}
	if m.IsInternalKeeper {
		i--
		if m.IsInternalKeeper {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x58
	}
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.LiquidationTimestamp, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.LiquidationTimestamp):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintLiquidate(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x52
	{
		size := m.CollateralToBeAuctioned.Size()
		i -= size
		if _, err := m.CollateralToBeAuctioned.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLiquidate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.CurrentCollaterlisationRatio.Size()
		i -= size
		if _, err := m.CurrentCollaterlisationRatio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLiquidate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.AmountOut.Size()
		i -= size
		if _, err := m.AmountOut.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLiquidate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.AmountIn.Size()
		i -= size
		if _, err := m.AmountIn.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLiquidate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintLiquidate(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x2a
	}
	if m.ExtendedPairId != 0 {
		i = encodeVarintLiquidate(dAtA, i, uint64(m.ExtendedPairId))
		i--
		dAtA[i] = 0x20
	}
	if m.OriginalVaultId != 0 {
		i = encodeVarintLiquidate(dAtA, i, uint64(m.OriginalVaultId))
		i--
		dAtA[i] = 0x18
	}
	if m.AppId != 0 {
		i = encodeVarintLiquidate(dAtA, i, uint64(m.AppId))
		i--
		dAtA[i] = 0x10
	}
	if m.LockedVaultId != 0 {
		i = encodeVarintLiquidate(dAtA, i, uint64(m.LockedVaultId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintLiquidate(dAtA []byte, offset int, v uint64) int {
	offset -= sovLiquidate(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *WhiteListing) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AppId != 0 {
		n += 1 + sovLiquidate(uint64(m.AppId))
	}
	if m.AuctionType != 0 {
		n += 1 + sovLiquidate(uint64(m.AuctionType))
	}
	return n
}

func (m *LockedVault) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.LockedVaultId != 0 {
		n += 1 + sovLiquidate(uint64(m.LockedVaultId))
	}
	if m.AppId != 0 {
		n += 1 + sovLiquidate(uint64(m.AppId))
	}
	if m.OriginalVaultId != 0 {
		n += 1 + sovLiquidate(uint64(m.OriginalVaultId))
	}
	if m.ExtendedPairId != 0 {
		n += 1 + sovLiquidate(uint64(m.ExtendedPairId))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovLiquidate(uint64(l))
	}
	l = m.AmountIn.Size()
	n += 1 + l + sovLiquidate(uint64(l))
	l = m.AmountOut.Size()
	n += 1 + l + sovLiquidate(uint64(l))
	l = m.CurrentCollaterlisationRatio.Size()
	n += 1 + l + sovLiquidate(uint64(l))
	l = m.CollateralToBeAuctioned.Size()
	n += 1 + l + sovLiquidate(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.LiquidationTimestamp)
	n += 1 + l + sovLiquidate(uint64(l))
	if m.IsInternalKeeper {
		n += 2
	}
	l = len(m.InternalKeeperAddress)
	if l > 0 {
		n += 1 + l + sovLiquidate(uint64(l))
	}
	l = len(m.IsExternalKeeper)
	if l > 0 {
		n += 1 + l + sovLiquidate(uint64(l))
	}
	l = len(m.ExternalKeeperAddress)
	if l > 0 {
		n += 1 + l + sovLiquidate(uint64(l))
	}
	return n
}

func sovLiquidate(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLiquidate(x uint64) (n int) {
	return sovLiquidate(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *WhiteListing) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquidate
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
			return fmt.Errorf("proto: WhiteListing: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WhiteListing: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AppId", wireType)
			}
			m.AppId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AppId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionType", wireType)
			}
			m.AuctionType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuctionType |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipLiquidate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquidate
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
func (m *LockedVault) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquidate
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
			return fmt.Errorf("proto: LockedVault: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LockedVault: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockedVaultId", wireType)
			}
			m.LockedVaultId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LockedVaultId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AppId", wireType)
			}
			m.AppId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AppId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginalVaultId", wireType)
			}
			m.OriginalVaultId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OriginalVaultId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtendedPairId", wireType)
			}
			m.ExtendedPairId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExtendedPairId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountIn", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountIn.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountOut", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentCollaterlisationRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CurrentCollaterlisationRatio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralToBeAuctioned", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralToBeAuctioned.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidationTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.LiquidationTimestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsInternalKeeper", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
			m.IsInternalKeeper = bool(v != 0)
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InternalKeeperAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InternalKeeperAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsExternalKeeper", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IsExternalKeeper = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExternalKeeperAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidate
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
				return ErrInvalidLengthLiquidate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExternalKeeperAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLiquidate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquidate
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
func skipLiquidate(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLiquidate
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
					return 0, ErrIntOverflowLiquidate
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
					return 0, ErrIntOverflowLiquidate
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
				return 0, ErrInvalidLengthLiquidate
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLiquidate
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLiquidate
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLiquidate        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLiquidate          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLiquidate = fmt.Errorf("proto: unexpected end of group")
)

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: friends.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AddFriendRequest struct {
	UserId               uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FriendId             uint64   `protobuf:"varint,2,opt,name=friend_id,json=friendId,proto3" json:"friend_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddFriendRequest) Reset()         { *m = AddFriendRequest{} }
func (m *AddFriendRequest) String() string { return proto.CompactTextString(m) }
func (*AddFriendRequest) ProtoMessage()    {}
func (*AddFriendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{0}
}

func (m *AddFriendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddFriendRequest.Unmarshal(m, b)
}
func (m *AddFriendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddFriendRequest.Marshal(b, m, deterministic)
}
func (m *AddFriendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddFriendRequest.Merge(m, src)
}
func (m *AddFriendRequest) XXX_Size() int {
	return xxx_messageInfo_AddFriendRequest.Size(m)
}
func (m *AddFriendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddFriendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddFriendRequest proto.InternalMessageInfo

func (m *AddFriendRequest) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *AddFriendRequest) GetFriendId() uint64 {
	if m != nil {
		return m.FriendId
	}
	return 0
}

type AddFriendResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddFriendResponse) Reset()         { *m = AddFriendResponse{} }
func (m *AddFriendResponse) String() string { return proto.CompactTextString(m) }
func (*AddFriendResponse) ProtoMessage()    {}
func (*AddFriendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{1}
}

func (m *AddFriendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddFriendResponse.Unmarshal(m, b)
}
func (m *AddFriendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddFriendResponse.Marshal(b, m, deterministic)
}
func (m *AddFriendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddFriendResponse.Merge(m, src)
}
func (m *AddFriendResponse) XXX_Size() int {
	return xxx_messageInfo_AddFriendResponse.Size(m)
}
func (m *AddFriendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddFriendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddFriendResponse proto.InternalMessageInfo

type DeleteFriendRequest struct {
	UserId               uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FriendId             uint64   `protobuf:"varint,2,opt,name=friend_id,json=friendId,proto3" json:"friend_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteFriendRequest) Reset()         { *m = DeleteFriendRequest{} }
func (m *DeleteFriendRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteFriendRequest) ProtoMessage()    {}
func (*DeleteFriendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{2}
}

func (m *DeleteFriendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteFriendRequest.Unmarshal(m, b)
}
func (m *DeleteFriendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteFriendRequest.Marshal(b, m, deterministic)
}
func (m *DeleteFriendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteFriendRequest.Merge(m, src)
}
func (m *DeleteFriendRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteFriendRequest.Size(m)
}
func (m *DeleteFriendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteFriendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteFriendRequest proto.InternalMessageInfo

func (m *DeleteFriendRequest) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *DeleteFriendRequest) GetFriendId() uint64 {
	if m != nil {
		return m.FriendId
	}
	return 0
}

type DeleteFriendResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteFriendResponse) Reset()         { *m = DeleteFriendResponse{} }
func (m *DeleteFriendResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteFriendResponse) ProtoMessage()    {}
func (*DeleteFriendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{3}
}

func (m *DeleteFriendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteFriendResponse.Unmarshal(m, b)
}
func (m *DeleteFriendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteFriendResponse.Marshal(b, m, deterministic)
}
func (m *DeleteFriendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteFriendResponse.Merge(m, src)
}
func (m *DeleteFriendResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteFriendResponse.Size(m)
}
func (m *DeleteFriendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteFriendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteFriendResponse proto.InternalMessageInfo

type GetFriendRequest struct {
	UserId               uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFriendRequest) Reset()         { *m = GetFriendRequest{} }
func (m *GetFriendRequest) String() string { return proto.CompactTextString(m) }
func (*GetFriendRequest) ProtoMessage()    {}
func (*GetFriendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{4}
}

func (m *GetFriendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFriendRequest.Unmarshal(m, b)
}
func (m *GetFriendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFriendRequest.Marshal(b, m, deterministic)
}
func (m *GetFriendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFriendRequest.Merge(m, src)
}
func (m *GetFriendRequest) XXX_Size() int {
	return xxx_messageInfo_GetFriendRequest.Size(m)
}
func (m *GetFriendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFriendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetFriendRequest proto.InternalMessageInfo

func (m *GetFriendRequest) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

type GetFriendResponse struct {
	Friends              []uint64 `protobuf:"varint,1,rep,packed,name=friends,proto3" json:"friends,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFriendResponse) Reset()         { *m = GetFriendResponse{} }
func (m *GetFriendResponse) String() string { return proto.CompactTextString(m) }
func (*GetFriendResponse) ProtoMessage()    {}
func (*GetFriendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{5}
}

func (m *GetFriendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFriendResponse.Unmarshal(m, b)
}
func (m *GetFriendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFriendResponse.Marshal(b, m, deterministic)
}
func (m *GetFriendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFriendResponse.Merge(m, src)
}
func (m *GetFriendResponse) XXX_Size() int {
	return xxx_messageInfo_GetFriendResponse.Size(m)
}
func (m *GetFriendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFriendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetFriendResponse proto.InternalMessageInfo

func (m *GetFriendResponse) GetFriends() []uint64 {
	if m != nil {
		return m.Friends
	}
	return nil
}

type GetRequestsFriendRequest struct {
	UserId               uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequestsFriendRequest) Reset()         { *m = GetRequestsFriendRequest{} }
func (m *GetRequestsFriendRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequestsFriendRequest) ProtoMessage()    {}
func (*GetRequestsFriendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{6}
}

func (m *GetRequestsFriendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequestsFriendRequest.Unmarshal(m, b)
}
func (m *GetRequestsFriendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequestsFriendRequest.Marshal(b, m, deterministic)
}
func (m *GetRequestsFriendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequestsFriendRequest.Merge(m, src)
}
func (m *GetRequestsFriendRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequestsFriendRequest.Size(m)
}
func (m *GetRequestsFriendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequestsFriendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequestsFriendRequest proto.InternalMessageInfo

func (m *GetRequestsFriendRequest) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

type GetRequestsFriendResponse struct {
	Requests             []uint64 `protobuf:"varint,1,rep,packed,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequestsFriendResponse) Reset()         { *m = GetRequestsFriendResponse{} }
func (m *GetRequestsFriendResponse) String() string { return proto.CompactTextString(m) }
func (*GetRequestsFriendResponse) ProtoMessage()    {}
func (*GetRequestsFriendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{7}
}

func (m *GetRequestsFriendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequestsFriendResponse.Unmarshal(m, b)
}
func (m *GetRequestsFriendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequestsFriendResponse.Marshal(b, m, deterministic)
}
func (m *GetRequestsFriendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequestsFriendResponse.Merge(m, src)
}
func (m *GetRequestsFriendResponse) XXX_Size() int {
	return xxx_messageInfo_GetRequestsFriendResponse.Size(m)
}
func (m *GetRequestsFriendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequestsFriendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequestsFriendResponse proto.InternalMessageInfo

func (m *GetRequestsFriendResponse) GetRequests() []uint64 {
	if m != nil {
		return m.Requests
	}
	return nil
}

type AcceptFriendRequest struct {
	UserId               uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FriendId             uint64   `protobuf:"varint,2,opt,name=friend_id,json=friendId,proto3" json:"friend_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AcceptFriendRequest) Reset()         { *m = AcceptFriendRequest{} }
func (m *AcceptFriendRequest) String() string { return proto.CompactTextString(m) }
func (*AcceptFriendRequest) ProtoMessage()    {}
func (*AcceptFriendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{8}
}

func (m *AcceptFriendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcceptFriendRequest.Unmarshal(m, b)
}
func (m *AcceptFriendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcceptFriendRequest.Marshal(b, m, deterministic)
}
func (m *AcceptFriendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcceptFriendRequest.Merge(m, src)
}
func (m *AcceptFriendRequest) XXX_Size() int {
	return xxx_messageInfo_AcceptFriendRequest.Size(m)
}
func (m *AcceptFriendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AcceptFriendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AcceptFriendRequest proto.InternalMessageInfo

func (m *AcceptFriendRequest) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *AcceptFriendRequest) GetFriendId() uint64 {
	if m != nil {
		return m.FriendId
	}
	return 0
}

type AcceptFriendResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AcceptFriendResponse) Reset()         { *m = AcceptFriendResponse{} }
func (m *AcceptFriendResponse) String() string { return proto.CompactTextString(m) }
func (*AcceptFriendResponse) ProtoMessage()    {}
func (*AcceptFriendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{9}
}

func (m *AcceptFriendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcceptFriendResponse.Unmarshal(m, b)
}
func (m *AcceptFriendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcceptFriendResponse.Marshal(b, m, deterministic)
}
func (m *AcceptFriendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcceptFriendResponse.Merge(m, src)
}
func (m *AcceptFriendResponse) XXX_Size() int {
	return xxx_messageInfo_AcceptFriendResponse.Size(m)
}
func (m *AcceptFriendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AcceptFriendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AcceptFriendResponse proto.InternalMessageInfo

type DenyFriendRequest struct {
	UserId               uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FriendId             uint64   `protobuf:"varint,2,opt,name=friend_id,json=friendId,proto3" json:"friend_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DenyFriendRequest) Reset()         { *m = DenyFriendRequest{} }
func (m *DenyFriendRequest) String() string { return proto.CompactTextString(m) }
func (*DenyFriendRequest) ProtoMessage()    {}
func (*DenyFriendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{10}
}

func (m *DenyFriendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DenyFriendRequest.Unmarshal(m, b)
}
func (m *DenyFriendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DenyFriendRequest.Marshal(b, m, deterministic)
}
func (m *DenyFriendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DenyFriendRequest.Merge(m, src)
}
func (m *DenyFriendRequest) XXX_Size() int {
	return xxx_messageInfo_DenyFriendRequest.Size(m)
}
func (m *DenyFriendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DenyFriendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DenyFriendRequest proto.InternalMessageInfo

func (m *DenyFriendRequest) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *DenyFriendRequest) GetFriendId() uint64 {
	if m != nil {
		return m.FriendId
	}
	return 0
}

type DenyFriendResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DenyFriendResponse) Reset()         { *m = DenyFriendResponse{} }
func (m *DenyFriendResponse) String() string { return proto.CompactTextString(m) }
func (*DenyFriendResponse) ProtoMessage()    {}
func (*DenyFriendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef2cbb5c12e56bfb, []int{11}
}

func (m *DenyFriendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DenyFriendResponse.Unmarshal(m, b)
}
func (m *DenyFriendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DenyFriendResponse.Marshal(b, m, deterministic)
}
func (m *DenyFriendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DenyFriendResponse.Merge(m, src)
}
func (m *DenyFriendResponse) XXX_Size() int {
	return xxx_messageInfo_DenyFriendResponse.Size(m)
}
func (m *DenyFriendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DenyFriendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DenyFriendResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AddFriendRequest)(nil), "pb.AddFriendRequest")
	proto.RegisterType((*AddFriendResponse)(nil), "pb.AddFriendResponse")
	proto.RegisterType((*DeleteFriendRequest)(nil), "pb.DeleteFriendRequest")
	proto.RegisterType((*DeleteFriendResponse)(nil), "pb.DeleteFriendResponse")
	proto.RegisterType((*GetFriendRequest)(nil), "pb.GetFriendRequest")
	proto.RegisterType((*GetFriendResponse)(nil), "pb.GetFriendResponse")
	proto.RegisterType((*GetRequestsFriendRequest)(nil), "pb.GetRequestsFriendRequest")
	proto.RegisterType((*GetRequestsFriendResponse)(nil), "pb.GetRequestsFriendResponse")
	proto.RegisterType((*AcceptFriendRequest)(nil), "pb.AcceptFriendRequest")
	proto.RegisterType((*AcceptFriendResponse)(nil), "pb.AcceptFriendResponse")
	proto.RegisterType((*DenyFriendRequest)(nil), "pb.DenyFriendRequest")
	proto.RegisterType((*DenyFriendResponse)(nil), "pb.DenyFriendResponse")
}

func init() { proto.RegisterFile("friends.proto", fileDescriptor_ef2cbb5c12e56bfb) }

var fileDescriptor_ef2cbb5c12e56bfb = []byte{
	// 332 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0x41, 0x6b, 0xb3, 0x40,
	0x10, 0x25, 0x2a, 0x9a, 0x6f, 0x3e, 0x0a, 0x71, 0x63, 0x12, 0x6b, 0x5b, 0x08, 0x9e, 0x02, 0xa5,
	0x1e, 0x0c, 0xa5, 0x87, 0x9e, 0x84, 0x50, 0x6b, 0x7b, 0xf3, 0x0f, 0x14, 0xcc, 0x4e, 0x21, 0x50,
	0xd4, 0xba, 0x9b, 0x43, 0xff, 0x74, 0x7f, 0x43, 0xd1, 0xd9, 0xb6, 0xc6, 0x58, 0xc8, 0xc1, 0xe3,
	0xce, 0x9b, 0x79, 0x6f, 0x7c, 0x6f, 0x84, 0xb3, 0xd7, 0x6a, 0x87, 0x39, 0x17, 0x41, 0x59, 0x15,
	0xb2, 0x60, 0x5a, 0x99, 0xf9, 0x8f, 0x30, 0x89, 0x38, 0x7f, 0x68, 0xea, 0x29, 0xbe, 0xef, 0x51,
	0x48, 0xb6, 0x00, 0x6b, 0x2f, 0xb0, 0x7a, 0xd9, 0x71, 0x77, 0xb4, 0x1c, 0xad, 0x8c, 0xd4, 0xac,
	0x9f, 0x09, 0x67, 0x17, 0xf0, 0x8f, 0x18, 0x6a, 0x48, 0x6b, 0xa0, 0x31, 0x15, 0x12, 0xee, 0x4f,
	0xc1, 0x6e, 0x31, 0x89, 0xb2, 0xc8, 0x05, 0xfa, 0xcf, 0x30, 0xdd, 0xe0, 0x1b, 0x4a, 0x1c, 0x42,
	0x61, 0x0e, 0xce, 0x21, 0x99, 0x12, 0xb9, 0x86, 0x49, 0x8c, 0xf2, 0x34, 0x05, 0xff, 0x06, 0xec,
	0x56, 0x33, 0x31, 0x30, 0x17, 0x2c, 0x65, 0x8d, 0x3b, 0x5a, 0xea, 0x2b, 0x23, 0xfd, 0x7e, 0xfa,
	0x6b, 0x70, 0x63, 0x94, 0x8a, 0x55, 0x9c, 0xa8, 0x71, 0x07, 0xe7, 0x3d, 0x43, 0x4a, 0xcb, 0x83,
	0x71, 0xa5, 0x10, 0x25, 0xf6, 0xf3, 0xae, 0xed, 0x8a, 0xb6, 0x5b, 0x2c, 0xe5, 0x40, 0x76, 0x1d,
	0x92, 0x29, 0xbb, 0x12, 0xb0, 0x37, 0x98, 0x7f, 0x0c, 0x21, 0xe1, 0x00, 0x6b, 0x53, 0x91, 0x40,
	0xf8, 0xa9, 0x81, 0x45, 0x25, 0xc1, 0x42, 0xd0, 0x23, 0xce, 0x99, 0x13, 0x94, 0x59, 0xd0, 0x3d,
	0x34, 0x6f, 0xd6, 0xa9, 0x2a, 0x87, 0xee, 0xc1, 0xa4, 0x9c, 0xd9, 0xa2, 0x6e, 0xe8, 0x39, 0x20,
	0xcf, 0x3d, 0x06, 0xd4, 0x70, 0x08, 0x7a, 0x8c, 0x92, 0x04, 0xbb, 0x57, 0x41, 0x82, 0xc7, 0xf1,
	0x3f, 0xc1, 0xff, 0x56, 0x5e, 0xec, 0x52, 0x75, 0xf5, 0xa6, 0xee, 0x5d, 0xfd, 0x81, 0xfe, 0x2e,
	0x4f, 0xae, 0xd3, 0xf2, 0x3d, 0x71, 0xd2, 0xf2, 0x7d, 0xd1, 0xb0, 0x5b, 0x30, 0x6a, 0x3f, 0xd9,
	0x8c, 0x3e, 0xaf, 0x13, 0x92, 0x37, 0xef, 0x96, 0x69, 0x2c, 0x33, 0x9b, 0xff, 0x79, 0xfd, 0x15,
	0x00, 0x00, 0xff, 0xff, 0x93, 0xf6, 0x01, 0x23, 0xe0, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FriendsClient is the client API for Friends service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FriendsClient interface {
	Add(ctx context.Context, in *AddFriendRequest, opts ...grpc.CallOption) (*AddFriendResponse, error)
	Delete(ctx context.Context, in *DeleteFriendRequest, opts ...grpc.CallOption) (*DeleteFriendResponse, error)
	Get(ctx context.Context, in *GetFriendRequest, opts ...grpc.CallOption) (*GetFriendResponse, error)
	GetRequests(ctx context.Context, in *GetRequestsFriendRequest, opts ...grpc.CallOption) (*GetRequestsFriendResponse, error)
	Accept(ctx context.Context, in *AcceptFriendRequest, opts ...grpc.CallOption) (*AcceptFriendResponse, error)
	Deny(ctx context.Context, in *DenyFriendRequest, opts ...grpc.CallOption) (*DenyFriendResponse, error)
}

type friendsClient struct {
	cc *grpc.ClientConn
}

func NewFriendsClient(cc *grpc.ClientConn) FriendsClient {
	return &friendsClient{cc}
}

func (c *friendsClient) Add(ctx context.Context, in *AddFriendRequest, opts ...grpc.CallOption) (*AddFriendResponse, error) {
	out := new(AddFriendResponse)
	err := c.cc.Invoke(ctx, "/pb.Friends/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsClient) Delete(ctx context.Context, in *DeleteFriendRequest, opts ...grpc.CallOption) (*DeleteFriendResponse, error) {
	out := new(DeleteFriendResponse)
	err := c.cc.Invoke(ctx, "/pb.Friends/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsClient) Get(ctx context.Context, in *GetFriendRequest, opts ...grpc.CallOption) (*GetFriendResponse, error) {
	out := new(GetFriendResponse)
	err := c.cc.Invoke(ctx, "/pb.Friends/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsClient) GetRequests(ctx context.Context, in *GetRequestsFriendRequest, opts ...grpc.CallOption) (*GetRequestsFriendResponse, error) {
	out := new(GetRequestsFriendResponse)
	err := c.cc.Invoke(ctx, "/pb.Friends/GetRequests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsClient) Accept(ctx context.Context, in *AcceptFriendRequest, opts ...grpc.CallOption) (*AcceptFriendResponse, error) {
	out := new(AcceptFriendResponse)
	err := c.cc.Invoke(ctx, "/pb.Friends/Accept", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsClient) Deny(ctx context.Context, in *DenyFriendRequest, opts ...grpc.CallOption) (*DenyFriendResponse, error) {
	out := new(DenyFriendResponse)
	err := c.cc.Invoke(ctx, "/pb.Friends/Deny", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FriendsServer is the server API for Friends service.
type FriendsServer interface {
	Add(context.Context, *AddFriendRequest) (*AddFriendResponse, error)
	Delete(context.Context, *DeleteFriendRequest) (*DeleteFriendResponse, error)
	Get(context.Context, *GetFriendRequest) (*GetFriendResponse, error)
	GetRequests(context.Context, *GetRequestsFriendRequest) (*GetRequestsFriendResponse, error)
	Accept(context.Context, *AcceptFriendRequest) (*AcceptFriendResponse, error)
	Deny(context.Context, *DenyFriendRequest) (*DenyFriendResponse, error)
}

// UnimplementedFriendsServer can be embedded to have forward compatible implementations.
type UnimplementedFriendsServer struct {
}

func (*UnimplementedFriendsServer) Add(ctx context.Context, req *AddFriendRequest) (*AddFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (*UnimplementedFriendsServer) Delete(ctx context.Context, req *DeleteFriendRequest) (*DeleteFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (*UnimplementedFriendsServer) Get(ctx context.Context, req *GetFriendRequest) (*GetFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedFriendsServer) GetRequests(ctx context.Context, req *GetRequestsFriendRequest) (*GetRequestsFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRequests not implemented")
}
func (*UnimplementedFriendsServer) Accept(ctx context.Context, req *AcceptFriendRequest) (*AcceptFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Accept not implemented")
}
func (*UnimplementedFriendsServer) Deny(ctx context.Context, req *DenyFriendRequest) (*DenyFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deny not implemented")
}

func RegisterFriendsServer(s *grpc.Server, srv FriendsServer) {
	s.RegisterService(&_Friends_serviceDesc, srv)
}

func _Friends_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Friends/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServer).Add(ctx, req.(*AddFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Friends_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Friends/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServer).Delete(ctx, req.(*DeleteFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Friends_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Friends/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServer).Get(ctx, req.(*GetFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Friends_GetRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequestsFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServer).GetRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Friends/GetRequests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServer).GetRequests(ctx, req.(*GetRequestsFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Friends_Accept_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServer).Accept(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Friends/Accept",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServer).Accept(ctx, req.(*AcceptFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Friends_Deny_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DenyFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServer).Deny(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Friends/Deny",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServer).Deny(ctx, req.(*DenyFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Friends_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Friends",
	HandlerType: (*FriendsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Friends_Add_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Friends_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Friends_Get_Handler,
		},
		{
			MethodName: "GetRequests",
			Handler:    _Friends_GetRequests_Handler,
		},
		{
			MethodName: "Accept",
			Handler:    _Friends_Accept_Handler,
		},
		{
			MethodName: "Deny",
			Handler:    _Friends_Deny_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "friends.proto",
}

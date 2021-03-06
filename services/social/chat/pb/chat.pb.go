// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type SendChatRequest struct {
	From                 uint64   `protobuf:"varint,1,opt,name=from,proto3" json:"from,omitempty"`
	To                   uint64   `protobuf:"varint,2,opt,name=to,proto3" json:"to,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendChatRequest) Reset()         { *m = SendChatRequest{} }
func (m *SendChatRequest) String() string { return proto.CompactTextString(m) }
func (*SendChatRequest) ProtoMessage()    {}
func (*SendChatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{0}
}

func (m *SendChatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendChatRequest.Unmarshal(m, b)
}
func (m *SendChatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendChatRequest.Marshal(b, m, deterministic)
}
func (m *SendChatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendChatRequest.Merge(m, src)
}
func (m *SendChatRequest) XXX_Size() int {
	return xxx_messageInfo_SendChatRequest.Size(m)
}
func (m *SendChatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendChatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendChatRequest proto.InternalMessageInfo

func (m *SendChatRequest) GetFrom() uint64 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *SendChatRequest) GetTo() uint64 {
	if m != nil {
		return m.To
	}
	return 0
}

func (m *SendChatRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type SendChatResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendChatResponse) Reset()         { *m = SendChatResponse{} }
func (m *SendChatResponse) String() string { return proto.CompactTextString(m) }
func (*SendChatResponse) ProtoMessage()    {}
func (*SendChatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{1}
}

func (m *SendChatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendChatResponse.Unmarshal(m, b)
}
func (m *SendChatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendChatResponse.Marshal(b, m, deterministic)
}
func (m *SendChatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendChatResponse.Merge(m, src)
}
func (m *SendChatResponse) XXX_Size() int {
	return xxx_messageInfo_SendChatResponse.Size(m)
}
func (m *SendChatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SendChatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SendChatResponse proto.InternalMessageInfo

type GetChatRequest struct {
	User1                uint64   `protobuf:"varint,1,opt,name=user1,proto3" json:"user1,omitempty"`
	User2                uint64   `protobuf:"varint,2,opt,name=user2,proto3" json:"user2,omitempty"`
	Total                uint64   `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`
	Skip                 uint64   `protobuf:"varint,4,opt,name=skip,proto3" json:"skip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetChatRequest) Reset()         { *m = GetChatRequest{} }
func (m *GetChatRequest) String() string { return proto.CompactTextString(m) }
func (*GetChatRequest) ProtoMessage()    {}
func (*GetChatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{2}
}

func (m *GetChatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetChatRequest.Unmarshal(m, b)
}
func (m *GetChatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetChatRequest.Marshal(b, m, deterministic)
}
func (m *GetChatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetChatRequest.Merge(m, src)
}
func (m *GetChatRequest) XXX_Size() int {
	return xxx_messageInfo_GetChatRequest.Size(m)
}
func (m *GetChatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetChatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetChatRequest proto.InternalMessageInfo

func (m *GetChatRequest) GetUser1() uint64 {
	if m != nil {
		return m.User1
	}
	return 0
}

func (m *GetChatRequest) GetUser2() uint64 {
	if m != nil {
		return m.User2
	}
	return 0
}

func (m *GetChatRequest) GetTotal() uint64 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *GetChatRequest) GetSkip() uint64 {
	if m != nil {
		return m.Skip
	}
	return 0
}

type GetChatResponse struct {
	Messages             []*ChatMessage `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetChatResponse) Reset()         { *m = GetChatResponse{} }
func (m *GetChatResponse) String() string { return proto.CompactTextString(m) }
func (*GetChatResponse) ProtoMessage()    {}
func (*GetChatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{3}
}

func (m *GetChatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetChatResponse.Unmarshal(m, b)
}
func (m *GetChatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetChatResponse.Marshal(b, m, deterministic)
}
func (m *GetChatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetChatResponse.Merge(m, src)
}
func (m *GetChatResponse) XXX_Size() int {
	return xxx_messageInfo_GetChatResponse.Size(m)
}
func (m *GetChatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetChatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetChatResponse proto.InternalMessageInfo

func (m *GetChatResponse) GetMessages() []*ChatMessage {
	if m != nil {
		return m.Messages
	}
	return nil
}

type ChatMessage struct {
	From                 uint64               `protobuf:"varint,1,opt,name=from,proto3" json:"from,omitempty"`
	To                   uint64               `protobuf:"varint,2,opt,name=to,proto3" json:"to,omitempty"`
	Message              string               `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	Time                 *timestamp.Timestamp `protobuf:"bytes,4,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ChatMessage) Reset()         { *m = ChatMessage{} }
func (m *ChatMessage) String() string { return proto.CompactTextString(m) }
func (*ChatMessage) ProtoMessage()    {}
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{4}
}

func (m *ChatMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChatMessage.Unmarshal(m, b)
}
func (m *ChatMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChatMessage.Marshal(b, m, deterministic)
}
func (m *ChatMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChatMessage.Merge(m, src)
}
func (m *ChatMessage) XXX_Size() int {
	return xxx_messageInfo_ChatMessage.Size(m)
}
func (m *ChatMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ChatMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ChatMessage proto.InternalMessageInfo

func (m *ChatMessage) GetFrom() uint64 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *ChatMessage) GetTo() uint64 {
	if m != nil {
		return m.To
	}
	return 0
}

func (m *ChatMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ChatMessage) GetTime() *timestamp.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func init() {
	proto.RegisterType((*SendChatRequest)(nil), "pb.SendChatRequest")
	proto.RegisterType((*SendChatResponse)(nil), "pb.SendChatResponse")
	proto.RegisterType((*GetChatRequest)(nil), "pb.GetChatRequest")
	proto.RegisterType((*GetChatResponse)(nil), "pb.GetChatResponse")
	proto.RegisterType((*ChatMessage)(nil), "pb.ChatMessage")
}

func init() { proto.RegisterFile("chat.proto", fileDescriptor_8c585a45e2093e54) }

var fileDescriptor_8c585a45e2093e54 = []byte{
	// 296 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x51, 0x3d, 0x4f, 0xc3, 0x30,
	0x14, 0x54, 0x52, 0xf3, 0xf5, 0x2a, 0xb5, 0xe8, 0xb5, 0x83, 0x95, 0x85, 0x2a, 0x53, 0x25, 0x24,
	0x57, 0x0d, 0x3b, 0x0b, 0x43, 0x27, 0x84, 0x64, 0xf8, 0x03, 0x09, 0x38, 0x69, 0x44, 0x53, 0x9b,
	0xf8, 0x65, 0xe2, 0xcf, 0x23, 0xdb, 0x49, 0x69, 0x58, 0xd9, 0x7c, 0xf7, 0x4e, 0xb9, 0xcb, 0x1d,
	0xc0, 0xfb, 0x3e, 0x27, 0x61, 0x5a, 0x4d, 0x1a, 0x63, 0x53, 0x24, 0x77, 0x95, 0xd6, 0xd5, 0x41,
	0x6d, 0x3c, 0x53, 0x74, 0xe5, 0x86, 0xea, 0x46, 0x59, 0xca, 0x1b, 0x13, 0x44, 0xe9, 0x0b, 0xcc,
	0x5f, 0xd5, 0xf1, 0xe3, 0x69, 0x9f, 0x93, 0x54, 0x5f, 0x9d, 0xb2, 0x84, 0x08, 0xac, 0x6c, 0x75,
	0xc3, 0xa3, 0x55, 0xb4, 0x66, 0xd2, 0xbf, 0x71, 0x06, 0x31, 0x69, 0x1e, 0x7b, 0x26, 0x26, 0x8d,
	0x1c, 0xae, 0x1a, 0x65, 0x6d, 0x5e, 0x29, 0x3e, 0x59, 0x45, 0xeb, 0x1b, 0x39, 0xc0, 0x14, 0xe1,
	0xf6, 0xf7, 0x83, 0xd6, 0xe8, 0xa3, 0x55, 0x69, 0x09, 0xb3, 0x9d, 0xa2, 0x73, 0x8f, 0x25, 0x5c,
	0x74, 0x56, 0xb5, 0xdb, 0xde, 0x24, 0x80, 0x81, 0xcd, 0x7a, 0xa3, 0x00, 0x1c, 0x4b, 0x9a, 0xf2,
	0x83, 0x77, 0x62, 0x32, 0x00, 0x97, 0xd2, 0x7e, 0xd6, 0x86, 0xb3, 0x90, 0xd2, 0xbd, 0xd3, 0x47,
	0x98, 0x9f, 0x7c, 0x82, 0x35, 0xde, 0xc3, 0x75, 0x9f, 0xcc, 0xf2, 0x68, 0x35, 0x59, 0x4f, 0xb3,
	0xb9, 0x30, 0x85, 0x70, 0x9a, 0xe7, 0xc0, 0xcb, 0x93, 0x20, 0xfd, 0x86, 0xe9, 0xd9, 0xe1, 0x7f,
	0x45, 0xa0, 0x00, 0xe6, 0xca, 0xf6, 0x01, 0xa7, 0x59, 0x22, 0xc2, 0x12, 0x62, 0x58, 0x42, 0xbc,
	0x0d, 0x4b, 0x48, 0xaf, 0xcb, 0x6a, 0x60, 0xce, 0x1c, 0xb7, 0xc0, 0x5c, 0x81, 0xb8, 0x70, 0x39,
	0xff, 0x6c, 0x93, 0x2c, 0xc7, 0x64, 0xff, 0x93, 0x02, 0x26, 0x3b, 0x45, 0x88, 0xee, 0x38, 0x2e,
	0x3a, 0x59, 0x8c, 0xb8, 0xa0, 0x2f, 0x2e, 0x7d, 0x88, 0x87, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x51, 0x49, 0xc3, 0x4e, 0x2e, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatClient interface {
	Send(ctx context.Context, in *SendChatRequest, opts ...grpc.CallOption) (*SendChatResponse, error)
	Get(ctx context.Context, in *GetChatRequest, opts ...grpc.CallOption) (*GetChatResponse, error)
}

type chatClient struct {
	cc *grpc.ClientConn
}

func NewChatClient(cc *grpc.ClientConn) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) Send(ctx context.Context, in *SendChatRequest, opts ...grpc.CallOption) (*SendChatResponse, error) {
	out := new(SendChatResponse)
	err := c.cc.Invoke(ctx, "/pb.Chat/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) Get(ctx context.Context, in *GetChatRequest, opts ...grpc.CallOption) (*GetChatResponse, error) {
	out := new(GetChatResponse)
	err := c.cc.Invoke(ctx, "/pb.Chat/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
type ChatServer interface {
	Send(context.Context, *SendChatRequest) (*SendChatResponse, error)
	Get(context.Context, *GetChatRequest) (*GetChatResponse, error)
}

// UnimplementedChatServer can be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (*UnimplementedChatServer) Send(ctx context.Context, req *SendChatRequest) (*SendChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (*UnimplementedChatServer) Get(ctx context.Context, req *GetChatRequest) (*GetChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func RegisterChatServer(s *grpc.Server, srv ChatServer) {
	s.RegisterService(&_Chat_serviceDesc, srv)
}

func _Chat_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Chat/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).Send(ctx, req.(*SendChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Chat/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).Get(ctx, req.(*GetChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Chat_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Chat_Send_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Chat_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat.proto",
}

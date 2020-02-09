// Code generated by protoc-gen-go. DO NOT EDIT.
// source: live.proto

package pb

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type GetLiveRequest struct {
	MatchId              uint64   `protobuf:"varint,1,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetLiveRequest) Reset()         { *m = GetLiveRequest{} }
func (m *GetLiveRequest) String() string { return proto.CompactTextString(m) }
func (*GetLiveRequest) ProtoMessage()    {}
func (*GetLiveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{0}
}

func (m *GetLiveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetLiveRequest.Unmarshal(m, b)
}
func (m *GetLiveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetLiveRequest.Marshal(b, m, deterministic)
}
func (m *GetLiveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetLiveRequest.Merge(m, src)
}
func (m *GetLiveRequest) XXX_Size() int {
	return xxx_messageInfo_GetLiveRequest.Size(m)
}
func (m *GetLiveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetLiveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetLiveRequest proto.InternalMessageInfo

func (m *GetLiveRequest) GetMatchId() uint64 {
	if m != nil {
		return m.MatchId
	}
	return 0
}

type GetLiveResponse struct {
	MatchId              uint64               `protobuf:"varint,1,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
	Team1                *TeamLiveInfo        `protobuf:"bytes,2,opt,name=team1,proto3" json:"team1,omitempty"`
	Team2                *TeamLiveInfo        `protobuf:"bytes,3,opt,name=team2,proto3" json:"team2,omitempty"`
	StartTime            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GetLiveResponse) Reset()         { *m = GetLiveResponse{} }
func (m *GetLiveResponse) String() string { return proto.CompactTextString(m) }
func (*GetLiveResponse) ProtoMessage()    {}
func (*GetLiveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{1}
}

func (m *GetLiveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetLiveResponse.Unmarshal(m, b)
}
func (m *GetLiveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetLiveResponse.Marshal(b, m, deterministic)
}
func (m *GetLiveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetLiveResponse.Merge(m, src)
}
func (m *GetLiveResponse) XXX_Size() int {
	return xxx_messageInfo_GetLiveResponse.Size(m)
}
func (m *GetLiveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetLiveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetLiveResponse proto.InternalMessageInfo

func (m *GetLiveResponse) GetMatchId() uint64 {
	if m != nil {
		return m.MatchId
	}
	return 0
}

func (m *GetLiveResponse) GetTeam1() *TeamLiveInfo {
	if m != nil {
		return m.Team1
	}
	return nil
}

func (m *GetLiveResponse) GetTeam2() *TeamLiveInfo {
	if m != nil {
		return m.Team2
	}
	return nil
}

func (m *GetLiveResponse) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

type FindMultipleLiveRequest struct {
	Total                uint32   `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	RatingUnder          uint32   `protobuf:"varint,2,opt,name=rating_under,json=ratingUnder,proto3" json:"rating_under,omitempty"`
	RatingOver           uint32   `protobuf:"varint,3,opt,name=rating_over,json=ratingOver,proto3" json:"rating_over,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindMultipleLiveRequest) Reset()         { *m = FindMultipleLiveRequest{} }
func (m *FindMultipleLiveRequest) String() string { return proto.CompactTextString(m) }
func (*FindMultipleLiveRequest) ProtoMessage()    {}
func (*FindMultipleLiveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{2}
}

func (m *FindMultipleLiveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindMultipleLiveRequest.Unmarshal(m, b)
}
func (m *FindMultipleLiveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindMultipleLiveRequest.Marshal(b, m, deterministic)
}
func (m *FindMultipleLiveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindMultipleLiveRequest.Merge(m, src)
}
func (m *FindMultipleLiveRequest) XXX_Size() int {
	return xxx_messageInfo_FindMultipleLiveRequest.Size(m)
}
func (m *FindMultipleLiveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FindMultipleLiveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FindMultipleLiveRequest proto.InternalMessageInfo

func (m *FindMultipleLiveRequest) GetTotal() uint32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *FindMultipleLiveRequest) GetRatingUnder() uint32 {
	if m != nil {
		return m.RatingUnder
	}
	return 0
}

func (m *FindMultipleLiveRequest) GetRatingOver() uint32 {
	if m != nil {
		return m.RatingOver
	}
	return 0
}

type FindMultipleLiveResponse struct {
	Matches              []*GetLiveResponse `protobuf:"bytes,1,rep,name=matches,proto3" json:"matches,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *FindMultipleLiveResponse) Reset()         { *m = FindMultipleLiveResponse{} }
func (m *FindMultipleLiveResponse) String() string { return proto.CompactTextString(m) }
func (*FindMultipleLiveResponse) ProtoMessage()    {}
func (*FindMultipleLiveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{3}
}

func (m *FindMultipleLiveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindMultipleLiveResponse.Unmarshal(m, b)
}
func (m *FindMultipleLiveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindMultipleLiveResponse.Marshal(b, m, deterministic)
}
func (m *FindMultipleLiveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindMultipleLiveResponse.Merge(m, src)
}
func (m *FindMultipleLiveResponse) XXX_Size() int {
	return xxx_messageInfo_FindMultipleLiveResponse.Size(m)
}
func (m *FindMultipleLiveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FindMultipleLiveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FindMultipleLiveResponse proto.InternalMessageInfo

func (m *FindMultipleLiveResponse) GetMatches() []*GetLiveResponse {
	if m != nil {
		return m.Matches
	}
	return nil
}

type FindUserLiveRequest struct {
	UserId               uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindUserLiveRequest) Reset()         { *m = FindUserLiveRequest{} }
func (m *FindUserLiveRequest) String() string { return proto.CompactTextString(m) }
func (*FindUserLiveRequest) ProtoMessage()    {}
func (*FindUserLiveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{4}
}

func (m *FindUserLiveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindUserLiveRequest.Unmarshal(m, b)
}
func (m *FindUserLiveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindUserLiveRequest.Marshal(b, m, deterministic)
}
func (m *FindUserLiveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindUserLiveRequest.Merge(m, src)
}
func (m *FindUserLiveRequest) XXX_Size() int {
	return xxx_messageInfo_FindUserLiveRequest.Size(m)
}
func (m *FindUserLiveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FindUserLiveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FindUserLiveRequest proto.InternalMessageInfo

func (m *FindUserLiveRequest) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

type FindUserLiveResponse struct {
	UserId               uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	MatchId              uint64   `protobuf:"varint,2,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindUserLiveResponse) Reset()         { *m = FindUserLiveResponse{} }
func (m *FindUserLiveResponse) String() string { return proto.CompactTextString(m) }
func (*FindUserLiveResponse) ProtoMessage()    {}
func (*FindUserLiveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{5}
}

func (m *FindUserLiveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindUserLiveResponse.Unmarshal(m, b)
}
func (m *FindUserLiveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindUserLiveResponse.Marshal(b, m, deterministic)
}
func (m *FindUserLiveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindUserLiveResponse.Merge(m, src)
}
func (m *FindUserLiveResponse) XXX_Size() int {
	return xxx_messageInfo_FindUserLiveResponse.Size(m)
}
func (m *FindUserLiveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FindUserLiveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FindUserLiveResponse proto.InternalMessageInfo

func (m *FindUserLiveResponse) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *FindUserLiveResponse) GetMatchId() uint64 {
	if m != nil {
		return m.MatchId
	}
	return 0
}

type GetTotalLiveRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTotalLiveRequest) Reset()         { *m = GetTotalLiveRequest{} }
func (m *GetTotalLiveRequest) String() string { return proto.CompactTextString(m) }
func (*GetTotalLiveRequest) ProtoMessage()    {}
func (*GetTotalLiveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{6}
}

func (m *GetTotalLiveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTotalLiveRequest.Unmarshal(m, b)
}
func (m *GetTotalLiveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTotalLiveRequest.Marshal(b, m, deterministic)
}
func (m *GetTotalLiveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTotalLiveRequest.Merge(m, src)
}
func (m *GetTotalLiveRequest) XXX_Size() int {
	return xxx_messageInfo_GetTotalLiveRequest.Size(m)
}
func (m *GetTotalLiveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTotalLiveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTotalLiveRequest proto.InternalMessageInfo

type GetTotalLiveResponse struct {
	Total                uint64   `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTotalLiveResponse) Reset()         { *m = GetTotalLiveResponse{} }
func (m *GetTotalLiveResponse) String() string { return proto.CompactTextString(m) }
func (*GetTotalLiveResponse) ProtoMessage()    {}
func (*GetTotalLiveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{7}
}

func (m *GetTotalLiveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTotalLiveResponse.Unmarshal(m, b)
}
func (m *GetTotalLiveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTotalLiveResponse.Marshal(b, m, deterministic)
}
func (m *GetTotalLiveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTotalLiveResponse.Merge(m, src)
}
func (m *GetTotalLiveResponse) XXX_Size() int {
	return xxx_messageInfo_GetTotalLiveResponse.Size(m)
}
func (m *GetTotalLiveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTotalLiveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTotalLiveResponse proto.InternalMessageInfo

func (m *GetTotalLiveResponse) GetTotal() uint64 {
	if m != nil {
		return m.Total
	}
	return 0
}

type TeamLiveInfo struct {
	Users                []uint64 `protobuf:"varint,1,rep,packed,name=users,proto3" json:"users,omitempty"`
	Score                uint32   `protobuf:"varint,2,opt,name=score,proto3" json:"score,omitempty"`
	AverageRating        uint64   `protobuf:"varint,3,opt,name=average_rating,json=averageRating,proto3" json:"average_rating,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamLiveInfo) Reset()         { *m = TeamLiveInfo{} }
func (m *TeamLiveInfo) String() string { return proto.CompactTextString(m) }
func (*TeamLiveInfo) ProtoMessage()    {}
func (*TeamLiveInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{8}
}

func (m *TeamLiveInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamLiveInfo.Unmarshal(m, b)
}
func (m *TeamLiveInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamLiveInfo.Marshal(b, m, deterministic)
}
func (m *TeamLiveInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamLiveInfo.Merge(m, src)
}
func (m *TeamLiveInfo) XXX_Size() int {
	return xxx_messageInfo_TeamLiveInfo.Size(m)
}
func (m *TeamLiveInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamLiveInfo.DiscardUnknown(m)
}

var xxx_messageInfo_TeamLiveInfo proto.InternalMessageInfo

func (m *TeamLiveInfo) GetUsers() []uint64 {
	if m != nil {
		return m.Users
	}
	return nil
}

func (m *TeamLiveInfo) GetScore() uint32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *TeamLiveInfo) GetAverageRating() uint64 {
	if m != nil {
		return m.AverageRating
	}
	return 0
}

type AddLiveRequest struct {
	MatchId              uint64               `protobuf:"varint,1,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
	Team1                *TeamLiveInfo        `protobuf:"bytes,2,opt,name=team1,proto3" json:"team1,omitempty"`
	Team2                *TeamLiveInfo        `protobuf:"bytes,3,opt,name=team2,proto3" json:"team2,omitempty"`
	StartTime            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *AddLiveRequest) Reset()         { *m = AddLiveRequest{} }
func (m *AddLiveRequest) String() string { return proto.CompactTextString(m) }
func (*AddLiveRequest) ProtoMessage()    {}
func (*AddLiveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{9}
}

func (m *AddLiveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddLiveRequest.Unmarshal(m, b)
}
func (m *AddLiveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddLiveRequest.Marshal(b, m, deterministic)
}
func (m *AddLiveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddLiveRequest.Merge(m, src)
}
func (m *AddLiveRequest) XXX_Size() int {
	return xxx_messageInfo_AddLiveRequest.Size(m)
}
func (m *AddLiveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddLiveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddLiveRequest proto.InternalMessageInfo

func (m *AddLiveRequest) GetMatchId() uint64 {
	if m != nil {
		return m.MatchId
	}
	return 0
}

func (m *AddLiveRequest) GetTeam1() *TeamLiveInfo {
	if m != nil {
		return m.Team1
	}
	return nil
}

func (m *AddLiveRequest) GetTeam2() *TeamLiveInfo {
	if m != nil {
		return m.Team2
	}
	return nil
}

func (m *AddLiveRequest) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

type AddLiveResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddLiveResponse) Reset()         { *m = AddLiveResponse{} }
func (m *AddLiveResponse) String() string { return proto.CompactTextString(m) }
func (*AddLiveResponse) ProtoMessage()    {}
func (*AddLiveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{10}
}

func (m *AddLiveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddLiveResponse.Unmarshal(m, b)
}
func (m *AddLiveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddLiveResponse.Marshal(b, m, deterministic)
}
func (m *AddLiveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddLiveResponse.Merge(m, src)
}
func (m *AddLiveResponse) XXX_Size() int {
	return xxx_messageInfo_AddLiveResponse.Size(m)
}
func (m *AddLiveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddLiveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddLiveResponse proto.InternalMessageInfo

type RemoveLiveRequest struct {
	MatchId              uint64   `protobuf:"varint,1,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveLiveRequest) Reset()         { *m = RemoveLiveRequest{} }
func (m *RemoveLiveRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveLiveRequest) ProtoMessage()    {}
func (*RemoveLiveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{11}
}

func (m *RemoveLiveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveLiveRequest.Unmarshal(m, b)
}
func (m *RemoveLiveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveLiveRequest.Marshal(b, m, deterministic)
}
func (m *RemoveLiveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveLiveRequest.Merge(m, src)
}
func (m *RemoveLiveRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveLiveRequest.Size(m)
}
func (m *RemoveLiveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveLiveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveLiveRequest proto.InternalMessageInfo

func (m *RemoveLiveRequest) GetMatchId() uint64 {
	if m != nil {
		return m.MatchId
	}
	return 0
}

type RemoveLiveResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveLiveResponse) Reset()         { *m = RemoveLiveResponse{} }
func (m *RemoveLiveResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveLiveResponse) ProtoMessage()    {}
func (*RemoveLiveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9b1f08afb3d995e, []int{12}
}

func (m *RemoveLiveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveLiveResponse.Unmarshal(m, b)
}
func (m *RemoveLiveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveLiveResponse.Marshal(b, m, deterministic)
}
func (m *RemoveLiveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveLiveResponse.Merge(m, src)
}
func (m *RemoveLiveResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveLiveResponse.Size(m)
}
func (m *RemoveLiveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveLiveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveLiveResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GetLiveRequest)(nil), "pb.GetLiveRequest")
	proto.RegisterType((*GetLiveResponse)(nil), "pb.GetLiveResponse")
	proto.RegisterType((*FindMultipleLiveRequest)(nil), "pb.FindMultipleLiveRequest")
	proto.RegisterType((*FindMultipleLiveResponse)(nil), "pb.FindMultipleLiveResponse")
	proto.RegisterType((*FindUserLiveRequest)(nil), "pb.FindUserLiveRequest")
	proto.RegisterType((*FindUserLiveResponse)(nil), "pb.FindUserLiveResponse")
	proto.RegisterType((*GetTotalLiveRequest)(nil), "pb.GetTotalLiveRequest")
	proto.RegisterType((*GetTotalLiveResponse)(nil), "pb.GetTotalLiveResponse")
	proto.RegisterType((*TeamLiveInfo)(nil), "pb.TeamLiveInfo")
	proto.RegisterType((*AddLiveRequest)(nil), "pb.AddLiveRequest")
	proto.RegisterType((*AddLiveResponse)(nil), "pb.AddLiveResponse")
	proto.RegisterType((*RemoveLiveRequest)(nil), "pb.RemoveLiveRequest")
	proto.RegisterType((*RemoveLiveResponse)(nil), "pb.RemoveLiveResponse")
}

func init() { proto.RegisterFile("live.proto", fileDescriptor_c9b1f08afb3d995e) }

var fileDescriptor_c9b1f08afb3d995e = []byte{
	// 537 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x53, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x55, 0xda, 0xac, 0x1b, 0xb7, 0x1f, 0x63, 0x6e, 0x47, 0x43, 0x40, 0x5a, 0x89, 0x04, 0xaa,
	0x04, 0x64, 0xa2, 0x3c, 0x20, 0x1e, 0x78, 0xd8, 0x0b, 0x55, 0x11, 0x08, 0xc9, 0xea, 0x9e, 0x2b,
	0x77, 0xb9, 0x2b, 0x91, 0x9a, 0x3a, 0xd8, 0x4e, 0x7e, 0x1a, 0xfc, 0x25, 0x7e, 0x06, 0xb2, 0x9d,
	0x6c, 0x49, 0xdb, 0x4d, 0x7b, 0xdc, 0xe3, 0x3d, 0x3e, 0xf6, 0xbd, 0xe7, 0xdc, 0x63, 0x80, 0x75,
	0x9c, 0x63, 0x98, 0x0a, 0xae, 0x38, 0x69, 0xa4, 0x4b, 0xff, 0x6c, 0xc5, 0xf9, 0x6a, 0x8d, 0xe7,
	0x06, 0x59, 0x66, 0xd7, 0xe7, 0x2a, 0x4e, 0x50, 0x2a, 0x96, 0xa4, 0x96, 0x14, 0xbc, 0x85, 0xde,
	0x14, 0xd5, 0xf7, 0x38, 0x47, 0x8a, 0xbf, 0x33, 0x94, 0x8a, 0x3c, 0x87, 0xa3, 0x84, 0xa9, 0xab,
	0x5f, 0x8b, 0x38, 0xf2, 0x9c, 0x91, 0x33, 0x76, 0xe9, 0xa1, 0xa9, 0x67, 0x51, 0xf0, 0xd7, 0x81,
	0xe3, 0x1b, 0xb6, 0x4c, 0xf9, 0x46, 0xe2, 0x3d, 0x74, 0xf2, 0x06, 0x0e, 0x14, 0xb2, 0xe4, 0x83,
	0xd7, 0x18, 0x39, 0xe3, 0xf6, 0xe4, 0x69, 0x98, 0x2e, 0xc3, 0x39, 0xb2, 0x44, 0xdf, 0x9f, 0x6d,
	0xae, 0x39, 0xb5, 0xc7, 0x25, 0x6f, 0xe2, 0x35, 0xef, 0xe3, 0x4d, 0xc8, 0x67, 0x00, 0xa9, 0x98,
	0x50, 0x0b, 0x2d, 0xc2, 0x73, 0x0d, 0xd9, 0x0f, 0xad, 0xc2, 0xb0, 0x54, 0x18, 0xce, 0x4b, 0x85,
	0xf4, 0x89, 0x61, 0xeb, 0x3a, 0x90, 0x30, 0xfc, 0x1a, 0x6f, 0xa2, 0x1f, 0xd9, 0x5a, 0xc5, 0xe9,
	0x1a, 0xab, 0x7a, 0x07, 0x70, 0xa0, 0xb8, 0x62, 0x6b, 0x33, 0x7d, 0x97, 0xda, 0x82, 0xbc, 0x82,
	0x8e, 0x60, 0x2a, 0xde, 0xac, 0x16, 0xd9, 0x26, 0x42, 0x61, 0x24, 0x74, 0x69, 0xdb, 0x62, 0x97,
	0x1a, 0x22, 0x67, 0x50, 0x94, 0x0b, 0x9e, 0xa3, 0x30, 0xc3, 0x77, 0x29, 0x58, 0xe8, 0x67, 0x8e,
	0x22, 0x98, 0x81, 0xb7, 0xdb, 0xb4, 0xb0, 0xed, 0x3d, 0x58, 0x9b, 0x50, 0x7a, 0xce, 0xa8, 0x39,
	0x6e, 0x4f, 0xfa, 0x5a, 0xf5, 0x96, 0xb9, 0xb4, 0xe4, 0x04, 0x21, 0xf4, 0xf5, 0x53, 0x97, 0x12,
	0x45, 0x75, 0xf6, 0x21, 0x1c, 0x66, 0x12, 0xc5, 0xad, 0xf7, 0x2d, 0x5d, 0xce, 0xa2, 0xe0, 0x1b,
	0x0c, 0xea, 0xfc, 0xa2, 0xed, 0x5d, 0x17, 0x6a, 0x6b, 0x6c, 0xd4, 0xb7, 0x7e, 0x0a, 0xfd, 0x29,
	0xaa, 0xb9, 0xb6, 0xa5, 0xd2, 0x3b, 0x78, 0x07, 0x83, 0x3a, 0x5c, 0xb4, 0xa8, 0xf9, 0xe9, 0x16,
	0x7e, 0x06, 0x0c, 0x3a, 0xd5, 0x95, 0x6a, 0x96, 0xee, 0x6c, 0xd5, 0xbb, 0xd4, 0x16, 0x1a, 0x95,
	0x57, 0x5c, 0x60, 0x61, 0xb7, 0x2d, 0xc8, 0x6b, 0xe8, 0xb1, 0x1c, 0x05, 0x5b, 0xe1, 0xc2, 0xba,
	0x6b, 0xbc, 0x76, 0x69, 0xb7, 0x40, 0xa9, 0x01, 0x83, 0x3f, 0x0e, 0xf4, 0x2e, 0xa2, 0xe8, 0x61,
	0x59, 0x7e, 0x4c, 0xe1, 0x3c, 0x81, 0xe3, 0x9b, 0xb9, 0xad, 0x89, 0x41, 0x08, 0x27, 0x14, 0x13,
	0x9e, 0xe3, 0x03, 0x7f, 0xe6, 0x00, 0x48, 0x95, 0x6f, 0x5f, 0x99, 0xfc, 0x6b, 0x80, 0xab, 0x01,
	0x12, 0x42, 0x73, 0x8a, 0x8a, 0x90, 0x5a, 0xc6, 0xcc, 0xa3, 0xfe, 0xbe, 0xdc, 0x91, 0x2f, 0x70,
	0x54, 0xee, 0x96, 0x0c, 0x0b, 0xc2, 0x76, 0x00, 0x7c, 0x6f, 0xf7, 0xa0, 0xb8, 0x3e, 0x83, 0x4e,
	0x35, 0xf8, 0xe4, 0x85, 0x66, 0xde, 0xf1, 0xff, 0xfc, 0x97, 0xfb, 0x0f, 0x6f, 0x27, 0x29, 0x83,
	0x6c, 0x27, 0xd9, 0xf3, 0x0d, 0xec, 0x24, 0x7b, 0xf3, 0x1e, 0x42, 0xf3, 0x22, 0x8a, 0xac, 0xf0,
	0x7a, 0x36, 0xac, 0xf0, 0x2d, 0xdf, 0xc9, 0x27, 0x68, 0x59, 0x1f, 0xc9, 0xa9, 0x3e, 0xde, 0xd9,
	0x81, 0xff, 0x6c, 0x1b, 0xb6, 0x17, 0x97, 0x2d, 0xb3, 0xe2, 0x8f, 0xff, 0x03, 0x00, 0x00, 0xff,
	0xff, 0xf9, 0x84, 0x00, 0x24, 0x81, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LiveClient is the client API for Live service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LiveClient interface {
	Get(ctx context.Context, in *GetLiveRequest, opts ...grpc.CallOption) (*GetLiveResponse, error)
	GetTotal(ctx context.Context, in *GetTotalLiveRequest, opts ...grpc.CallOption) (*GetTotalLiveResponse, error)
	FindMultiple(ctx context.Context, in *FindMultipleLiveRequest, opts ...grpc.CallOption) (*FindMultipleLiveResponse, error)
	FindUser(ctx context.Context, in *FindUserLiveRequest, opts ...grpc.CallOption) (*FindUserLiveResponse, error)
	Add(ctx context.Context, in *AddLiveRequest, opts ...grpc.CallOption) (*AddLiveResponse, error)
	Remove(ctx context.Context, in *RemoveLiveRequest, opts ...grpc.CallOption) (*RemoveLiveResponse, error)
}

type liveClient struct {
	cc *grpc.ClientConn
}

func NewLiveClient(cc *grpc.ClientConn) LiveClient {
	return &liveClient{cc}
}

func (c *liveClient) Get(ctx context.Context, in *GetLiveRequest, opts ...grpc.CallOption) (*GetLiveResponse, error) {
	out := new(GetLiveResponse)
	err := c.cc.Invoke(ctx, "/pb.Live/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClient) GetTotal(ctx context.Context, in *GetTotalLiveRequest, opts ...grpc.CallOption) (*GetTotalLiveResponse, error) {
	out := new(GetTotalLiveResponse)
	err := c.cc.Invoke(ctx, "/pb.Live/GetTotal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClient) FindMultiple(ctx context.Context, in *FindMultipleLiveRequest, opts ...grpc.CallOption) (*FindMultipleLiveResponse, error) {
	out := new(FindMultipleLiveResponse)
	err := c.cc.Invoke(ctx, "/pb.Live/FindMultiple", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClient) FindUser(ctx context.Context, in *FindUserLiveRequest, opts ...grpc.CallOption) (*FindUserLiveResponse, error) {
	out := new(FindUserLiveResponse)
	err := c.cc.Invoke(ctx, "/pb.Live/FindUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClient) Add(ctx context.Context, in *AddLiveRequest, opts ...grpc.CallOption) (*AddLiveResponse, error) {
	out := new(AddLiveResponse)
	err := c.cc.Invoke(ctx, "/pb.Live/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClient) Remove(ctx context.Context, in *RemoveLiveRequest, opts ...grpc.CallOption) (*RemoveLiveResponse, error) {
	out := new(RemoveLiveResponse)
	err := c.cc.Invoke(ctx, "/pb.Live/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LiveServer is the server API for Live service.
type LiveServer interface {
	Get(context.Context, *GetLiveRequest) (*GetLiveResponse, error)
	GetTotal(context.Context, *GetTotalLiveRequest) (*GetTotalLiveResponse, error)
	FindMultiple(context.Context, *FindMultipleLiveRequest) (*FindMultipleLiveResponse, error)
	FindUser(context.Context, *FindUserLiveRequest) (*FindUserLiveResponse, error)
	Add(context.Context, *AddLiveRequest) (*AddLiveResponse, error)
	Remove(context.Context, *RemoveLiveRequest) (*RemoveLiveResponse, error)
}

// UnimplementedLiveServer can be embedded to have forward compatible implementations.
type UnimplementedLiveServer struct {
}

func (*UnimplementedLiveServer) Get(ctx context.Context, req *GetLiveRequest) (*GetLiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedLiveServer) GetTotal(ctx context.Context, req *GetTotalLiveRequest) (*GetTotalLiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotal not implemented")
}
func (*UnimplementedLiveServer) FindMultiple(ctx context.Context, req *FindMultipleLiveRequest) (*FindMultipleLiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindMultiple not implemented")
}
func (*UnimplementedLiveServer) FindUser(ctx context.Context, req *FindUserLiveRequest) (*FindUserLiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindUser not implemented")
}
func (*UnimplementedLiveServer) Add(ctx context.Context, req *AddLiveRequest) (*AddLiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (*UnimplementedLiveServer) Remove(ctx context.Context, req *RemoveLiveRequest) (*RemoveLiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

func RegisterLiveServer(s *grpc.Server, srv LiveServer) {
	s.RegisterService(&_Live_serviceDesc, srv)
}

func _Live_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Live/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveServer).Get(ctx, req.(*GetLiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Live_GetTotal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTotalLiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveServer).GetTotal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Live/GetTotal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveServer).GetTotal(ctx, req.(*GetTotalLiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Live_FindMultiple_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindMultipleLiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveServer).FindMultiple(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Live/FindMultiple",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveServer).FindMultiple(ctx, req.(*FindMultipleLiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Live_FindUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindUserLiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveServer).FindUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Live/FindUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveServer).FindUser(ctx, req.(*FindUserLiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Live_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddLiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Live/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveServer).Add(ctx, req.(*AddLiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Live_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveLiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Live/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveServer).Remove(ctx, req.(*RemoveLiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Live_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Live",
	HandlerType: (*LiveServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Live_Get_Handler,
		},
		{
			MethodName: "GetTotal",
			Handler:    _Live_GetTotal_Handler,
		},
		{
			MethodName: "FindMultiple",
			Handler:    _Live_FindMultiple_Handler,
		},
		{
			MethodName: "FindUser",
			Handler:    _Live_FindUser_Handler,
		},
		{
			MethodName: "Add",
			Handler:    _Live_Add_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Live_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "live.proto",
}

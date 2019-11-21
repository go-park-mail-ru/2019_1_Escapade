// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat.proto

/*
Package chat is a generated protocol buffer package.

It is generated from these files:
	chat.proto

It has these top-level messages:
	ChatID
	MessageID
	MessagesID
	Result
	User
	Message
	Messages
	Chat
	ChatWithUsers
	UserInGroup
*/
package proto

import (
	fmt "fmt"

	proto "github.com/golang/protobuf/proto"

	math "math"

	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

	context "golang.org/x/net/context"

	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Status int32

const (
	Status_NO       Status = 0
	Status_OBSERVER Status = 1
	Status_PLAYER   Status = 2
	Status_ADMIN    Status = 3
)

var Status_name = map[int32]string{
	0: "NO",
	1: "OBSERVER",
	2: "PLAYER",
	3: "ADMIN",
}
var Status_value = map[string]int32{
	"NO":       0,
	"OBSERVER": 1,
	"PLAYER":   2,
	"ADMIN":    3,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}
func (Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ChatID struct {
	Value int32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *ChatID) Reset()                    { *m = ChatID{} }
func (m *ChatID) String() string            { return proto.CompactTextString(m) }
func (*ChatID) ProtoMessage()               {}
func (*ChatID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ChatID) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type MessageID struct {
	Value int32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *MessageID) Reset()                    { *m = MessageID{} }
func (m *MessageID) String() string            { return proto.CompactTextString(m) }
func (*MessageID) ProtoMessage()               {}
func (*MessageID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MessageID) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type MessagesID struct {
	Values []*MessageID `protobuf:"bytes,1,rep,name=values" json:"values,omitempty"`
}

func (m *MessagesID) Reset()                    { *m = MessagesID{} }
func (m *MessagesID) String() string            { return proto.CompactTextString(m) }
func (*MessagesID) ProtoMessage()               {}
func (*MessagesID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MessagesID) GetValues() []*MessageID {
	if m != nil {
		return m.Values
	}
	return nil
}

type Result struct {
	Done bool `protobuf:"varint,1,opt,name=done" json:"done,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Result) GetDone() bool {
	if m != nil {
		return m.Done
	}
	return false
}

type User struct {
	Id     int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Photo  string `protobuf:"bytes,3,opt,name=photo" json:"photo,omitempty"`
	Status Status `protobuf:"varint,4,opt,name=status,enum=chat.Status" json:"status,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *User) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetPhoto() string {
	if m != nil {
		return m.Photo
	}
	return ""
}

func (m *User) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_NO
}

type Message struct {
	Id     int32                      `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Answer *Message                   `protobuf:"bytes,2,opt,name=answer" json:"answer,omitempty"`
	Text   string                     `protobuf:"bytes,3,opt,name=text" json:"text,omitempty"`
	From   *User                      `protobuf:"bytes,4,opt,name=from" json:"from,omitempty"`
	To     *User                      `protobuf:"bytes,5,opt,name=to" json:"to,omitempty"`
	ChatId int32                      `protobuf:"varint,6,opt,name=chat_id,json=chatId" json:"chat_id,omitempty"`
	Time   *google_protobuf.Timestamp `protobuf:"bytes,7,opt,name=time" json:"time,omitempty"`
	Edited bool                       `protobuf:"varint,8,opt,name=edited" json:"edited,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Message) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Message) GetAnswer() *Message {
	if m != nil {
		return m.Answer
	}
	return nil
}

func (m *Message) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Message) GetFrom() *User {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Message) GetTo() *User {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Message) GetChatId() int32 {
	if m != nil {
		return m.ChatId
	}
	return 0
}

func (m *Message) GetTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *Message) GetEdited() bool {
	if m != nil {
		return m.Edited
	}
	return false
}

type Messages struct {
	Messages    []*Message `protobuf:"bytes,1,rep,name=messages" json:"messages,omitempty"`
	BlockSize   int32      `protobuf:"varint,2,opt,name=block_size,json=blockSize" json:"block_size,omitempty"`
	BlockAmount int32      `protobuf:"varint,3,opt,name=block_amount,json=blockAmount" json:"block_amount,omitempty"`
	BlockNumber int32      `protobuf:"varint,4,opt,name=block_number,json=blockNumber" json:"block_number,omitempty"`
}

func (m *Messages) Reset()                    { *m = Messages{} }
func (m *Messages) String() string            { return proto.CompactTextString(m) }
func (*Messages) ProtoMessage()               {}
func (*Messages) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Messages) GetMessages() []*Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *Messages) GetBlockSize() int32 {
	if m != nil {
		return m.BlockSize
	}
	return 0
}

func (m *Messages) GetBlockAmount() int32 {
	if m != nil {
		return m.BlockAmount
	}
	return 0
}

func (m *Messages) GetBlockNumber() int32 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}

type Chat struct {
	Id       int32       `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Type     int32       `protobuf:"varint,2,opt,name=type" json:"type,omitempty"`
	TypeId   int32       `protobuf:"varint,3,opt,name=type_id,json=typeId" json:"type_id,omitempty"`
	Messages []*Messages `protobuf:"bytes,4,rep,name=messages" json:"messages,omitempty"`
}

func (m *Chat) Reset()                    { *m = Chat{} }
func (m *Chat) String() string            { return proto.CompactTextString(m) }
func (*Chat) ProtoMessage()               {}
func (*Chat) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Chat) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Chat) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Chat) GetTypeId() int32 {
	if m != nil {
		return m.TypeId
	}
	return 0
}

func (m *Chat) GetMessages() []*Messages {
	if m != nil {
		return m.Messages
	}
	return nil
}

type ChatWithUsers struct {
	Type   int32   `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	TypeId int32   `protobuf:"varint,2,opt,name=type_id,json=typeId" json:"type_id,omitempty"`
	Users  []*User `protobuf:"bytes,3,rep,name=users" json:"users,omitempty"`
}

func (m *ChatWithUsers) Reset()                    { *m = ChatWithUsers{} }
func (m *ChatWithUsers) String() string            { return proto.CompactTextString(m) }
func (*ChatWithUsers) ProtoMessage()               {}
func (*ChatWithUsers) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ChatWithUsers) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *ChatWithUsers) GetTypeId() int32 {
	if m != nil {
		return m.TypeId
	}
	return 0
}

func (m *ChatWithUsers) GetUsers() []*User {
	if m != nil {
		return m.Users
	}
	return nil
}

type UserInGroup struct {
	User *User `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
	Chat *Chat `protobuf:"bytes,2,opt,name=chat" json:"chat,omitempty"`
}

func (m *UserInGroup) Reset()                    { *m = UserInGroup{} }
func (m *UserInGroup) String() string            { return proto.CompactTextString(m) }
func (*UserInGroup) ProtoMessage()               {}
func (*UserInGroup) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *UserInGroup) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *UserInGroup) GetChat() *Chat {
	if m != nil {
		return m.Chat
	}
	return nil
}

func init() {
	proto.RegisterType((*ChatID)(nil), "chat.ChatID")
	proto.RegisterType((*MessageID)(nil), "chat.MessageID")
	proto.RegisterType((*MessagesID)(nil), "chat.MessagesID")
	proto.RegisterType((*Result)(nil), "chat.Result")
	proto.RegisterType((*User)(nil), "chat.User")
	proto.RegisterType((*Message)(nil), "chat.Message")
	proto.RegisterType((*Messages)(nil), "chat.Messages")
	proto.RegisterType((*Chat)(nil), "chat.Chat")
	proto.RegisterType((*ChatWithUsers)(nil), "chat.ChatWithUsers")
	proto.RegisterType((*UserInGroup)(nil), "chat.UserInGroup")
	proto.RegisterEnum("chat.Status", Status_name, Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ChatService service

type ChatServiceClient interface {
	CreateChat(ctx context.Context, in *ChatWithUsers, opts ...grpc.CallOption) (*ChatID, error)
	GetChat(ctx context.Context, in *Chat, opts ...grpc.CallOption) (*ChatID, error)
	InviteToChat(ctx context.Context, in *UserInGroup, opts ...grpc.CallOption) (*Result, error)
	LeaveChat(ctx context.Context, in *UserInGroup, opts ...grpc.CallOption) (*Result, error)
	AppendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*MessageID, error)
	AppendMessages(ctx context.Context, in *Messages, opts ...grpc.CallOption) (*MessagesID, error)
	UpdateMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Result, error)
	DeleteMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Result, error)
	GetChatMessages(ctx context.Context, in *ChatID, opts ...grpc.CallOption) (*Messages, error)
}

type chatServiceClient struct {
	cc *grpc.ClientConn
}

func NewChatServiceClient(cc *grpc.ClientConn) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) CreateChat(ctx context.Context, in *ChatWithUsers, opts ...grpc.CallOption) (*ChatID, error) {
	out := new(ChatID)
	err := grpc.Invoke(ctx, "/chat.ChatService/CreateChat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetChat(ctx context.Context, in *Chat, opts ...grpc.CallOption) (*ChatID, error) {
	out := new(ChatID)
	err := grpc.Invoke(ctx, "/chat.ChatService/GetChat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) InviteToChat(ctx context.Context, in *UserInGroup, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/chat.ChatService/InviteToChat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) LeaveChat(ctx context.Context, in *UserInGroup, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/chat.ChatService/LeaveChat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) AppendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*MessageID, error) {
	out := new(MessageID)
	err := grpc.Invoke(ctx, "/chat.ChatService/AppendMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) AppendMessages(ctx context.Context, in *Messages, opts ...grpc.CallOption) (*MessagesID, error) {
	out := new(MessagesID)
	err := grpc.Invoke(ctx, "/chat.ChatService/AppendMessages", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) UpdateMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/chat.ChatService/UpdateMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) DeleteMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/chat.ChatService/DeleteMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetChatMessages(ctx context.Context, in *ChatID, opts ...grpc.CallOption) (*Messages, error) {
	out := new(Messages)
	err := grpc.Invoke(ctx, "/chat.ChatService/GetChatMessages", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ChatService service

type ChatServiceServer interface {
	CreateChat(context.Context, *ChatWithUsers) (*ChatID, error)
	GetChat(context.Context, *Chat) (*ChatID, error)
	InviteToChat(context.Context, *UserInGroup) (*Result, error)
	LeaveChat(context.Context, *UserInGroup) (*Result, error)
	AppendMessage(context.Context, *Message) (*MessageID, error)
	AppendMessages(context.Context, *Messages) (*MessagesID, error)
	UpdateMessage(context.Context, *Message) (*Result, error)
	DeleteMessage(context.Context, *Message) (*Result, error)
	GetChatMessages(context.Context, *ChatID) (*Messages, error)
}

func RegisterChatServiceServer(s *grpc.Server, srv ChatServiceServer) {
	s.RegisterService(&_ChatService_serviceDesc, srv)
}

func _ChatService_CreateChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatWithUsers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).CreateChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/CreateChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).CreateChat(ctx, req.(*ChatWithUsers))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Chat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/GetChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetChat(ctx, req.(*Chat))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_InviteToChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInGroup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).InviteToChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/InviteToChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).InviteToChat(ctx, req.(*UserInGroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_LeaveChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInGroup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).LeaveChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/LeaveChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).LeaveChat(ctx, req.(*UserInGroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_AppendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).AppendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/AppendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).AppendMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_AppendMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Messages)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).AppendMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/AppendMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).AppendMessages(ctx, req.(*Messages))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_UpdateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).UpdateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/UpdateMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).UpdateMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_DeleteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).DeleteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/DeleteMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).DeleteMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetChatMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetChatMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/GetChatMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetChatMessages(ctx, req.(*ChatID))
	}
	return interceptor(ctx, in, info, handler)
}

var _ChatService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chat.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateChat",
			Handler:    _ChatService_CreateChat_Handler,
		},
		{
			MethodName: "GetChat",
			Handler:    _ChatService_GetChat_Handler,
		},
		{
			MethodName: "InviteToChat",
			Handler:    _ChatService_InviteToChat_Handler,
		},
		{
			MethodName: "LeaveChat",
			Handler:    _ChatService_LeaveChat_Handler,
		},
		{
			MethodName: "AppendMessage",
			Handler:    _ChatService_AppendMessage_Handler,
		},
		{
			MethodName: "AppendMessages",
			Handler:    _ChatService_AppendMessages_Handler,
		},
		{
			MethodName: "UpdateMessage",
			Handler:    _ChatService_UpdateMessage_Handler,
		},
		{
			MethodName: "DeleteMessage",
			Handler:    _ChatService_DeleteMessage_Handler,
		},
		{
			MethodName: "GetChatMessages",
			Handler:    _ChatService_GetChatMessages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat.proto",
}

func init() { proto.RegisterFile("chat.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 723 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xdd, 0x6e, 0xe2, 0x46,
	0x14, 0xc6, 0xc6, 0x18, 0x38, 0xfc, 0x84, 0x4e, 0xab, 0xd6, 0x42, 0x6d, 0x44, 0xac, 0x46, 0xa5,
	0xb9, 0x30, 0x82, 0x36, 0x57, 0xbd, 0xa2, 0x21, 0x8a, 0x2c, 0xe5, 0x67, 0x65, 0x92, 0x5d, 0xed,
	0xcd, 0x46, 0x06, 0x4f, 0xc0, 0x5a, 0xec, 0xb1, 0x3c, 0x63, 0x76, 0x37, 0x8f, 0xb2, 0x0f, 0xb4,
	0xef, 0xb3, 0x6f, 0xb0, 0x9a, 0xe3, 0x01, 0x02, 0x24, 0x52, 0xae, 0x38, 0xe7, 0x3b, 0xdf, 0x99,
	0xef, 0xfc, 0x19, 0x80, 0xe9, 0xdc, 0x17, 0x4e, 0x92, 0x32, 0xc1, 0x88, 0x21, 0xed, 0xf6, 0x7f,
	0xb3, 0x50, 0xcc, 0xb3, 0x89, 0x33, 0x65, 0x51, 0x6f, 0xc6, 0x16, 0x7e, 0x3c, 0xeb, 0x61, 0x78,
	0x92, 0x3d, 0xf4, 0x12, 0xf1, 0x25, 0xa1, 0xbc, 0x27, 0xc2, 0x88, 0x72, 0xe1, 0x47, 0xc9, 0xc6,
	0xca, 0x9f, 0xb0, 0x0f, 0xc1, 0x3c, 0x9b, 0xfb, 0xc2, 0x1d, 0x91, 0x5f, 0xa0, 0xb4, 0xf4, 0x17,
	0x19, 0xb5, 0xb4, 0x8e, 0xd6, 0x2d, 0x79, 0xb9, 0x63, 0x1f, 0x41, 0xf5, 0x8a, 0x72, 0xee, 0xcf,
	0xe8, 0x8b, 0x94, 0x53, 0x00, 0x45, 0xe1, 0xee, 0x88, 0xfc, 0x05, 0x26, 0xc2, 0xdc, 0xd2, 0x3a,
	0xc5, 0x6e, 0x6d, 0x70, 0xe0, 0x60, 0xc1, 0xeb, 0x47, 0x3c, 0x15, 0xb6, 0x7f, 0x07, 0xd3, 0xa3,
	0x3c, 0x5b, 0x08, 0x42, 0xc0, 0x08, 0x58, 0x9c, 0xbf, 0x5a, 0xf1, 0xd0, 0xb6, 0x1f, 0xc0, 0xb8,
	0xe3, 0x34, 0x25, 0x4d, 0xd0, 0xc3, 0x40, 0xe9, 0xe9, 0x61, 0x20, 0xb9, 0xb1, 0x1f, 0x51, 0x4b,
	0xef, 0x68, 0xdd, 0xaa, 0x87, 0xb6, 0x2c, 0x2b, 0x99, 0x33, 0xc1, 0xac, 0x22, 0x82, 0xb9, 0x43,
	0xfe, 0x04, 0x93, 0x0b, 0x5f, 0x64, 0xdc, 0x32, 0x3a, 0x5a, 0xb7, 0x39, 0xa8, 0xe7, 0x85, 0x8c,
	0x11, 0xf3, 0x54, 0xcc, 0xfe, 0xae, 0x41, 0x59, 0xd5, 0xb6, 0xa7, 0x75, 0x0c, 0xa6, 0x1f, 0xf3,
	0x4f, 0x34, 0x45, 0xb5, 0xda, 0xa0, 0xb1, 0xd5, 0x8a, 0xa7, 0x82, 0xb2, 0x24, 0x41, 0x3f, 0x0b,
	0xa5, 0x8e, 0x36, 0x39, 0x04, 0xe3, 0x21, 0x65, 0x11, 0x4a, 0xd7, 0x06, 0x90, 0x27, 0xca, 0x86,
	0x3c, 0xc4, 0x49, 0x1b, 0x74, 0xc1, 0xac, 0xd2, 0x5e, 0x54, 0x17, 0x8c, 0xfc, 0x06, 0x65, 0x09,
	0xdc, 0x87, 0x81, 0x65, 0x62, 0x2d, 0xa6, 0x74, 0xdd, 0x80, 0x38, 0x60, 0xc8, 0xf5, 0x59, 0x65,
	0x4c, 0x6b, 0x3b, 0x33, 0xc6, 0x66, 0x0b, 0xea, 0xac, 0x96, 0xed, 0xdc, 0xae, 0x76, 0xeb, 0x21,
	0x8f, 0xfc, 0x0a, 0x26, 0x0d, 0x42, 0x41, 0x03, 0xab, 0x82, 0x93, 0x55, 0x9e, 0xfd, 0x55, 0x83,
	0xca, 0x6a, 0x63, 0xe4, 0x6f, 0xa8, 0x44, 0xca, 0x56, 0x1b, 0xdb, 0x69, 0x73, 0x1d, 0x26, 0x7f,
	0x00, 0x4c, 0x16, 0x6c, 0xfa, 0xf1, 0x9e, 0x87, 0x8f, 0xf9, 0x06, 0x4a, 0x5e, 0x15, 0x91, 0x71,
	0xf8, 0x48, 0xc9, 0x11, 0xd4, 0xf3, 0xb0, 0x1f, 0xb1, 0x2c, 0xce, 0xe7, 0x51, 0xf2, 0x6a, 0x88,
	0x0d, 0x11, 0xda, 0x50, 0xe2, 0x2c, 0x9a, 0xd0, 0x14, 0xc7, 0xb3, 0xa2, 0x5c, 0x23, 0x64, 0x33,
	0x30, 0xe4, 0x41, 0x3e, 0xb7, 0x78, 0x79, 0xcb, 0x4a, 0x16, 0x6d, 0x39, 0x29, 0xf9, 0x2b, 0x27,
	0x95, 0x8b, 0x99, 0xd2, 0x75, 0x03, 0x72, 0xf2, 0xa4, 0x29, 0x03, 0x9b, 0x6a, 0x6e, 0x35, 0xc5,
	0x37, 0x5d, 0xd9, 0x1f, 0xa0, 0x21, 0x05, 0xdf, 0x85, 0x62, 0x2e, 0x57, 0xc0, 0xd7, 0x4a, 0xda,
	0xf3, 0x4a, 0xfa, 0x96, 0x52, 0x07, 0x4a, 0x99, 0xcc, 0xb2, 0x8a, 0x28, 0xf3, 0x74, 0x97, 0x79,
	0xc0, 0xbe, 0x82, 0x9a, 0x74, 0xdd, 0xf8, 0x22, 0x65, 0x59, 0x22, 0x2f, 0x43, 0xe2, 0xf8, 0xfa,
	0xce, 0x65, 0x48, 0x5c, 0xc6, 0x25, 0xa4, 0x4e, 0x4e, 0xc5, 0x65, 0x81, 0x1e, 0xe2, 0x27, 0xa7,
	0x60, 0xe6, 0x27, 0x4c, 0x4c, 0xd0, 0xaf, 0x6f, 0x5a, 0x05, 0x52, 0x87, 0xca, 0xcd, 0xff, 0xe3,
	0x73, 0xef, 0xed, 0xb9, 0xd7, 0xd2, 0x08, 0x80, 0xf9, 0xe6, 0x72, 0xf8, 0xfe, 0xdc, 0x6b, 0xe9,
	0xa4, 0x0a, 0xa5, 0xe1, 0xe8, 0xca, 0xbd, 0x6e, 0x15, 0x07, 0xdf, 0x8a, 0x50, 0x93, 0xaf, 0x8c,
	0x69, 0xba, 0x0c, 0xa7, 0x94, 0xf4, 0x01, 0xce, 0x52, 0xea, 0x0b, 0x8a, 0xc3, 0xfe, 0x79, 0x23,
	0xb3, 0x9e, 0x43, 0xbb, 0xbe, 0x01, 0xdd, 0x91, 0x5d, 0x20, 0xc7, 0x50, 0xbe, 0xa0, 0x02, 0xf9,
	0x4f, 0xca, 0xda, 0xa3, 0xf5, 0xa1, 0xee, 0xc6, 0xcb, 0x50, 0xd0, 0x5b, 0x86, 0xdc, 0x9f, 0x36,
	0x2d, 0xaa, 0x19, 0xac, 0x52, 0xf2, 0xcf, 0xdf, 0x2e, 0x10, 0x07, 0xaa, 0x97, 0xd4, 0x5f, 0xd2,
	0xd7, 0xf2, 0xfb, 0xd0, 0x18, 0x26, 0x09, 0x8d, 0x83, 0xd5, 0x97, 0xbb, 0x7d, 0xb2, 0xed, 0xdd,
	0xff, 0x1c, 0xbb, 0x40, 0xfe, 0x85, 0xe6, 0x56, 0x0a, 0x27, 0x3b, 0x17, 0xd1, 0x6e, 0x6d, 0xfb,
	0x98, 0xe5, 0x40, 0xe3, 0x2e, 0x09, 0x7c, 0x41, 0x5f, 0x10, 0xda, 0x6f, 0xa4, 0x31, 0xa2, 0x0b,
	0xfa, 0x6a, 0x7e, 0x1f, 0x0e, 0xd4, 0x48, 0xd7, 0x65, 0x6d, 0x8d, 0xb3, 0xbd, 0x53, 0xa4, 0x5d,
	0x98, 0x98, 0xf8, 0xb9, 0xff, 0xf3, 0x23, 0x00, 0x00, 0xff, 0xff, 0x2e, 0x13, 0x7d, 0x76, 0x08,
	0x06, 0x00, 0x00,
}
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.12.4
// source: base.proto

package base

import (
	_ "go-social-network/biz/model/api"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ErrCode int32

const (
	ErrCode_Success                ErrCode = 0
	ErrCode_ArgumentError          ErrCode = 1
	ErrCode_CreateUserError        ErrCode = 2
	ErrCode_CopierError            ErrCode = 3
	ErrCode_GetUserInfoError       ErrCode = 4
	ErrCode_PostStatusError        ErrCode = 5
	ErrCode_GetTimelineError       ErrCode = 6
	ErrCode_FollowAndUnfollowError ErrCode = 7
	ErrCode_DeleteStatusError      ErrCode = 8
	ErrCode_CreateChatError        ErrCode = 9
	ErrCode_UnauthorizedError      ErrCode = 10
	ErrCode_PostMessageError       ErrCode = 11
	ErrCode_GetPendingMessageError ErrCode = 12
	ErrCode_LeaveChatError         ErrCode = 13
	ErrCode_SearchUserError        ErrCode = 14
	ErrCode_GetFollowingsError     ErrCode = 15
	ErrCode_GetFollowersError      ErrCode = 16
	ErrCode_GetFriendsError        ErrCode = 17
	ErrCode_GetChatListError       ErrCode = 18
	ErrCode_ToggleLikeStatusError  ErrCode = 19
	ErrCode_GetHotError            ErrCode = 20
)

// Enum value maps for ErrCode.
var (
	ErrCode_name = map[int32]string{
		0:  "Success",
		1:  "ArgumentError",
		2:  "CreateUserError",
		3:  "CopierError",
		4:  "GetUserInfoError",
		5:  "PostStatusError",
		6:  "GetTimelineError",
		7:  "FollowAndUnfollowError",
		8:  "DeleteStatusError",
		9:  "CreateChatError",
		10: "UnauthorizedError",
		11: "PostMessageError",
		12: "GetPendingMessageError",
		13: "LeaveChatError",
		14: "SearchUserError",
		15: "GetFollowingsError",
		16: "GetFollowersError",
		17: "GetFriendsError",
		18: "GetChatListError",
		19: "ToggleLikeStatusError",
		20: "GetHotError",
	}
	ErrCode_value = map[string]int32{
		"Success":                0,
		"ArgumentError":          1,
		"CreateUserError":        2,
		"CopierError":            3,
		"GetUserInfoError":       4,
		"PostStatusError":        5,
		"GetTimelineError":       6,
		"FollowAndUnfollowError": 7,
		"DeleteStatusError":      8,
		"CreateChatError":        9,
		"UnauthorizedError":      10,
		"PostMessageError":       11,
		"GetPendingMessageError": 12,
		"LeaveChatError":         13,
		"SearchUserError":        14,
		"GetFollowingsError":     15,
		"GetFollowersError":      16,
		"GetFriendsError":        17,
		"GetChatListError":       18,
		"ToggleLikeStatusError":  19,
		"GetHotError":            20,
	}
)

func (x ErrCode) Enum() *ErrCode {
	p := new(ErrCode)
	*p = x
	return p
}

func (x ErrCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrCode) Descriptor() protoreflect.EnumDescriptor {
	return file_base_proto_enumTypes[0].Descriptor()
}

func (ErrCode) Type() protoreflect.EnumType {
	return &file_base_proto_enumTypes[0]
}

func (x ErrCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrCode.Descriptor instead.
func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{0}
}

type BaseResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrCode ErrCode `protobuf:"varint,1,opt,name=errCode,proto3,enum=base.ErrCode" json:"errCode,omitempty" form:"errCode" query:"errCode"`
	ErrMsg  string  `protobuf:"bytes,2,opt,name=errMsg,proto3" json:"errMsg,omitempty" form:"errMsg" query:"errMsg"`
}

func (x *BaseResp) Reset() {
	*x = BaseResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_base_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseResp) ProtoMessage() {}

func (x *BaseResp) ProtoReflect() protoreflect.Message {
	mi := &file_base_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseResp.ProtoReflect.Descriptor instead.
func (*BaseResp) Descriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{0}
}

func (x *BaseResp) GetErrCode() ErrCode {
	if x != nil {
		return x.ErrCode
	}
	return ErrCode_Success
}

func (x *BaseResp) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

var File_base_proto protoreflect.FileDescriptor

var file_base_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x62, 0x61,
	0x73, 0x65, 0x1a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4b, 0x0a,
	0x08, 0x42, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x27, 0x0a, 0x07, 0x65, 0x72, 0x72,
	0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x2a, 0xd0, 0x03, 0x0a, 0x07, 0x45,
	0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x43,
	0x6f, 0x70, 0x69, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x10, 0x04, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x05, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x06, 0x12, 0x1a, 0x0a,
	0x16, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x41, 0x6e, 0x64, 0x55, 0x6e, 0x66, 0x6f, 0x6c, 0x6c,
	0x6f, 0x77, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x07, 0x12, 0x15, 0x0a, 0x11, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x08,
	0x12, 0x13, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x74, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x10, 0x09, 0x12, 0x15, 0x0a, 0x11, 0x55, 0x6e, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x65, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x0a, 0x12, 0x14, 0x0a, 0x10,
	0x50, 0x6f, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x10, 0x0b, 0x12, 0x1a, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x0c, 0x12, 0x12,
	0x0a, 0x0e, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x43, 0x68, 0x61, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x10, 0x0d, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x55, 0x73, 0x65, 0x72,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x0e, 0x12, 0x16, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x46, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x69, 0x6e, 0x67, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x0f, 0x12,
	0x15, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x73, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x10, 0x10, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x46, 0x72, 0x69,
	0x65, 0x6e, 0x64, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x11, 0x12, 0x14, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x43, 0x68, 0x61, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10,
	0x12, 0x12, 0x19, 0x0a, 0x15, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x13, 0x12, 0x0f, 0x0a, 0x0b,
	0x47, 0x65, 0x74, 0x48, 0x6f, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x14, 0x42, 0x22, 0x5a,
	0x20, 0x67, 0x6f, 0x2d, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x6c, 0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x2f, 0x62, 0x69, 0x7a, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x62, 0x61, 0x73,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_base_proto_rawDescOnce sync.Once
	file_base_proto_rawDescData = file_base_proto_rawDesc
)

func file_base_proto_rawDescGZIP() []byte {
	file_base_proto_rawDescOnce.Do(func() {
		file_base_proto_rawDescData = protoimpl.X.CompressGZIP(file_base_proto_rawDescData)
	})
	return file_base_proto_rawDescData
}

var file_base_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_base_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_base_proto_goTypes = []interface{}{
	(ErrCode)(0),     // 0: base.ErrCode
	(*BaseResp)(nil), // 1: base.BaseResp
}
var file_base_proto_depIdxs = []int32{
	0, // 0: base.BaseResp.errCode:type_name -> base.ErrCode
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_base_proto_init() }
func file_base_proto_init() {
	if File_base_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_base_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_base_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_base_proto_goTypes,
		DependencyIndexes: file_base_proto_depIdxs,
		EnumInfos:         file_base_proto_enumTypes,
		MessageInfos:      file_base_proto_msgTypes,
	}.Build()
	File_base_proto = out.File
	file_base_proto_rawDesc = nil
	file_base_proto_goTypes = nil
	file_base_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: api/serviceCenter/v1/user.proto

package v1

import (
	v1 "gitee.com/moyusir/util/api/util/v1"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// 注册请求
type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 用户注册信息
	User *v1.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	// 配置注册信息
	DeviceConfigRegisterInfos []*v1.DeviceConfigRegisterInfo `protobuf:"bytes,2,rep,name=device_config_register_infos,json=deviceConfigRegisterInfos,proto3" json:"device_config_register_infos,omitempty"`
	// 设备状态及预警规则注册信息
	DeviceStateRegisterInfos []*v1.DeviceStateRegisterInfo `protobuf:"bytes,3,rep,name=device_state_register_infos,json=deviceStateRegisterInfos,proto3" json:"device_state_register_infos,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_serviceCenter_v1_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_serviceCenter_v1_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_api_serviceCenter_v1_user_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterRequest) GetUser() *v1.User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *RegisterRequest) GetDeviceConfigRegisterInfos() []*v1.DeviceConfigRegisterInfo {
	if x != nil {
		return x.DeviceConfigRegisterInfos
	}
	return nil
}

func (x *RegisterRequest) GetDeviceStateRegisterInfos() []*v1.DeviceStateRegisterInfo {
	if x != nil {
		return x.DeviceStateRegisterInfos
	}
	return nil
}

// 注册响应
type RegisterReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 表示是否注册成功
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	// 用户的token
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *RegisterReply) Reset() {
	*x = RegisterReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_serviceCenter_v1_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterReply) ProtoMessage() {}

func (x *RegisterReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_serviceCenter_v1_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterReply.ProtoReflect.Descriptor instead.
func (*RegisterReply) Descriptor() ([]byte, []int) {
	return file_api_serviceCenter_v1_user_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *RegisterReply) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// 登录响应
type LoginReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Token   string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *LoginReply) Reset() {
	*x = LoginReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_serviceCenter_v1_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginReply) ProtoMessage() {}

func (x *LoginReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_serviceCenter_v1_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginReply.ProtoReflect.Descriptor instead.
func (*LoginReply) Descriptor() ([]byte, []int) {
	return file_api_serviceCenter_v1_user_proto_rawDescGZIP(), []int{2}
}

func (x *LoginReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *LoginReply) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// 注销响应
type UnregisterReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *UnregisterReply) Reset() {
	*x = UnregisterReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_serviceCenter_v1_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnregisterReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterReply) ProtoMessage() {}

func (x *UnregisterReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_serviceCenter_v1_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterReply.ProtoReflect.Descriptor instead.
func (*UnregisterReply) Descriptor() ([]byte, []int) {
	return file_api_serviceCenter_v1_user_proto_rawDescGZIP(), []int{3}
}

func (x *UnregisterReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_api_serviceCenter_v1_user_proto protoreflect.FileDescriptor

var file_api_serviceCenter_v1_user_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x65, 0x6e,
	0x74, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x14, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x65,
	0x6e, 0x74, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x75, 0x74, 0x69, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x02, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x75, 0x74,
	0x69, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72,
	0x12, 0x66, 0x0a, 0x1c, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x75, 0x74, 0x69,
	0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x19, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x12, 0x63, 0x0a, 0x1b, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x18, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x22, 0x3f, 0x0a,
	0x0d, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3c,
	0x0a, 0x0a, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2b, 0x0a, 0x0f,
	0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0x97, 0x02, 0x0a, 0x04, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x69, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x25,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x65, 0x6e, 0x74,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x43, 0x65, 0x6e, 0x74, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x0b, 0x22, 0x06, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x3a, 0x01, 0x2a, 0x12, 0x4c, 0x0a,
	0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x75, 0x74, 0x69,
	0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x20, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x65, 0x6e, 0x74, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x0e, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x08, 0x12, 0x06, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x56, 0x0a, 0x0a, 0x55,
	0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x75, 0x74, 0x69, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x25, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x65, 0x6e, 0x74, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x0e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x08, 0x2a, 0x06, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x42, 0x65, 0x0a, 0x27, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x69, 0x74, 0x65, 0x65,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x79, 0x75, 0x73, 0x69, 0x72, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2d, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x50, 0x01,
	0x5a, 0x38, 0x67, 0x69, 0x74, 0x65, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x79, 0x75,
	0x73, 0x69, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2d, 0x63, 0x65, 0x6e, 0x74,
	0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_api_serviceCenter_v1_user_proto_rawDescOnce sync.Once
	file_api_serviceCenter_v1_user_proto_rawDescData = file_api_serviceCenter_v1_user_proto_rawDesc
)

func file_api_serviceCenter_v1_user_proto_rawDescGZIP() []byte {
	file_api_serviceCenter_v1_user_proto_rawDescOnce.Do(func() {
		file_api_serviceCenter_v1_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_serviceCenter_v1_user_proto_rawDescData)
	})
	return file_api_serviceCenter_v1_user_proto_rawDescData
}

var file_api_serviceCenter_v1_user_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_serviceCenter_v1_user_proto_goTypes = []interface{}{
	(*RegisterRequest)(nil),             // 0: api.serviceCentre.v1.RegisterRequest
	(*RegisterReply)(nil),               // 1: api.serviceCentre.v1.RegisterReply
	(*LoginReply)(nil),                  // 2: api.serviceCentre.v1.LoginReply
	(*UnregisterReply)(nil),             // 3: api.serviceCentre.v1.UnregisterReply
	(*v1.User)(nil),                     // 4: api.util.v1.User
	(*v1.DeviceConfigRegisterInfo)(nil), // 5: api.util.v1.DeviceConfigRegisterInfo
	(*v1.DeviceStateRegisterInfo)(nil),  // 6: api.util.v1.DeviceStateRegisterInfo
}
var file_api_serviceCenter_v1_user_proto_depIdxs = []int32{
	4, // 0: api.serviceCentre.v1.RegisterRequest.user:type_name -> api.util.v1.User
	5, // 1: api.serviceCentre.v1.RegisterRequest.device_config_register_infos:type_name -> api.util.v1.DeviceConfigRegisterInfo
	6, // 2: api.serviceCentre.v1.RegisterRequest.device_state_register_infos:type_name -> api.util.v1.DeviceStateRegisterInfo
	0, // 3: api.serviceCentre.v1.User.Register:input_type -> api.serviceCentre.v1.RegisterRequest
	4, // 4: api.serviceCentre.v1.User.Login:input_type -> api.util.v1.User
	4, // 5: api.serviceCentre.v1.User.Unregister:input_type -> api.util.v1.User
	1, // 6: api.serviceCentre.v1.User.Register:output_type -> api.serviceCentre.v1.RegisterReply
	2, // 7: api.serviceCentre.v1.User.Login:output_type -> api.serviceCentre.v1.LoginReply
	3, // 8: api.serviceCentre.v1.User.Unregister:output_type -> api.serviceCentre.v1.UnregisterReply
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_serviceCenter_v1_user_proto_init() }
func file_api_serviceCenter_v1_user_proto_init() {
	if File_api_serviceCenter_v1_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_serviceCenter_v1_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRequest); i {
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
		file_api_serviceCenter_v1_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterReply); i {
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
		file_api_serviceCenter_v1_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginReply); i {
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
		file_api_serviceCenter_v1_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnregisterReply); i {
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
			RawDescriptor: file_api_serviceCenter_v1_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_serviceCenter_v1_user_proto_goTypes,
		DependencyIndexes: file_api_serviceCenter_v1_user_proto_depIdxs,
		MessageInfos:      file_api_serviceCenter_v1_user_proto_msgTypes,
	}.Build()
	File_api_serviceCenter_v1_user_proto = out.File
	file_api_serviceCenter_v1_user_proto_rawDesc = nil
	file_api_serviceCenter_v1_user_proto_goTypes = nil
	file_api_serviceCenter_v1_user_proto_depIdxs = nil
}

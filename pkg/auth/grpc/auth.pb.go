// Command to generate grpc
// $ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\auth.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: auth.proto

package auth

import (
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

type GetAdminTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAdminTokenRequest) Reset() {
	*x = GetAdminTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAdminTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAdminTokenRequest) ProtoMessage() {}

func (x *GetAdminTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAdminTokenRequest.ProtoReflect.Descriptor instead.
func (*GetAdminTokenRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{0}
}

type GetAdminTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *GetAdminTokenResponse) Reset() {
	*x = GetAdminTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAdminTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAdminTokenResponse) ProtoMessage() {}

func (x *GetAdminTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAdminTokenResponse.ProtoReflect.Descriptor instead.
func (*GetAdminTokenResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{1}
}

func (x *GetAdminTokenResponse) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

type GetUserTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetUserTokenRequest) Reset() {
	*x = GetUserTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserTokenRequest) ProtoMessage() {}

func (x *GetUserTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserTokenRequest.ProtoReflect.Descriptor instead.
func (*GetUserTokenRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserTokenRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetUserTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *GetUserTokenResponse) Reset() {
	*x = GetUserTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserTokenResponse) ProtoMessage() {}

func (x *GetUserTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserTokenResponse.ProtoReflect.Descriptor instead.
func (*GetUserTokenResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserTokenResponse) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

type GetThingTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ThingId string `protobuf:"bytes,1,opt,name=thingId,proto3" json:"thingId,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetThingTokenRequest) Reset() {
	*x = GetThingTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetThingTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetThingTokenRequest) ProtoMessage() {}

func (x *GetThingTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetThingTokenRequest.ProtoReflect.Descriptor instead.
func (*GetThingTokenRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{4}
}

func (x *GetThingTokenRequest) GetThingId() string {
	if x != nil {
		return x.ThingId
	}
	return ""
}

func (x *GetThingTokenRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetThingTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *GetThingTokenResponse) Reset() {
	*x = GetThingTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetThingTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetThingTokenResponse) ProtoMessage() {}

func (x *GetThingTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetThingTokenResponse.ProtoReflect.Descriptor instead.
func (*GetThingTokenResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{5}
}

func (x *GetThingTokenResponse) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

type GetClaimsUserTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *GetClaimsUserTokenRequest) Reset() {
	*x = GetClaimsUserTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClaimsUserTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClaimsUserTokenRequest) ProtoMessage() {}

func (x *GetClaimsUserTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClaimsUserTokenRequest.ProtoReflect.Descriptor instead.
func (*GetClaimsUserTokenRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{6}
}

func (x *GetClaimsUserTokenRequest) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

type GetClaimsUserTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role      string `protobuf:"bytes,1,opt,name=Role,proto3" json:"Role,omitempty"`
	Subject   string `protobuf:"bytes,2,opt,name=Subject,proto3" json:"Subject,omitempty"`
	IssuedAt  int64  `protobuf:"varint,3,opt,name=IssuedAt,proto3" json:"IssuedAt,omitempty"`
	ExpiresAt int64  `protobuf:"varint,4,opt,name=ExpiresAt,proto3" json:"ExpiresAt,omitempty"`
	Issuer    string `protobuf:"bytes,5,opt,name=Issuer,proto3" json:"Issuer,omitempty"`
}

func (x *GetClaimsUserTokenResponse) Reset() {
	*x = GetClaimsUserTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClaimsUserTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClaimsUserTokenResponse) ProtoMessage() {}

func (x *GetClaimsUserTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClaimsUserTokenResponse.ProtoReflect.Descriptor instead.
func (*GetClaimsUserTokenResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{7}
}

func (x *GetClaimsUserTokenResponse) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *GetClaimsUserTokenResponse) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *GetClaimsUserTokenResponse) GetIssuedAt() int64 {
	if x != nil {
		return x.IssuedAt
	}
	return 0
}

func (x *GetClaimsUserTokenResponse) GetExpiresAt() int64 {
	if x != nil {
		return x.ExpiresAt
	}
	return 0
}

func (x *GetClaimsUserTokenResponse) GetIssuer() string {
	if x != nil {
		return x.Issuer
	}
	return ""
}

type GetClaimsThingTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
}

func (x *GetClaimsThingTokenRequest) Reset() {
	*x = GetClaimsThingTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClaimsThingTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClaimsThingTokenRequest) ProtoMessage() {}

func (x *GetClaimsThingTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClaimsThingTokenRequest.ProtoReflect.Descriptor instead.
func (*GetClaimsThingTokenRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{8}
}

func (x *GetClaimsThingTokenRequest) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

type GetClaimsThingTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Subject   string `protobuf:"bytes,2,opt,name=Subject,proto3" json:"Subject,omitempty"`
	IssuedAt  int64  `protobuf:"varint,3,opt,name=IssuedAt,proto3" json:"IssuedAt,omitempty"`
	ExpiresAt int64  `protobuf:"varint,4,opt,name=ExpiresAt,proto3" json:"ExpiresAt,omitempty"`
	Issuer    string `protobuf:"bytes,5,opt,name=Issuer,proto3" json:"Issuer,omitempty"`
}

func (x *GetClaimsThingTokenResponse) Reset() {
	*x = GetClaimsThingTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClaimsThingTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClaimsThingTokenResponse) ProtoMessage() {}

func (x *GetClaimsThingTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClaimsThingTokenResponse.ProtoReflect.Descriptor instead.
func (*GetClaimsThingTokenResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{9}
}

func (x *GetClaimsThingTokenResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetClaimsThingTokenResponse) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *GetClaimsThingTokenResponse) GetIssuedAt() int64 {
	if x != nil {
		return x.IssuedAt
	}
	return 0
}

func (x *GetClaimsThingTokenResponse) GetExpiresAt() int64 {
	if x != nil {
		return x.ExpiresAt
	}
	return 0
}

func (x *GetClaimsThingTokenResponse) GetIssuer() string {
	if x != nil {
		return x.Issuer
	}
	return ""
}

var File_auth_proto protoreflect.FileDescriptor

var file_auth_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x61, 0x75,
	0x74, 0x68, 0x22, 0x16, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x29, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6a, 0x77, 0x74, 0x22, 0x2d, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x22, 0x28, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x22, 0x48,
	0x0a, 0x14, 0x47, 0x65, 0x74, 0x54, 0x68, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x29, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x54,
	0x68, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6a, 0x77, 0x74, 0x22, 0x2d, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73,
	0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6a,
	0x77, 0x74, 0x22, 0x9c, 0x01, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73,
	0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x45,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x49, 0x73, 0x73,
	0x75, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x49, 0x73, 0x73, 0x75, 0x65,
	0x72, 0x22, 0x2e, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x54, 0x68,
	0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6a, 0x77,
	0x74, 0x22, 0xa1, 0x01, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x54,
	0x68, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64, 0x41, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x1c, 0x0a, 0x09, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x49,
	0x73, 0x73, 0x75, 0x65, 0x72, 0x32, 0xa0, 0x03, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x4a,
	0x0a, 0x0d, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x1a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x54, 0x68, 0x69, 0x6e, 0x67, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x54,
	0x68, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1b, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x68, 0x69, 0x6e, 0x67,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x59, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x55, 0x73, 0x65, 0x72,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1f, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74,
	0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65,
	0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5c, 0x0a, 0x13, 0x47, 0x65,
	0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x54, 0x68, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x20, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69,
	0x6d, 0x73, 0x54, 0x68, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c,
	0x61, 0x69, 0x6d, 0x73, 0x54, 0x68, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x69, 0x6e, 0x63,
	0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x2f, 0x42, 0x61, 0x72, 0x69, 0x6f, 0x74, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x75, 0x74, 0x68,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_proto_rawDescOnce sync.Once
	file_auth_proto_rawDescData = file_auth_proto_rawDesc
)

func file_auth_proto_rawDescGZIP() []byte {
	file_auth_proto_rawDescOnce.Do(func() {
		file_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_proto_rawDescData)
	})
	return file_auth_proto_rawDescData
}

var file_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_auth_proto_goTypes = []interface{}{
	(*GetAdminTokenRequest)(nil),        // 0: auth.GetAdminTokenRequest
	(*GetAdminTokenResponse)(nil),       // 1: auth.GetAdminTokenResponse
	(*GetUserTokenRequest)(nil),         // 2: auth.GetUserTokenRequest
	(*GetUserTokenResponse)(nil),        // 3: auth.GetUserTokenResponse
	(*GetThingTokenRequest)(nil),        // 4: auth.GetThingTokenRequest
	(*GetThingTokenResponse)(nil),       // 5: auth.GetThingTokenResponse
	(*GetClaimsUserTokenRequest)(nil),   // 6: auth.GetClaimsUserTokenRequest
	(*GetClaimsUserTokenResponse)(nil),  // 7: auth.GetClaimsUserTokenResponse
	(*GetClaimsThingTokenRequest)(nil),  // 8: auth.GetClaimsThingTokenRequest
	(*GetClaimsThingTokenResponse)(nil), // 9: auth.GetClaimsThingTokenResponse
}
var file_auth_proto_depIdxs = []int32{
	0, // 0: auth.Auth.GetAdminToken:input_type -> auth.GetAdminTokenRequest
	2, // 1: auth.Auth.GetUserToken:input_type -> auth.GetUserTokenRequest
	4, // 2: auth.Auth.GetThingToken:input_type -> auth.GetThingTokenRequest
	6, // 3: auth.Auth.GetClaimsUserToken:input_type -> auth.GetClaimsUserTokenRequest
	8, // 4: auth.Auth.GetClaimsThingToken:input_type -> auth.GetClaimsThingTokenRequest
	1, // 5: auth.Auth.GetAdminToken:output_type -> auth.GetAdminTokenResponse
	3, // 6: auth.Auth.GetUserToken:output_type -> auth.GetUserTokenResponse
	5, // 7: auth.Auth.GetThingToken:output_type -> auth.GetThingTokenResponse
	7, // 8: auth.Auth.GetClaimsUserToken:output_type -> auth.GetClaimsUserTokenResponse
	9, // 9: auth.Auth.GetClaimsThingToken:output_type -> auth.GetClaimsThingTokenResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_auth_proto_init() }
func file_auth_proto_init() {
	if File_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAdminTokenRequest); i {
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
		file_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAdminTokenResponse); i {
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
		file_auth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserTokenRequest); i {
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
		file_auth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserTokenResponse); i {
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
		file_auth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetThingTokenRequest); i {
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
		file_auth_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetThingTokenResponse); i {
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
		file_auth_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClaimsUserTokenRequest); i {
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
		file_auth_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClaimsUserTokenResponse); i {
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
		file_auth_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClaimsThingTokenRequest); i {
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
		file_auth_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClaimsThingTokenResponse); i {
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
			RawDescriptor: file_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_proto_goTypes,
		DependencyIndexes: file_auth_proto_depIdxs,
		MessageInfos:      file_auth_proto_msgTypes,
	}.Build()
	File_auth_proto = out.File
	file_auth_proto_rawDesc = nil
	file_auth_proto_goTypes = nil
	file_auth_proto_depIdxs = nil
}

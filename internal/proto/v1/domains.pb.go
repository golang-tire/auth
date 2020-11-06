// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.8.0
// source: api/proto/v1/domains.proto

package auth

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Domain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid      string               `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name      string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Enable    bool                 `protobuf:"varint,3,opt,name=enable,proto3" json:"enable,omitempty"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Domain) Reset() {
	*x = Domain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_v1_domains_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Domain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Domain) ProtoMessage() {}

func (x *Domain) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_domains_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Domain.ProtoReflect.Descriptor instead.
func (*Domain) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_domains_proto_rawDescGZIP(), []int{0}
}

func (x *Domain) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Domain) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Domain) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

func (x *Domain) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Domain) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type ListDomainsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit  int64 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset int64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *ListDomainsRequest) Reset() {
	*x = ListDomainsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_v1_domains_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDomainsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDomainsRequest) ProtoMessage() {}

func (x *ListDomainsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_domains_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDomainsRequest.ProtoReflect.Descriptor instead.
func (*ListDomainsRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_domains_proto_rawDescGZIP(), []int{1}
}

func (x *ListDomainsRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListDomainsRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type ListDomainsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Domains    []*Domain `protobuf:"bytes,1,rep,name=domains,proto3" json:"domains,omitempty"`
	TotalCount int64     `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	Limit      int64     `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset     int64     `protobuf:"varint,4,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *ListDomainsResponse) Reset() {
	*x = ListDomainsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_v1_domains_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDomainsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDomainsResponse) ProtoMessage() {}

func (x *ListDomainsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_domains_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDomainsResponse.ProtoReflect.Descriptor instead.
func (*ListDomainsResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_domains_proto_rawDescGZIP(), []int{2}
}

func (x *ListDomainsResponse) GetDomains() []*Domain {
	if x != nil {
		return x.Domains
	}
	return nil
}

func (x *ListDomainsResponse) GetTotalCount() int64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *ListDomainsResponse) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListDomainsResponse) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type GetDomainRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *GetDomainRequest) Reset() {
	*x = GetDomainRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_v1_domains_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDomainRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDomainRequest) ProtoMessage() {}

func (x *GetDomainRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_domains_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDomainRequest.ProtoReflect.Descriptor instead.
func (*GetDomainRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_domains_proto_rawDescGZIP(), []int{3}
}

func (x *GetDomainRequest) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type CreateDomainRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Enable bool   `protobuf:"varint,2,opt,name=enable,proto3" json:"enable,omitempty"`
}

func (x *CreateDomainRequest) Reset() {
	*x = CreateDomainRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_v1_domains_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateDomainRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDomainRequest) ProtoMessage() {}

func (x *CreateDomainRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_domains_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDomainRequest.ProtoReflect.Descriptor instead.
func (*CreateDomainRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_domains_proto_rawDescGZIP(), []int{4}
}

func (x *CreateDomainRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateDomainRequest) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

type UpdateDomainRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid   string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Enable bool   `protobuf:"varint,3,opt,name=enable,proto3" json:"enable,omitempty"`
}

func (x *UpdateDomainRequest) Reset() {
	*x = UpdateDomainRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_v1_domains_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateDomainRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDomainRequest) ProtoMessage() {}

func (x *UpdateDomainRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_domains_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDomainRequest.ProtoReflect.Descriptor instead.
func (*UpdateDomainRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_domains_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateDomainRequest) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *UpdateDomainRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateDomainRequest) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

type DeleteDomainRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *DeleteDomainRequest) Reset() {
	*x = DeleteDomainRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_v1_domains_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteDomainRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDomainRequest) ProtoMessage() {}

func (x *DeleteDomainRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_domains_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDomainRequest.ProtoReflect.Descriptor instead.
func (*DeleteDomainRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_domains_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteDomainRequest) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

var File_api_proto_v1_domains_proto protoreflect.FileDescriptor

var file_api_proto_v1_domains_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x75,
	0x74, 0x68, 0x56, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xbe, 0x01, 0x0a, 0x06, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x22, 0x42, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6f,
	0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x8e, 0x01, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a,
	0x07, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x61, 0x75, 0x74, 0x68, 0x56, 0x31, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x52, 0x07,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x26, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x22, 0x41,
	0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x61,
	0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c,
	0x65, 0x22, 0x55, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x29, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75,
	0x75, 0x69, 0x64, 0x32, 0xd1, 0x03, 0x0a, 0x0d, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x73, 0x12, 0x1a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x56, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1b, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x56, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x13, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x73, 0x12, 0x51, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12,
	0x18, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x56, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x56, 0x31, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x14, 0x12, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x2f, 0x7b,
	0x75, 0x75, 0x69, 0x64, 0x7d, 0x12, 0x53, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1b, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x56, 0x31, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x56, 0x31, 0x2e, 0x44, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x22, 0x0b, 0x2f, 0x76, 0x31, 0x2f,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x3a, 0x01, 0x2a, 0x12, 0x5a, 0x0a, 0x0c, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1b, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x56, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x56, 0x31,
	0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x1a,
	0x12, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x2f, 0x7b, 0x75, 0x75,
	0x69, 0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x5f, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1b, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x56, 0x31, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1a, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x14, 0x2a, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73,
	0x2f, 0x7b, 0x75, 0x75, 0x69, 0x64, 0x7d, 0x42, 0x18, 0x5a, 0x16, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x75, 0x74,
	0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_v1_domains_proto_rawDescOnce sync.Once
	file_api_proto_v1_domains_proto_rawDescData = file_api_proto_v1_domains_proto_rawDesc
)

func file_api_proto_v1_domains_proto_rawDescGZIP() []byte {
	file_api_proto_v1_domains_proto_rawDescOnce.Do(func() {
		file_api_proto_v1_domains_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_v1_domains_proto_rawDescData)
	})
	return file_api_proto_v1_domains_proto_rawDescData
}

var file_api_proto_v1_domains_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_proto_v1_domains_proto_goTypes = []interface{}{
	(*Domain)(nil),              // 0: authV1.Domain
	(*ListDomainsRequest)(nil),  // 1: authV1.ListDomainsRequest
	(*ListDomainsResponse)(nil), // 2: authV1.ListDomainsResponse
	(*GetDomainRequest)(nil),    // 3: authV1.GetDomainRequest
	(*CreateDomainRequest)(nil), // 4: authV1.CreateDomainRequest
	(*UpdateDomainRequest)(nil), // 5: authV1.UpdateDomainRequest
	(*DeleteDomainRequest)(nil), // 6: authV1.DeleteDomainRequest
	(*timestamp.Timestamp)(nil), // 7: google.protobuf.Timestamp
	(*empty.Empty)(nil),         // 8: google.protobuf.Empty
}
var file_api_proto_v1_domains_proto_depIdxs = []int32{
	7, // 0: authV1.Domain.created_at:type_name -> google.protobuf.Timestamp
	7, // 1: authV1.Domain.updated_at:type_name -> google.protobuf.Timestamp
	0, // 2: authV1.ListDomainsResponse.domains:type_name -> authV1.Domain
	1, // 3: authV1.DomainService.ListDomains:input_type -> authV1.ListDomainsRequest
	3, // 4: authV1.DomainService.GetDomain:input_type -> authV1.GetDomainRequest
	4, // 5: authV1.DomainService.CreateDomain:input_type -> authV1.CreateDomainRequest
	5, // 6: authV1.DomainService.UpdateDomain:input_type -> authV1.UpdateDomainRequest
	6, // 7: authV1.DomainService.DeleteDomain:input_type -> authV1.DeleteDomainRequest
	2, // 8: authV1.DomainService.ListDomains:output_type -> authV1.ListDomainsResponse
	0, // 9: authV1.DomainService.GetDomain:output_type -> authV1.Domain
	0, // 10: authV1.DomainService.CreateDomain:output_type -> authV1.Domain
	0, // 11: authV1.DomainService.UpdateDomain:output_type -> authV1.Domain
	8, // 12: authV1.DomainService.DeleteDomain:output_type -> google.protobuf.Empty
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_proto_v1_domains_proto_init() }
func file_api_proto_v1_domains_proto_init() {
	if File_api_proto_v1_domains_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_v1_domains_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Domain); i {
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
		file_api_proto_v1_domains_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDomainsRequest); i {
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
		file_api_proto_v1_domains_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDomainsResponse); i {
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
		file_api_proto_v1_domains_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDomainRequest); i {
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
		file_api_proto_v1_domains_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateDomainRequest); i {
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
		file_api_proto_v1_domains_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateDomainRequest); i {
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
		file_api_proto_v1_domains_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteDomainRequest); i {
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
			RawDescriptor: file_api_proto_v1_domains_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_v1_domains_proto_goTypes,
		DependencyIndexes: file_api_proto_v1_domains_proto_depIdxs,
		MessageInfos:      file_api_proto_v1_domains_proto_msgTypes,
	}.Build()
	File_api_proto_v1_domains_proto = out.File
	file_api_proto_v1_domains_proto_rawDesc = nil
	file_api_proto_v1_domains_proto_goTypes = nil
	file_api_proto_v1_domains_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DomainServiceClient is the client API for DomainService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DomainServiceClient interface {
	// List Domains
	ListDomains(ctx context.Context, in *ListDomainsRequest, opts ...grpc.CallOption) (*ListDomainsResponse, error)
	// Get Domain
	GetDomain(ctx context.Context, in *GetDomainRequest, opts ...grpc.CallOption) (*Domain, error)
	// Create Domain object request
	CreateDomain(ctx context.Context, in *CreateDomainRequest, opts ...grpc.CallOption) (*Domain, error)
	// Update Domain object request
	UpdateDomain(ctx context.Context, in *UpdateDomainRequest, opts ...grpc.CallOption) (*Domain, error)
	// Delete Domain object request
	DeleteDomain(ctx context.Context, in *DeleteDomainRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type domainServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDomainServiceClient(cc grpc.ClientConnInterface) DomainServiceClient {
	return &domainServiceClient{cc}
}

func (c *domainServiceClient) ListDomains(ctx context.Context, in *ListDomainsRequest, opts ...grpc.CallOption) (*ListDomainsResponse, error) {
	out := new(ListDomainsResponse)
	err := c.cc.Invoke(ctx, "/authV1.DomainService/ListDomains", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *domainServiceClient) GetDomain(ctx context.Context, in *GetDomainRequest, opts ...grpc.CallOption) (*Domain, error) {
	out := new(Domain)
	err := c.cc.Invoke(ctx, "/authV1.DomainService/GetDomain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *domainServiceClient) CreateDomain(ctx context.Context, in *CreateDomainRequest, opts ...grpc.CallOption) (*Domain, error) {
	out := new(Domain)
	err := c.cc.Invoke(ctx, "/authV1.DomainService/CreateDomain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *domainServiceClient) UpdateDomain(ctx context.Context, in *UpdateDomainRequest, opts ...grpc.CallOption) (*Domain, error) {
	out := new(Domain)
	err := c.cc.Invoke(ctx, "/authV1.DomainService/UpdateDomain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *domainServiceClient) DeleteDomain(ctx context.Context, in *DeleteDomainRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/authV1.DomainService/DeleteDomain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DomainServiceServer is the server API for DomainService service.
type DomainServiceServer interface {
	// List Domains
	ListDomains(context.Context, *ListDomainsRequest) (*ListDomainsResponse, error)
	// Get Domain
	GetDomain(context.Context, *GetDomainRequest) (*Domain, error)
	// Create Domain object request
	CreateDomain(context.Context, *CreateDomainRequest) (*Domain, error)
	// Update Domain object request
	UpdateDomain(context.Context, *UpdateDomainRequest) (*Domain, error)
	// Delete Domain object request
	DeleteDomain(context.Context, *DeleteDomainRequest) (*empty.Empty, error)
}

// UnimplementedDomainServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDomainServiceServer struct {
}

func (*UnimplementedDomainServiceServer) ListDomains(context.Context, *ListDomainsRequest) (*ListDomainsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDomains not implemented")
}
func (*UnimplementedDomainServiceServer) GetDomain(context.Context, *GetDomainRequest) (*Domain, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDomain not implemented")
}
func (*UnimplementedDomainServiceServer) CreateDomain(context.Context, *CreateDomainRequest) (*Domain, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDomain not implemented")
}
func (*UnimplementedDomainServiceServer) UpdateDomain(context.Context, *UpdateDomainRequest) (*Domain, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDomain not implemented")
}
func (*UnimplementedDomainServiceServer) DeleteDomain(context.Context, *DeleteDomainRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDomain not implemented")
}

func RegisterDomainServiceServer(s *grpc.Server, srv DomainServiceServer) {
	s.RegisterService(&_DomainService_serviceDesc, srv)
}

func _DomainService_ListDomains_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDomainsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainServiceServer).ListDomains(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authV1.DomainService/ListDomains",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainServiceServer).ListDomains(ctx, req.(*ListDomainsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DomainService_GetDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDomainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainServiceServer).GetDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authV1.DomainService/GetDomain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainServiceServer).GetDomain(ctx, req.(*GetDomainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DomainService_CreateDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDomainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainServiceServer).CreateDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authV1.DomainService/CreateDomain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainServiceServer).CreateDomain(ctx, req.(*CreateDomainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DomainService_UpdateDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDomainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainServiceServer).UpdateDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authV1.DomainService/UpdateDomain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainServiceServer).UpdateDomain(ctx, req.(*UpdateDomainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DomainService_DeleteDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDomainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainServiceServer).DeleteDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authV1.DomainService/DeleteDomain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainServiceServer).DeleteDomain(ctx, req.(*DeleteDomainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DomainService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "authV1.DomainService",
	HandlerType: (*DomainServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDomains",
			Handler:    _DomainService_ListDomains_Handler,
		},
		{
			MethodName: "GetDomain",
			Handler:    _DomainService_GetDomain_Handler,
		},
		{
			MethodName: "CreateDomain",
			Handler:    _DomainService_CreateDomain_Handler,
		},
		{
			MethodName: "UpdateDomain",
			Handler:    _DomainService_UpdateDomain_Handler,
		},
		{
			MethodName: "DeleteDomain",
			Handler:    _DomainService_DeleteDomain_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/v1/domains.proto",
}

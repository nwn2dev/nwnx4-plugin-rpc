// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v4.22.2
// source: scorco.proto

package proto

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

type SCORCOGetGFFSizeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SVarName string `protobuf:"bytes,1,opt,name=sVarName,proto3" json:"sVarName,omitempty"`
}

func (x *SCORCOGetGFFSizeRequest) Reset() {
	*x = SCORCOGetGFFSizeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scorco_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SCORCOGetGFFSizeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCORCOGetGFFSizeRequest) ProtoMessage() {}

func (x *SCORCOGetGFFSizeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scorco_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCORCOGetGFFSizeRequest.ProtoReflect.Descriptor instead.
func (*SCORCOGetGFFSizeRequest) Descriptor() ([]byte, []int) {
	return file_scorco_proto_rawDescGZIP(), []int{0}
}

func (x *SCORCOGetGFFSizeRequest) GetSVarName() string {
	if x != nil {
		return x.SVarName
	}
	return ""
}

type SCORCOGetGFFSizeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Size uint32 `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *SCORCOGetGFFSizeResponse) Reset() {
	*x = SCORCOGetGFFSizeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scorco_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SCORCOGetGFFSizeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCORCOGetGFFSizeResponse) ProtoMessage() {}

func (x *SCORCOGetGFFSizeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scorco_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCORCOGetGFFSizeResponse.ProtoReflect.Descriptor instead.
func (*SCORCOGetGFFSizeResponse) Descriptor() ([]byte, []int) {
	return file_scorco_proto_rawDescGZIP(), []int{1}
}

func (x *SCORCOGetGFFSizeResponse) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

type SCORCOGetGFFRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SVarName string `protobuf:"bytes,1,opt,name=sVarName,proto3" json:"sVarName,omitempty"`
}

func (x *SCORCOGetGFFRequest) Reset() {
	*x = SCORCOGetGFFRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scorco_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SCORCOGetGFFRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCORCOGetGFFRequest) ProtoMessage() {}

func (x *SCORCOGetGFFRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scorco_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCORCOGetGFFRequest.ProtoReflect.Descriptor instead.
func (*SCORCOGetGFFRequest) Descriptor() ([]byte, []int) {
	return file_scorco_proto_rawDescGZIP(), []int{2}
}

func (x *SCORCOGetGFFRequest) GetSVarName() string {
	if x != nil {
		return x.SVarName
	}
	return ""
}

type SCORCOGetGFFResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GffData []byte `protobuf:"bytes,1,opt,name=gffData,proto3" json:"gffData,omitempty"`
}

func (x *SCORCOGetGFFResponse) Reset() {
	*x = SCORCOGetGFFResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scorco_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SCORCOGetGFFResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCORCOGetGFFResponse) ProtoMessage() {}

func (x *SCORCOGetGFFResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scorco_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCORCOGetGFFResponse.ProtoReflect.Descriptor instead.
func (*SCORCOGetGFFResponse) Descriptor() ([]byte, []int) {
	return file_scorco_proto_rawDescGZIP(), []int{3}
}

func (x *SCORCOGetGFFResponse) GetGffData() []byte {
	if x != nil {
		return x.GffData
	}
	return nil
}

type SCORCOSetGFFRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SVarName    string `protobuf:"bytes,1,opt,name=sVarName,proto3" json:"sVarName,omitempty"`
	GffData     []byte `protobuf:"bytes,2,opt,name=gffData,proto3" json:"gffData,omitempty"`
	GffDataSize uint32 `protobuf:"varint,3,opt,name=gffDataSize,proto3" json:"gffDataSize,omitempty"`
}

func (x *SCORCOSetGFFRequest) Reset() {
	*x = SCORCOSetGFFRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scorco_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SCORCOSetGFFRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCORCOSetGFFRequest) ProtoMessage() {}

func (x *SCORCOSetGFFRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scorco_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCORCOSetGFFRequest.ProtoReflect.Descriptor instead.
func (*SCORCOSetGFFRequest) Descriptor() ([]byte, []int) {
	return file_scorco_proto_rawDescGZIP(), []int{4}
}

func (x *SCORCOSetGFFRequest) GetSVarName() string {
	if x != nil {
		return x.SVarName
	}
	return ""
}

func (x *SCORCOSetGFFRequest) GetGffData() []byte {
	if x != nil {
		return x.GffData
	}
	return nil
}

func (x *SCORCOSetGFFRequest) GetGffDataSize() uint32 {
	if x != nil {
		return x.GffDataSize
	}
	return 0
}

type SCORCOSetGFFResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SCORCOSetGFFResponse) Reset() {
	*x = SCORCOSetGFFResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scorco_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SCORCOSetGFFResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCORCOSetGFFResponse) ProtoMessage() {}

func (x *SCORCOSetGFFResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scorco_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCORCOSetGFFResponse.ProtoReflect.Descriptor instead.
func (*SCORCOSetGFFResponse) Descriptor() ([]byte, []int) {
	return file_scorco_proto_rawDescGZIP(), []int{5}
}

var File_scorco_proto protoreflect.FileDescriptor

var file_scorco_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x63, 0x6f, 0x72, 0x63, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x4e, 0x57, 0x4e, 0x58, 0x34, 0x2e, 0x52, 0x50, 0x43, 0x22, 0x35, 0x0a, 0x17, 0x53, 0x43, 0x4f,
	0x52, 0x43, 0x4f, 0x47, 0x65, 0x74, 0x47, 0x46, 0x46, 0x53, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x56, 0x61, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x56, 0x61, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x2e, 0x0a, 0x18, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65, 0x74, 0x47, 0x46, 0x46,
	0x53, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x22, 0x31, 0x0a, 0x13, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65, 0x74, 0x47, 0x46, 0x46,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x56, 0x61, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x56, 0x61, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x30, 0x0a, 0x14, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65, 0x74,
	0x47, 0x46, 0x46, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x67,
	0x66, 0x66, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x67, 0x66,
	0x66, 0x44, 0x61, 0x74, 0x61, 0x22, 0x6d, 0x0a, 0x13, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x53,
	0x65, 0x74, 0x47, 0x46, 0x46, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x73, 0x56, 0x61, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x73, 0x56, 0x61, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x66, 0x66, 0x44,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x67, 0x66, 0x66, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x20, 0x0a, 0x0b, 0x67, 0x66, 0x66, 0x44, 0x61, 0x74, 0x61, 0x53, 0x69, 0x7a,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x67, 0x66, 0x66, 0x44, 0x61, 0x74, 0x61,
	0x53, 0x69, 0x7a, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x53, 0x65,
	0x74, 0x47, 0x46, 0x46, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x8e, 0x02, 0x0a,
	0x0d, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b,
	0x0a, 0x10, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65, 0x74, 0x47, 0x46, 0x46, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x22, 0x2e, 0x4e, 0x57, 0x4e, 0x58, 0x34, 0x2e, 0x52, 0x50, 0x43, 0x2e, 0x53,
	0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65, 0x74, 0x47, 0x46, 0x46, 0x53, 0x69, 0x7a, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x4e, 0x57, 0x4e, 0x58, 0x34, 0x2e, 0x52,
	0x50, 0x43, 0x2e, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65, 0x74, 0x47, 0x46, 0x46, 0x53,
	0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0c, 0x53,
	0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65, 0x74, 0x47, 0x46, 0x46, 0x12, 0x1e, 0x2e, 0x4e, 0x57,
	0x4e, 0x58, 0x34, 0x2e, 0x52, 0x50, 0x43, 0x2e, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65,
	0x74, 0x47, 0x46, 0x46, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x4e, 0x57,
	0x4e, 0x58, 0x34, 0x2e, 0x52, 0x50, 0x43, 0x2e, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x47, 0x65,
	0x74, 0x47, 0x46, 0x46, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0c,
	0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x53, 0x65, 0x74, 0x47, 0x46, 0x46, 0x12, 0x1e, 0x2e, 0x4e,
	0x57, 0x4e, 0x58, 0x34, 0x2e, 0x52, 0x50, 0x43, 0x2e, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x53,
	0x65, 0x74, 0x47, 0x46, 0x46, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x4e,
	0x57, 0x4e, 0x58, 0x34, 0x2e, 0x52, 0x50, 0x43, 0x2e, 0x53, 0x43, 0x4f, 0x52, 0x43, 0x4f, 0x53,
	0x65, 0x74, 0x47, 0x46, 0x46, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0e, 0x5a,
	0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_scorco_proto_rawDescOnce sync.Once
	file_scorco_proto_rawDescData = file_scorco_proto_rawDesc
)

func file_scorco_proto_rawDescGZIP() []byte {
	file_scorco_proto_rawDescOnce.Do(func() {
		file_scorco_proto_rawDescData = protoimpl.X.CompressGZIP(file_scorco_proto_rawDescData)
	})
	return file_scorco_proto_rawDescData
}

var file_scorco_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_scorco_proto_goTypes = []interface{}{
	(*SCORCOGetGFFSizeRequest)(nil),  // 0: NWNX4.RPC.SCORCOGetGFFSizeRequest
	(*SCORCOGetGFFSizeResponse)(nil), // 1: NWNX4.RPC.SCORCOGetGFFSizeResponse
	(*SCORCOGetGFFRequest)(nil),      // 2: NWNX4.RPC.SCORCOGetGFFRequest
	(*SCORCOGetGFFResponse)(nil),     // 3: NWNX4.RPC.SCORCOGetGFFResponse
	(*SCORCOSetGFFRequest)(nil),      // 4: NWNX4.RPC.SCORCOSetGFFRequest
	(*SCORCOSetGFFResponse)(nil),     // 5: NWNX4.RPC.SCORCOSetGFFResponse
}
var file_scorco_proto_depIdxs = []int32{
	0, // 0: NWNX4.RPC.SCORCOService.SCORCOGetGFFSize:input_type -> NWNX4.RPC.SCORCOGetGFFSizeRequest
	2, // 1: NWNX4.RPC.SCORCOService.SCORCOGetGFF:input_type -> NWNX4.RPC.SCORCOGetGFFRequest
	4, // 2: NWNX4.RPC.SCORCOService.SCORCOSetGFF:input_type -> NWNX4.RPC.SCORCOSetGFFRequest
	1, // 3: NWNX4.RPC.SCORCOService.SCORCOGetGFFSize:output_type -> NWNX4.RPC.SCORCOGetGFFSizeResponse
	3, // 4: NWNX4.RPC.SCORCOService.SCORCOGetGFF:output_type -> NWNX4.RPC.SCORCOGetGFFResponse
	5, // 5: NWNX4.RPC.SCORCOService.SCORCOSetGFF:output_type -> NWNX4.RPC.SCORCOSetGFFResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_scorco_proto_init() }
func file_scorco_proto_init() {
	if File_scorco_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_scorco_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SCORCOGetGFFSizeRequest); i {
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
		file_scorco_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SCORCOGetGFFSizeResponse); i {
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
		file_scorco_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SCORCOGetGFFRequest); i {
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
		file_scorco_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SCORCOGetGFFResponse); i {
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
		file_scorco_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SCORCOSetGFFRequest); i {
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
		file_scorco_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SCORCOSetGFFResponse); i {
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
			RawDescriptor: file_scorco_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_scorco_proto_goTypes,
		DependencyIndexes: file_scorco_proto_depIdxs,
		MessageInfos:      file_scorco_proto_msgTypes,
	}.Build()
	File_scorco_proto = out.File
	file_scorco_proto_rawDesc = nil
	file_scorco_proto_goTypes = nil
	file_scorco_proto_depIdxs = nil
}
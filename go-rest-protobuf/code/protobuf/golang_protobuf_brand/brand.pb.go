// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: protobuf/brand.proto

package golang_protobuf_brand

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

type ProtoBrandRepo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Brands []*ProtoBrandRepo_ProtoBrand `protobuf:"bytes,1,rep,name=brands,proto3" json:"brands,omitempty"`
}

func (x *ProtoBrandRepo) Reset() {
	*x = ProtoBrandRepo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_brand_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoBrandRepo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoBrandRepo) ProtoMessage() {}

func (x *ProtoBrandRepo) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_brand_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoBrandRepo.ProtoReflect.Descriptor instead.
func (*ProtoBrandRepo) Descriptor() ([]byte, []int) {
	return file_protobuf_brand_proto_rawDescGZIP(), []int{0}
}

func (x *ProtoBrandRepo) GetBrands() []*ProtoBrandRepo_ProtoBrand {
	if x != nil {
		return x.Brands
	}
	return nil
}

type ProtoBrandRepo_ProtoBrand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID   uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Year uint32 `protobuf:"varint,3,opt,name=Year,proto3" json:"Year,omitempty"`
}

func (x *ProtoBrandRepo_ProtoBrand) Reset() {
	*x = ProtoBrandRepo_ProtoBrand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_brand_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoBrandRepo_ProtoBrand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoBrandRepo_ProtoBrand) ProtoMessage() {}

func (x *ProtoBrandRepo_ProtoBrand) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_brand_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoBrandRepo_ProtoBrand.ProtoReflect.Descriptor instead.
func (*ProtoBrandRepo_ProtoBrand) Descriptor() ([]byte, []int) {
	return file_protobuf_brand_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ProtoBrandRepo_ProtoBrand) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *ProtoBrandRepo_ProtoBrand) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProtoBrandRepo_ProtoBrand) GetYear() uint32 {
	if x != nil {
		return x.Year
	}
	return 0
}

var File_protobuf_brand_proto protoreflect.FileDescriptor

var file_protobuf_brand_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x62, 0x72, 0x61, 0x6e, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x5f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x22, 0xa0, 0x01,
	0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6f,
	0x12, 0x48, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x30, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x72,
	0x61, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x72, 0x61,
	0x6e, 0x64, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x73, 0x1a, 0x44, 0x0a, 0x0a, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x59, 0x65, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x59, 0x65, 0x61, 0x72,
	0x42, 0x18, 0x5a, 0x16, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_protobuf_brand_proto_rawDescOnce sync.Once
	file_protobuf_brand_proto_rawDescData = file_protobuf_brand_proto_rawDesc
)

func file_protobuf_brand_proto_rawDescGZIP() []byte {
	file_protobuf_brand_proto_rawDescOnce.Do(func() {
		file_protobuf_brand_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_brand_proto_rawDescData)
	})
	return file_protobuf_brand_proto_rawDescData
}

var file_protobuf_brand_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protobuf_brand_proto_goTypes = []interface{}{
	(*ProtoBrandRepo)(nil),            // 0: golang_protobuf_brand.ProtoBrandRepo
	(*ProtoBrandRepo_ProtoBrand)(nil), // 1: golang_protobuf_brand.ProtoBrandRepo.ProtoBrand
}
var file_protobuf_brand_proto_depIdxs = []int32{
	1, // 0: golang_protobuf_brand.ProtoBrandRepo.brands:type_name -> golang_protobuf_brand.ProtoBrandRepo.ProtoBrand
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protobuf_brand_proto_init() }
func file_protobuf_brand_proto_init() {
	if File_protobuf_brand_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_brand_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtoBrandRepo); i {
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
		file_protobuf_brand_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtoBrandRepo_ProtoBrand); i {
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
			RawDescriptor: file_protobuf_brand_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protobuf_brand_proto_goTypes,
		DependencyIndexes: file_protobuf_brand_proto_depIdxs,
		MessageInfos:      file_protobuf_brand_proto_msgTypes,
	}.Build()
	File_protobuf_brand_proto = out.File
	file_protobuf_brand_proto_rawDesc = nil
	file_protobuf_brand_proto_goTypes = nil
	file_protobuf_brand_proto_depIdxs = nil
}

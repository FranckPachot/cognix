// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: connector_messages.proto

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

type ConnectorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int64             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Params map[string]string `protobuf:"bytes,2,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ConnectorRequest) Reset() {
	*x = ConnectorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectorRequest) ProtoMessage() {}

func (x *ConnectorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_connector_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectorRequest.ProtoReflect.Descriptor instead.
func (*ConnectorRequest) Descriptor() ([]byte, []int) {
	return file_connector_messages_proto_rawDescGZIP(), []int{0}
}

func (x *ConnectorRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ConnectorRequest) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

type EmbeddAsyncRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DocumentId int64  `protobuf:"varint,1,opt,name=document_id,json=documentId,proto3" json:"document_id,omitempty"`
	ChunkId    int64  `protobuf:"varint,2,opt,name=chunk_id,json=chunkId,proto3" json:"chunk_id,omitempty"`
	Content    string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	ModelId    string `protobuf:"bytes,4,opt,name=model_id,json=modelId,proto3" json:"model_id,omitempty"`
}

func (x *EmbeddAsyncRequest) Reset() {
	*x = EmbeddAsyncRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmbeddAsyncRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmbeddAsyncRequest) ProtoMessage() {}

func (x *EmbeddAsyncRequest) ProtoReflect() protoreflect.Message {
	mi := &file_connector_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmbeddAsyncRequest.ProtoReflect.Descriptor instead.
func (*EmbeddAsyncRequest) Descriptor() ([]byte, []int) {
	return file_connector_messages_proto_rawDescGZIP(), []int{1}
}

func (x *EmbeddAsyncRequest) GetDocumentId() int64 {
	if x != nil {
		return x.DocumentId
	}
	return 0
}

func (x *EmbeddAsyncRequest) GetChunkId() int64 {
	if x != nil {
		return x.ChunkId
	}
	return 0
}

func (x *EmbeddAsyncRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *EmbeddAsyncRequest) GetModelId() string {
	if x != nil {
		return x.ModelId
	}
	return ""
}

type EmbeddAsyncResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DocumentId int64     `protobuf:"varint,1,opt,name=document_id,json=documentId,proto3" json:"document_id,omitempty"`
	ChunkId    int64     `protobuf:"varint,2,opt,name=chunk_id,json=chunkId,proto3" json:"chunk_id,omitempty"`
	Content    string    `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	Vector     []float32 `protobuf:"fixed32,4,rep,packed,name=vector,proto3" json:"vector,omitempty"`
}

func (x *EmbeddAsyncResponse) Reset() {
	*x = EmbeddAsyncResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmbeddAsyncResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmbeddAsyncResponse) ProtoMessage() {}

func (x *EmbeddAsyncResponse) ProtoReflect() protoreflect.Message {
	mi := &file_connector_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmbeddAsyncResponse.ProtoReflect.Descriptor instead.
func (*EmbeddAsyncResponse) Descriptor() ([]byte, []int) {
	return file_connector_messages_proto_rawDescGZIP(), []int{2}
}

func (x *EmbeddAsyncResponse) GetDocumentId() int64 {
	if x != nil {
		return x.DocumentId
	}
	return 0
}

func (x *EmbeddAsyncResponse) GetChunkId() int64 {
	if x != nil {
		return x.ChunkId
	}
	return 0
}

func (x *EmbeddAsyncResponse) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *EmbeddAsyncResponse) GetVector() []float32 {
	if x != nil {
		return x.Vector
	}
	return nil
}

var File_connector_messages_proto protoreflect.FileDescriptor

var file_connector_messages_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x63, 0x6f, 0x6d, 0x2e,
	0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x22, 0x9f, 0x01, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x40, 0x0a, 0x06, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x39, 0x0a,
	0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x85, 0x01, 0x0a, 0x12, 0x45, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x41, 0x73, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x19, 0x0a, 0x08, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x49, 0x64,
	0x22, 0x83, 0x01, 0x0a, 0x13, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x41, 0x73, 0x79, 0x6e, 0x63,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x64,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x68, 0x75,
	0x6e, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x68, 0x75,
	0x6e, 0x6b, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x03, 0x28, 0x02, 0x52, 0x06,
	0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x42, 0x19, 0x5a, 0x17, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67,
	0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_connector_messages_proto_rawDescOnce sync.Once
	file_connector_messages_proto_rawDescData = file_connector_messages_proto_rawDesc
)

func file_connector_messages_proto_rawDescGZIP() []byte {
	file_connector_messages_proto_rawDescOnce.Do(func() {
		file_connector_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_connector_messages_proto_rawDescData)
	})
	return file_connector_messages_proto_rawDescData
}

var file_connector_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_connector_messages_proto_goTypes = []interface{}{
	(*ConnectorRequest)(nil),    // 0: com.embedd.ConnectorRequest
	(*EmbeddAsyncRequest)(nil),  // 1: com.embedd.EmbeddAsyncRequest
	(*EmbeddAsyncResponse)(nil), // 2: com.embedd.EmbeddAsyncResponse
	nil,                         // 3: com.embedd.ConnectorRequest.ParamsEntry
}
var file_connector_messages_proto_depIdxs = []int32{
	3, // 0: com.embedd.ConnectorRequest.params:type_name -> com.embedd.ConnectorRequest.ParamsEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_connector_messages_proto_init() }
func file_connector_messages_proto_init() {
	if File_connector_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_connector_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectorRequest); i {
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
		file_connector_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmbeddAsyncRequest); i {
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
		file_connector_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmbeddAsyncResponse); i {
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
			RawDescriptor: file_connector_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_connector_messages_proto_goTypes,
		DependencyIndexes: file_connector_messages_proto_depIdxs,
		MessageInfos:      file_connector_messages_proto_msgTypes,
	}.Build()
	File_connector_messages_proto = out.File
	file_connector_messages_proto_rawDesc = nil
	file_connector_messages_proto_goTypes = nil
	file_connector_messages_proto_depIdxs = nil
}
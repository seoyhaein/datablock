// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: fileblock.proto

package protos

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

// 단일 파일 블럭을 나타내는 메시지
type FileBlockData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockId       string   `protobuf:"bytes,1,opt,name=block_id,json=blockId,proto3" json:"block_id,omitempty"`                   // 블록을 구분하기 위한 고유 ID
	ColumnHeaders []string `protobuf:"bytes,2,rep,name=column_headers,json=columnHeaders,proto3" json:"column_headers,omitempty"` // 컬럼 이름들
	Rows          []*Row   `protobuf:"bytes,3,rep,name=rows,proto3" json:"rows,omitempty"`                                        // 행 데이터
}

func (x *FileBlockData) Reset() {
	*x = FileBlockData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileblock_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileBlockData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileBlockData) ProtoMessage() {}

func (x *FileBlockData) ProtoReflect() protoreflect.Message {
	mi := &file_fileblock_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileBlockData.ProtoReflect.Descriptor instead.
func (*FileBlockData) Descriptor() ([]byte, []int) {
	return file_fileblock_proto_rawDescGZIP(), []int{0}
}

func (x *FileBlockData) GetBlockId() string {
	if x != nil {
		return x.BlockId
	}
	return ""
}

func (x *FileBlockData) GetColumnHeaders() []string {
	if x != nil {
		return x.ColumnHeaders
	}
	return nil
}

func (x *FileBlockData) GetRows() []*Row {
	if x != nil {
		return x.Rows
	}
	return nil
}

// 하나의 행(row)을 나타내며, 행 번호와 헤더-값 매핑을 포함
type Row struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RowNumber   int32             `protobuf:"varint,1,opt,name=row_number,json=rowNumber,proto3" json:"row_number,omitempty"`                                                                                              // 행 번호
	CellColumns map[string]string `protobuf:"bytes,2,rep,name=cell_columns,json=cellColumns,proto3" json:"cell_columns,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // 헤더 이름과 셀 값 매핑
}

func (x *Row) Reset() {
	*x = Row{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileblock_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Row) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Row) ProtoMessage() {}

func (x *Row) ProtoReflect() protoreflect.Message {
	mi := &file_fileblock_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Row.ProtoReflect.Descriptor instead.
func (*Row) Descriptor() ([]byte, []int) {
	return file_fileblock_proto_rawDescGZIP(), []int{1}
}

func (x *Row) GetRowNumber() int32 {
	if x != nil {
		return x.RowNumber
	}
	return 0
}

func (x *Row) GetCellColumns() map[string]string {
	if x != nil {
		return x.CellColumns
	}
	return nil
}

var File_fileblock_proto protoreflect.FileDescriptor

var file_fileblock_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x66, 0x69, 0x6c, 0x65, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x22, 0x72, 0x0a, 0x0d, 0x46, 0x69, 0x6c,
	0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x5f,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x63,
	0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x1f, 0x0a, 0x04,
	0x72, 0x6f, 0x77, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2e, 0x52, 0x6f, 0x77, 0x52, 0x04, 0x72, 0x6f, 0x77, 0x73, 0x22, 0xa5, 0x01,
	0x0a, 0x03, 0x52, 0x6f, 0x77, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x6f, 0x77, 0x5f, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x72, 0x6f, 0x77, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x0c, 0x63, 0x65, 0x6c, 0x6c, 0x5f, 0x63, 0x6f, 0x6c,
	0x75, 0x6d, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2e, 0x52, 0x6f, 0x77, 0x2e, 0x43, 0x65, 0x6c, 0x6c, 0x43, 0x6f, 0x6c, 0x75,
	0x6d, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x63, 0x65, 0x6c, 0x6c, 0x43, 0x6f,
	0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x1a, 0x3e, 0x0a, 0x10, 0x43, 0x65, 0x6c, 0x6c, 0x43, 0x6f, 0x6c,
	0x75, 0x6d, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x6f, 0x79, 0x68, 0x61, 0x65, 0x69, 0x6e, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fileblock_proto_rawDescOnce sync.Once
	file_fileblock_proto_rawDescData = file_fileblock_proto_rawDesc
)

func file_fileblock_proto_rawDescGZIP() []byte {
	file_fileblock_proto_rawDescOnce.Do(func() {
		file_fileblock_proto_rawDescData = protoimpl.X.CompressGZIP(file_fileblock_proto_rawDescData)
	})
	return file_fileblock_proto_rawDescData
}

var file_fileblock_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_fileblock_proto_goTypes = []interface{}{
	(*FileBlockData)(nil), // 0: protos.FileBlockData
	(*Row)(nil),           // 1: protos.Row
	nil,                   // 2: protos.Row.CellColumnsEntry
}
var file_fileblock_proto_depIdxs = []int32{
	1, // 0: protos.FileBlockData.rows:type_name -> protos.Row
	2, // 1: protos.Row.cell_columns:type_name -> protos.Row.CellColumnsEntry
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_fileblock_proto_init() }
func file_fileblock_proto_init() {
	if File_fileblock_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fileblock_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileBlockData); i {
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
		file_fileblock_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Row); i {
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
			RawDescriptor: file_fileblock_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fileblock_proto_goTypes,
		DependencyIndexes: file_fileblock_proto_depIdxs,
		MessageInfos:      file_fileblock_proto_msgTypes,
	}.Build()
	File_fileblock_proto = out.File
	file_fileblock_proto_rawDesc = nil
	file_fileblock_proto_goTypes = nil
	file_fileblock_proto_depIdxs = nil
}
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.2
// source: genre_service.proto

package genres_pb

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

// Responses
type GetAllResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Categories []*CategoryModel `protobuf:"bytes,1,rep,name=categories,proto3" json:"categories,omitempty"`
}

func (x *GetAllResponse) Reset() {
	*x = GetAllResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genre_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllResponse) ProtoMessage() {}

func (x *GetAllResponse) ProtoReflect() protoreflect.Message {
	mi := &file_genre_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllResponse.ProtoReflect.Descriptor instead.
func (*GetAllResponse) Descriptor() ([]byte, []int) {
	return file_genre_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetAllResponse) GetCategories() []*CategoryModel {
	if x != nil {
		return x.Categories
	}
	return nil
}

// Models
type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genre_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_genre_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_genre_service_proto_rawDescGZIP(), []int{1}
}

type CategoryModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	TranslitName string        `protobuf:"bytes,2,opt,name=translitName,proto3" json:"translitName,omitempty"`
	Genres       []*GenreModel `protobuf:"bytes,3,rep,name=genres,proto3" json:"genres,omitempty"`
}

func (x *CategoryModel) Reset() {
	*x = CategoryModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genre_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CategoryModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CategoryModel) ProtoMessage() {}

func (x *CategoryModel) ProtoReflect() protoreflect.Message {
	mi := &file_genre_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CategoryModel.ProtoReflect.Descriptor instead.
func (*CategoryModel) Descriptor() ([]byte, []int) {
	return file_genre_service_proto_rawDescGZIP(), []int{2}
}

func (x *CategoryModel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CategoryModel) GetTranslitName() string {
	if x != nil {
		return x.TranslitName
	}
	return ""
}

func (x *CategoryModel) GetGenres() []*GenreModel {
	if x != nil {
		return x.Genres
	}
	return nil
}

type GenreModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	TranslitName string `protobuf:"bytes,2,opt,name=translitName,proto3" json:"translitName,omitempty"`
	BookCount    int64  `protobuf:"varint,3,opt,name=bookCount,proto3" json:"bookCount,omitempty"`
}

func (x *GenreModel) Reset() {
	*x = GenreModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genre_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenreModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenreModel) ProtoMessage() {}

func (x *GenreModel) ProtoReflect() protoreflect.Message {
	mi := &file_genre_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenreModel.ProtoReflect.Descriptor instead.
func (*GenreModel) Descriptor() ([]byte, []int) {
	return file_genre_service_proto_rawDescGZIP(), []int{3}
}

func (x *GenreModel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GenreModel) GetTranslitName() string {
	if x != nil {
		return x.TranslitName
	}
	return ""
}

func (x *GenreModel) GetBookCount() int64 {
	if x != nil {
		return x.BookCount
	}
	return 0
}

var File_genre_service_proto protoreflect.FileDescriptor

var file_genre_service_proto_rawDesc = []byte{
	0x0a, 0x13, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x22, 0x47, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x35, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x73, 0x0a, 0x0d, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x69, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x6c, 0x69, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x72,
	0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x65, 0x6e, 0x72, 0x65,
	0x73, 0x2e, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x06, 0x67, 0x65,
	0x6e, 0x72, 0x65, 0x73, 0x22, 0x62, 0x0a, 0x0a, 0x47, 0x65, 0x6e, 0x72, 0x65, 0x4d, 0x6f, 0x64,
	0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c,
	0x69, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x6c, 0x69, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6f,
	0x6f, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x62,
	0x6f, 0x6f, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0x38, 0x0a, 0x05, 0x47, 0x65, 0x6e, 0x72,
	0x65, 0x12, 0x2f, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x0d, 0x2e, 0x67, 0x65,
	0x6e, 0x72, 0x65, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x65, 0x6e,
	0x72, 0x65, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x1c, 0x5a, 0x1a, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x3b, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x73, 0x5f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_genre_service_proto_rawDescOnce sync.Once
	file_genre_service_proto_rawDescData = file_genre_service_proto_rawDesc
)

func file_genre_service_proto_rawDescGZIP() []byte {
	file_genre_service_proto_rawDescOnce.Do(func() {
		file_genre_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_genre_service_proto_rawDescData)
	})
	return file_genre_service_proto_rawDescData
}

var file_genre_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_genre_service_proto_goTypes = []any{
	(*GetAllResponse)(nil), // 0: genres.GetAllResponse
	(*Empty)(nil),          // 1: genres.Empty
	(*CategoryModel)(nil),  // 2: genres.CategoryModel
	(*GenreModel)(nil),     // 3: genres.GenreModel
}
var file_genre_service_proto_depIdxs = []int32{
	2, // 0: genres.GetAllResponse.categories:type_name -> genres.CategoryModel
	3, // 1: genres.CategoryModel.genres:type_name -> genres.GenreModel
	1, // 2: genres.Genre.GetAll:input_type -> genres.Empty
	0, // 3: genres.Genre.GetAll:output_type -> genres.GetAllResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_genre_service_proto_init() }
func file_genre_service_proto_init() {
	if File_genre_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_genre_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetAllResponse); i {
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
		file_genre_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Empty); i {
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
		file_genre_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CategoryModel); i {
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
		file_genre_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GenreModel); i {
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
			RawDescriptor: file_genre_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_genre_service_proto_goTypes,
		DependencyIndexes: file_genre_service_proto_depIdxs,
		MessageInfos:      file_genre_service_proto_msgTypes,
	}.Build()
	File_genre_service_proto = out.File
	file_genre_service_proto_rawDesc = nil
	file_genre_service_proto_goTypes = nil
	file_genre_service_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: author_service.proto

package authors_pb

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

// Requests
type GetAuthorsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// array of authors id to find
	Id []string `protobuf:"bytes,1,rep,name=id,proto3" json:"id,omitempty" form:"id" validate:"primitiveid,required_without_all=Translit"`  
	// array of translit names to find
	Translit []string `protobuf:"bytes,2,rep,name=translit,proto3" json:"translit,omitempty" form:"translit" validate:"required_without_all=Id"`  
}

func (x *GetAuthorsRequest) Reset() {
	*x = GetAuthorsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthorsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthorsRequest) ProtoMessage() {}

func (x *GetAuthorsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthorsRequest.ProtoReflect.Descriptor instead.
func (*GetAuthorsRequest) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetAuthorsRequest) GetId() []string {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *GetAuthorsRequest) GetTranslit() []string {
	if x != nil {
		return x.Translit
	}
	return nil
}

// Responses
type GetAuthorsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Authors []*AuthorModel `protobuf:"bytes,1,rep,name=authors,proto3" json:"authors,omitempty"`
}

func (x *GetAuthorsResponse) Reset() {
	*x = GetAuthorsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthorsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthorsResponse) ProtoMessage() {}

func (x *GetAuthorsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthorsResponse.ProtoReflect.Descriptor instead.
func (*GetAuthorsResponse) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetAuthorsResponse) GetAuthors() []*AuthorModel {
	if x != nil {
		return x.Authors
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
		mi := &file_author_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[2]
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
	return file_author_service_proto_rawDescGZIP(), []int{2}
}

type AuthorModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name           string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Translitname   string  `protobuf:"bytes,3,opt,name=translitname,proto3" json:"translitname,omitempty"`
	Profilepicture string  `protobuf:"bytes,4,opt,name=profilepicture,proto3" json:"profilepicture,omitempty"`
	About          string  `protobuf:"bytes,5,opt,name=about,proto3" json:"about,omitempty"`
	Rating         float32 `protobuf:"fixed32,6,opt,name=rating,proto3" json:"rating,omitempty"`
}

func (x *AuthorModel) Reset() {
	*x = AuthorModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorModel) ProtoMessage() {}

func (x *AuthorModel) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorModel.ProtoReflect.Descriptor instead.
func (*AuthorModel) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{3}
}

func (x *AuthorModel) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AuthorModel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AuthorModel) GetTranslitname() string {
	if x != nil {
		return x.Translitname
	}
	return ""
}

func (x *AuthorModel) GetProfilepicture() string {
	if x != nil {
		return x.Profilepicture
	}
	return ""
}

func (x *AuthorModel) GetAbout() string {
	if x != nil {
		return x.About
	}
	return ""
}

func (x *AuthorModel) GetRating() float32 {
	if x != nil {
		return x.Rating
	}
	return 0
}

var File_author_service_proto protoreflect.FileDescriptor

var file_author_service_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x22,
	0x3f, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x69, 0x74,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x69, 0x74,
	0x22, 0x44, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x73, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x07, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0xab, 0x01, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x69, 0x74, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x69, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x70, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x70, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x61, 0x62, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x61, 0x62, 0x6f, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x32, 0x4f, 0x0a,
	0x06, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x45, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x73, 0x12, 0x1a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1e,
	0x5a, 0x1c, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x73, 0x3b, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x5f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_author_service_proto_rawDescOnce sync.Once
	file_author_service_proto_rawDescData = file_author_service_proto_rawDesc
)

func file_author_service_proto_rawDescGZIP() []byte {
	file_author_service_proto_rawDescOnce.Do(func() {
		file_author_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_author_service_proto_rawDescData)
	})
	return file_author_service_proto_rawDescData
}

var file_author_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_author_service_proto_goTypes = []any{
	(*GetAuthorsRequest)(nil),  // 0: authors.GetAuthorsRequest
	(*GetAuthorsResponse)(nil), // 1: authors.GetAuthorsResponse
	(*Empty)(nil),              // 2: authors.Empty
	(*AuthorModel)(nil),        // 3: authors.AuthorModel
}
var file_author_service_proto_depIdxs = []int32{
	3, // 0: authors.GetAuthorsResponse.authors:type_name -> authors.AuthorModel
	0, // 1: authors.Author.GetAuthors:input_type -> authors.GetAuthorsRequest
	1, // 2: authors.Author.GetAuthors:output_type -> authors.GetAuthorsResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_author_service_proto_init() }
func file_author_service_proto_init() {
	if File_author_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_author_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetAuthorsRequest); i {
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
		file_author_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetAuthorsResponse); i {
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
		file_author_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
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
		file_author_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*AuthorModel); i {
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
			RawDescriptor: file_author_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_author_service_proto_goTypes,
		DependencyIndexes: file_author_service_proto_depIdxs,
		MessageInfos:      file_author_service_proto_msgTypes,
	}.Build()
	File_author_service_proto = out.File
	file_author_service_proto_rawDesc = nil
	file_author_service_proto_goTypes = nil
	file_author_service_proto_depIdxs = nil
}

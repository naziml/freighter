// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.14.0
// source: proto/freighter.proto

package freighter

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

type FileType int32

const (
	FileType_FILE    FileType = 0
	FileType_DIR     FileType = 1
	FileType_SYMLINK FileType = 2
	FileType_OTHER   FileType = 3
)

// Enum value maps for FileType.
var (
	FileType_name = map[int32]string{
		0: "FILE",
		1: "DIR",
		2: "SYMLINK",
		3: "OTHER",
	}
	FileType_value = map[string]int32{
		"FILE":    0,
		"DIR":     1,
		"SYMLINK": 2,
		"OTHER":   3,
	}
)

func (x FileType) Enum() *FileType {
	p := new(FileType)
	*p = x
	return p
}

func (x FileType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FileType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_freighter_proto_enumTypes[0].Descriptor()
}

func (FileType) Type() protoreflect.EnumType {
	return &file_proto_freighter_proto_enumTypes[0]
}

func (x FileType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FileType.Descriptor instead.
func (FileType) EnumDescriptor() ([]byte, []int) {
	return file_proto_freighter_proto_rawDescGZIP(), []int{0}
}

type TreeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repository string `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	Target     string `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
}

func (x *TreeRequest) Reset() {
	*x = TreeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_freighter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TreeRequest) ProtoMessage() {}

func (x *TreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_freighter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TreeRequest.ProtoReflect.Descriptor instead.
func (*TreeRequest) Descriptor() ([]byte, []int) {
	return file_proto_freighter_proto_rawDescGZIP(), []int{0}
}

func (x *TreeRequest) GetRepository() string {
	if x != nil {
		return x.Repository
	}
	return ""
}

func (x *TreeRequest) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

type TreeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files []*FileInfo `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *TreeReply) Reset() {
	*x = TreeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_freighter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TreeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TreeReply) ProtoMessage() {}

func (x *TreeReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_freighter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TreeReply.ProtoReflect.Descriptor instead.
func (*TreeReply) Descriptor() ([]byte, []int) {
	return file_proto_freighter_proto_rawDescGZIP(), []int{1}
}

func (x *TreeReply) GetFiles() []*FileInfo {
	if x != nil {
		return x.Files
	}
	return nil
}

type DirRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repository string `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	Target     string `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Path       string `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *DirRequest) Reset() {
	*x = DirRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_freighter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirRequest) ProtoMessage() {}

func (x *DirRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_freighter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirRequest.ProtoReflect.Descriptor instead.
func (*DirRequest) Descriptor() ([]byte, []int) {
	return file_proto_freighter_proto_rawDescGZIP(), []int{2}
}

func (x *DirRequest) GetRepository() string {
	if x != nil {
		return x.Repository
	}
	return ""
}

func (x *DirRequest) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *DirRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type DirReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files []*FileInfo `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *DirReply) Reset() {
	*x = DirReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_freighter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirReply) ProtoMessage() {}

func (x *DirReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_freighter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirReply.ProtoReflect.Descriptor instead.
func (*DirReply) Descriptor() ([]byte, []int) {
	return file_proto_freighter_proto_rawDescGZIP(), []int{3}
}

func (x *DirReply) GetFiles() []*FileInfo {
	if x != nil {
		return x.Files
	}
	return nil
}

type FileInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Path       string   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	IsDir      bool     `protobuf:"varint,3,opt,name=isDir,proto3" json:"isDir,omitempty"`
	Size       uint64   `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	ModTime    uint64   `protobuf:"varint,5,opt,name=modTime,proto3" json:"modTime,omitempty"`
	Mode       uint32   `protobuf:"varint,6,opt,name=mode,proto3" json:"mode,omitempty"`
	AccessTime uint64   `protobuf:"varint,7,opt,name=accessTime,proto3" json:"accessTime,omitempty"`
	ChangeTime uint64   `protobuf:"varint,8,opt,name=changeTime,proto3" json:"changeTime,omitempty"`
	Type       FileType `protobuf:"varint,9,opt,name=type,proto3,enum=freighter.FileType" json:"type,omitempty"`
	ExtraData  string   `protobuf:"bytes,10,opt,name=extraData,proto3" json:"extraData,omitempty"`
}

func (x *FileInfo) Reset() {
	*x = FileInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_freighter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileInfo) ProtoMessage() {}

func (x *FileInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_freighter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileInfo.ProtoReflect.Descriptor instead.
func (*FileInfo) Descriptor() ([]byte, []int) {
	return file_proto_freighter_proto_rawDescGZIP(), []int{4}
}

func (x *FileInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FileInfo) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *FileInfo) GetIsDir() bool {
	if x != nil {
		return x.IsDir
	}
	return false
}

func (x *FileInfo) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *FileInfo) GetModTime() uint64 {
	if x != nil {
		return x.ModTime
	}
	return 0
}

func (x *FileInfo) GetMode() uint32 {
	if x != nil {
		return x.Mode
	}
	return 0
}

func (x *FileInfo) GetAccessTime() uint64 {
	if x != nil {
		return x.AccessTime
	}
	return 0
}

func (x *FileInfo) GetChangeTime() uint64 {
	if x != nil {
		return x.ChangeTime
	}
	return 0
}

func (x *FileInfo) GetType() FileType {
	if x != nil {
		return x.Type
	}
	return FileType_FILE
}

func (x *FileInfo) GetExtraData() string {
	if x != nil {
		return x.ExtraData
	}
	return ""
}

type FileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repository string `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	Target     string `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Path       string `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *FileRequest) Reset() {
	*x = FileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_freighter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileRequest) ProtoMessage() {}

func (x *FileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_freighter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileRequest.ProtoReflect.Descriptor instead.
func (*FileRequest) Descriptor() ([]byte, []int) {
	return file_proto_freighter_proto_rawDescGZIP(), []int{5}
}

func (x *FileRequest) GetRepository() string {
	if x != nil {
		return x.Repository
	}
	return ""
}

func (x *FileRequest) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *FileRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type FileReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *FileReply) Reset() {
	*x = FileReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_freighter_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileReply) ProtoMessage() {}

func (x *FileReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_freighter_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileReply.ProtoReflect.Descriptor instead.
func (*FileReply) Descriptor() ([]byte, []int) {
	return file_proto_freighter_proto_rawDescGZIP(), []int{6}
}

func (x *FileReply) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_freighter_proto protoreflect.FileDescriptor

var file_proto_freighter_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x65, 0x72, 0x22, 0x45, 0x0a, 0x0b, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x22, 0x36, 0x0a, 0x09, 0x54, 0x72, 0x65,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65,
	0x72, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x22, 0x58, 0x0a, 0x0a, 0x44, 0x69, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x35, 0x0a, 0x08, 0x44,
	0x69, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x65, 0x72, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x22, 0x91, 0x02, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x73, 0x44, 0x69, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69, 0x73, 0x44, 0x69, 0x72, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x6f, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x6d, 0x6f, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6d,
	0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e,
	0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x72, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x72,
	0x61, 0x44, 0x61, 0x74, 0x61, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x78, 0x74,
	0x72, 0x61, 0x44, 0x61, 0x74, 0x61, 0x22, 0x59, 0x0a, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74,
	0x68, 0x22, 0x1f, 0x0a, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x2a, 0x35, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08,
	0x0a, 0x04, 0x46, 0x49, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x49, 0x52, 0x10,
	0x01, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x59, 0x4d, 0x4c, 0x49, 0x4e, 0x4b, 0x10, 0x02, 0x12, 0x09,
	0x0a, 0x05, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10, 0x03, 0x32, 0xb9, 0x01, 0x0a, 0x09, 0x46, 0x72,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x72, 0x12, 0x36, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x44, 0x69,
	0x72, 0x12, 0x15, 0x2e, 0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x69,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x66, 0x72, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12,
	0x39, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x66, 0x72, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x65, 0x72, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x14, 0x2e, 0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x72, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x07, 0x47, 0x65,
	0x74, 0x54, 0x72, 0x65, 0x65, 0x12, 0x16, 0x2e, 0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65,
	0x72, 0x2e, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x68, 0x6e, 0x65, 0x77, 0x61, 0x72, 0x74, 0x2f, 0x66, 0x72,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x72, 0x2f, 0x66, 0x72, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65,
	0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_freighter_proto_rawDescOnce sync.Once
	file_proto_freighter_proto_rawDescData = file_proto_freighter_proto_rawDesc
)

func file_proto_freighter_proto_rawDescGZIP() []byte {
	file_proto_freighter_proto_rawDescOnce.Do(func() {
		file_proto_freighter_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_freighter_proto_rawDescData)
	})
	return file_proto_freighter_proto_rawDescData
}

var file_proto_freighter_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_freighter_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_freighter_proto_goTypes = []interface{}{
	(FileType)(0),       // 0: freighter.FileType
	(*TreeRequest)(nil), // 1: freighter.TreeRequest
	(*TreeReply)(nil),   // 2: freighter.TreeReply
	(*DirRequest)(nil),  // 3: freighter.DirRequest
	(*DirReply)(nil),    // 4: freighter.DirReply
	(*FileInfo)(nil),    // 5: freighter.FileInfo
	(*FileRequest)(nil), // 6: freighter.FileRequest
	(*FileReply)(nil),   // 7: freighter.FileReply
}
var file_proto_freighter_proto_depIdxs = []int32{
	5, // 0: freighter.TreeReply.files:type_name -> freighter.FileInfo
	5, // 1: freighter.DirReply.files:type_name -> freighter.FileInfo
	0, // 2: freighter.FileInfo.type:type_name -> freighter.FileType
	3, // 3: freighter.Freighter.GetDir:input_type -> freighter.DirRequest
	6, // 4: freighter.Freighter.GetFile:input_type -> freighter.FileRequest
	1, // 5: freighter.Freighter.GetTree:input_type -> freighter.TreeRequest
	4, // 6: freighter.Freighter.GetDir:output_type -> freighter.DirReply
	7, // 7: freighter.Freighter.GetFile:output_type -> freighter.FileReply
	2, // 8: freighter.Freighter.GetTree:output_type -> freighter.TreeReply
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_freighter_proto_init() }
func file_proto_freighter_proto_init() {
	if File_proto_freighter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_freighter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TreeRequest); i {
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
		file_proto_freighter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TreeReply); i {
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
		file_proto_freighter_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirRequest); i {
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
		file_proto_freighter_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirReply); i {
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
		file_proto_freighter_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileInfo); i {
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
		file_proto_freighter_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileRequest); i {
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
		file_proto_freighter_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileReply); i {
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
			RawDescriptor: file_proto_freighter_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_freighter_proto_goTypes,
		DependencyIndexes: file_proto_freighter_proto_depIdxs,
		EnumInfos:         file_proto_freighter_proto_enumTypes,
		MessageInfos:      file_proto_freighter_proto_msgTypes,
	}.Build()
	File_proto_freighter_proto = out.File
	file_proto_freighter_proto_rawDesc = nil
	file_proto_freighter_proto_goTypes = nil
	file_proto_freighter_proto_depIdxs = nil
}

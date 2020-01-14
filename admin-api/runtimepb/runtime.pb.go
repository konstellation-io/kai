// Code generated by protoc-gen-go. DO NOT EDIT.
// source: runtimepb/runtime.proto

package runtimepb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Entrypoint struct {
	ProtoFile            string   `protobuf:"bytes,1,opt,name=proto_file,json=protoFile,proto3" json:"proto_file,omitempty"`
	Image                string   `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	Src                  string   `protobuf:"bytes,3,opt,name=src,proto3" json:"src,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Entrypoint) Reset()         { *m = Entrypoint{} }
func (m *Entrypoint) String() string { return proto.CompactTextString(m) }
func (*Entrypoint) ProtoMessage()    {}
func (*Entrypoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{0}
}

func (m *Entrypoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entrypoint.Unmarshal(m, b)
}
func (m *Entrypoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entrypoint.Marshal(b, m, deterministic)
}
func (m *Entrypoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entrypoint.Merge(m, src)
}
func (m *Entrypoint) XXX_Size() int {
	return xxx_messageInfo_Entrypoint.Size(m)
}
func (m *Entrypoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Entrypoint.DiscardUnknown(m)
}

var xxx_messageInfo_Entrypoint proto.InternalMessageInfo

func (m *Entrypoint) GetProtoFile() string {
	if m != nil {
		return m.ProtoFile
	}
	return ""
}

func (m *Entrypoint) GetImage() string {
	if m != nil {
		return m.Image
	}
	return ""
}

func (m *Entrypoint) GetSrc() string {
	if m != nil {
		return m.Src
	}
	return ""
}

type Workflow struct {
	Name                 string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Entrypoint           string           `protobuf:"bytes,2,opt,name=entrypoint,proto3" json:"entrypoint,omitempty"`
	Nodes                []*Workflow_Node `protobuf:"bytes,3,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Edges                []*Workflow_Edge `protobuf:"bytes,4,rep,name=edges,proto3" json:"edges,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Workflow) Reset()         { *m = Workflow{} }
func (m *Workflow) String() string { return proto.CompactTextString(m) }
func (*Workflow) ProtoMessage()    {}
func (*Workflow) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{1}
}

func (m *Workflow) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow.Unmarshal(m, b)
}
func (m *Workflow) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow.Marshal(b, m, deterministic)
}
func (m *Workflow) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow.Merge(m, src)
}
func (m *Workflow) XXX_Size() int {
	return xxx_messageInfo_Workflow.Size(m)
}
func (m *Workflow) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow proto.InternalMessageInfo

func (m *Workflow) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Workflow) GetEntrypoint() string {
	if m != nil {
		return m.Entrypoint
	}
	return ""
}

func (m *Workflow) GetNodes() []*Workflow_Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *Workflow) GetEdges() []*Workflow_Edge {
	if m != nil {
		return m.Edges
	}
	return nil
}

type Workflow_Node struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Image                string   `protobuf:"bytes,3,opt,name=Image,proto3" json:"Image,omitempty"`
	Src                  string   `protobuf:"bytes,4,opt,name=Src,proto3" json:"Src,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Node) Reset()         { *m = Workflow_Node{} }
func (m *Workflow_Node) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node) ProtoMessage()    {}
func (*Workflow_Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{1, 0}
}

func (m *Workflow_Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node.Unmarshal(m, b)
}
func (m *Workflow_Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node.Marshal(b, m, deterministic)
}
func (m *Workflow_Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node.Merge(m, src)
}
func (m *Workflow_Node) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node.Size(m)
}
func (m *Workflow_Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node proto.InternalMessageInfo

func (m *Workflow_Node) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Workflow_Node) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Workflow_Node) GetImage() string {
	if m != nil {
		return m.Image
	}
	return ""
}

func (m *Workflow_Node) GetSrc() string {
	if m != nil {
		return m.Src
	}
	return ""
}

type Workflow_Edge struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	FromNode             string   `protobuf:"bytes,2,opt,name=FromNode,proto3" json:"FromNode,omitempty"`
	ToNode               string   `protobuf:"bytes,3,opt,name=ToNode,proto3" json:"ToNode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Edge) Reset()         { *m = Workflow_Edge{} }
func (m *Workflow_Edge) String() string { return proto.CompactTextString(m) }
func (*Workflow_Edge) ProtoMessage()    {}
func (*Workflow_Edge) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{1, 1}
}

func (m *Workflow_Edge) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Edge.Unmarshal(m, b)
}
func (m *Workflow_Edge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Edge.Marshal(b, m, deterministic)
}
func (m *Workflow_Edge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Edge.Merge(m, src)
}
func (m *Workflow_Edge) XXX_Size() int {
	return xxx_messageInfo_Workflow_Edge.Size(m)
}
func (m *Workflow_Edge) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Edge.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Edge proto.InternalMessageInfo

func (m *Workflow_Edge) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Workflow_Edge) GetFromNode() string {
	if m != nil {
		return m.FromNode
	}
	return ""
}

func (m *Workflow_Edge) GetToNode() string {
	if m != nil {
		return m.ToNode
	}
	return ""
}

type Version struct {
	Name                 string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Entrypoint           *Entrypoint       `protobuf:"bytes,2,opt,name=entrypoint,proto3" json:"entrypoint,omitempty"`
	Config               []*Version_Config `protobuf:"bytes,3,rep,name=config,proto3" json:"config,omitempty"`
	Workflows            []*Workflow       `protobuf:"bytes,4,rep,name=workflows,proto3" json:"workflows,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{2}
}

func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (m *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(m, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Version) GetEntrypoint() *Entrypoint {
	if m != nil {
		return m.Entrypoint
	}
	return nil
}

func (m *Version) GetConfig() []*Version_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *Version) GetWorkflows() []*Workflow {
	if m != nil {
		return m.Workflows
	}
	return nil
}

type Version_Config struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Version_Config) Reset()         { *m = Version_Config{} }
func (m *Version_Config) String() string { return proto.CompactTextString(m) }
func (*Version_Config) ProtoMessage()    {}
func (*Version_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{2, 0}
}

func (m *Version_Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version_Config.Unmarshal(m, b)
}
func (m *Version_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version_Config.Marshal(b, m, deterministic)
}
func (m *Version_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version_Config.Merge(m, src)
}
func (m *Version_Config) XXX_Size() int {
	return xxx_messageInfo_Version_Config.Size(m)
}
func (m *Version_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Version_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Version_Config proto.InternalMessageInfo

func (m *Version_Config) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Version_Config) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type DeployVersionRequest struct {
	Version              *Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeployVersionRequest) Reset()         { *m = DeployVersionRequest{} }
func (m *DeployVersionRequest) String() string { return proto.CompactTextString(m) }
func (*DeployVersionRequest) ProtoMessage()    {}
func (*DeployVersionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{3}
}

func (m *DeployVersionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeployVersionRequest.Unmarshal(m, b)
}
func (m *DeployVersionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeployVersionRequest.Marshal(b, m, deterministic)
}
func (m *DeployVersionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployVersionRequest.Merge(m, src)
}
func (m *DeployVersionRequest) XXX_Size() int {
	return xxx_messageInfo_DeployVersionRequest.Size(m)
}
func (m *DeployVersionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployVersionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeployVersionRequest proto.InternalMessageInfo

func (m *DeployVersionRequest) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

type DeployVersionResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeployVersionResponse) Reset()         { *m = DeployVersionResponse{} }
func (m *DeployVersionResponse) String() string { return proto.CompactTextString(m) }
func (*DeployVersionResponse) ProtoMessage()    {}
func (*DeployVersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{4}
}

func (m *DeployVersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeployVersionResponse.Unmarshal(m, b)
}
func (m *DeployVersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeployVersionResponse.Marshal(b, m, deterministic)
}
func (m *DeployVersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployVersionResponse.Merge(m, src)
}
func (m *DeployVersionResponse) XXX_Size() int {
	return xxx_messageInfo_DeployVersionResponse.Size(m)
}
func (m *DeployVersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployVersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeployVersionResponse proto.InternalMessageInfo

func (m *DeployVersionResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *DeployVersionResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type StopVersionRequest struct {
	Version              *Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopVersionRequest) Reset()         { *m = StopVersionRequest{} }
func (m *StopVersionRequest) String() string { return proto.CompactTextString(m) }
func (*StopVersionRequest) ProtoMessage()    {}
func (*StopVersionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{5}
}

func (m *StopVersionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopVersionRequest.Unmarshal(m, b)
}
func (m *StopVersionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopVersionRequest.Marshal(b, m, deterministic)
}
func (m *StopVersionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopVersionRequest.Merge(m, src)
}
func (m *StopVersionRequest) XXX_Size() int {
	return xxx_messageInfo_StopVersionRequest.Size(m)
}
func (m *StopVersionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StopVersionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StopVersionRequest proto.InternalMessageInfo

func (m *StopVersionRequest) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

type StopVersionResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopVersionResponse) Reset()         { *m = StopVersionResponse{} }
func (m *StopVersionResponse) String() string { return proto.CompactTextString(m) }
func (*StopVersionResponse) ProtoMessage()    {}
func (*StopVersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{6}
}

func (m *StopVersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopVersionResponse.Unmarshal(m, b)
}
func (m *StopVersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopVersionResponse.Marshal(b, m, deterministic)
}
func (m *StopVersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopVersionResponse.Merge(m, src)
}
func (m *StopVersionResponse) XXX_Size() int {
	return xxx_messageInfo_StopVersionResponse.Size(m)
}
func (m *StopVersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StopVersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StopVersionResponse proto.InternalMessageInfo

func (m *StopVersionResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *StopVersionResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type DeactivateVersionRequest struct {
	Version              *Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeactivateVersionRequest) Reset()         { *m = DeactivateVersionRequest{} }
func (m *DeactivateVersionRequest) String() string { return proto.CompactTextString(m) }
func (*DeactivateVersionRequest) ProtoMessage()    {}
func (*DeactivateVersionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{7}
}

func (m *DeactivateVersionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeactivateVersionRequest.Unmarshal(m, b)
}
func (m *DeactivateVersionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeactivateVersionRequest.Marshal(b, m, deterministic)
}
func (m *DeactivateVersionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeactivateVersionRequest.Merge(m, src)
}
func (m *DeactivateVersionRequest) XXX_Size() int {
	return xxx_messageInfo_DeactivateVersionRequest.Size(m)
}
func (m *DeactivateVersionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeactivateVersionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeactivateVersionRequest proto.InternalMessageInfo

func (m *DeactivateVersionRequest) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

type DeactivateVersionResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeactivateVersionResponse) Reset()         { *m = DeactivateVersionResponse{} }
func (m *DeactivateVersionResponse) String() string { return proto.CompactTextString(m) }
func (*DeactivateVersionResponse) ProtoMessage()    {}
func (*DeactivateVersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{8}
}

func (m *DeactivateVersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeactivateVersionResponse.Unmarshal(m, b)
}
func (m *DeactivateVersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeactivateVersionResponse.Marshal(b, m, deterministic)
}
func (m *DeactivateVersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeactivateVersionResponse.Merge(m, src)
}
func (m *DeactivateVersionResponse) XXX_Size() int {
	return xxx_messageInfo_DeactivateVersionResponse.Size(m)
}
func (m *DeactivateVersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeactivateVersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeactivateVersionResponse proto.InternalMessageInfo

func (m *DeactivateVersionResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *DeactivateVersionResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type ActivateVersionRequest struct {
	Version              *Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ActivateVersionRequest) Reset()         { *m = ActivateVersionRequest{} }
func (m *ActivateVersionRequest) String() string { return proto.CompactTextString(m) }
func (*ActivateVersionRequest) ProtoMessage()    {}
func (*ActivateVersionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{9}
}

func (m *ActivateVersionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ActivateVersionRequest.Unmarshal(m, b)
}
func (m *ActivateVersionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ActivateVersionRequest.Marshal(b, m, deterministic)
}
func (m *ActivateVersionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ActivateVersionRequest.Merge(m, src)
}
func (m *ActivateVersionRequest) XXX_Size() int {
	return xxx_messageInfo_ActivateVersionRequest.Size(m)
}
func (m *ActivateVersionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ActivateVersionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ActivateVersionRequest proto.InternalMessageInfo

func (m *ActivateVersionRequest) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

type ActivateVersionResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ActivateVersionResponse) Reset()         { *m = ActivateVersionResponse{} }
func (m *ActivateVersionResponse) String() string { return proto.CompactTextString(m) }
func (*ActivateVersionResponse) ProtoMessage()    {}
func (*ActivateVersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{10}
}

func (m *ActivateVersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ActivateVersionResponse.Unmarshal(m, b)
}
func (m *ActivateVersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ActivateVersionResponse.Marshal(b, m, deterministic)
}
func (m *ActivateVersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ActivateVersionResponse.Merge(m, src)
}
func (m *ActivateVersionResponse) XXX_Size() int {
	return xxx_messageInfo_ActivateVersionResponse.Size(m)
}
func (m *ActivateVersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ActivateVersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ActivateVersionResponse proto.InternalMessageInfo

func (m *ActivateVersionResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *ActivateVersionResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*Entrypoint)(nil), "runtime.Entrypoint")
	proto.RegisterType((*Workflow)(nil), "runtime.Workflow")
	proto.RegisterType((*Workflow_Node)(nil), "runtime.Workflow.Node")
	proto.RegisterType((*Workflow_Edge)(nil), "runtime.Workflow.Edge")
	proto.RegisterType((*Version)(nil), "runtime.Version")
	proto.RegisterType((*Version_Config)(nil), "runtime.Version.Config")
	proto.RegisterType((*DeployVersionRequest)(nil), "runtime.DeployVersionRequest")
	proto.RegisterType((*DeployVersionResponse)(nil), "runtime.DeployVersionResponse")
	proto.RegisterType((*StopVersionRequest)(nil), "runtime.StopVersionRequest")
	proto.RegisterType((*StopVersionResponse)(nil), "runtime.StopVersionResponse")
	proto.RegisterType((*DeactivateVersionRequest)(nil), "runtime.DeactivateVersionRequest")
	proto.RegisterType((*DeactivateVersionResponse)(nil), "runtime.DeactivateVersionResponse")
	proto.RegisterType((*ActivateVersionRequest)(nil), "runtime.ActivateVersionRequest")
	proto.RegisterType((*ActivateVersionResponse)(nil), "runtime.ActivateVersionResponse")
}

func init() { proto.RegisterFile("runtimepb/runtime.proto", fileDescriptor_d0e5095094a8d27f) }

var fileDescriptor_d0e5095094a8d27f = []byte{
	// 565 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xdb, 0x6e, 0xd3, 0x40,
	0x10, 0x25, 0x97, 0xe6, 0x32, 0x11, 0xa5, 0xdd, 0x96, 0xc4, 0x18, 0x5a, 0x82, 0x9f, 0x2a, 0x84,
	0x12, 0x94, 0xfe, 0x00, 0x2d, 0x69, 0xa4, 0x14, 0x51, 0x90, 0x83, 0x8a, 0x84, 0x90, 0x50, 0x6a,
	0x4f, 0x22, 0xab, 0x89, 0xd7, 0xec, 0x3a, 0xa9, 0xf2, 0x6b, 0x7c, 0x0f, 0xe2, 0x3b, 0xd0, 0x5e,
	0xbc, 0xb9, 0x39, 0xbc, 0x98, 0xb7, 0x9d, 0x9d, 0xb3, 0x67, 0xce, 0x9c, 0x19, 0x69, 0xa1, 0xc1,
	0x66, 0x61, 0x1c, 0x4c, 0x31, 0xba, 0x6b, 0xeb, 0x53, 0x2b, 0x62, 0x34, 0xa6, 0xa4, 0xac, 0x43,
	0x67, 0x00, 0x70, 0x15, 0xc6, 0x6c, 0x11, 0xd1, 0x20, 0x8c, 0xc9, 0x09, 0x80, 0xcc, 0xff, 0x18,
	0x05, 0x13, 0xb4, 0x72, 0xcd, 0xdc, 0x59, 0xd5, 0xad, 0xca, 0x9b, 0x5e, 0x30, 0x41, 0x72, 0x0c,
	0x7b, 0xc1, 0x74, 0x38, 0x46, 0x2b, 0x2f, 0x33, 0x2a, 0x20, 0x07, 0x50, 0xe0, 0xcc, 0xb3, 0x0a,
	0xf2, 0x4e, 0x1c, 0x9d, 0x5f, 0x79, 0xa8, 0x7c, 0xa5, 0xec, 0x7e, 0x34, 0xa1, 0x0f, 0x84, 0x40,
	0x31, 0x1c, 0x4e, 0x13, 0x36, 0x79, 0x26, 0xa7, 0x00, 0x68, 0xaa, 0x6a, 0xb6, 0x95, 0x1b, 0xf2,
	0x06, 0xf6, 0x42, 0xea, 0x23, 0xb7, 0x0a, 0xcd, 0xc2, 0x59, 0xad, 0x53, 0x6f, 0x25, 0xea, 0x13,
	0xd6, 0xd6, 0x0d, 0xf5, 0xd1, 0x55, 0x20, 0x81, 0x46, 0x7f, 0x8c, 0xdc, 0x2a, 0xee, 0x42, 0x5f,
	0xf9, 0x63, 0x74, 0x15, 0xc8, 0x76, 0xa1, 0x28, 0x1e, 0x93, 0x7d, 0xc8, 0xf7, 0x7d, 0xad, 0x2a,
	0xdf, 0xf7, 0x85, 0xce, 0x1b, 0xa1, 0x53, 0xa9, 0x91, 0x67, 0xd1, 0x70, 0x5f, 0x36, 0xac, 0x9a,
	0x53, 0x81, 0x68, 0x78, 0xc0, 0x3c, 0xab, 0xa8, 0x1a, 0x1e, 0x30, 0xcf, 0xbe, 0x86, 0xa2, 0x28,
	0xb1, 0xc5, 0x69, 0x43, 0xa5, 0xc7, 0xe8, 0x54, 0xd4, 0xd3, 0xbc, 0x26, 0x26, 0x75, 0x28, 0x7d,
	0xa1, 0x32, 0xa3, 0xc8, 0x75, 0xe4, 0xfc, 0xc9, 0x41, 0xf9, 0x16, 0x19, 0x0f, 0x68, 0x98, 0xea,
	0xdd, 0xf9, 0x96, 0x77, 0xb5, 0xce, 0x91, 0x69, 0x79, 0x39, 0xcc, 0x35, 0x43, 0xdb, 0x50, 0xf2,
	0x68, 0x38, 0x0a, 0xc6, 0xda, 0xd1, 0x86, 0x79, 0xa0, 0x4b, 0xb5, 0xde, 0xcb, 0xb4, 0xab, 0x61,
	0xa4, 0x0d, 0xd5, 0x07, 0xed, 0x5e, 0xe2, 0xeb, 0xe1, 0x96, 0xaf, 0xee, 0x12, 0x63, 0xbf, 0x85,
	0x92, 0xa2, 0x10, 0xf6, 0xdc, 0xe3, 0x42, 0x6b, 0x16, 0x47, 0x61, 0xe3, 0x7c, 0x38, 0x99, 0x99,
	0xbd, 0x91, 0x81, 0x73, 0x09, 0xc7, 0x5d, 0x8c, 0x26, 0x74, 0xa1, 0x25, 0xb8, 0xf8, 0x73, 0x86,
	0x3c, 0x26, 0xaf, 0xa1, 0x3c, 0x57, 0x37, 0x92, 0xa3, 0xd6, 0x39, 0xd8, 0x14, 0xeb, 0x26, 0x00,
	0xe7, 0x03, 0x3c, 0xdd, 0xe0, 0xe0, 0x11, 0x0d, 0x39, 0x12, 0x0b, 0xca, 0x7c, 0xe6, 0x79, 0xc8,
	0xb9, 0x24, 0xa9, 0xb8, 0x49, 0x28, 0x32, 0x53, 0xe4, 0x7c, 0xb9, 0xc6, 0x49, 0xe8, 0xbc, 0x03,
	0x32, 0x88, 0x69, 0x94, 0x41, 0x4e, 0x1f, 0x8e, 0xd6, 0x18, 0x32, 0x88, 0xe9, 0x81, 0xd5, 0xc5,
	0xa1, 0x17, 0x07, 0xf3, 0x61, 0x8c, 0x19, 0x24, 0x7d, 0x82, 0x67, 0x29, 0x3c, 0x19, 0x84, 0x75,
	0xa1, 0x7e, 0x91, 0x5d, 0xd6, 0x47, 0x68, 0x5c, 0xfc, 0x3f, 0x51, 0x9d, 0xdf, 0x79, 0xd8, 0x77,
	0x55, 0xad, 0x01, 0xb2, 0x79, 0xe0, 0x21, 0xf9, 0x0c, 0x8f, 0xd7, 0x56, 0x83, 0x9c, 0x18, 0x35,
	0x69, 0x6b, 0x67, 0x9f, 0xee, 0x4a, 0x2b, 0x59, 0xce, 0x23, 0x72, 0x0d, 0xb5, 0x95, 0xe9, 0x92,
	0xe7, 0xe6, 0xc1, 0xf6, 0xd6, 0xd8, 0x2f, 0xd2, 0x93, 0x86, 0xeb, 0x3b, 0x1c, 0x6e, 0x8d, 0x85,
	0xbc, 0x5a, 0x91, 0x90, 0x3e, 0x7a, 0xdb, 0xf9, 0x17, 0xc4, 0xb0, 0xdf, 0xc2, 0x93, 0x0d, 0x77,
	0xc9, 0x4b, 0xf3, 0x30, 0x7d, 0x7a, 0x76, 0x73, 0x37, 0x20, 0xe1, 0xbd, 0xac, 0x7d, 0xab, 0x9a,
	0x1f, 0xe5, 0xae, 0x24, 0x3f, 0x86, 0xf3, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1d, 0x07, 0x97,
	0x12, 0x65, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RuntimeServiceClient is the client API for RuntimeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RuntimeServiceClient interface {
	DeployVersion(ctx context.Context, in *DeployVersionRequest, opts ...grpc.CallOption) (*DeployVersionResponse, error)
	StopVersion(ctx context.Context, in *StopVersionRequest, opts ...grpc.CallOption) (*StopVersionResponse, error)
	DeactivateVersion(ctx context.Context, in *DeactivateVersionRequest, opts ...grpc.CallOption) (*DeactivateVersionResponse, error)
	ActivateVersion(ctx context.Context, in *ActivateVersionRequest, opts ...grpc.CallOption) (*ActivateVersionResponse, error)
}

type runtimeServiceClient struct {
	cc *grpc.ClientConn
}

func NewRuntimeServiceClient(cc *grpc.ClientConn) RuntimeServiceClient {
	return &runtimeServiceClient{cc}
}

func (c *runtimeServiceClient) DeployVersion(ctx context.Context, in *DeployVersionRequest, opts ...grpc.CallOption) (*DeployVersionResponse, error) {
	out := new(DeployVersionResponse)
	err := c.cc.Invoke(ctx, "/runtime.RuntimeService/DeployVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) StopVersion(ctx context.Context, in *StopVersionRequest, opts ...grpc.CallOption) (*StopVersionResponse, error) {
	out := new(StopVersionResponse)
	err := c.cc.Invoke(ctx, "/runtime.RuntimeService/StopVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) DeactivateVersion(ctx context.Context, in *DeactivateVersionRequest, opts ...grpc.CallOption) (*DeactivateVersionResponse, error) {
	out := new(DeactivateVersionResponse)
	err := c.cc.Invoke(ctx, "/runtime.RuntimeService/DeactivateVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) ActivateVersion(ctx context.Context, in *ActivateVersionRequest, opts ...grpc.CallOption) (*ActivateVersionResponse, error) {
	out := new(ActivateVersionResponse)
	err := c.cc.Invoke(ctx, "/runtime.RuntimeService/ActivateVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RuntimeServiceServer is the server API for RuntimeService service.
type RuntimeServiceServer interface {
	DeployVersion(context.Context, *DeployVersionRequest) (*DeployVersionResponse, error)
	StopVersion(context.Context, *StopVersionRequest) (*StopVersionResponse, error)
	DeactivateVersion(context.Context, *DeactivateVersionRequest) (*DeactivateVersionResponse, error)
	ActivateVersion(context.Context, *ActivateVersionRequest) (*ActivateVersionResponse, error)
}

// UnimplementedRuntimeServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRuntimeServiceServer struct {
}

func (*UnimplementedRuntimeServiceServer) DeployVersion(ctx context.Context, req *DeployVersionRequest) (*DeployVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployVersion not implemented")
}
func (*UnimplementedRuntimeServiceServer) StopVersion(ctx context.Context, req *StopVersionRequest) (*StopVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopVersion not implemented")
}
func (*UnimplementedRuntimeServiceServer) DeactivateVersion(ctx context.Context, req *DeactivateVersionRequest) (*DeactivateVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateVersion not implemented")
}
func (*UnimplementedRuntimeServiceServer) ActivateVersion(ctx context.Context, req *ActivateVersionRequest) (*ActivateVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateVersion not implemented")
}

func RegisterRuntimeServiceServer(s *grpc.Server, srv RuntimeServiceServer) {
	s.RegisterService(&_RuntimeService_serviceDesc, srv)
}

func _RuntimeService_DeployVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).DeployVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runtime.RuntimeService/DeployVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).DeployVersion(ctx, req.(*DeployVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_StopVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).StopVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runtime.RuntimeService/StopVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).StopVersion(ctx, req.(*StopVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_DeactivateVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeactivateVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).DeactivateVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runtime.RuntimeService/DeactivateVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).DeactivateVersion(ctx, req.(*DeactivateVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_ActivateVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).ActivateVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runtime.RuntimeService/ActivateVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).ActivateVersion(ctx, req.(*ActivateVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RuntimeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "runtime.RuntimeService",
	HandlerType: (*RuntimeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeployVersion",
			Handler:    _RuntimeService_DeployVersion_Handler,
		},
		{
			MethodName: "StopVersion",
			Handler:    _RuntimeService_StopVersion_Handler,
		},
		{
			MethodName: "DeactivateVersion",
			Handler:    _RuntimeService_DeactivateVersion_Handler,
		},
		{
			MethodName: "ActivateVersion",
			Handler:    _RuntimeService_ActivateVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "runtimepb/runtime.proto",
}

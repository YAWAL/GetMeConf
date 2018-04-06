// Code generated by protoc-gen-go. DO NOT EDIT.
// source: getmeconfig.proto

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	getmeconfig.proto

It has these top-level messages:
	GetConfigByNameRequest
	DeleteConfigRequest
	GetConfigsByTypeRequest
	GetConfigResponce
	Config
	Responce
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetConfigByNameRequest struct {
	ConfigName string `protobuf:"bytes,1,opt,name=ConfigName" json:"ConfigName,omitempty"`
	ConfigType string `protobuf:"bytes,2,opt,name=ConfigType" json:"ConfigType,omitempty"`
}

func (m *GetConfigByNameRequest) Reset()                    { *m = GetConfigByNameRequest{} }
func (m *GetConfigByNameRequest) String() string            { return proto.CompactTextString(m) }
func (*GetConfigByNameRequest) ProtoMessage()               {}
func (*GetConfigByNameRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetConfigByNameRequest) GetConfigName() string {
	if m != nil {
		return m.ConfigName
	}
	return ""
}

func (m *GetConfigByNameRequest) GetConfigType() string {
	if m != nil {
		return m.ConfigType
	}
	return ""
}

type DeleteConfigRequest struct {
	ConfigName string `protobuf:"bytes,1,opt,name=ConfigName" json:"ConfigName,omitempty"`
	ConfigType string `protobuf:"bytes,2,opt,name=ConfigType" json:"ConfigType,omitempty"`
}

func (m *DeleteConfigRequest) Reset()                    { *m = DeleteConfigRequest{} }
func (m *DeleteConfigRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteConfigRequest) ProtoMessage()               {}
func (*DeleteConfigRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DeleteConfigRequest) GetConfigName() string {
	if m != nil {
		return m.ConfigName
	}
	return ""
}

func (m *DeleteConfigRequest) GetConfigType() string {
	if m != nil {
		return m.ConfigType
	}
	return ""
}

type GetConfigsByTypeRequest struct {
	ConfigType string `protobuf:"bytes,1,opt,name=ConfigType" json:"ConfigType,omitempty"`
}

func (m *GetConfigsByTypeRequest) Reset()                    { *m = GetConfigsByTypeRequest{} }
func (m *GetConfigsByTypeRequest) String() string            { return proto.CompactTextString(m) }
func (*GetConfigsByTypeRequest) ProtoMessage()               {}
func (*GetConfigsByTypeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetConfigsByTypeRequest) GetConfigType() string {
	if m != nil {
		return m.ConfigType
	}
	return ""
}

type GetConfigResponce struct {
	Config []byte `protobuf:"bytes,1,opt,name=Config,proto3" json:"Config,omitempty"`
}

func (m *GetConfigResponce) Reset()                    { *m = GetConfigResponce{} }
func (m *GetConfigResponce) String() string            { return proto.CompactTextString(m) }
func (*GetConfigResponce) ProtoMessage()               {}
func (*GetConfigResponce) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetConfigResponce) GetConfig() []byte {
	if m != nil {
		return m.Config
	}
	return nil
}

type Config struct {
	Config     []byte `protobuf:"bytes,1,opt,name=Config,proto3" json:"Config,omitempty"`
	ConfigType string `protobuf:"bytes,2,opt,name=ConfigType" json:"ConfigType,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Config) GetConfig() []byte {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *Config) GetConfigType() string {
	if m != nil {
		return m.ConfigType
	}
	return ""
}

type Responce struct {
	Status string `protobuf:"bytes,1,opt,name=Status" json:"Status,omitempty"`
}

func (m *Responce) Reset()                    { *m = Responce{} }
func (m *Responce) String() string            { return proto.CompactTextString(m) }
func (*Responce) ProtoMessage()               {}
func (*Responce) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Responce) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*GetConfigByNameRequest)(nil), "api.GetConfigByNameRequest")
	proto.RegisterType((*DeleteConfigRequest)(nil), "api.DeleteConfigRequest")
	proto.RegisterType((*GetConfigsByTypeRequest)(nil), "api.GetConfigsByTypeRequest")
	proto.RegisterType((*GetConfigResponce)(nil), "api.GetConfigResponce")
	proto.RegisterType((*Config)(nil), "api.Config")
	proto.RegisterType((*Responce)(nil), "api.Responce")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ConfigService service

type ConfigServiceClient interface {
	GetConfigByName(ctx context.Context, in *GetConfigByNameRequest, opts ...client.CallOption) (*GetConfigResponce, error)
	GetConfigsByType(ctx context.Context, in *GetConfigsByTypeRequest, opts ...client.CallOption) (ConfigService_GetConfigsByTypeClient, error)
	CreateConfig(ctx context.Context, in *Config, opts ...client.CallOption) (*Responce, error)
	DeleteConfig(ctx context.Context, in *DeleteConfigRequest, opts ...client.CallOption) (*Responce, error)
	UpdateConfig(ctx context.Context, in *Config, opts ...client.CallOption) (*Responce, error)
}

type configServiceClient struct {
	c           client.Client
	serviceName string
}

func NewConfigServiceClient(serviceName string, c client.Client) ConfigServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "api"
	}
	return &configServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *configServiceClient) GetConfigByName(ctx context.Context, in *GetConfigByNameRequest, opts ...client.CallOption) (*GetConfigResponce, error) {
	req := c.c.NewRequest(c.serviceName, "ConfigService.GetConfigByName", in)
	out := new(GetConfigResponce)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceClient) GetConfigsByType(ctx context.Context, in *GetConfigsByTypeRequest, opts ...client.CallOption) (ConfigService_GetConfigsByTypeClient, error) {
	req := c.c.NewRequest(c.serviceName, "ConfigService.GetConfigsByType", &GetConfigsByTypeRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &configServiceGetConfigsByTypeClient{stream}, nil
}

type ConfigService_GetConfigsByTypeClient interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*GetConfigResponce, error)
}

type configServiceGetConfigsByTypeClient struct {
	stream client.Streamer
}

func (x *configServiceGetConfigsByTypeClient) Close() error {
	return x.stream.Close()
}

func (x *configServiceGetConfigsByTypeClient) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *configServiceGetConfigsByTypeClient) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *configServiceGetConfigsByTypeClient) Recv() (*GetConfigResponce, error) {
	m := new(GetConfigResponce)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *configServiceClient) CreateConfig(ctx context.Context, in *Config, opts ...client.CallOption) (*Responce, error) {
	req := c.c.NewRequest(c.serviceName, "ConfigService.CreateConfig", in)
	out := new(Responce)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceClient) DeleteConfig(ctx context.Context, in *DeleteConfigRequest, opts ...client.CallOption) (*Responce, error) {
	req := c.c.NewRequest(c.serviceName, "ConfigService.DeleteConfig", in)
	out := new(Responce)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceClient) UpdateConfig(ctx context.Context, in *Config, opts ...client.CallOption) (*Responce, error) {
	req := c.c.NewRequest(c.serviceName, "ConfigService.UpdateConfig", in)
	out := new(Responce)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ConfigService service

type ConfigServiceHandler interface {
	GetConfigByName(context.Context, *GetConfigByNameRequest, *GetConfigResponce) error
	GetConfigsByType(context.Context, *GetConfigsByTypeRequest, ConfigService_GetConfigsByTypeStream) error
	CreateConfig(context.Context, *Config, *Responce) error
	DeleteConfig(context.Context, *DeleteConfigRequest, *Responce) error
	UpdateConfig(context.Context, *Config, *Responce) error
}

func RegisterConfigServiceHandler(s server.Server, hdlr ConfigServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&ConfigService{hdlr}, opts...))
}

type ConfigService struct {
	ConfigServiceHandler
}

func (h *ConfigService) GetConfigByName(ctx context.Context, in *GetConfigByNameRequest, out *GetConfigResponce) error {
	return h.ConfigServiceHandler.GetConfigByName(ctx, in, out)
}

func (h *ConfigService) GetConfigsByType(ctx context.Context, stream server.Streamer) error {
	m := new(GetConfigsByTypeRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.ConfigServiceHandler.GetConfigsByType(ctx, m, &configServiceGetConfigsByTypeStream{stream})
}

type ConfigService_GetConfigsByTypeStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*GetConfigResponce) error
}

type configServiceGetConfigsByTypeStream struct {
	stream server.Streamer
}

func (x *configServiceGetConfigsByTypeStream) Close() error {
	return x.stream.Close()
}

func (x *configServiceGetConfigsByTypeStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *configServiceGetConfigsByTypeStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *configServiceGetConfigsByTypeStream) Send(m *GetConfigResponce) error {
	return x.stream.Send(m)
}

func (h *ConfigService) CreateConfig(ctx context.Context, in *Config, out *Responce) error {
	return h.ConfigServiceHandler.CreateConfig(ctx, in, out)
}

func (h *ConfigService) DeleteConfig(ctx context.Context, in *DeleteConfigRequest, out *Responce) error {
	return h.ConfigServiceHandler.DeleteConfig(ctx, in, out)
}

func (h *ConfigService) UpdateConfig(ctx context.Context, in *Config, out *Responce) error {
	return h.ConfigServiceHandler.UpdateConfig(ctx, in, out)
}

func init() { proto.RegisterFile("getmeconfig.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0x4f, 0x2d, 0xc9,
	0x4d, 0x4d, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e,
	0x2c, 0xc8, 0x54, 0x8a, 0xe0, 0x12, 0x73, 0x4f, 0x2d, 0x71, 0x06, 0x8b, 0x3b, 0x55, 0xfa, 0x25,
	0xe6, 0xa6, 0x06, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0xc9, 0x71, 0x71, 0x41, 0x84, 0x41,
	0x82, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x48, 0x22, 0x08, 0xf9, 0x90, 0xca, 0x82, 0x54,
	0x09, 0x26, 0x64, 0x79, 0x90, 0x88, 0x52, 0x28, 0x97, 0xb0, 0x4b, 0x6a, 0x4e, 0x6a, 0x49, 0x2a,
	0x44, 0x8c, 0x5a, 0xc6, 0x5a, 0x72, 0x89, 0xc3, 0x1d, 0x5c, 0xec, 0x54, 0x09, 0x12, 0xc3, 0x30,
	0x1a, 0xac, 0x95, 0x11, 0x43, 0xab, 0x36, 0x97, 0x20, 0x5c, 0x6b, 0x50, 0x6a, 0x71, 0x41, 0x7e,
	0x5e, 0x72, 0xaa, 0x90, 0x18, 0x17, 0x1b, 0x44, 0x04, 0xac, 0x81, 0x27, 0x08, 0xca, 0x53, 0x72,
	0x80, 0x89, 0xe3, 0x52, 0x41, 0xd0, 0xa5, 0x4a, 0x5c, 0x1c, 0xc8, 0xb6, 0x04, 0x97, 0x24, 0x96,
	0x94, 0x16, 0x43, 0x9d, 0x05, 0xe5, 0x19, 0x6d, 0x61, 0xe2, 0xe2, 0x85, 0x68, 0x09, 0x4e, 0x2d,
	0x2a, 0xcb, 0x4c, 0x4e, 0x15, 0x72, 0xe3, 0xe2, 0x47, 0x8b, 0x10, 0x21, 0x69, 0xbd, 0xc4, 0x82,
	0x4c, 0x3d, 0xec, 0xd1, 0x24, 0x25, 0x86, 0x2a, 0x09, 0xb7, 0xd1, 0x8b, 0x4b, 0x00, 0x3d, 0x9c,
	0x84, 0x64, 0x50, 0xd5, 0xa2, 0x06, 0x1f, 0x2e, 0x93, 0x0c, 0x18, 0x85, 0xb4, 0xb8, 0x78, 0x9c,
	0x8b, 0x52, 0x13, 0x61, 0x51, 0x29, 0xc4, 0x0d, 0x56, 0x09, 0xe1, 0x48, 0xf1, 0x82, 0x39, 0x70,
	0x7b, 0xcd, 0xb9, 0x78, 0x90, 0xa3, 0x5d, 0x48, 0x02, 0x2c, 0x8d, 0x25, 0x25, 0xa0, 0x6b, 0xd4,
	0xe2, 0xe2, 0x09, 0x2d, 0x48, 0x21, 0xca, 0x92, 0x24, 0x36, 0x70, 0x0a, 0x36, 0x06, 0x04, 0x00,
	0x00, 0xff, 0xff, 0x70, 0x8d, 0xa1, 0xb3, 0xd6, 0x02, 0x00, 0x00,
}
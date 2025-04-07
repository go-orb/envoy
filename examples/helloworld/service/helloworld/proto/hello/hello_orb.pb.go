// Code generated by protoc-gen-go-orb. DO NOT EDIT.
//
// version:
// - protoc-gen-go-orb        v0.0.1
// - protoc                   v6.30.1
//
// Proto source: hello/hello.proto

package hello_v1

import (
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/server"

	"google.golang.org/protobuf/proto"
	"storj.io/drpc"

	grpc "google.golang.org/grpc"

	mdrpc "github.com/go-orb/plugins/server/drpc"
	memory "github.com/go-orb/plugins/server/memory"

	mhttp "github.com/go-orb/plugins/server/http"
)

// HandlerHello is the name of a service, it's here to static type/reference.
const HandlerHello = "hello.v1.Hello"
const EndpointHelloHello = "/hello.v1.Hello/Hello"

// orbEncoding_Hello_proto is a protobuf encoder for the hello.v1.Hello service.
type orbEncoding_Hello_proto struct{}

// Marshal implements the drpc.Encoding interface.
func (orbEncoding_Hello_proto) Marshal(msg drpc.Message) ([]byte, error) {
	m, ok := msg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("message is not a proto.Message: %T", msg)
	}
	return proto.Marshal(m)
}

// Unmarshal implements the drpc.Encoding interface.
func (orbEncoding_Hello_proto) Unmarshal(data []byte, msg drpc.Message) error {
	m, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("message is not a proto.Message: %T", msg)
	}
	return proto.Unmarshal(data, m)
}

// Name implements the drpc.Encoding interface.
func (orbEncoding_Hello_proto) Name() string {
	return "proto"
}

// HelloClient is the client for hello.v1.Hello
type HelloClient struct {
	client client.Client
}

// NewHelloClient creates a new client for hello.v1.Hello
func NewHelloClient(client client.Client) *HelloClient {
	return &HelloClient{client: client}
}

// Hello requests Hello.
func (c *HelloClient) Hello(ctx context.Context, service string, req *emptypb.Empty, opts ...client.CallOption) (*HelloResponse, error) {
	return client.Request[HelloResponse](ctx, c.client, service, EndpointHelloHello, req, opts...)
}

// HelloHandler is the Handler for hello.v1.Hello
type HelloHandler interface {
	Hello(ctx context.Context, req *emptypb.Empty) (*HelloResponse, error)
}

// orbGRPCHello provides the adapter to convert a HelloHandler to a gRPC HelloServer.
type orbGRPCHello struct {
	handler HelloHandler
}

// Hello implements the HelloServer interface by adapting to the HelloHandler.
func (s *orbGRPCHello) Hello(ctx context.Context, req *emptypb.Empty) (*HelloResponse, error) {
	return s.handler.Hello(ctx, req)
}

// Stream adapters to convert gRPC streams to ORB streams.

// Verification that our adapters implement the required interfaces.
var _ HelloServer = (*orbGRPCHello)(nil)

// registerHelloGRPCServerHandler registers the service to a gRPC server.
func registerHelloGRPCServerHandler(srv grpc.ServiceRegistrar, handler HelloHandler) {
	// Create the adapter to convert from HelloHandler to HelloServer
	grpcHandler := &orbGRPCHello{handler: handler}

	srv.RegisterService(&Hello_ServiceDesc, grpcHandler)
}

// orbDRPCHelloHandler wraps a HelloHandler to implement DRPCHelloServer.
type orbDRPCHelloHandler struct {
	handler HelloHandler
}

// Hello implements the DRPCHelloServer interface by adapting to the HelloHandler.
func (w *orbDRPCHelloHandler) Hello(ctx context.Context, req *emptypb.Empty) (*HelloResponse, error) {
	return w.handler.Hello(ctx, req)
}

// Stream adapters to convert DRPC streams to ORB streams.

// Verification that our adapters implement the required interfaces.
var _ DRPCHelloServer = (*orbDRPCHelloHandler)(nil)

// registerHelloDRPCHandler registers the service to an dRPC server.
func registerHelloDRPCHandler(srv *mdrpc.Server, handler HelloHandler) error {
	desc := DRPCHelloDescription{}

	// Wrap the ORB handler with our adapter to make it compatible with DRPC.
	drpcHandler := &orbDRPCHelloHandler{handler: handler}

	// Register with the server/drpc(.Mux).
	err := srv.Router().Register(drpcHandler, desc)
	if err != nil {
		return err
	}

	// Add each endpoint name of this handler to the orb drpc server.
	srv.AddEndpoint("/hello.v1.Hello/Hello")

	return nil
}

// registerHelloMemoryHandler registers the service to a memory server.
func registerHelloMemoryHandler(srv *memory.Server, handler HelloHandler) error {
	desc := DRPCHelloDescription{}

	// Wrap the ORB handler with our adapter to make it compatible with DRPC.
	drpcHandler := &orbDRPCHelloHandler{handler: handler}

	// Register with the server/drpc(.Mux).
	err := srv.Router().Register(drpcHandler, desc)
	if err != nil {
		return err
	}

	// Add each endpoint name of this handler to the orb drpc server.
	srv.AddEndpoint("/hello.v1.Hello/Hello")

	return nil
}

// registerHelloHTTPHandler registers the service to an HTTP server.
func registerHelloHTTPHandler(srv *mhttp.Server, handler HelloHandler) {
	srv.Router().Post("/hello.v1.Hello/Hello", mhttp.NewGRPCHandler(srv, handler.Hello, HandlerHello, "Hello"))
}

// RegisterHelloHandler will return a registration function that can be
// provided to entrypoints as a handler registration.
func RegisterHelloHandler(handler any) server.RegistrationFunc {
	return func(s any) {
		switch srv := s.(type) {

		case grpc.ServiceRegistrar:
			registerHelloGRPCServerHandler(srv, handler.(HelloHandler))
		case *mdrpc.Server:
			registerHelloDRPCHandler(srv, handler.(HelloHandler))
		case *memory.Server:
			registerHelloMemoryHandler(srv, handler.(HelloHandler))
		case *mhttp.Server:
			registerHelloHTTPHandler(srv, handler.(HelloHandler))
		default:
			log.Warn("No provider for this server found", "proto", "hello/hello.proto", "handler", "Hello", "server", s)
		}
	}
}

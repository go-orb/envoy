// Package echo provdes a echo handler.
package echo

import (
	"context"

	helloV1Proto "github.com/go-orb/envoy/examples/helloworld/service/helloworld/proto/hello"
	"github.com/go-orb/go-orb/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ helloV1Proto.HelloHandler = (*Handler)(nil)

// Handler is a test handler.
type Handler struct {
	logger log.Logger
}

// New creates a new Handler.
func New(logger log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Hello implements the Hello method.
func (h *Handler) Hello(_ context.Context, _ *emptypb.Empty) (*helloV1Proto.HelloResponse, error) {
	resp := &helloV1Proto.HelloResponse{
		Message: "Hello World",
	}

	return resp, nil
}

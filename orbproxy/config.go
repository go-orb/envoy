// Package main contains the orbproxy filter.
package main

import (
	"errors"
	"fmt"

	xds "github.com/cncf/xds/go/xds/type/v3"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/http"

	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb_transport/drpc"
	_ "github.com/go-orb/plugins/client/orb_transport/grpc"
	_ "github.com/go-orb/plugins/client/orb_transport/http"

	_ "github.com/go-orb/envoy/envoylog"

	_ "github.com/go-orb/plugins/codecs/form"
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/msgpack"
	_ "github.com/go-orb/plugins/codecs/proto"

	_ "github.com/go-orb/plugins/kvstore/natsjs"

	_ "github.com/go-orb/plugins/registry/kvstore"
)

const name = "orbproxy"

func init() {
	http.RegisterHttpFilterFactoryAndConfigParser(name, filterFactory, &parser{})
}

type config struct {
	service  string
	endpoint string
	stream   bool

	config string
}

type parser struct {
}

// Parse the filter configuration. We can call the ConfigCallbackHandler to control the filter's
// behavior
func (p *parser) Parse(any *anypb.Any, _ api.ConfigCallbackHandler) (any, error) {
	configStruct := &xds.TypedStruct{}
	if err := any.UnmarshalTo(configStruct); err != nil {
		return nil, err
	}

	v := configStruct.Value
	conf := &config{}
	service, ok := v.AsMap()["service"]
	if !ok {
		return nil, errors.New("missing service")
	}
	if str, ok := service.(string); ok {
		conf.service = str
	} else {
		return nil, fmt.Errorf("service: expect string while got %T", service)
	}

	endpoint, ok := v.AsMap()["endpoint"]
	if !ok {
		return nil, errors.New("missing endpoint")
	}
	if str, ok := endpoint.(string); ok {
		conf.endpoint = str
	} else {
		return nil, fmt.Errorf("endpoint: expect string while got %T", endpoint)
	}

	stream, ok := v.AsMap()["stream"]
	if !ok {
		return nil, errors.New("missing stream")
	}
	if streamBool, ok := stream.(bool); ok {
		conf.stream = streamBool
	}

	config, ok := v.AsMap()["config"]
	if !ok {
		return nil, errors.New("missing config")
	}
	if str, ok := config.(string); ok {
		conf.config = str
	} else {
		return nil, fmt.Errorf("config: expect string while got %T", config)
	}

	return conf, nil
}

// Merge configuration from the inherited parent configuration
func (p *parser) Merge(parent any, child any) any {
	parentConfig := parent.(*config)
	childConfig := child.(*config)

	// copy one, do not update parentConfig directly.
	newConfig := *parentConfig

	if childConfig.service != "" {
		newConfig.service = childConfig.service
	}

	if childConfig.endpoint != "" {
		newConfig.endpoint = childConfig.endpoint
	}

	newConfig.stream = childConfig.stream

	if childConfig.config != "" {
		newConfig.config = childConfig.config
	}

	return &newConfig
}

func main() {}

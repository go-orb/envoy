module github.com/go-orb/envoy/orbproxy

go 1.24

require (
	github.com/cncf/xds/go v0.0.0-20250326154945-ae57f3c0d45f
	github.com/envoyproxy/envoy v1.33.2
	github.com/go-orb/envoy/envoylog v0.0.0-20250407202818-54fbdebb0420
	github.com/go-orb/go-orb v0.4.1
	github.com/go-orb/plugins/client/orb v0.2.1
	github.com/go-orb/plugins/client/orb_transport/drpc v0.1.0
	github.com/go-orb/plugins/client/orb_transport/grpc v0.1.0
	github.com/go-orb/plugins/client/orb_transport/http v0.3.0
	github.com/go-orb/plugins/codecs/form v0.2.0
	github.com/go-orb/plugins/codecs/goccyjson v0.2.0
	github.com/go-orb/plugins/codecs/msgpack v0.1.0
	github.com/go-orb/plugins/codecs/proto v0.2.0
	github.com/go-orb/plugins/kvstore/natsjs v0.2.1
	github.com/go-orb/plugins/registry/kvstore v0.2.2
	google.golang.org/protobuf v1.36.6
)

require (
	cel.dev/expr v0.23.1 // indirect
	dario.cat/mergo v1.0.1 // indirect
	github.com/cornelk/hashmap v1.0.8 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/go-orb/plugins/registry/regutil v0.2.0 // indirect
	github.com/go-orb/plugins/server/drpc v0.2.1 // indirect
	github.com/go-playground/form/v4 v4.2.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/google/pprof v0.0.0-20250403155104-27863c87afa6 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nats.go v1.41.0 // indirect
	github.com/nats-io/nkeys v0.4.10 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/onsi/ginkgo/v2 v2.23.4 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.50.1 // indirect
	github.com/shamaton/msgpack/v2 v2.2.3 // indirect
	github.com/zeebo/errs v1.4.0 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	go.uber.org/mock v0.5.1 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250407143221-ac9807e6c755 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250407143221-ac9807e6c755 // indirect
	google.golang.org/grpc v1.71.1 // indirect
	storj.io/drpc v0.0.34 // indirect
)

replace github.com/go-orb/envoy/envoylog => ../envoylog

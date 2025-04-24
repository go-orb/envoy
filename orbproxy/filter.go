package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/orberrors"
)

var _ api.StreamFilter = (*filter)(nil)

var globalClientOnce sync.Once
var globalClient client.Type

type filter struct {
	api.PassThroughStreamFilter

	callbacks api.FilterCallbackHandler
	config    *config

	requestHeaders map[string]string

	responseHeaders map[string]string
	responseBody    []byte

	responseError *orberrors.Error
}

func filterFactory(c any, callbacks api.FilterCallbackHandler) api.StreamFilter {
	conf, ok := c.(*config)
	if !ok {
		panic("unexpected config type")
	}
	f := &filter{
		callbacks:       callbacks,
		config:          conf,
		requestHeaders:  make(map[string]string),
		responseHeaders: make(map[string]string),
		responseBody:    make([]byte, 0),
	}

	codec, err := codecs.GetMime(codecs.MimeJSON)
	if err != nil {
		panic(err)
	}

	globalClientOnce.Do(func() {
		components := types.NewComponents()
		cfg := map[string]any{}
		if err := codec.Unmarshal([]byte(f.config.config), &cfg); err != nil {
			panic(err)
		}

		if _, ok := cfg["logger"]; !ok {
			cfg["logger"] = map[string]any{
				"plugin": "envoylog",
			}
		} else if _, ok := cfg["logger"].(map[string]any)["plugin"]; !ok {
			cfg["logger"].(map[string]any)["plugin"] = "envoylog"
		}

		logger, err := log.NewConfigDatas([]string{}, cfg)
		if err != nil {
			panic(err)
		}

		registry, err := registry.New(cfg, components, logger)
		if err != nil {
			panic(err)
		}

		globalClient, err = client.New(cfg, components, logger, registry)
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		if err := globalClient.Start(ctx); err != nil {
			panic(err)
		}
		if err := registry.Start(ctx); err != nil {
			panic(err)
		}
	})

	return f
}

func (f *filter) request(reqBytes []byte) {
	defer f.callbacks.DecoderFilterCallbacks().RecoverPanic()

	contentType, ok := f.requestHeaders["content-type"]
	if !ok {
		contentType = codecs.MimeJSON
	}

	api.LogDebugf("Requesting %s%s", f.config.service, f.config.endpoint)

	bPointer, err := client.Request[[]byte](
		context.Background(),
		globalClient,
		f.config.service,
		f.config.endpoint,
		reqBytes,
		client.WithContentType(contentType),
		client.WithMetadata(f.requestHeaders),
		client.WithResponseMetadata(f.responseHeaders),
	)

	if err != nil {
		f.responseError = orberrors.From(err)
		api.LogError(f.responseError.Error())
		f.callbacks.DecoderFilterCallbacks().Continue(api.Continue)
		return
	}

	b := *bPointer

	accept, ok := f.requestHeaders["accept"]
	if !ok || accept == "*/*" || accept == "" {
		accept = codecs.MimeJSON
	}

	f.responseHeaders["Content-Type"] = accept
	f.responseHeaders["Content-Length"] = fmt.Sprintf("%d", len(b))

	f.responseBody = b

	f.callbacks.DecoderFilterCallbacks().Continue(api.Continue)
}

// Callbacks which are called in request path
// The endStream is true if the request doesn't have body
func (f *filter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	api.LogDebug("Decoding headers")

	// Copy the headers from the request.
	if len(f.requestHeaders) == 0 {
		header.RangeWithCopy(func(key, value string) bool {
			// Skip pseudo headers
			if key[0:1] == ":" {
				return true
			}

			f.requestHeaders[strings.ToLower(key)] = value
			return true
		})
	}

	if endStream {
		// If the request is empty, request with empty JSON object
		api.LogDebug("Request is empty, requesting empty JSON object")
		go f.request([]byte("{}"))

		return api.Running
	}

	return api.Continue
}

// DecodeData might be called multiple times during handling the request body.
// The endStream is true when handling the last piece of the body.
func (f *filter) DecodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	api.LogDebug("Decoding data")

	go f.request(buffer.Bytes())

	return api.Running
}

// Callbacks which are called in response path
// The endStream is true if the response doesn't have body
func (f *filter) EncodeHeaders(header api.ResponseHeaderMap, endStream bool) api.StatusType {
	api.LogDebug("Encoding headers")

	if f.responseError != nil {
		api.LogError(f.responseError.Error())

		code := f.responseError.Code
		if code == 0 {
			code = 500
		}

		body := fmt.Sprintf("%s\r\n", f.responseError.Error())

		f.callbacks.EncoderFilterCallbacks().SendLocalReply(code, body, nil, 0, "")
		return api.LocalReply
	}

	for k, v := range f.responseHeaders {
		header.Set(k, v)
	}

	return api.Continue
}

// EncodeData might be called multiple times during handling the response body.
// The endStream is true when handling the last piece of the body.
func (f *filter) EncodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	api.LogDebugf("Encoding data, responseError: %v, endStream: %v, responseBody: %v", f.responseError != nil, endStream, string(f.responseBody))

	if f.responseError != nil {
		return api.Continue
	}

	buffer.Reset()
	buffer.Write(f.responseBody)

	return api.Continue
}

func (f *filter) EncodeTrailers(trailers api.ResponseTrailerMap) api.StatusType {
	return api.Continue
}

// OnLog is called when the HTTP stream is ended on HTTP Connection Manager filter.
func (f *filter) OnLog(reqHeader api.RequestHeaderMap, reqTrailer api.RequestTrailerMap, respHeader api.ResponseHeaderMap, respTrailer api.ResponseTrailerMap) {
}

// OnLogDownstreamStart is called when HTTP Connection Manager filter receives a new HTTP request
// (required the corresponding access log type is enabled)
func (f *filter) OnLogDownstreamStart(reqHeader api.RequestHeaderMap) {
	// also support kicking off a goroutine here, like OnLog.
}

// OnLogDownstreamPeriodic is called on any HTTP Connection Manager periodic log record
// (required the corresponding access log type is enabled)
func (f *filter) OnLogDownstreamPeriodic(reqHeader api.RequestHeaderMap, reqTrailer api.RequestTrailerMap, respHeader api.ResponseHeaderMap, respTrailer api.ResponseTrailerMap) {
	// also support kicking off a goroutine here, like OnLog.
}

func (f *filter) OnDestroy(reason api.DestroyReason) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Printf("http: panic serving: %v\n%s", p, buf)
			}
		}()
	}()
}

func (f *filter) OnStreamComplete() {}

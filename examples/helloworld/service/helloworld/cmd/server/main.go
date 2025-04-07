// Package main contains a fake login server.
package main

import (
	"fmt"
	"os"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/server"

	handler "github.com/go-orb/envoy/examples/helloworld/service/helloworld/handler/hello"
	proto "github.com/go-orb/envoy/examples/helloworld/service/helloworld/proto/hello"

	mhttp "github.com/go-orb/plugins/server/http"

	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"

	_ "github.com/go-orb/plugins/kvstore/natsjs"
	_ "github.com/go-orb/plugins/registry/kvstore"
)

// provideServerOpts provides options for the go-orb server.
//
//nolint:unparam
func provideServerOpts() ([]server.ConfigOption, error) {
	opts := []server.ConfigOption{}

	opts = append(opts, server.WithEntrypointConfig(
		"http",
		mhttp.NewConfig(
			mhttp.WithInsecure(),
			mhttp.WithAddress(":10000"),
		),
	))

	return opts, nil
}

// provideServerConfigured configures the go-orb server(s).
//
//nolint:unparam
func provideServerConfigured(logger log.Logger, srv server.Server) (serverConfigured, error) {
	// Register server Handlers.
	hInstance := handler.New(logger)
	hRegister := proto.RegisterHelloHandler(hInstance)

	// Add our server handler to all entrypoints.
	srv.GetEntrypoints().Range(func(_ string, entrypoint server.Entrypoint) bool {
		entrypoint.AddHandler(hRegister)

		return true
	})

	return serverConfigured{}, nil
}

func runner(
	svcCtx *cli.ServiceContextWithConfig,
	logger log.Logger,
) error {
	logger.Info("Started", "name", svcCtx.Name(), "version", svcCtx.Version())

	// Blocks until the process receives a signal.
	<-svcCtx.Context().Done()

	logger.Info("Stopping", "name", svcCtx.Name(), "version", svcCtx.Version())

	return nil
}

func main() {
	app := cli.App{
		Name:     "helloworld",
		Version:  "",
		Usage:    "Hello world service",
		NoAction: false,
		Flags: []*cli.Flag{
			{
				Name:        "log-level",
				Default:     "INFO",
				EnvVars:     []string{"LOG_LEVEL"},
				ConfigPaths: [][]string{{"logger", "level"}},
				Usage:       "Set the log level, one of TRACE, DEBUG, INFO, WARN, ERROR",
			},
		},
		Commands: []*cli.Command{},

		NoMultiServiceConfig: true,
		HardcodedConfigs:     []cli.HardcodedConfig{{Format: "json", Data: config}},
	}

	appContext := cli.NewAppContext(&app)

	_, err := run(appContext, os.Args, runner)
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}

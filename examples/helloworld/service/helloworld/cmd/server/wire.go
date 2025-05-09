//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/cli/urfave"

	"github.com/go-orb/wire"
)

// wireRunResult is here so "wire" has a type for the return value of wireRun.
// wire needs a explicit type for each provider including "wireRun".
type wireRunResult string

// wireRunCallback is the actual code that runs the business logic.
type wireRunCallback func(
	svcCtx *cli.ServiceContextWithConfig,
	logger log.Logger,
) error

type serverConfigured struct{}

func wireRun(
	serviceContext *cli.ServiceContextWithConfig,
	components *types.Components,
	logger log.Logger,
	_ serverConfigured,
	cb wireRunCallback,
) (wireRunResult, error) {
	// Orb start
	for _, c := range components.Iterate(false) {
		err := c.Start(serviceContext.Context())
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return "", err
		}
	}

	//
	// Actual code
	runErr := cb(serviceContext, logger)

	// Orb shutdown.
	ctx := context.Background()

	for _, c := range components.Iterate(true) {
		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return "", runErr
}

func run(
	appContext *cli.AppContext,
	args []string,
	cb wireRunCallback,
) (wireRunResult, error) {
	panic(wire.Build(
		urfave.ProvideParser,
		cli.ProvideParsedFlagsFromArgs,

		cli.ProvideAppConfigData,
		cli.ProvideServiceConfigData,

		cli.ProvideSingleServiceContext,
		types.ProvideComponents,

		log.ProvideNoOpts,
		registry.ProvideNoOpts,

		provideServerConfigured,
		provideServerOpts,
		server.Provide,

		wireRun,
	))
}

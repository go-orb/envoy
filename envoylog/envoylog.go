// Package envoylog provides the envoylog handler.
package envoylog

import (
	"context"
	"errors"
	"fmt"
	"os"

	"log/slog"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
)

// Name is this providers name.
const Name = "envoylog"

// The register.
func init() {
	log.Register(Name, Factory)
}

// Config is the config struct for slog.
type Config struct {
	log.Config
}

// NewConfig creates a new config.
func NewConfig(opts ...log.Option) Config {
	cfg := Config{
		Config: log.NewConfig(),
	}

	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

var _ (log.Provider) = (*Provider)(nil)

// Provider is the provider for slog.
type Provider struct {
	config Config

	file    *os.File
	handler slog.Handler
}

// Start configures the slog Handler.
func (p *Provider) Start() error {

	p.handler = &Handler{}

	return nil
}

// Stop closes if required a open log file.
func (p *Provider) Stop(_ context.Context) error {
	if p.file != nil {
		return p.file.Close()
	}

	return nil
}

// Handler returns the configure handler.
func (p *Provider) Handler() (slog.Handler, error) {
	return p.handler, nil
}

// Key returns an identifier for this handler provider with its config.
func (p *Provider) Key() string {
	return fmt.Sprintf("__%s__", Name)
}

// Factory is the factory for a slog provider.
func Factory(sections []string, configs map[string]any, opts ...log.Option) (log.ProviderType, error) {
	cfg := NewConfig(opts...)

	if err := config.Parse(sections, "logger", configs, &cfg); err != nil && !errors.Is(err, config.ErrNoSuchKey) {
		return log.ProviderType{}, err
	}

	return log.ProviderType{
		Provider: &Provider{
			config: cfg,
		},
	}, nil
}

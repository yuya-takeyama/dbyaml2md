package main

import (
	"github.com/codegangsta/cli"
)

type Context interface {
	OutDirectory() string
	Config() *Config
}

type AppContext struct {
	app    *cli.Context
	config *Config
}

func NewAppContext(app *cli.Context, config *Config) Context {
	return &AppContext{
		app,
		config,
	}
}

func (context *AppContext) OutDirectory() string {
	if context.app.IsSet("out") {
		return context.app.String("out")
	}

	if context.config.Out != "" {
		return context.config.Out
	} else {
		return "./dbyaml2md_out"
	}
}

func (context *AppContext) Config() *Config {
	return context.config
}

type EmptyContext struct {
	config *Config
}

func NewEmptyContext() Context {
	return &EmptyContext{
		NewEmptyConfig(),
	}
}

func (context *EmptyContext) OutDirectory() string {
	return "./empty_context_out"
}

func (context *EmptyContext) Config() *Config {
	return context.config
}

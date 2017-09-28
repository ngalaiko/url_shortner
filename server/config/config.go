package config

import (
	"context"
	"flag"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/ngalayko/url_shortner/server/logger"
)

const (
	ctxKey configCtxKey = "config_ctx_key"
)

type configCtxKey string

// Config is a application config struct
type Config struct {
	Db  DbConfig  `yaml:"Db"`
	Web WebConfig `yaml:"Web"`
	Facebook FacebookConfig `yaml:"Facebook"`
}

var (
	configPath = flag.String("config", "", "path to app Config")
)

// NewContext stores config in context
func NewContext(ctx context.Context, config interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := config.(*Config); !ok {
		config = newConfig(ctx)
	}

	return context.WithValue(ctx, ctxKey, config)
}

// FromContext returns config from context
func FromContext(ctx context.Context) *Config {
	if config, ok := ctx.Value(ctxKey).(*Config); ok {
		return config
	}

	return newConfig(ctx)
}

func newConfig(ctx context.Context) *Config {
	l := logger.FromContext(ctx)

	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		l.Panic("error while open Config file",
			zap.String("file", *configPath),
			zap.String("err", err.Error()),
		)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		l.Panic("error while reading Config file",
			zap.String("file", file.Name()),
			zap.String("err", err.Error()),
		)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		l.Panic("error while unmarshal yaml Config",
			zap.Error(err),
		)
	}

	l.Info("config parsed")
	return config
}

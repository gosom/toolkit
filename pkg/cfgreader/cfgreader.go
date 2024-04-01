package cfgreader

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig returns a new Config readinf from env variables.
func NewConfig[C any](prefix string) (C, error) {
	var cfg C
	err := envconfig.Process(prefix, &cfg)

	return cfg, err
}

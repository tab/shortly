package config

import (
	"shortly/internal/app/helpers"
	"shortly/internal/app/store"
)

type AppConfig struct {
	SecureRandom helpers.SecureRandomGenerator
	Store        *store.URLStore
	Flags        Flags
}

func NewAppConfig(SecureRandom helpers.SecureRandomGenerator, Store *store.URLStore, flags Flags) *AppConfig {
	return &AppConfig{
		SecureRandom: SecureRandom,
		Store:        Store,
		Flags:        flags,
	}
}

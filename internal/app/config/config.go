package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"shortly/internal/logger"
)

// BaseURL is the base URL of the application
const BaseURL = "http://localhost:8080"

// ServerAddress is the address and port to run the server
const ServerAddress = "localhost:8080"

// ProfilerAddress is the address and port to run the profiler
const ProfilerAddress = "localhost:2080"

// Config is the application configuration
type Config struct {
	AppEnv          string `json:"env"`
	Addr            string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	ClientURL       string `json:"client_url"`
	ProfilerAddr    string `json:"profiler_address"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	SecretKey       string `json:"secret_key"`
	EnableHTTPS     bool   `json:"enable_https"`
	Certificate     string `json:"certificate_path"`
	PrivateKey      string `json:"certificate_key_path"`
	ConfigFilePath  string
}

// Flags is the flags for the configuration
type Flags struct {
	Addr            string
	BaseURL         string
	ProfilerAddr    string
	FileStoragePath string
	DatabaseDSN     string
	SecretKey       string
	EnableHTTPS     bool
	ConfigFilePath  string
}

// Builder is a builder for the Config
type Builder struct {
	cfg *Config
}

// NewConfigBuilder creates a new Builder instance
func NewConfigBuilder() *Builder {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	return &Builder{
		cfg: &Config{
			AppEnv: env,
		},
	}
}

// ParseFlags parses application flags
func ParseFlags() Flags {
	flagAddr := flag.String("a", ServerAddress, "address and port to run server")
	flagBaseURL := flag.String("b", BaseURL, "base address of the resulting shortened URL")
	flagProfilerAddr := flag.String("p", ProfilerAddress, "address and port to run profiler")
	flagFileStoragePath := flag.String("f", "", "path to the file storage")
	flagDatabaseDSN := flag.String("d", "", "database DSN")
	flagSecretKey := flag.String("k", "", "JWT secret key")
	flagConfigFilePath := flag.String("c", "", "path to the config file")
	flagAliasConfigFilePath := flag.String("config", "", "path to the config file")
	flag.Parse()

	if *flagConfigFilePath == "" {
		*flagConfigFilePath = *flagAliasConfigFilePath
	}

	return Flags{
		Addr:            *flagAddr,
		BaseURL:         *flagBaseURL,
		ProfilerAddr:    *flagProfilerAddr,
		FileStoragePath: *flagFileStoragePath,
		DatabaseDSN:     *flagDatabaseDSN,
		SecretKey:       *flagSecretKey,
		ConfigFilePath:  *flagConfigFilePath,
	}
}

// LoadConfig loads the application configuration, priority: env > flags > file
func LoadConfig() *Config {
	appLogger := logger.NewLogger()

	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	envFiles := []string{
		".env",
		fmt.Sprintf(".env.%s", env),
		fmt.Sprintf(".env.%s.local", env),
	}
	for _, file := range envFiles {
		err := godotenv.Overload(file)
		if err == nil {
			appLogger.Info().Msgf("Loaded %s file", file)
		}
	}

	flags := ParseFlags()

	builder := NewConfigBuilder()
	builder.cfg.ConfigFilePath = flags.ConfigFilePath

	cfg := builder.
		WithFile().
		WithFlags(flags).
		WithEnv().
		Build()

	return cfg
}

// WithFile loads the configuration from the JSON file
func (b *Builder) WithFile() *Builder {
	file, err := os.Open(b.cfg.ConfigFilePath)
	if err != nil {
		return b
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(b.cfg)
	if err != nil {
		return b
	}

	return b
}

// WithFlags loads the configuration from the flags
func (b *Builder) WithFlags(f Flags) *Builder {
	if f.ConfigFilePath != "" {
		b.cfg.ConfigFilePath = f.ConfigFilePath
	}
	if f.Addr != "" {
		b.cfg.Addr = f.Addr
	}
	if f.BaseURL != "" {
		b.cfg.BaseURL = f.BaseURL
	}
	if f.ProfilerAddr != "" {
		b.cfg.ProfilerAddr = f.ProfilerAddr
	}
	if f.FileStoragePath != "" {
		b.cfg.FileStoragePath = f.FileStoragePath
	}
	if f.DatabaseDSN != "" {
		b.cfg.DatabaseDSN = f.DatabaseDSN
	}
	if f.SecretKey != "" {
		b.cfg.SecretKey = f.SecretKey
	}
	b.cfg.EnableHTTPS = f.EnableHTTPS

	return b
}

// WithEnv loads the configuration from the environment variables
func (b *Builder) WithEnv() *Builder {
	if v, ok := os.LookupEnv("SERVER_ADDRESS"); ok && v != "" {
		b.cfg.Addr = v
	}
	if v, ok := os.LookupEnv("BASE_URL"); ok && v != "" {
		b.cfg.BaseURL = v
	}
	if v, ok := os.LookupEnv("CLIENT_URL"); ok && v != "" {
		b.cfg.ClientURL = v
	}
	if v, ok := os.LookupEnv("PROFILER_ADDRESS"); ok && v != "" {
		b.cfg.ProfilerAddr = v
	}
	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok && v != "" {
		b.cfg.FileStoragePath = v
	}
	if v, ok := os.LookupEnv("DATABASE_DSN"); ok && v != "" {
		b.cfg.DatabaseDSN = v
	}
	if v, ok := os.LookupEnv("SECRET_KEY"); ok && v != "" {
		b.cfg.SecretKey = v
	}
	if v, ok := os.LookupEnv("ENABLE_HTTPS"); ok && v != "" {
		b.cfg.EnableHTTPS = (v == "true")
	}
	if v, ok := os.LookupEnv("CERTIFICATE_PATH"); ok && v != "" {
		b.cfg.Certificate = v
	}
	if v, ok := os.LookupEnv("CERTIFICATE_KEY_PATH"); ok && v != "" {
		b.cfg.PrivateKey = v
	}
	if v, ok := os.LookupEnv("CONFIG"); ok && v != "" {
		b.cfg.ConfigFilePath = v
	}

	return b
}

// Build builds the Config
func (b *Builder) Build() *Config {
	return b.cfg
}

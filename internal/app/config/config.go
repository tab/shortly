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

// GRPCServerAddress is the address and port to run the gRPC server
const GRPCServerAddress = "localhost:50051"

// GRPCGatewayAddress is the address and port to run the gRPC gateway
const GRPCGatewayAddress = "localhost:8081"

// ProfilerAddress is the address and port to run the profiler
const ProfilerAddress = "localhost:2080"

// Config is the application configuration
type Config struct {
	AppEnv          string `json:"env"`
	Addr            string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	ClientURL       string `json:"client_url"`
	GRPCServerAddr  string `json:"grpc_server_address"`
	GRPCSecretKey   string `json:"grpc_secret_key"`
	GRPCGatewayAddr string `json:"grpc_gateway_address"`
	ProfilerAddr    string `json:"profiler_address"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	SecretKey       string `json:"secret_key"`
	EnableHTTPS     bool   `json:"enable_https"`
	Certificate     string `json:"certificate_path"`
	PrivateKey      string `json:"certificate_key_path"`
	ConfigFilePath  string
	TrustedSubnet   string `json:"trusted_subnet"`
}

// Flags is the flags for the configuration
type Flags struct {
	Addr            string
	BaseURL         string
	GRPCServerAddr  string
	GRPCSecretKey   string
	GRPCGatewayAddr string
	ProfilerAddr    string
	FileStoragePath string
	DatabaseDSN     string
	SecretKey       string
	EnableHTTPS     bool
	ConfigFilePath  string
	TrustedSubnet   string
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
	flagGRPCServerAddr := flag.String("g", GRPCServerAddress, "address and port to run gRPC server")
	flagGRPCSecretKey := flag.String("s", "", "gRPC secret key")
	flagGRPCGatewayAddr := flag.String("w", GRPCGatewayAddress, "address and port to run gRPC gateway")
	flagProfilerAddr := flag.String("p", ProfilerAddress, "address and port to run profiler")
	flagFileStoragePath := flag.String("f", "", "path to the file storage")
	flagDatabaseDSN := flag.String("d", "", "database DSN")
	flagSecretKey := flag.String("k", "", "JWT secret key")
	flagConfigFilePath := flag.String("c", "", "path to the config file")
	flagAliasConfigFilePath := flag.String("config", "", "path to the config file")
	flagTrustedSubnet := flag.String("t", "", "trusted subnet")
	flag.Parse()

	if *flagConfigFilePath == "" {
		*flagConfigFilePath = *flagAliasConfigFilePath
	}

	return Flags{
		Addr:            *flagAddr,
		BaseURL:         *flagBaseURL,
		GRPCServerAddr:  *flagGRPCServerAddr,
		GRPCSecretKey:   *flagGRPCSecretKey,
		GRPCGatewayAddr: *flagGRPCGatewayAddr,
		ProfilerAddr:    *flagProfilerAddr,
		FileStoragePath: *flagFileStoragePath,
		DatabaseDSN:     *flagDatabaseDSN,
		SecretKey:       *flagSecretKey,
		ConfigFilePath:  *flagConfigFilePath,
		TrustedSubnet:   *flagTrustedSubnet,
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
	set := func(val string, target *string) {
		if val != "" {
			*target = val
		}
	}

	set(f.ConfigFilePath, &b.cfg.ConfigFilePath)
	set(f.Addr, &b.cfg.Addr)
	set(f.BaseURL, &b.cfg.BaseURL)
	set(f.GRPCServerAddr, &b.cfg.GRPCServerAddr)
	set(f.GRPCSecretKey, &b.cfg.GRPCSecretKey)
	set(f.GRPCGatewayAddr, &b.cfg.GRPCGatewayAddr)
	set(f.ProfilerAddr, &b.cfg.ProfilerAddr)
	set(f.FileStoragePath, &b.cfg.FileStoragePath)
	set(f.DatabaseDSN, &b.cfg.DatabaseDSN)
	set(f.SecretKey, &b.cfg.SecretKey)
	set(f.TrustedSubnet, &b.cfg.TrustedSubnet)

	b.cfg.EnableHTTPS = f.EnableHTTPS

	return b
}

// WithEnv loads the configuration from the environment variables
func (b *Builder) WithEnv() *Builder {
	set := func(envName string, target *string) {
		if v, ok := os.LookupEnv(envName); ok && v != "" {
			*target = v
		}
	}

	list := []struct {
		env    string
		target *string
	}{
		{"SERVER_ADDRESS", &b.cfg.Addr},
		{"BASE_URL", &b.cfg.BaseURL},
		{"CLIENT_URL", &b.cfg.ClientURL},
		{"GRPC_SERVER_ADDRESS", &b.cfg.GRPCServerAddr},
		{"GRPC_SECRET_KEY", &b.cfg.GRPCSecretKey},
		{"GRPC_GATEWAY_ADDRESS", &b.cfg.GRPCGatewayAddr},
		{"PROFILER_ADDRESS", &b.cfg.ProfilerAddr},
		{"FILE_STORAGE_PATH", &b.cfg.FileStoragePath},
		{"DATABASE_DSN", &b.cfg.DatabaseDSN},
		{"SECRET_KEY", &b.cfg.SecretKey},
		{"CONFIG", &b.cfg.ConfigFilePath},
		{"TRUSTED_SUBNET", &b.cfg.TrustedSubnet},
	}

	for _, item := range list {
		set(item.env, item.target)
	}

	if v, ok := os.LookupEnv("ENABLE_HTTPS"); ok && v != "" {
		b.cfg.EnableHTTPS = (v == "true")
	}
	set("CERTIFICATE_PATH", &b.cfg.Certificate)
	set("CERTIFICATE_KEY_PATH", &b.cfg.PrivateKey)

	return b
}

// Build builds the Config
func (b *Builder) Build() *Config {
	return b.cfg
}

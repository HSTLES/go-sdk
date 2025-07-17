// go-sdk/config/config.go
package core_config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	// names of env vars
	EnvAppNameVar    = "APP_NAME"
	EnvAppDomainsVar = "APP_DOMAINS"
	EnvAppEnvVar     = "APP_ENV"
	EnvDevEnvFileVar = "DEV_ENV_FILE"
)

// AppEnv is the value of APP_ENV
type AppEnv string

const (
	Development AppEnv = "development"
	Production  AppEnv = "production"
)

// CoreConfig holds only the truly common settings
type CoreConfig struct {
	AppName    string
	AppDomains string // comma-separated list of domains
	AppEnv     AppEnv

	TursoCoreDBURL     string
	TursoAuthToken     string
	TursoPlatformToken string
}

// loadEnvFiles applies .env → dev-override
func LoadEnvFiles() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env, using system env")
	}
	if AppEnv(os.Getenv(EnvAppEnvVar)) == Development {
		if devFile := os.Getenv(EnvDevEnvFileVar); devFile != "" {
			if err := godotenv.Overload(devFile); err != nil {
				log.Printf("failed loading %q: %v", devFile, err)
			}
		}
	}
}

// LoadCoreConfig reads APP_NAME, APP_ENV, TURSO_CORE_DB_URL and TURSO_HSTLES_KEY
// and fatals if any are missing.
func LoadCoreConfig() *CoreConfig {
	LoadEnvFiles()
	return &CoreConfig{
		AppName:            RequireEnv(EnvAppNameVar),
		AppDomains:         RequireEnvDefault(EnvAppDomainsVar, ""),
		AppEnv:             AppEnv(RequireEnv(EnvAppEnvVar)),
		TursoCoreDBURL:     RequireEnv("TURSO_CORE_DB_URL"),
		TursoAuthToken:     RequireEnv("TURSO_HSTLES_KEY"),
		TursoPlatformToken: RequireEnv("TURSO_PLATFORM_TOKEN"),
	}
}

// --- helpers for all apps ---

// requireEnv fatals if key is missing
func RequireEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("env var %s is required", key)
	return "" // unreachable
}

// requireEnvInt falls back to def if missing/invalid
func RequireEnvInt(key string, def int) int {
	if s := os.Getenv(key); s != "" {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
		log.Printf("invalid %s=%q, using %d", key, s, def)
	}
	return def
}

// requireEnvDefault falls back to def if missing
func RequireEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

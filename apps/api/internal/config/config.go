package config

import "os"

type Config struct {
	Port                   string
	Env                    string
	DatabaseURL            string
	RedisURL               string
	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string
	SentryDSN              string
}

func Load() *Config {
	return &Config{
		Port:                   getEnv("API_PORT", "8080"),
		Env:                    getEnv("API_ENV", "development"),
		DatabaseURL:            getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:54322/postgres"),
		RedisURL:               getEnv("REDIS_URL", "redis://localhost:6379"),
		SupabaseURL:            getEnv("SUPABASE_URL", "http://127.0.0.1:54321"),
		SupabaseAnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		SentryDSN:              getEnv("SENTRY_DSN", ""),
	}
}

func (c *Config) IsProd() bool {
	return c.Env == "production"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

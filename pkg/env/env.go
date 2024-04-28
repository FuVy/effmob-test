package env

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetOrPanicOnEmpty(env string) string {
	e := os.Getenv(env)
	if e == "" {
		panic(fmt.Sprintf("environment variable '%s' cannot be empty", env))
	}
	return e
}

func GetAsURLOrPanicOnEmpty(env string) string {
	host := GetOrPanicOnEmpty(env)

	_, err := url.ParseRequestURI(host)
	if err != nil {
		panic(err)
	}

	return strings.Trim(host, "/")
}
func WithDefault(env string, defaultValue string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultValue
	}

	return e
}

func Uint64WithDefault(env string, defaultValue uint64) uint64 {
	envValue := os.Getenv(env)
	if envValue == "" {
		return defaultValue
	}
	v, err := strconv.ParseUint(envValue, 10, 64)
	if err != nil {
		log.WithError(err).Fatalf("parse env '%s' as uint64, got value: %s", env, envValue)
	}

	return v
}

func Int64WithDefault(env string, defaultValue int64) int64 {
	envValue := os.Getenv(env)
	if envValue == "" {
		return defaultValue
	}
	v, err := strconv.ParseInt(envValue, 10, 64)
	if err != nil {
		log.WithError(err).Fatalf("parse env '%s' as int64, got value: %s", env, envValue)
	}

	return v
}

func StringWithDefault(env string, defaultValue string) string {
	envValue := os.Getenv(env)
	if envValue == "" {
		return defaultValue
	}

	return envValue
}

func MakePostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		GetOrPanicOnEmpty("POSTGRES_HOST"),
		GetOrPanicOnEmpty("POSTGRES_USER"),
		GetOrPanicOnEmpty("POSTGRES_PASSWORD"),
		GetOrPanicOnEmpty("POSTGRES_DB"),
		GetOrPanicOnEmpty("POSTGRES_PORT"),
		WithDefault("POSTGRES_SSLMODE", "prefer"),
		WithDefault("TZ", "UTC"),
	)
}

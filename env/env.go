package env

import (
	"os"
	"strings"
	"log"
)

var Env = getenvName("ENV", Dev) // TODO: set this up in production
var JWTsecret = getenvBytes("JWT_SECRET", "super-duper-secret-key") // TODO: set this up in production
var Port = getenv("PORT", "3000")
var DatabaseURL = getenv("DATABASE_URL", "postgres://localhost:5432")

func getenv(key string, default_ string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return default_
	}
	return value
}

func getenvName(key string, default_ EnvName) EnvName {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return default_
	}

	switch value {
	case string(Dev):
		return Dev
	case string(Prod):
		return Prod
	}

	log.Fatalf("unknown ENV name: %s", value)
	return Dev
}

func getenvBytes(key string, default_ string) []byte {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return []byte(default_)
	}
	return []byte(value)
}

type EnvName string

const (
	Dev EnvName = "dev"
	Prod EnvName = "prod"
)

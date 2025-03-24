package env

import (
	"os"
	"strings"
)

var Port = getenv("PORT", "3000")
var DatabaseURL = getenv("DATABASE_URL", "postgres://localhost:5432")

func getenv(key string, default_ string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return default_
	}
	return value
}

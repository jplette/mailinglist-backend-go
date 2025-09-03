package envcfg

import (
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func Value(key string) string {
	return os.Getenv(key)
}

func Exists(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func Values(key string) []string {
	return strings.Split(Value(key), ",")
}

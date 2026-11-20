package configuration

import (
	"os"
)

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// @Injectable
func Load() *Configuration {
	config := &Configuration{
		DataDir: GetEnvOrDefault("NAGARE_DATA_DIR", DefaultConfiguration.DataDir),
	}

	return config
}

package configuration

import (
	"os"
	"strings"
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
	allowOrigins := DefaultConfiguration.AllowOrigins
	if o := strings.Split(GetEnvOrDefault("NAGARE_ALLOW_ORIGINS", ""), ","); len(o) > 0 {
		allowOrigins = o
	}

	config := &Configuration{
		DataDir:      GetEnvOrDefault("NAGARE_DATA_DIR", DefaultConfiguration.DataDir),
		Port:         GetEnvOrDefault("NAGARE_PORT", DefaultConfiguration.Port),
		Host:         GetEnvOrDefault("NAGARE_HOST", DefaultConfiguration.Host),
		AllowOrigins: allowOrigins,
	}

	return config
}

package config

import "os"

func Address() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8000"
}

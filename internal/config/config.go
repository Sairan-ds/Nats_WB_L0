package config

import (
	"log"
	"os"
)

// ConfigSetup - настройка конфигурации для нашего приложения
func ConfigSetup() {
	// Database settings
	os.Setenv("DB_USERNAME", "nats_admin")
	os.Setenv("DB_PASSWORD", "qwerty")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "wb_l0_ns")
	os.Setenv("DB_PORT", "5432")

	os.Setenv("NATS_CLUSTER_ID", "test-cluster")
	os.Setenv("NATS_CLIENT_ID", "testUser")

	log.Println("Setting up config")
}

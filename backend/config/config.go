package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	Google     GoogleConfig
	WhatsApp   WhatsAppConfig
	Resend     ResendConfig
	Kafka      KafkaConfig
	Upload     UploadConfig
	Encryption EncryptionConfig
}

type EncryptionConfig struct {
	Key string // 32-byte hex-encoded encryption key
}

// Note: CSRF protection uses Go 1.25's http.CrossOriginProtection
// which is header-based (Sec-Fetch-Site, Origin) and doesn't require keys

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type GoogleConfig struct {
	ClientID        string
	ClientSecret    string
	RedirectURL     string
	StaffEmailDomain string
}

type WhatsAppConfig struct {
	APIURL   string
	APIToken string
}

type ResendConfig struct {
	APIKey string
	From   string
}

type KafkaConfig struct {
	Brokers      string
	TopicPayment string
}

type UploadConfig struct {
	Dir           string
	MaxSizeMB     int
}

func Load() (*Config, error) {
	jwtExpHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "168"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	maxUploadSize, err := strconv.Atoi(getEnv("MAX_UPLOAD_SIZE_MB", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_UPLOAD_SIZE_MB: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnvRequired("DATABASE_HOST"),
			Port:     getEnvRequired("DATABASE_PORT"),
			User:     getEnvRequired("DATABASE_USER"),
			Password: getEnvRequired("DATABASE_PASSWORD"),
			Name:     getEnvRequired("DATABASE_NAME"),
			SSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnvRequired("JWT_SECRET"),
			ExpirationHours: jwtExpHours,
		},
		Google: GoogleConfig{
			ClientID:        getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret:    getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:     getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			StaffEmailDomain: getEnv("STAFF_EMAIL_DOMAIN", "tazkia.ac.id"),
		},
		WhatsApp: WhatsAppConfig{
			APIURL:   getEnv("WHATSAPP_API_URL", ""),
			APIToken: getEnv("WHATSAPP_API_TOKEN", ""),
		},
		Resend: ResendConfig{
			APIKey: getEnv("RESEND_API_KEY", ""),
			From:   getEnv("RESEND_FROM", ""),
		},
		Kafka: KafkaConfig{
			Brokers:      getEnv("KAFKA_BROKERS", ""),
			TopicPayment: getEnv("KAFKA_TOPIC_PAYMENT", ""),
		},
		Upload: UploadConfig{
			Dir:       getEnv("UPLOAD_DIR", "./uploads"),
			MaxSizeMB: maxUploadSize,
		},
		Encryption: EncryptionConfig{
			Key: getEnvRequired("ENCRYPTION_KEY"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

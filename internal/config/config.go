package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	S3       S3Config
	SMTP     SMTPConfig
	CORS     CORSConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	SSLMode         string
}

type RedisConfig struct {
	Host string
	Port string
	DB   int
}

type JWTConfig struct {
	Secret                     string
	AccessTokenExpireMinutes   int
	RefreshTokenExpireDays     int
}

type S3Config struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	Bucket     string
	Region     string
	UseSSL     bool
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type CORSConfig struct {
	Origins []string
}

type LoggerConfig struct {
	Level    string
	Encoding string
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()

	// Set defaults
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("GIN_MODE", "debug")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("JWT_ACCESS_TOKEN_EXPIRE_MINUTES", 15)
	v.SetDefault("JWT_REFRESH_TOKEN_EXPIRE_DAYS", 7)
	v.SetDefault("S3_REGION", "us-east-1")
	v.SetDefault("S3_USE_SSL", false)
	v.SetDefault("SMTP_PORT", 587)
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_ENCODING", "json")

	if err := v.ReadInConfig(); err != nil && configPath != "" {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: v.GetString("SERVER_PORT"),
			Mode: v.GetString("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetString("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			Name:            v.GetString("DB_NAME"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			SSLMode:         v.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host: v.GetString("REDIS_HOST"),
			Port: v.GetString("REDIS_PORT"),
			DB:   v.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:                   v.GetString("JWT_SECRET"),
			AccessTokenExpireMinutes: v.GetInt("JWT_ACCESS_TOKEN_EXPIRE_MINUTES"),
			RefreshTokenExpireDays:   v.GetInt("JWT_REFRESH_TOKEN_EXPIRE_DAYS"),
		},
		S3: S3Config{
			Endpoint:  v.GetString("S3_ENDPOINT"),
			AccessKey: v.GetString("S3_ACCESS_KEY"),
			SecretKey: v.GetString("S3_SECRET_KEY"),
			Bucket:    v.GetString("S3_BUCKET"),
			Region:    v.GetString("S3_REGION"),
			UseSSL:    v.GetBool("S3_USE_SSL"),
		},
		SMTP: SMTPConfig{
			Host:     v.GetString("SMTP_HOST"),
			Port:     v.GetInt("SMTP_PORT"),
			User:     v.GetString("SMTP_USER"),
			Password: v.GetString("SMTP_PASSWORD"),
			From:     v.GetString("SMTP_FROM"),
		},
		CORS: CORSConfig{
			Origins: v.GetStringSlice("CORS_ORIGINS"),
		},
		Logger: LoggerConfig{
			Level:    v.GetString("LOG_LEVEL"),
			Encoding: v.GetString("LOG_ENCODING"),
		},
	}

	return cfg, nil
}

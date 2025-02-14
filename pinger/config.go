package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	BackendURL   string
	EnvName      string
	PingInterval time.Duration
	PingTimeout  time.Duration
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"BackendURL: %s\n"+
			"EnvName: %s\n"+
			"PingInterval: %s\n"+
			"PingTimeout: %s\n", c.BackendURL, c.EnvName, c.PingInterval, c.PingTimeout)
}

func loadConfig() Config {
	url := os.Getenv("BACKEND_API_URL")
	if url == "" {
		panic("BACKEND_API_URL should be defined in env configuration")
	}

	env := os.Getenv("ENV_NAME")

	pingInterval := 10 * time.Second
	if v := os.Getenv("PING_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			pingInterval = d
		} else {
			log.Printf("Ошибка парсинга PING_INTERVAL, используем значение по умолчанию: %v", pingInterval)
		}
	}

	pingTimeout := 2 * time.Second
	if v := os.Getenv("PING_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			pingTimeout = d
		} else {
			log.Printf("Ошибка парсинга PING_TIMEOUT, используем значение по умолчанию: %v", pingTimeout)
		}
	}

	return Config{
		BackendURL:   url,
		EnvName:      env,
		PingInterval: pingInterval,
		PingTimeout:  pingTimeout,
	}
}

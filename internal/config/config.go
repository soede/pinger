package config

import (
	"docker/pkg/utils"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const configPath = "config.yml"

type Config struct {
	App `yaml:"app"`
	DB  `yaml:"db"`
}

type App struct {
	EnvName      string `yaml:"env_name"`
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

type DB struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Name        string `yaml:"name"`
	SSL         string `yaml:"ssl"`
	SSLRootCert string `yaml:"sslrootcert"`
	SSLCert     string `yaml:"sslcert"`
	SSLKey      string `yaml:"sslkey"`
}

func LoadConfig() *Config {

	if utils.FileExists(configPath) {
		log.Println("Loading config from", configPath)
		return loadFromFile(configPath)
	} else {

	}
	log.Println("Config file not found. Loading from environment variables.")
	return mustEnvConfig()
}

func loadFromFile(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Sprintf("Error parsing YAML: %v", err))
	}

	return &cfg
}

func mustEnvConfig() *Config {
	env := os.Getenv("env")
	if env == "" {
		env = "local"
	}

	port := os.Getenv("port")
	if port == "" {
		panic("port should be defined in env configuration")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		panic("DB_HOST should be defined in env configuration")
	}

	DBport := os.Getenv("DB_PORT")
	if DBport == "" {
		panic("DB_PORT should be defined in env configuration")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		panic("DB_USER should be defined in env configuration")
	}

	password := os.Getenv("DB_PASSWORD")

	name := os.Getenv("DB_NAME")
	if name == "" {
		panic("DB_NAME should be defined in env configuration")
	}

	ssl := os.Getenv("SSL")
	if ssl == "" {
		ssl = "disable"
	}
	SSLRootCert := os.Getenv("SSLRootCert")
	SSLCert := os.Getenv("SSLCert")
	SSLKey := os.Getenv("SSLKey")
	return &Config{
		App: App{
			env,
			port,
			5,
			5,
		},
		DB: DB{
			host,
			DBport,
			user,
			password,
			name,
			ssl,
			SSLRootCert,
			SSLCert,
			SSLKey,
		},
	}
}

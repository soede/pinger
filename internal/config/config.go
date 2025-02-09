package config

import (
	"os"
)

type Config struct {
	App
	DB
}

type App struct {
	EnvName string
	Port    string
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func MustEnvConfig() Config {
	env := os.Getenv("env")
	if env == "" {
		env = "local"
		/*return Config{
			App{
				env,
				"8080",
			},
			DB{
				"localhost",
				"5432",
				"postgres",
				"",
				"postgres",
			},
		}*/
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

	return Config{
		App{
			env,
			port,
		},
		DB{
			host,
			DBport,
			user,
			password,
			name,
		},
	}
}

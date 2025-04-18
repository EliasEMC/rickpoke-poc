package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Port     string
	RickURL  string
	PokeURL  string
	Timeout  time.Duration
	LogLevel string
}

func Load() Config {
	timeout := 3 * time.Second
	if t := os.Getenv("HTTP_TIMEOUT_MS"); t != "" {
		if ms, err := time.ParseDuration(t + "ms"); err == nil {
			timeout = ms
		}
	}
	return Config{
		Port:     getenv("PORT", "8080"),
		RickURL:  getenv("RICK_URL", "https://rickandmortyapi.com/api"),
		PokeURL:  getenv("POKE_URL", "https://pokeapi.co/api/v2"),
		Timeout:  timeout,
		LogLevel: getenv("LOG_LEVEL", "info"),
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

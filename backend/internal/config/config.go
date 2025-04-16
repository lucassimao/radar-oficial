package config

import "os"

func Env() string {
	if e := os.Getenv("ENV"); e != "" {
		return e
	}
	return "development"
}

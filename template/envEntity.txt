package entity

import "os"

type Env struct {
	DB_USER          string
	DB_PASS          string
	DB_PORT          string
	DB_HOST          string
	DB_NAME          string
	JWT_SECRET_TOKEN string
}

func SetEnv() Env {
	env := Env{
		DB_USER:          os.Getenv("DB_USER"),
		DB_PASS:          os.Getenv("DB_PASS"),
		DB_PORT:          os.Getenv("DB_PORT"),
		DB_HOST:          os.Getenv("DB_HOST"),
		DB_NAME:          os.Getenv("DB_NAME"),
		JWT_SECRET_TOKEN: os.Getenv("JWT_SECRET_TOKEN"),
	}
	return env
}

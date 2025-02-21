package models

type Config struct {
	DBHost          string
	DBUser          string
	DBPass          string
	DBName          string
	Port            string
	DBPort          int
	MetricsInterval int
}

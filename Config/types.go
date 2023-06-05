package config

type LogLevelType string
type RedisURLType string
type NSType struct {
	Name   string
	Length uint
}
type SettingsType struct {
	NameSpaces   []NSType     `env:"NAMESPACE,notEmpty"`
	RedisURL     RedisURLType `env:"REDIS_URL,required"`
	Listen       string       `env:"LISTEN" envDefault:":8080"`
	ReserveCount uint         `env:"RESERVE_COUNT" envDefault:"100"`
	ReserveMin   uint         `env:"RESERVE_MIN" envDefault:"10"`
	LogLevel     LogLevelType `env:"LOG_LEVEL" envDefault:"warning"`
	FixSeed      bool         `env:"FIX_SEED" envDefault:"false"`
}

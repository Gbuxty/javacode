package config

import "github.com/spf13/viper"

const (
	pathConfigFile = ".env"
	ConfigType     = "dotenv"
)

type Config struct {
	HTTP     HTTPServer
	Postgres Postgres
}

func New() (*Config, error) {
	viper.SetConfigFile(pathConfigFile)
	viper.SetConfigType(ConfigType)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		HTTP:     newHTTPServer(),
		Postgres: newPostgres(),
	}, nil
}

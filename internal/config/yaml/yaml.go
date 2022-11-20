package yaml

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	pgModel "github.com/iTukaev/news_service/internal/repo/models"
)

type Config struct{}

func New() (*Config, error) {
	log.Println("Init config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "config init")
	}
	return &Config{}, nil
}

func (*Config) ServiceURL() string {
	return viper.GetString("service")
}

func (*Config) ClientHTTP() string {
	return viper.GetString("client")
}

func (*Config) LogLevel() string {
	return viper.GetString("log")
}

func (*Config) PGConfig() pgModel.Config {
	var pg pgModel.Config
	if err := viper.UnmarshalKey("pg", &pg); err != nil {
		log.Fatalf("Postgres config unmarshal error: %v\n", err)
	}
	return pg
}

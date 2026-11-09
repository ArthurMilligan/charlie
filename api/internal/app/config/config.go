package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string         `yaml:"env" env-default:"local"`
	APIConfig        APIConfig      `yaml:"api" env-required:"true"`
	DB               PostgresConfig `yaml:"postgres" env-required:"true"`
	TelegramBotToken string         `yaml:"tg_bot_token" env-required:"true"`
	GPTToken         string         `yaml:"gpt_token" env-required:"true"`
}

type APIConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	Host           string `yaml:"host" env-required:"true"`
	Port           int    `yaml:"port" env-required:"true"`
	User           string `yaml:"user" env-required:"true"`
	Password       string `yaml:"password" env-required:"true"`
	DBName         string `yaml:"db_name" env-required:"true"`
	MigrationsPath string
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

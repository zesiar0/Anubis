package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	// Server
	Name       string `mapstructure:"NAME"`
	BindHost   string `mapstructure:"BIND_HOST"`
	SSHPort    string `mapstructure:"SSH_PORT"`
	SSHTimeout int    `mapstructure:"SSH_TIMEOUT"`

	// Logger
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// SSH
	SSHVersion  string `mapstructure:"SSH_VERSION"`
	HostKeyFile string `mapstructure:"HOST_KEY_FILE"`

	// Database
	DBType    string `mapstructure:"DB_TYPE"`
	DBAddress string `mapstructure:"DB_ADDRESS"`
	DBUser    string `mapstructure:"DB_USER"`
	DBPass    string `mapstructure:"DB_PASS"`
	DBPort    string `mapstructure:"DB_PORT"`
	DBTable   string `mapstructure:"DB_TABLE"`
	// Mysql params
	DBParams map[string]string `mapstructure:"DB_PARAMS"`
}

var GlobalConfig *Config

func Setup(configPath string) {
	var conf Config
	loadConfigFromFile(configPath, &conf)
	GlobalConfig = &conf
	log.Printf("Set up global configuration\n")
}

func loadConfigFromFile(path string, config *Config) {
	var err error
	if valid(path) {
		fileViper := viper.New()
		fileViper.SetConfigFile(path)
		if err = fileViper.ReadInConfig(); err == nil {
			if err = fileViper.Unmarshal(config); err == nil {
				log.Printf("Load config from %s success\n", path)
			}
		}
	}
	if err != nil {
		log.Fatalf("Load config from %s failed: %s\n", path, err)
	}
}

func valid(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

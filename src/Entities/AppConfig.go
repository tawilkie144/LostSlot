package Entities

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	UserServiceDataSource string `mapstructure:"user_service_data_source"`
	PostgresUser          string `mapstructure:"postgres_user"`
	PostgresPassword      string `mapstructure:"postgres_password"`
	PostgresDB            string `mapstructure:"postgres_db"`
	PostgresHost          string `mapstructure:"postgres_host"`
	PostgresPort          int    `mapstructure:"postgres_port"`
}

type appConfig struct {
	config      *Config
	environment string
	initialized bool
}

var aConfig *appConfig = &appConfig{}

func init() {
	pflag.StringP("env", "e", "error", "Environment")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("Error binding to command line: %s", err)
	}
	env := viper.GetString("env")
	if env == "error" {
		log.Fatalf("error binding to command line")
	}
	log.Println("Loading config for environment: ", env)
	viper.AddConfigPath("./environment")
	viper.SetConfigName(env)
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&aConfig.config); err != nil {
		log.Fatalf("unable to decode into struct: %w", err)
	}

	if aConfig.config.PostgresUser == "" || aConfig.config.PostgresPassword == "" {
		log.Fatalf("unable to decode into struct")
	}
	aConfig.initialized = true
}

func GetConfig() (Config, error) {
	if aConfig.initialized {
		return *aConfig.config, nil
	}
	log.Fatalf("For right now, just crashing. Later I'll figure out what to do in this case")

	return *aConfig.config, nil
}

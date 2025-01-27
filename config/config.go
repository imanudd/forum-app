package config

import (
	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	Service Service `mapstructure:"service"`
	Db      Mysql   `mapstructure:"db"`
}

type Service struct {
	Name         string `mapstructure:"name"`
	Port         string `mapstructure:"port"`
	Environment  string `mapstructure:"environment"`
	MigrationDir string `mapstructure:"migrationDir"`
	SignatureKey string `mapstructure:"signatureKey"`
}

type Mysql struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
}

type option struct {
	configFolders []string
	configFile    string
	configType    string
}

type Option func(*option)

func getDefaultConfigOption() *option {
	return &option{
		configFolders: []string{"./config"},
		configFile:    "config",
		configType:    "yaml",
	}
}

func InitConfig(options ...Option) error {
	opts := getDefaultConfigOption()

	for _, overrideOpt := range options {
		overrideOpt(opts)
	}

	for _, configFolder := range opts.configFolders {
		viper.AddConfigPath(configFolder)
	}

	viper.SetConfigName(opts.configFile)
	viper.SetConfigType(opts.configType)
	viper.AutomaticEnv()

	config = new(Config)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(config)
}

func WithConfigFolder(configFolder []string) Option {
	return func(o *option) {
		o.configFolders = configFolder
	}
}

func WithConfigFile(configFile string) Option {
	return func(o *option) {
		o.configFile = configFile
	}
}

func WithConfigType(configType string) Option {
	return func(o *option) {
		o.configType = configType
	}
}

func Get() *Config {
	if config == nil {
		return &Config{}
	}

	return config
}

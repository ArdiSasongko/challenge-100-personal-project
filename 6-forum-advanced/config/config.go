package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var config *Config

type option struct {
	configFolder []string
	configFile   string
	configTypes  string
}

type Option func(*option)

func Init(opts ...Option) error {
	opt := &option{
		configFolder: getDefaultFolder(),
		configFile:   getDefaultFile(),
		configTypes:  getDefaultType(),
	}

	for _, o := range opts {
		o(opt)
	}

	for _, folder := range opt.configFolder {
		viper.AddConfigPath(folder)
	}

	viper.SetConfigName(opt.configFile)
	viper.SetConfigType(opt.configTypes)
	viper.AutomaticEnv()

	config = new(Config)

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithField("viper config", err.Error()).Error(err.Error())
		return err
	}

	return viper.Unmarshal(config)
}

func getDefaultFolder() []string {
	return []string{"./config"}
}

func getDefaultFile() string {
	return "config"
}

func getDefaultType() string {
	return "yaml"
}

func ConfigFolder(folder []string) Option {
	return func(o *option) {
		o.configFolder = folder
	}
}

func ConfigFile(file string) Option {
	return func(o *option) {
		o.configFile = file
	}
}

func ConfigType(types string) Option {
	return func(o *option) {
		o.configTypes = types
	}
}

func GetConfig() *Config {
	if config == nil {
		return &Config{}
	}

	return config
}

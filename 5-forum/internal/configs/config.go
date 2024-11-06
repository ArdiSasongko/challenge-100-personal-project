package configs

import "github.com/spf13/viper"

var config *Config

type option struct {
	configFolder []string
	configFile   string
	configType   string
}

type Option func(*option)

func Init(opts ...Option) error {
	opt := &option{
		configFolder: getDefaultFolder(),
		configFile:   getDefaultFile(),
		configType:   getDefaultType(),
	}

	for _, o := range opts {
		o(opt)
	}

	for _, folder := range opt.configFolder {
		viper.AddConfigPath(folder)
	}

	viper.SetConfigName(opt.configFile)
	viper.SetConfigType(opt.configType)
	viper.AutomaticEnv()

	config = new(Config)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&config)
}

func getDefaultFolder() []string {
	return []string{"./configs"}
}

func getDefaultFile() string {
	return "config"
}

func getDefaultType() string {
	return "yaml"
}

func ConfigWithFolder(configFolder []string) Option {
	return func(o *option) {
		o.configFolder = configFolder
	}
}

func ConfigWithFile(configFile string) Option {
	return func(o *option) {
		o.configFile = configFile
	}
}

func ConfigWithType(configType string) Option {
	return func(o *option) {
		o.configType = configType
	}
}

func GetConfig() *Config {
	if config == nil {
		return &Config{}
	}
	return config
}

package config

type (
	Config struct {
		Service  Service  `mapstructure:"service"`
		Database Database `mapstructure:"database"`
	}

	Service struct {
		Port          string `mapstructure:"port"`
		SecretJWT     string `mapstructure:"secret_jwt"`
		CloudInaryURL string `mapstructure:"cloud_inary_url"`
	}

	Database struct {
		SourceName string `mapstructure:"source_name"`
	}
)

package configs

type (
	Config struct {
		Service  Service  `mapstructure:"service"`
		Database Database `mapstructre:"database"`
	}

	Service struct {
		Port          string `mapstructure:"port"`
		CloudInaryURL string `mapstructure:"cloud_inary_url"`
	}

	Database struct {
		DataSourceName string `mapstructure:"data_source_name"`
	}
)

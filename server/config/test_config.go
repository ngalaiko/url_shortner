package config

// NewTestConfig returns test config
func NewTestConfig() *Config {
	return &Config{
		Db: DbConfig{
			Driver:       "postgres",
			Connect:      "host=localhost user=url_short_test dbname=url_short_test sslmode=disable password=secret",
			MaxIdleConns: 5,
			MaxOpenConns: 5,
		},
	}
}

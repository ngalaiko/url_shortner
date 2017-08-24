package config

type DbConfig struct {
	Driver       string `yaml:"Driver"`
	Connect      string `yaml:"Connect"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
}

package config

type FacebookConfig struct {
	FacebookAppID     string `yaml:"FacebookAppID"`
	FacebookAppSecret string `yaml:"FacebookAppSecret"`
	FacebookApiSDK    string `yaml:"FacebookApiSDK"`
	FacebookLoginURL  string `yaml:"FacebookLoginURL"`
}

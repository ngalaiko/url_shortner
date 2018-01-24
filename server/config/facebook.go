package config

// FacebookConfig holds config for facebook
type FacebookConfig struct {
	FacebookAppID     string `yaml:"FacebookAppID"`
	FacebookAppSecret string `yaml:"FacebookAppSecret"`
	FacebookAPISDK    string `yaml:"FacebookApiSDK"`
	FacebookLoginURL  string `yaml:"FacebookLoginURL"`
}

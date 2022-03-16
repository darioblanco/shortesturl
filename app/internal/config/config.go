package config

import (
	"time"

	"github.com/spf13/viper"
)

// The AppName that identifies the application
const AppName = "shortesturl"

// A Values struct that holds all the loaded configuration variables for the app
type Values struct {
	Environment          string
	HttpHost             string
	HttpPort             int
	HttpScheme           string
	IsDevelopment        bool
	RedisHost            string
	RedisPort            string
	UrlLength            int
	UrlExpirationInHours time.Duration
	Version              string
}

// New loads config variables from file paths
func New(configName string, configPaths ...string) (*Values, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.SetEnvPrefix(AppName)
	v.AutomaticEnv()
	// Bind multi-word environment variables
	v.BindEnv("environment", "SHORTESTURL_ENVIRONMENT")
	v.BindEnv("httpHost", "SHORTESTURL_HTTP_HOST")
	v.BindEnv("httpPort", "SHORTESTURL_HTTP_PORT")
	v.BindEnv("httpScheme", "SHORTESTURL_HTTP_SCHEME")
	v.BindEnv("redisHost", "SHORTESTURL_REDIS_HOST")
	v.BindEnv("redisPort", "SHORTESTURL_REDIS_PORT")
	v.BindEnv("urlLength", "SHORTESTURL_URL_LENGTH")
	v.BindEnv("urlLength", "SHORTESTURL_URL_EXPIRATION_IN_HOURS")
	v.BindEnv("version", "SHORTESTURL_VERSION")
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	} else {
		var conf Values
		err := v.Unmarshal(&conf)
		// Dynamic load of configuration variables
		conf.IsDevelopment = conf.Environment == "dev"
		return &conf, err
	}
}

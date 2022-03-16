package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	conf, err := New("config.default", "../../../cmd")
	assert.NoError(t, err)

	assert.Equal(t, Values{
		Environment:          "dev",
		HttpHost:             "localhost",
		HttpPort:             3000,
		HttpScheme:           "http",
		IsDevelopment:        true,
		RedisHost:            "localhost",
		RedisPort:            "6379",
		UrlLength:            6,
		UrlExpirationInHours: 0,
		Version:              "unknown",
	}, *conf)
}

func TestNew_ReadError(t *testing.T) {
	conf, err := New("config.default")
	assert.Error(t, err)
	assert.Nil(t, conf)
}

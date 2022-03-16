package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/darioblanco/shortesturl/app/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	ctx   context.Context
	cache Cache
	mr    *miniredis.Miniredis
}

func (s *TestSuite) SetupSuite() {
	s.ctx = context.Background()

	mr, _ := miniredis.Run()
	s.mr = mr
	conf := &config.Values{
		RedisHost: mr.Host(),
		RedisPort: mr.Port(),
	}

	c, err := New(
		context.Background(),
		conf,
		logging.NewTest(s.T()),
	)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(),
		fmt.Sprintf("%s:%s", mr.Host(), mr.Port()),
		c.(*cache).client.Options().Addr,
	)
	s.cache = c
}

func TestCache(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func TestNew(t *testing.T) {
	mr, _ := miniredis.Run()
	c, err := New(
		context.Background(),
		&config.Values{
			RedisHost: mr.Host(),
			RedisPort: mr.Port(),
		},
		logging.NewTest(t),
	)
	assert.NoError(t, err)
	assert.Equal(t,
		fmt.Sprintf("%s:%s", mr.Host(), mr.Port()),
		c.(*cache).client.Options().Addr,
	)
}

func TestNew_Development(t *testing.T) {
	_, err := New(
		context.Background(),
		&config.Values{IsDevelopment: true},
		logging.NewTest(t),
	)
	assert.NoError(t, err)
}

func (s *TestSuite) TestNew_PingError() {
	r, err := New(
		s.ctx,
		&config.Values{
			RedisHost: "fake",
			RedisPort: "678912",
		},
		logging.NewTest(s.T()),
	)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), r)
}

func (s *TestSuite) TestGet() {
	s.mr.Set("key1", "val1")
	val, err := s.cache.Get(s.ctx, "key1")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "val1", val)
}

func (s *TestSuite) TestGet_KeyNotFound() {
	val, err := s.cache.Get(s.ctx, "key2")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "", val)
}

func (s *TestSuite) TestSetIfNotExists() {
	res, err := s.cache.SetIfNotExists(s.ctx, "key3", "value3", 0)
	assert.NoError(s.T(), err)
	assert.True(s.T(), res)
	val, _ := s.mr.Get("key3")
	assert.Equal(s.T(), "value3", val)
}

func (s *TestSuite) TestSetIfNotExists_AlreadyExistsWithSameValue() {
	s.mr.Set("key4", "val4")
	res, err := s.cache.SetIfNotExists(s.ctx, "key4", "val4", 0)
	assert.NoError(s.T(), err)
	assert.True(s.T(), res)
}

func (s *TestSuite) TestSetIfNotExists_AlreadyExistsWithDifferentValue() {
	s.mr.Set("key5", "valDifferent")
	res, err := s.cache.SetIfNotExists(s.ctx, "key5", "val5", 0)
	assert.NoError(s.T(), err)
	assert.False(s.T(), res)
}

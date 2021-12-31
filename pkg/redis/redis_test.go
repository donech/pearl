package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert := require.New(t)
	c := Config{
		Name:     Default,
		Addr:     ":6379",
		Password: "",
		DB:       0,
	}
	New(c)
	redis := Redis()
	s := redis.Set("solar", "solar", time.Second*30)
	assert.NoError(s.Err())
	r, err := redis.Get("solar").Result()
	assert.NoError(err)
	assert.Equal("solar", r)
}

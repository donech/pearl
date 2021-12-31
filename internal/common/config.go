package common

import (
	"context"

	"github.com/donech/pearl/pkg/ent/db"
	"github.com/donech/pearl/pkg/entry/grpc"
	"github.com/donech/pearl/pkg/jwt"
	"github.com/donech/pearl/pkg/log"
	"github.com/donech/pearl/pkg/redis"
	"github.com/spf13/viper"
)

var cfg *Config = &Config{}

type Config struct {
	SystemInfo SystemInfo     `yaml:"systemInfo"`
	Grpc       grpc.Config    `yaml:"grpc"`
	Redis      redis.Config   `yaml:"redis"`
	DB         []*db.DBConfig `yaml:"db"`
	Jwt        jwt.Config     `yaml:"jwt"`
	Log        log.Config     `yaml:"log"`
}

type SystemInfo struct {
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	Desc      string `yaml:"desc"`
	Href      string `yaml:"href"`
	Copyright string `yaml:"copyright"`
}

func InitConfig() *Config {
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.S(context.Background()).Fatalf("init config error")
	}
	return cfg
}

func C() *Config {
	return cfg
}

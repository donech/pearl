package common

import (
	"context"

	"github.com/donech/pearl/internal/data"
	"github.com/donech/pearl/pkg/ent/db"
	"github.com/donech/pearl/pkg/log"
	"github.com/donech/pearl/pkg/redis"
)

//InitGlobal init global vat and return a clean up func.
func InitGlobal() {
	initLogger()
	initRedis()
	initDB()
}

func initLogger() {
	_, err := log.New(cfg.Log)
	if err != nil {
		log.S(context.Background()).Panicf("can't init zap logger :", err)
	}
	log.S(context.Background()).Info("init logger done")
}

func initRedis() {
	redis.New(cfg.Redis)
}

func initDB() {
	if cfg.DB != nil {
		for _, dbC := range cfg.DB {
			db.New(dbC)
		}
	}
	data.NewDB()
}

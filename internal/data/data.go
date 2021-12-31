package data

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/donech/pearl/internal/data/ent"
	"github.com/donech/pearl/pkg/ent/db"
	"github.com/donech/pearl/pkg/log"
)

var DB *ent.Client

func NewDB() *ent.Client {
	logger := func(ctx context.Context, args ...interface{}) {
		log.S(ctx).Infof("ent: args=%+v", args...)
	}
	DB = ent.NewClient(ent.Driver(dialect.DebugWithContext(db.EntDriver(), logger)))
	return DB
}

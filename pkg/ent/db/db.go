package db

import (
	"context"
	"database/sql"
	"sync"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/donech/pearl/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

const Default = "default"

// DBDriver db driver
type DBDriver string

const (
	MySQL    DBDriver = "mysql"
	Postgres DBDriver = "postgres"
	SQLite   DBDriver = "sqlite3"
)

// DBConfig keeps the settings to setup db connection.
type DBConfig struct {
	Name   string `json:"name" yaml:"name"`
	Driver string `json:"driver" yaml:"driver"`
	// DSN data source name
	//
	// [MySQL] username:password@tcp(localhost:3306)/dbname?timeout=10s&charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local
	//
	// [Postgres] host=localhost port=5432 user=root password=secret dbname=test connect_timeout=10 sslmode=disable
	//
	// [SQLite] file::memory:?cache=shared
	DSN string `json:"dsn" yaml:"dsn"`

	// Options optional settings to setup db connection.
	Options *DBOptions `json:"options" yaml:"options"`
}

// DBOptions optional settings to setup db connection.
type DBOptions struct {
	// MaxOpenConns is the maximum number of open connections to the database.
	// Use value -1 for no timeout and 0 for default.
	// Default is 20.
	MaxOpenConns int `json:"max_open_conns" yaml:"maxOpenConns"`

	// MaxIdleConns is the maximum number of connections in the idle connection pool.
	// Use value -1 for no timeout and 0 for default.
	// Default is 10.
	MaxIdleConns int `json:"max_idle_conns" yaml:"maxIdleConns"`

	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	// Use value -1 for no timeout and 0 for default.
	// Default is 10 minutes.
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"connMaxLifetime"`

	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	// Use value -1 for no timeout and 0 for default.
	// Default is 5 minutes.
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"connMaxIdleTime"`
}

func (o *DBOptions) rebuild(opt *DBOptions) {
	if opt.MaxOpenConns > 0 {
		o.MaxOpenConns = opt.MaxOpenConns
	} else {
		if opt.MaxOpenConns == -1 {
			o.MaxOpenConns = 0
		}
	}

	if opt.MaxIdleConns > 0 {
		o.MaxIdleConns = opt.MaxIdleConns
	} else {
		if opt.MaxIdleConns == -1 {
			o.MaxIdleConns = 0
		}
	}

	if opt.ConnMaxLifetime > 0 {
		o.ConnMaxLifetime = opt.ConnMaxLifetime
	} else {
		if opt.ConnMaxLifetime == -1 {
			o.ConnMaxLifetime = 0
		}
	}

	if opt.ConnMaxIdleTime > 0 {
		o.ConnMaxIdleTime = opt.ConnMaxIdleTime
	} else {
		if opt.ConnMaxIdleTime == -1 {
			o.ConnMaxIdleTime = 0
		}
	}
}

var (
	defaultDB *sql.DB
	dbmap     sync.Map

	defaultEntDriver *entsql.Driver
	entmap           sync.Map
)

func New(cfg *DBConfig) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)

	if err != nil {
		log.S(context.Background()).Panic("db init error", zap.String("name", cfg.Name), zap.Error(err))
	}

	if err = db.Ping(); err != nil {
		db.Close()

		log.S(context.Background()).Panic("db init error", zap.String("name", cfg.Name), zap.Error(err))
	}

	opt := &DBOptions{
		MaxOpenConns:    20,
		MaxIdleConns:    10,
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	if cfg.Options != nil {
		opt.rebuild(cfg.Options)
	}

	db.SetMaxOpenConns(opt.MaxOpenConns)
	db.SetMaxIdleConns(opt.MaxIdleConns)
	db.SetConnMaxLifetime(opt.ConnMaxLifetime)
	db.SetConnMaxIdleTime(opt.ConnMaxIdleTime)

	entDriver := entsql.OpenDB(cfg.Driver, db)

	if cfg.Name == Default {
		defaultDB = db
		defaultEntDriver = entDriver
	}

	dbmap.Store(cfg.Name, db)
	entmap.Store(cfg.Name, entDriver)
	log.S(context.Background()).Infof("init db.%s success", cfg.Name)
}

// EntDriver returns an ent dialect.Driver.
func EntDriver(name ...string) *entsql.Driver {
	if len(name) == 0 || name[0] == Default {
		if defaultEntDriver == nil {
			log.S(context.Background()).Panicf("unknown db.%s (forgotten configure?)", Default)
		}

		return defaultEntDriver
	}

	v, ok := entmap.Load(name[0])

	if !ok {
		log.S(context.Background()).Panicf("unknown db.%s (forgotten configure?)", name[0])
	}

	return v.(*entsql.Driver)
}

// DB returns a db.
func DB(name ...string) *sql.DB {
	if len(name) == 0 || name[0] == Default {
		if defaultDB == nil {
			log.S(context.Background()).Panicf("unknown db.%s (forgotten configure?)", Default)
		}

		return defaultDB
	}

	v, ok := dbmap.Load(name[0])

	if !ok {
		log.S(context.Background()).Panicf("unknown db.%s (forgotten configure?)", name[0])
	}

	return v.(*sql.DB)
}

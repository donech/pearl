package gin

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/donech/pearl/pkg/entry/gin/middleware"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func NewEntry(conf *Config, router Router, logger *zap.Logger) *Entry {
	engine := gin.New()
	engine.Use(middleware.GinZap(time.RFC3339, true, conf.Mod))
	engine.Use(middleware.RecoveryWithZap(true))
	return &Entry{
		conf:   conf,
		engine: engine,
		router: router,
		logger: logger,
	}
}

type Entry struct {
	conf   *Config
	engine *gin.Engine
	router Router
	logger *zap.Logger
	srv    *http.Server
}

func (e *Entry) Engine() *gin.Engine {
	return e.engine
}

func (e *Entry) Run() error {
	srv := &http.Server{
		Addr:    e.conf.Addr,
		Handler: e.engine,
	}
	e.router.Init(e.engine)
	go func() {
		log.Println("start http server at ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	e.srv = srv
	return nil
}

func (e *Entry) Stop(ctx context.Context) error {
	err := e.srv.Shutdown(ctx)
	if err != nil {
		log.Println("Http Server Shutdown failed: ", err)
		return err
	}
	log.Println("Http Server exiting")
	return nil
}

type Controller interface {
	RegisterRoute(engine *gin.RouterGroup)
}

type Router interface {
	Init(engine *gin.Engine)
	RegisterController(c Controller)
}

type DefaultRouter struct {
	controllers []Controller
}

func (r *DefaultRouter) RegisterController(c Controller) {
	if r.controllers == nil {
		r.controllers = make([]Controller, 0, 2)
	}
	r.controllers = append(r.controllers, c)
}

func (r DefaultRouter) Init(engine *gin.Engine) {
	rootGroup := engine.Group("/")
	for _, c := range r.controllers {
		c.RegisterRoute(rootGroup)
	}
}

package grpc

import (
	"context"
	"net"

	"github.com/donech/pearl/pkg/jwt"
	"github.com/donech/pearl/pkg/log"

	"google.golang.org/grpc/reflection"

	"github.com/donech/pearl/pkg/entry/grpc/interceptor"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"
)

func New(config Config, options ...Option) *Entry {
	en := &Entry{
		config: config,
	}
	for _, option := range options {
		option(en)
	}
	return en
}

type Option func(entry *Entry)

func WithRegisteServer(server RegisteServer) Option {
	return func(entry *Entry) {
		entry.registeServer = server
	}
}

func WithRegisteWebHandler(handler RegisteWebHandler) Option {
	return func(entry *Entry) {
		entry.registeWebHandler = handler
	}
}

func WithJwtFactory(jwtFactory *jwt.JWTFactory) Option {
	return func(entry *Entry) {
		entry.jwtFactory = jwtFactory
	}
}

func WithJumpMethods(jumps map[string]bool) Option {
	return func(entry *Entry) {
		entry.jumpMethods = jumps
	}
}

type RegisteServer func(server *grpc.Server)
type RegisteWebHandler func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

type Entry struct {
	config            Config
	srv               *grpc.Server
	registeServer     RegisteServer
	registeWebHandler RegisteWebHandler
	jwtFactory        *jwt.JWTFactory
	jumpMethods       map[string]bool
}

func (e *Entry) Run() error {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcopentracing.UnaryServerInterceptor(grpcopentracing.WithTracer(opentracing.GlobalTracer())),
			interceptor.Log(),
			interceptor.Recover(),
			interceptor.Jwt(e.jwtFactory, e.jumpMethods),
		)))
	if e.registeServer != nil {
		e.registeServer(srv)
		e.srv = srv
	}
	if e.config.EnableReflect {
		log.S(context.Background()).Info("enable grpc reflect")
		reflection.Register(srv)
	}
	listen, err := net.Listen("tcp", e.config.Port)
	if err != nil {
		log.S(context.Background()).Errorf("listen tcp error: %s", err)
		return err
	}
	log.S(context.Background()).Infof("listening tcp port: %s", e.config.Port)

	go func() {
		log.S(context.Background()).Info("start grpc server at ", e.config.Port)
		if err = srv.Serve(listen); err != nil {
			log.S(context.Background()).Fatalf("grpc serve listen error: %s", err)
		}
	}()
	return nil
}

func (e *Entry) Stop(ctx context.Context) error {
	e.srv.GracefulStop()
	return nil
}

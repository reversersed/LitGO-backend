package app

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/reversersed/go-grpc/tree/main/api_gateway/docs"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/config"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers"
	freecache "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/cache"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/logging/logrus"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/middleware"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/shutdown"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API
// @version 1.0

// @host localhost:9000
// @BasePath /api/v1/

// @scheme http
// @accept json

// @securityDefinitions.apiKey ApiKeyAuth
// @in Cookie
// @name Authorization
func New() (*app, error) {
	logger, err := logrus.GetLogger()
	if err != nil {
		return nil, err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	server := &app{
		router: gin.New(),
		config: cfg,
		logger: logger,
		cache:  freecache.NewFreeCache(104857600),
	}
	server.logger.Info("setting up gin router...")
	gin.SetMode(server.config.Server.Environment)
	server.router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PATCH", "DELETE"},
	}))
	server.router.Use(gin.LoggerWithWriter(logger.Writer()))
	server.router.Use(middleware.ErrorHandler)
	server.logger.Info("router has been set up")

	return server, nil
}
func (a *app) Run() error {
	go shutdown.Graceful(a)

	a.logger.Info("setting up middleware...")
	jwt := middleware.NewJwtMiddleware(a.logger, a.cache, a.config.Server.JwtSecret)

	a.logger.Info("setting up handlers...")
	if userHandler, err := handlers.NewUserHandler(a.config.Url, a.logger, jwt); err != nil {
		return err
	} else {
		a.handlers = append(a.handlers, userHandler)
		userHandler.RegisterRouter(a.router)
	}
	if a.config.Server.Environment == "debug" {
		a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	a.logger.Infof("starting listening address %s:%d...", a.config.Server.Host, a.config.Server.Port)
	if err := a.router.Run(fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)); err != nil {
		return err
	}
	return nil
}
func (a *app) Close() error {
	for _, i := range a.handlers {
		if err := i.Close(); err != nil {
			return err
		}
	}
	return nil
}

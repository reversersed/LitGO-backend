package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/reversersed/LitGO-backend-pkg/logging/logrus"
	"github.com/reversersed/LitGO-backend-pkg/middleware"
	"github.com/reversersed/LitGO-backend-pkg/shutdown"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/config"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/endpoint"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/swagger"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
	}
	server.logger.Info("setting up gin router...")
	gin.SetMode(server.config.Server.Environment)
	server.router.Use(cors.New(cors.Config{
		AllowWildcard:    true,
		AllowAllOrigins:  false,
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost", "https://*.ngrok-free.app", ("http://localhost:" + strconv.Itoa(server.config.Server.Port)), "http://localhost:7000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Accept", "Cookie", "Access-Control-Expose-Headers"},
	}))
	server.router.Use(gin.LoggerWithWriter(logger.Writer()))
	server.router.Use(middleware.ErrorHandler)
	server.router.Use(gin.CustomRecoveryWithWriter(logger.Writer(), middleware.RecoveryMiddleware))
	server.logger.Info("router has been set up")

	return server, nil
}
func (a *app) Run() error {
	ctx := context.Background()
	a.logger.Info("setting up jwt middlware...")

	conn, err := grpc.NewClient(a.config.Url.UserServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	if jwtMiddleware, err := middleware.NewJwtMiddleware(a.logger, a.config.Server.JwtSecret, users_pb.NewUserClient(conn)); err != nil {
		return err
	} else {
		a.router.Use(jwtMiddleware.Middleware)
	}

	a.logger.Info("setting up grpc clients...")

	mux, err := endpoint.RegisterEndpoints(ctx, a.config.Url)
	if err != nil {
		return err
	}

	a.router.Any("/api/v1/*any", gin.WrapH(mux))
	if a.config.Server.Environment == "debug" {
		a.logger.Infof("creating swagger page at %s:%d/api/swagger/index.html", a.config.Server.Host, a.config.Server.Port)
		swagger.InitiateSwagger(a.router)
	}
	a.logger.Infof("starting listening address %s:%d...", a.config.Server.Host, a.config.Server.Port)

	go shutdown.Graceful(a, conn)

	if err := a.router.Run(fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)); err != nil {
		return err
	}
	ctx.Done()
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

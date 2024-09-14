package app

import (
	"fmt"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/reversersed/go-grpc/tree/main/api_gateway/docs"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/config"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers/author"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers/book"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers/genre"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers/user"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/logging/logrus"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/middleware"
	authors_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/authors"
	books_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/books"
	genres_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/genres"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/users"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/shutdown"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	}
	server.logger.Info("setting up gin router...")
	gin.SetMode(server.config.Server.Environment)
	server.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowCredentials: true,
		AllowOrigins:     []string{("http://localhost:" + strconv.Itoa(server.config.Server.Port)), "http://localhost:7000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Accept", "Cookie", "Access-Control-Expose-Headers"},
	}))
	server.router.Use(gin.LoggerWithWriter(logger.Writer()))
	server.router.Use(middleware.ErrorHandler)
	server.logger.Info("router has been set up")

	return server, nil
}
func (a *app) Run() error {
	a.logger.Info("setting up grpc clients...")
	userConnection, err := grpc.NewClient(a.config.Url.UserServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	userClient := users_pb.NewUserClient(userConnection)

	genreConnection, err := grpc.NewClient(a.config.Url.GenreServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	genreClient := genres_pb.NewGenreClient(genreConnection)

	authorConnection, err := grpc.NewClient(a.config.Url.AuthorServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	authorClient := authors_pb.NewAuthorClient(authorConnection)

	bookConnection, err := grpc.NewClient(a.config.Url.BookServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	bookClient := books_pb.NewBookClient(bookConnection)

	a.logger.Info("setting up middleware...")
	jwt := middleware.NewJwtMiddleware(a.logger, a.config.Server.JwtSecret)

	a.logger.Info("setting up handlers...")

	// users
	userHandler := user.New(userClient, a.logger, jwt)
	jwt.ApplyUserServer(userClient)
	a.handlers = append(a.handlers, userHandler)
	userHandler.RegisterRouter(a.router)

	// genres
	genreHandler := genre.New(genreClient, a.logger, jwt)
	a.handlers = append(a.handlers, genreHandler)
	genreHandler.RegisterRouter(a.router)

	// authors
	authorHandler := author.New(authorClient, a.logger, jwt)
	a.handlers = append(a.handlers, authorHandler)
	authorHandler.RegisterRouter(a.router)

	// books
	bookHandler := book.New(bookClient, a.logger, jwt)
	a.handlers = append(a.handlers, bookHandler)
	bookHandler.RegisterRouter(a.router)

	if a.config.Server.Environment == "debug" {
		a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	a.logger.Infof("starting listening address %s:%d...", a.config.Server.Host, a.config.Server.Port)
	go shutdown.Graceful(a, userConnection)
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

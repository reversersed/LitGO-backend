package app

import (
	"fmt"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/reversersed/LitGO-backend/tree/main/api_gateway/docs"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/config"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers/author"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers/book"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers/collection"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers/genre"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers/review"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers/user"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/pkg/logging/logrus"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/pkg/middleware"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/pkg/shutdown"
	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	collections_pb "github.com/reversersed/LitGO-proto/gen/go/collections"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
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
// @accept x-www-form-urlencoded

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

	reviewConnection, err := grpc.NewClient(a.config.Url.ReviewServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	reviewClient := reviews_pb.NewReviewClient(reviewConnection)

	collectionConnection, err := grpc.NewClient(a.config.Url.CollectionServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	collectionClient := collections_pb.NewCollectionClient(collectionConnection)

	a.logger.Info("setting up middleware...")
	jwt := middleware.NewJwtMiddleware(a.logger, a.config.Server.JwtSecret)

	a.logger.Info("setting up handlers...")

	// users
	userHandler := user.New(userClient, a.logger, jwt)
	jwt.ApplyUserServer(userClient)
	a.handlers = append(a.handlers, userHandler)

	// genres
	a.handlers = append(a.handlers, genre.New(genreClient, a.logger, jwt))

	// authors
	a.handlers = append(a.handlers, author.New(authorClient, a.logger, jwt))

	// books
	a.handlers = append(a.handlers, book.New(bookClient, a.logger, jwt))

	// reviews
	a.handlers = append(a.handlers, review.New(reviewClient, a.logger, jwt))

	// collections
	a.handlers = append(a.handlers, collection.New(collectionClient, a.logger, jwt))

	for _, hander := range a.handlers {
		hander.RegisterRouter(a.router)
	}

	if a.config.Server.Environment == "debug" {
		a.router.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	a.logger.Infof("starting listening address %s:%d...", a.config.Server.Host, a.config.Server.Port)

	go shutdown.Graceful(a, userConnection, genreConnection, authorConnection, bookConnection, reviewConnection, collectionConnection)

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

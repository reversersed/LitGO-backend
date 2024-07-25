package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/middleware"
	_ "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary      Authenticates user
// @Description  check if current user has legit token
// @Tags         users
// @Produce      json
// @Success      200  {object}  handlers.UserAuthenticate.UserResponse "User successfully authorized"
// @Failure      401  {object}  middleware.CustomError "User does not authorized"
// @Failure      410 {object}  middleware.CustomError "Service does not responding (maybe crush)"
// @Router       /users/auth [get]
func (h *userHandler) UserAuthenticate(c *gin.Context) {
	id, exist := c.Get(string(middleware.UserIdKey))
	if !exist {
		c.Error(status.Error(codes.Unauthenticated, "no user credentials found"))
		return
	}
	reply, err := h.client.GetUserById(c.Request.Context(), &users_pb.UserIdRequest{Id: id.(string)})
	if err != nil {
		c.Error(err)
		return
	}
	type UserResponse struct {
		Login string   `json:"login"`
		Roles []string `json:"roles"`
	}
	c.JSON(http.StatusOK, UserResponse{
		Login: reply.Login,
		Roles: reply.Roles,
	})
}

// @Summary      Authorizes user
// @Description  log in user with provided login and password
// @Tags         users
// @Produce      json
// @Param        request body users_pb.LoginRequest true "Request body"
// @Success      200  {object}  handlers.UserLogin.UserResponse "User successfully authorized"
// @Failure      400  {object}  middleware.CustomError "Invalid request data"
// @Failure      410 {object}  middleware.CustomError "Service does not responding (maybe crush)"
// @Router       /users/login [post]
func (h *userHandler) UserLogin(c *gin.Context) {
	var request users_pb.LoginRequest
	if err := c.BindJSON(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	type UserResponse struct {
		Login string   `json:"login"`
		Roles []string `json:"roles"`
	}
	reply, err := h.client.Login(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}
	c.SetCookie(middleware.TokenCookieName, reply.Token, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)
	c.SetCookie(middleware.RefreshCookieName, reply.Refreshtoken, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)

	c.JSON(http.StatusOK, UserResponse{
		Login: reply.Login,
		Roles: reply.Roles,
	})
}

// @Summary      Registration
// @Description  creates new user and authorizes it
// @Tags         users
// @Produce      json
// @Param        request body users_pb.RegistrationRequest true "Request body"
// @Success      200  {object}  handlers.UserRegister.UserResponse "User successfully authorized"
// @Failure      400  {object}  middleware.CustomError "Invalid request data"
// @Failure      410  {object}  middleware.CustomError "Service does not responding (maybe crush)"
// @Failure      500  {object}  middleware.CustomError "Some internal error occured"
// @Router       /users/signin [post]
func (h *userHandler) UserRegister(c *gin.Context) {
	var request users_pb.RegistrationRequest
	if err := c.BindJSON(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	type UserResponse struct {
		Login string   `json:"login"`
		Roles []string `json:"roles"`
	}
	reply, err := h.client.RegisterUser(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}
	c.SetCookie(middleware.TokenCookieName, reply.Token, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)
	c.SetCookie(middleware.RefreshCookieName, reply.Refreshtoken, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)

	c.JSON(http.StatusCreated, UserResponse{
		Login: reply.Login,
		Roles: reply.Roles,
	})
}

package handlers

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/middleware"
	_ "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary      Find user by credentials
// @Description  params goes in specific order: id -> login -> email
// @Description  first found user will be returned. If no user found, there'll be an error with details
// @Tags         users
// @Produce      json
// @Param		 id      query     string 		false 		"User Id"
// @Param		 login   query     string 		false 		"User login"
// @Param		 email   query     string 		false 		"User email" Format(email)
// @Success      200  {object}  users_pb.UserModel 		"User DTO model"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Request's field was not in a correct format"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"User not found"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Service does not responding (maybe crush)"
// @Router       /users [get]
func (h *userHandler) UserSearch(c *gin.Context) {

	var request users_pb.UserRequest
	if err := c.BindQuery(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	// need to made
	reply, err := h.Client.GetUser(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reply)
}

// @Summary      Authenticates user
// @Description  check if current user has legit token
// @Tags         users
// @Produce      json
// @Success      200  {object}  handlers.UserAuthenticate.UserResponse "User successfully authorized"
// @Failure      401  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "User does not authorized"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "User does not exists in database"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Service does not responding (maybe crush)"
// @Security 	 ApiKeyAuth
// @Router       /users/auth [get]
func (h *userHandler) UserAuthenticate(c *gin.Context) {
	user, err := middleware.GetCredentialsFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	reply, err := h.Client.GetUser(c.Request.Context(), &users_pb.UserRequest{Id: user.Id})
	if err != nil {
		c.Error(err)
		return
	}
	if !reflect.DeepEqual(user.Roles, reply.Roles) {
		h.logger.Infof("user's %s rights has changed, regenerating token...", reply.Login)
		refreshCookie, _ := c.Cookie(middleware.RefreshCookieName)
		tokenReply, err := h.Client.UpdateToken(c.Request.Context(), &users_pb.TokenRequest{Refreshtoken: refreshCookie})
		if err != nil {
			c.SetCookie(middleware.TokenCookieName, "", -1, "/", "", true, true)
			c.SetCookie(middleware.RefreshCookieName, "", -1, "/", "", true, true)
			c.Error(err)
			c.Abort()
			return
		}
		c.SetCookie(middleware.TokenCookieName, tokenReply.Token, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)
		c.SetCookie(middleware.RefreshCookieName, tokenReply.Refreshtoken, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)
	}
	h.logger.Infof("user %s authenticated with token and %v rights", reply.Login, reply.Roles)
	type UserResponse struct {
		Login string   `json:"login" example:"admin"`
		Roles []string `json:"roles" example:"user"`
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
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Invalid request data"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Service does not responding (maybe crush)"
// @Router       /users/login [post]
func (h *userHandler) UserLogin(c *gin.Context) {
	var request users_pb.LoginRequest
	if err := c.BindJSON(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	type UserResponse struct {
		Login string   `json:"login" example:"admin"`
		Roles []string `json:"roles" example:"user"`
	}
	reply, err := h.Client.Login(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}
	h.logger.Infof("user %s authoirized via login and password", request.Login)
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
// @Success      201  {object}  handlers.UserRegister.UserResponse "User registered and authorized"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Invalid request data"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Some internal error occured"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Service does not responding (maybe crush)"
// @Router       /users/signin [post]
func (h *userHandler) UserRegister(c *gin.Context) {
	var request users_pb.RegistrationRequest
	if err := c.BindJSON(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	type UserResponse struct {
		Login string   `json:"login" example:"admin"`
		Roles []string `json:"roles" example:"user"`
	}
	reply, err := h.Client.RegisterUser(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}
	h.logger.Infof("user %s registered with email %s", reply.Login, request.Email)
	c.SetCookie(middleware.TokenCookieName, reply.Token, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)
	c.SetCookie(middleware.RefreshCookieName, reply.Refreshtoken, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)

	c.JSON(http.StatusCreated, UserResponse{
		Login: reply.Login,
		Roles: reply.Roles,
	})
}

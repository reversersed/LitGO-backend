package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomError struct {
	Code      int32  `json:"code"`
	NamedCode string `json:"type"`
	Message   string `json:"message"`
	Details   []any  `json:"details"`
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	if lastError := c.Errors.Last(); lastError != nil {
		err, valid := status.FromError(lastError.Unwrap())
		if !valid {
			c.JSON(http.StatusInternalServerError, lastError.Error())
		} else {
			custom := CustomError{
				Code:      err.Proto().Code,
				NamedCode: err.Code().String(),
				Message:   err.Message(),
				Details:   err.Details(),
			}
			c.JSON(rpgCodeToHttpStatus(err.Code()), custom)
		}
	}
}

func rpgCodeToHttpStatus(code codes.Code) int {
	table := map[codes.Code]int{
		codes.OK:                 http.StatusOK,
		codes.Canceled:           http.StatusGone,
		codes.Unknown:            http.StatusInternalServerError,
		codes.InvalidArgument:    http.StatusBadRequest,
		codes.DeadlineExceeded:   http.StatusGatewayTimeout,
		codes.NotFound:           http.StatusNotFound,
		codes.AlreadyExists:      http.StatusConflict,
		codes.PermissionDenied:   http.StatusForbidden,
		codes.ResourceExhausted:  http.StatusTooManyRequests,
		codes.FailedPrecondition: http.StatusBadRequest,
		codes.Aborted:            http.StatusConflict,
		codes.OutOfRange:         http.StatusBadRequest,
		codes.Unimplemented:      http.StatusNotImplemented,
		codes.Internal:           http.StatusInternalServerError,
		codes.Unavailable:        http.StatusServiceUnavailable,
		codes.DataLoss:           http.StatusInternalServerError,
		codes.Unauthenticated:    http.StatusUnauthorized,
	}
	status, ok := table[code]
	if !ok {
		return http.StatusInternalServerError
	}
	return status
}

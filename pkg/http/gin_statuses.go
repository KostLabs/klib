package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Status struct {
	Code  int    `json:"code"`
	Enum  string `json:"enum"`
	Error string `json:"error,omitempty"`
}

type Response struct {
	Status Status `json:"status"`
	Data   any    `json:"data,omitempty"`
}

func StatusOK(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{
		Status: Status{
			Code: http.StatusOK,
			Enum: "SUCCESS",
		},
		Data: data,
	})
}

func StatusCreated(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusCreated, Response{
		Status: Status{
			Code: http.StatusCreated,
			Enum: "CREATED",
		},
		Data: data,
	})
}

func StatusAccepted(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusAccepted, Response{
		Status: Status{
			Code: http.StatusAccepted,
			Enum: "ACCEPTED",
		},
		Data: data,
	})
}

func StatusNoContent(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, Response{
		Status: Status{
			Code: http.StatusNoContent,
			Enum: "NO_CONTENT",
		},
	})
}

func StatusMovedPermanently(ctx *gin.Context, location string) {
	ctx.Header("Location", location)
	ctx.JSON(http.StatusMovedPermanently, Response{
		Status: Status{
			Code: http.StatusMovedPermanently,
			Enum: "MOVED_PERMANENTLY",
		},
	})
}

func StatusFound(ctx *gin.Context, location string) {
	ctx.Header("Location", location)
	ctx.JSON(http.StatusFound, Response{
		Status: Status{
			Code: http.StatusFound,
			Enum: "FOUND",
		},
	})
}

func StatusNotModified(ctx *gin.Context) {
	ctx.JSON(http.StatusNotModified, Response{
		Status: Status{
			Code: http.StatusNotModified,
			Enum: "NOT_MODIFIED",
		},
	})
}

func StatusTemporaryRedirect(ctx *gin.Context, location string) {
	ctx.Header("Location", location)
	ctx.JSON(http.StatusTemporaryRedirect, Response{
		Status: Status{
			Code: http.StatusTemporaryRedirect,
			Enum: "TEMPORARY_REDIRECT",
		},
	})
}

func StatusPermanentRedirect(ctx *gin.Context, location string) {
	ctx.Header("Location", location)
	ctx.JSON(http.StatusPermanentRedirect, Response{
		Status: Status{
			Code: http.StatusPermanentRedirect,
			Enum: "PERMANENT_REDIRECT",
		},
	})
}

func StatusBadRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, Response{
		Status: Status{
			Code:  http.StatusBadRequest,
			Enum:  "BAD_REQUEST",
			Error: err.Error(),
		},
	})
}

func StatusUnauthorized(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Status: Status{
			Code:  http.StatusUnauthorized,
			Enum:  "UNAUTHORIZED",
			Error: err.Error(),
		},
	})
}

func StatusForbidden(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusForbidden, Response{
		Status: Status{
			Code:  http.StatusForbidden,
			Enum:  "FORBIDDEN",
			Error: err.Error(),
		},
	})
}

func StatusNotFound(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusNotFound, Response{
		Status: Status{
			Code:  http.StatusNotFound,
			Enum:  "NOT_FOUND",
			Error: err.Error(),
		},
	})
}

func StatusMethodNotAllowed(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusMethodNotAllowed, Response{
		Status: Status{
			Code:  http.StatusMethodNotAllowed,
			Enum:  "METHOD_NOT_ALLOWED",
			Error: err.Error(),
		},
	})
}

func StatusConflict(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusConflict, Response{
		Status: Status{
			Code:  http.StatusConflict,
			Enum:  "CONFLICT",
			Error: err.Error(),
		},
	})
}

func StatusInternalServerError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, Response{
		Status: Status{
			Code:  http.StatusInternalServerError,
			Enum:  "INTERNAL_SERVER_ERROR",
			Error: err.Error(),
		},
	})
}

func StatusServiceUnavailable(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusServiceUnavailable, Response{
		Status: Status{
			Code:  http.StatusServiceUnavailable,
			Enum:  "SERVICE_UNAVAILABLE",
			Error: err.Error(),
		},
	})
}

func StatusGatewayTimeout(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusGatewayTimeout, Response{
		Status: Status{
			Code:  http.StatusGatewayTimeout,
			Enum:  "GATEWAY_TIMEOUT",
			Error: err.Error(),
		},
	})
}

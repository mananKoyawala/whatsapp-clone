package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var Unauthorized error = errors.New("permisson denied")

type ApiError struct { // to return error
	Error string `json:"error"`
}

type ApiMessage struct { // to return message
	Message string `json:"message"`
}

type ApiData struct { // to return data
	Data any `json:"data"`
}

type Handler interface {
	WriteJSON(c *gin.Context, status int, v any) error
}

type apiFunc func(c *gin.Context) (int, error)

func MakeHTTPHandleFunc(f apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if statusCode, err := f(c); err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(statusCode, ApiError{Error: err.Error()})
		}
	}
}

func WriteError(c *gin.Context, status int, v any) (int, error) {
	c.JSON(status, gin.H{"error": v})
	return 0, nil
}

func WriteMessage(c *gin.Context, status int, v string) (int, error) {
	c.JSON(status, ApiMessage{Message: v})
	return 0, nil
}

func WriteData(c *gin.Context, status int, v any) (int, error) {
	c.JSON(status, ApiData{Data: v})
	return 0, nil
}

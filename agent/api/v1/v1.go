package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SUCCESS = 200
	ERROR   = 500
)

var MessageFlags = map[int]string{
	SUCCESS: "ok",
	ERROR:   "failed",
}

func GetMsg(code int) string {
	if msg, ok := MessageFlags[code]; ok {
		return msg
	}
	return "ERROR NOT DEFINE"
}

type Controller struct {
}

func Setup() *Controller {
	return &Controller{}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Success int         `json:"success"`
	Data    interface{} `json:"data"`
}

func (c *Controller) OK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{Code: SUCCESS, Message: GetMsg(SUCCESS)})
}

func (c *Controller) Error(ctx *gin.Context, httpCode, errorCode int, message string) {
	ctx.JSON(httpCode, Response{Code: errorCode, Message: message})
}

func (c *Controller) Reply(ctx *gin.Context, httpCode int, errorCode int, data interface{}) {
	ctx.JSON(httpCode, Response{
		Code:    httpCode,
		Message: GetMsg(errorCode),
		Success: 1,
		Data:    data,
	})
}

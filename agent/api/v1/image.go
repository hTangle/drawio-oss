package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"super-markdown-editor-web/model"
	"super-markdown-editor-web/service"
)

func (c *Controller) NewDefaultOssData(ctx *gin.Context) {
	key := ctx.Query("key")
	if key == "" {
		c.Error(ctx, ERROR, ERROR, "error")
		return
	}
	url, ok := service.NewDefaultOssData(key)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "",
		"success": ok,
		"data":    url,
	})
}

func (c *Controller) GetOssDrawList(ctx *gin.Context) {
	key := ctx.Query("key")
	result, err := service.GetOssObjects(key)
	if err != nil {
		if key == "" {
			c.Error(ctx, ERROR, ERROR, err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "",
		"data":    result,
	})
}

func (c *Controller) GetImageBase64Data(ctx *gin.Context) {
	key := ctx.Query("key")
	if key == "" {
		c.Error(ctx, ERROR, ERROR, "error")
		return
	}
	result, err := service.GetOssObject(key)
	if err != nil {
		if key == "" {
			c.Error(ctx, ERROR, ERROR, err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "",
		"data":    result,
	})
}

func (c *Controller) SaveImage(ctx *gin.Context) {
	var req model.Image
	if err := ctx.ShouldBind(&req); err == nil {
		service.UploadDrawIODataToOss(ctx, req)
	} else {
		c.Error(ctx, ERROR, ERROR, "error")
	}
}

package agent

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	v1 "super-markdown-editor-web/agent/api/v1"
	"super-markdown-editor-web/model"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Static("/api/css", "./static/css")
	r.Static("/api/js", "./static/js")

	r.Use(gin.Logger())
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})
	r.LoadHTMLGlob("templates/**")
	r.GET("/draw", func(context *gin.Context) {
		context.HTML(http.StatusOK, "draw.html", nil)
	})

	r.GET("/login", func(context *gin.Context) {
		context.HTML(http.StatusOK, "login.html", nil)
	})

	controller := v1.Setup()
	loginInApi := r.Group("/login")
	loginInApi.POST("/auth", controller.Login)
	loginInApi.GET("/auth", func(context *gin.Context) {
		context.HTML(http.StatusOK, "login.html", nil)
	})

	apiV1 := r.Group(model.GroupPrefix)
	apiV1.Use(JWTAuth())
	apiV1.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
			"code":    200,
		})
	})
	apiV1.POST("/image", controller.SaveImage)
	apiV1.GET("/image", controller.GetImageBase64Data)
	apiV1.PUT("/image", controller.NewDefaultOssData)
	apiV1.GET("/images", controller.GetOssDrawList)
	return r
}

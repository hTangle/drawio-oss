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

	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.Static("/svg", "./static/svg")
	r.Static("/themes", "./static/themes")
	r.Static("/fonts", "./static/fonts")
	r.Static("/image", "./static/image")

	r.Use(gin.Logger())
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})
	r.LoadHTMLGlob("templates/**")
	r.GET("/main", func(context *gin.Context) {
		context.HTML(http.StatusOK, "main.html", nil)
	})
	r.GET("/show", func(context *gin.Context) {
		context.HTML(http.StatusOK, "show.html", nil)
	})

	r.GET("/list", func(context *gin.Context) {
		context.HTML(http.StatusOK, "list.html", nil)
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
	apiV1.Static("/draw", "./static/drawio")

	apiV1.GET("/content", controller.GetNote)
	apiV1.POST("/content", controller.WriteNote)
	apiV1.POST("/image", controller.SaveImage)
	apiV1.GET("/tree", controller.GetBookTree)
	apiV1.POST("/rename", controller.RenameBookOrNote)
	apiV1.POST("/create", controller.CreateNote)

	return r
}

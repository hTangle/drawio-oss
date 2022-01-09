package agent

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"super-markdown-editor-web/model"
)

func RunAgent(debug bool, conf *model.EditorConf) {
	if debug {
		gin.SetMode("debug")
	}
	routerInit := InitRouter()
	endPoint := fmt.Sprintf(":%d", conf.Http.Port)
	logrus.Warnf("start to lister: %s", endPoint)
	maxHeaderBytes := 1 << 20
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routerInit,
		MaxHeaderBytes: maxHeaderBytes,
	}
	if err := server.ListenAndServe(); err != nil {
		logrus.Errorf("close http server error: %v", err)
	}
}

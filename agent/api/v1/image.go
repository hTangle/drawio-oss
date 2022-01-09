package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"strings"
	"super-markdown-editor-web/model"
)

func (c *Controller) SaveImage(ctx *gin.Context) {
	//file, err := ctx.FormFile("editormd-image-file")
	file, err := ctx.FormFile("file")
	if err != nil {
		c.Error(ctx, ERROR, ERROR, fmt.Sprintf("get form error: %s", err.Error()))
		return
	}
	basePath := model.GetLocalEditorConf("").ImageDir
	arr := strings.Split(file.Filename, ".")
	fileType := arr[len(arr)-1]
	logrus.Debugf("file type is: %s", fileType)
	//newFileName := fmt.Sprintf("%s.%s", strings.ReplaceAll(uuid.New().String(), "-", ""), fileType)
	newFileName := fmt.Sprintf("%s.%s", strings.ReplaceAll(uuid.New().String(), "-", ""), fileType)
	filename := path.Join(basePath, newFileName)
	logrus.Warnf("save file to %s", newFileName)
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		c.Error(ctx, ERROR, ERROR, fmt.Sprintf("upload file error: %s", err.Error()))
		return
	}
	logrus.Warnf("upload file success! %s", filename)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": 1,
		"code":    200,
		"data": map[string]interface{}{
			"errFiles": []string{},
			"succMap": map[string]string{
				newFileName: "/image/" + newFileName,
			},
		},
		"message": "success",
	})

}

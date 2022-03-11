package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		logrus.Warnf("alt: %s, xml: %s data: %x", req.Alt, req.Xml, req.Data)
		service.UploadDrawIODataToOss(ctx, req)
	} else {
		c.Error(ctx, ERROR, ERROR, "error")
	}
	//file, err := ctx.FormFile("editormd-image-file")
	//file, err := ctx.FormFile("file")
	//if err != nil {
	//	c.Error(ctx, ERROR, ERROR, fmt.Sprintf("get form error: %s", err.Error()))
	//	return
	//}
	//basePath := model.GetLocalEditorConf("").ImageDir
	//arr := strings.Split(file.Filename, ".")
	//fileType := arr[len(arr)-1]
	//logrus.Debugf("file type is: %s", fileType)
	//newFileName := fmt.Sprintf("%s.%s", strings.ReplaceAll(uuid.New().String(), "-", ""), fileType)
	//newFileName := fmt.Sprintf("%s.%s", strings.ReplaceAll(uuid.New().String(), "-", ""), fileType)
	//filename := path.Join(basePath, newFileName)
	//logrus.Warnf("save file to %s", newFileName)
	//if service.UploadImageToOss(ctx, file, newFileName) {
	//	logrus.Debugf("upload file to s3 success")
	//	ctx.JSON(http.StatusOK, map[string]interface{}{
	//		"success": 1,
	//		"code":    200,
	//		"data": map[string]interface{}{
	//			"errFiles": []string{},
	//			"succMap": map[string]string{
	//				newFileName: "https://image.ahsup.top/image/" + newFileName,
	//			},
	//		},
	//		"message": "success",
	//	})
	//	return
	//}
	//
	//if err := ctx.SaveUploadedFile(file, filename); err != nil {
	//	c.Error(ctx, ERROR, ERROR, fmt.Sprintf("upload file error: %s", err.Error()))
	//	return
	//}
	//logrus.Warnf("upload file success! %s", filename)
	//ctx.JSON(http.StatusOK, map[string]interface{}{
	//	"success": 1,
	//	"code":    200,
	//	"data": map[string]interface{}{
	//		"errFiles": []string{},
	//		"succMap": map[string]string{
	//			newFileName: "/image/" + newFileName,
	//		},
	//	},
	//	"message": "success",
	//})

}

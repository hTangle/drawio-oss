package service

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"super-markdown-editor-web/model"
)

func UploadImageToOss(ctx context.Context, file *multipart.FileHeader, newFileName string) bool {
	if model.LocalEditorConf.Oss.IsEmpty() {
		return false
	}
	client, err := oss.New(model.LocalEditorConf.Oss.Endpoint, model.LocalEditorConf.Oss.Ak, model.LocalEditorConf.Oss.Sk)
	if err != nil {
		logrus.Errorf("init client error: %v", err)
		return false
	}
	bucket, err := client.Bucket(model.LocalEditorConf.Oss.Bucket)
	if err != nil {
		logrus.Errorf("init bucket error: %v", err)
		return false
	}
	fileContent, _ := file.Open()
	defer fileContent.Close()
	err = bucket.PutObject("image/"+newFileName, fileContent)
	if err != nil {
		logrus.Errorf("put file to bucket error: %v", err)
		return false
	}
	return true
}

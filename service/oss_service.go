package service

import (
	"context"
	"encoding/base64"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"super-markdown-editor-web/model"
	"sync"
)

var (
	ossClient *oss.Client
	ossOnce   = &sync.Once{}
)

const (
	OssDataPath = "./data"
)

var (
	TmpPath      = "/tmp"
	ossDataCache = &model.OssListCache{}
)

func init() {
	if _, err := os.Stat(OssDataPath); err != nil {
		os.Mkdir(OssDataPath, 0777)
	}
	dev := os.Getenv("IS_DEV")
	if dev != "" {
		TmpPath = "./tmp"
	}
	if _, err := os.Stat(TmpPath); err != nil {
		os.Mkdir(TmpPath, 0777)
	}

}

func getOssHttpsPrefix() string {
	return model.LocalEditorConf.Oss.Url
}

func getOssClient() *oss.Client {
	if ossClient == nil {
		ossOnce.Do(func() {
			var err error
			ossClient, err = oss.New(model.LocalEditorConf.Oss.Endpoint, model.LocalEditorConf.Oss.Ak, model.LocalEditorConf.Oss.Sk)
			if err != nil {
				panic(err)
			}
			objs, err := GetOssObjectsFromRemote("")
			ossDataCache.InitCache(objs)
		})
	}
	return ossClient
}

func GetOssDrawIOData(filter string) {
	client := getOssClient()
	bucket, err := client.Bucket(model.LocalEditorConf.Oss.Bucket)
	if err != nil {
		logrus.Errorf("init bucket error: %v", err)
		return
	}
	bucket.ListObjectsV2(oss.Prefix(""))
}

func tranKeyToUrl(key string) string {
	return getOssHttpsPrefix() + "/draw/" + key + ".png"
}

func NewDefaultOssData(key string) (url string, ok bool) {
	if ossDataCache.HasObject(key) {
		return "", false
	}
	image := model.NewDefaultImage(key)
	ossDataCache.PutOssObject(key, getOssHttpsPrefix()+"/draw/"+key+".png")
	return getOssHttpsPrefix() + "/draw/" + key + ".png", UploadDrawIODataToOss(context.Background(), image)
}

func GetOssObjects(key string) (results map[string]string, err error) {
	_ = getOssClient()
	return ossDataCache.GetOssObjects(key), nil
}

func GetOssObjectsFromRemote(key string) (results map[string]string, err error) {
	results = map[string]string{}
	client := getOssClient()
	bucket, err := client.Bucket(model.LocalEditorConf.Oss.Bucket)
	if err != nil {
		logrus.Errorf("init bucket error: %v", err)
		return nil, err
	}
	objs, err := bucket.ListObjectsV2(oss.Prefix("draw/" + key))
	if err != nil || objs.Objects == nil {
		logrus.Errorf("init bucket error: %v", err)
		return nil, err
	}
	for _, obj := range objs.Objects {
		if strings.HasSuffix(obj.Key, ".png") {
			results[strings.ReplaceAll(strings.ReplaceAll(obj.Key, "draw/", ""), ".png", "")] = getOssHttpsPrefix() + obj.Key
		}
	}
	return
}

func GetOssObject(key string) (string, error) {
	client := getOssClient()
	bucket, err := client.Bucket(model.LocalEditorConf.Oss.Bucket)
	if err != nil {
		logrus.Errorf("init bucket error: %v", err)
		return "", err
	}
	obj, err := bucket.GetObject("draw/" + key + ".data")
	if err != nil {
		logrus.Errorf("init bucket error: %v", err)
		return "", err
	}
	defer obj.Close()
	data, err := ioutil.ReadAll(obj)
	return string(data), err
}

func UploadDrawIODataToOss(ctx context.Context, image model.Image) bool {
	if model.LocalEditorConf.Oss.IsEmpty() {
		return false
	}
	if image.Alt == "" {
		return false
	}
	if image.Data == "" && image.Xml == "" {
		return false
	}
	client := getOssClient()
	bucket, err := client.Bucket(model.LocalEditorConf.Oss.Bucket)
	if err != nil {
		logrus.Errorf("init bucket error: %v", err)
		return false
	}
	if image.Data != "" {
		tranStringToBucket(bucket, image.Alt, image.Data, "data")
		add, _ := base64.StdEncoding.DecodeString(image.Data[22:])
		tranStringToBucket(bucket, image.Alt, string(add), "png")
	}
	if image.Xml != "" {
		tranStringToBucket(bucket, image.Alt, image.Xml, "xml")
	}
	ossDataCache.PutOssObject(image.Alt, tranKeyToUrl(image.Alt))
	return true
}

func tranStringToBucket(bucket *oss.Bucket, name string, data string, type_ string) {
	if _, err := os.Stat(TmpPath); err != nil {
		os.Mkdir(TmpPath, 0777)
	}
	p := path.Join(TmpPath, name+"."+type_)
	err := ioutil.WriteFile(p, []byte(data), 0666)
	if err != nil {
		logrus.Errorf("write %s error: %v", data, err)
	}
	err = bucket.PutObjectFromFile("draw/"+name+"."+type_, p)
	if err != nil {
		logrus.Errorf("put %s to oss error: %v", data, err)
	}
}

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

package model

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type HttpConf struct {
	Port      int    `json:"port" yaml:"port"`
	SignedKey string `json:"signed_key" yaml:"signed_key"`
}

type OssConf struct {
	Ak       string `json:"ak"`
	Sk       string `json:"sk"`
	Bucket   string `json:"bucket"`
	Region   string `json:"region"`
	Endpoint string `json:"endpoint"`
}

func (o *OssConf) IsEmpty() bool {
	return o.Ak == "" || o.Sk == "" || o.Bucket == "" || o.Region == "" || o.Endpoint == ""
}

type EditorConf struct {
	Http      HttpConf          `json:"http" yaml:"http"`
	WorkDir   string            `json:"work_dir" yaml:"work_dir"`
	ImageDir  string            `json:"image_dir" yaml:"image_dir"`
	UserInfos map[string]string `json:"user_infos" yaml:"user"`
	Oss       OssConf           `json:"oss"`
}

func GetDefaultEditorConf() *EditorConf {
	workBase, _ := os.Getwd()
	logrus.Warnf("workdir: %s", workBase)
	return &EditorConf{
		Http: HttpConf{
			Port: 8080,
		},
		WorkDir:  "./app",
		ImageDir: path.Join(workBase, "static", "image"),
	}
}

func (e *EditorConf) GetWorkDir() string {
	return e.WorkDir
}

func (e *EditorConf) SignedKey() string {
	return e.Http.SignedKey
}

func (e *EditorConf) IsValidUser(user, pwd string) bool {
	if val, ok := e.UserInfos[user]; ok && pwd == val {
		return true
	}
	return false
}

var (
	LocalEditorConf *EditorConf
	ConfOnce        sync.Once
)

func init() {
	confPath := "./config.json"
	if _, err := os.Stat(confPath); err != nil {
		LocalEditorConf = GetDefaultEditorConf()
	} else {
		if data, err := ioutil.ReadFile(confPath); err == nil {
			if err = json.Unmarshal(data, &LocalEditorConf); err == nil {
				logrus.Warnf("read conf from %s", confPath)
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	if _, err := os.Stat(LocalEditorConf.GetWorkDir()); err != nil {
		os.MkdirAll(LocalEditorConf.GetWorkDir(), 0777)
	}
	hpp := path.Join(LocalEditorConf.GetWorkDir(), BlogHtmlPath)
	if _, err := os.Stat(hpp); err != nil {
		os.MkdirAll(hpp, 0766)
	}
	Blogs.SetHtmlPath(hpp)

	pp := path.Join(LocalEditorConf.GetWorkDir(), PublisherPath)
	if _, err := os.Stat(pp); err != nil {
		err = ioutil.WriteFile(pp, []byte("{\"blogs\":[]}"), 0766)
		if err != nil {
			panic(err)
		}
	}
	data, err := ioutil.ReadFile(pp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &Blogs)
	if err != nil {
		panic(err)
	}
	Blogs.SetPath(pp)
	Blogs.InitBlog()
}

func GetLocalEditorConf(confPath string) *EditorConf {
	return LocalEditorConf
}

func GetLocalEditorConfig() *EditorConf {
	return LocalEditorConf
}

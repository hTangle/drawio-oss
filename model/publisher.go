package model

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"strings"
	"super-markdown-editor-web/util"
	"sync"
	"time"
)

const (
	PublisherPath  = "blog_info.json"
	BlogHtmlPath   = "html"
	ShortCutLength = 128
)

var Blogs BlogList

type Publisher struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	ReleaseTime int64  `json:"release_time"`
	UpdateTime  string `json:"update_time"`
	Cover       string `json:"cover"`
	ShortCut    string `json:"short_cut"`
}

type BlogList struct {
	Blogs   []*Publisher          `json:"blogs"` //按照时间戳从大到小排列
	BlogMap map[string]*Publisher `json:"-"`
	sync.RWMutex
	Path     string `json:"-"`
	HtmlPath string
}

func (b *BlogList) AddPostPage(title, id, html string) {
	util.SavePostToHTML(id, title, html, b.HtmlPath)
}

func (b *BlogList) SetHtmlPath(p string) {
	b.HtmlPath = p
}

func (b *BlogList) SetPath(p string) {
	b.Path = p
}
func (b *BlogList) getShortCutAndCover(content string) (cover string, shortCut string) {
	re, _ := regexp.Compile("!\\[[A-Za-z0-9-.]{0,40}\\]\\([A-Za-z0-9/.:]{30,70}\\.png\\)") //![](.png) ![](/image/ebb6b2a03dbf4d3fafaa185e1ea40e04.png)
	one := string(re.Find([]byte(content)))
	index := strings.Index(one, "](")
	if index > 0 {
		cover = one[index+2 : len(one)-1]
	} else {
		logrus.Warnf("no cover found in content")
	}
	endLine := ShortCutLength
	runStr := []rune(content)

	if len(runStr) < endLine {
		endLine = len(runStr)
	}
	shortCut = string(runStr[0:endLine])
	return
}

func (b *BlogList) insertPublisher(p *Publisher) {
	b.Blogs = append(b.Blogs, p)
	b.BlogMap[p.Id] = p
}

func (b *BlogList) AddAPublisher(title, id, content string) {
	b.Lock()
	defer b.Unlock()
	cover, shortCut := b.getShortCutAndCover(content)
	if _, ok := b.BlogMap[id]; ok {
		b.BlogMap[id].Cover = cover
		b.BlogMap[id].ShortCut = shortCut
		b.BlogMap[id].Title = title
		b.BlogMap[id].UpdateTime = time.Now().Format("2006-01-02 15:04")
	} else {
		pub := &Publisher{
			Id:          id,
			Title:       title,
			ReleaseDate: time.Now().Format("2006-01-02 15:04"),
			ReleaseTime: time.Now().UnixNano(),
			UpdateTime:  time.Now().Format("2006-01-02 15:04"),
			Cover:       cover,
			ShortCut:    shortCut,
		}
		b.insertPublisher(pub)
	}
	b.syncToDisk()
}

func (b *BlogList) syncToDisk() {
	data, err := json.Marshal(&b)
	if err != nil {
		logrus.Errorf("Marshal blogs error: %v", err)
	}
	err = ioutil.WriteFile(b.Path, data, 0766)
	if err != nil {
		logrus.Errorf("sync blogs to disk error: %v", err)
	}
}

func (b *BlogList) Length() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.Blogs)
}

func (b *BlogList) GetBlogsPure(offset, size int) []*Publisher {
	ps, _ := b.GetBlogs(offset, size)
	return ps
}

func (b *BlogList) GetBlogs(offset, size int) ([]*Publisher, int) {
	b.RLock()
	defer b.RUnlock()
	length := len(b.Blogs)
	begin := length - offset
	if begin < 0 {
		return []*Publisher{}, length
	}
	end := begin - size
	if end < 0 {
		end = 0
	}
	if end > begin {
		end = begin
	}
	return b.Blogs[end:begin], length

}

func (b *BlogList) InitBlog() {
	if b.Blogs == nil {
		b.Blogs = []*Publisher{}
	}
	b.BlogMap = map[string]*Publisher{}
	for _, p := range b.Blogs {
		b.BlogMap[p.Id] = p
	}
}

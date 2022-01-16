package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"super-markdown-editor-web/model"
	"super-markdown-editor-web/service"
)

type PublishNoteRequest struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Html    string `json:"html"`
}

func (c *Controller) PublishNote(ctx *gin.Context) {
	var request PublishNoteRequest
	if err := ctx.ShouldBind(&request); err == nil {
		if request.Id != "" {
			service.WriteNote(request.Id, request.Content)
			service.GenerateHtml(model.GetLocalShowTrees().GetNoteTitle(request.Id), request.Id, request.Content)
		} else {
			logrus.Warnf("request id is null")
		}
		c.OK(ctx)
	} else {
		c.Error(ctx, ERROR, ERROR, "params error")
	}
}

func (c *Controller) GetBlogs(ctx *gin.Context) {
	size := ctx.DefaultQuery("size", "12")
	offset := ctx.DefaultQuery("offset", "0")
	size_, err := strconv.Atoi(size)
	if err != nil {
		logrus.Errorf("tran size:%s to int error", size)
		size_ = 12
	}
	offset_, err := strconv.Atoi(offset)
	if err != nil {
		logrus.Errorf("tran offset:%s to int error", offset)
		offset_ = 0
	}
	if offset_ < 0 {
		c.Error(ctx, ERROR, ERROR, "offset should great than -1")
		return
	}
	blogs, length := service.GetBlogs(offset_, size_)
	c.Reply(ctx, SUCCESS, SUCCESS, map[string]interface{}{
		"blogs":     blogs,
		"length":    length,
		"page_size": size_,
	})
}

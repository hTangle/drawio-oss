package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"super-markdown-editor-web/service"
)

func (c *Controller) GetNote(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		c.Error(ctx, ERROR, ERROR, "params error")
		return
	}
	if note, err := service.GetNote(id); err == nil {
		c.Reply(ctx, SUCCESS, SUCCESS, map[string]string{
			"data":  note,
			"title": service.GetNoteShowName(id),
		})
	} else {
		c.Error(ctx, ERROR, ERROR, "params error")
	}
}

type WriteNoteRequest struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func (c *Controller) WriteNote(ctx *gin.Context) {
	var request WriteNoteRequest
	if err := ctx.ShouldBind(&request); err == nil {
		if request.Id != "" {
			service.WriteNote(request.Id, request.Content)
		} else {
			logrus.Warnf("request id is null")
		}
		c.OK(ctx)
	} else {
		c.Error(ctx, ERROR, ERROR, "params error")
	}
}

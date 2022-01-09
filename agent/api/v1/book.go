package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"super-markdown-editor-web/model"
	"super-markdown-editor-web/service"
)

func (c *Controller) GetBookTree(ctx *gin.Context) {
	c.Reply(ctx, SUCCESS, SUCCESS, service.GetBookShowTree())
}

type RenameBookRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (c *Controller) RenameBookOrNote(ctx *gin.Context) {
	var req RenameBookRequest
	if err := ctx.ShouldBind(&req); err == nil {
		if req.Id != "" && req.Name != "" {
			logrus.Debugf("rename %s -> %s", req.Id, req.Name)
			service.RenameBookOrNote(req.Id, req.Name)
		}
	}
	c.OK(ctx)
}

type CreateNodeRequest struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	ParentId string `json:"parent_id"`
}

func (c *Controller) CreateNote(ctx *gin.Context) {
	var req CreateNodeRequest
	if err := ctx.ShouldBind(&req); err == nil {
		if req.Id != "" && req.Name != "" && req.Type != "" {
			switch req.Type {
			case model.TypeDefault:
				service.CreateBook(req.Id, req.Name, req.ParentId)
			case model.TypeFile:
				service.AddNewNote(req.ParentId, req.Id, req.Name)
			}
		}
	}
	c.OK(ctx)
}

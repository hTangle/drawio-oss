package service

import "super-markdown-editor-web/model"

func PublishNote(title, id, content string) {
	model.Blogs.AddAPublisher(title, id, content)
}

func GenerateHtml(title, id, html string) {
	model.Blogs.AddPostPage(title, id, html)
}

func GetBlogs(offset, size int) ([]*model.Publisher, int) {
	return model.Blogs.GetBlogs(offset, size)
}

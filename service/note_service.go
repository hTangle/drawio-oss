package service

import "super-markdown-editor-web/model"

func GetNote(id string) (string, error) {
	return model.GetLocalShowTrees().ReadNote(id)
}

func GetNoteShowName(id string) string {
	return model.GetLocalShowTrees().GetNoteTitle(id)
}

func AddNewNote(parentId, id, showName string) {
	model.GetLocalShowTrees().AddNewNote(parentId, id, showName)
}

func WriteNote(id, content string) {
	model.GetLocalShowTrees().WriteNote(id, content)
}

func CreateNote(id, name, parentId string) {
	model.GetLocalShowTrees().AddNewNote(parentId, id, name)
}

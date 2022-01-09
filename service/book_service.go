package service

import "super-markdown-editor-web/model"

func GetBookShowTree() []*model.TreeNode {
	return model.GetLocalShowTrees().GenerateShowTrees()
}

func RenameBookOrNote(id, newName string) {
	model.GetLocalShowTrees().RenameNote(id, newName)
}

func CreateBook(id, name, parentId string) {
	model.GetLocalShowTrees().AddNewBook(parentId, id, name)
}

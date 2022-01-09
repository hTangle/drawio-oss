package model

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

const (
	TypeFile = "file"
	BookPath = "book"
	NotePath = "note"

	TypeDefault = "default"
	WelcomeData = `
### 主要特性
- 支持“标准”Markdown / CommonMark和Github风格的语法，也可变身为代码编辑器；
- 支持实时预览、图片（跨域）上传、预格式文本/代码/表格插入、代码折叠、搜索替换、只读模式、自定义样式主题和多语言语法高亮等功能；
- 支持ToC（Table of Contents）、Emoji表情、Task lists、@链接等Markdown扩展语法；
- 支持TeX科学公式（基于KaTeX）、流程图 Flowchart 和 时序图 Sequence Diagram;
- 支持识别和解析HTML标签，并且支持自定义过滤标签解析，具有可靠的安全性和几乎无限的扩展性；
- 支持 AMD / CMD 模块化加载（支持 Require.js & Sea.js），并且支持自定义扩展插件；
- 兼容主流的浏览器（IE8+）和Zepto.js，且支持iPad等平板设备；
- 支持自定义主题样式；`
	DefaultBookId = "00000000-0000-0000-0000-000000000000"
	DefaultNoteId = "00000000-0000-0000-0000-000000000001"
)

type Book struct {
	ParentId   string    `json:"parent_id"`
	ShowName   string    `json:"show_name"`
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	UpdateTime time.Time `json:"update_time"`
	CreateTime time.Time `json:"create_time"`
}

func (b *Book) UpdateShowName(showName string) {
	b.ShowName = showName
	b.UpdateTime = time.Now()
}

type TreeNode struct {
	Type     string      `json:"type"`
	Text     string      `json:"text"`
	Id       string      `json:"id"`
	Children []*TreeNode `json:"children,omitempty"`
}

func (t *TreeNode) SetType(type_ string) {
	t.Type = type_
}

func (t *TreeNode) SetText(text string) {
	t.Text = text
}

func (t *TreeNode) SetId(id string) {
	t.Id = id
}

func (t *TreeNode) AddChildren(tn *TreeNode) {
	if t.Children == nil {
		t.Children = []*TreeNode{}
	}
	t.Children = append(t.Children, tn)
}

type ShowTrees struct {
	Tree          []*TreeNode
	IDInfo        map[string]*Book
	IDNameMapping map[string]map[string]bool
	RootNode      map[string]bool
	sync.RWMutex
}

var (
	LocalShowTree *ShowTrees
	once          sync.Once
)

func (s *ShowTrees) GetTree() []*TreeNode {
	s.RLock()
	defer s.RUnlock()
	return s.Tree
}

func (s *ShowTrees) AddIDInfoWithoutLock(id string, book *Book) {
	s.IDInfo[id] = book
	if book.ParentId != "" {
		s.AddIdNameMappingWithoutLock(id, book.ParentId)
	} else {
		s.AddRootNodeWithoutLock(id)
	}
}

func (s *ShowTrees) AddRootNodeWithoutLock(id string) {
	if s.RootNode == nil {
		s.RootNode = map[string]bool{}
	}
	s.RootNode[id] = true
}

func (s *ShowTrees) AddIdNameMappingWithoutLock(id, parentId string) {
	if s.IDNameMapping == nil {
		s.IDNameMapping = map[string]map[string]bool{}
	}
	if _, ok := s.IDNameMapping[parentId]; !ok {
		s.IDNameMapping[parentId] = map[string]bool{
			id: true,
		}
	} else {
		s.IDNameMapping[parentId][id] = true
	}
}

func (s *ShowTrees) AddNewBook(parentId, id, showName string) {
	s.Lock()
	defer s.Unlock()
	if parentId == "#" {
		parentId = ""
	}
	book := &Book{
		ParentId:   parentId,
		ShowName:   showName,
		ID:         id,
		Type:       TypeDefault,
		UpdateTime: time.Now(),
		CreateTime: time.Now(),
	}
	s.addBookType(parentId, id, book)
}

func (s *ShowTrees) addBookType(parentId, id string, book *Book) {
	s.IDInfo[id] = book
	s.syncChangeToDisk(id)
	s.addIdNameMapping(parentId, id)
}

func (s *ShowTrees) RenameNote(id, newName string) {
	s.Lock()
	defer s.Unlock()
	s.IDInfo[id].UpdateShowName(newName)
	s.syncChangeToDisk(id)
}

func (s *ShowTrees) getNotePath(id string) string {
	return path.Join(GetLocalEditorConf("").GetWorkDir(), NotePath, id)
}

func (s *ShowTrees) getBookPath(id string) string {
	return path.Join(GetLocalEditorConf("").GetWorkDir(), BookPath, id)
}

func (s *ShowTrees) GetNoteTitle(id string) string {
	if val, ok := s.IDInfo[id]; ok {
		return val.ShowName
	}
	return ""
}

func (s *ShowTrees) ReadNote(id string) (string, error) {
	if data, err := os.ReadFile(s.getNotePath(id)); err == nil {
		return string(data), nil
	} else {
		return "", err
	}
}

func (s *ShowTrees) AddNewNote(parentId, id, showName string) {
	s.Lock()
	defer s.Unlock()
	if parentId == "#" {
		parentId = ""
	}
	book := &Book{
		ParentId:   parentId,
		ShowName:   showName,
		ID:         id,
		Type:       TypeFile,
		UpdateTime: time.Now(),
		CreateTime: time.Now(),
	}
	s.addBookType(parentId, id, book)
	s.writeNote(id, "")
}
func (s *ShowTrees) syncChangeToDisk(id string) {
	savePath := path.Join(GetLocalEditorConf("").GetWorkDir(), BookPath, id)
	if data, err := json.Marshal(s.IDInfo[id]); err == nil {
		if err = ioutil.WriteFile(savePath, data, 0766); err != nil {
			logrus.Errorf("write book error: %s, error: %v", string(data), err)
		}
	} else {
		logrus.Errorf("Marshal book error: %v", err)
	}
}

func (s *ShowTrees) WriteNote(id, content string) {
	s.writeNote(id, content)
}

func (s *ShowTrees) writeNote(id, content string) {
	savePath := path.Join(GetLocalEditorConf("").GetWorkDir(), NotePath, id)
	if err := ioutil.WriteFile(savePath, []byte(content), 0766); err != nil {
		logrus.Errorf("write note error: %s, error: %v", content, err)
	}
}

func (s *ShowTrees) addIdNameMapping(parentId, id string) {
	if parentId == "" {
		s.RootNode[id] = true
	} else {
		if _, ok := s.IDNameMapping[parentId]; !ok {
			s.IDNameMapping[parentId] = map[string]bool{}
		}
		s.IDNameMapping[parentId][id] = true
	}
}

func (s *ShowTrees) InitWelcomeData() {
	if s.RootNode == nil {
		s.RootNode = map[string]bool{}
	}

	if len(s.RootNode) != 0 {
		return
	}
	logrus.Warnf("should init welcome data")
	os.MkdirAll(path.Join(GetLocalEditorConf("").GetWorkDir(), NotePath), 0777)
	os.MkdirAll(path.Join(GetLocalEditorConf("").GetWorkDir(), BookPath), 0777)
	s.AddNewBook("", DefaultBookId, "default")
	s.AddNewNote(DefaultBookId, DefaultNoteId, "welcome")
	s.writeNote(DefaultNoteId, WelcomeData)
	s.addRootNote(DefaultBookId)
	s.RootNode[DefaultBookId] = true
	s.addIdNameMapping(DefaultBookId, DefaultNoteId)
}

func (s *ShowTrees) GenerateShowTrees() []*TreeNode {
	s.RLock()
	defer s.RUnlock()
	return s.generateShowTree(s.RootNode)
}

func (s *ShowTrees) generateShowTree(rootNode map[string]bool) []*TreeNode {
	logrus.Debugf("start to get tree")
	var tree []*TreeNode
	for key, _ := range rootNode {
		logrus.Debugf("current node is: %s", key)
		if s.IDInfo[key].Type == TypeFile {
			tree = append(tree, &TreeNode{
				Type:     TypeFile,
				Text:     s.IDInfo[key].ShowName,
				Id:       key,
				Children: nil,
			})
		} else if v_, ok := s.IDNameMapping[key]; ok {
			logrus.WithField("id_name_mapping", v_).Debugf("get from id name mapping")
			childrenTree := s.generateShowTree(v_)
			if len(childrenTree) > 0 {
				tree = append(tree, &TreeNode{
					Type:     TypeDefault,
					Text:     s.IDInfo[key].ShowName,
					Id:       key,
					Children: childrenTree,
				})
			} else {
				tree = append(tree, &TreeNode{
					Type: TypeDefault,
					Text: s.IDInfo[key].ShowName,
					Id:   key,
				})
			}
		} else if v_, ok := s.IDInfo[key]; ok {
			tree = append(tree, &TreeNode{
				Type: TypeDefault,
				Text: v_.ShowName,
				Id:   key,
			})
		} else {
			logrus.Errorf("%s not exist in id_name_mapping", key)
		}
	}
	return tree
}

func (s *ShowTrees) addRootNote(id string) {
	if s.RootNode == nil {
		s.RootNode = map[string]bool{}
	}
	s.RootNode[id] = true
}

func GetLocalShowTrees() *ShowTrees {
	if LocalShowTree == nil {
		once.Do(func() {
			if LocalShowTree == nil {
				LocalShowTree = InitLocalShowTreeByPath()
			}
		})
	}
	return LocalShowTree
}

func InitLocalShowTreeByPath() *ShowTrees {
	showTree := &ShowTrees{
		Tree:          []*TreeNode{},
		IDInfo:        map[string]*Book{},
		IDNameMapping: map[string]map[string]bool{},
		RootNode:      map[string]bool{},
	}
	workDir := GetLocalEditorConf("").GetWorkDir()
	bookDir := path.Join(workDir, BookPath)
	if books, err := os.ReadDir(bookDir); err == nil {
		for _, book := range books {
			if !book.IsDir() {
				if data, err := ioutil.ReadFile(path.Join(bookDir, book.Name())); err == nil {
					var book_ *Book //
					if err = json.Unmarshal(data, &book_); err == nil {
						showTree.IDInfo[book.Name()] = book_
						if book_.ParentId != "" && book_.ParentId != "#" {
							showTree.AddIdNameMappingWithoutLock(book_.ID, book_.ParentId)
						} else {
							showTree.AddRootNodeWithoutLock(book_.ID)
						}
					} else {
						logrus.Errorf("Unmarshal book error: %v", err)
					}
				} else {
					logrus.Errorf("read book error: %s, error: %v", path.Join(bookDir, book.Name()), err)
				}
			}
		}
	} else {
		logrus.Errorf("read dir error: %s, error: %v", workDir, err)
	}
	showTree.InitWelcomeData()
	return showTree
}

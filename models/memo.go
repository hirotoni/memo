package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hirotoni/memo/markdown"
)

type Memo struct {
	Filepath string
	Title    string
	Content  string
}

func NewMemo(filepath, title, content string) *Memo {
	return &Memo{
		Filepath: filepath,
		Title:    title,
		Content:  content,
	}
}

func (m *Memo) SearchKey() string {
	return filepath.Base(m.Filepath) + "#" + markdown.Text2tag(m.Title)
}

func (m *Memo) Link() string {
	link := ".." + string(os.PathSeparator) + m.Filepath + "#" + markdown.Text2tag(m.Title)
	return link
}

func (m *Memo) Print() {
	fmt.Print(m.Title)
	fmt.Println(m.Content)
}

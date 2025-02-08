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

// SearchKey returns a key to search for a link in the memo content.
func (m *Memo) SearchKey() string {
	return filepath.Base(m.Filepath) + "#" + markdown.Text2tag(m.Title)
}

// Link returns a relative path to the memo file with the title as a tag.
func (m *Memo) Link() string {
	link := ".." + string(os.PathSeparator) + m.Filepath + "#" + markdown.Text2tag(m.Title)
	return link
}

// Print prints the memo title and content.
func (m *Memo) Print() {
	fmt.Print(m.Title)
	fmt.Println(m.Content)
}

package models

import "github.com/hirotoni/memo/markdown"

type Template struct {
	Headings []markdown.Heading
}

func NewTemplate(headings []markdown.Heading) Template {
	return Template{Headings: headings}
}

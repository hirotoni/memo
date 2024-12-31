package models

type Template struct {
	Headings []Heading
}

func NewTemplate(headings []Heading) Template {
	return Template{Headings: headings}
}

type Heading struct {
	Level int
	Text  string
}

func NewHeading(level int, text string) Heading {
	return Heading{Level: level, Text: text}
}

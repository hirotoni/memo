package models

type Template struct {
	Headings []Heading
}

type Heading struct {
	Level int
	Text  string
}

func NewHeading(level int, text string) Heading {
	return Heading{Level: level, Text: text}
}

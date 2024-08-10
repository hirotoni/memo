package models

type Template struct {
	Headings []Heading
}

type Heading struct {
	Level int
	Text  string
}

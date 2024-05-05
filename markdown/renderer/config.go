package renderer

import "github.com/yuin/goldmark/renderer"

/* this file does nothing, but keep it for potential use in the future. */

const (
	optSampleOption  renderer.OptionName = "SampleOption"
	optSampleOption2 renderer.OptionName = "SampleOption2"
)

type MarkdownRendererConfig struct {
	sampleOption  bool
	sampleOption2 bool
}

func NewMarkdownRendererConfig() *MarkdownRendererConfig {
	return &MarkdownRendererConfig{
		sampleOption:  true,
		sampleOption2: true,
	}
}

func (c *MarkdownRendererConfig) SetOption(name renderer.OptionName, value interface{}) {
	switch name {
	case optSampleOption:
		c.sampleOption = value.(bool)
	case optSampleOption2:
		c.sampleOption2 = value.(bool)
	}
}

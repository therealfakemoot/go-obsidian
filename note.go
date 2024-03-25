package obsidian

import (
// "go.abhg.dev/goldmark/frontmatter"
)

type Note struct {
	Name, Path string
	Tags       []Tag
	Properties any
}

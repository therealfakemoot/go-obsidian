package obsidian

import (
	"time"
	// "go.abhg.dev/goldmark/frontmatter"
)

type Note struct {
	Name, Path string
	Tags       []string
	Aliases    []string
	CSSClasses []string
	Date       time.Time
	Properties any
}

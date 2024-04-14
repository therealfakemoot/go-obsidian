package obsidian

import (
	"time"

	"github.com/goccy/go-graphviz/cgraph"
)

type Note struct {
	*cgraph.Node
	Name, Path string
	Names      []string
	Tags       []string
	Aliases    []string
	CSSClasses []string
	Date       time.Time
	Links      []Link
	Properties any
}

type Link struct {
	*cgraph.Edge
}

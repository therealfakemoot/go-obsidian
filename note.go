package obsidian

import (
	"time"
)

type Note struct {
	Tags                        []Tag
	Created, Modified, Accessed time.Time
	Properties                  any
}

package obsidian

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/frontmatter"
)

type Vault struct {
	Notes map[string]Note
	gm    goldmark.Markdown
}

func (v *Vault) walk(path string, d fs.DirEntry, err error) error {
	if strings.HasSuffix(path, ".md") && !d.IsDir() {
		var n Note
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("could not build absolute path for %q: %w", path, err)
		}
		n.Path = absPath
		filename := filepath.Base(path)
		if err != nil {
			return fmt.Errorf("could not split file path: %w", err)
		}
		n.Name = filename[:len(filename)-3]

		v.Notes[path] = n
	}

	return nil
}

func NewVault(root fs.FS) (*Vault, error) {
	v := &Vault{}
	v.Notes = make(map[string]Note)
	v.gm = goldmark.New(
		goldmark.WithExtensions(
			// ...
			&frontmatter.Extender{},
		),
	)

	err := fs.WalkDir(root, ".", v.walk)
	if err != nil {
		return v, fmt.Errorf("error walking vault: %w", err)
	}
	return v, nil
}

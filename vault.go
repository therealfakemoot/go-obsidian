package obsidian

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	// "github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

type Vault struct {
	Root  string
	Notes map[string]Note
	gm    goldmark.Markdown
}

func (v *Vault) walk(path string, d fs.DirEntry, err error) error {
	if strings.HasSuffix(path, ".md") && !d.IsDir() {
		var n Note
		absPath := filepath.Join(v.Root, path)
		if err != nil {
			return fmt.Errorf("could not build absolute path for %q: %w", path, err)
		}
		n.Path = absPath
		filename := filepath.Base(path)
		if err != nil {
			return fmt.Errorf("could not split file path: %w", err)
		}
		n.Name = filename[:len(filename)-3]

		f, err := os.Open(absPath)
		if err != nil {
			return fmt.Errorf("could not open note %q: %w", path, err)
		}
		var b bytes.Buffer
		io.Copy(&b, f)

		ctx := parser.NewContext()
		v.gm.Convert(b.Bytes(), io.Discard, parser.WithContext(ctx))

		raw := frontmatter.Get(ctx)

		var meta NoteMeta
		if err := raw.Decode(&meta); err != nil {
			return fmt.Errorf("couldn't decode frontmatter for %q: %w", path, err)
		}

		n.Tags = meta.Tags
		n.Aliases = meta.Aliases
		n.CSSClasses = meta.CSSClasses

		stamp, err := time.Parse("2006-01-02", meta.Date)
		if err != nil {
			return fmt.Errorf("could not parse Date field %q: %w", meta.Date, err)
		}

		n.Date = stamp
		v.Notes[path] = n
	}

	return nil
}

func NewVault(root string) (*Vault, error) {
	v := &Vault{}
	v.Root = root

	v.Notes = make(map[string]Note)
	v.gm = goldmark.New(
		goldmark.WithExtensions(&frontmatter.Extender{
			Mode: frontmatter.SetMetadata,
		}),
	)

	absRoot, err := filepath.Abs(v.Root)
	if err != nil {
		log.Fatalf("couldn't transform absolute path: %s\n", err)
	}

	rootFS := os.DirFS(absRoot)
	err = fs.WalkDir(rootFS, ".", v.walk)
	if err != nil {
		return v, fmt.Errorf("error walking vault: %w", err)
	}
	return v, nil
}

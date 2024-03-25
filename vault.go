package obsidian

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type Vault struct {
	Notes map[string]Note
}

func (v *Vault) walk(path string, d fs.DirEntry, err error) error {
	log.Printf("walking path %q\n", path)
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
		log.Printf("%#+v\n", n)
	}

	return nil
}

func NewVault(root fs.FS) (*Vault, error) {
	v := &Vault{}
	v.Notes = make(map[string]Note)

	err := fs.WalkDir(root, ".", v.walk)
	if err != nil {
		return v, fmt.Errorf("error walking vault: %w", err)
	}
	return v, nil
}

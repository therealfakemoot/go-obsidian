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

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/wikilink"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Vault struct {
	*cgraph.Graph
	Root   string
	Notes  map[string]Note
	gm     goldmark.Markdown
	Logger *zap.Logger
}

func (v *Vault) walk(path string, d fs.DirEntry, err error) error {
	v.Logger.Info("walking",
		zap.String("path", path),
	)

	if d.IsDir() && d.Name() == ".git" {
		return filepath.SkipDir
	}

	if strings.HasSuffix(path, ".md") && !d.IsDir() {
		v.Logger.Info("found note",
			zap.String("filename", path),
		)

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
		node, err := v.Graph.CreateNode(n.Name)
		if err != nil {
			return fmt.Errorf("couldn't create graph node for %q: %w", path, err)
		}

		n.Node = node
		n.Names = append(n.Names, n.Name)

		f, err := os.Open(absPath)
		if err != nil {
			return fmt.Errorf("could not open note %q: %w", path, err)
		}
		var b bytes.Buffer
		io.Copy(&b, f)

		ctx := parser.NewContext()
		links := LinkCollector{
			Links:  make([]*wikilink.Node, 0),
			Logger: v.Logger.Named("link-collector").With(zap.String("name", n.Name)),
		}
		gm := goldmark.New(
			goldmark.WithExtensions(
				&frontmatter.Extender{
					Mode: frontmatter.SetMetadata,
				},
				&wikilink.Extender{
					Resolver: &links,
				},
			),
		)
		v.Logger.Info("links in note",
			zap.String("name", n.Name),
			zap.Any("links", links.Links),
		)

		gm.Convert(b.Bytes(), io.Discard, parser.WithContext(ctx))

		raw := frontmatter.Get(ctx)
		if raw != nil {

			var meta NoteMeta
			if err := raw.Decode(&meta); err != nil {
				return fmt.Errorf("couldn't decode frontmatter for %q: %w", path, err)
			}

			n.Tags = meta.Tags
			n.Aliases = meta.Aliases
			n.CSSClasses = meta.CSSClasses

			for _, alias := range n.Aliases {
				n.Names = append(n.Names, alias)
			}

			var stamp time.Time
			if meta.Date != "" {
				stamp, err = time.Parse("2006-01-02", meta.Date)
				if err != nil {
					return fmt.Errorf("could not parse Date field %q: %w", meta.Date, err)
				}
			}

			n.Date = stamp
		}
		v.Notes[n.Name] = n
	}

	return nil
}

func NewVault(root string, g *cgraph.Graph) (*Vault, error) {
	prodConfig := zap.NewProductionConfig()
	prodConfig.Encoding = "console"
	prodConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	prodConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	logger, _ := prodConfig.Build()

	logger.Info("entering NewVault()")

	v := &Vault{Graph: g}
	v.Logger = logger
	v.Root = root

	v.Notes = make(map[string]Note)
	v.Logger.Info("building goldmark")
	v.Logger.Info("building absRoot")
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

type LinkCollector struct {
	Links  []*wikilink.Node
	Logger *zap.Logger
}

func (lc *LinkCollector) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	lc.Logger.Info("received node",
		zap.Any("node", n),
	)
	lc.Links = append(lc.Links, n)

	return n.Target, nil
}

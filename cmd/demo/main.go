package main

import (
	"flag"
	"log"

	"github.com/goccy/go-graphviz"

	"github.com/therealfakemoot/go-obsidian"
)

func main() {
	var (
		version bool
		root    string
	)

	flag.StringVar(&root, "root", ".", "Directory to use as vault root")
	flag.BoolVar(&version, "version", false, "Print version info")

	flag.Parse()

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatalf("couldn't create graphviz graph: %s\n", err)
	}

	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatalf("error closing graph: %s\n", err)
		}
	}()

	vault, err := obsidian.NewVault(root, graph)
	if err != nil {
		log.Fatalf("couldn't walk vault: %s\n", err)
	}

	defer vault.Logger.Sync()

	/*
		for k, v := range vault.Notes {
			vault.Logger.Info("indexed note",
				zap.String("path", k),
				zap.Any("note", v),
			)
		}
	*/
}

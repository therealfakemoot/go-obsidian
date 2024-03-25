package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

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

	absRoot, err := filepath.Abs(root)
	if err != nil {
		log.Fatalf("couldn't transform absolute path: %s\n", err)
	}

	rootFS := os.DirFS(absRoot)

	v, err := obsidian.NewVault(rootFS)
	if err != nil {
		log.Fatalf("couldn't walk vault: %s\n", err)
	}

	log.Printf("%#+v\n", v)
}

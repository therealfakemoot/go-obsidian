package main

import (
	"flag"
	"log"

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

	_, err := obsidian.NewVault(root)
	if err != nil {
		log.Fatalf("couldn't walk vault: %s\n", err)
	}
}

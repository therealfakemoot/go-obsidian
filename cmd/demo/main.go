package main

import (
	"flag"
	"log"

	"go.uber.org/zap"

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

	vault, err := obsidian.NewVault(root)
	if err != nil {
		log.Fatalf("couldn't walk vault: %s\n", err)
	}

	defer vault.Logger.Sync()

	for k, v := range vault.Notes {
		vault.Logger.Info("indexed note",
			zap.String("path", k),
			zap.Any("note", v),
		)
	}
}

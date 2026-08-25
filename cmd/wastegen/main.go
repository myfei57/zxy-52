package main

import (
	"flag"
	"log"

	"wastegen/internal/console"
	"wastegen/internal/store"
)

func main() {
	root := flag.String("root", "./data", "file persistence root directory")
	addr := flag.String("addr", "127.0.0.1:8901", "console listen address")
	flag.Parse()
	st, err := store.Open(*root)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	srv := console.NewServer(st)
	if err := srv.Start(*addr); err != nil {
		log.Fatalf("start console: %v", err)
	}
}

package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/axeprpr/n2nGUI/internal/app"
)

//go:embed frontend/*
var staticFiles embed.FS

func main() {
	listen := flag.String("listen", "127.0.0.1:8787", "HTTP listen address")
	base := flag.String("base", ".", "base directory containing n2n binaries and config")
	flag.Parse()

	baseDir, err := filepath.Abs(*base)
	if err != nil {
		log.Fatalf("resolve base directory: %v", err)
	}

	if _, err := os.Stat(filepath.Join(baseDir, "n2n")); err != nil {
		log.Printf("warning: n2n directory not found under %s", baseDir)
	}

	server := app.NewServer(baseDir)
	handler, err := server.Handler(staticFiles)
	if err != nil {
		log.Fatalf("build HTTP handler: %v", err)
	}

	httpServer := &http.Server{
		Addr:              *listen,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("n2nGUI listening on http://%s (base=%s)", *listen, baseDir)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %v", err)
	}
}

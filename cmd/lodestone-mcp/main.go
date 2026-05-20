package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmt-labs/lodestone/internal/lodestone/mcp"
)

var (
	version = "dev"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Fprintf(os.Stdout, "lodestone-mcp %s\n", version)
		return
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cwd:", err)
		os.Exit(1)
	}
	storeRoot := filepath.Join(root, ".lodestone")
	if err := os.MkdirAll(storeRoot, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir .lodestone:", err)
		os.Exit(1)
	}

	reg := mcp.NewToolRegistry()
	mcp.RegisterBuiltins(reg, mcp.DefaultDeps(root, storeRoot))

	srv := mcp.NewServer("lodestone-mcp", version, reg)
	if err := srv.Serve(context.Background(), os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "serve:", err)
		os.Exit(1)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/simulacre7/todo-rank/internal/scan"
)

const (
	defaultRoot     = "."
	defaultIgnore   = ".git,node_modules,dist"
	defaultFormat   = "text"
	defaultMinScore = 0
	defaultTags     = "TODO,FIXME,@next"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]

	// Handle "scan" subcommand: strip it if present
	if len(args) > 0 && args[0] == "scan" {
		args = args[1:]
	}

	fs := flag.NewFlagSet("todo-rank", flag.ContinueOnError)

	root := fs.String("root", defaultRoot, "scan start directory")
	ignore := fs.String("ignore", defaultIgnore, "directories to ignore (comma-separated)")
	format := fs.String("format", defaultFormat, "output format (text|md)")
	out := fs.String("out", "", "output file path (default: stdout)")
	minScore := fs.Int("min-score", defaultMinScore, "minimum score filter")
	tags := fs.String("tags", defaultTags, "tags to scan (comma-separated)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate --format
	if *format != "text" && *format != "md" {
		return fmt.Errorf("invalid format: %q (must be text or md)", *format)
	}

	// Validate --root exists
	info, err := os.Stat(*root)
	if err != nil {
		return fmt.Errorf("root directory error: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("root is not a directory: %s", *root)
	}

	// Validate --out directory exists (if specified)
	if *out != "" {
		dir := dirOf(*out)
		if dir != "" {
			info, err := os.Stat(dir)
			if err != nil {
				return fmt.Errorf("output directory error: %w", err)
			}
			if !info.IsDir() {
				return fmt.Errorf("output path parent is not a directory: %s", dir)
			}
		}
	}

	opts := scan.ScanOptions{
		Root:     *root,
		Ignore:   splitCSV(*ignore),
		Format:   *format,
		OutPath:  *out,
		MinScore: *minScore,
		Tags:     splitCSV(*tags),
	}

	return scan.Run(opts)
}

// splitCSV splits a comma-separated string into a slice.
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// dirOf returns the directory portion of a path.
func dirOf(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == os.PathSeparator {
			return path[:i]
		}
	}
	return ""
}

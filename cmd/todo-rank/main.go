package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/simulacre7/todo-rank/internal/initialize"
	"github.com/simulacre7/todo-rank/internal/render"
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
	if len(os.Args) < 2 {
		return runScan(os.Args[1:])
	}

	switch os.Args[1] {
	case "init":
		return runInit(os.Args[2:])
	case "scan":
		return runScan(os.Args[2:])
	default:
		// No subcommand, treat as scan with flags
		return runScan(os.Args[1:])
	}
}

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	force := fs.Bool("force", false, "overwrite existing files")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: todo-rank init [--force] <agent-type>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Agent types:")
		fmt.Fprintln(os.Stderr, "  llms     Generate llms.txt")
		fmt.Fprintln(os.Stderr, "  claude   Generate CLAUDE.md")
		fmt.Fprintln(os.Stderr, "  cursor   Generate .cursorrules")
		fmt.Fprintln(os.Stderr, "  all      Generate all files")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		fs.Usage()
		return fmt.Errorf("agent type required")
	}

	agentType := fs.Arg(0)

	results, err := initialize.Run(agentType, *force)
	if err != nil {
		return err
	}

	for _, r := range results {
		if r.Error != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %v\n", r.File, r.Error)
		} else if r.Skipped {
			fmt.Printf("skipped: %s (already exists, use --force to overwrite)\n", r.File)
		} else if r.Created {
			fmt.Printf("created: %s\n", r.File)
		}
	}

	return nil
}

func runScan(args []string) error {
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

	results, err := scan.Run(opts)
	if err != nil {
		return err
	}

	// Determine output writer
	var w *os.File
	if *out == "" {
		w = os.Stdout
	} else {
		f, err := os.Create(*out)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return render.Render(w, results, *format)
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

package scan

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/todo-rank/internal/parse"
	"github.com/todo-rank/internal/score"
)

// Run executes the scan with the given options.
// It walks the directory tree, parses TODO comments, scores them,
// and returns the collected results filtered by MinScore.
func Run(opts ScanOptions) ([]score.ScoredTodo, error) {
	var results []score.ScoredTodo

	ignoreSet := makeIgnoreSet(opts.Ignore)

	err := filepath.WalkDir(opts.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored directories
		if d.IsDir() {
			if shouldIgnore(d.Name(), ignoreSet) {
				return filepath.SkipDir
			}
			return nil
		}

		// Process files
		todos, err := scanFile(path, opts.Root)
		if err != nil {
			// Skip files that cannot be read
			return nil
		}

		for _, item := range todos {
			scored := score.ScoreTodo(item)
			if scored.Score >= opts.MinScore {
				results = append(results, scored)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// makeIgnoreSet creates a set of directory names to ignore.
func makeIgnoreSet(ignore []string) map[string]struct{} {
	set := make(map[string]struct{}, len(ignore))
	for _, name := range ignore {
		set[name] = struct{}{}
	}
	return set
}

// shouldIgnore checks if a directory name should be skipped.
func shouldIgnore(name string, ignoreSet map[string]struct{}) bool {
	_, ok := ignoreSet[name]
	return ok
}

// scanFile reads a file line-by-line and extracts TODO items.
func scanFile(path string, root string) ([]score.TodoItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var items []score.TodoItem
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		parsed, ok := parse.ParseLine(line)
		if !ok {
			continue
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			relPath = path
		}

		item := score.TodoItem{
			Tag:      parsed.Tag,
			Priority: parsed.Priority,
			Message:  parsed.Message,
			Path:     relPath,
			Line:     lineNum,
		}
		items = append(items, item)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

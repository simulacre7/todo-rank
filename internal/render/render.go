package render

import (
	"fmt"
	"io"
	"sort"

	"github.com/simulacre7/todo-rank/internal/score"
)

type ScoredTodo = score.ScoredTodo

var levelOrder = []string{"P0", "P1", "P2", "P3"}

var levelLabels = map[string]string{
	"P0": "P0 (Now)",
	"P1": "P1 (Soon)",
	"P2": "P2 (Later)",
	"P3": "P3 (Cleanup)",
}

func sortItems(items []ScoredTodo) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Score != items[j].Score {
			return items[i].Score > items[j].Score
		}
		if items[i].Path != items[j].Path {
			return items[i].Path < items[j].Path
		}
		return items[i].Line < items[j].Line
	})
}

func groupByLevel(items []ScoredTodo) map[string][]ScoredTodo {
	groups := make(map[string][]ScoredTodo)
	for _, item := range items {
		groups[item.Level] = append(groups[item.Level], item)
	}
	return groups
}

func formatTag(tag string, priority *int) string {
	if priority == nil {
		return tag
	}
	return fmt.Sprintf("%s[P%d]", tag, *priority)
}

func RenderText(w io.Writer, items []ScoredTodo) error {
	sortItems(items)
	groups := groupByLevel(items)

	first := true
	for _, level := range levelOrder {
		group, ok := groups[level]
		if !ok || len(group) == 0 {
			continue
		}

		if !first {
			fmt.Fprintln(w)
		}
		first = false

		fmt.Fprintln(w, levelLabels[level])

		for _, item := range group {
			fmt.Fprintf(w, "(%d) %s:%d\n", item.Score, item.Path, item.Line)
			fmt.Fprintf(w, "  %s: %s\n", formatTag(item.Tag, item.Priority), item.Message)
		}
	}

	return nil
}

func RenderMarkdown(w io.Writer, items []ScoredTodo) error {
	sortItems(items)
	groups := groupByLevel(items)

	first := true
	for _, level := range levelOrder {
		group, ok := groups[level]
		if !ok || len(group) == 0 {
			continue
		}

		if !first {
			fmt.Fprintln(w)
		}
		first = false

		fmt.Fprintf(w, "## %s\n", levelLabels[level])

		for _, item := range group {
			fmt.Fprintf(w, "- [ ] %s:%d  \n", item.Path, item.Line)
			fmt.Fprintf(w, "  %s: %s\n", formatTag(item.Tag, item.Priority), item.Message)
		}
	}

	return nil
}

func Render(w io.Writer, items []ScoredTodo, format string) error {
	switch format {
	case "md":
		return RenderMarkdown(w, items)
	default:
		return RenderText(w, items)
	}
}

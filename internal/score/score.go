package score

import (
	"path/filepath"
	"strings"
)

type TodoItem struct {
	Tag      string
	Priority *int
	Message  string
	Path     string
	Line     int
}

type ScoredTodo struct {
	TodoItem
	Score int
	Level string
}

func CalcTagScore(tag string) int {
	switch tag {
	case "FIXME":
		return 100
	case "TODO":
		return 50
	case "@next":
		return 30
	default:
		return 0
	}
}

func CalcPriorityScore(p *int) int {
	if p == nil {
		return 0
	}
	switch *p {
	case 0:
		return 100
	case 1:
		return 70
	case 2:
		return 40
	case 3:
		return 10
	default:
		return 0
	}
}

func CalcPathBonus(path string) int {
	if strings.Contains(path, "cmd/") {
		return 20
	}
	if filepath.Base(path) == "main.go" {
		return 20
	}
	return 0
}

func CalcTestPenalty(path string) int {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if strings.HasSuffix(name, "_test") {
		return -20
	}
	return 0
}

func CalcLevel(score int) string {
	if score >= 120 {
		return "P0"
	}
	if score >= 80 {
		return "P1"
	}
	if score >= 40 {
		return "P2"
	}
	return "P3"
}

func ScoreTodo(item TodoItem) ScoredTodo {
	score := CalcTagScore(item.Tag) +
		CalcPriorityScore(item.Priority) +
		CalcPathBonus(item.Path) +
		CalcTestPenalty(item.Path)

	return ScoredTodo{
		TodoItem: item,
		Score:    score,
		Level:    CalcLevel(score),
	}
}

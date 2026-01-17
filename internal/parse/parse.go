package parse

import (
	"regexp"
	"strconv"
)

type ParsedTodo struct {
	Tag      string
	Priority *int
	Message  string
}

var todoRegex = regexp.MustCompile(`^\s*(?://|#|/\*)?\s*(TODO|FIXME|@next)(?:\[(P[0-3])\])?\s*[:\-]\s*(.+?)\s*(?:\*/)?$`)

func ParseLine(line string) (todo ParsedTodo, ok bool) {
	matches := todoRegex.FindStringSubmatch(line)
	if matches == nil {
		return ParsedTodo{}, false
	}

	todo.Tag = matches[1]
	todo.Message = matches[3]

	if matches[2] != "" {
		p, _ := strconv.Atoi(string(matches[2][1]))
		todo.Priority = &p
	}

	return todo, true
}

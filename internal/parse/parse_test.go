package parse

import (
	"testing"
)

func intPtr(i int) *int {
	return &i
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantOk   bool
		wantTodo ParsedTodo
	}{
		{
			name:   "simple TODO with colon",
			input:  "// TODO: something",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:     "TODO",
				Message: "something",
			},
		},
		{
			name:   "TODO with priority",
			input:  "// TODO[P1]: improve error handling",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:      "TODO",
				Priority: intPtr(1),
				Message:  "improve error handling",
			},
		},
		{
			name:   "FIXME with priority",
			input:  "// FIXME[P0]: data race here",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:      "FIXME",
				Priority: intPtr(0),
				Message:  "data race here",
			},
		},
		{
			name:   "@next without priority",
			input:  "// @next: refactor naming",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:     "@next",
				Message: "refactor naming",
			},
		},
		{
			name:   "@next with priority",
			input:  "// @next[P2]: cleanup naming",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:      "@next",
				Priority: intPtr(2),
				Message:  "cleanup naming",
			},
		},
		{
			name:   "hash comment marker",
			input:  "# TODO: something",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:     "TODO",
				Message: "something",
			},
		},
		{
			name:   "block comment marker",
			input:  "/* TODO: something */",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:     "TODO",
				Message: "something",
			},
		},
		{
			name:   "no comment marker",
			input:  "TODO: something",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:     "TODO",
				Message: "something",
			},
		},
		{
			name:   "dash separator",
			input:  "// TODO - alternative separator is allowed",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:     "TODO",
				Message: "alternative separator is allowed",
			},
		},
		{
			name:   "P3 priority",
			input:  "// TODO[P3]: low priority task",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:      "TODO",
				Priority: intPtr(3),
				Message:  "low priority task",
			},
		},
		{
			name:   "no separator - should fail",
			input:  "TODO something",
			wantOk: false,
		},
		{
			name:   "wrong tag - should fail",
			input:  "TODOS: something",
			wantOk: false,
		},
		{
			name:   "space before priority - should fail",
			input:  "TODO [P1]: something",
			wantOk: false,
		},
		{
			name:   "empty message - should fail",
			input:  "// TODO:",
			wantOk: false,
		},
		{
			name:   "lowercase tag - should fail",
			input:  "// todo: something",
			wantOk: false,
		},
		{
			name:   "invalid priority - should fail",
			input:  "// TODO[P4]: something",
			wantOk: false,
		},
		{
			name:   "leading whitespace",
			input:  "    // TODO: something",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:     "TODO",
				Message: "something",
			},
		},
		{
			name:   "trailing whitespace in message trimmed",
			input:  "// TODO: something   ",
			wantOk: true,
			wantTodo: ParsedTodo{
				Tag:     "TODO",
				Message: "something",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseLine(tt.input)
			if ok != tt.wantOk {
				t.Errorf("ParseLine() ok = %v, want %v", ok, tt.wantOk)
				return
			}
			if !tt.wantOk {
				return
			}
			if got.Tag != tt.wantTodo.Tag {
				t.Errorf("Tag = %q, want %q", got.Tag, tt.wantTodo.Tag)
			}
			if got.Message != tt.wantTodo.Message {
				t.Errorf("Message = %q, want %q", got.Message, tt.wantTodo.Message)
			}
			if (got.Priority == nil) != (tt.wantTodo.Priority == nil) {
				t.Errorf("Priority nil mismatch: got %v, want %v", got.Priority, tt.wantTodo.Priority)
			} else if got.Priority != nil && *got.Priority != *tt.wantTodo.Priority {
				t.Errorf("Priority = %d, want %d", *got.Priority, *tt.wantTodo.Priority)
			}
		})
	}
}

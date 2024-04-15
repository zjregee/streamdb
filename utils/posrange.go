package utils

import "fmt"

type Pos int

type PositionRange struct {
	Start Pos
	End   Pos
}

func (p PositionRange) StartPosInput(query string, lineOffset int) string {
	if query == "" {
		return "unknown position"
	}
	pos := int(p.Start)
	if pos < 0 || pos > len(query) {
		return "invalid position"
	}
	lastLineBreak := -1
	line := lineOffset + 1
	for i, c := range query[:pos] {
		if c == '\n' {
			lastLineBreak = i
			line += 1
		}
	}
	col := pos - lastLineBreak
	return fmt.Sprintf("%d:%d", line, col)
}

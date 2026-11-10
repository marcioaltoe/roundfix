package spec

import (
	"fmt"
	"strings"
)

// QAResultsRow is one data row with its line in the whole QA Report.
type QAResultsRow struct {
	Line  int
	Cells []string
}

// QAResultsTable holds a header-led table in the Results section.
type QAResultsTable struct {
	HeaderLine int
	Header     []string
	Rows       []QAResultsRow
}

// QAResults holds the Results tables and the first table line outside them.
type QAResults struct {
	Tables []QAResultsTable
	Break  *QAResultsBreak
}

// ReadQAResults reads the first Results section, outside frontmatter and
// fences, stopping at its next level-one or level-two heading or an orphan row.
func ReadQAResults(path string, content []byte) QAResults {
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	var results QAResults
	index := 0
	if lines[0] == "---" {
		for index = 1; index < len(lines); index++ {
			if lines[index] == "---" {
				index++
				break
			}
		}
	}
	inResults, fence := false, ""
	lastTableEnd := -1
	for index < len(lines) {
		text := strings.TrimSpace(lines[index])
		if marker, ok := markdownFenceMarker(text); ok {
			if fence == "" {
				fence = marker
			} else if fence == marker {
				fence = ""
			}
			index++
			continue
		}
		if fence != "" {
			index++
			continue
		}
		if depth, heading, ok := markdownHeading(text); ok {
			if inResults && depth <= 2 {
				return results
			}
			if depth == 2 && strings.EqualFold(heading, "Results") {
				inResults = true
			}
			index++
			continue
		}
		cells := markdownTableCells(text)
		if !inResults || len(cells) == 0 {
			index++
			continue
		}
		if index+1 < len(lines) && markdownTableSeparator(markdownTableCells(lines[index+1])) {
			table := QAResultsTable{HeaderLine: index + 1, Header: cells}
			index += 2
			for index < len(lines) {
				cells := markdownTableCells(lines[index])
				if len(cells) == 0 {
					break
				}
				if !markdownTableSeparator(cells) {
					table.Rows = append(table.Rows, QAResultsRow{Line: index + 1, Cells: cells})
				}
				index++
			}
			results.Tables = append(results.Tables, table)
			lastTableEnd = index - 1
			continue
		}
		brk := &QAResultsBreak{Path: path, Line: index + 1}
		if len(results.Tables) > 0 {
			table := results.Tables[len(results.Tables)-1]
			brk.TableLine, brk.LastRowLine = table.HeaderLine, table.HeaderLine
			if len(table.Rows) > 0 {
				brk.LastRowLine = table.Rows[len(table.Rows)-1].Line
			}
			blank := lastTableEnd+1 < index
			for _, line := range lines[lastTableEnd+1 : index] {
				if strings.TrimSpace(line) != "" {
					blank = false
					break
				}
			}
			if blank {
				brk.BlankFrom, brk.BlankTo = lastTableEnd+2, index
			}
		}
		results.Break = brk
		return results
	}
	return results
}

// StatusTables selects the tables whose own header declares a Status column.
func (results QAResults) StatusTables() []QAResultsTable {
	var tables []QAResultsTable
	for _, table := range results.Tables {
		for _, cell := range table.Header {
			if strings.EqualFold(strings.TrimSpace(cell), "Status") {
				tables = append(tables, table)
				break
			}
		}
	}
	return tables
}

// QAResultsBreak names the first Results table line no header leads and the
// last table above it. BlankFrom and BlankTo are zero for a non-blank gap.
type QAResultsBreak struct {
	Path                                             string
	Line, TableLine, LastRowLine, BlankFrom, BlankTo int
}

func (b QAResultsBreak) Detail() string {
	if b.TableLine == 0 {
		return fmt.Sprintf("Results row at line %d is outside every Results table", b.Line)
	}
	return fmt.Sprintf("Results row at line %d is outside the Results table that starts at line %d", b.Line, b.TableLine)
}

func (b QAResultsBreak) Repair() string {
	switch {
	case b.TableLine == 0:
		return fmt.Sprintf("add a separator line directly below the table header at %s:%d", b.Path, b.Line)
	case b.BlankFrom > 0 && b.BlankFrom == b.BlankTo:
		return fmt.Sprintf("remove the blank line at %s:%d so the row continues the table", b.Path, b.BlankFrom)
	case b.BlankFrom > 0:
		return fmt.Sprintf("remove the blank lines at %s:%d-%d so the row continues the table", b.Path, b.BlankFrom, b.BlankTo)
	default:
		return fmt.Sprintf("move the row directly below the last row of the table, at %s:%d", b.Path, b.LastRowLine)
	}
}

func (b QAResultsBreak) Error() string {
	return fmt.Sprintf("QA Report %s: %s; %s", b.Path, b.Detail(), b.Repair())
}

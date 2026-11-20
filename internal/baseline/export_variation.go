package baseline

import "strings"

// ExportVariationDepartures reads only the repository's lines outside a valid
// standard block. A comment supplies a reason until a blank or marker line.
func ExportVariationDepartures(current []byte, facts ExportListFacts) []DocumentationDeparture {
	var result []DocumentationDeparture
	lines := strings.Split(string(current), "\n")
	begin, end, begins, ends := -1, -1, 0, 0
	for i, line := range lines {
		switch strings.TrimSpace(line) {
		case ExportListBeginMarker:
			begin = i
			begins++
		case ExportListEndMarker:
			end = i
			ends++
		}
	}
	reason := false
	for i, line := range lines {
		entry := strings.TrimSpace(line)
		if begins == 1 && ends == 1 && begin < end && i >= begin && i <= end {
			reason = false
			continue
		}
		switch {
		case entry == "", entry == ExportListBeginMarker, entry == ExportListEndMarker:
			reason = false
		case strings.HasPrefix(entry, "#"):
			reason = true
		default:
			if !strings.HasPrefix(entry, "!") && (exportProcess(entry, facts) || exportContainsProcess(entry, facts)) {
				result = append(result, DocumentationDeparture{Code: "DR-EXPORT-PROCESS-PATH", Path: entry, Line: i + 1})
			}
			if !reason {
				result = append(result, DocumentationDeparture{Code: "DR-EXPORT-REASON-MISSING", Path: entry, Line: i + 1})
			}
		}
	}
	return result
}

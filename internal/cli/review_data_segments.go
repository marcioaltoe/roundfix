package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type reviewDataSegment struct {
	Paths                 []string `json:"paths"`
	DiffBytes             int      `json:"diffBytes"`
	EstimatedPromptTokens int      `json:"estimatedPromptTokens,omitempty"`
}

// reviewDataFit reports whether a prompt of promptBytes carrying diffBytes of diff fits.
type reviewDataFit func(promptBytes, diffBytes int) bool

func planReviewDataSegments(diff, from, head string, fit reviewDataFit) ([]reviewDataSegment, []string, []string, error) {
	var segments []reviewDataSegment
	var diffs, prompts []string
	// Every admitted index and total is one digit. buildReviewPrompt adds one
	// newline for an empty diff; account for the actual last byte when measuring.
	overhead := len(reviewDataSegmentHeader(1, reviewSegmentLimit, from, head)) + len(buildReviewPrompt(from, head, "")) - 1
	fits := func(bytes int, newline bool) bool {
		promptBytes := overhead + bytes
		if !newline {
			promptBytes++
		}
		return fit(promptBytes, bytes)
	}
	var current strings.Builder
	var paths []string
	lastPath := ""
	flush := func() {
		if current.Len() == 0 {
			return
		}
		segments = append(segments, reviewDataSegment{Paths: paths, DiffBytes: current.Len()})
		diffs = append(diffs, current.String())
		current = strings.Builder{}
		paths = nil
		lastPath = ""
	}
	appendPiece := func(path, header, piece string) {
		addition := piece
		// Adjacent hunks of one file share a file header in the packed segment.
		// This also keeps all its hunk ranges in the finding anchor index.
		if lastPath == path {
			addition = strings.TrimPrefix(piece, header)
		}
		if current.Len() > 0 && !fits(current.Len()+len(addition), strings.HasSuffix(addition, "\n")) {
			flush()
			addition = piece
		}
		if lastPath != path {
			paths = append(paths, path)
		}
		current.WriteString(addition)
		lastPath = path
	}
	for _, section := range splitReviewDataSections(diff, "diff --git ") {
		if section == "" {
			continue
		}
		at := strings.Index(section, "\n@@ ")
		header := section
		if at >= 0 {
			header = section[:at+1]
		}
		path := ""
		for name := range indexReviewDiff(header) {
			path = name
		}
		if fits(len(section), strings.HasSuffix(section, "\n")) {
			appendPiece(path, header, section)
			continue
		}
		if at < 0 {
			return nil, nil, nil, fmt.Errorf("review data too large: line 0 of %s does not fit one data review segment", path)
		}
		for _, hunk := range splitReviewDataSections(section[at+1:], "@@ ") {
			err := cutReviewDataHunk(path, header, hunk, fits, func(piece string) { appendPiece(path, header, piece) })
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}
	flush()
	if len(segments) > reviewSegmentLimit {
		return nil, nil, nil, fmt.Errorf("review data too large: %d bytes of pre_pr_review.data_paths needs %d data review segments; a review allows at most %d", len(diff), len(segments), reviewSegmentLimit)
	}
	for i, piece := range diffs {
		prompts = append(prompts, reviewDataSegmentHeader(i+1, len(segments), from, head)+buildReviewPrompt(from, head, piece))
	}
	return segments, diffs, prompts, nil
}

func reviewDataSegmentHeader(index, total int, from, head string) string {
	return fmt.Sprintf("Data Review Segment %d of %d: pre_pr_review.data_paths changes %s..%s of the candidate; answer for this segment only.\n", index, total, from, head) +
		"These files are data the project writes by hand and declares in pre_pr_review.data_paths, such as rule catalogs. Check each changed value against the Spec contexts, the code that reads it and the rest of the file; report a wrong, missing, duplicated or contradictory value.\n"
}

// Keep slices of the original diff, including its line endings.
func splitReviewDataSections(diff, prefix string) []string {
	var sections []string
	for diff != "" {
		at := strings.Index(diff, "\n"+prefix)
		if at < 0 {
			sections = append(sections, diff)
			break
		}
		sections = append(sections, diff[:at+1])
		diff = diff[at+1:]
	}
	return sections
}

func reviewDataHunkRange(value string) (int, error) {
	start, count, hasCount := strings.Cut(value[1:], ",")
	line, err := strconv.Atoi(start)
	if err != nil {
		return 0, fmt.Errorf("parse data review hunk range: %w", err)
	}
	if hasCount && count == "0" {
		line++
	}
	return line, nil
}

func cutReviewDataHunk(path, header, hunk string, fits func(int, bool) bool, emit func(string)) error {
	first, body, _ := strings.Cut(hunk, "\n")
	fields := strings.Fields(first)
	if len(fields) < 4 {
		return fmt.Errorf("parse data review hunk: invalid header %q", first)
	}
	old, err := reviewDataHunkRange(fields[1])
	if err != nil {
		return err
	}
	newLine, err := reviewDataHunkRange(fields[2])
	if err != nil {
		return err
	}
	start, end := 0, 0
	oldStart, newStart := old, newLine
	oldCount, newCount := 0, 0
	emitPiece := func() {
		o, n := oldStart, newStart
		if oldCount == 0 {
			o--
		}
		if newCount == 0 {
			n--
		}
		emit(header + fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", o, oldCount, n, newCount) + body[start:end])
	}
	for end < len(body) {
		next := reviewDataLineEnd(body, end)
		// Git's no-newline marker belongs to the preceding content line.
		if strings.HasPrefix(body[next:], "\\ No newline at end of file") {
			next = reviewDataLineEnd(body, next)
		}
		oc, nc := 0, 0
		switch body[end] {
		case ' ':
			oc, nc = 1, 1
		case '-':
			oc = 1
		case '+':
			nc = 1
		}
		o, n := oldStart, newStart
		if oldCount+oc == 0 {
			o--
		}
		if newCount+nc == 0 {
			n--
		}
		bytes := len(header) + reviewDataHunkHeaderBytes(o, oldCount+oc, n, newCount+nc) + next - start
		if !fits(bytes, body[next-1] == '\n') {
			if end == start {
				return fmt.Errorf("review data too large: line %d of %s does not fit one data review segment", newLine, path)
			}
			emitPiece()
			start = end
			oldStart, newStart = old, newLine
			oldCount, newCount = 0, 0
			continue
		}
		oldCount += oc
		newCount += nc
		old += oc
		newLine += nc
		end = next
	}
	if end > start {
		emitPiece()
	}
	return nil
}

func reviewDataLineEnd(body string, start int) int {
	if at := strings.IndexByte(body[start:], '\n'); at >= 0 {
		return start + at + 1
	}
	return len(body)
}

// Only the four decimal fields vary; never rebuild the growing body to size a line.
func reviewDataHunkHeaderBytes(old, oldCount, newLine, newCount int) int {
	size := len("@@ -, +, @@\n")
	for _, value := range []int{old, oldCount, newLine, newCount} {
		size++
		for value >= 10 {
			size++
			value /= 10
		}
	}
	return size
}

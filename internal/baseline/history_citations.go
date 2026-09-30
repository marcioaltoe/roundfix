package baseline

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const (
	relocationCitationCode          = "baseline.history.citation"
	relocationCitationOmittedCode   = "baseline.history.citation.omitted"
	relocationCitationUnscannedCode = "baseline.history.citation.unscanned"
	relocationCitationFileLimit     = 200
	relocationCitationPerFileLimit  = 3
	relocationCitationMaxFileSize   = 4 * 1024 * 1024
	relocationCitationBinaryProbe   = 8000
)

type relocationCitation struct {
	line        int
	text        string
	before      string
	after       string
	destination string
}

type markdownDestination struct {
	raw   string
	start int
	end   int
}

type relocationResolution struct {
	existsBefore map[string]struct{}
	existsAfter  map[string]struct{}
	afterPath    map[string]string
	movedFiles   map[string]string
	movedUnits   map[string]string
}

type trackedReadResult uint8

const (
	trackedReadSkipped trackedReadResult = iota
	trackedReadScanned
	trackedReadUnscanned
)

// relocationCitationFindings returns one warning per tracked file whose
// citations resolve before the given History Relocations and not after them.
func relocationCitationFindings(
	ctx context.Context,
	root string,
	moves []HistoryMove,
	refused map[string]bool,
) ([]Finding, error) {
	if len(moves) == 0 {
		return nil, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	tracked, err := listTrackedPaths(ctx, root)
	if err != nil {
		return nil, err
	}
	resolution := buildRelocationResolution(tracked, moves, refused)
	directorySafety := make(map[string]bool)
	findings := make([]Finding, 0)
	unscanned := make([]string, 0)
	omitted := 0
	for _, relative := range tracked {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		content, result := readTrackedCitationFile(root, relative, directorySafety)
		switch result {
		case trackedReadSkipped:
			continue
		case trackedReadUnscanned:
			unscanned = append(unscanned, relative)
			continue
		}

		citations := citationsBrokenByRelocation(relative, content, resolution)
		if len(citations) == 0 {
			continue
		}
		if len(findings) >= relocationCitationFileLimit {
			omitted++
			continue
		}
		findings = append(findings, Finding{
			Code:    relocationCitationCode,
			Path:    relative,
			Message: relocationCitationMessage(citations),
		})
	}

	if omitted > 0 {
		findings = append(findings, Finding{
			Code:    relocationCitationOmittedCode,
			Path:    ".",
			Message: fmt.Sprintf("%d more tracked files cite paths this plan relocates and are not listed", omitted),
		})
	}
	if len(unscanned) > 0 {
		listed := unscanned
		if len(listed) > 3 {
			listed = listed[:3]
		}
		findings = append(findings, Finding{
			Code: relocationCitationUnscannedCode,
			Path: ".",
			Message: fmt.Sprintf(
				"%d tracked files were not scanned for citations (larger than 4 MiB or unreadable): %s",
				len(unscanned),
				strings.Join(listed, ", "),
			),
		})
	}
	return findings, nil
}

// listTrackedPaths returns safe paths from the Git index, deduplicated and sorted.
func listTrackedPaths(ctx context.Context, root string) ([]string, error) {
	output, err := (ExecGitRunner{}).RunGit(ctx, root, "ls-files", "-z", "--cached", "--full-name")
	if err != nil {
		return nil, fmt.Errorf("list tracked paths for relocation citations: %w", err)
	}
	seen := make(map[string]struct{})
	for _, relative := range strings.Split(output, "\x00") {
		if relative == "" || !repositoryPathIsSafe(relative) || containsControlCharacter(relative) {
			continue
		}
		seen[relative] = struct{}{}
	}
	paths := make([]string, 0, len(seen))
	for relative := range seen {
		paths = append(paths, relative)
	}
	sort.Strings(paths)
	return paths, nil
}

func buildRelocationResolution(tracked []string, moves []HistoryMove, refused map[string]bool) relocationResolution {
	filesBefore := citationStringSet(tracked)
	filesAfter := citationStringSet(tracked)
	unitDirectories := make(map[string]struct{})
	afterPath := make(map[string]string)
	movedFiles := make(map[string]string)
	movedUnits := make(map[string]string)

	for _, move := range moves {
		sourceRoot := historyMoveSourceRoot(move.From)
		sourceDirectory := path.Dir(move.From)
		destinationDirectory := path.Dir(move.To)
		for sourceDirectory == sourceRoot || strings.HasPrefix(sourceDirectory, sourceRoot+"/") {
			if sourceDirectory == sourceRoot && !historyCitationSourceRootIsUnit(sourceRoot) {
				break
			}
			unitDirectories[sourceDirectory] = struct{}{}
			if !refused[move.From] {
				unitDirectories[destinationDirectory] = struct{}{}
				movedUnits[sourceDirectory] = destinationDirectory
			}
			sourceDirectory = path.Dir(sourceDirectory)
			destinationDirectory = path.Dir(destinationDirectory)
			if sourceDirectory == "." {
				break
			}
		}
		if refused[move.From] {
			continue
		}
		delete(filesAfter, move.From)
		filesAfter[move.To] = struct{}{}
		afterPath[move.From] = move.To
		movedFiles[move.From] = move.To
	}

	existsBefore := cloneStringSet(filesBefore)
	existsAfter := cloneStringSet(filesAfter)
	markExistingUnitDirectories(existsBefore, filesBefore, unitDirectories)
	markExistingUnitDirectories(existsAfter, filesAfter, unitDirectories)
	return relocationResolution{
		existsBefore: existsBefore,
		existsAfter:  existsAfter,
		afterPath:    afterPath,
		movedFiles:   movedFiles,
		movedUnits:   movedUnits,
	}
}

func historyCitationSourceRootIsUnit(sourceRoot string) bool {
	return strings.HasPrefix(sourceRoot, "docs/specs/_reviews/") ||
		strings.HasPrefix(sourceRoot, "docs/specs/reviews/")
}

func citationStringSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func cloneStringSet(values map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for value := range values {
		result[value] = struct{}{}
	}
	return result
}

func markExistingUnitDirectories(
	exists map[string]struct{},
	files map[string]struct{},
	unitDirectories map[string]struct{},
) {
	for file := range files {
		for directory := path.Dir(file); directory != "."; directory = path.Dir(directory) {
			if _, ok := unitDirectories[directory]; ok {
				exists[directory] = struct{}{}
			}
		}
	}
}

func readTrackedCitationFile(
	root string,
	relative string,
	directorySafety map[string]bool,
) ([]byte, trackedReadResult) {
	info, safe := lstatTrackedCitationPath(root, relative, directorySafety)
	if !safe || info == nil || !info.Mode().IsRegular() {
		return nil, trackedReadSkipped
	}
	if info.Size() > relocationCitationMaxFileSize {
		return nil, trackedReadUnscanned
	}

	absolute := filepath.Join(root, filepath.FromSlash(relative))
	file, err := openCitationFileNoFollow(absolute)
	if err != nil {
		current, currentSafe := lstatTrackedCitationPath(root, relative, make(map[string]bool))
		if !currentSafe || current == nil || !current.Mode().IsRegular() || !os.SameFile(info, current) {
			return nil, trackedReadSkipped
		}
		return nil, trackedReadUnscanned
	}
	defer file.Close()

	current, currentSafe := lstatTrackedCitationPath(root, relative, make(map[string]bool))
	if !currentSafe || current == nil || !current.Mode().IsRegular() || !os.SameFile(info, current) {
		return nil, trackedReadSkipped
	}
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, trackedReadUnscanned
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(current, openedInfo) {
		return nil, trackedReadSkipped
	}
	content, err := io.ReadAll(io.LimitReader(file, relocationCitationMaxFileSize+1))
	if err != nil || len(content) > relocationCitationMaxFileSize {
		return nil, trackedReadUnscanned
	}
	probe := content
	if len(probe) > relocationCitationBinaryProbe {
		probe = probe[:relocationCitationBinaryProbe]
	}
	if bytes.IndexByte(probe, 0) >= 0 {
		return nil, trackedReadSkipped
	}
	return content, trackedReadScanned
}

func lstatTrackedCitationPath(
	root string,
	relative string,
	directorySafety map[string]bool,
) (fs.FileInfo, bool) {
	rootInfo, err := os.Lstat(root)
	if err != nil || !rootInfo.IsDir() || rootInfo.Mode()&fs.ModeSymlink != 0 {
		return nil, false
	}
	parts := strings.Split(relative, "/")
	for index := 1; index < len(parts); index++ {
		directory := strings.Join(parts[:index], "/")
		if safe, known := directorySafety[directory]; known {
			if !safe {
				return nil, false
			}
			continue
		}
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(directory)))
		safe := err == nil && info.IsDir() && info.Mode()&fs.ModeSymlink == 0
		directorySafety[directory] = safe
		if !safe {
			return nil, false
		}
	}
	info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(relative)))
	if err != nil || info.Mode()&fs.ModeSymlink != 0 {
		return nil, false
	}
	return info, true
}

func citationsBrokenByRelocation(
	citingPath string,
	content []byte,
	resolution relocationResolution,
) []relocationCitation {
	afterCitingPath := citingPath
	if moved, ok := resolution.afterPath[citingPath]; ok {
		afterCitingPath = moved
	}
	markdown := strings.EqualFold(path.Ext(citingPath), ".md") || strings.EqualFold(path.Ext(citingPath), ".markdown")
	seen := make(map[string]struct{})
	citations := make([]relocationCitation, 0)
	fenceMarker := byte(0)
	fenceLength := 0
	for lineIndex, line := range strings.Split(string(content), "\n") {
		lineNumber := lineIndex + 1
		insideFence := fenceMarker != 0
		tokenLine := line
		if markdown {
			destinations := markdownDestinations(line)
			tokenLine = maskMarkdownDestinations(line, destinations)
			marker, length, fence := markdownFence(line)
			if fence {
				if fenceMarker == 0 {
					fenceMarker, fenceLength = marker, length
				} else if marker == fenceMarker && length >= fenceLength {
					fenceMarker, fenceLength = 0, 0
				}
				insideFence = true
			}
			if !insideFence {
				for _, destination := range destinations {
					text, target, ok := normalizeMarkdownDestination(destination.raw)
					if !ok {
						continue
					}
					before, ok := resolveMarkdownDestination(path.Dir(citingPath), target)
					if !ok {
						continue
					}
					after, ok := resolveMarkdownDestination(path.Dir(afterCitingPath), target)
					if !ok {
						continue
					}
					appendBrokenCitation(&citations, seen, lineNumber, text, before, after, resolution)
				}
			}
		}
		for _, token := range repositoryPathTokens(tokenLine) {
			appendBrokenCitation(&citations, seen, lineNumber, token, token, token, resolution)
		}
	}
	sort.Slice(citations, func(left int, right int) bool {
		if citations[left].line != citations[right].line {
			return citations[left].line < citations[right].line
		}
		if citations[left].text != citations[right].text {
			return citations[left].text < citations[right].text
		}
		return citations[left].before < citations[right].before
	})
	return citations
}

func appendBrokenCitation(
	citations *[]relocationCitation,
	seen map[string]struct{},
	line int,
	text string,
	before string,
	after string,
	resolution relocationResolution,
) {
	if _, exists := resolution.existsBefore[before]; !exists {
		return
	}
	if _, exists := resolution.existsAfter[after]; exists {
		return
	}
	key := fmt.Sprintf("%d\x00%s", line, before)
	if _, duplicate := seen[key]; duplicate {
		return
	}
	seen[key] = struct{}{}
	*citations = append(*citations, relocationCitation{
		line:        line,
		text:        text,
		before:      before,
		after:       after,
		destination: relocationDestination(before, resolution),
	})
}

func relocationDestination(before string, resolution relocationResolution) string {
	if destination, ok := resolution.movedFiles[before]; ok {
		return destination
	}
	bestSource := ""
	bestDestination := ""
	for source, destination := range resolution.movedUnits {
		if before != source && !strings.HasPrefix(before, source+"/") {
			continue
		}
		if len(source) > len(bestSource) {
			bestSource = source
			bestDestination = destination
		}
	}
	if bestSource == "" {
		return ""
	}
	remainder := strings.TrimPrefix(before, bestSource)
	return bestDestination + remainder
}

func relocationCitationMessage(citations []relocationCitation) string {
	listed := citations
	if len(listed) > relocationCitationPerFileLimit {
		listed = listed[:relocationCitationPerFileLimit]
	}
	parts := make([]string, 0, len(listed)+1)
	for _, citation := range listed {
		after := "nothing exists there"
		if citation.after != citation.before {
			after = fmt.Sprintf("it resolves to %s, where nothing exists", citation.after)
		}
		message := fmt.Sprintf(
			"line %d cites %s, which resolves to %s; after this plan %s",
			citation.line,
			citation.text,
			citation.before,
			after,
		)
		if citation.destination != "" {
			message += fmt.Sprintf(" (it moves to %s)", citation.destination)
		}
		parts = append(parts, message)
	}
	if omitted := len(citations) - len(listed); omitted > 0 {
		parts = append(parts, fmt.Sprintf("and %d more", omitted))
	}
	return strings.Join(parts, "; ")
}

func markdownFence(line string) (byte, int, bool) {
	indent := 0
	for indent < len(line) && indent < 4 && line[indent] == ' ' {
		indent++
	}
	if indent > 3 || indent >= len(line) || (line[indent] != '`' && line[indent] != '~') {
		return 0, 0, false
	}
	marker := line[indent]
	length := 0
	for indent+length < len(line) && line[indent+length] == marker {
		length++
	}
	return marker, length, length >= 3
}

func markdownDestinations(line string) []markdownDestination {
	destinations := make([]markdownDestination, 0)
	for offset := 0; offset < len(line); {
		index := strings.Index(line[offset:], "](")
		if index < 0 {
			break
		}
		start := offset + index + 2
		end := markdownClosingParenthesis(line, start)
		if end < 0 {
			break
		}
		destinations = append(destinations, markdownDestination{raw: line[start:end], start: start, end: end})
		offset = end + 1
	}

	trimmed := strings.TrimLeft(line, " ")
	indent := len(line) - len(trimmed)
	if indent <= 3 && strings.HasPrefix(trimmed, "[") {
		if labelEnd := strings.Index(trimmed, "]:"); labelEnd >= 1 {
			start := indent + labelEnd + 2
			destinations = append(destinations, markdownDestination{raw: line[start:], start: start, end: len(line)})
		}
	}
	return destinations
}

func maskMarkdownDestinations(line string, destinations []markdownDestination) string {
	if len(destinations) == 0 {
		return line
	}
	masked := []byte(line)
	for _, destination := range destinations {
		for index := destination.start; index < destination.end; index++ {
			masked[index] = ' '
		}
	}
	return string(masked)
}

func markdownClosingParenthesis(line string, start int) int {
	depth := 0
	escaped := false
	angle := false
	for index := start; index < len(line); index++ {
		character := line[index]
		if escaped {
			escaped = false
			continue
		}
		if character == '\\' {
			escaped = true
			continue
		}
		if character == '<' && depth == 0 {
			angle = true
			continue
		}
		if character == '>' && angle {
			angle = false
			continue
		}
		if angle {
			continue
		}
		switch character {
		case '(':
			depth++
		case ')':
			if depth == 0 {
				return index
			}
			depth--
		}
	}
	return -1
}

func normalizeMarkdownDestination(raw string) (string, string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.HasPrefix(raw, "#") || containsControlCharacter(raw) {
		return "", "", false
	}
	destination := raw
	if strings.HasPrefix(destination, "<") {
		end := strings.Index(destination, ">")
		if end < 0 {
			return "", "", false
		}
		destination = destination[1:end]
	} else if index := strings.IndexFunc(destination, unicode.IsSpace); index >= 0 {
		destination = destination[:index]
	}
	if destination == "" || strings.HasPrefix(destination, "//") || hasURLScheme(destination) {
		return "", "", false
	}
	if index := strings.IndexAny(destination, "?#"); index >= 0 {
		destination = destination[:index]
	}
	if destination == "" || containsControlCharacter(destination) {
		return "", "", false
	}
	decoded, err := url.PathUnescape(destination)
	if err != nil {
		decoded = destination
	}
	if decoded == "" || containsControlCharacter(decoded) || strings.HasPrefix(decoded, "//") {
		return "", "", false
	}
	return decoded, decoded, true
}

func hasURLScheme(value string) bool {
	colon := strings.IndexByte(value, ':')
	if colon <= 0 {
		return false
	}
	for index := 0; index < colon; index++ {
		character := value[index]
		if index == 0 {
			if (character < 'A' || character > 'Z') && (character < 'a' || character > 'z') {
				return false
			}
			continue
		}
		if (character >= 'A' && character <= 'Z') || (character >= 'a' && character <= 'z') ||
			(character >= '0' && character <= '9') || character == '+' || character == '-' || character == '.' {
			continue
		}
		return false
	}
	return true
}

func resolveMarkdownDestination(base string, destination string) (string, bool) {
	var resolved string
	if strings.HasPrefix(destination, "/") {
		resolved = strings.TrimPrefix(destination, "/")
	} else {
		resolved = path.Join(base, destination)
	}
	resolved = path.Clean(resolved)
	if !repositoryPathIsSafe(resolved) {
		return "", false
	}
	return resolved, true
}

func repositoryPathTokens(line string) []string {
	tokens := make([]string, 0)
	var run strings.Builder
	appendRun := func() {
		token := run.String()
		run.Reset()
		if !strings.Contains(token, "/") {
			return
		}
		if strings.HasPrefix(token, "./") {
			token = strings.TrimPrefix(token, "./")
		}
		token = strings.TrimRight(token, ".,:;/")
		if repositoryPathIsSafe(token) && !containsControlCharacter(token) {
			tokens = append(tokens, token)
		}
	}
	for _, character := range line {
		if repositoryPathTokenRune(character) {
			run.WriteRune(character)
			continue
		}
		appendRun()
	}
	appendRun()
	return tokens
}

func repositoryPathTokenRune(character rune) bool {
	return unicode.IsLetter(character) || unicode.IsDigit(character) ||
		strings.ContainsRune("._~/@+-", character)
}

func containsControlCharacter(value string) bool {
	for _, character := range value {
		if character <= '\u001f' || (character >= '\u007f' && character <= '\u009f') {
			return true
		}
	}
	return false
}

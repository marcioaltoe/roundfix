package spec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type CancelRequest struct {
	SpecsRoot, Slug, Reason, SourceRevision, RepositoryRoot string
	BuiltInRoot                                             bool
	CancelledAt                                             time.Time
}

type ReturnedSource struct{ From, To, Status string }

type CancelResult struct {
	RecordPath         string
	Completed, Pending []string
	Returned           []ReturnedSource
	RemovedFiles       int
	RemovedBytes       int64
	RemovedPaths       []string
}

type cancelCopy struct {
	path    string
	content []byte
}

// Cancel records the current Task state and returns adopted intent and evidence
// before removing the Spec. Validation and planning never modify the tree.
func Cancel(req CancelRequest) (CancelResult, error) {
	if req.Slug == "" || filepath.Base(req.Slug) != req.Slug || req.Slug == "." || req.Slug == ".." || strings.Contains(req.Slug, "\\") {
		return CancelResult{}, fmt.Errorf("unsafe Spec slug %q", req.Slug)
	}
	sourceDir := filepath.Join(req.SpecsRoot, req.Slug)
	info, err := os.Lstat(sourceDir)
	if err != nil {
		return CancelResult{}, fmt.Errorf("read cancellation Spec folder: %w", err)
	}
	if !info.IsDir() {
		return CancelResult{}, fmt.Errorf("cancellation Spec %q is not a directory", sourceDir)
	}
	prd, err := os.ReadFile(filepath.Join(sourceDir, "_prd.md"))
	if err != nil {
		return CancelResult{}, fmt.Errorf("read cancellation PRD: %w", err)
	}
	if err := cancelVacant(filepath.Join(sourceDir, SupersessionFilename)); err != nil {
		return CancelResult{}, fmt.Errorf("cannot cancel a Spec with a supersession: %w", err)
	}
	reason := strings.Join(strings.Fields(req.Reason), " ")
	if reason == "" {
		return CancelResult{}, errors.New("cancellation requires a non-blank reason")
	}
	if !regexp.MustCompile(`^[0-9a-fA-F]{40}$`).MatchString(req.SourceRevision) {
		return CancelResult{}, errors.New("cancellation requires a 40-hex SourceRevision")
	}
	result := CancelResult{Completed: []string{}, Pending: []string{}}
	if _, err := os.Lstat(filepath.Join(sourceDir, "_tasks.md")); err == nil {
		graph, err := Load(req.SpecsRoot, req.Slug)
		if err != nil {
			return CancelResult{}, fmt.Errorf("read cancellation Task Graph: %w", err)
		}
		for _, task := range graph.Tasks {
			if task.Status == StatusCompleted {
				result.Completed = append(result.Completed, task.ID)
			} else {
				result.Pending = append(result.Pending, task.ID)
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return CancelResult{}, fmt.Errorf("stat cancellation Task Graph: %w", err)
	}
	archiveRoot := ArchiveSpecRoot(req.SpecsRoot, req.BuiltInRoot)
	result.RecordPath = ArchiveRecordPath(archiveRoot, req.Slug)
	for _, destination := range []string{result.RecordPath, filepath.Join(archiveRoot, req.Slug)} {
		if err := cancelVacant(destination); err != nil {
			return CancelResult{}, err
		}
	}
	if req.RepositoryRoot == "" {
		return CancelResult{}, errors.New("cancellation requires RepositoryRoot")
	}
	relative := func(path string) (string, error) {
		rel, err := RepositoryPath(req.RepositoryRoot, path)
		if err != nil || !cancelRelativePath(rel) {
			return "", fmt.Errorf("cancellation path %q is outside the repository", path)
		}
		return filepath.ToSlash(rel), nil
	}
	sourceRel, err := relative(sourceDir)
	if err != nil {
		return CancelResult{}, err
	}
	date := archiveDate(req.CancelledAt)
	copies, returned, err := planSourceReturn(req.RepositoryRoot, sourceDir, date, req.Slug, reason)
	if err != nil {
		return CancelResult{}, err
	}
	result.Returned = returned
	specRootRel, err := relative(req.SpecsRoot)
	if err != nil {
		return CancelResult{}, err
	}
	archiveRootRel, err := relative(archiveRoot)
	if err != nil {
		return CancelResult{}, err
	}
	recordRel, err := relative(result.RecordPath)
	if err != nil {
		return CancelResult{}, err
	}
	ctx := context.Background()
	relink, err := PlanArchiveRelink(ctx, req.RepositoryRoot, specRootRel, archiveRootRel, sourceRel, recordRel, nil)
	if err != nil {
		return CancelResult{}, fmt.Errorf("plan cancellation relink: %w", err)
	}
	originals := make(map[string][]byte)
	var relinked, promoted []string
	for path := range relink.Relinked {
		relinked = append(relinked, path)
	}
	sort.Strings(relinked)
	for _, path := range relinked {
		changes, err := archiveGitRead(ctx, req.RepositoryRoot, "status", "--porcelain", "--untracked-files=all", "--", path)
		if err != nil {
			return CancelResult{}, fmt.Errorf("inspect cancellation relink %q: %w", path, err)
		}
		head, err := archiveGitRead(ctx, req.RepositoryRoot, "show", "HEAD:"+path)
		if err != nil {
			return CancelResult{}, fmt.Errorf("read cancellation relink HEAD %q: %w", path, err)
		}
		absolute := filepath.Join(req.RepositoryRoot, filepath.FromSlash(path))
		info, err := os.Lstat(absolute)
		if err != nil {
			return CancelResult{}, fmt.Errorf("stat cancellation relink %q: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return CancelResult{}, fmt.Errorf("cancellation relink path %q is not a regular file", path)
		}
		current, err := os.ReadFile(absolute)
		if err != nil {
			return CancelResult{}, fmt.Errorf("read cancellation relink %q: %w", path, err)
		}
		if len(bytes.TrimSpace(changes)) > 0 || !bytes.Equal(current, head) {
			return CancelResult{}, fmt.Errorf("Spec %q cannot cancel: %s has changes not committed at HEAD; commit them so the cancellation can relink it", req.Slug, path)
		}
		originals[path] = current
	}
	destinations := make(map[string]bool)
	for _, copy := range copies {
		destinations[copy.path] = true
	}
	for _, path := range relink.Promote {
		if !cancelRelativePath(path) || archiveCoreArtifact(path) {
			return CancelResult{}, fmt.Errorf("cannot promote cancellation path %q", path)
		}
		source := filepath.Join(sourceDir, filepath.FromSlash(path))
		if err := cancelRegularPath(sourceDir, path); err != nil {
			return CancelResult{}, err
		}
		destination := filepath.Join(req.RepositoryRoot, "docs", "references", filepath.Base(path))
		if _, err := RepositoryPath(req.RepositoryRoot, destination); err != nil {
			return CancelResult{}, fmt.Errorf("cancellation path %q is outside the repository", destination)
		}
		if err := cancelVacant(destination); err != nil {
			return CancelResult{}, err
		}
		if destinations[destination] {
			return CancelResult{}, fmt.Errorf("duplicate promotion destination %q", destination)
		}
		destinations[destination] = true
		data, err := os.ReadFile(source)
		if err != nil {
			return CancelResult{}, fmt.Errorf("read cancellation promotion: %w", err)
		}
		copies = append(copies, cancelCopy{destination, data})
		promoted = append(promoted, "docs/references/"+filepath.Base(path))
	}
	// BuildArchiveRecord observes QA metadata; cancellation always keeps the PRD's intent.
	record, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: sourceDir, Slug: req.Slug, Source: sourceRel, SourceRevision: req.SourceRevision, Archived: date, Promoted: promoted})
	if err != nil {
		return CancelResult{}, fmt.Errorf("build cancelled archive record: %w", err)
	}
	_, body, err := splitFrontmatter(prd)
	if err != nil {
		return CancelResult{}, fmt.Errorf("parse cancellation PRD: %w", err)
	}
	record.Disposition, record.CancelReason = ArchiveCancelled, reason
	record.CompletedTasks, record.PendingTasks = result.Completed, result.Pending
	record.ReturnedSources = []string{}
	for _, source := range returned {
		record.ReturnedSources = append(record.ReturnedSources, source.To)
	}
	record.Outcome, record.Archived, record.Relinked = archiveParagraph(string(body), "# "), date, relinked
	content, err := RenderArchiveRecord(record)
	if err != nil {
		return CancelResult{}, fmt.Errorf("render cancelled archive record: %w", err)
	}
	if err := filepath.WalkDir(sourceDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, err := relative(path)
		if err != nil {
			return err
		}
		result.RemovedPaths = append(result.RemovedPaths, rel)
		result.RemovedFiles++
		result.RemovedBytes += info.Size()
		return nil
	}); err != nil {
		return CancelResult{}, fmt.Errorf("inspect cancelled Spec files: %w", err)
	}
	removalFolder, err := filepath.Abs(sourceDir)
	if err != nil {
		return CancelResult{}, fmt.Errorf("resolve Spec folder for removal: %w", err)
	}
	// Record and copies are exclusive creations. Track directories too so rollback
	// restores even a repository that had no history or return directories.
	copies = append([]cancelCopy{{result.RecordPath, content}}, copies...)
	var written, directories, rewritten []string
	rollback := func(cause error) error {
		for i := len(rewritten) - 1; i >= 0; i-- {
			path := rewritten[i]
			if err := os.WriteFile(filepath.Join(req.RepositoryRoot, filepath.FromSlash(path)), originals[path], 0o644); err != nil {
				cause = errors.Join(cause, fmt.Errorf("restore cancellation relink %q: %w", path, err))
			}
		}
		for i := len(written) - 1; i >= 0; i-- {
			if err := os.Remove(written[i]); err != nil {
				cause = errors.Join(cause, fmt.Errorf("remove cancellation write %q: %w", written[i], err))
			}
		}
		for i := len(directories) - 1; i >= 0; i-- {
			if err := os.Remove(directories[i]); err != nil {
				cause = errors.Join(cause, fmt.Errorf("remove cancellation directory %q: %w", directories[i], err))
			}
		}
		return cause
	}
	var mkdir func(string) error
	mkdir = func(path string) error {
		info, err := os.Stat(path)
		if err == nil {
			if !info.IsDir() {
				return fmt.Errorf("cancellation destination parent %q is not a directory", path)
			}
			return nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := mkdir(filepath.Dir(path)); err != nil {
			return err
		}
		if err := os.Mkdir(path, 0o755); err != nil {
			return err
		}
		directories = append(directories, path)
		return nil
	}
	for _, copy := range copies {
		if err := mkdir(filepath.Dir(copy.path)); err != nil {
			return CancelResult{}, rollback(fmt.Errorf("create cancellation destination directory: %w", err))
		}
		file, err := os.OpenFile(copy.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err != nil {
			return CancelResult{}, rollback(fmt.Errorf("create cancellation destination %q: %w", copy.path, err))
		}
		written = append(written, copy.path)
		n, writeErr := file.Write(copy.content)
		if writeErr == nil && n != len(copy.content) {
			writeErr = io.ErrShortWrite
		}
		if err := errors.Join(writeErr, file.Close()); err != nil {
			return CancelResult{}, rollback(fmt.Errorf("write cancellation destination %q: %w", copy.path, err))
		}
	}
	for _, path := range relinked {
		// A failed write may already have truncated the live document.
		rewritten = append(rewritten, path)
		if err := os.WriteFile(filepath.Join(req.RepositoryRoot, filepath.FromSlash(path)), relink.Relinked[path], 0o644); err != nil {
			return CancelResult{}, rollback(fmt.Errorf("write cancellation relink %q: %w", path, err))
		}
	}
	if err := os.RemoveAll(sourceDir); err != nil {
		return result, &FolderRemovalError{Folder: removalFolder, Err: err}
	}
	return result, nil
}

func cancelVacant(path string) error {
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("cancellation destination %q already exists", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat cancellation destination %q: %w", path, err)
	}
	return nil
}

func cancelRelativePath(path string) bool {
	clean := filepath.Clean(path)
	return path != "" && !filepath.IsAbs(path) && clean != "." && clean != ".." && !strings.HasPrefix(clean, ".."+string(filepath.Separator)) && !strings.ContainsAny(path, "\\\x00")
}

func cancelRegularPath(root, path string) error {
	if !cancelRelativePath(path) {
		return fmt.Errorf("indexed cancellation path %q is unsafe", path)
	}
	current := root
	for _, part := range strings.Split(filepath.Clean(path), string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return fmt.Errorf("stat indexed cancellation path %q: %w", path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("indexed cancellation path %q is a symbolic link", path)
		}
	}
	info, err := os.Lstat(current)
	if err != nil {
		return fmt.Errorf("stat indexed cancellation source %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("indexed cancellation source %q is not a regular file", path)
	}
	return nil
}

func planSourceReturn(root, specDir, date, slug, reason string) ([]cancelCopy, []ReturnedSource, error) {
	index, err := os.ReadFile(filepath.Join(specDir, "references", "_index.md"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("read cancellation source index: %w", err)
	}
	pathColumn, typeColumn := -1, -1
	var copies []cancelCopy
	var returned []ReturnedSource
	destinations := make(map[string]bool)
	for _, line := range strings.Split(string(index), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		for i := range cells {
			cells[i] = strings.Trim(strings.TrimSpace(cells[i]), "`")
		}
		if pathColumn < 0 {
			for i, cell := range cells {
				if strings.EqualFold(cell, "path") {
					pathColumn = i
				}
				if strings.EqualFold(cell, "type") {
					typeColumn = i
				}
			}
			continue
		}
		if pathColumn >= len(cells) {
			continue
		}
		path := cells[pathColumn]
		if path == "" || strings.HasPrefix(path, "-") || strings.HasPrefix(path, ":") {
			continue
		}
		if err := cancelRegularPath(specDir, filepath.Join("references", path)); err != nil {
			return nil, nil, err
		}
		// Validate the unprefixed path too: joining must not conceal an escape.
		if !cancelRelativePath(path) {
			return nil, nil, fmt.Errorf("indexed cancellation path %q is unsafe", path)
		}
		if typeColumn < 0 || typeColumn >= len(cells) {
			continue
		}
		kind := cells[typeColumn]
		if kind != "backlog" && kind != "finding" {
			continue
		}
		family, status := "backlog", "open"
		if kind == "finding" {
			family, status = "findings", "pending"
		}
		to := filepath.ToSlash(filepath.Join("docs", family, path))
		destination := filepath.Join(root, filepath.FromSlash(to))
		if _, err := RepositoryPath(root, destination); err != nil {
			return nil, nil, fmt.Errorf("cancellation path %q is outside the repository", destination)
		}
		if err := cancelVacant(destination); err != nil {
			return nil, nil, err
		}
		if destinations[destination] {
			return nil, nil, fmt.Errorf("duplicate return destination %q", to)
		}
		destinations[destination] = true
		from := filepath.Join(specDir, "references", path)
		content, err := os.ReadFile(from)
		if err != nil {
			return nil, nil, fmt.Errorf("read returned source %q: %w", from, err)
		}
		next, err := cancelSourceContent(content, kind, date, slug, reason)
		if err != nil {
			return nil, nil, fmt.Errorf("return source %q: %w", from, err)
		}
		rel, err := RepositoryPath(root, from)
		if err != nil || !cancelRelativePath(rel) {
			return nil, nil, fmt.Errorf("cancellation path %q is outside the repository", from)
		}
		copies = append(copies, cancelCopy{destination, next})
		returned = append(returned, ReturnedSource{filepath.ToSlash(rel), to, status})
	}
	return copies, returned, nil
}

// Edit physical front matter lines instead of re-encoding YAML: comments,
// unrelated formatting and the entire original observation keep their bytes.
func cancelSourceContent(content []byte, kind, date, slug, reason string) ([]byte, error) {
	lines := strings.SplitAfter(string(content), "\n")
	if len(lines) == 0 || strings.TrimRight(lines[0], "\r\n") != "---" {
		return nil, errors.New("returned source missing YAML front matter")
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], "\r\n") == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return nil, errors.New("returned source missing YAML closing marker")
	}
	var node yaml.Node
	frontmatter := strings.Join(lines[1:end], "")
	if err := yaml.Unmarshal([]byte(frontmatter), &node); err != nil {
		return nil, fmt.Errorf("parse returned source front matter: %w", err)
	}
	mapping := archiveFrontmatterMapping(&node)
	if mapping == nil {
		return nil, errors.New("returned source front matter must be a mapping")
	}
	// Decode once to reject duplicate keys before changing physical lines.
	var fields map[string]yaml.Node
	if err := yaml.Unmarshal([]byte(frontmatter), &fields); err != nil {
		return nil, fmt.Errorf("parse returned source fields: %w", err)
	}
	if _, ok := fields["status"]; !ok {
		return nil, errors.New("returned source missing status key")
	}
	edits := map[int]string{}
	hasSpec, hasUpdated := false, false
	statusLine := 0
	newline := "\n"
	if strings.HasSuffix(lines[0], "\r\n") {
		newline = "\r\n"
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		key, value := mapping.Content[i], mapping.Content[i+1]
		replacement, change := "", false
		switch key.Value {
		case "status":
			statusLine = key.Line
			replacement, change = "open", true
			if kind == "finding" {
				replacement = "pending"
			}
		case "spec":
			hasSpec = true
			if kind == "backlog" {
				replacement, change = "null", true
			}
		case "updated_at":
			hasUpdated = true
			if kind == "finding" {
				replacement, change = date, true
			}
		case "absorbed_by":
			if kind == "finding" {
				edits[key.Line] = ""
				continue
			}
		}
		if !change {
			continue
		}
		if value.Kind != yaml.ScalarNode || value.Line != key.Line || value.Style&yaml.LiteralStyle != 0 || value.Style&yaml.FoldedStyle != 0 {
			return nil, fmt.Errorf("returned source %s must occupy one front matter line", key.Value)
		}
		comment := ""
		if value.LineComment != "" {
			line := strings.TrimRight(lines[key.Line], "\r\n")
			at := strings.LastIndex(line, value.LineComment)
			if at >= 0 {
				start := at
				for start > 0 && (line[start-1] == ' ' || line[start-1] == '\t') {
					start--
				}
				comment = line[start:]
			}
		}
		edits[key.Line] = key.Value + ": " + replacement + comment + newline
	}
	var out strings.Builder
	for i, line := range lines {
		if i == end && kind == "finding" && !hasUpdated {
			out.WriteString("updated_at: " + date + newline)
		}
		if replacement, ok := edits[i]; ok {
			out.WriteString(replacement)
		} else {
			out.WriteString(line)
		}
		if i == statusLine && kind == "backlog" && !hasSpec {
			out.WriteString("spec: null" + newline)
		}
	}
	next := out.String()
	if kind == "finding" {
		// Keep existing trailing bytes; add only the separation needed by the addendum.
		if !strings.HasSuffix(next, newline) {
			next += newline
		}
		if !strings.HasSuffix(next, newline+newline) {
			next += newline
		}
		next += "## Addendum — " + date + " — Spec " + slug + " cancelled" + newline + newline + "Spec " + slug + " was cancelled on " + date + " before it implemented this finding: " + reason + newline
	}
	return []byte(next), nil
}

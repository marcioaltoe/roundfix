package spec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type LegacyDelivery struct{ Commit, PullRequest, Date string }

type LegacyConversionRequest struct {
	RepositoryRoot string
	ArchiveRoot    string
	Slug           string
	SourceRevision string
	Delivery       LegacyDelivery
	Promote        []string
}

type LegacyConversion struct {
	Slug, Folder, RecordPath string
	Record                   ArchiveRecord
	Rendered                 []byte
	Files                    []string
	Bytes                    int64
	Promoted                 []string
	legacyCopies             map[string][]byte
}

// LegacyArchiveFolders returns sorted slugs, independent of their PRD status.
func LegacyArchiveFolders(archiveRoot string) ([]string, error) {
	entries, err := os.ReadDir(archiveRoot)
	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list legacy archives: %w", err)
	}
	folders := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		info, err := os.Lstat(filepath.Join(archiveRoot, entry.Name(), "_prd.md"))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("inspect legacy PRD: %w", err)
		}
		if info.Mode().IsRegular() {
			folders = append(folders, entry.Name())
		}
	}
	sort.Strings(folders)
	return folders, nil
}

func legacyGit(ctx context.Context, root string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return nil, fmt.Errorf("read legacy delivery with git %s: %w: %s", args[0], err, out)
	}
	return out, nil
}

// FindLegacyDelivery uses only HEAD's first-parent history, never rename guesses.
func FindLegacyDelivery(ctx context.Context, repoRoot, specRoot, archiveRoot, slug string) (LegacyDelivery, error) {
	if !legacySlug(slug) || !validArchiveSource(specRoot, strings.Repeat("0", 40)) || !validArchiveSource(archiveRoot, strings.Repeat("0", 40)) {
		return LegacyDelivery{}, errors.New("unsafe legacy delivery path")
	}
	deleted, err := legacyGit(ctx, repoRoot, "log", "--first-parent", "--no-renames", "--diff-filter=D", "-1", "--format=%H", "HEAD", "--", path.Join(specRoot, slug, "_prd.md"))
	if err != nil {
		return LegacyDelivery{}, err
	}
	commit := strings.TrimSpace(string(deleted))
	if commit == "" {
		added, err := legacyGit(ctx, repoRoot, "log", "--first-parent", "--no-renames", "--diff-filter=A", "--format=%H", "HEAD", "--", path.Join(archiveRoot, slug, "_prd.md"))
		if err != nil {
			return LegacyDelivery{}, err
		}
		for _, candidate := range strings.Fields(string(added)) {
			// Explicitly diff the first parent, including merge commits; a root has no parent.
			parents, err := legacyGit(ctx, repoRoot, "rev-list", "--parents", "-1", candidate)
			if err != nil {
				return LegacyDelivery{}, err
			}
			fields := strings.Fields(string(parents))
			args := []string{"diff-tree", "--root", "--no-renames", "--no-commit-id", "-r", "--diff-filter=D", "--name-only", "-z"}
			if len(fields) > 1 {
				args = append(args, fields[1])
			}
			args = append(args, candidate)
			removed, err := legacyGit(ctx, repoRoot, args...)
			if err != nil {
				return LegacyDelivery{}, err
			}
			relocation := false
			for _, name := range strings.Split(string(removed), "\x00") {
				if path.Base(name) == "_prd.md" && path.Base(path.Dir(name)) == slug {
					relocation = true
					break
				}
			}
			if !relocation {
				commit = candidate
				break
			}
		}
	}
	if commit == "" {
		return LegacyDelivery{}, nil
	}
	metadata, err := legacyGit(ctx, repoRoot, "show", "-s", "--format=%H%n%as%n%s", commit)
	if err != nil {
		return LegacyDelivery{}, err
	}
	fields := strings.SplitN(strings.TrimSuffix(string(metadata), "\n"), "\n", 3)
	if len(fields) != 3 {
		return LegacyDelivery{}, errors.New("invalid legacy delivery metadata")
	}
	delivery := LegacyDelivery{Commit: fields[0], Date: fields[1]}
	if match := regexp.MustCompile(`\(#([0-9]+)\)$`).FindStringSubmatch(fields[2]); match != nil {
		delivery.PullRequest = match[1]
	}
	return delivery, nil
}

func legacySlug(slug string) bool {
	return slug != "" && slug != "." && slug != ".." && path.Base(slug) == slug && !strings.ContainsAny(slug, "\\\x00")
}

// legacyPath confines paths lexically and rejects symlinks in existing components.
func legacyPath(root, relative string) (string, error) {
	if !filepath.IsAbs(root) || !validArchiveSource(relative, strings.Repeat("0", 40)) {
		return "", fmt.Errorf("unsafe legacy path %q", relative)
	}
	current := root
	for _, component := range strings.Split(relative, "/") {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("inspect legacy path %q: %w", relative, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("legacy path %q contains a symlink", relative)
		}
	}
	return current, nil
}

func legacyAbsent(name string) error {
	_, err := os.Lstat(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect legacy destination: %w", err)
	}
	return fmt.Errorf("legacy destination %q already exists", name)
}

// PlanLegacyConversion reads the folder and builds a validated record without writes.
func PlanLegacyConversion(req LegacyConversionRequest) (LegacyConversion, error) {
	c := LegacyConversion{Slug: req.Slug, legacyCopies: map[string][]byte{}}
	if !legacySlug(req.Slug) || !filepath.IsAbs(req.ArchiveRoot) {
		return c, errors.New("unsafe legacy folder")
	}
	relative, err := filepath.Rel(req.RepositoryRoot, filepath.Join(req.ArchiveRoot, req.Slug))
	if err != nil {
		return c, fmt.Errorf("resolve legacy folder: %w", err)
	}
	c.Folder = filepath.ToSlash(relative)
	if !validArchiveSource(c.Folder, req.SourceRevision) {
		return c, errors.New("unsafe legacy source or source_revision")
	}
	folder, err := legacyPath(req.RepositoryRoot, c.Folder)
	if err != nil {
		return c, err
	}
	c.RecordPath = c.Folder + ".md"
	recordPath, err := legacyPath(req.RepositoryRoot, c.RecordPath)
	if err != nil {
		return c, err
	}
	if err := legacyAbsent(recordPath); err != nil {
		return c, err
	}
	err = filepath.WalkDir(folder, func(name string, entry fs.DirEntry, walkErr error) error {
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
		if !info.Mode().IsRegular() {
			return fmt.Errorf("legacy archive file %q is not regular", name)
		}
		rel, err := filepath.Rel(folder, name)
		if err != nil {
			return err
		}
		c.Files = append(c.Files, filepath.ToSlash(rel))
		c.Bytes += info.Size()
		return nil
	})
	if err != nil {
		return c, fmt.Errorf("inventory legacy folder: %w", err)
	}
	sort.Strings(c.Files)
	for _, promotion := range req.Promote {
		if !strings.HasPrefix(promotion, c.Folder+"/") {
			return c, fmt.Errorf("promotion %q is outside legacy folder", promotion)
		}
		source, err := legacyPath(req.RepositoryRoot, promotion)
		if err != nil {
			return c, err
		}
		rel := strings.TrimPrefix(promotion, c.Folder+"/")
		base := path.Base(rel)
		if (path.Dir(rel) == "." && strings.HasPrefix(base, "_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "task_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "qa-report-") && strings.HasSuffix(base, ".md")) {
			return c, fmt.Errorf("cannot promote core artifact %q", promotion)
		}
		info, err := os.Lstat(source)
		if err != nil {
			return c, fmt.Errorf("inspect promotion: %w", err)
		}
		if !info.Mode().IsRegular() {
			return c, fmt.Errorf("promotion %q is not regular", promotion)
		}
		destination := path.Join("docs/references", base)
		target, err := legacyPath(req.RepositoryRoot, destination)
		if err != nil {
			return c, err
		}
		if err := legacyAbsent(target); err != nil {
			return c, err
		}
		if _, exists := c.legacyCopies[destination]; exists {
			return c, fmt.Errorf("duplicate promotion destination %q", destination)
		}
		content, err := os.ReadFile(source)
		if err != nil {
			return c, fmt.Errorf("read promotion: %w", err)
		}
		c.legacyCopies[destination] = content
		c.Promoted = append(c.Promoted, destination)
	}
	sort.Strings(c.Promoted)
	c.Record, err = BuildArchiveRecord(ArchiveRecordInput{SpecDir: folder, Slug: req.Slug, Source: c.Folder, SourceRevision: req.SourceRevision, Promoted: c.Promoted})
	if err != nil {
		return c, fmt.Errorf("build legacy record: %w", err)
	}
	c.Record.PullRequest, c.Record.DeliveryCommit = req.Delivery.PullRequest, req.Delivery.Commit
	if c.Record.Archived == "" {
		c.Record.Archived = req.Delivery.Date
	}
	c.Rendered, err = RenderArchiveRecord(c.Record)
	if err != nil {
		return c, fmt.Errorf("render legacy record: %w", err)
	}
	parsed, err := ParseArchiveRecord(c.Rendered)
	if err != nil {
		return c, fmt.Errorf("parse legacy record: %w", err)
	}
	// Compare the YAML nodes so nil and empty lists share the renderer's representation.
	if !reflect.DeepEqual(archiveRecordMapping(c.Record), archiveRecordMapping(parsed)) {
		return c, errors.New("legacy record metadata does not round-trip")
	}
	rerendered, err := RenderArchiveRecord(parsed)
	if err != nil {
		return c, fmt.Errorf("render parsed legacy record: %w", err)
	}
	if !bytes.Equal(c.Rendered, rerendered) {
		return c, errors.New("legacy record does not round-trip")
	}
	c.Record = parsed
	return c, nil
}

// ApplyLegacyConversion rolls back new files on any failure before removal.
func ApplyLegacyConversion(repositoryRoot string, c LegacyConversion) (resultErr error) {
	folder, err := legacyPath(repositoryRoot, c.Folder)
	if err != nil {
		return err
	}
	written := []string{}
	defer func() {
		if resultErr == nil {
			return
		}
		for i := len(written) - 1; i >= 0; i-- {
			if err := os.Remove(written[i]); err != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("rollback legacy output: %w", err))
			}
		}
	}()
	write := func(relative string, content []byte) error {
		target, err := legacyPath(repositoryRoot, relative)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create legacy destination directory: %w", err)
		}
		file, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("create legacy output: %w", err)
		}
		written = append(written, target)
		_, writeErr := file.Write(content)
		return errors.Join(writeErr, file.Close())
	}
	if err := write(c.RecordPath, c.Rendered); err != nil {
		return fmt.Errorf("write legacy record: %w", err)
	}
	for _, destination := range c.Promoted {
		content, ok := c.legacyCopies[destination]
		if !ok {
			return fmt.Errorf("missing legacy promotion %q", destination)
		}
		if err := write(destination, content); err != nil {
			return fmt.Errorf("copy legacy promotion: %w", err)
		}
	}
	if err := os.RemoveAll(folder); err != nil {
		// Removal has begun: retain provenance and copies for operator recovery.
		written = nil
		return fmt.Errorf("remove legacy folder: %w", err)
	}
	return nil
}

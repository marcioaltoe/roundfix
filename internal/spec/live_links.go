package spec

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// PlanLiveLinks repairs missing relative links in the same live pages scanned by
// SC-HISTORY-RETAINED. Archive Records have their own record-links unit.
func PlanLiveLinks(ctx context.Context, root, specRootRel, archiveRootRel string, skip map[string]bool) (RecordLinkPlan, error) {
	p := RecordLinkPlan{Files: map[string][]byte{}, rewrites: map[string]map[string]GitLocator{}}
	if !historyRepositoryPath(specRootRel) || !historyRepositoryPath(archiveRootRel) {
		return p, fmt.Errorf("invalid live-links root")
	}
	seen := map[string]bool{}
	for _, scanRoot := range []string{"docs", specRootRel} {
		err := filepath.WalkDir(filepath.Join(root, scanRoot), func(full string, entry fs.DirEntry, walkErr error) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			if errors.Is(walkErr, os.ErrNotExist) {
				return nil
			}
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || !entry.Type().IsRegular() {
				return nil
			}
			path, err := filepath.Rel(root, full)
			if err != nil {
				return err
			}
			path = filepath.ToSlash(path)
			if seen[path] || skip[path] || !strings.HasSuffix(path, ".md") {
				return nil
			}
			seen[path] = true
			if !IsLiveDocument(path, specRootRel, archiveRootRel) && !strings.HasPrefix(path, "docs/user-guide/") {
				return nil
			}
			if err := historyRewriteRegularFile(root, path); err != nil {
				return err
			}
			content, err := os.ReadFile(full)
			if err != nil {
				return fmt.Errorf("read live links %s: %w", path, err)
			}
			if filepath.ToSlash(filepath.Dir(path)) == archiveRootRel && IsArchiveRecord(entry.Name(), content) {
				return nil
			}
			locators := map[string]GitLocator{}
			for _, link := range MissingRelativeLinks(root, path, content) {
				if link.Self {
					continue
				}
				locator, err := ResolveGitLocator(ctx, root, link.Target, "", "")
				if err != nil {
					return err
				}
				locators[link.Target] = locator
				p.Citations = append(p.Citations, HistoryCitation{Path: path, Line: link.Line, Target: link.Target, Rewrite: locator.Text(false)})
			}
			if len(locators) == 0 {
				return nil
			}
			next, count, err := rewriteRecordLinks(content, path, locators)
			if err != nil {
				return err
			}
			p.Files[path], p.rewrites[path] = next, locators
			p.Links += count
			return nil
		})
		if err != nil {
			return p, fmt.Errorf("plan live links: %w", err)
		}
	}
	return p, nil
}

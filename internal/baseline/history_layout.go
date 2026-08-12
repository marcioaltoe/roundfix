package baseline

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/spec"
)

// HistoryRelocation is one retired file that belongs under the history root
// and is not there. ContentIdentity is read before any mutation is planned.
type HistoryRelocation struct {
	From            string
	To              string
	ContentIdentity string
}

// HistoryCollision is one retired file whose destination is occupied or is
// also claimed by another legacy source.
type HistoryCollision struct {
	From            string
	To              string
	ContentIdentity string
	Reason          string
}

type historyLayoutSource struct {
	from string
	to   string
}

type historyLayoutTree struct {
	from string
	kind spec.ArchiveKind
}

var legacyHistoryLayoutTrees = []historyLayoutTree{
	{from: "docs/specs/_archived", kind: spec.ArchiveKindSpec},
	{from: "docs/findings/_archived", kind: spec.ArchiveKindFinding},
	{from: "_archived/specs", kind: spec.ArchiveKindSpec},
	{from: "_archived/findings", kind: spec.ArchiveKindFinding},
	{from: "_archived/adr", kind: spec.ArchiveKindADR},
	{from: "_archived/backlog", kind: spec.ArchiveKindBacklog},
}

// DiscoverHistoryLayout returns the relocations that bring a repository to the
// current layout, sorted by From, and the collisions that refuse to move. It
// reads repository and local Git state but never changes either.
func DiscoverHistoryLayout(root string) ([]HistoryRelocation, []HistoryCollision, error) {
	root = filepath.Clean(root)
	info, err := os.Stat(root)
	if err != nil {
		return nil, nil, fmt.Errorf("stat repository root %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("repository root %q is not a directory", root)
	}

	sources := make([]historyLayoutSource, 0)
	for _, tree := range legacyHistoryLayoutTrees {
		to := spec.ArchiveDir(tree.kind)
		if err := historyLayoutAppendTree(root, tree.from, to, &sources); err != nil {
			return nil, nil, err
		}
	}

	if err := historyLayoutAppendRetiredDocuments(
		root,
		"docs/adr",
		spec.ArchiveDir(spec.ArchiveKindADR),
		spec.ClassifyADR,
		&sources,
	); err != nil {
		return nil, nil, err
	}
	if err := historyLayoutAppendRetiredDocuments(
		root,
		"docs/backlog",
		spec.ArchiveDir(spec.ArchiveKindBacklog),
		spec.ClassifyBacklogEntry,
		&sources,
	); err != nil {
		return nil, nil, err
	}

	for _, reviewRoot := range []string{"docs/specs/_reviews", "docs/specs/reviews"} {
		if err := historyLayoutAppendFinishedReviews(root, reviewRoot, &sources); err != nil {
			return nil, nil, err
		}
	}

	return historyLayoutClassifySources(root, sources)
}

func historyLayoutAppendTree(root string, fromRoot string, toRoot string, sources *[]historyLayoutSource) error {
	absoluteRoot := filepath.Join(root, filepath.FromSlash(fromRoot))
	info, err := os.Stat(absoluteRoot)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat history source %q: %w", fromRoot, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("history source %q is not a directory", fromRoot)
	}

	err = filepath.WalkDir(absoluteRoot, func(absolutePath string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		entryInfo, err := entry.Info()
		if err != nil {
			return err
		}
		if !entryInfo.Mode().IsRegular() {
			return fmt.Errorf("history source %q is not a regular file", absolutePath)
		}
		relative, err := filepath.Rel(absoluteRoot, absolutePath)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		*sources = append(*sources, historyLayoutSource{
			from: path.Join(fromRoot, relative),
			to:   path.Join(toRoot, relative),
		})
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk history source %q: %w", fromRoot, err)
	}
	return nil
}

func historyLayoutAppendRetiredDocuments(
	root string,
	fromRoot string,
	toRoot string,
	classify func([]byte) spec.Retirement,
	sources *[]historyLayoutSource,
) error {
	candidates := make([]historyLayoutSource, 0)
	if err := historyLayoutAppendTree(root, fromRoot, toRoot, &candidates); err != nil {
		return err
	}
	for _, candidate := range candidates {
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(candidate.from)))
		if err != nil {
			return fmt.Errorf("read active document %q: %w", candidate.from, err)
		}
		if classify(content).Retired {
			*sources = append(*sources, candidate)
		}
	}
	return nil
}

func historyLayoutAppendFinishedReviews(root string, reviewRoot string, sources *[]historyLayoutSource) error {
	absoluteRoot := filepath.Join(root, filepath.FromSlash(reviewRoot))
	entries, err := os.ReadDir(absoluteRoot)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read live Review Artifact root %q: %w", reviewRoot, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		reviewDir := filepath.Join(absoluteRoot, entry.Name())
		liveness, _, err := spec.ClassifyReview(root, reviewDir)
		if err != nil {
			return fmt.Errorf("classify Review Artifact %q: %w", path.Join(reviewRoot, entry.Name()), err)
		}
		if liveness != spec.ReviewFinished {
			continue
		}
		from := path.Join(reviewRoot, entry.Name())
		to := path.Join(spec.ArchiveDir(spec.ArchiveKindReview), entry.Name())
		if err := historyLayoutAppendTree(root, from, to, sources); err != nil {
			return err
		}
	}
	return nil
}

func historyLayoutClassifySources(root string, sources []historyLayoutSource) ([]HistoryRelocation, []HistoryCollision, error) {
	sort.Slice(sources, func(left int, right int) bool {
		if sources[left].from == sources[right].from {
			return sources[left].to < sources[right].to
		}
		return sources[left].from < sources[right].from
	})

	destinationSources := make(map[string][]string, len(sources))
	for _, source := range sources {
		destinationSources[source.to] = append(destinationSources[source.to], source.from)
	}

	relocations := make([]HistoryRelocation, 0, len(sources))
	collisions := make([]HistoryCollision, 0)
	for _, source := range sources {
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(source.from)))
		if err != nil {
			return nil, nil, fmt.Errorf("read history source %q: %w", source.from, err)
		}
		identity := planContentIdentity(content)

		if competing := destinationSources[source.to]; len(competing) > 1 {
			other := historyLayoutOtherSource(competing, source.from)
			collisions = append(collisions, HistoryCollision{
				From:            source.from,
				To:              source.to,
				ContentIdentity: identity,
				Reason:          fmt.Sprintf("destination is also claimed by legacy source %q", other),
			})
			continue
		}

		_, err = os.Lstat(filepath.Join(root, filepath.FromSlash(source.to)))
		switch {
		case err == nil:
			collisions = append(collisions, HistoryCollision{
				From:            source.from,
				To:              source.to,
				ContentIdentity: identity,
				Reason:          "destination already exists",
			})
			continue
		case !errors.Is(err, os.ErrNotExist):
			return nil, nil, fmt.Errorf("inspect history destination %q: %w", source.to, err)
		}

		relocations = append(relocations, HistoryRelocation{
			From:            source.from,
			To:              source.to,
			ContentIdentity: identity,
		})
	}

	sort.Slice(collisions, func(left int, right int) bool {
		if collisions[left].From == collisions[right].From {
			return collisions[left].To < collisions[right].To
		}
		return collisions[left].From < collisions[right].From
	})
	return relocations, collisions, nil
}

func historyLayoutOtherSource(sources []string, current string) string {
	for _, source := range sources {
		if source != current {
			return source
		}
	}
	return strings.Join(sources, ", ")
}

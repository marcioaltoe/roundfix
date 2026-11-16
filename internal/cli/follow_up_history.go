package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"roundfix/internal/preflight"
	"roundfix/internal/queueforecast"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type deliveryHistoryInput struct {
	Runs                                    []store.Run
	Checkout, Branch, SpecRoot, ArchiveRoot string
	RetentionDays                           int
	Now                                     time.Time
	Git                                     preflight.GitRunner
}

// deliveryHistoryGitRead holds what the Git-only read found on the default
// branch head: the commit time that added each path and the content of each
// archive record at that head, never the working tree's copy.
type deliveryHistoryGitRead struct {
	Times   map[string]time.Time
	Records map[string][]byte
}

func readDeliveryHistory(ctx context.Context, in deliveryHistoryInput) ([]queueforecast.SpecTimes, followUpSource) {
	in = relativeDeliveryHistoryInput(in)
	read, err := readDeliveryHistoryGit(ctx, in)
	return deliveryHistoryFromTimes(in, read, err)
}

func deliveryHistoryFromTimes(in deliveryHistoryInput, read deliveryHistoryGitRead, err error) ([]queueforecast.SpecTimes, followUpSource) {
	in = relativeDeliveryHistoryInput(in)
	if err != nil {
		return []queueforecast.SpecTimes{}, followUpSourceResult("delivery-history", 0, err)
	}
	starts := map[string]time.Time{}
	for _, run := range in.Runs {
		if run.Kind != store.KindImplement || run.SpecSlug == "" {
			continue
		}
		old, found := starts[run.SpecSlug]
		if !found || run.CreatedAt.Before(old) {
			starts[run.SpecSlug] = run.CreatedAt
		}
	}
	history := make([]queueforecast.SpecTimes, 0, len(starts))
	delivered := 0
	for slug, start := range starts {
		row := queueforecast.SpecTimes{Spec: slug, Start: start, Plan: read.Times[filepath.ToSlash(filepath.Join(in.SpecRoot, slug, "_prd.md"))]}
		archivePath := filepath.ToSlash(filepath.Join(in.ArchiveRoot, slug+".md"))
		if end := read.Times[archivePath]; !end.IsZero() {
			if data, found := read.Records[archivePath]; found {
				if record, parseErr := spec.ParseArchiveRecord(data); parseErr == nil && record.Spec == slug && record.Disposition != spec.ArchiveCancelled {
					row.End = end
					delivered++
				}
			}
		}
		history = append(history, row)
	}
	sort.Slice(history, func(i, j int) bool { return history[i].Spec < history[j].Spec })
	return history, followUpSourceResult("delivery-history", delivered, nil)
}

// This Git-only read is separate from Run processing so the status view can run
// it alongside its forge read without sharing the Run Database reader.
func readDeliveryHistoryGit(ctx context.Context, in deliveryHistoryInput) (deliveryHistoryGitRead, error) {
	in = relativeDeliveryHistoryInput(in)
	runner := in.Git
	if runner == nil {
		runner = preflight.ExecGitRunner{}
	}
	head, found, err := preflight.ResolveDefaultBranchHead(ctx, in.Checkout, in.Branch, runner)
	if err != nil {
		return deliveryHistoryGitRead{}, err
	}
	if !found {
		return deliveryHistoryGitRead{}, fmt.Errorf("default branch head is unavailable")
	}
	since := in.Now.AddDate(0, 0, -in.RetentionDays-1).UTC().Format(time.RFC3339)
	raw, err := runner.RunGit(ctx, in.Checkout, "log", head.Ref, "--since="+since, "--no-renames", "--diff-filter=A", "--name-only", "--format=%x00%cI", "--", filepath.ToSlash(filepath.Join(in.ArchiveRoot, "*.md")), filepath.ToSlash(filepath.Join(in.SpecRoot, "*", "_prd.md")))
	if err != nil {
		return deliveryHistoryGitRead{}, fmt.Errorf("read delivery history: %w", err)
	}
	times, err := parseDeliveryHistoryLog(raw)
	if err != nil {
		return deliveryHistoryGitRead{}, err
	}
	// The archive record is read from the default branch head, so a branch, an
	// uncommitted edit or a symlink in the checkout cannot add or remove
	// delivered history. A record absent at the head is not delivered.
	archivePrefix := filepath.ToSlash(in.ArchiveRoot) + "/"
	records := map[string][]byte{}
	for path := range times {
		if !strings.HasPrefix(path, archivePrefix) {
			continue
		}
		content, showErr := runner.RunGit(ctx, in.Checkout, "cat-file", "blob", head.Ref+":"+path)
		if showErr == nil {
			records[path] = []byte(content + "\n")
		}
	}
	return deliveryHistoryGitRead{Times: times, Records: records}, nil
}

func parseDeliveryHistoryLog(raw string) (map[string]time.Time, error) {
	times := map[string]time.Time{}
	chunks := strings.Split(raw, "\x00")
	if strings.TrimSpace(chunks[0]) != "" {
		return nil, fmt.Errorf("parse delivery history: expected a commit time")
	}
	for _, chunk := range chunks[1:] {
		stamp, paths, _ := strings.Cut(chunk, "\n")
		at, err := time.Parse(time.RFC3339, strings.TrimSpace(stamp))
		if err != nil {
			return nil, fmt.Errorf("parse delivery history commit time: %w", err)
		}
		for _, path := range strings.Split(paths, "\n") {
			path = strings.TrimSpace(path)
			if path == "" {
				continue
			}
			if _, seen := times[path]; !seen {
				times[path] = at
			}
		}
	}
	return times, nil
}

func relativeDeliveryHistoryInput(in deliveryHistoryInput) deliveryHistoryInput {
	if filepath.IsAbs(in.SpecRoot) {
		in.SpecRoot = deliverPlanRelativePath(in.Checkout, in.SpecRoot)
	}
	if filepath.IsAbs(in.ArchiveRoot) {
		in.ArchiveRoot = deliverPlanRelativePath(in.Checkout, in.ArchiveRoot)
	}
	return in
}

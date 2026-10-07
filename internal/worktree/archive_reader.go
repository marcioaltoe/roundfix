package worktree

import (
	"context"
	"fmt"
	"os"
	"path"

	"roundfix/internal/spec"
)

func archivedSpecAtHead(ctx context.Context, runner gitRunner, root, head, slug string) (spec.ArchivedSpec, error) {
	return spec.ReadArchivedSpecAt(spec.ArchiveDir(spec.ArchiveKindSpec), slug, func(file string) ([]byte, error) {
		if _, err := runner.Run(ctx, root, "cat-file", "-e", head+":"+file); err != nil {
			return nil, os.ErrNotExist
		}
		return gitBlobAtHead(ctx, runner, root, head, file)
	})
}

func archivedSourceAtHead(ctx context.Context, runner gitRunner, root, head, slug string) (string, string, error) {
	archived, err := archivedSpecAtHead(ctx, runner, root, head, slug)
	if err != nil {
		return "", "", err
	}
	if archived.Form == spec.ArchivedFolder {
		return head, archived.Path, nil
	}
	r := archived.Record
	if _, err := runner.Run(ctx, root, "cat-file", "-e", r.SourceRevision+"^{commit}"); err != nil {
		return "", "", fmt.Errorf("source_revision %s is unavailable; keep the worktree for Spec %s: %w", r.SourceRevision, slug, err)
	}
	return r.SourceRevision, r.Source, nil
}

func archiveSourceUnavailable(ctx context.Context, runner gitRunner, root, head, slug string) error {
	a, err := archivedSpecAtHead(ctx, runner, root, head, slug)
	if err != nil || a.Form != spec.ArchivedRecord {
		return err
	}
	_, _, err = archivedSourceAtHead(ctx, runner, root, head, slug)
	return err
}

func qaReportRevision(ctx context.Context, runner gitRunner, root, head, slug, report string) string {
	a, err := archivedSpecAtHead(ctx, runner, root, head, slug)
	if err == nil && a.Form == spec.ArchivedRecord && pathUnderAnyGitDirectory(report, []string{path.Join(a.Record.Source, "qa")}) {
		return a.Record.SourceRevision
	}
	return head
}

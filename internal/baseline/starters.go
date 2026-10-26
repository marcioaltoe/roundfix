package baseline

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const StartersSchemaVersion = "roundfix/baseline-starters/v1"

var ErrStartersStale = errors.New("starters.plan.stale: the confirmed digest no longer matches; preview again")

type StartersRequest struct {
	Repository  string
	ProfileID   string
	ConfirmPlan string
}

type StarterEntry struct {
	Path            string `json:"path"`
	Action          string `json:"action"`
	ContentIdentity string `json:"contentIdentity"`
}

type StartersPayload struct {
	SchemaVersion string         `json:"schemaVersion"`
	Profile       string         `json:"profile"`
	State         string         `json:"state"`
	Starters      []StarterEntry `json:"starters"`
	PlanDigest    string         `json:"planDigest"`
}

// Starters previews repository-owned files and creates only confirmed absent paths.
func Starters(ctx context.Context, request StartersRequest, catalog *Catalog) (StartersPayload, error) {
	var payload StartersPayload
	if err := ctx.Err(); err != nil {
		return payload, err
	}
	repository, err := cleanRepositoryRoot(request.Repository)
	if err != nil {
		return payload, fmt.Errorf("resolve starters repository: %w", err)
	}
	profile, err := ResolveProfile(repository, request.ProfileID, catalog)
	if err != nil {
		return payload, err
	}
	files, err := catalog.StarterFiles(profile.Modules)
	if err != nil {
		return payload, err
	}
	root, err := os.OpenRoot(repository)
	if err != nil {
		return payload, fmt.Errorf("open starters repository: %w", err)
	}
	defer root.Close()
	payload = StartersPayload{SchemaVersion: StartersSchemaVersion, Profile: profile.ID, State: "current", Starters: make([]StarterEntry, 0, len(files))}
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return payload, err
		}
		action := "keep"
		if _, err := root.Lstat(filepath.FromSlash(file.Path)); err != nil {
			action = "create"
			payload.State = "preview"
		}
		payload.Starters = append(payload.Starters, StarterEntry{Path: file.Path, Action: action, ContentIdentity: file.ContentIdentity})
	}
	digest, err := canonicalSHA256(map[string]any{"profile": payload.Profile, "starters": payload.Starters})
	if err != nil {
		return payload, fmt.Errorf("digest starters preview: %w", err)
	}
	payload.PlanDigest = "sha256:" + digest
	if request.ConfirmPlan == "" || payload.State == "current" {
		return payload, nil
	}
	if request.ConfirmPlan != payload.PlanDigest {
		return payload, fmt.Errorf("confirm starters: %w", ErrStartersStale)
	}
	created := make([]string, 0, len(files))
	rollback := func(relative string, cause error) (StartersPayload, error) {
		failure := fmt.Errorf("starters.write.failed: %s: %w", relative, cause)
		for i := len(created) - 1; i >= 0; i-- {
			if err := root.Remove(created[i]); err != nil {
				failure = errors.Join(failure, fmt.Errorf("remove created starter %s: %w", created[i], err))
			}
		}
		return payload, failure
	}
	for index, file := range files {
		if payload.Starters[index].Action != "create" {
			continue
		}
		relative := filepath.FromSlash(file.Path)
		if err := ctx.Err(); err != nil {
			return rollback(file.Path, err)
		}
		if err := root.MkdirAll(filepath.Dir(relative), 0o755); err != nil {
			return rollback(file.Path, err)
		}
		output, err := root.OpenFile(relative, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return rollback(file.Path, err)
		}
		created = append(created, relative)
		_, writeErr := output.Write(file.Content)
		closeErr := output.Close()
		if err := errors.Join(writeErr, closeErr); err != nil {
			return rollback(file.Path, err)
		}
	}
	payload.State = "applied"
	return payload, nil
}

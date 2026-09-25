package speccheck

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/spec"
)

type declaredGovernedTouch struct {
	path string
	line int
}

type undeclaredGrant struct {
	granted     bool
	paths       map[string]bool
	regenerated map[string]bool
	location    Location
}

func detectUndeclaredGovernedPaths(
	result *Result,
	specsRoot string,
	repoRoot string,
	graph *spec.Graph,
	artifacts []constraintArtifact,
) error {
	grant, err := readUndeclaredGrant(repoRoot, graph.Spec.Slug, artifacts)
	if err != nil {
		return fmt.Errorf("read undeclared Governed Path authorization: %w", err)
	}
	rows := presentToolingRows(artifacts)
	for _, task := range graph.Tasks {
		if task.Status == spec.StatusCompleted || task.Type == spec.TaskTypeQA {
			continue
		}
		taskPath := filepath.Join(specsRoot, filepath.FromSlash(task.File))
		content, err := os.ReadFile(taskPath)
		if err != nil {
			return fmt.Errorf("read Task %q for undeclared Governed Paths: %w", task.File, err)
		}
		touches, err := declaredGovernedTouches(repoRoot, task, content)
		if err != nil {
			return fmt.Errorf("read Task %q Verification operands: %w", task.File, err)
		}
		for _, touch := range touches {
			if !GovernedPath(touch.path) || grant.regenerated[touch.path] {
				continue
			}
			missingRecord := !grant.granted || !grant.paths[touch.path]
			var missingRows []constraintArtifact
			for _, artifact := range rows {
				row := artifact.rows[strings.ToLower(constraintTooling)]
				if !containsPath(row.BoundedPaths, touch.path) {
					missingRows = append(missingRows, artifact)
				}
			}
			if !missingRecord && len(missingRows) == 0 {
				continue
			}

			taskDisplayPath := artifactDisplayPath(repoRoot, taskPath)
			locations := []Location{{Path: taskDisplayPath, Line: touch.line}}
			var declarations []string
			if missingRecord {
				locations = append(locations, grant.location)
				declarations = append(declarations, grant.location.Path)
			}
			for _, artifact := range missingRows {
				row := artifact.rows[strings.ToLower(constraintTooling)]
				locations = append(locations, Location{Path: artifact.displayPath, Line: row.Line})
				declarations = append(declarations, artifact.displayPath+" Tooling authority row")
			}
			result.Findings = append(result.Findings, Finding{
				Code:     CodeToolingUndeclared,
				Severity: SeverityError,
				Summary:  taskDisplayPath + " declares Governed Path " + touch.path + ", but its tooling authority omits it",
				Where:    locations,
				Fix:      "Add `" + touch.path + "` to " + strings.Join(declarations, ", ") + ".",
			})
		}
	}
	return nil
}

func readUndeclaredGrant(repoRoot, slug string, artifacts []constraintArtifact) (undeclaredGrant, error) {
	grant := undeclaredGrant{
		paths:       make(map[string]bool),
		regenerated: make(map[string]bool),
		location:    Location{Path: filepath.ToSlash(filepath.Join("docs", "specs", slug, "_authorization.md")), Line: 1},
	}
	for _, artifact := range artifacts {
		row, present := artifact.rows[strings.ToLower(constraintTooling)]
		if !present {
			continue
		}
		if !row.Authorization.Selected {
			if grant.location.Path == filepath.ToSlash(filepath.Join("docs", "specs", slug, "_authorization.md")) {
				grant.location = Location{Path: artifact.displayPath, Line: row.Line}
			}
			continue
		}
		reference := row.Authorization.Reference
		grant.location = Location{Path: reference.Path, Line: 1}
		resolution := spec.ReadAuthorization(context.Background(), spec.AuthorizationReadRequest{
			RepoRoot:   reference.ReadRoot,
			RecordPath: reference.ReadPath,
			Role:       authorizationRole(reference.Path),
			AskingSpec: slug,
		})
		if resolution.Outcome != spec.AuthorizationGranted {
			return grant, nil
		}
		grant.granted = true
		for _, path := range resolution.Record.Paths {
			grant.paths[path] = true
		}
		outputs, err := mechanicalRegenerationOutputs(repoRoot, resolution.Record.Regenerations)
		if err != nil {
			return undeclaredGrant{}, err
		}
		grant.regenerated = outputs
		return grant, nil
	}
	return grant, nil
}

func presentToolingRows(artifacts []constraintArtifact) []constraintArtifact {
	var rows []constraintArtifact
	for _, artifact := range artifacts {
		if _, present := artifact.rows[strings.ToLower(constraintTooling)]; present {
			rows = append(rows, artifact)
		}
	}
	return rows
}

func declaredGovernedTouches(repoRoot string, task spec.Task, content []byte) ([]declaredGovernedTouch, error) {
	byPath := make(map[string]declaredGovernedTouch)
	for _, ref := range task.Context {
		if ref.Kind != spec.ContextKindInterface && ref.Kind != spec.ContextKindCreates {
			continue
		}
		byPath[ref.Path] = declaredGovernedTouch{
			path: ref.Path,
			line: sectionLineContaining(content, "Context", ref.Path),
		}
	}
	for _, command := range task.Verification {
		files, err := spec.TaskVerificationFiles(repoRoot, spec.Task{
			Context:      task.Context,
			Verification: []string{command},
		})
		if err != nil {
			return nil, err
		}
		for _, path := range files {
			if _, present := byPath[path]; present {
				continue
			}
			byPath[path] = declaredGovernedTouch{
				path: path,
				line: sectionLineContaining(content, "Verification", command),
			}
		}
	}
	paths := make([]string, 0, len(byPath))
	for path := range byPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	touches := make([]declaredGovernedTouch, 0, len(paths))
	for _, path := range paths {
		touches = append(touches, byPath[path])
	}
	return touches, nil
}

func containsPath(paths []string, target string) bool {
	for _, path := range paths {
		if path == target {
			return true
		}
	}
	return false
}

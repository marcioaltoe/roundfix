package speccheck

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"roundfix/internal/spec"
)

const (
	// CodeOrdinalClaimed identifies an ADR ordinal held by another tree path
	// or claimed for another path by an active Spec Task.
	CodeOrdinalClaimed = "SC-ORDINAL-CLAIMED"
)

type ordinalClaim struct {
	number   string
	path     string
	specSlug string
	taskPath string
}

func detectOrdinalClaims(result *Result, specsRoot, repoRoot string, graph *spec.Graph) error {
	claims := ordinalClaims(graph, specsRoot, repoRoot)
	if len(claims) == 0 {
		return nil
	}
	heldPaths, err := heldADRPaths(repoRoot)
	if err != nil {
		return err
	}
	for index, claim := range claims {
		for _, other := range claims[:index] {
			if claim.number == other.number && claim.path != other.path {
				if ordinalClaimIsHeld(claim, heldPaths) || ordinalClaimIsHeld(other, heldPaths) {
					continue
				}
				result.Findings = append(result.Findings, ordinalSameSpecFinding(claim, other))
			}
		}
	}

	otherClaims, err := activeOrdinalClaims(specsRoot, repoRoot, graph.Spec.Slug)
	if err != nil {
		return err
	}

	for _, claim := range claims {
		fulfilled := false
		for _, heldPath := range heldPaths[claim.number] {
			if heldPath == claim.path {
				fulfilled = true
				continue
			}
			result.Findings = append(result.Findings, ordinalTreeFinding(claim, heldPath))
		}
		if fulfilled {
			continue
		}
		for _, other := range otherClaims[claim.number] {
			result.Findings = append(result.Findings, ordinalSpecFinding(claim, other))
		}
	}
	return nil
}

func ordinalClaimIsHeld(claim ordinalClaim, heldPaths map[string][]string) bool {
	for _, heldPath := range heldPaths[claim.number] {
		if heldPath == claim.path {
			return true
		}
	}
	return false
}

func ordinalClaims(graph *spec.Graph, specsRoot, repoRoot string) []ordinalClaim {
	var claims []ordinalClaim
	for _, task := range graph.Tasks {
		for _, ref := range task.Context {
			if ref.Kind != spec.ContextKindCreates {
				continue
			}
			number, ok := claimedADRNumber(ref.Path)
			if !ok {
				continue
			}
			claims = append(claims, ordinalClaim{
				number:   number,
				path:     ref.Path,
				specSlug: graph.Spec.Slug,
				taskPath: artifactDisplayPath(repoRoot, filepath.Join(specsRoot, filepath.FromSlash(task.File))),
			})
		}
	}
	return claims
}

func claimedADRNumber(claimPath string) (string, bool) {
	if path.Dir(claimPath) != "docs/adr" {
		return "", false
	}
	match := adrFilenamePattern.FindStringSubmatch(path.Base(claimPath))
	if len(match) != 2 {
		return "", false
	}
	return match[1], true
}

func heldADRPaths(repoRoot string) (map[string][]string, error) {
	held := map[string][]string{}
	dir := filepath.Join(repoRoot, "docs", "adr")
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return held, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read ADR directory %q: %w", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		match := adrFilenamePattern.FindStringSubmatch(entry.Name())
		if len(match) != 2 {
			continue
		}
		held[match[1]] = append(held[match[1]], path.Join("docs/adr", entry.Name()))
	}
	return held, nil
}

func activeOrdinalClaims(specsRoot, repoRoot, currentSlug string) (map[string][]ordinalClaim, error) {
	claims := map[string][]ordinalClaim{}
	activeSpecs, err := spec.ListActive(specsRoot)
	if err != nil {
		return nil, fmt.Errorf("list active Specs for ADR ordinal claims: %w", err)
	}
	for _, activeSpec := range activeSpecs {
		if activeSpec.Slug == currentSlug {
			continue
		}
		manifestPath := filepath.Join(activeSpec.Dir, "_tasks.md")
		if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("read Task Graph %q for ADR ordinal claims: %w", manifestPath, err)
		}
		graph, loadErr := spec.LoadForRecovery(specsRoot, activeSpec.Slug)
		if loadErr != nil {
			continue
		}
		for _, claim := range ordinalClaims(graph, specsRoot, repoRoot) {
			claims[claim.number] = append(claims[claim.number], claim)
		}
	}
	return claims, nil
}

func ordinalTreeFinding(claim ordinalClaim, heldPath string) Finding {
	return Finding{
		Code:     CodeOrdinalClaimed,
		Severity: SeverityError,
		Summary: "ADR ordinal " + claim.number + " claimed by " + claim.path +
			" in active Spec " + claim.specSlug + " is already held by " + heldPath,
		Where: []Location{
			{Path: claim.taskPath, Line: 1},
			{Path: heldPath, Line: 1},
		},
		Fix: "Choose an unclaimed ADR ordinal and rename the `creates:` path in " + claim.taskPath + ".",
	}
}

func ordinalSpecFinding(claim, other ordinalClaim) Finding {
	return Finding{
		Code:     CodeOrdinalClaimed,
		Severity: SeverityError,
		Summary: "ADR ordinal " + claim.number + " claimed by " + claim.path +
			" in active Spec " + claim.specSlug + " is also claimed by " + other.path +
			" in active Spec " + other.specSlug,
		Where: []Location{
			{Path: claim.taskPath, Line: 1},
			{Path: other.taskPath, Line: 1},
		},
		Fix: "Give one active Spec an unclaimed ADR ordinal and update its `creates:` path.",
	}
}

func ordinalSameSpecFinding(claim, other ordinalClaim) Finding {
	return Finding{
		Code:     CodeOrdinalClaimed,
		Severity: SeverityError,
		Summary: "ADR ordinal " + claim.number + " is claimed by both " + claim.path +
			" and " + other.path + " in active Spec " + claim.specSlug,
		Where: []Location{
			{Path: claim.taskPath, Line: 1},
			{Path: other.taskPath, Line: 1},
		},
		Fix: "Renumber one Task's `creates:` path to an unclaimed ADR ordinal.",
	}
}

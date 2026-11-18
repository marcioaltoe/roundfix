package baseline

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const ClaudeInstructionLinkPath, ClaudeInstructionLinkTarget = "CLAUDE.md", "AGENTS.md"

// ClaudeInstructionLinkShape describes the root instruction carrier on disk.
type ClaudeInstructionLinkShape string

const (
	ClaudeInstructionLinkLinked           ClaudeInstructionLinkShape = "link"
	ClaudeInstructionLinkAbsent           ClaudeInstructionLinkShape = "absent"
	ClaudeInstructionLinkFile             ClaudeInstructionLinkShape = "file"
	ClaudeInstructionLinkElsewhere        ClaudeInstructionLinkShape = "link-elsewhere"
	ClaudeInstructionLinkUnlinkedCheckout ClaudeInstructionLinkShape = "unlinked-checkout"
	ClaudeInstructionLinkAgentsNotRegular ClaudeInstructionLinkShape = "agents-not-regular"
)

// ClaudeInstructionLink records the carrier's shape and any passages of its own.
type ClaudeInstructionLink struct {
	Shape    ClaudeInstructionLinkShape
	Target   string
	Passages []MisplacedRepositoryRule
}

// Convertible reports whether replacing the carrier with the link loses no rule.
func (link ClaudeInstructionLink) Convertible() bool {
	return link.Shape == ClaudeInstructionLinkAbsent ||
		(link.Shape == ClaudeInstructionLinkFile || link.Shape == ClaudeInstructionLinkElsewhere) && len(link.Passages) == 0
}

// ReadClaudeInstructionLink reads the root carrier without changing files or
// taking locks. It consults the Git index only for the bytes of an unlinked checkout.
func ReadClaudeInstructionLink(ctx context.Context, root string, git GitRunner) (ClaudeInstructionLink, error) {
	var link ClaudeInstructionLink
	if err := ctx.Err(); err != nil {
		return link, fmt.Errorf("read %s: %w", ClaudeInstructionLinkPath, err)
	}
	info, err := os.Lstat(filepath.Join(root, ClaudeInstructionLinkPath))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return link, fmt.Errorf("inspect %s: %w", ClaudeInstructionLinkPath, err)
	}
	missing := errors.Is(err, fs.ErrNotExist)
	target := ""
	relative := ClaudeInstructionLinkPath
	if !missing && info.Mode()&fs.ModeSymlink != 0 {
		target, err = os.Readlink(filepath.Join(root, ClaudeInstructionLinkPath))
		if err != nil {
			return link, fmt.Errorf("read %s target: %w", ClaudeInstructionLinkPath, err)
		}
		anchored, err := os.OpenRoot(root)
		if err != nil {
			return link, fmt.Errorf("open repository for %s: %w", ClaudeInstructionLinkPath, err)
		}
		defer anchored.Close()
		builder := inventoryBuilder{root: anchored, preimages: make(map[string]Preimage)}
		relative, err = builder.resolveAlias(ClaudeInstructionLinkPath, target)
		if err != nil {
			return link, fmt.Errorf("resolve %s: %w", ClaudeInstructionLinkPath, err)
		}
		if path.Clean(target) == ClaudeInstructionLinkTarget {
			return ClaudeInstructionLink{Shape: ClaudeInstructionLinkLinked, Target: target}, nil
		}
	}
	agentsInfo, err := os.Lstat(filepath.Join(root, ClaudeInstructionLinkTarget))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return link, fmt.Errorf("inspect %s for %s: %w", ClaudeInstructionLinkTarget, ClaudeInstructionLinkPath, err)
	}
	if errors.Is(err, fs.ErrNotExist) || !agentsInfo.Mode().IsRegular() {
		return ClaudeInstructionLink{Shape: ClaudeInstructionLinkAgentsNotRegular}, nil
	}
	if missing {
		return ClaudeInstructionLink{Shape: ClaudeInstructionLinkAbsent}, nil
	}
	shape := ClaudeInstructionLinkFile
	if target != "" {
		shape = ClaudeInstructionLinkElsewhere
	} else if !info.Mode().IsRegular() {
		return link, fmt.Errorf("inspect %s: expected a regular file or safe symbolic link", ClaudeInstructionLinkPath)
	}
	content, err := readOptionalRegular(root, relative)
	if err != nil {
		return link, fmt.Errorf("read %s: %w", ClaudeInstructionLinkPath, err)
	}
	if content == nil {
		return link, fmt.Errorf("read %s target %q: %w", ClaudeInstructionLinkPath, relative, fs.ErrNotExist)
	}
	link = ClaudeInstructionLink{Shape: shape, Target: target}
	if shape == ClaudeInstructionLinkFile && string(content) == ClaudeInstructionLinkTarget {
		if git == nil {
			git = ExecGitRunner{}
		}
		output, err := git.RunGit(ctx, root, "ls-files", "-s", "--", ClaudeInstructionLinkPath)
		if err != nil {
			return ClaudeInstructionLink{}, fmt.Errorf("read %s Git index mode: %w", ClaudeInstructionLinkPath, err)
		}
		for _, entry := range strings.Split(output, "\n") {
			fields := strings.Fields(entry)
			if len(fields) == 4 && fields[0] == "120000" && fields[3] == ClaudeInstructionLinkPath {
				link.Shape = ClaudeInstructionLinkUnlinkedCheckout
				break
			}
		}
		return link, nil
	}
	agents, err := readOptionalRegular(root, ClaudeInstructionLinkTarget)
	if err != nil {
		return ClaudeInstructionLink{}, fmt.Errorf("read %s for %s: %w", ClaudeInstructionLinkTarget, ClaudeInstructionLinkPath, err)
	}
	if bytes.Equal(content, agents) {
		return link, nil
	}
	var artifacts map[string]ManifestArtifact
	if shape == ClaudeInstructionLinkFile {
		manifestBytes, err := readOptionalRegular(root, manifestPath)
		if err != nil {
			return ClaudeInstructionLink{}, fmt.Errorf("read %s Setup Manifest: %w", ClaudeInstructionLinkPath, err)
		}
		if manifestBytes != nil {
			manifest, valid := parseManagedSetupManifest(manifestBytes)
			if !valid {
				return ClaudeInstructionLink{}, fmt.Errorf("read %s Setup Manifest: baseline.preservation.managed-refresh.manifest.invalid", ClaudeInstructionLinkPath)
			}
			artifacts = manifestArtifactsByPath(manifest.ManagedArtifacts)[ClaudeInstructionLinkPath]
		}
	}
	entries := partitionRootSource(relative, content)
	managedCounts := make(map[string]int)
	for _, entry := range entries {
		if entry.Kind == "managed-block" {
			id, _ := entry.StructuralProvenance["managedId"].(string)
			managedCounts[id]++
		}
	}
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return ClaudeInstructionLink{}, fmt.Errorf("read %s passages: %w", ClaudeInstructionLinkPath, err)
		}
		if entry.Kind == "unmarked-span" {
			// The import exception belongs to what CLAUDE.md reads, even when
			// its passages are reported at a different target path.
			for _, passage := range misplacedUnmarkedPassages(ClaudeInstructionLinkPath, content, entry) {
				passage.Path = relative
				link.Passages = append(link.Passages, passage)
			}
			continue
		}
		id, _ := entry.StructuralProvenance["managedId"].(string)
		if artifact, recorded := artifacts[id]; recorded && managedCounts[id] == 1 && !managedEntryMatchesManifest(entry, artifact) {
			link.Passages = append(link.Passages, misplacedManagedPassage(relative, content, entry))
		}
	}
	return link, nil
}

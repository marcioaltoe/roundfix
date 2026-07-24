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
	"regexp"
	"strings"
)

var assetsSyncGitHubRemote = regexp.MustCompile(
	`^(?:https://github\.com/|git@github\.com:|ssh://git@github\.com/)([A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+?)(?:\.git)?$`,
)

func inspectAssetsSyncCheckout(
	ctx context.Context,
	sourceDir string,
) (assetsSyncCheckout, error) {
	rootOutput, err := runRestoreGit(ctx, "-C", sourceDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return assetsSyncCheckout{}, assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			filepath.ToSlash(sourceDir),
			"setups",
			fmt.Sprintf("Canonical setup source is not a readable Git checkout: %v.", err),
			"Fix the canonical setup source before synchronizing snapshots.",
			err,
		)
	}
	root, err := filepath.Abs(strings.TrimSpace(string(rootOutput)))
	if err != nil {
		return assetsSyncCheckout{}, err
	}
	if resolved, resolveErr := filepath.EvalSymlinks(root); resolveErr == nil {
		root = resolved
	}
	revisionOutput, err := runRestoreGit(ctx, "-C", root, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return assetsSyncCheckout{}, assetsSyncCheckoutError(sourceDir, err)
	}
	revision := strings.TrimSpace(string(revisionOutput))
	if !immutableGitID.MatchString(revision) {
		return assetsSyncCheckout{}, assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			filepath.ToSlash(sourceDir),
			"setups",
			"Canonical setup source did not resolve to a full immutable commit.",
			"Fix the canonical setup source before synchronizing snapshots.",
			errors.New("source revision is not immutable"),
		)
	}
	status, err := runRestoreGit(ctx, "-C", root, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return assetsSyncCheckout{}, assetsSyncCheckoutError(sourceDir, err)
	}
	if len(status) != 0 {
		return assetsSyncCheckout{}, assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			filepath.ToSlash(sourceDir),
			"setups",
			"Canonical setup Git checkout contains dirty or untracked source bytes.",
			"Fix the canonical setup source before synchronizing snapshots.",
			errors.New("source checkout is dirty"),
		)
	}
	remote, err := runRestoreGit(ctx, "-C", root, "remote", "get-url", "origin")
	if err != nil {
		return assetsSyncCheckout{}, assetsSyncCheckoutError(sourceDir, err)
	}
	match := assetsSyncGitHubRemote.FindStringSubmatch(strings.TrimSpace(string(remote)))
	if len(match) != 2 {
		return assetsSyncCheckout{}, assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			filepath.ToSlash(sourceDir),
			"setups",
			"Canonical setup Git origin is not a portable GitHub repository identity.",
			"Fix the canonical setup source before synchronizing snapshots.",
			errors.New("source origin is not a GitHub repository"),
		)
	}
	return assetsSyncCheckout{
		root:       filepath.Clean(root),
		repository: match[1],
		revision:   revision,
	}, nil
}

func assetsSyncCheckoutError(sourceDir string, err error) error {
	return assetsSyncError(
		AssetsSyncInvalid,
		"skills.setup-snapshot.drift",
		"error",
		filepath.ToSlash(sourceDir),
		"setups",
		fmt.Sprintf("Canonical setup source is not a readable Git checkout: %v.", err),
		"Fix the canonical setup source before synchronizing snapshots.",
		err,
	)
}

func verifyAssetsSyncCommittedFile(
	ctx context.Context,
	checkout assetsSyncCheckout,
	sourcePath string,
) error {
	relative, err := filepath.Rel(checkout.root, sourcePath)
	if err != nil || !safeRelative(filepath.ToSlash(relative)) {
		return errors.New("Canonical setup source file is outside the Git checkout.")
	}
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("Canonical setup source file cannot be read: %w", err)
	}
	committed, err := runRestoreGit(
		ctx,
		"-C",
		checkout.root,
		"show",
		checkout.revision+":"+filepath.ToSlash(relative),
	)
	if err != nil {
		return fmt.Errorf("Canonical setup source file is not committed: %w", err)
	}
	if !bytes.Equal(data, committed) {
		return errors.New("Canonical setup source file bytes do not match the declared commit.")
	}
	return nil
}

func validateAssetsSyncDeclaredSource(
	ctx context.Context,
	raw any,
	checkout assetsSyncCheckout,
	sourcePath string,
) string {
	source, ok := raw.(map[string]any)
	if !ok || len(source) == 0 {
		return "Canonical setup skill has empty external provenance."
	}
	if len(source) != 4 {
		return "Canonical setup skill provenance fields are incomplete or machine-local."
	}
	for _, key := range []string{"type", "repository", "ref", "path"} {
		if _, ok := source[key]; !ok {
			return "Canonical setup skill provenance fields are incomplete or machine-local."
		}
	}
	if source["type"] != "github" {
		return "Canonical setup skill provider must be github."
	}
	if source["repository"] != checkout.repository {
		return "Canonical setup skill repository does not match the Git origin."
	}
	revision, ok := source["ref"].(string)
	if !ok || !immutableGitID.MatchString(revision) {
		return "Canonical setup skill ref must be a full immutable commit."
	}
	resolved, err := runRestoreGit(
		ctx,
		"-C",
		checkout.root,
		"rev-parse",
		"--verify",
		revision+"^{commit}",
	)
	if err != nil {
		return "Canonical setup skill ref is not available in the Git checkout."
	}
	if strings.TrimSpace(string(resolved)) != revision {
		return "Canonical setup skill ref did not resolve to the declared full commit."
	}
	declaredPath, ok := source["path"].(string)
	if !ok || !safeRelative(declaredPath) {
		return "Canonical setup skill source path is not portable."
	}
	if declaredPath != sourcePath {
		return "Canonical setup skill source path does not match its setup path."
	}
	return ""
}

func assetsSyncFilesystemTreeDigest(root string) (string, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return "", fmt.Errorf("cannot inspect skill tree %s: %w", root, err)
	}
	if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 {
		return "", fmt.Errorf("skill tree root is not a regular directory: %s", root)
	}
	files := []restoreFile{}
	err = filepath.WalkDir(root, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if current == root {
			return nil
		}
		relative, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		info, err := os.Lstat(current)
		if err != nil {
			return fmt.Errorf("cannot inspect skill tree entry %s: %w", relative, err)
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return fmt.Errorf("skill tree entry is a symbolic link: %s", relative)
		}
		if info.IsDir() {
			if entry.Name() == ".git" || entry.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("skill tree entry is not a regular file: %s", relative)
		}
		if fileLinkCount(info) != 1 {
			return fmt.Errorf("skill tree entry is hard linked: %s", relative)
		}
		content, err := os.ReadFile(current)
		if err != nil {
			return fmt.Errorf("cannot read skill tree file %s: %w", relative, err)
		}
		files = append(files, restoreFile{Path: relative, Content: content})
		return nil
	})
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("skill tree has no regular files: %s", root)
	}
	return portableRestoreDigest(files), nil
}

func assetsSyncCommittedTreeDigest(
	ctx context.Context,
	checkoutRoot, revision, sourcePath string,
) (string, error) {
	output, err := runRestoreGit(
		ctx,
		"-C",
		checkoutRoot,
		"ls-tree",
		"-r",
		"-z",
		"--full-tree",
		revision,
		"--",
		sourcePath,
	)
	if err != nil {
		return "", err
	}
	prefix := []byte(sourcePath + "/")
	files := []restoreFile{}
	for _, raw := range bytes.Split(output, []byte{0}) {
		if len(raw) == 0 {
			continue
		}
		header, rawPath, ok := bytes.Cut(raw, []byte{'\t'})
		parts := bytes.Fields(header)
		if !ok || len(parts) != 3 || !bytes.HasPrefix(rawPath, prefix) {
			return "", fmt.Errorf("Git tree entry escapes source path %s", sourcePath)
		}
		relative := string(rawPath[len(prefix):])
		if assetsSyncIgnoredTreePath(relative) {
			continue
		}
		if !safeRestoreRelative(relative) || path.Clean(relative) != relative {
			return "", fmt.Errorf("Git tree entry has unsafe path %q", relative)
		}
		if string(parts[1]) != "blob" ||
			(string(parts[0]) != "100644" && string(parts[0]) != "100755") {
			return "", fmt.Errorf("Git tree entry is not a regular file: %s", relative)
		}
		content, err := runRestoreGit(
			ctx,
			"-C",
			checkoutRoot,
			"cat-file",
			"blob",
			string(parts[2]),
		)
		if err != nil {
			return "", err
		}
		files = append(files, restoreFile{Path: relative, Content: content})
	}
	if len(files) == 0 {
		return "", fmt.Errorf("Git commit has no regular files under %s", sourcePath)
	}
	return portableRestoreDigest(files), nil
}

func assetsSyncIgnoredTreePath(relative string) bool {
	parts := strings.Split(relative, "/")
	for _, part := range parts[:len(parts)-1] {
		if part == ".git" || part == "node_modules" {
			return true
		}
	}
	return false
}

package cli

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"gopkg.in/yaml.v3"

	"roundfix/internal/preflight"
	"roundfix/internal/spec"
)

const deliveryConventionsVersion = "roundfix/delivery-conventions/v2"

type deliveryConvention struct {
	ID     string
	Prompt string
}

func deliveryConventions() []deliveryConvention {
	return []deliveryConvention{
		{"C1", "A Spec's QA Report records the head it audited and is committed after that head, so it never names the commit that records it."},
		{"C2", "The Daemon writes a Task file's status and its `## Result`, `## Recorded paths` and `## Carry-forward provenance` sections after the Task's Verification passes; a Result that calls status Daemon-owned agrees with a `completed` status."},
		{"C3", "The archive commit adds an Archive Record at the archive root and removes the Spec folder, whose bytes stay in Git at source_revision; legacy archives move the folder and stamp its front matter."},
		{"C4", "A planning candidate authors a Spec whose Tasks are all pending and which has no QA Report; that Spec's own delivery implements it and is reviewed then."},
		{"C5", "A Spec archived through the QA Archive Override records `qa_override: true`, `qa_override_approval` and `qa_override_reason` in its archived `_prd.md` front matter; its QA Task keeps its observed status and its QA Report its observed verdict. For an Archive Record, these fields are in the record and the unchanged QA Task and Report stay in Git at source_revision."},
	}
}

// reviewRepository reads only the immutable candidate tree, never checkout files.
type reviewRepository struct {
	Root, Head, Base string
	SpecRoots        []string
	Git              preflight.GitRunner
}

func conventionRegions(ctx context.Context, repo reviewRepository, anchor reviewFindingAnchor) ([]string, error) {
	var rules []string
	if anchor.StartLine < 1 || anchor.EndLine < anchor.StartLine || path.Clean(anchor.Path) != anchor.Path {
		return rules, nil
	}
	for index, root := range repo.SpecRoots {
		root = strings.TrimSuffix(root, "/")
		if !strings.HasPrefix(anchor.Path, root+"/") {
			continue
		}
		if index > 0 {
			rules = append(rules, "C3")
		}
		rel := strings.TrimPrefix(anchor.Path, root+"/")
		parts := strings.Split(rel, "/")
		if len(parts) == 1 && index > 0 && strings.HasSuffix(parts[0], ".md") {
			override, err := reviewSpecQAOverride(ctx, repo, root+"/"+strings.TrimSuffix(parts[0], ".md"))
			if err != nil {
				return nil, err
			}
			if override {
				rules = append(rules, "C5")
			}
		}
		if len(parts) < 2 {
			continue
		}
		if index > 0 {
			override, err := reviewSpecQAOverride(ctx, repo, root+"/"+parts[0])
			if err != nil {
				return nil, err
			}
			if override {
				rules = append(rules, "C5")
			}
		}
		if len(parts) > 2 && parts[1] == "qa" {
			rules = append(rules, "C1")
		}
		if len(parts) != 2 {
			continue
		}
		isTask := strings.HasPrefix(parts[1], "task_") && strings.HasSuffix(parts[1], ".md")
		var lines []string
		if isTask {
			body, err := repo.Git.RunGit(ctx, repo.Root, "show", repo.Head+":"+anchor.Path)
			if err != nil {
				return nil, fmt.Errorf("read convention Task %s: %w", anchor.Path, err)
			}
			lines = reviewFileLines(body)
			if reviewTaskSettlementRegion(lines, anchor) {
				rules = append(rules, "C2")
			}
		}
		if index != 0 || (!isTask && parts[1] != "_tasks.md") {
			continue
		}
		pending, err := reviewSpecAllPending(ctx, repo, root+"/"+parts[0])
		if err != nil {
			return nil, err
		}
		if pending && (parts[1] == "_tasks.md" || reviewRangeInside(anchor, 1, reviewFrontMatterEnd(lines))) {
			rules = append(rules, "C4")
		}
	}
	return rules, nil
}

func reviewSpecQAOverride(ctx context.Context, repo reviewRepository, specPath string) (bool, error) {
	archived, err := readArchivedSpecAt(ctx, repo.Git, repo.Root, repo.Head, path.Dir(specPath), path.Base(specPath))
	if err == nil && archived.Form == spec.ArchivedRecord {
		override := archived.Record.QAOverride
		return override != nil && strings.TrimSpace(override.Approval) != "" && strings.TrimSpace(override.Reason) != "", nil
	}
	if err != nil && !errors.Is(err, spec.ErrNotArchived) {
		return false, err
	}
	prdPath := specPath + "/_prd.md"
	files, err := repo.Git.RunGit(ctx, repo.Root, "ls-tree", "--name-only", repo.Head, "--", prdPath)
	if err != nil {
		return false, fmt.Errorf("inspect override PRD %s: %w", prdPath, err)
	}
	if strings.TrimSpace(files) == "" {
		return false, nil
	}
	body, err := repo.Git.RunGit(ctx, repo.Root, "show", repo.Head+":"+prdPath)
	if err != nil {
		return false, fmt.Errorf("read override PRD %s: %w", prdPath, err)
	}
	lines := reviewFileLines(body)
	end := reviewFrontMatterEnd(lines)
	if end == 0 {
		return false, nil
	}
	var record struct {
		Override bool   `yaml:"qa_override"`
		Approval string `yaml:"qa_override_approval"`
		Reason   string `yaml:"qa_override_reason"`
	}
	if err := yaml.Unmarshal([]byte(strings.Join(lines[1:end-1], "\n")), &record); err != nil {
		return false, fmt.Errorf("parse override PRD %s: %w", prdPath, err)
	}
	return record.Override && strings.TrimSpace(record.Approval) != "" && strings.TrimSpace(record.Reason) != "", nil
}

func reviewFileLines(body string) []string {
	return strings.Split(strings.TrimSuffix(body, "\n"), "\n")
}

func reviewRangeInside(anchor reviewFindingAnchor, start, end int) bool {
	return start > 0 && anchor.StartLine >= start && anchor.EndLine <= end
}

func reviewFrontMatterEnd(lines []string) int {
	if len(lines) == 0 || lines[0] != "---" {
		return 0
	}
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			return i + 1
		}
	}
	return 0
}

func reviewTaskSettlementRegion(lines []string, anchor reviewFindingAnchor) bool {
	if anchor.StartLine < 1 || anchor.EndLine > len(lines) {
		return false
	}
	frontMatterEnd := reviewFrontMatterEnd(lines)
	settlement := false
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			settlement = line == "## Result" || line == "## Recorded paths" || line == "## Carry-forward provenance"
		}
		if i+1 >= anchor.StartLine && i+1 <= anchor.EndLine && i+1 > frontMatterEnd && !settlement {
			return false
		}
	}
	return true
}

func reviewSpecAllPending(ctx context.Context, repo reviewRepository, slug string) (bool, error) {
	files, err := repo.Git.RunGit(ctx, repo.Root, "ls-tree", "-r", "--name-only", "-z", repo.Head, "--", slug+"/")
	if err != nil {
		return false, fmt.Errorf("list planning Spec: %w", err)
	}
	for _, file := range strings.Split(files, "\x00") {
		if strings.HasPrefix(file, slug+"/qa/") {
			return false, nil
		}
		if path.Dir(file) != slug || !strings.HasPrefix(path.Base(file), "task_") || !strings.HasSuffix(file, ".md") {
			continue
		}
		body, err := repo.Git.RunGit(ctx, repo.Root, "show", repo.Head+":"+file)
		if err != nil {
			return false, fmt.Errorf("read planning Task %s: %w", file, err)
		}
		lines := reviewFileLines(body)
		end := reviewFrontMatterEnd(lines)
		if end == 0 {
			return false, nil
		}
		statuses := 0
		for _, line := range lines[1 : end-1] {
			if strings.HasPrefix(line, "status:") {
				statuses++
				if strings.TrimSpace(strings.TrimPrefix(line, "status:")) != "pending" {
					return false, nil
				}
			}
		}
		if statuses != 1 {
			return false, nil
		}
	}
	return true, nil
}

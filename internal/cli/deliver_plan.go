package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const deliverPlanSchema = "roundfix-deliver-plan/v1"

var deliveryOperations = []spec.AuthorizationOperation{
	spec.AuthorizationOperationImplement,
	spec.AuthorizationOperationCommit,
	spec.AuthorizationOperationPush,
	spec.AuthorizationOperationPullRequest,
	spec.AuthorizationOperationMerge,
}

type deliverPlanDocument struct {
	Schema string                  `json:"schema"`
	Specs  []deliverPlanSpec       `json:"specs"`
	Intent []deliverPlanIntentItem `json:"intent"`
}

type deliverPlanSpec struct {
	Slug           string                     `json:"slug"`
	Verdict        string                     `json:"verdict"`
	Tasks          int                        `json:"tasks"`
	Unfinished     int                        `json:"unfinishedTasks"`
	Reasons        []string                   `json:"reasons"`
	SharedPremises []deliverPlanSharedPremise `json:"sharedPremises"`
	premises       []string
}

type deliverPlanSharedPremise struct {
	With  string   `json:"with"`
	Paths []string `json:"paths"`
}

type deliverPlanIntentItem struct {
	Kind   string `json:"kind"`
	Path   string `json:"path"`
	Status string `json:"status"`
}

func strictSpecFindings(specsRoot, repoRoot, specSlug string) ([]speccheck.Finding, error) {
	checked, err := speccheck.Check(specsRoot, repoRoot, specSlug)
	if err != nil {
		return nil, err
	}
	speccheck.PromoteGaps(&checked)
	return speccheck.GatePrecondition(checked).Findings, nil
}

func productionPremises(graph *spec.Graph) []string {
	if graph == nil {
		return nil
	}
	unique := make(map[string]struct{})
	for _, task := range graph.Tasks {
		if task.Type == spec.TaskTypeQA {
			continue
		}
		for _, ref := range task.Context {
			if ref.Kind != spec.ContextKindInterface || !strings.HasSuffix(ref.Path, ".go") || strings.HasSuffix(ref.Path, "_test.go") {
				continue
			}
			unique[ref.Path] = struct{}{}
		}
	}
	paths := make([]string, 0, len(unique))
	for path := range unique {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func deliveryAuthorizationReasons(
	ctx context.Context,
	loaded roundconfig.Loaded,
	specsRoot roundconfig.SpecsRoot,
	slug string,
) []string {
	resolution := spec.ReadSpecAuthorization(ctx, loaded.GitRoot, specsRoot.Path, slug, "")
	if resolution.Outcome != spec.AuthorizationGranted {
		return []string{fmt.Sprintf("authorization %s: %s", resolution.Outcome, resolution.Reason.Code)}
	}
	missing := make([]string, 0, len(deliveryOperations))
	for _, operation := range deliveryOperations {
		if !resolution.Permits(operation) {
			missing = append(missing, string(operation))
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return []string{"authorization lacks " + strings.Join(missing, ", ")}
}

func runDeliverPlan(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	request, err := parseDeliverPlan(args)
	if err != nil {
		return printDeliverFailure("plan", err, stderr)
	}
	loaded, specsRoot, err := loadDeliveryCommand(ctx, environment, stderr)
	if err != nil {
		return printDeliverFailure("plan", err, stderr)
	}

	slugs := request.slugs
	if len(slugs) == 0 {
		active, skipped, listErr := spec.ListActiveDetailed(specsRoot.Path)
		if listErr != nil {
			return printDeliverFailure("plan", listErr, stderr)
		}
		printSkippedSpecDiagnostics(stderr, skipped)
		slugs = make([]string, 0, len(active))
		for _, activeSpec := range active {
			slugs = append(slugs, activeSpec.Slug)
		}
	}

	document := deliverPlanDocument{
		Schema: deliverPlanSchema,
		Specs:  make([]deliverPlanSpec, 0, len(slugs)),
		Intent: []deliverPlanIntentItem{},
	}
	blocked := false
	for _, slug := range slugs {
		graph, loadErr := spec.Load(specsRoot.Path, slug)
		if loadErr != nil {
			return printDeliverFailure("plan", loadErr, stderr)
		}
		planned := deliverPlanSpec{
			Slug:           slug,
			Verdict:        "approved",
			Tasks:          len(graph.Tasks),
			Reasons:        []string{},
			SharedPremises: []deliverPlanSharedPremise{},
			premises:       productionPremises(graph),
		}
		for _, task := range graph.Tasks {
			if task.Status != spec.StatusCompleted {
				planned.Unfinished++
			}
		}
		planned.Reasons = append(planned.Reasons, deliveryAuthorizationReasons(ctx, loaded, specsRoot, slug)...)
		findings, checkErr := strictSpecFindings(specsRoot.Path, loaded.GitRoot, slug)
		if checkErr != nil {
			return printDeliverFailure("plan", checkErr, stderr)
		}
		if codes := findingCodes(findings); len(codes) > 0 {
			planned.Reasons = append(planned.Reasons, "spec check: "+strings.Join(codes, ", "))
		}
		for _, earlier := range document.Specs {
			if shared := intersectSortedPaths(planned.premises, earlier.premises); len(shared) > 0 {
				planned.SharedPremises = append(planned.SharedPremises, deliverPlanSharedPremise{
					With:  earlier.Slug,
					Paths: shared,
				})
			}
		}
		if len(planned.Reasons) > 0 {
			planned.Verdict = "blocked"
			blocked = true
		}
		document.Specs = append(document.Specs, planned)
	}

	document.Intent, err = deliveryPlanIntent(loaded.GitRoot)
	if err != nil {
		return printDeliverFailure("plan", err, stderr)
	}
	if request.json {
		encoder := json.NewEncoder(stdout)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(document); err != nil {
			return printDeliverFailure("plan", fmt.Errorf("write JSON plan: %w", err), stderr)
		}
	} else {
		printDeliverPlanText(stdout, document)
	}
	if blocked {
		return exitRunFailed
	}
	return exitOK
}

type deliverPlanRequest struct {
	json  bool
	slugs []string
}

func parseDeliverPlan(args []string) (deliverPlanRequest, error) {
	request := deliverPlanRequest{}
	fs := flagSet("deliver plan")
	fs.BoolVar(&request.json, "json", false, "Print one machine-readable plan")
	if err := fs.Parse(hoistCommandFlags(args, nil)); err != nil {
		return request, validationError{message: err.Error()}
	}
	request.slugs = fs.Args()
	for index := range request.slugs {
		request.slugs[index] = strings.TrimSpace(request.slugs[index])
		if request.slugs[index] == "" {
			return deliverPlanRequest{}, validationError{message: "Spec slug cannot be empty"}
		}
	}
	return request, nil
}

func findingCodes(findings []speccheck.Finding) []string {
	unique := make(map[string]struct{}, len(findings))
	for _, finding := range findings {
		code := strings.TrimSpace(finding.Code)
		if code != "" {
			unique[code] = struct{}{}
		}
	}
	codes := make([]string, 0, len(unique))
	for code := range unique {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}

func intersectSortedPaths(first, second []string) []string {
	shared := make([]string, 0)
	left, right := 0, 0
	for left < len(first) && right < len(second) {
		switch {
		case first[left] < second[right]:
			left++
		case first[left] > second[right]:
			right++
		default:
			shared = append(shared, first[left])
			left++
			right++
		}
	}
	return shared
}

func printDeliverPlanText(output io.Writer, document deliverPlanDocument) {
	for _, planned := range document.Specs {
		reasons := "-"
		if len(planned.Reasons) > 0 {
			reasons = strings.Join(planned.Reasons, "; ")
		}
		fmt.Fprintf(output, "spec\t%s\t%s\t%d/%d\t%s\n", planned.Slug, planned.Verdict, planned.Unfinished, planned.Tasks, reasons)
		for _, shared := range planned.SharedPremises {
			fmt.Fprintf(output, "shared\t%s\t%s\t%s\n", planned.Slug, shared.With, strings.Join(shared.Paths, ", "))
		}
	}
	for _, intent := range document.Intent {
		fmt.Fprintf(output, "%s\t%s\t%s\n", intent.Kind, intent.Path, intent.Status)
	}
}

func deliveryPlanIntent(repoRoot string) ([]deliverPlanIntentItem, error) {
	var intent []deliverPlanIntentItem
	for _, source := range []struct {
		kind string
		dir  string
	}{
		{kind: "backlog", dir: "docs/backlog"},
		{kind: "finding", dir: "docs/findings"},
	} {
		paths, err := filepath.Glob(filepath.Join(repoRoot, filepath.FromSlash(source.dir), "*.md"))
		if err != nil {
			return nil, fmt.Errorf("list %s intent: %w", source.kind, err)
		}
		for _, path := range paths {
			status, err := readIntentStatus(path)
			if err != nil {
				return nil, err
			}
			intent = append(intent, deliverPlanIntentItem{
				Kind:   source.kind,
				Path:   deliverPlanRelativePath(repoRoot, path),
				Status: status,
			})
		}
	}
	inboxRoot := filepath.Join(repoRoot, "docs", "_inbox")
	err := filepath.WalkDir(inboxRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type().IsRegular() {
			intent = append(intent, deliverPlanIntentItem{
				Kind:   "inbox",
				Path:   deliverPlanRelativePath(repoRoot, path),
				Status: "-",
			})
		}
		return nil
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("list inbox intent: %w", err)
	}
	kindOrder := map[string]int{"backlog": 0, "finding": 1, "inbox": 2}
	sort.Slice(intent, func(left, right int) bool {
		if kindOrder[intent[left].Kind] != kindOrder[intent[right].Kind] {
			return kindOrder[intent[left].Kind] < kindOrder[intent[right].Kind]
		}
		return intent[left].Path < intent[right].Path
	})
	if intent == nil {
		return []deliverPlanIntentItem{}, nil
	}
	return intent, nil
}

func readIntentStatus(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("read intent %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return "", fmt.Errorf("read intent %q: missing frontmatter", path)
	}
	var frontmatter strings.Builder
	closed := false
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "---" {
			closed = true
			break
		}
		frontmatter.WriteString(scanner.Text())
		frontmatter.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read intent %q: %w", path, err)
	}
	if !closed {
		return "", fmt.Errorf("read intent %q: unterminated frontmatter", path)
	}
	var fields struct {
		Status string `yaml:"status"`
	}
	if err := yaml.Unmarshal([]byte(frontmatter.String()), &fields); err != nil {
		return "", fmt.Errorf("parse intent %q frontmatter: %w", path, err)
	}
	fields.Status = strings.TrimSpace(fields.Status)
	if fields.Status == "" {
		return "", fmt.Errorf("read intent %q: frontmatter has no status", path)
	}
	return fields.Status, nil
}

func deliverPlanRelativePath(repoRoot, path string) string {
	relative, err := filepath.Rel(repoRoot, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relative)
}

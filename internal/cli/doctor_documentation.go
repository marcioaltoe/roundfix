package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"roundfix/internal/baseline"
	"roundfix/internal/spec"
)

const (
	doctorExportNext = "run roundfix baseline update, then edit the repository's own lines"
	doctorDocsNext   = "see each line; roundfix baseline update fills an Amendment List and removes a retired ADR"
)

// doctorDocumentationResults reads repository documentation only. Both checks
// are advisory, so neither can change the exit determined by other checks.
func doctorDocumentationResults(root string) []CheckResult {
	results := []CheckResult{
		{Name: HealthCheckBaselineExport, Status: CheckStatusSkipped, Detail: "no Setup Manifest at docs/agents/setup-context.json"},
		{Name: HealthCheckBaselineDocs, Status: CheckStatusSkipped, Detail: "no Setup Manifest at docs/agents/setup-context.json"},
	}
	if root == "" {
		return results
	}
	if _, err := os.Stat(filepath.Join(root, baselineSetupManifestPath)); errors.Is(err, os.ErrNotExist) {
		return results
	}
	facts, err := gatherDocumentationFacts(root)
	if errors.Is(err, baseline.ErrNoManifest) {
		return results
	}
	if err != nil {
		results[0] = doctorDocumentationUnread(results[0].Name, err)
		results[1] = doctorDocumentationUnread(results[1].Name, err)
		return results
	}
	anchored, err := os.OpenRoot(root)
	if err != nil {
		results[0] = doctorDocumentationUnread(results[0].Name, err)
		results[1] = doctorDocumentationUnread(results[1].Name, err)
		return results
	}
	defer anchored.Close()
	if slices.Contains(facts.Modules, "secondbrain") {
		results[0] = doctorExportResult(anchored, facts)
	} else {
		results[0].Detail = "the Secondbrain is turned off in the Setup Manifest"
	}
	results[1] = doctorDocsResult(anchored, facts)
	return results
}

func doctorDocumentationUnread(name string, err error) CheckResult {
	return CheckResult{Name: name, Status: CheckStatusWarn, Detail: "DR-EXPORT-UNREAD: " + err.Error(), NextAction: "re-run roundfix doctor after fixing the Setup Manifest"}
}

func doctorDocumentationReport(name, ok, rules, next string, departures []baseline.DocumentationDeparture) CheckResult {
	result := CheckResult{Name: name, Status: CheckStatusOK, Detail: ok}
	if len(departures) == 0 {
		return result
	}
	result.Status = CheckStatusWarn
	result.Detail = fmt.Sprintf("%d departure(s) from the %s rules", len(departures), rules)
	result.NextAction = next
	counts := map[string]int{}
	for _, departure := range departures {
		counts[departure.Code]++
	}
	printed := map[string]int{}
	for _, departure := range departures {
		code := departure.Code
		printed[code]++
		if printed[code] <= 20 {
			result.Lines = append(result.Lines, code+": "+departure.Detail)
		}
		if printed[code] == 20 && counts[code] > 20 {
			result.Lines = append(result.Lines, fmt.Sprintf("%s: %d more", code, counts[code]-20))
		}
	}
	return result
}

func doctorExportResult(root *os.Root, facts documentationFacts) CheckResult {
	current, err := root.ReadFile(baseline.ExportListPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return doctorDocumentationUnread(HealthCheckBaselineExport, fmt.Errorf("read export list: %w", err))
	}
	plan := baseline.PlanExportList(current, err == nil, facts.Export)
	var departures []baseline.DocumentationDeparture
	add := func(code, detail string) {
		departures = append(departures, baseline.DocumentationDeparture{Code: code, Detail: detail})
	}
	switch {
	case plan.Action == "refused":
		add("DR-EXPORT-MARKER", plan.Refusal+"; restore the pair of markers or delete both")
	case !slices.ContainsFunc(strings.Split(string(current), "\n"), func(line string) bool { return strings.TrimSpace(line) == baseline.ExportListBeginMarker }):
		add("DR-EXPORT-LIST-MISSING", ".secondbrain-export has no Standard Export List")
	case plan.Action == "update":
		add("DR-EXPORT-LIST-STALE", fmt.Sprintf("the Standard Export List differs from what roundfix baseline update would write: %d line(s) to add, %d to remove", len(plan.Added), len(plan.Removed)))
	}
	for _, d := range baseline.ExportVariationDepartures(current, facts.Export) {
		detail := fmt.Sprintf(".secondbrain-export:%d: %s", d.Line, d.Path)
		if d.Code == "DR-EXPORT-PROCESS-PATH" {
			detail += " publishes process; remove the line"
		} else {
			detail += " has no reason; add a comment line directly above it"
		}
		add(d.Code, detail)
	}
	total, unmarked := 0, 0
	for _, adr := range facts.ADRs {
		if adr.Retirement.Retired || doctorADRProposed(adr.Content) {
			continue
		}
		total++
		if adr.Header.Knowledge == nil {
			unmarked++
		}
	}
	if unmarked > 0 {
		add("DR-EXPORT-ADR-UNMARKED", fmt.Sprintf("%d of %d ADR(s) carry no Knowledge Mark and are not exported", unmarked, total))
	}
	// Count entries outside the standard block, including exclusions.
	own, inStandard := 0, false
	for _, line := range strings.Split(string(current), "\n") {
		line = strings.TrimSpace(line)
		switch line {
		case baseline.ExportListBeginMarker:
			inStandard = true
		case baseline.ExportListEndMarker:
			inStandard = false
		default:
			if !inStandard && line != "" && !strings.HasPrefix(line, "#") {
				own++
			}
		}
	}
	return doctorDocumentationReport(HealthCheckBaselineExport, fmt.Sprintf("the Standard Export List is current; %d line(s) of the repository's own", own), "export", doctorExportNext, departures)
}

var doctorLegacyProposed = regexp.MustCompile(`(?im)^\s*(?:\*\*)?status(?:\*\*)?:\s*proposed\b`)

// ReadADRHeader owns frontmatter and ClassifyADR owns retirement. The legacy
// status reader keeps a proposed ADR without frontmatter advisory-free too.
func doctorADRProposed(content []byte) bool {
	if header := spec.ReadADRHeader(content); header.Status != "" {
		return header.Status == "proposed"
	}
	lines := strings.SplitN(string(content), "\n", 21)
	if len(lines) > 20 {
		lines = lines[:20]
	}
	return doctorLegacyProposed.MatchString(strings.Join(lines, "\n"))
}

func doctorDocsResult(root *os.Root, facts documentationFacts) CheckResult {
	f := baseline.DocsLayoutFacts{DocsEntries: facts.Export.DocsEntries, SpecRoot: facts.SpecRoot, History: facts.HistoryRoot, GlossaryFiles: facts.Export.GlossaryFiles}
	entries, err := fs.ReadDir(root.FS(), ".")
	if err != nil {
		return doctorDocumentationUnread(HealthCheckBaselineDocs, fmt.Errorf("read root documents: %w", err))
	}
	for _, entry := range entries {
		if entry.Name() == "README.md" && !entry.IsDir() {
			f.ReadmePresent = true
		}
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			f.RootDocuments = append(f.RootDocuments, entry.Name())
		}
	}
	f.SpecificRepository, err = root.ReadFile(baseline.SpecificRepositoryPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return doctorDocumentationUnread(HealthCheckBaselineDocs, fmt.Errorf("read repository rules: %w", err))
	}
	departures := baseline.DocsLayoutDepartures(f)
	for i, d := range departures {
		if d.Code == "DR-DOCS-OUTSIDE-LAYOUT" {
			departures[i].Detail = d.Path + ": " + d.Detail + " the docs-layout guide does not name; move its content or name it in docs/agents/specific-repository.md"
		}
	}
	for _, adr := range facts.ADRs {
		if adr.Retirement.Retired {
			departures = append(departures, baseline.DocumentationDeparture{Code: "DR-DOCS-ADR-RETIRED", Detail: fmt.Sprintf("%s: a %s ADR is still in the tree", adr.Path, adr.Retirement.Reason)})
			continue
		}
		if doctorADRProposed(adr.Content) {
			continue
		}
		var missing []string
		for _, id := range documentationNotices(adr) {
			if !slices.Contains(adr.Header.AmendedBy, id) {
				missing = append(missing, id)
			}
		}
		if len(missing) > 0 {
			departures = append(departures, baseline.DocumentationDeparture{Code: "DR-DOCS-ADR-AMENDMENT-UNLISTED", Detail: adr.Path + ": the body names " + strings.Join(missing, ", ") + " and amended_by does not"})
		}
	}
	err = fs.WalkDir(root.FS(), "docs/lessons", func(path string, entry fs.DirEntry, walkErr error) error {
		if path == "docs/lessons" && errors.Is(walkErr, os.ErrNotExist) {
			return nil
		}
		if walkErr != nil {
			return fmt.Errorf("read lessons %s: %w", path, walkErr)
		}
		if entry.IsDir() {
			return nil
		}
		content, err := root.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read lesson %s: %w", path, err)
		}
		for _, problem := range spec.LessonLearnedProblems(filepath.Base(path), content) {
			departures = append(departures, baseline.DocumentationDeparture{Code: "DR-DOCS-LESSON-FORM", Detail: path + ": " + problem})
		}
		return nil
	})
	if err != nil {
		return doctorDocumentationUnread(HealthCheckBaselineDocs, err)
	}
	return doctorDocumentationReport(HealthCheckBaselineDocs, "the documentation follows the docs-layout guide", "documentation", doctorDocsNext, departures)
}

package releaseplan

import (
	"path"
	"strings"
)

// ClassifyCommit extracts the supported Conventional Commit evidence subset
// and maintenance-only path boundary for one normalized commit.
func ClassifyCommit(commit Commit) ChangeEvidence {
	conventionalType, subjectBreaking := parseConventionalSubject(commit.Subject)
	breaking := subjectBreaking || hasBreakingFooter(commit.Body)
	impact := ImpactNone
	switch {
	case breaking:
		impact = ImpactMajor
	case conventionalType == "feat":
		impact = ImpactMinor
	case conventionalType == "fix" || conventionalType == "perf":
		impact = ImpactPatch
	}

	return ChangeEvidence{
		CommitSHA:                      commit.SHA,
		Subject:                        commit.Subject,
		ConventionalType:               conventionalType,
		Breaking:                       breaking,
		AutomaticImpact:                impact,
		CrossesMaintenanceOnlyBoundary: crossesMaintenanceOnlyBoundary(commit.ChangedPaths),
	}
}

// ClassifyChanges aggregates commit evidence, validates any manual
// classification, and returns the conservative maximum required impact.
func ClassifyChanges(request ClassifyRequest) (ClassificationResult, error) {
	changes := make([]ChangeEvidence, 0, len(request.Commits))
	sources := map[ClassificationSource]bool{}
	automaticMinimum := ImpactNone
	breaking := false
	var blockingCommits []string

	for _, commit := range request.Commits {
		evidence := ClassifyCommit(commit)
		changes = append(changes, evidence)

		switch {
		case !evidence.CrossesMaintenanceOnlyBoundary:
			sources[SourceMaintenanceOnly] = true
		case evidence.AutomaticImpact != ImpactNone:
			automaticMinimum = MaxImpact(automaticMinimum, evidence.AutomaticImpact)
			breaking = breaking || evidence.Breaking
			sources[SourceConventionalCommit] = true
		default:
			blockingCommits = append(blockingCommits, evidence.CommitSHA)
		}
	}

	manualProvided := request.ManualImpact != ""
	if !manualProvided && strings.TrimSpace(request.ManualReason) != "" {
		return ClassificationResult{}, ManualImpactError{
			Reason:     request.ManualReason,
			NextAction: "provide --impact with --reason",
			Err:        ErrManualImpactRequired,
		}
	}
	if manualProvided {
		if err := ValidateManualImpact(request.ManualImpact, request.ManualReason, automaticMinimum); err != nil {
			return ClassificationResult{}, err
		}
		sources[SourceManual] = true
	}

	manualRequired := len(blockingCommits) > 0 && !manualProvided
	impact := automaticMinimum
	if manualProvided {
		impact = MaxImpact(impact, request.ManualImpact)
	}
	source := classificationSource(sources)
	if manualRequired {
		source = SourceMixed
	}
	classification := Classification{
		Source:                       source,
		Impact:                       impact,
		Breaking:                     breaking,
		ManualReason:                 strings.TrimSpace(request.ManualReason),
		ManualClassificationRequired: manualRequired,
		BlockingCommits:              append([]string(nil), blockingCommits...),
	}
	if manualRequired {
		classification.Impact = automaticMinimum
		return ClassificationResult{
			State:          StateManualClassificationRequired,
			Classification: classification,
			Changes:        changes,
		}, nil
	}

	state := StateReady
	if impact == ImpactNone {
		state = StateNoRelease
	}
	return ClassificationResult{
		State:          state,
		Classification: classification,
		Changes:        changes,
	}, nil
}

// ValidateManualImpact checks a user-supplied classification before it can
// participate in maximum-impact selection.
func ValidateManualImpact(impact Impact, reason string, minimum Impact) error {
	if !AllowedImpact(impact) {
		return UnknownImpactError{Impact: impact}
	}
	if strings.TrimSpace(reason) == "" {
		return ManualImpactError{
			Impact:     impact,
			Minimum:    minimum,
			NextAction: "provide a non-empty --reason explaining the manual classification",
			Err:        ErrManualReasonRequired,
		}
	}
	if CompareImpact(impact, minimum) < 0 {
		return ManualImpactError{
			Impact:     impact,
			Minimum:    minimum,
			Reason:     reason,
			NextAction: "choose an impact that is at least the automatic minimum",
			Err:        ErrManualImpactTooLow,
		}
	}
	return nil
}

// MaintenanceOnly reports whether a changed path is inside the no-release
// maintenance boundary documented by the Release Plan TechSpec.
func MaintenanceOnly(changedPath string) bool {
	clean := cleanChangedPath(changedPath)
	if clean == "" {
		return false
	}
	if isTestOrFixturePath(clean) {
		return true
	}
	if isPlanningEvidencePath(clean) {
		return true
	}
	if isNonReleaseCIPath(clean) {
		return true
	}
	return false
}

func parseConventionalSubject(subject string) (string, bool) {
	header, _, ok := strings.Cut(subject, ":")
	if !ok {
		return "", false
	}
	header = strings.TrimSpace(header)
	if header == "" {
		return "", false
	}
	breaking := strings.Contains(header, "!")
	typeEnd := len(header)
	for index, char := range header {
		if char == '(' || char == '!' {
			typeEnd = index
			break
		}
	}
	conventionalType := header[:typeEnd]
	if !validConventionalType(conventionalType) {
		return "", false
	}
	return conventionalType, breaking
}

func validConventionalType(conventionalType string) bool {
	if conventionalType == "" {
		return false
	}
	for _, char := range conventionalType {
		if char < 'a' || char > 'z' {
			return false
		}
	}
	return true
}

func hasBreakingFooter(body string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "BREAKING CHANGE:") {
			return true
		}
	}
	return false
}

func crossesMaintenanceOnlyBoundary(paths []string) bool {
	if len(paths) == 0 {
		return true
	}
	for _, changedPath := range paths {
		if !MaintenanceOnly(changedPath) {
			return true
		}
	}
	return false
}

func cleanChangedPath(changedPath string) string {
	clean := path.Clean(strings.TrimSpace(strings.ReplaceAll(changedPath, "\\", "/")))
	if clean == "." || strings.HasPrefix(clean, "../") || clean == ".." || strings.HasPrefix(clean, "/") {
		return ""
	}
	return clean
}

func isPlanningEvidencePath(changedPath string) bool {
	for _, prefix := range []string{
		"docs/specs/",
		"docs/adr/",
		"docs/findings/",
		"docs/handoffs/",
	} {
		if strings.HasPrefix(changedPath, prefix) {
			return true
		}
	}
	return false
}

func isTestOrFixturePath(changedPath string) bool {
	fileName := path.Base(changedPath)
	if strings.HasSuffix(fileName, "_test.go") ||
		strings.HasSuffix(fileName, "_test.py") ||
		strings.HasSuffix(fileName, ".test.js") ||
		strings.HasSuffix(fileName, ".test.ts") ||
		strings.HasSuffix(fileName, ".spec.js") ||
		strings.HasSuffix(fileName, ".spec.ts") ||
		strings.HasSuffix(fileName, ".golden") {
		return true
	}
	for _, segment := range strings.Split(changedPath, "/") {
		switch segment {
		case "testdata", "tests", "fixtures", "fixture":
			return true
		}
	}
	return false
}

func isNonReleaseCIPath(changedPath string) bool {
	if !strings.HasPrefix(changedPath, ".github/workflows/") {
		return false
	}
	fileName := path.Base(changedPath)
	return fileName != "release.yml" && fileName != "release.yaml"
}

func classificationSource(sources map[ClassificationSource]bool) ClassificationSource {
	if len(sources) == 0 {
		return SourceMaintenanceOnly
	}
	if len(sources) > 1 {
		return SourceMixed
	}
	for source := range sources {
		return source
	}
	return SourceMaintenanceOnly
}

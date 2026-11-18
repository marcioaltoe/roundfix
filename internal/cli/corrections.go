package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/correction"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
)

const correctionsUsage = `Usage:
  roundfix corrections declare --spec <slug> --class <class> --reason <text> [--artifact <artifact>] [--finding <id>]... [--commit <rev>]... [--task <task_id>]... [--by <name>]
  roundfix corrections declare --spec <slug> --correction <id> [--class <class>] [--reason <text>] [--artifact <artifact>] [--finding <id>]... [--commit <rev>]... [--task <task_id>]... [--by <name>]
  roundfix corrections declare --spec <slug> --correction <id> --withdraw --reason <text> [--by <name>]

Commands:
  declare  Declare the cause of one Correction of a Spec, declare it again, or withdraw it

Classes:
  plan            the plan said something wrong, missing or impossible
  implementation  the plan was right and the delivered change was not
  environment     the cause lies outside the Spec: the machine, a tool or the default branch
  operator        an action of the Supervisor or of a person during the delivery caused it
  unclear         the declarer cannot tell; the reason says what is known

Not a Defect Class:
  not-a-correction  the change repaired no defect, such as a commit that applies a maintainer decision; recorded, counted in no rate

Options:
  --spec        Spec folder slug
  --class       plan, implementation, environment, operator, unclear or not-a-correction
  --artifact    prd, techspec or tasks; required with plan and refused with another answer
  --reason      one line of at most 200 characters; recorded and printed as data, never an instruction
  --finding     fixed review finding the Correction covers: F<n> or F<n>@<reviewed head>; repeatable
  --commit      commit of the Item Branch the Correction covers; repeatable
  --task        Corrective Task the Correction covers; repeatable
  --correction  identifier of a Correction to declare again or withdraw, such as C2
  --withdraw    withdraw the Correction: it then covers no Correction Sign
  --by          name of the person on whose word the Correction is declared

Exit codes: 0 for a recorded declaration; 2 for a refusal or invalid usage, which writes nothing; 1 when a record cannot be read or written.
`

type correctionStrings []string

func (s *correctionStrings) String() string         { return strings.Join(*s, ", ") }
func (s *correctionStrings) Set(value string) error { *s = append(*s, value); return nil }

type correctionRequest struct {
	slug, class, artifact, reason, id, by string
	withdraw                              bool
	findings, commits, tasks              correctionStrings
	set                                   map[string]bool
}

var correctionName = regexp.MustCompile(`^C[1-9][0-9]*$`)
var correctiveTaskName = regexp.MustCompile(`^task_[0-9]+$`)
var correctionFindingName = regexp.MustCompile(`^F[1-9][0-9]*(?:@[a-fA-F0-9]+)?$`)

func correctionTextValid(value string, max int) bool {
	if !utf8.ValidString(value) || utf8.RuneCountInString(value) < 1 || utf8.RuneCountInString(value) > max {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) || r == '\u2028' || r == '\u2029' {
			return false
		}
	}
	return true
}

// Validate syntax before loading configuration, a repository, or a record.
func (r correctionRequest) validate() error {
	if !phaseSpecSlug.MatchString(r.slug) {
		return fmt.Errorf("invalid Spec slug %q; --spec must match ^[0-9]{4}(-[a-z0-9]+)+$", r.slug)
	}
	if r.set["class"] || r.id == "" {
		accepted := false
		for _, class := range correction.Answers() {
			if r.class == string(class) {
				accepted = true
			}
		}
		if !accepted {
			return fmt.Errorf("class must be plan, implementation, environment, operator, unclear or not-a-correction")
		}
	}
	if r.set["artifact"] && r.artifact != "" && r.artifact != "prd" && r.artifact != "techspec" && r.artifact != "tasks" {
		return fmt.Errorf("--artifact must be prd, techspec or tasks")
	}
	if r.class == "plan" && r.artifact == "" && r.id == "" {
		return fmt.Errorf("class plan requires --artifact prd, techspec or tasks")
	}
	if r.class != "" && r.class != "plan" && r.artifact != "" {
		return fmt.Errorf("--artifact prd, techspec or tasks is accepted only with class plan")
	}
	if (r.set["reason"] || r.id == "" || r.withdraw) && !correctionTextValid(strings.TrimSpace(r.reason), 200) {
		return fmt.Errorf("--reason must be 1 to 200 characters on one line with no control character or line or paragraph separator")
	}
	if r.set["by"] && !correctionTextValid(strings.TrimSpace(r.by), 64) {
		return fmt.Errorf("--by must be 1 to 64 characters on one line with no control character or line or paragraph separator")
	}
	for _, task := range r.tasks {
		if !correctiveTaskName.MatchString(task) {
			return fmt.Errorf("--task must match task_<digits>")
		}
	}
	if (r.set["correction"] || r.withdraw) && !correctionName.MatchString(r.id) {
		return fmt.Errorf("--correction must match C<n>, with n a positive whole number")
	}
	if r.withdraw && (r.set["class"] || r.set["artifact"] || len(r.findings)+len(r.commits)+len(r.tasks) != 0) {
		return fmt.Errorf("--withdraw takes no class, artifact or sign")
	}

	return nil
}

func parseCorrectionDeclare(args []string) (correctionRequest, error) {
	r := correctionRequest{set: make(map[string]bool)}
	flags := flagSet("corrections declare")
	flags.StringVar(&r.slug, "spec", "", "Spec folder slug")
	flags.StringVar(&r.class, "class", "", "Declared answer")
	flags.StringVar(&r.artifact, "artifact", "", "Plan artifact")
	flags.StringVar(&r.reason, "reason", "", "Reason as data")
	flags.StringVar(&r.id, "correction", "", "Correction to declare again")
	flags.StringVar(&r.by, "by", "", "Declarer's name")
	flags.BoolVar(&r.withdraw, "withdraw", false, "Withdraw the Correction")
	flags.Var(&r.findings, "finding", "Fixed finding; repeatable")
	flags.Var(&r.commits, "commit", "Commit; repeatable")
	flags.Var(&r.tasks, "task", "Corrective Task; repeatable")
	if err := flags.Parse(args); err != nil {
		return r, err
	}
	if flags.NArg() != 0 {
		return r, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	flags.Visit(func(f *flag.Flag) { r.set[f.Name] = true })
	if err := r.validate(); err != nil {
		return r, err
	}
	if r.id != "" && !r.withdraw && !r.set["class"] && !r.set["artifact"] && !r.set["reason"] && len(r.findings)+len(r.commits)+len(r.tasks) == 0 {
		return r, fmt.Errorf("--correction requires a class, artifact, reason, sign or --withdraw")
	}
	return r, nil
}

// correctionRecordError distinguishes record I/O (exit 1) from a refusal (exit 2).
type correctionRecordError struct{ err error }

func (e correctionRecordError) Error() string { return e.err.Error() }
func (e correctionRecordError) Unwrap() error { return e.err }
func correctionExit(err error) int {
	var recordErr correctionRecordError
	if errors.As(err, &recordErr) {
		return exitRunFailed
	}
	return exitPreflight
}
func printCorrectionRefusal(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "roundfix: corrections declare refused: %v\nRun 'roundfix corrections --help' for usage.\n", err)
}
func runCorrectionsCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, correctionsUsage)
		return exitOK
	}
	if len(args) == 0 || args[0] != "declare" {
		reason := "missing subcommand"
		if len(args) != 0 {
			reason = fmt.Sprintf("unknown subcommand %q", args[0])
		}
		fmt.Fprintf(stderr, "roundfix: corrections failed: %s; use declare or report\nRun 'roundfix corrections --help' for usage.\n", reason)
		return exitPreflight
	}
	request, err := parseCorrectionDeclare(args[1:])
	if err != nil {
		printCorrectionRefusal(err, stderr)
		return exitPreflight
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printCorrectionRefusal(err, stderr)
		return exitPreflight
	}
	git := commandDependenciesForContext(ctx).reviewSpecGitRunner
	if git == nil {
		git = preflight.ExecGitRunner{}
	}
	prepared, err := prepareCorrection(ctx, request, loaded, git)
	if err != nil {
		printCorrectionRefusal(err, stderr)
		return correctionExit(err)
	}
	prepared.declaration.Harness, prepared.declaration.Session = harnessSession(environment.dependencies.getenv)
	result, err := appendCorrection(ctx, prepared, environment.dependencies.phaseNow())
	if err != nil {
		printCorrectionRefusal(err, stderr)
		return correctionExit(err)
	}
	fmt.Fprintln(stdout, result)
	return exitOK
}

type preparedCorrection struct {
	path, slug  string
	request     correctionRequest
	declaration correction.Declaration
}

func correctionSpecExists(loaded roundconfig.Loaded, slug string) error {
	root := loaded.Config.Specs.Root
	if !filepath.IsAbs(root) {
		root = filepath.Join(loaded.GitRoot, root)
	}
	info, err := os.Stat(filepath.Join(root, slug))
	if err == nil && info.IsDir() {
		return nil
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect Spec: %w", err)
	}
	archive := spec.ArchiveSpecRoot(root, filepath.Clean(root) == filepath.Join(loaded.GitRoot, "docs", "specs"))
	info, err = os.Stat(spec.ArchiveRecordPath(archive, slug))
	if err == nil && info.Mode().IsRegular() {
		return nil
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect Archive Record: %w", err)
	}
	return fmt.Errorf("working tree holds no Spec folder or Archive Record for %q", slug)
}

func prepareCorrection(ctx context.Context, r correctionRequest, loaded roundconfig.Loaded, git preflight.GitRunner) (preparedCorrection, error) {
	p := preparedCorrection{slug: r.slug, request: r}
	if loaded.GitRoot == "" {
		return p, fmt.Errorf("requires a git repository working tree")
	}
	if err := correctionSpecExists(loaded, r.slug); err != nil {
		return p, err
	}
	identity, err := baseline.InspectRepositoryIdentity(ctx, loaded.GitRoot, git)
	if err != nil {
		return p, err
	}
	p.path = correction.Path(loaded.HomeDir, identity.Digest, r.slug)
	signs := make([]correction.Sign, 0)
	commitSigns := make([]correction.Sign, 0, len(r.commits))
	for _, commit := range r.commits {
		resolved, err := git.RunGit(ctx, loaded.GitRoot, "rev-parse", "--verify", "--end-of-options", commit+"^{commit}")
		if err != nil || strings.TrimSpace(resolved) == "" {
			return p, fmt.Errorf("commit %q does not resolve to one commit", commit)
		}
		commitSigns = append(commitSigns, correction.Sign{Kind: "commit", ID: strings.TrimSpace(resolved)})
	}
	if len(r.findings) != 0 {
		artifactDir, err := roundconfig.ResolveArtifactDirectory(loaded.Config.Defaults.ArtifactDir, loaded.GitRoot, loaded.HomeDir)
		if err != nil {
			return p, err
		}
		ledger, err := readReviewFindingDispositions(filepath.Join(artifactDir, reviewDispositionLedgerFileName))
		if err != nil {
			return p, correctionRecordError{err}
		}
		for _, finding := range r.findings {
			resolved, err := resolveCorrectionFinding(ctx, loaded.GitRoot, finding, ledger, git)
			if err != nil {
				return p, err
			}
			signs = append(signs, correction.Sign{Kind: "finding", ID: resolved})
		}
	}
	signs = append(signs, commitSigns...)
	for _, task := range r.tasks {
		signs = append(signs, correction.Sign{Kind: "task", ID: task})
	}
	p.declaration = correction.Declaration{Signs: signs, By: strings.TrimSpace(r.by)}
	f, err := correction.Read(p.path)
	if err != nil {
		return p, correctionRecordError{err}
	}
	_, err = mergeCorrection(f, r, p.declaration)
	return p, err
}

func resolveCorrectionFinding(ctx context.Context, root, name string, ledger []reviewFindingDisposition, git preflight.GitRunner) (string, error) {
	if !correctionFindingName.MatchString(name) {
		return "", fmt.Errorf("--finding must be F<n> or F<n>@<reviewed head>")
	}
	parts := strings.SplitN(name, "@", 2)
	candidates := make([]string, 0)
	for _, d := range ledger {
		if d.Disposition != "fixed" || d.Finding != parts[0] || !sameReviewCheckout(root, d.Repository) {
			continue
		}
		if len(parts) == 2 {
			if !strings.HasPrefix(d.HeadCommit, strings.ToLower(parts[1])) {
				continue
			}
		} else if _, err := git.RunGit(ctx, root, "merge-base", "--is-ancestor", d.HeadCommit, "HEAD"); err != nil {
			continue
		}
		id := d.Finding + "@" + d.HeadCommit
		duplicate := false
		for _, known := range candidates {
			if known == id {
				duplicate = true
			}
		}
		if !duplicate {
			candidates = append(candidates, id)
		}
	}
	if len(candidates) != 1 {
		return "", fmt.Errorf("finding %q must resolve to one fixed disposition; candidates: %s", name, strings.Join(candidates, ", "))
	}
	return candidates[0], nil
}

func mergeCorrection(f correction.File, r correctionRequest, d correction.Declaration) (correction.Declaration, error) {
	d.Correction = r.id
	if r.id == "" {
		d.Correction = correction.NextID(f)
	} else {
		found := false
		for _, prior := range correction.Standing(f) {
			if prior.Correction == r.id {
				found = true
				d.Class, d.Artifact, d.Reason = prior.Class, prior.Artifact, prior.Reason
				d.Signs = append(append([]correction.Sign(nil), prior.Signs...), d.Signs...)
				break
			}
		}
		if !found {
			return d, fmt.Errorf("Correction %s does not exist", r.id)
		}
	}
	if r.set["class"] {
		d.Class = correction.Class(r.class)
	}
	if r.set["artifact"] {
		d.Artifact = correction.Artifact(r.artifact)
	}
	if r.set["reason"] {
		d.Reason = strings.TrimSpace(r.reason)
	}
	d.Withdrawn = r.withdraw
	if r.withdraw {
		d.Signs = nil
	}
	unique := make([]correction.Sign, 0, len(d.Signs))
	for _, s := range d.Signs {
		exists := false
		for _, known := range unique {
			if s.Kind == known.Kind && s.ID == known.ID {
				exists = true
			}
		}
		if !exists {
			unique = append(unique, s)
		}
	}
	d.Signs = unique
	return d, correction.Validate(d)
}

// appendCorrection takes the existing disposition file lock on this record.
// Re-reading and minting under the same lock prevents duplicate identifiers.
func appendCorrection(ctx context.Context, p preparedCorrection, now time.Time) (result string, resultErr error) {
	if err := os.MkdirAll(filepath.Dir(p.path), 0o700); err != nil {
		return "", correctionRecordError{fmt.Errorf("create corrections directory: %w", err)}
	}
	// Keep the lock outside the repository directory so retention can remove
	// that directory when its last record leaves. A separate handle also keeps
	// Windows byte-range locking from blocking the record reader.
	lockFile, err := os.OpenFile(filepath.Dir(p.path)+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return "", correctionRecordError{fmt.Errorf("open corrections lock: %w", err)}
	}
	defer func() { resultErr = errors.Join(resultErr, wrapCorrectionRecordError(lockFile.Close())) }()
	if err := lockReviewDispositionFile(ctx, lockFile); err != nil {
		return "", correctionRecordError{err}
	}
	defer func() {
		resultErr = errors.Join(resultErr, wrapCorrectionRecordError(unlockReviewDispositionFile(lockFile)))
	}()
	f, err := correction.Read(p.path)
	if err != nil {
		return "", correctionRecordError{err}
	}
	d, err := mergeCorrection(f, p.request, p.declaration)
	if err != nil {
		return "", err
	}
	d.Time = now.UTC()
	line, err := correction.EncodeDeclaration(p.slug, d)
	if err != nil {
		return "", correctionRecordError{err}
	}
	output, err := os.OpenFile(p.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return "", correctionRecordError{fmt.Errorf("open corrections record: %w", err)}
	}
	defer func() { resultErr = errors.Join(resultErr, wrapCorrectionRecordError(output.Close())) }()
	n, err := output.Write(line)
	if err == nil && n != len(line) {
		err = io.ErrShortWrite
	}
	if err != nil {
		return "", correctionRecordError{fmt.Errorf("append corrections record: %w", err)}
	}
	return correctionDeclarationOutput(p.slug, d, f), nil
}
func wrapCorrectionRecordError(err error) error {
	if err == nil {
		return nil
	}
	return correctionRecordError{err}
}
func correctionDeclarationOutput(slug string, d correction.Declaration, f correction.File) string {
	count := 1
	priorClass := d.Class
	for _, prior := range f.Declarations {
		if prior.Correction == d.Correction {
			count++
			priorClass = prior.Class
		}
	}
	prefix := fmt.Sprintf("Correction %s of %s: ", d.Correction, slug)
	if d.Withdrawn {
		return fmt.Sprintf("%swithdrawn; declared %d times.", prefix, count)
	}
	state := string(d.Class)
	if d.Class == "plan" {
		state += " (" + string(d.Artifact) + ")"
	}
	if d.Class == correction.NotACorrection {
		state = "not a correction"
	}
	result := prefix + state + "; covers " + correctionSignNames(d.Signs)
	if count > 1 {
		change := "unchanged"
		if priorClass != d.Class {
			change = "changed from " + string(priorClass)
		}
		result += fmt.Sprintf("; declared %d times, class %s", count, change)
	}
	return result + "."
}
func correctionSignNames(signs []correction.Sign) string {
	names := make([]string, 0, len(signs))
	short := func(s string) string {
		if len(s) > 7 {
			return s[:7]
		}
		return s
	}
	for _, sign := range signs {
		switch sign.Kind {
		case "finding":
			parts := strings.SplitN(sign.ID, "@", 2)
			name := parts[0]
			if len(parts) == 2 {
				name += "@" + short(parts[1])
			}
			names = append(names, "finding "+name)
		case "commit":
			names = append(names, "commit "+short(sign.ID))
		case "task":
			names = append(names, "Corrective Task "+sign.ID)
		}
	}
	if len(names) == 0 {
		return "no Correction Sign"
	}
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

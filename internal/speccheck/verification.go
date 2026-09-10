package speccheck

import (
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

const (
	// CodeVerifyWorkIndependent identifies a Verification that cannot distinguish Task work from no work.
	CodeVerifyWorkIndependent = "SC-VERIFY-WORK-INDEPENDENT"
	// CodeVerifyInvertedExit identifies a Verification whose shell status reverses or ignores its asserted condition.
	CodeVerifyInvertedExit = "SC-VERIFY-INVERTED-EXIT"
	// CodeVerifyNonHermetic identifies a Verification that depends on state outside the repository.
	CodeVerifyNonHermetic = "SC-VERIFY-NON-HERMETIC"
	// CodeVerifyVacuousCommand identifies one Verification command that already
	// passes against the unchanged tree, whatever its siblings prove.
	CodeVerifyVacuousCommand = "SC-VERIFY-VACUOUS-COMMAND"
)

var (
	makeVerifyPattern = regexp.MustCompile(`^(?:rtk\s+)?(?:env\s+(?:\S+=\S+\s+)*)?make\s+verify$`)
	goWideGatePattern = regexp.MustCompile(`\bgo\s+(?:build|test|vet)\b`)
	// An assertion over which paths changed.
	workingTreeStatePattern = regexp.MustCompile(`git\s+(?:status|diff)\b[^|&;]*(?:--porcelain|--short|--check|--quiet|--exit-code|--name-only|--name-status|--stat)`)
	// A terminal predicate that succeeds on empty input, which is what an
	// unchanged tree produces. `git diff --name-only | grep -q .` reads the
	// same paths and its terminal `grep -q .` fails on empty output, so it is
	// not vacuous; `| cat` succeeds on empty output and is. Matching the git
	// flags alone is never enough: `--exit-code` on the git command does not
	// change what the final predicate does when it consumes no input.
	//
	// Every alternative is anchored to the whole terminal command so that
	// text inside another command's quoted argument is never taken for a
	// success form: `grep -q 'exit 0'` is a grep that fails on empty input,
	// not an `exit 0` success predicate.
	//
	// A `test -z` / `[ -z` empty-string test is a guaranteed success only when
	// the whole length operand is a command substitution that reads the working
	// tree at test time, so an unchanged tree yields an empty value. A
	// substitution of arbitrary output (`$(printf x)`), a single-quoted form
	// (which the shell treats as literal text, not a substitution), a bare
	// variable, or a literal operand is not proof: the predicate can fail
	// rather than pass. Those forms fall through to the emptyInputUnknown
	// branch and are never reported vacuous.
	emptyOutputSucceedsPattern = regexp.MustCompile(
		`^(?:rtk\s+)?(?:cat|true|:)\s*$` + // | cat | true | : succeed on empty input
			`|^(?:rtk\s+)?exit\s+0\s*$`, // | exit 0 explicitly succeeds
	)
	// emptyTestPattern and emptyTestBracketPattern capture the single length
	// operand of an empty-string test so it can be judged by
	// emptyTestOperandPasses rather than matched greedily off the `$(` opener.
	emptyTestPattern        = regexp.MustCompile(`^(?:rtk\s+)?test\s+-z\s+(.+)$`)
	emptyTestBracketPattern = regexp.MustCompile(`^(?:rtk\s+)?\[\s*-z\s+(.+?)\s*\]\s*$`)
	// A terminal command that provably fails on empty input. grep finds no
	// lines in an unchanged tree's empty output and exits nonzero, so it is
	// honest Verification no matter what text sits in its arguments.
	emptyOutputFailsPattern = regexp.MustCompile(`^(?:rtk\s+)?grep\b`)
	// These patterns deliberately cover only the reversed forms measured in
	// authored Verification. They match command positions, not quoted examples,
	// and leave general shell correctness to command execution.
	grepCountExitPattern = regexp.MustCompile(
		`(?:^|(?:&&|\|\||;)\s*)(?:rtk\s+)?grep\s+` +
			`(?:(?:-[a-z]+|--[a-z-]+)\s+)*(?:-[a-z]*c[a-z]*|--count)(?:\s|$)`,
	)
	filteredCountExitPattern = regexp.MustCompile(
		`(?:^|(?:&&|\|\||;)\s*)(?:rtk\s+)?grep\s+` +
			`(?:(?:-[a-z]+|--[a-z-]+)\s+)*(?:-[a-z]*v[a-z]*|--invert-match)(?:\s|$)` +
			`[^;&]*\|\s*(?:rtk\s+)?wc\s+` +
			`(?:(?:-[a-z]+|--[a-z-]+)\s+)*(?:-[a-z]*l[a-z]*|--lines)(?:\s|$)`,
	)
	bareTestSubstitutionPattern = regexp.MustCompile(
		`(?:^|(?:&&|\|\||;)\s*)(?:rtk\s+)?test\s+` +
			`(?:"\$\([^)]*\)"|'\$\([^)]*\)'|\$\([^)]*\))\s*(?:$|&&|\|\||;)`,
	)
	testComparisonPattern = regexp.MustCompile(
		`(?:^|(?:&&|\|\||;)\s*)(?:rtk\s+)?test\b[^;&|]*` +
			`(?:-(?:eq|ne|gt|ge|lt|le)\b|(?:^|\s)(?:=|==|!=)(?:\s|$))`,
	)
	shellAssignmentPattern  = regexp.MustCompile(`(?:^|[;&]\s*|\s)([A-Za-z_][A-Za-z0-9_]*)=`)
	shellReadPattern        = regexp.MustCompile(`(?:^|[;&|{]\s*)read(?:\s+-[A-Za-z]+)*\s+([A-Za-z_][A-Za-z0-9_]*)`)
	shellForPattern         = regexp.MustCompile(`(?:^|[;&|{]\s*)for\s+([A-Za-z_][A-Za-z0-9_]*)\s+in\b`)
	environmentGuardPattern = regexp.MustCompile(
		`(?:^|(?:&&|\|\||;)\s*)(?:rtk\s+)?(?:test\s+-n|\[\s*-n)\s+` +
			`"?\$(?:\{([A-Za-z_][A-Za-z0-9_]*)[^}]*\}|([A-Za-z_][A-Za-z0-9_]*))`,
	)
	// commandChainOpPattern locates the top-level short-circuit and sequence
	// operators of a shell chain. A single `|` pipe is deliberately excluded:
	// a pipeline's exit status is decided by its last member, and the
	// previous members are counted by pipelineExitOutcome.
	commandChainOpPattern = regexp.MustCompile(`\|\||&&|;`)
)

// WorkIndependentVerification reports a Task whose declared Verification
// cannot distinguish the Task having run from the Task having done nothing.
func WorkIndependentVerification(task spec.Task) (Finding, bool) {
	if len(task.Verification) == 0 {
		return Finding{}, false
	}
	for _, command := range task.Verification {
		if !workIndependentVerificationCommand(command) {
			return Finding{}, false
		}
	}

	return Finding{
		Code:     CodeVerifyWorkIndependent,
		Severity: SeverityError,
		Summary:  task.File + " declares only repository-wide gates and working-tree cleanliness checks, so its Verification cannot distinguish Task work from no work",
		Where:    []Location{{Path: task.File, Line: 1}},
		Fix:      "Add a declared Verification command that asserts this Task's own effect.",
	}, true
}

func workIndependentVerificationCommand(command string) bool {
	normalized := strings.Join(strings.Fields(strings.ToLower(command)), " ")
	return repositoryWideGate(normalized) || workingTreeCleanlinessCheck(normalized)
}

func repositoryWideGate(command string) bool {
	if makeVerifyPattern.MatchString(command) {
		return true
	}
	if !goWideGatePattern.MatchString(command) || !strings.Contains(command, "./...") {
		return false
	}
	return !strings.Contains(command, " -run ") && !strings.Contains(command, " -run=")
}

type invertedExitForm struct {
	name        string
	replacement string
}

type nonHermeticForm struct {
	name string
	fix  string
}

// InvertedExitVerification reports authored commands whose effective shell
// status is a known reversal of the condition their output appears to assert.
func InvertedExitVerification(task spec.Task) []Finding {
	var findings []Finding
	for _, command := range task.Verification {
		form, matched := invertedExitVerificationCommand(command)
		if !matched {
			continue
		}
		findings = append(findings, Finding{
			Code:     CodeVerifyInvertedExit,
			Severity: SeverityError,
			Summary:  task.File + " declares a Verification command using the " + form.name + " form: " + command,
			Where:    []Location{{Path: task.File, Line: 1}},
			Fix:      "Replace it with the working form " + form.replacement + ".",
		})
	}
	return findings
}

func invertedExitVerificationCommand(command string) (invertedExitForm, bool) {
	normalized := strings.Join(strings.Fields(strings.ToLower(command)), " ")
	switch {
	case bareTestSubstitutionPattern.MatchString(normalized):
		return invertedExitForm{
			name:        "test $(...) without a comparison",
			replacement: "`test -z \"$(cmd)\"`",
		}, true
	case filteredCountExitPattern.MatchString(normalized):
		return invertedExitForm{
			name:        "grep -v ... | wc -l filtered count",
			replacement: "`test \"$(grep -v ... | wc -l)\" -eq 0`",
		}, true
	case grepCountExitPattern.MatchString(normalized) && !testComparisonPattern.MatchString(normalized):
		return invertedExitForm{
			name:        "grep -c count-and-exit",
			replacement: "`grep -q ...` for presence, or `test \"$(grep -c ...)\" -eq <expected>` for a count",
		}, true
	default:
		return invertedExitForm{}, false
	}
}

// NonHermeticVerification reports authored commands whose success depends on
// undeclared environment state or a path outside the repository. A path first
// created by the same command is command-local state and remains permitted.
func NonHermeticVerification(task spec.Task) []Finding {
	// A qa Task's effective command is supplied by Roundfix, not authored by
	// the Task, so authoring-policy checks do not govern it.
	if task.Type == spec.TaskTypeQA {
		return nil
	}
	var findings []Finding
	createdPaths := make(map[string]bool)
	for _, command := range task.Verification {
		form, matched := nonHermeticVerificationCommand(command, createdPaths)
		if !matched {
			continue
		}
		findings = append(findings, nonHermeticFinding(task.File, command, form))
	}
	return findings
}

func nonHermeticFinding(taskFile, command string, form nonHermeticForm) Finding {
	return Finding{
		Code:     CodeVerifyNonHermetic,
		Severity: SeverityError,
		Summary:  taskFile + " declares a Verification command using the " + form.name + " form: " + command,
		Where:    []Location{{Path: taskFile, Line: 1}},
		Fix:      form.fix,
	}
}

func nonHermeticVerificationCommand(command string, createdPaths map[string]bool) (nonHermeticForm, bool) {
	assignments := shellAssignments(command)
	var environmentForm nonHermeticForm
	if variable, index, ok := environmentGuard(command); ok && !declaredBefore(assignments, variable, index) {
		environmentForm = nonHermeticForm{
			name: "environment-presence guard " + variable,
			fix:  "Remove the presence guard and make the command run from repository-carried state.",
		}
	}
	for _, reference := range shellVariableReferences(command) {
		if environmentForm.name != "" {
			break
		}
		if declaredBefore(assignments, reference.name, reference.index) {
			continue
		}
		environmentForm = nonHermeticForm{
			name: "undeclared environment variable " + reference.name,
			fix:  "Declare " + reference.name + " within the command or replace it with repository-carried input.",
		}
	}
	path, pathDependent := externalDependencyPath(command, createdPaths)
	if environmentForm.name != "" {
		return environmentForm, true
	}
	if pathDependent {
		return nonHermeticForm{
			name: "external path " + path,
			fix:  "Create " + path + " earlier in the same command before reading it, or use a repository path.",
		}, true
	}
	return nonHermeticForm{}, false
}

type shellVariableReference struct {
	name  string
	index int
}

func shellVariableReferences(command string) []shellVariableReference {
	var references []shellVariableReference
	singleQuoted := false
	doubleQuoted := false
	for index := 0; index < len(command); index++ {
		switch command[index] {
		case '\\':
			if !singleQuoted && index+1 < len(command) {
				index++
			}
			continue
		case '\'':
			if !doubleQuoted {
				singleQuoted = !singleQuoted
			}
			continue
		case '"':
			if !singleQuoted {
				doubleQuoted = !doubleQuoted
			}
			continue
		case '$':
			if singleQuoted {
				continue
			}
		}
		if command[index] != '$' || index+1 >= len(command) {
			continue
		}
		nameStart := index + 1
		if command[nameStart] == '{' {
			nameStart++
		}
		if nameStart >= len(command) || !shellNameStart(command[nameStart]) {
			continue
		}
		nameEnd := nameStart + 1
		for nameEnd < len(command) && shellNamePart(command[nameEnd]) {
			nameEnd++
		}
		references = append(references, shellVariableReference{
			name:  command[nameStart:nameEnd],
			index: index,
		})
	}
	return references
}

func shellNameStart(char byte) bool {
	return char == '_' || char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z'
}

func shellNamePart(char byte) bool {
	return shellNameStart(char) || char >= '0' && char <= '9'
}

func shellAssignments(command string) map[string]int {
	assignments := make(map[string]int)
	for _, pattern := range []*regexp.Regexp{shellAssignmentPattern, shellReadPattern, shellForPattern} {
		for _, match := range pattern.FindAllStringSubmatchIndex(command, -1) {
			name := command[match[2]:match[3]]
			if previous, exists := assignments[name]; !exists || match[2] < previous {
				assignments[name] = match[2]
			}
		}
	}
	return assignments
}

func environmentGuard(command string) (string, int, bool) {
	match := environmentGuardPattern.FindStringSubmatchIndex(command)
	if match == nil {
		return "", 0, false
	}
	return firstCapturedText(command, match, 2, 4), match[0], true
}

func firstCapturedText(input string, indexes []int, pairs ...int) string {
	for _, pair := range pairs {
		if indexes[pair] >= 0 {
			return input[indexes[pair]:indexes[pair+1]]
		}
	}
	return ""
}

func declaredBefore(assignments map[string]int, name string, reference int) bool {
	assignment, ok := assignments[name]
	return ok && assignment < reference
}

type shellWord struct {
	text string
}

func externalDependencyPath(command string, created map[string]bool) (string, bool) {
	words := shellWords(command)
	createCommand := false
	for index, word := range words {
		switch word.text {
		case ";", "&&", "||", "|":
			createCommand = false
			continue
		case "mkdir", "touch", "mktemp":
			createCommand = true
			continue
		}

		path, ok := externalPath(word.text)
		if !ok {
			continue
		}
		redirectTarget := index > 0 && (words[index-1].text == ">" || words[index-1].text == ">>" ||
			words[index-1].text == "-o" || words[index-1].text == "--output")
		if redirectTarget || createCommand {
			created[path] = true
			continue
		}
		if !created[path] {
			return path, true
		}
	}
	return "", false
}

func externalPath(word string) (string, bool) {
	if equals := strings.IndexByte(word, '='); equals >= 0 {
		word = word[equals+1:]
	}
	if strings.HasPrefix(word, "/dev/") || strings.HasPrefix(word, "/proc/self/fd/") {
		return "", false
	}
	if strings.HasPrefix(word, "/") || word == ".." || strings.HasPrefix(word, "../") || word == "~" || strings.HasPrefix(word, "~/") {
		return word, true
	}
	return "", false
}

func shellWords(command string) []shellWord {
	var words []shellWord
	var current strings.Builder
	quote := rune(0)
	flush := func() {
		if current.Len() == 0 {
			return
		}
		words = append(words, shellWord{text: current.String()})
		current.Reset()
	}
	for _, char := range command {
		if quote != 0 {
			if char == quote {
				quote = 0
				continue
			}
			current.WriteRune(char)
			continue
		}
		switch char {
		case '\'', '"':
			quote = char
		case ' ', '\t', '\n', '\r':
			flush()
		case ';', '|', '&', '<', '>', '(', ')', '{', '}':
			flush()
			operator := string(char)
			if len(words) > 0 && ((char == '>' && words[len(words)-1].text == ">") ||
				(char == '|' && words[len(words)-1].text == "|") ||
				(char == '&' && words[len(words)-1].text == "&")) {
				words[len(words)-1].text += operator
			} else {
				words = append(words, shellWord{text: operator})
			}
		default:
			current.WriteRune(char)
		}
	}
	flush()
	return words
}

// workingTreeCleanlinessCheck reports a command that asserts over the set of
// changed paths and passes when that set is empty, which is what an unchanged
// tree produces. Naming the git flags is not enough: `git diff --name-only |
// grep -q .` reads the same paths and its terminal predicate exits 1 on an
// unchanged tree, so it fails rather than passing and is honest Verification.
// Classification therefore reads the final predicate of the whole pipeline or
// chain, never a git flag off the earlier command.
func workingTreeCleanlinessCheck(command string) bool {
	if !workingTreeStatePattern.MatchString(command) {
		return false
	}
	if !strings.ContainsAny(command, "|&;") {
		// Nothing consumes the output: git prints nothing and exits zero on an
		// unchanged tree, so the command passes before any work happens.
		return true
	}
	// Something consumes it, and whether the command passes depends on what:
	// `| grep -q .` fails on empty output, `| cat` and `|| exit 0` succeed on
	// it. Only the effective terminal predicate decides the exit status.
	return chainPassesOnEmptyOutput(command)
}

// emptyInputOutcome is how a single command segment behaves on empty input.
type emptyInputOutcome int

const (
	emptyInputUnknown emptyInputOutcome = iota
	emptyInputPasses                    // exits zero on empty input
	emptyInputFails                     // exits nonzero on empty input
)

// pipelineExitOutcome returns the empty-input outcome of a pipeline: the exit
// status of a `|` pipeline is decided by its last member. A bare command with
// no pipe is its own single-member pipeline.
func pipelineExitOutcome(segment string) emptyInputOutcome {
	if idx := strings.LastIndex(segment, "|"); idx >= 0 {
		segment = segment[idx+1:]
	}
	segment = strings.TrimSpace(segment)
	switch {
	case emptyOutputFailsPattern.MatchString(segment):
		return emptyInputFails
	case emptyOutputSucceedsPattern.MatchString(segment):
		return emptyInputPasses
	case emptyTestPattern.MatchString(segment) || emptyTestBracketPattern.MatchString(segment):
		if emptyTestOperandPasses(segment) {
			return emptyInputPasses
		}
		return emptyInputUnknown
	default:
		return emptyInputUnknown
	}
}

// emptyTestOperandPasses reports whether a `test -z <operand>` or `[ -z
// <operand> ]` terminal predicate is a guaranteed success on an unchanged tree.
// Only a whole operand that is a command substitution over the working tree
// qualifies: an unchanged tree yields an empty value from it. A non-empty
// substitution of fixed output, a single-quoted form (literal text rather than
// substitution), a variable, or a literal operand can fail and is never treated
// as vacuous.
func emptyTestOperandPasses(segment string) bool {
	var operand string
	if m := emptyTestBracketPattern.FindStringSubmatch(segment); m != nil {
		operand = m[1]
	} else if m := emptyTestPattern.FindStringSubmatch(segment); m != nil {
		operand = m[1]
	} else {
		return false
	}

	cond := strings.TrimSpace(operand)
	if len(cond) >= 2 && (cond[0] == '"' || cond[0] == '\'') && cond[len(cond)-1] == cond[0] {
		cond = strings.TrimSpace(cond[1 : len(cond)-1])
	}
	// Only a working-tree command substitution proves emptiness on an
	// unchanged tree; anything else may hold output.
	if !strings.HasPrefix(cond, "$(") || !strings.HasSuffix(cond, ")") {
		return false
	}
	inner := strings.TrimSpace(cond[2 : len(cond)-1])
	return workingTreeStatePattern.MatchString(strings.Join(strings.Fields(inner), " "))
}

// chainPassesOnEmptyOutput evaluates a whole command's `&&`, `||`, and `;`
// chain over empty input, conservatively. A chain passes on empty output only
// when control flow provably reaches a consuming predicate that succeeds on
// empty output; short-circuit operators can skip a later segment entirely. An
// unknown intermediate never lets an unexecuted tail decide the result, so a
// chain whose passing tail is not guaranteed to run is never reported vacuous.
func chainPassesOnEmptyOutput(command string) bool {
	indices := commandChainOpPattern.FindAllStringIndex(command, -1)
	if len(indices) == 0 {
		return pipelineExitOutcome(command) == emptyInputPasses
	}

	var parts []string
	var ops []string
	start := 0
	for _, idx := range indices {
		parts = append(parts, command[start:idx[0]])
		ops = append(ops, command[idx[0]:idx[1]])
		start = idx[1]
	}
	parts = append(parts, command[start:])

	outcome := pipelineExitOutcome(parts[0])
	for i, op := range ops {
		right := pipelineExitOutcome(parts[i+1])
		switch op {
		case ";":
			// Sequence: the right side always runs; only it decides the result.
			outcome = right
		case "&&":
			// And: the right side runs only when the left passed; otherwise it
			// is skipped and the left outcome stands.
			if outcome == emptyInputPasses {
				outcome = right
			} else if outcome != emptyInputFails {
				outcome = emptyInputUnknown
			}
		case "||":
			// Or: the right side runs only when the left failed; otherwise it
			// is skipped and the left outcome stands.
			if outcome == emptyInputFails {
				outcome = right
			} else if outcome != emptyInputPasses {
				outcome = emptyInputUnknown
			}
		}
	}
	return outcome == emptyInputPasses
}

// VacuousVerificationCommands reports each declared command that already passes
// against the unchanged tree. WorkIndependentVerification judges the Task as a
// whole and stays silent when one honest command sits beside a vacuous one;
// the Daemon's pre-work probe judges each command on its own and refuses the
// Task for any single one, so this check applies the Daemon's unit.
func VacuousVerificationCommands(task spec.Task) []string {
	var vacuous []string
	for _, command := range task.Verification {
		normalized := strings.Join(strings.Fields(strings.ToLower(command)), " ")
		if workingTreeCleanlinessCheck(normalized) {
			vacuous = append(vacuous, command)
		}
	}
	return vacuous
}

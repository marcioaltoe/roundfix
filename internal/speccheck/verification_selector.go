package speccheck

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

// CodeVerifySubsumedSelector identifies a Verification command that passes a
// test selector together with a collection argument that contains it.
const CodeVerifySubsumedSelector = "SC-VERIFY-SUBSUMED-SELECTOR"

var pythonInterpreterPattern = regexp.MustCompile(`^python[0-9.]*$`)

// pytestValueOptions take the next word as their value, never as a
// collection argument.
var pytestValueOptions = map[string]bool{
	"-k": true, "-m": true, "-p": true, "-c": true, "-o": true, "-W": true, "-n": true, "-r": true,
	"--deselect": true, "--ignore": true, "--ignore-glob": true, "--rootdir": true, "--confcutdir": true,
	"--basetemp": true, "--maxfail": true, "--tb": true, "--durations": true, "--import-mode": true,
	"--junitxml": true, "--junit-xml": true, "--cov": true, "--cov-report": true, "--log-level": true,
	"--override-ini": true,
}

type collectionArgument struct {
	text  string
	path  string
	parts []string
}

// detectSubsumedSelectors reports a pytest node id that another collection
// argument of the same invocation contains: pytest 9 runs the containing
// argument whole and drops the selector, so a missing test passes unnoticed
// (ADR-0315). The command is read, never run.
func detectSubsumedSelectors(result *Result, repoRoot, specDir string, graph *spec.Graph) error {
	if graph == nil {
		addSkip(result, CodeVerifySubsumedSelector, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_tasks.md")))
		return nil
	}
	for _, task := range graph.Tasks {
		if task.Status != spec.StatusPending || task.ID == graph.QATaskID || task.Type == spec.TaskTypeQA {
			continue
		}
		taskPath := filepath.Join(filepath.Dir(specDir), filepath.FromSlash(task.File))
		var content []byte
		for _, command := range task.Verification {
			bodies := []string{command}
			for _, substitution := range commandSubstitutions(command) {
				bodies = append(bodies, substitution.content)
			}
			reported := make(map[string]bool)
			for _, body := range bodies {
				for _, invocation := range pytestInvocations(body) {
					for _, pair := range subsumedSelectors(invocation) {
						selector, argument := pair[0], pair[1]
						if reported[selector+"\x00"+argument] {
							continue
						}
						reported[selector+"\x00"+argument] = true
						if content == nil {
							read, err := os.ReadFile(taskPath)
							if err != nil {
								return fmt.Errorf("read Verification Task %q: %w", taskPath, err)
							}
							content = read
						}
						result.Findings = append(result.Findings, Finding{
							Code: CodeVerifySubsumedSelector, Severity: SeverityError,
							Summary: fmt.Sprintf("Task %s Verification passes the test selector %s together with %s, which contains it; pytest runs %s whole and ignores the selector, so a test that does not exist yet goes unnoticed", task.ID, selector, argument, argument),
							Where:   []Location{{Path: artifactDisplayPath(repoRoot, taskPath), Line: sectionLineContaining(content, "Verification", command)}},
							Fix:     fmt.Sprintf("Name only the selectors in one command and run %s in a separate command, following the write-tasks rule \"Test selectors and their containing path\".", argument),
						})
					}
				}
			}
		}
	}
	return nil
}

// pytestInvocations returns the collection arguments of each pytest
// invocation in one command body. An invocation with --keep-duplicates
// collects every argument, so it is left out.
func pytestInvocations(body string) [][]collectionArgument {
	words := shellWords(body)
	var invocations [][]collectionArgument
	start := 0
	for index := 0; index <= len(words); index++ {
		if index < len(words) {
			switch words[index].text {
			case ";", "&&", "||", "|", "&", "(", ")", "{", "}":
			default:
				continue
			}
		}
		if arguments, ok := pytestCollectionArguments(words[start:index]); ok {
			invocations = append(invocations, arguments)
		}
		start = index + 1
	}
	return invocations
}

func pytestCollectionArguments(words []shellWord) ([]collectionArgument, bool) {
	index := 0
	for index < len(words) {
		word := words[index].text
		if word == "!" || word == "rtk" || word == "env" || word == "exec" || shellAssignmentWord(word) {
			index++
			continue
		}
		break
	}
	if index == len(words) {
		return nil, false
	}
	next, ok := pytestRunnerEnd(words, index)
	if !ok {
		switch filepath.Base(words[index].text) {
		case "uv", "poetry", "pdm", "pipenv", "hatch":
		default:
			return nil, false
		}
		run := index + 1
		for run < len(words) && words[run].text != "run" {
			run++
		}
		for runner := run + 1; run < len(words) && runner < len(words); runner++ {
			if next, ok = pytestRunnerEnd(words, runner); ok {
				break
			}
		}
		if !ok {
			return nil, false
		}
	}
	var arguments []collectionArgument
	for i := next; i < len(words); i++ {
		word := words[i].text
		switch {
		case word == "--keep-duplicates":
			return nil, false
		case word == "<" || word == ">" || word == ">>":
			i++
		case pytestValueOptions[word]:
			i++
		case strings.HasPrefix(word, "-"):
		case !strings.Contains(word, "::") && strings.ContainsAny(word, "*?[$`"):
		default:
			arguments = append(arguments, parseCollectionArgument(word))
		}
	}
	return arguments, true
}

// pytestRunnerEnd reports whether words[index] starts a pytest runner and
// returns the index of the first word after it.
func pytestRunnerEnd(words []shellWord, index int) (int, bool) {
	base := filepath.Base(words[index].text)
	if base == "pytest" || base == "py.test" {
		return index + 1, true
	}
	if pythonInterpreterPattern.MatchString(base) && index+2 < len(words) && words[index+1].text == "-m" && words[index+2].text == "pytest" {
		return index + 3, true
	}
	return 0, false
}

func shellAssignmentWord(word string) bool {
	name, _, found := strings.Cut(word, "=")
	if !found || name == "" || !shellNameStart(name[0]) {
		return false
	}
	for index := 1; index < len(name); index++ {
		if !shellNamePart(name[index]) {
			return false
		}
	}
	return true
}

func parseCollectionArgument(word string) collectionArgument {
	parts := strings.Split(word, "::")
	location := strings.TrimSuffix(parts[0], "/")
	if location == "" {
		location = parts[0]
	}
	return collectionArgument{text: word, path: path.Clean(location), parts: parts[1:]}
}

// subsumedSelectors pairs each node id with the first argument in command
// order that contains it, under pytest 9's is_collection_argument_subsumed_by.
func subsumedSelectors(arguments []collectionArgument) [][2]string {
	var pairs [][2]string
	for _, selector := range arguments {
		if len(selector.parts) == 0 {
			continue
		}
		for _, by := range arguments {
			if by.text != selector.text && collectionArgumentContains(by, selector) {
				pairs = append(pairs, [2]string{selector.text, by.text})
				break
			}
		}
	}
	return pairs
}

func collectionArgumentContains(by, selector collectionArgument) bool {
	if len(by.parts) == 0 {
		return by.path == "." || by.path == selector.path || strings.HasPrefix(selector.path, by.path+"/")
	}
	if by.path != selector.path || len(by.parts) > len(selector.parts) {
		return false
	}
	last := len(by.parts) - 1
	for index := 0; index < last; index++ {
		if by.parts[index] != selector.parts[index] {
			return false
		}
	}
	part := selector.parts[last]
	if part == by.parts[last] {
		return true
	}
	name, _, parametrized := strings.Cut(part, "[")
	return parametrized && name == by.parts[last]
}

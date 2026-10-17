//go:build docscontract

//verify:always

package docscontract

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"unicode"
	"unicode/utf8"

	"roundfix/internal/agent"
	"roundfix/internal/cli"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
	"roundfix/internal/worktree"
)

// These are names of external products, standards, protocols, or artifact sections.
// Rule 1 also checks that removing any row would expose an uncovered phrase.
var glossaryPhraseAllowlist = []struct{ phrase, reason string }{
	{"Better Auth", "Authentication product named by the adopter decision."},
	{"Claude Code", "Coding runtime product named by the Supervisor and background-task entries."},
	{"Conventional Commits", "External commit-message convention used by Backlog Entry types."},
	{"Agent Client Protocol", "External protocol used by ACP Runtimes."},
	{"Coverage Map", "TechSpec section label read by Coverage Cross-Read."},
	{"Success Metrics", "PRD section label checked by Metric Declaration Gap."},
	{"API Contracts", "TechSpec section label checked by Contract Declaration Gap."},
}

var glossaryClosedLists = []struct {
	term       string
	values     []string
	exceptions []string
}{
	{term: "Task Type", values: func() []string {
		var values []string
		for _, taskType := range spec.AllowedTaskTypes() {
			values = append(values, string(taskType))
		}
		return values
	}()},
	{term: "Mechanical Refusal Code", values: []string{
		speccheck.CodeMechanicalAuthPaths, speccheck.CodeMechanicalConsequentOrder,
		speccheck.CodeMechanicalReportShape, speccheck.CodeMechanicalEvidencePath,
	}, exceptions: []string{"QA-PRECONDITION"}}, // A row label, not a refusal code.
	{term: "Run Worktree Reconciliation", values: []string{
		string(worktree.ReconciliationSafe), string(worktree.ReconciliationSuperseded),
		string(worktree.ReconciliationUnintegrated), string(worktree.ReconciliationDirty),
		string(worktree.ReconciliationUnknown), string(worktree.ReconciliationReleased),
	}},
	{term: "Selection Encoding", values: []string{
		agent.SelectionEncodingIndependent, agent.SelectionEncodingModelVariant,
		agent.SelectionEncodingModelManaged, agent.SelectionEncodingRuntimeManaged,
		agent.SelectionEncodingRuntimeDeferred,
	}},
}

// brokenLinks returns file:line: destination for missing relative Markdown links.
func brokenLinks(repoRoot, file, content string) []string {
	links := regexp.MustCompile(`\[[^\]]*\]\((<[^>]*>[^)]*|[^)]*)\)`)
	scheme := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)
	var broken []string
	var fence byte
	var fenceLength int
	for index, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) >= 3 && (trimmed[0] == '`' || trimmed[0] == '~') {
			length := len(trimmed) - len(strings.TrimLeft(trimmed, string(trimmed[0])))
			if length >= 3 {
				if fence == 0 {
					fence, fenceLength = trimmed[0], length
				} else if trimmed[0] == fence && length >= fenceLength && strings.TrimSpace(trimmed[length:]) == "" {
					fence, fenceLength = 0, 0
				}
				continue
			}
		}
		if fence != 0 {
			continue
		}
		for _, match := range links.FindAllStringSubmatch(line, -1) {
			destination := strings.TrimSpace(match[1])
			if strings.HasPrefix(destination, "<") {
				end := strings.Index(destination, ">")
				destination = destination[1:end]
			} else if fields := strings.Fields(destination); len(fields) > 0 {
				destination = fields[0]
			}
			if destination == "" || strings.HasPrefix(destination, "#") || scheme.MatchString(destination) {
				continue
			}
			destination = strings.SplitN(strings.SplitN(destination, "#", 2)[0], "?", 2)[0]
			if _, err := os.Stat(filepath.Join(repoRoot, filepath.Dir(file), filepath.FromSlash(destination))); err != nil {
				broken = append(broken, fmt.Sprintf("%s:%d: %s", file, index+1, destination))
			}
		}
	}
	return broken
}

func userGuideFiles(t *testing.T, includeREADME bool) []string {
	t.Helper()
	root := baselineDocumentationRepoRoot()
	files, err := filepath.Glob(filepath.Join(root, "docs", "user-guide", "*.md"))
	if err != nil {
		t.Fatalf("list user guide: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no user guide files found")
	}
	companions, err := filepath.Glob(filepath.Join(root, "docs", "user-guide", "commands", "*.md"))
	if err != nil {
		t.Fatalf("list command guides: %v", err)
	}
	if includeREADME {
		files = append(files, companions...)
		files = append(files, filepath.Join(root, "README.md"))
	}
	return files
}

func TestEveryCommandIsNamedInTheUserGuide(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("root help exited %d: %s", code, stderr.String())
	}
	var text strings.Builder
	for _, file := range userGuideFiles(t, false) {
		text.WriteString(readContractDocument(t, file))
		text.WriteByte('\n')
	}
	paths := commandPaths(stdout.String())
	if !slices.Contains(paths, "runs causes") {
		t.Fatal("root help omitted runs causes")
	}
	if missing := undocumentedCommands(paths, text.String()); len(missing) > 0 {
		t.Fatalf("user guide does not name command paths: %v", missing)
	}
}

func TestUserGuideLinksResolve(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	for _, file := range userGuideFiles(t, true) {
		relative, err := filepath.Rel(root, file)
		if err != nil {
			t.Fatalf("relative guide path: %v", err)
		}
		for _, broken := range brokenLinks(root, relative, mustRead(t, file)) {
			t.Error(broken)
		}
	}
}

func TestABrokenUserGuideLinkIsReported(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "guide"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "live.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	content := "[dead](missing.md)\n[live](../live.md)\n[URL](https://example.com)\n```markdown\n[fenced](also-missing.md)\n```\n"
	want := []string{"guide/index.md:1: missing.md"}
	if got := brokenLinks(root, "guide/index.md", content); !reflect.DeepEqual(got, want) {
		t.Fatalf("brokenLinks() = %v, want %v", got, want)
	}
}

func TestBrokenLinksHandlesDestinationDecorationsAndFences(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"live.md", "live file.md"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	content := "[live](<live file.md?view=1#section> \"title\") [plain](live.md?view=1#section \"title\")\n[anchor](#section)\n[mail](mailto:user@example.com)\n~~~markdown\n[ignored](missing.md)\n```\n[still ignored](missing.md)\n~~~\n[dead](<missing.md?view=1#section> \"title\")\n"
	want := []string{"index.md:9: missing.md"}
	if got := brokenLinks(root, "index.md", content); !reflect.DeepEqual(got, want) {
		t.Fatalf("brokenLinks() = %v, want %v", got, want)
	}
}

func TestUserGuideNamesNoRefusedQAFlag(t *testing.T) {
	for _, file := range userGuideFiles(t, true) {
		for index, line := range strings.Split(mustRead(t, file), "\n") {
			for _, suffix := range strings.Split(line, "--qa")[1:] {
				if !strings.HasPrefix(suffix, "-override") {
					t.Errorf("%s:%d: refused QA flag", file, index+1)
				}
			}
		}
	}
}

// The five rules below share this parser for both the repository and planted
// violations. Line numbers remain relative to CONTEXT.md after normalization.
type glossaryContractEntry struct {
	term  string
	line  int
	lines []string
}

type glossaryContract struct {
	text    string
	entries []glossaryContractEntry
}

func parseGlossaryContract(text string) glossaryContract {
	definition := regexp.MustCompile(`^\*\*([^*]+)\*\*:`)
	heading := regexp.MustCompile(`^#{1,6}(\s|$)`)
	doc := glossaryContract{text: text}
	current := -1
	for i, line := range strings.Split(text, "\n") {
		if match := definition.FindStringSubmatchIndex(line); match != nil {
			doc.entries = append(doc.entries, glossaryContractEntry{
				term: line[match[2]:match[3]], line: i + 1,
				lines: []string{line[match[1]:]},
			})
			current = len(doc.entries) - 1
		} else if heading.MatchString(line) {
			current = -1
		} else if current >= 0 {
			doc.entries[current].lines = append(doc.entries[current].lines, line)
		}
	}
	return doc
}

func (entry glossaryContractEntry) body() string {
	lines := slices.Clone(entry.lines)
	for i, line := range lines {
		if strings.HasPrefix(line, "_Avoid_:") {
			lines[i] = ""
		}
	}
	return strings.Join(lines, "\n")
}

type glossaryContractCodeSpan struct {
	start, end int
	value      string
}

func glossaryContractCodeSpans(text string) []glossaryContractCodeSpan {
	var spans []glossaryContractCodeSpan
	for offset := 0; offset < len(text); {
		index := strings.IndexByte(text[offset:], '`')
		if index < 0 {
			break
		}
		start := offset + index
		width := len(text[start:]) - len(strings.TrimLeft(text[start:], "`"))
		offset = start + width
		for next := offset; next < len(text); {
			index = strings.IndexByte(text[next:], '`')
			if index < 0 {
				break
			}
			next += index
			closing := len(text[next:]) - len(strings.TrimLeft(text[next:], "`"))
			if closing == width {
				spans = append(spans, glossaryContractCodeSpan{start, next + width, text[start+width : next]})
				offset = next + width
				break
			}
			next += closing
		}
	}
	return spans
}

func glossaryContractText(text string) string {
	// Strip destinations before code spans, so backticks in link labels are
	// still breaks and punctuation in destinations cannot affect sentences.
	link := regexp.MustCompile(`\[([^\]]*)\]\((?:<[^>]*>|[^)\s]*)(?:\s+"[^"\n]*")?\)`)
	text = link.ReplaceAllString(text, "$1")
	var body strings.Builder
	previous := 0
	for _, span := range glossaryContractCodeSpans(text) {
		body.WriteString(text[previous:span.start])
		body.WriteByte(0)
		body.WriteString(strings.Repeat("\n", strings.Count(text[span.start:span.end], "\n")))
		previous = span.end
	}
	body.WriteString(text[previous:])
	return body.String()
}

func glossaryContractKey(text string) string {
	return strings.ToLower(strings.Join(strings.Fields(text), " "))
}

func glossaryPhraseVocabulary(doc glossaryContract, omitted string) map[string]bool {
	known := make(map[string]bool)
	for _, entry := range doc.entries {
		key := glossaryContractKey(entry.term)
		known[key], known[key+"s"], known[key+"es"] = true, true, true
		if strings.HasSuffix(key, "y") {
			known[strings.TrimSuffix(key, "y")+"ies"] = true
		}
	}
	for _, allowed := range glossaryPhraseAllowlist {
		if allowed.phrase != omitted {
			known[glossaryContractKey(allowed.phrase)] = true
		}
	}
	return known
}

func glossaryPhraseCovered(words []string, known map[string]bool) bool {
	if known[glossaryContractKey(strings.Join(words, " "))] {
		return true
	}
	// A split must cover the whole run; no one-word piece can hide a gap.
	covered := make([]bool, len(words)+1)
	covered[0] = true
	for end := 2; end <= len(words); end++ {
		for size := 2; size <= 5 && size <= end; size++ {
			if covered[end-size] && known[glossaryContractKey(strings.Join(words[end-size:end], " "))] {
				covered[end] = true
				break
			}
		}
	}
	return covered[len(words)]
}

func glossaryPhraseFindings(doc glossaryContract, known map[string]bool) []string {
	word := regexp.MustCompile(`[\p{L}\p{N}]+(?:-[\p{Lu}\p{N}][\p{L}\p{N}]*)*(?:['’]s)?`)
	leading := strings.Fields("A An The One Each Every Its This That These Those What Like Anything No Any Some Both When If Only")
	var findings []string
	for _, entry := range doc.entries {
		body := glossaryContractText(entry.body())
		var run []string
		start, previousEnd := 0, 0
		flush := func() {
			words := run
			if len(words) > 0 && slices.Contains(leading, words[0]) {
				words = words[1:]
			}
			if len(words) >= 2 && !glossaryPhraseCovered(words, known) {
				findings = append(findings, fmt.Sprintf("CONTEXT.md:%d: %s: undefined phrase %q",
					entry.line+strings.Count(body[:start], "\n"), entry.term, strings.Join(words, " ")))
			}
			run = nil
		}
		for _, bounds := range word.FindAllStringIndex(body, -1) {
			if strings.TrimSpace(body[previousEnd:bounds[0]]) != "" {
				flush()
			}
			value := body[bounds[0]:bounds[1]]
			first, _ := utf8.DecodeRuneInString(value)
			if unicode.IsUpper(first) {
				if len(run) == 0 {
					start = bounds[0]
				}
				possessive := strings.HasSuffix(value, "'s") || strings.HasSuffix(value, "’s")
				value = strings.TrimSuffix(strings.TrimSuffix(value, "'s"), "’s")
				run = append(run, value)
				if possessive {
					flush()
				}
			} else {
				flush()
			}
			previousEnd = bounds[1]
		}
		flush()
	}
	return findings
}

func glossaryLinkFindings(doc glossaryContract, files fs.FS) []string {
	link := regexp.MustCompile(`\[[^\]]*\]\(\s*(?:<([^>]*)>|([^\s)]*))(?:\s+"[^"\n]*")?\s*\)|(?m)^\s*\[[^\]]+\]:\s*(?:<([^>]+)>|([^\s]+))`)
	var findings []string
	for _, bounds := range link.FindAllStringSubmatchIndex(doc.text, -1) {
		destination := ""
		for i := 2; i < len(bounds); i += 2 {
			if bounds[i] >= 0 {
				destination = doc.text[bounds[i]:bounds[i+1]]
				break
			}
		}
		line := 1 + strings.Count(doc.text[:bounds[0]], "\n")
		owner := "CONTEXT.md"
		for _, entry := range doc.entries {
			if line >= entry.line && line < entry.line+len(entry.lines) {
				owner = entry.term
				break
			}
		}
		problem := ""
		if !fs.ValidPath(destination) || filepath.ToSlash(filepath.Dir(destination)) != "docs/adr" {
			problem = "destination must be a file directly under docs/adr/"
		} else if content, err := fs.ReadFile(files, destination); err != nil {
			problem = err.Error()
		} else {
			lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
			if lines[0] == "---" {
				status, closed := "", false
				for _, text := range lines[1:] {
					if text == "---" {
						closed = true
						break
					}
					key, value, ok := strings.Cut(text, ":")
					if ok && key == "status" {
						status = strings.Trim(strings.TrimSpace(strings.SplitN(value, " #", 2)[0]), "\"'")
					}
				}
				if !closed || status != "accepted" {
					problem = fmt.Sprintf("ADR lifecycle status is %q, want accepted", status)
				}
			}
		}
		if problem != "" {
			findings = append(findings, fmt.Sprintf("CONTEXT.md:%d: %s: link %q: %s", line, owner, destination, problem))
		}
	}
	return findings
}

func glossaryClosedListFindings(doc glossaryContract) []string {
	var findings []string
	for _, mapped := range glossaryClosedLists {
		var values []string
		found := false
		for _, entry := range doc.entries {
			if entry.term == mapped.term {
				found = true
				for _, span := range glossaryContractCodeSpans(entry.body()) {
					values = append(values, span.value)
				}
			}
		}
		want := append(slices.Clone(mapped.values), mapped.exceptions...)
		slices.Sort(values)
		slices.Sort(want)
		values, want = slices.Compact(values), slices.Compact(want)
		if !found || !slices.Equal(values, want) {
			findings = append(findings, fmt.Sprintf("%s: code spans = %q, want %q", mapped.term, values, want))
		}
	}
	return findings
}

func glossaryAvoidFindings(doc glossaryContract) []string {
	var findings []string
	for _, entry := range doc.entries {
		count, last, avoid := 0, -1, -1
		valid := true
		for i, line := range entry.lines {
			if strings.TrimSpace(line) != "" {
				last = i
			}
			if strings.HasPrefix(line, "_Avoid_:") {
				count++
				avoid = i
				if !strings.HasPrefix(line, "_Avoid_: ") {
					valid = false
				}
				for _, item := range strings.Split(strings.TrimPrefix(line, "_Avoid_: "), ", ") {
					item = strings.TrimSpace(item)
					if item == "" || strings.ContainsAny(item, ";—():") || strings.HasSuffix(item, ".") {
						valid = false
					}
				}
			}
		}
		if count != 1 || avoid != last || !valid {
			findings = append(findings, fmt.Sprintf("CONTEXT.md:%d: %s: want exactly one final _Avoid_: plain list", entry.line, entry.term))
		}
	}
	return findings
}

func glossarySentenceFindings(doc glossaryContract) []string {
	boundary := regexp.MustCompile(`([.;?!]) `)
	owners := make(map[string]string)
	var findings []string
	for _, entry := range doc.entries {
		body := strings.Join(strings.Fields(glossaryContractText(entry.body())), " ")
		body = boundary.ReplaceAllString(body, "$1\n")
		for _, sentence := range strings.FieldsFunc(body, func(r rune) bool { return r == '\n' || r == 0 }) {
			sentence = strings.TrimSpace(sentence)
			if utf8.RuneCountInString(sentence) <= 60 {
				continue
			}
			if owner, ok := owners[sentence]; ok && owner != entry.term {
				findings = append(findings, fmt.Sprintf("%s: repeats sentence owned by %s: %q", entry.term, owner, sentence))
			} else {
				owners[sentence] = entry.term
			}
		}
	}
	return findings
}

func readGlossaryContract(t *testing.T) glossaryContract {
	t.Helper()
	doc := parseGlossaryContract(mustRead(t, filepath.Join(baselineDocumentationRepoRoot(), "CONTEXT.md")))
	if len(doc.entries) == 0 {
		t.Fatal("CONTEXT.md has no glossary entries")
	}
	return doc
}

func TestGlossaryPhrasesAreDefinedTerms(t *testing.T) {
	t.Parallel()
	doc := readGlossaryContract(t)
	findings := glossaryPhraseFindings(doc, glossaryPhraseVocabulary(doc, ""))
	for _, finding := range findings {
		t.Error(finding)
	}
	for _, allowed := range glossaryPhraseAllowlist {
		if strings.TrimSpace(allowed.phrase) == "" || strings.TrimSpace(allowed.reason) == "" {
			t.Errorf("allowlist entry %q needs a phrase and reason", allowed.phrase)
		}
		if len(glossaryPhraseFindings(doc, glossaryPhraseVocabulary(doc, allowed.phrase))) <= len(findings) {
			t.Errorf("allowlist entry %q is no longer needed by rule 1", allowed.phrase)
		}
	}
}

func TestGlossaryLinksNameAcceptedADRs(t *testing.T) {
	t.Parallel()
	for _, finding := range glossaryLinkFindings(readGlossaryContract(t), os.DirFS(baselineDocumentationRepoRoot())) {
		t.Error(finding)
	}
}

func TestGlossaryClosedListsEqualTheirConstants(t *testing.T) {
	t.Parallel()
	for _, finding := range glossaryClosedListFindings(readGlossaryContract(t)) {
		t.Error(finding)
	}
}

func TestGlossaryAvoidLinesArePlainLists(t *testing.T) {
	t.Parallel()
	for _, finding := range glossaryAvoidFindings(readGlossaryContract(t)) {
		t.Error(finding)
	}
}

func TestGlossaryStatesEachRuleOnce(t *testing.T) {
	t.Parallel()
	for _, finding := range glossarySentenceFindings(readGlossaryContract(t)) {
		t.Error(finding)
	}
}

func TestGlossaryContractReportsAPlantedViolationOfEachRule(t *testing.T) {
	t.Parallel()
	repeated := "A deliberately repeated sentence that is longer than sixty characters and belongs to one entry."
	text := "**Phrase Owner**:\nThe Spec Runner Lane carries work.\n_Avoid_: alias\n\n" +
		"**Proposed Link Owner**:\nSee [decision](docs/adr/proposed.md).\n_Avoid_: alias\n\n" +
		"**Spec Link Owner**:\nSee [plan](docs/specs/ephemeral/_prd.md).\n_Avoid_: alias\n\n" +
		"**Task Type**:\n"
	for _, taskType := range spec.AllowedTaskTypes() {
		if taskType != spec.TaskTypeQA {
			text += "`" + string(taskType) + "` "
		}
	}
	text += "\n_Avoid_: alias\n\n**Avoid Owner**:\nA definition.\n_Avoid_: This is a sentence.\n\n" +
		"**First Sentence Owner**:\n" + repeated + "\n_Avoid_: alias\n\n" +
		"**Second Sentence Owner**:\n" + repeated + "\n_Avoid_: alias\n"
	for _, mapped := range glossaryClosedLists {
		if mapped.term == "Task Type" {
			continue
		}
		text += "\n**" + mapped.term + "**:\n"
		for _, value := range append(slices.Clone(mapped.values), mapped.exceptions...) {
			text += "`" + value + "` "
		}
		text += "\n_Avoid_: alias\n"
	}
	doc := parseGlossaryContract(text)
	files := fstest.MapFS{"docs/adr/proposed.md": &fstest.MapFile{Data: []byte("---\nstatus: proposed\n---\n# Decision\n")}}
	for _, scenario := range []struct {
		name     string
		findings []string
		owners   []string
		want     int
	}{
		{"phrases", glossaryPhraseFindings(doc, glossaryPhraseVocabulary(doc, "")), []string{"Phrase Owner"}, 1},
		{"links", glossaryLinkFindings(doc, files), []string{"Proposed Link Owner", "Spec Link Owner"}, 2},
		{"closed lists", glossaryClosedListFindings(doc), []string{"Task Type"}, 1},
		{"avoid lines", glossaryAvoidFindings(doc), []string{"Avoid Owner"}, 1},
		{"one owner", glossarySentenceFindings(doc), []string{"First Sentence Owner", "Second Sentence Owner"}, 1},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			if len(scenario.findings) != scenario.want {
				t.Fatalf("got %d findings, want %d: %v", len(scenario.findings), scenario.want, scenario.findings)
			}
			for _, owner := range scenario.owners {
				if !strings.Contains(strings.Join(scenario.findings, "\n"), owner+":") &&
					!strings.Contains(strings.Join(scenario.findings, "\n"), "by "+owner+":") {
					t.Errorf("rule did not name planted entry %q: %v", owner, scenario.findings)
				}
			}
		})
	}

	// Probe the boundaries that distinguish a whole undefined run from valid
	// adjacent terms, and sentences from links, code spans, and wrapping.
	t.Run("phrase boundaries", func(t *testing.T) {
		definitions := "**Spec Runner**:\n_Avoid_: alias\n**Agent Selection**:\n_Avoid_: alias\n" +
			"**Policy Entry**:\n_Avoid_: alias\n**Six Word Term That Is Defined**:\n_Avoid_: alias\n"
		for _, scenario := range []struct {
			name, body string
			want       int
		}{
			{"wrapped and split", "The Spec Runner\nAgent Selection applies.", 0},
			{"plurals", "Spec Runners, Spec Runneres and Policy Entries apply.", 0},
			{"whole term over five words", "Six Word Term That Is Defined applies.", 0},
			{"possessive ends run", "Spec Runner's Agent Selection applies.", 0},
			{"lowercase suffix ends word", "Spec Runner-facing text applies.", 0},
			{"code is a break", "Spec `Ignored Unknown Phrase` Runner applies.", 0},
			{"code with nested backticks", "Spec ``Ignored `Unknown Phrase` text`` Runner applies.", 0},
			{"link text is retained", "[Spec Runner](docs/adr/ignored.md) applies.", 0},
			{"uncovered singleton piece", "Spec Runner Lane applies.", 1},
			{"unknown possessive", "Unknown Lane's text applies.", 1},
			{"punctuation separates runs", "Spec. Runner Lane applies.", 1},
		} {
			t.Run(scenario.name, func(t *testing.T) {
				doc := parseGlossaryContract(definitions + "**Boundary Owner**:\n" + scenario.body + "\n_Avoid_: alias\n")
				if got := glossaryPhraseFindings(doc, glossaryPhraseVocabulary(doc, "")); len(got) != scenario.want {
					t.Fatalf("findings = %v, want %d", got, scenario.want)
				}
			})
		}
	})
	t.Run("link boundaries", func(t *testing.T) {
		for _, status := range []string{"accepted", "proposed", "rejected", "deprecated", "superseded", "legacy"} {
			t.Run(status, func(t *testing.T) {
				content := "# Legacy decision\n"
				if status != "legacy" {
					content = "---\nstatus: " + status + "\n---\n# Decision\n"
				}
				files := fstest.MapFS{"docs/adr/decision.md": &fstest.MapFile{Data: []byte(content)}}
				doc := parseGlossaryContract("**Link Owner**:\n[decision](docs/adr/decision.md)\n_Avoid_: alias\n")
				got := glossaryLinkFindings(doc, files)
				want := 1
				if status == "accepted" || status == "legacy" {
					want = 0
				}
				if len(got) != want {
					t.Fatalf("findings = %v, want %d", got, want)
				}
			})
		}
	})
	t.Run("avoid boundaries", func(t *testing.T) {
		for _, scenario := range []struct {
			name, lines string
			want        int
		}{
			{"defined term allowed", "_Avoid_: Task Type, alias\n\n", 0},
			{"missing", "", 1},
			{"duplicate", "_Avoid_: alias\n_Avoid_: other\n", 1},
			{"not final", "_Avoid_: alias\nMore text.\n", 1},
			{"blank item", "_Avoid_: alias, \n", 1},
			{"semicolon", "_Avoid_: alias; other\n", 1},
			{"dash", "_Avoid_: alias — other\n", 1},
			{"parenthesis", "_Avoid_: alias (other)\n", 1},
			{"colon", "_Avoid_: alias: other\n", 1},
			{"period", "_Avoid_: alias.\n", 1},
		} {
			t.Run(scenario.name, func(t *testing.T) {
				doc := parseGlossaryContract("**Avoid Owner**:\nA definition.\n" + scenario.lines + "## End\nOutside prose.\n")
				if got := glossaryAvoidFindings(doc); len(got) != scenario.want {
					t.Fatalf("findings = %v, want %d", got, scenario.want)
				}
			})
		}
	})
	t.Run("sentence boundaries", func(t *testing.T) {
		// The first copy ends its entry, the second precedes another sentence.
		// Link destinations and line wrapping must not hide their shared text.
		first := "**First Owner**:\n" + repeated + "\n_Avoid_: alias\n"
		second := strings.Replace(repeated, "repeated sentence", "repeated\n[sentence](docs/adr/decision.md)", 1)
		doc := parseGlossaryContract(first + "**Second Owner**:\n" + second + " More prose.\n_Avoid_: alias\n")
		if got := glossarySentenceFindings(doc); len(got) != 1 || !strings.Contains(got[0], "First Owner") || !strings.Contains(got[0], "Second Owner") {
			t.Fatalf("wrapped linked sentence was not reported: %v", got)
		}
		for _, length := range []int{60, 61} {
			body := strings.Repeat("x", length)
			doc := parseGlossaryContract("**First Owner**:\n" + body + "\n_Avoid_: alias\n**Second Owner**:\n" + body + "\n_Avoid_: alias\n")
			want := 0
			if length > 60 {
				want = 1
			}
			if got := glossarySentenceFindings(doc); len(got) != want {
				t.Fatalf("sentence length %d: findings = %v, want %d", length, got, want)
			}
		}
	})

}

//go:build docscontract

//verify:always

package docscontract

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/cli"
	"roundfix/internal/delivery"
	"roundfix/internal/runevent"
	"roundfix/internal/secretscan"
	"roundfix/internal/speccheck"
	"roundfix/internal/store"
)

func TestTheTemplateNoMutationToolingRowPassesTheSpecCheck(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()
	for _, template := range []struct {
		path     string
		stage    speccheck.Stage
		artifact string
	}{
		{".agents/skills/write-prd/references/prd-template.md", speccheck.StagePRD, "_prd.md"},
		{".agents/skills/write-techspec/references/techspec-template.md", speccheck.StageTechSpec, "_techspec.md"},
	} {
		t.Run(template.path, func(t *testing.T) {
			t.Parallel()

			content := readContractDocument(t, filepath.Join(root, template.path))
			normalized := strings.Join(strings.Fields(content), " ")
			reason := regexp.MustCompile("With no protected tooling mutation, record: `([^`]+)`").FindStringSubmatch(normalized)
			if len(reason) != 2 {
				t.Fatalf("%s: missing no-mutation reason after With no protected tooling mutation, record:", template.path)
			}
			ending := regexp.MustCompile("end the row with `(bounded files: [^`]+)`").FindStringSubmatch(normalized)
			if len(ending) != 2 {
				t.Fatalf("%s: missing bounded-files ending after end the row with", template.path)
			}
			for _, grant := range []struct {
				name          string
				paths         string
				wantUnbounded int
			}{
				{"empty paths", "[]", 0},
				{"Makefile path", "\n  - Makefile", 1},
			} {
				t.Run(grant.name, func(t *testing.T) {
					t.Parallel()

					const slug = "0114-tooling-row"
					const specPath = "docs/specs/" + slug
					const recordPath = specPath + "/_authorization.md"
					repoRoot := t.TempDir()
					write := func(path, content string) {
						t.Helper()
						fullPath := filepath.Join(repoRoot, filepath.FromSlash(path))
						if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
							t.Fatal(err)
						}
						if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
							t.Fatal(err)
						}
					}
					write("docs/agents/agent-instructions.md", "# Agent instructions\n")
					constraints := "## Project Constraints\n\n" +
						"- Identifier strategy: not applicable — no identifier change. Source: `docs/agents/agent-instructions.md`.\n" +
						"- Authentication and HTTP: not applicable — no network boundary. Source: `docs/agents/agent-instructions.md`.\n" +
						"- Active ADR obligations: not applicable — no ADR applies. Source: `docs/agents/agent-instructions.md`.\n"
					if template.stage == speccheck.StageTechSpec {
						write(specPath+"/_prd.md", "# Tooling row\n\n"+constraints+
							"- Tooling authority: applicable — no protected tooling mutation proposed or authorized. Source: `docs/agents/agent-instructions.md`.\n")
					}
					row := "- Tooling authority: " + reason[1] + ". Source: `docs/agents/agent-instructions.md`; Spec-contained authorization record: `" + recordPath + "`; " + ending[1] + ".\n"
					write(specPath+"/"+template.artifact, "# Tooling row\n\n"+constraints+row)
					write(recordPath, "---\nstatus: approved\ngranted: 2026-10-30\naction: deliver this Spec\noperations: [implement, commit, push, pull_request, merge]\nconsuming: "+slug+"\npaths: "+grant.paths+"\n---\n")

					result, err := speccheck.CheckStage(filepath.Join(repoRoot, "docs", "specs"), repoRoot, slug, template.stage)
					if err != nil {
						t.Fatalf("%s: CheckStage(%s): %v", template.path, template.stage, err)
					}
					var unbounded []speccheck.Finding
					for _, finding := range result.Findings {
						if finding.Code == speccheck.CodeToolingUnbounded {
							unbounded = append(unbounded, finding)
						}
					}
					if len(unbounded) != grant.wantUnbounded {
						t.Fatalf("%s: CheckStage(%s) unbounded findings = %#v, want %d; all findings: %#v", template.path, template.stage, unbounded, grant.wantUnbounded, result.Findings)
					}
				})
			}
		})
	}
}

// commandPaths returns command paths in root help's first-seen order.
func commandPaths(help string) []string {
	var paths []string
	seen := make(map[string]bool)
	add := func(path string) {
		if !seen[path] {
			seen[path] = true
			paths = append(paths, path)
		}
	}
	for _, line := range strings.Split(help, "\n") {
		if !strings.HasPrefix(line, "  roundfix ") {
			continue
		}
		var words []string
		for _, token := range strings.Fields(strings.TrimPrefix(line, "  roundfix ")) {
			if strings.HasPrefix(token, "[") || strings.HasPrefix(token, "--") || strings.HasPrefix(token, "(") || strings.HasPrefix(token, "<") {
				if len(words) > 0 && strings.HasPrefix(token, "<") && strings.HasSuffix(token, ">") && strings.Contains(token, "|") {
					for _, alternative := range strings.Split(token[1:len(token)-1], "|") {
						add(strings.Join(words, " ") + " " + alternative)
					}
					words = nil
				}
				break
			}
			words = append(words, token)
		}
		if len(words) > 0 {
			add(strings.Join(words, " "))
		}
	}
	return paths
}

func undocumentedCommands(paths []string, text string) []string {
	var missing []string
	for _, path := range paths {
		if !strings.Contains(text, "roundfix "+path) {
			missing = append(missing, path)
		}
	}
	return missing
}

func TestCommandPathsReadSubcommandsAndAlternatives(t *testing.T) {
	help := `Usage:
  roundfix window <set|show|clear>
  roundfix archive <slug>
  roundfix baseline capabilities check [--profile <id>]
  roundfix archive <slug> [--qa-override]
  roundfix --help
roundfix ignored
  roundfix events --run <id>
  roundfix doctor (read-only)
`
	want := []string{"window set", "window show", "window clear", "archive", "baseline capabilities check", "events", "doctor"}
	if got := commandPaths(help); !reflect.DeepEqual(got, want) {
		t.Fatalf("commandPaths() = %v, want %v", got, want)
	}
}

func TestEveryCommandIsNamedInTheRoundfixSkill(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("root help exited %d: %s", code, stderr.String())
	}
	paths := commandPaths(stdout.String())
	if len(paths) < 40 {
		t.Fatalf("root help yielded %d command paths, want at least 40: %v", len(paths), paths)
	}
	for _, anchor := range []string{"deliver retry", "window clear", "baseline capabilities check", "runs causes"} {
		if !slices.Contains(paths, anchor) {
			t.Fatalf("root help did not yield anchor command %q: %v", anchor, paths)
		}
	}
	text := readContractDocument(t, filepath.Join(baselineDocumentationRepoRoot(), ".agents/skills/roundfix/SKILL.md"))
	if missing := undocumentedCommands(paths, text); len(missing) > 0 {
		t.Fatalf("Roundfix skill does not name command paths: %v", missing)
	}
}

func TestAnUndocumentedCommandIsReported(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("root help exited %d: %s", code, stderr.String())
	}
	text := readContractDocument(t, filepath.Join(baselineDocumentationRepoRoot(), ".agents/skills/roundfix/SKILL.md"))
	text = strings.ReplaceAll(text, "roundfix reopen", "")
	if got, want := undocumentedCommands(commandPaths(stdout.String()), text), []string{"reopen"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("undocumentedCommands() = %v, want %v", got, want)
	}
}

func TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem(t *testing.T) {
	t.Parallel()

	if err := checkParkActionDocumentation(parkActionDocuments(t, baselineDocumentationRepoRoot())); err != nil {
		t.Fatal(err)
	}
}

func TestTheParkActionContractReportsAnActionInTheClassTable(t *testing.T) {
	t.Parallel()

	doc := parkActionDocuments(t, baselineDocumentationRepoRoot())[0]
	row, found := parkFindingRow(doc.content)
	if !found {
		t.Fatalf("%s: finding row is missing", doc.path)
	}
	joinedRow := strings.Replace(row, parkFindingBlockerCell(), parkFindingBlockerCell()+": "+parkActionNext(delivery.BlockerSecretFound), 1)
	joined := documentedText{path: doc.path, content: strings.Replace(doc.content, row, joinedRow, 1)}
	if err := checkParkActionDocumentation([]documentedText{joined}); err == nil || !strings.Contains(err.Error(), doc.path) || !strings.Contains(err.Error(), "finding row") {
		t.Fatalf("joined finding row must be rejected for %s, got %v", doc.path, err)
	}

	rewordedSentence := strings.Replace(parkActionSentence(delivery.BlockerSecretFound), parkActionNext(delivery.BlockerSecretFound), "a different next action", 1)
	reworded := documentedText{path: doc.path, content: strings.Replace(doc.content, parkActionSentence(delivery.BlockerSecretFound), rewordedSentence, 1)}
	if err := checkParkActionDocumentation([]documentedText{reworded}); err == nil || !strings.Contains(err.Error(), doc.path) || !strings.Contains(err.Error(), delivery.BlockerSecretFound) {
		t.Fatalf("reworded secret-found action must be rejected for %s, got %v", doc.path, err)
	}
}

func parkActionDocuments(t *testing.T, root string) []documentedText {
	t.Helper()
	var documents []documentedText
	for _, path := range []string{
		"docs/user-guide/commands/deliver.md",
		".agents/skills/roundfix/references/deliver.md",
		"skills/roundfix/references/deliver.md",
	} {
		documents = append(documents, documentedText{path: path, content: readContractDocument(t, filepath.Join(root, path))})
	}
	return documents
}

func checkParkActionDocumentation(documents []documentedText) error {
	for _, document := range documents {
		row, found := parkFindingRow(document.content)
		if !found {
			return fmt.Errorf("%s: finding row is missing", document.path)
		}
		cells := strings.Split(row, "|")
		if len(cells) < 4 || strings.Trim(strings.TrimSpace(cells[1]), "`") != string(delivery.ParkClassFinding) {
			return fmt.Errorf("%s: malformed finding row %q", document.path, row)
		}
		if got, want := strings.TrimSpace(cells[2]), parkFindingBlockerCell(); got != want {
			return fmt.Errorf("%s: finding row lists %q; want blocker names only %q", document.path, got, want)
		}

		normalized := strings.Join(strings.Fields(document.content), " ")
		for _, blocker := range []string{delivery.BlockerArchivedItemUnproven, delivery.BlockerSecretFound} {
			want := parkActionSentence(blocker)
			if count := strings.Count(normalized, want); count != 1 {
				return fmt.Errorf("%s: documented action for %q occurred %d times; want exactly once: %q", document.path, blocker, count, want)
			}
		}
	}
	return nil
}

func parkFindingRow(content string) (string, bool) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "| `"+string(delivery.ParkClassFinding)+"` |") {
			return line, true
		}
	}
	return "", false
}

func parkFindingBlockerCell() string {
	blockers := []string{
		delivery.BlockerRunUnresolved,
		delivery.BlockerReviewFindings,
		delivery.BlockerCorrectiveSpecRequired,
		delivery.BlockerGateFailed,
		delivery.BlockerChecksFailed,
		delivery.BlockerRevalidationFailed,
		delivery.BlockerArchivedItemUnproven,
		delivery.BlockerSecretFound,
	}
	for i := range blockers {
		blockers[i] = "`" + blockers[i] + "`"
	}
	return strings.Join(blockers, ", ")
}

func parkActionNext(blocker string) string {
	return delivery.ClassifyPark(store.DeliveryQueue{}, store.DeliveryQueueItem{
		SpecSlug: "<slug>",
		Worktree: "<worktree>",
		Blocker:  blocker,
	}).Next
}

func parkActionSentence(blocker string) string {
	article := "A"
	if strings.HasPrefix(blocker, "a") {
		article = "An"
	}
	return fmt.Sprintf("%s `%s` park's next action is `%s`.", article, blocker, parkActionNext(blocker))
}

func TestTheQASecretRefusalIsDocumentedAsTheDaemonPrintsIt(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()
	documented := qaSecretRefusalDocumentation(t, root)
	if err := checkQASecretRefusalDocumentation(qaSecretRefusalDocuments(t, root), documented); err != nil {
		t.Fatal(err)
	}
}

func TestTheQASecretRefusalContractRejectsTheParenthesizedForm(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()
	documented := qaSecretRefusalDocumentation(t, root)
	parenthesized := strings.Replace(documented, ":<line>: <kind>", ":<line> (<kind>)", 1)
	err := checkQASecretRefusalDocumentation([]qaSecretRefusalDocument{{
		path:    "docs/user-guide/commands/implement.md",
		content: parenthesized,
	}}, documented)
	if err == nil {
		t.Fatal("parenthesized QA secret refusal was accepted")
	}
	if !strings.Contains(err.Error(), "docs/user-guide/commands/implement.md") || !strings.Contains(err.Error(), parenthesized) {
		t.Fatalf("parenthesized QA secret refusal failure does not name the old sentence: %v", err)
	}
}

type qaSecretRefusalDocument struct {
	path    string
	content string
}

func qaSecretRefusalDocumentation(t *testing.T, root string) string {
	t.Helper()
	const prefix = "QA evidence carries a secret"
	var literals []string
	if err := filepath.WalkDir(filepath.Join(root, "internal", "daemon"), func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(file) != ".go" || strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		source, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), file, source, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", file, err)
		}
		ast.Inspect(parsed, func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}
			value, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Fatalf("unquote %s: %v", file, err)
			}
			if strings.HasPrefix(value, prefix) {
				literals = append(literals, value)
			}
			return true
		})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(literals) != 1 {
		t.Fatalf("found %d Daemon QA secret refusal literals, want exactly one", len(literals))
	}

	sample := secretscan.Match{Path: "<path>", Line: 1, Kind: "<kind>"}.String()
	documented, ok := strings.CutPrefix(strings.Replace(sample, "<path>:1:", "<path>:<line>:", 1), "<path>:<line>:")
	if !ok {
		t.Fatalf("Secret Match.String() did not produce the expected location: %q", sample)
	}
	return literals[0] + "<path>:<line>:" + documented
}

func qaSecretRefusalDocuments(t *testing.T, root string) []qaSecretRefusalDocument {
	t.Helper()
	var documents []qaSecretRefusalDocument
	directories := []string{
		"docs/user-guide",
		"docs/agents",
		".agents/skills",
		"skills",
	}
	for _, directory := range directories {
		if err := filepath.WalkDir(filepath.Join(root, directory), func(file string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || filepath.Ext(file) != ".md" {
				return nil
			}
			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read %s: %w", file, err)
			}
			relative, err := filepath.Rel(root, file)
			if err != nil {
				return fmt.Errorf("relativize %s: %w", file, err)
			}
			documents = append(documents, qaSecretRefusalDocument{
				path:    filepath.ToSlash(relative),
				content: string(content),
			})
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	content, err := os.ReadFile(filepath.Join(root, "CONTEXT.md"))
	if err != nil {
		t.Fatalf("read CONTEXT.md: %v", err)
	}
	documents = append(documents, qaSecretRefusalDocument{path: "CONTEXT.md", content: string(content)})
	sort.Slice(documents, func(i, j int) bool { return documents[i].path < documents[j].path })
	return documents
}

func checkQASecretRefusalDocumentation(documents []qaSecretRefusalDocument, documented string) error {
	const prefix = "QA evidence carries a secret"
	required := map[string]int{
		".agents/skills/roundfix/references/implement.md": 0,
		"skills/roundfix/references/implement.md":         0,
		"docs/user-guide/commands/implement.md":           0,
	}
	for _, document := range documents {
		normalized := strings.Join(strings.Fields(document.content), " ")
		for offset := 0; ; {
			index := strings.Index(normalized[offset:], prefix)
			if index < 0 {
				break
			}
			index += offset
			occurrence := normalized[index:]
			if !strings.HasPrefix(occurrence, documented) {
				return fmt.Errorf("%s: found QA secret refusal %q, want prefix %q", document.path, occurrence, documented)
			}
			if _, ok := required[document.path]; ok {
				required[document.path]++
			}
			offset = index + len(prefix)
		}
	}
	for path, count := range required {
		if count != 1 {
			return fmt.Errorf("%s: found %d QA secret refusal occurrences, want exactly one", path, count)
		}
	}
	return nil
}

func TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()
	version, conventions := deliveryConventionSource(t, root)
	var docs []documentedText
	for _, path := range []string{
		".agents/skills/roundfix/references/review.md",
		"skills/roundfix/references/review.md",
		"docs/user-guide/commands/review.md",
	} {
		docs = append(docs, documentedText{path: path, content: readContractDocument(t, filepath.Join(root, path))})
	}
	if err := checkDeliveryConventionDocumentation(docs, version, conventions); err != nil {
		t.Fatal(err)
	}
}

func TestTheDeliveryConventionContractReportsAStaleConvention(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()
	version, conventions := deliveryConventionSource(t, root)
	const path = ".agents/skills/roundfix/references/review.md"
	text := strings.Join(strings.Fields(readContractDocument(t, filepath.Join(root, path))), " ")
	var prompt string
	for _, convention := range conventions {
		if convention.id == "C3" {
			prompt = strings.Join(strings.Fields(convention.prompt), " ")
		}
	}
	if prompt == "" || !strings.Contains(text, prompt) {
		t.Fatalf("%s: cannot plant stale C3: parsed prompt is missing", path)
	}
	stale := strings.Replace(text, prompt, "The archive commit moves a completed Spec's directory to the archive root and stamps its archive front matter.", 1)
	err := checkDeliveryConventionDocumentation([]documentedText{{path: path, content: stale}}, version, conventions)
	if err == nil || !strings.Contains(err.Error(), path) || !strings.Contains(err.Error(), "C3") {
		t.Fatalf("stale C3 must report %s and C3, got %v", path, err)
	}
}

type deliveryConventionText struct {
	id, prompt string
}

type documentedText struct {
	path, content string
}

func deliveryConventionSource(t *testing.T, root string) (string, []deliveryConventionText) {
	t.Helper()
	const path = "internal/cli/review_conventions.go"
	parsed, err := parser.ParseFile(token.NewFileSet(), filepath.Join(root, path), nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	stringValue := func(expr ast.Expr) (string, bool) {
		literal, ok := expr.(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return "", false
		}
		value, err := strconv.Unquote(literal.Value)
		if err != nil {
			t.Fatalf("unquote %s: %v", path, err)
		}
		return value, true
	}
	var version string
	var conventions []deliveryConventionText
	idPattern := regexp.MustCompile(`^C[0-9]+$`)
	for _, declaration := range parsed.Decls {
		switch declaration := declaration.(type) {
		case *ast.GenDecl:
			if declaration.Tok != token.CONST {
				continue
			}
			for _, spec := range declaration.Specs {
				value := spec.(*ast.ValueSpec)
				for i, name := range value.Names {
					if name.Name == "deliveryConventionsVersion" && i < len(value.Values) {
						version, _ = stringValue(value.Values[i])
					}
				}
			}
		case *ast.FuncDecl:
			if declaration.Name.Name != "deliveryConventions" || declaration.Body == nil {
				continue
			}
			ast.Inspect(declaration.Body, func(node ast.Node) bool {
				literal, ok := node.(*ast.CompositeLit)
				if !ok || len(literal.Elts) != 2 {
					return true
				}
				id, idOK := stringValue(literal.Elts[0])
				prompt, promptOK := stringValue(literal.Elts[1])
				if idOK && promptOK && idPattern.MatchString(id) {
					conventions = append(conventions, deliveryConventionText{id: id, prompt: prompt})
				}
				return true
			})
		}
	}
	if version == "" || len(conventions) == 0 {
		t.Fatalf("%s: deliveryConventionsVersion = %q, found %d conventions; want a non-empty version and conventions", path, version, len(conventions))
	}
	return version, conventions
}

func checkDeliveryConventionDocumentation(docs []documentedText, version string, conventions []deliveryConventionText) error {
	items := regexp.MustCompile(`- C[0-9]+\. `)
	for _, doc := range docs {
		text := strings.Join(strings.Fields(doc.content), " ")
		if !strings.Contains(text, "Delivery Conventions version `"+version+"`") {
			return fmt.Errorf("%s: missing Delivery Conventions version %q", doc.path, version)
		}
		for _, convention := range conventions {
			prompt := strings.Join(strings.Fields(convention.prompt), " ")
			if !strings.Contains(text, "- "+convention.id+". "+prompt) {
				return fmt.Errorf("%s: missing or stale %s prompt", doc.path, convention.id)
			}
		}
		if found := items.FindAllString(text, -1); len(found) != len(conventions) {
			return fmt.Errorf("%s: convention items %q: got %d, want %d", doc.path, found, len(found), len(conventions))
		}
	}
	return nil
}

func TestTheEventCategoriesAreDocumentedAsTheStreamAcceptsThem(t *testing.T) {
	t.Parallel()

	if err := checkEventCategoryDocumentation(eventCategoryDocuments(t), documentedStreamCategories()); err != nil {
		t.Fatal(err)
	}
}

func TestTheEventCategoryContractReportsAMissingCategory(t *testing.T) {
	t.Parallel()

	doc := eventCategoryDocuments(t)[0]
	before, after, found := strings.Cut(doc.content, doc.anchor)
	end := strings.IndexAny(after, ".;")
	if !found || end < 0 {
		t.Fatalf("%s: cannot find category sentence", doc.path)
	}
	sentence := strings.Join(strings.Fields(after[:end]), " ")
	const removed = ", and `usage`"
	if !strings.Contains(sentence, removed) {
		t.Fatalf("%s: cannot remove usage from category sentence", doc.path)
	}
	doc.content = before + doc.anchor + " " + strings.Replace(sentence, removed, "", 1) + after[end:]
	err := checkEventCategoryDocumentation([]anchoredText{doc}, documentedStreamCategories())
	if err == nil || !strings.Contains(err.Error(), doc.path) || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("missing usage must report %s and usage, got %v", doc.path, err)
	}
}

type anchoredText struct {
	documentedText
	anchor string
	table  bool
}

func documentedStreamCategories() []string {
	var categories []string
	for category := range runevent.AllStreamCategories() {
		categories = append(categories, string(category))
	}
	slices.Sort(categories)
	return categories
}

func eventCategoryDocuments(t *testing.T) []anchoredText {
	t.Helper()
	root := baselineDocumentationRepoRoot()
	docs := []anchoredText{
		{documentedText: documentedText{path: ".agents/skills/roundfix/references/events.md"}, anchor: "Default replay emits these public categories in journal cursor order:", table: true},
		{documentedText: documentedText{path: "skills/roundfix/references/events.md"}, anchor: "Default replay emits these public categories in journal cursor order:", table: true},
		{documentedText: documentedText{path: "docs/user-guide/commands/events.md"}, anchor: "The public Supervisor categories, in journal cursor order, are"},
	}
	for i := range docs {
		docs[i].content = readContractDocument(t, filepath.Join(root, docs[i].path))
	}
	return docs
}

func checkEventCategoryDocumentation(docs []anchoredText, categories []string) error {
	codeSpan := regexp.MustCompile("`([^`]+)`")
	for _, doc := range docs {
		text := strings.Join(strings.Fields(doc.content), " ")
		_, sentence, found := strings.Cut(text, doc.anchor)
		if !found {
			return fmt.Errorf("%s: missing category anchor %q", doc.path, doc.anchor)
		}
		if end := strings.IndexAny(sentence, ".;"); end >= 0 {
			sentence = sentence[:end]
		}
		var got []string
		for _, span := range codeSpan.FindAllStringSubmatch(sentence, -1) {
			got = append(got, span[1])
		}
		slices.Sort(got)
		if !slices.Equal(got, categories) {
			missing := slices.Clone(categories)
			var extra []string
			for _, category := range got {
				if i := slices.Index(missing, category); i >= 0 {
					missing = slices.Delete(missing, i, i+1)
				} else {
					extra = append(extra, category)
				}
			}
			return fmt.Errorf("%s: category sentence missing %q; extra %q", doc.path, missing, extra)
		}
		if doc.table {
			var missing []string
			for _, category := range categories {
				prefix := "| `" + category + "` |"
				found := false
				for _, line := range strings.Split(doc.content, "\n") {
					if strings.HasPrefix(line, prefix) {
						found = true
						break
					}
				}
				if !found {
					missing = append(missing, category)
				}
			}
			if len(missing) > 0 {
				return fmt.Errorf("%s: stable-fields table missing categories %q", doc.path, missing)
			}
		}
	}
	return nil
}

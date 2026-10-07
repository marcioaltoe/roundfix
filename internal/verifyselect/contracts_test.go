// Suite: Repository Contract Test selection.
// Invariant: tagged contracts declare valid relevance and changes select the documented invocations.
// Boundary IN: temporary Go sources, local Git changes, and verify-select streams and exit codes.
// Boundary OUT: real repository contracts and execution of the printed invocations.
package verifyselect_test

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/verifyselect"
)

func contractSource(tag, header, declarations string) string {
	return "//go:build " + tag + "\n\n" + header + "\npackage fixture\nimport \"testing\"\n" + declarations + "\n"
}

func TestContractDiscoveryReadsTagsAndDirectives(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	writeFile(t, repo, "z/docs_test.go", contractSource("docscontract", "//verify:always\n\n// a separate header comment", `func TestZ(t *testing.T) {}
//verify:relevant inputs/ assets/*.json
func TestA(t *testing.T) {}
//verify:boundary regenerates artifacts
func TestBoundary(t *testing.T) {}
func Testhelper(t *testing.T) {}
func TestWrong(t string) {}
func TestResult(t *testing.T) int { return 1 }
func TestMany(t, other *testing.T) {}
type example struct{}
func (example) TestMethod(t *testing.T) {}
`))
	writeFile(t, repo, "a/repo_test.go", contractSource("repocontract", "", "func TestPackage(t *testing.T) {}"))
	writeFile(t, repo, "a/ignored_test.go", "package fixture\nimport \"testing\"\n//verify:sometimes\nfunc TestIgnored(t *testing.T) {}\n")
	for _, directory := range []string{"testdata", "vendor", "node_modules", ".git", ".hidden", "z/testdata"} {
		writeFile(t, repo, directory+"/ignored_test.go", "this is not Go")
	}
	got, err := verifyselect.DiscoverContracts(repo)
	if err != nil {
		t.Fatal(err)
	}
	want := []verifyselect.ContractTest{
		{Name: "TestPackage", Package: "a", Tag: "repocontract", File: "a/repo_test.go", Class: verifyselect.ContractPackage},
		{Name: "TestA", Package: "z", Tag: "docscontract", File: "z/docs_test.go", Class: verifyselect.ContractRelevant, Paths: []string{"inputs/", "assets/*.json"}},
		{Name: "TestBoundary", Package: "z", Tag: "docscontract", File: "z/docs_test.go", Class: verifyselect.ContractBoundary, Reason: "regenerates artifacts"},
		{Name: "TestZ", Package: "z", Tag: "docscontract", File: "z/docs_test.go", Class: verifyselect.ContractAlways},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DiscoverContracts() = %#v, want %#v", got, want)
	}
}

func TestContractDirectiveRefusesMalformedDeclarations(t *testing.T) {
	t.Parallel()
	cases := []struct{ name, directive, message string }{
		{"unknown", "//verify:sometimes", `unknown Contract Relevance class "sometimes"`},
		{"empty class", "//verify:", `unknown Contract Relevance class ""`},
		{"explicit package", "//verify:package", `unknown Contract Relevance class "package"`},
		{"missing paths", "//verify:relevant", "relevant needs at least one path"},
		{"missing reason", "//verify:boundary", "boundary needs a reason"},
		{"multiple", "//verify:always\n//verify:relevant inputs/", "more than one Contract Relevance directive"},
	}
	for _, pattern := range []string{"/abs", "./x", "a/../b", "a/**", "[", `a\b`} {
		cases = append(cases, struct{ name, directive, message string }{"path " + pattern, "//verify:relevant " + pattern, fmt.Sprintf("invalid relevant path %q", pattern)})
	}
	for _, test := range cases {
		for _, scope := range []string{"header", "test"} {
			t.Run(test.name+"/"+scope, func(t *testing.T) {
				t.Parallel()
				repo := t.TempDir()
				header, doc, prefix := test.directive, "", "pkg/contract_test.go: "
				if scope == "test" {
					header, doc = "//verify:always", test.directive+"\n"
					prefix += "TestFixture: "
				}
				writeFile(t, repo, "pkg/contract_test.go", contractSource("repocontract", header, doc+"func TestFixture(t *testing.T) {}"))
				_, err := verifyselect.DiscoverContracts(repo)
				if err == nil || err.Error() != prefix+test.message {
					t.Fatalf("error = %v, want %q", err, prefix+test.message)
				}
			})
		}
	}
}

func TestContractDiscoveryValidatesBuildConstraints(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		expression string
		wantTag    string
		invalid    bool
	}{
		{"docscontract && !repocontract", "docscontract", false},
		{"repocontract && !docscontract", "repocontract", false},
		{"integration", "", false},
		{"docscontract || repocontract", "", true},
		{"docscontract && repocontract", "", true},
		{"!docscontract", "", true},
		{"!integration || docscontract", "", true},
		{"docscontract && linux", "", true},
	} {
		t.Run(test.expression, func(t *testing.T) {
			t.Parallel()
			repo := t.TempDir()
			writeFile(t, repo, "contract_test.go", contractSource(test.expression, "", "func TestFixture(t *testing.T) {}"))
			got, err := verifyselect.DiscoverContracts(repo)
			if test.invalid {
				if err == nil || !strings.Contains(err.Error(), "contract_test.go:") {
					t.Fatalf("discovery error = %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if test.wantTag == "" {
				if len(got) != 0 {
					t.Fatalf("contracts = %v", got)
				}
				return
			}
			if len(got) != 1 || got[0].Tag != test.wantTag {
				t.Fatalf("contracts = %v, want tag %s", got, test.wantTag)
			}
		})
	}
	t.Run("legacy and testing alias", func(t *testing.T) {
		repo := t.TempDir()
		writeFile(t, repo, "contract_test.go", "// +build repocontract\n\npackage fixture\nimport check \"testing\"\nfunc TestAlias(t *check.T) {}\n")
		got, err := verifyselect.DiscoverContracts(repo)
		if err != nil || len(got) != 1 || got[0].Tag != "repocontract" {
			t.Fatalf("contracts = %v, error = %v", got, err)
		}
	})
}

func selectionContracts() []verifyselect.ContractTest {
	return []verifyselect.ContractTest{
		{Name: "TestAlways", Package: "always", Class: verifyselect.ContractAlways},
		{Name: "TestPackage", Package: "pkg", Class: verifyselect.ContractPackage},
		{Name: "TestRelevant", Package: "relevant", Class: verifyselect.ContractRelevant, Paths: []string{"inputs/", "assets/*.json"}},
		{Name: "TestBoundary", Package: "pkg", Class: verifyselect.ContractBoundary},
	}
}

func TestContractSelectionFollowsChangedPaths(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		changed []string
		indexes []int
	}{
		{"empty", nil, []int{0}},
		{"direct package file", []string{"pkg/code.go"}, []int{0, 1}},
		{"package subdirectory", []string{"pkg/sub/code.go"}, []int{0}},
		{"relevant package", []string{"relevant/code.go"}, []int{0, 2}},
		{"directory pattern", []string{"inputs/nested/data.txt"}, []int{0, 2}},
		{"glob", []string{"assets/data.json"}, []int{0, 2}},
		{"glob subdirectory", []string{"assets/sub/data.json"}, []int{0}},
		{"directory sibling", []string{"inputs-other/data.txt"}, []int{0}},
		{"duplicate paths", []string{"pkg/code.go", "pkg/code.go"}, []int{0, 1}},
	}
	for _, file := range []string{"go.mod", "go.sum", "Makefile"} {
		cases = append(cases, struct {
			name    string
			changed []string
			indexes []int
		}{file, []string{file}, []int{0, 1, 2}})
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			all := selectionContracts()
			var want []verifyselect.ContractTest
			for _, index := range test.indexes {
				want = append(want, all[index])
			}
			if got := verifyselect.SelectContractPaths(all, test.changed); !reflect.DeepEqual(got, want) {
				t.Fatalf("selection = %v, want %v", got, want)
			}
		})
	}
}

func TestContractSelectionFailsSafeWithoutAChangeList(t *testing.T) {
	repo := newGitRepository(t)
	writeFile(t, repo, "pkg/contract_test.go", contractSource("repocontract", "", `//verify:always
func TestAlways(t *testing.T) {}
func TestPackage(t *testing.T) {}
//verify:relevant inputs/
func TestRelevant(t *testing.T) {}
//verify:boundary manual gate
func TestBoundary(t *testing.T) {}`))
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-q", "-m", "fixture")
	all, selected, err := verifyselect.SelectContracts(t.Context(), repo, "missing-base")
	if err == nil || !strings.Contains(err.Error(), "resolve merge base") {
		t.Fatalf("error = %v, want merge-base failure", err)
	}
	if len(all) != 4 || len(selected) != 3 {
		t.Fatalf("all = %v, selected = %v", all, selected)
	}
	for _, contract := range selected {
		if contract.Class == verifyselect.ContractBoundary {
			t.Fatal("selected boundary")
		}
	}
	var stdout, stderr bytes.Buffer
	code := verifyselect.Run(t.Context(), []string{"-repo", repo, "-contracts", "-base", "missing-base"}, &stdout, &stderr)
	if code != 0 || stdout.String() != "repocontract ^(TestAlways|TestPackage|TestRelevant)$ ./pkg\n" || !strings.Contains(stderr.String(), "; selecting every contract that is not boundary\n") {
		t.Fatalf("Run() = %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
	t.Run("empty repository fails safe", func(t *testing.T) {
		var out, diagnostic bytes.Buffer
		code := verifyselect.Run(t.Context(), []string{"-repo", t.TempDir(), "-contracts"}, &out, &diagnostic)
		if code != 0 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "; selecting every contract that is not boundary\n") {
			t.Fatalf("Run() = %d, stdout %q, stderr %q", code, out.String(), diagnostic.String())
		}
	})
}

func TestContractInvocationsSortAndDeduplicate(t *testing.T) {
	t.Parallel()
	selected := []verifyselect.ContractTest{
		{Name: "TestZ", Tag: "repocontract", Package: "z"},
		{Name: "TestZ", Tag: "repocontract", Package: "a"},
		{Name: "TestA", Tag: "repocontract", Package: "a"},
		{Name: "TestA", Tag: "repocontract", Package: "a"},
		{Name: "Test.+", Tag: "docscontract", Package: "docs"},
	}
	want := []string{`docscontract ^(Test\.\+)$ ./docs`, "repocontract ^(TestA|TestZ)$ ./a ./z"}
	if got := verifyselect.ContractInvocations(selected); !reflect.DeepEqual(got, want) {
		t.Fatalf("invocations = %q, want %q", got, want)
	}
	if got := verifyselect.ContractInvocations(nil); len(got) != 0 {
		t.Fatalf("empty invocations = %v", got)
	}
}

func TestRunPrintsContractInvocations(t *testing.T) {
	repo := newGitRepository(t)
	writeFile(t, repo, "docs/contract_test.go", contractSource("docscontract", "//verify:always", "func TestDocs(t *testing.T) {}"))
	writeFile(t, repo, "pkg/contract_test.go", contractSource("repocontract", "", `func TestPackage(t *testing.T) {}
//verify:always
func TestAlways(t *testing.T) {}
//verify:boundary manual gate
func TestBoundary(t *testing.T) {}`))
	writeFile(t, repo, "other/contract_test.go", contractSource("repocontract", "", "func TestOther(t *testing.T) {}"))
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-q", "-m", "base contracts")
	base := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))
	writeFile(t, repo, "pkg/code.go", "package fixture\n")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-q", "-m", "package change")
	args := []string{"-repo", repo, "-contracts", "-base", base}
	var stdout, stderr bytes.Buffer
	code := verifyselect.Run(t.Context(), args, &stdout, &stderr)
	wantOut := "docscontract ^(TestDocs)$ ./docs\nrepocontract ^(TestAlways|TestPackage)$ ./pkg\n"
	wantErr := "verify-select: contracts: 3 selected (2 always, 1 by change), 2 not selected, 1 boundary\n"
	if code != 0 || stdout.String() != wantOut || stderr.String() != wantErr {
		t.Fatalf("Run() = %d, stdout %q, stderr %q; want 0, %q, %q", code, stdout.String(), stderr.String(), wantOut, wantErr)
	}
	t.Run("malformed directive", func(t *testing.T) {
		writeFile(t, repo, "bad/contract_test.go", contractSource("repocontract", "//verify:sometimes", "func TestBad(t *testing.T) {}"))
		var out, diagnostic bytes.Buffer
		code := verifyselect.Run(t.Context(), args, &out, &diagnostic)
		want := "verify-select: bad/contract_test.go: unknown Contract Relevance class \"sometimes\"\n"
		if code != 1 || out.Len() != 0 || diagnostic.String() != want {
			t.Fatalf("Run() = %d, stdout %q, stderr %q, want diagnostic %q", code, out.String(), diagnostic.String(), want)
		}
	})
	for _, conflict := range [][]string{{"-packages", "core"}, {"-packages="}, {"-baseline-cli-pattern"}, {"-baseline-cli-pattern=false"}} {
		t.Run("conflict "+strings.Join(conflict, " "), func(t *testing.T) {
			var out, diagnostic bytes.Buffer
			code := verifyselect.Run(t.Context(), append(append([]string{}, args...), conflict...), &out, &diagnostic)
			want := "verify-select: -contracts cannot be combined with -packages or -baseline-cli-pattern\n"
			if code != 2 || out.Len() != 0 || diagnostic.String() != want {
				t.Fatalf("Run() = %d, stdout %q, stderr %q", code, out.String(), diagnostic.String())
			}
		})
	}
}

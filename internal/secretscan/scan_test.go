package secretscan_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/secretscan"
)

func secretSamples() []struct{ kind, value string } {
	headerStart, headerEnd := "-----BEGIN ", "-----"
	return []struct{ kind, value string }{
		{"private-key-block", headerStart + "RSA " + "PRIVATE KEY" + headerEnd},
		{"github-token", "ghp_" + strings.Repeat("G7", 18)},
		{"openrouter-key", "sk-or-" + "v1-" + strings.Repeat("a7", 32)},
		{"anthropic-key", "sk-ant-" + strings.Repeat("C4", 16)},
		{"openai-key", "sk-" + "proj-" + strings.Repeat("D3", 16)},
		{"aws-access-key-id", "AKIA" + strings.Repeat("E2", 8)},
		{"slack-token", "xoxb-" + strings.Repeat("F1", 5)},
		{"google-api-key", "AIza" + strings.Repeat("H6", 17) + "Z"},
		{"stripe-secret-key", "sk_live_" + strings.Repeat("J5", 12)},
		{"credential-assignment", "SERVICE_API_KEY" + " := " + strings.Repeat("K8", 10)},
	}
}

func assertNoSampleRun(t *testing.T, output, sample string) {
	t.Helper()
	for index := 0; index+8 <= len(sample); index++ {
		if strings.Contains(output, sample[index:index+8]) {
			t.Fatal("diagnostic contains a sample byte run")
		}
	}
}

func TestSecretPatternsMatchTheirKindAndNeverEchoTheValue(t *testing.T) {
	t.Parallel()
	samples := secretSamples()
	var wantKinds []string
	for _, sample := range samples {
		wantKinds = append(wantKinds, sample.kind)
		t.Run(sample.kind, func(t *testing.T) {
			t.Parallel()
			matches := secretscan.ScanContent("evidence.txt", []byte(sample.value))
			want := []secretscan.Match{{Path: "evidence.txt", Line: 1, Kind: sample.kind}}
			if !reflect.DeepEqual(matches, want) {
				t.Fatalf("matches = %+v, want %+v", matches, want)
			}
			diagnostic := (&secretscan.MatchesError{Matches: matches}).Error()
			if diagnostic != "evidence.txt:1: "+sample.kind || matches[0].String() != diagnostic {
				t.Fatalf("diagnostic = %q", diagnostic)
			}
			encoded, err := json.Marshal(matches)
			if err != nil {
				t.Fatal(err)
			}
			for _, output := range []string{diagnostic, string(encoded), matches[0].String()} {
				assertNoSampleRun(t, output, sample.value)
			}
		})
	}
	if got := secretscan.Kinds(); !reflect.DeepEqual(got, wantKinds) {
		t.Fatalf("kinds = %v, want %v", got, wantKinds)
	}
	kinds := secretscan.Kinds()
	kinds[0] = "changed"
	if !reflect.DeepEqual(secretscan.Kinds(), wantKinds) {
		t.Fatal("Kinds exposed mutable table state")
	}
	headerStart, headerEnd := "-----BEGIN ", "-----"
	variants := []struct{ name, kind, value string }{
		{"private header", "private-key-block", headerStart + "PRIVATE KEY" + headerEnd},
		{"encrypted header", "private-key-block", headerStart + "ENCRYPTED PRIVATE KEY" + headerEnd},
		{"fine grained", "github-token", "github_pat_" + strings.Repeat("Z9_", 27) + "Q"},
		{"plain openai", "openai-key", "sk-" + strings.Repeat("V2", 16)},
		{"router shape falls through", "openai-key", samples[2].value + "g"},
		{"session aws", "aws-access-key-id", "ASIA" + strings.Repeat("N5", 8)},
		{"restricted stripe", "stripe-secret-key", "rk_live_" + strings.Repeat("M4", 12)},
	}
	for _, prefix := range []string{"gho_", "ghu_", "ghs_", "ghr_"} {
		variants = append(variants, struct{ name, kind, value string }{prefix, "github-token", prefix + strings.Repeat("P2", 18)})
	}
	for _, prefix := range []string{"xoxa-", "xoxp-", "xoxr-", "xoxs-"} {
		variants = append(variants, struct{ name, kind, value string }{prefix, "slack-token", prefix + strings.Repeat("S4", 5)})
	}
	for _, name := range []string{"API_KEY", "TOKEN", "SECRET", "PASSWORD", "MYTOKEN"} {
		for _, operator := range []string{"=", ":=", "?=", ":"} {
			variants = append(variants, struct{ name, kind, value string }{name + operator, "credential-assignment", name + " " + operator + " \"" + strings.Repeat("T3", 10) + "\""})
		}
	}
	for _, variant := range variants {
		t.Run(variant.name, func(t *testing.T) {
			t.Parallel()
			matches := secretscan.ScanContent("variants.txt", []byte(variant.value))
			if len(matches) != 1 || matches[0].Kind != variant.kind {
				t.Fatalf("matches = %+v, want kind %s", matches, variant.kind)
			}
			assertNoSampleRun(t, (&secretscan.MatchesError{Matches: matches}).Error(), variant.value)
		})
	}
}

func TestScanContentReportsPathLineAndKind(t *testing.T) {
	t.Parallel()
	samples := secretSamples()
	content := "clean\n" + samples[2].value + " " + samples[1].value + "\n\n" + samples[2].value
	want := []secretscan.Match{
		{Path: "qa/evidence.txt", Line: 2, Kind: "github-token"},
		{Path: "qa/evidence.txt", Line: 4, Kind: "openrouter-key"},
	}
	matches := secretscan.ScanContent("qa/evidence.txt", []byte(content))
	if !reflect.DeepEqual(matches, want) {
		t.Fatalf("matches = %+v, want %+v", matches, want)
	}
	if got := (&secretscan.MatchesError{Matches: matches}).Error(); got != "qa/evidence.txt:2: github-token\nqa/evidence.txt:4: openrouter-key" {
		t.Fatalf("error = %q", got)
	}
	// A line exceeding bufio.Scanner's default limit must still be inspected.
	matches = secretscan.ScanContent("long.txt", []byte(strings.Repeat(" ", 70000)+samples[2].value))
	if len(matches) != 1 || matches[0].Line != 1 {
		t.Fatalf("long-line matches = %+v", matches)
	}
}

func TestScanSkipsBinaryAndPlaceholders(t *testing.T) {
	t.Parallel()
	samples := secretSamples()
	cases := []struct{ name, content string }{
		{"empty", ""},
		{"binary", "\x00\n" + samples[2].value},
		{"binary boundary", strings.Repeat(" ", 7999) + "\x00\n" + samples[2].value},
		{"variable", "TOKEN" + "=$VAR"},
		{"angle", "API_KEY" + "=<value>"},
		{"redacted", "SECRET" + "=[redacted]"},
		{"workflow placeholder", "PASSWORD" + "=${{ secrets.X }}"},
		{"short", "TOKEN" + "=" + strings.Repeat("a1", 9)},
		{"letters only", "TOKEN" + "=" + strings.Repeat("z", 25)},
		{"digits only", "TOKEN" + "=" + strings.Repeat("1", 25)},
		{"starts with equals", "TOKEN" + "==" + strings.Repeat("a1", 12)},
		{"lowercase name", "service_token" + "=" + strings.Repeat("a1", 12)},
		{"prefixed token", "prefix_" + samples[1].value},
		{"suffixed fixed token", samples[1].value + "_suffix"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := secretscan.ScanContent("ignored.txt", []byte(test.content)); len(got) != 0 {
				t.Fatalf("unexpected matches: %+v", got)
			}
		})
	}
	matches := secretscan.ScanContent("text.txt", []byte(strings.Repeat(" ", 8000)+"\x00\n"+samples[2].value))
	if len(matches) != 1 || matches[0].Line != 2 {
		t.Fatalf("NUL beyond binary window matches = %+v", matches)
	}
	// An invalid first assignment must not consume the delimiter that a later
	// assignment needs, even when a single space separates them.
	for _, separator := range []string{"; ", " "} {
		matches = secretscan.ScanContent("assignments.txt", []byte("TOKEN"+"="+strings.Repeat("z", 25)+separator+samples[9].value))
		if len(matches) != 1 || matches[0].Kind != "credential-assignment" {
			t.Fatalf("later assignment matches = %+v", matches)
		}
	}
}

func TestScanPushedHistoryReportsCommitPathAndLine(t *testing.T) {
	t.Parallel()
	wantArgs := []string{"log", "-p", "-U0", "--reverse", "--no-color", "--no-ext-diff", "--no-renames", "--format=commit %H", "candidate", "--not", "--remotes=origin"}
	if got := secretscan.PushedHistoryArgs("candidate", "origin"); !reflect.DeepEqual(got, wantArgs) {
		t.Fatalf("args = %v, want %v", got, wantArgs)
	}
	sample := secretSamples()[2].value
	first, second, third := strings.Repeat("a", 40), strings.Repeat("b", 40), strings.Repeat("c", 40)
	// Synthetic output of API Contract 3. The second commit removes the value,
	// and the third is a merge with no patch, so neither contributes a Match.
	log := "commit " + first + "\n\ndiff --git a/qa/evidence.txt b/qa/evidence.txt\n--- /dev/null\n+++ b/qa/evidence.txt\n@@ -0,0 +1,2 @@\n+clean\n+" + sample +
		"\n\\ No newline at end of file\n\ncommit " + second + "\n\ndiff --git a/qa/evidence.txt b/qa/evidence.txt\n--- a/qa/evidence.txt\n+++ b/qa/evidence.txt\n@@ -2 +1,0 @@\n-" + sample + "\n\ncommit " + third + "\n"
	want := []secretscan.Match{{Path: "qa/evidence.txt", Line: 2, Kind: "openrouter-key", Commit: first}}
	if got := secretscan.ScanPushedHistory([]byte(log)); !reflect.DeepEqual(got, want) {
		t.Fatalf("matches = %+v, want %+v", got, want)
	}
	if got := want[0].String(); got != "qa/evidence.txt:2: openrouter-key (commit "+first[:12]+")" {
		t.Fatalf("commit diagnostic = %q", got)
	}
	assertNoSampleRun(t, (&secretscan.MatchesError{Matches: want}).Error(), sample)
	t.Run("patch syntax and paths", func(t *testing.T) {
		t.Parallel()
		patch := "commit " + first + "\n\ndiff --git a/old b/new\n--- a/old\n+++ \"b/qa/a\\tfile.txt\"\n@@ -3,2 +5,3 @@ function\n-old\n+new\n+" + sample + "\n unchanged\n@@ -10,0 +12 @@\n+++" + sample +
			"\ndiff --git a/binary b/binary\nBinary files /dev/null and b/binary differ\ndiff --git a/deleted b/deleted\n--- a/deleted\n+++ /dev/null\n@@ -1 +0,0 @@\n-" + sample +
			"\ndiff --git a/space b/space\n--- /dev/null\n+++ b/qa/a space.txt\t\n@@ -0,0 +1 @@\n+" + sample + "\n"
		want := []secretscan.Match{
			{Path: "qa/a\tfile.txt", Line: 6, Kind: "openrouter-key", Commit: first},
			{Path: "qa/a\tfile.txt", Line: 12, Kind: "openrouter-key", Commit: first},
			{Path: "qa/a space.txt", Line: 1, Kind: "openrouter-key", Commit: first},
		}
		if got := secretscan.ScanPushedHistory([]byte(patch)); !reflect.DeepEqual(got, want) {
			t.Fatalf("matches = %+v, want %+v", got, want)
		}
	})
}

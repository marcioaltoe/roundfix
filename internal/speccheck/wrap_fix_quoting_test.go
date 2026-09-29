package speccheck_test

// Suite: wrap-fragile remediation shell quoting
// Invariant: the suggested grep receives the detected phrase and Markdown path as literal operands.
// Boundary IN: WrapFragileVerification remediation text executed by a POSIX shell.
// Boundary OUT: general shell parsing and Task Verification execution.

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrapFragileRemediationQuotesAHostilePhrase(t *testing.T) {
	directory := t.TempDir()
	substitutionMarker := filepath.Join(directory, "substitution-ran")
	backtickMarker := filepath.Join(directory, "backtick-ran")
	phrase := "literal $(touch " + substitutionMarker + ") `touch " + backtickMarker + "`; spaced ' quote"
	path := filepath.Join(directory, "hostile.md")
	if err := os.WriteFile(path, []byte(phrase+"\n"), 0o644); err != nil {
		t.Fatalf("write hostile-phrase fixture: %v", err)
	}

	findings := wrapFragileFindings(`grep -q "` + phrase + `" "` + path + `"`)
	if len(findings) != 1 {
		t.Fatalf("WrapFragileVerification() = %#v, want one finding", findings)
	}
	fix := findings[0].Fix
	if output, err := exec.Command("sh", "-c", fix).CombinedOutput(); err != nil {
		t.Fatalf("sh -c suggested remediation: %v; output=%q; fix=%q", err, output, fix)
	}
	for _, marker := range []string{substitutionMarker, backtickMarker} {
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("suggested remediation executed phrase content and created %q: %v", marker, err)
		}
	}

	if err := os.WriteFile(path, []byte(strings.TrimSuffix(phrase, "e")+"\n"), 0o644); err != nil {
		t.Fatalf("write near-match fixture: %v", err)
	}
	if output, err := exec.Command("sh", "-c", fix).CombinedOutput(); err == nil {
		t.Fatalf("suggested remediation matched a different phrase; output=%q; fix=%q", output, fix)
	}
}

func TestWrapFragileRemediationQuotesAPathWithASpace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "guide with space.md")
	if err := os.WriteFile(path, []byte(wrapFragilePhrase+"\n"), 0o644); err != nil {
		t.Fatalf("write spaced-path fixture: %v", err)
	}

	findings := wrapFragileFindings(`grep -q "` + wrapFragilePhrase + `" "` + path + `"`)
	if len(findings) != 1 {
		t.Fatalf("WrapFragileVerification() = %#v, want one finding", findings)
	}
	fix := findings[0].Fix
	if !strings.Contains(fix, "< '"+path+"'") {
		t.Fatalf("suggested remediation = %q, want single-quoted file operand %q", fix, path)
	}
	if output, err := exec.Command("sh", "-c", fix).CombinedOutput(); err != nil {
		t.Fatalf("sh -c suggested remediation: %v; output=%q; fix=%q", err, output, fix)
	}
}

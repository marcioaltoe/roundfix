package secretscan_test

import (
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/secretscan"
)

func TestScanAddedLinesReportsOnlyAddedLines(t *testing.T) {
	t.Parallel()
	sample := "SERVICE_API_KEY" + " = " + strings.Repeat("K8", 15)
	patch := "diff --git a/fixture.txt b/fixture.txt\n--- a/fixture.txt\n+++ b/fixture.txt\n@@ -8,2 +12,2 @@\n " + sample + "\n-" + sample + "\n+" + sample + "\n"
	want := []secretscan.Match{{Path: "fixture.txt", Line: 13, Kind: "credential-assignment"}}
	if got := secretscan.ScanAddedLines([]byte(patch)); !reflect.DeepEqual(got, want) {
		t.Fatalf("matches = %+v, want %+v", got, want)
	}
}

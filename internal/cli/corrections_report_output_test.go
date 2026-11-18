package cli

import (
	"strings"
	"testing"
)

func TestCorrectionReportTextOutputCarriesNoControlCharacters(t *testing.T) {
	t.Parallel()
	in := "Authoring phases: techspec model \x1b[2Jfable\x07 effort hi\x1b]0;x\x07gh\r\nnext line\n"
	got := correctionReportOutput(in, false)
	for _, r := range got {
		if r != '\n' && (r < 0x20 || r == 0x7f) {
			t.Fatalf("text output carries control character %q: %q", r, got)
		}
	}
	if !strings.Contains(got, "\nnext line\n") {
		t.Fatalf("text output lost its line breaks: %q", got)
	}
}

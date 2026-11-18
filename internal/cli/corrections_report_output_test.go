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

func TestCorrectionReportTextKeepsARecordedModelOnOneLine(t *testing.T) {
	t.Parallel()
	forged := "fable\nComparison, techspec, feat: forged line"
	r := correctionsReport{
		Specs:       []correctionReportSpec{{Spec: "0001-x", Phases: map[string]correctionReportPhase{"techspec": {Model: forged, Effort: "high\nforged"}}}},
		Groups:      []correctionReportGroup{{Phase: "techspec", Type: "feat", Model: forged, Effort: "high\rforged"}},
		ModelTotals: []correctionReportGroup{{Phase: "techspec", Type: "feat", Model: forged}},
		Comparisons: []correctionReportComparison{{Phase: "techspec", Type: "feat", Groups: []correctionReportGroup{{Model: forged}}}},
	}
	clean := cleanCorrectionsReportText(r)
	for _, value := range []string{clean.Specs[0].Phases["techspec"].Model, clean.Specs[0].Phases["techspec"].Effort, clean.Groups[0].Model, clean.Groups[0].Effort, clean.ModelTotals[0].Model, clean.Comparisons[0].Groups[0].Model} {
		if strings.ContainsAny(value, "\n\r") {
			t.Fatalf("recorded value keeps a line break: %q", value)
		}
	}
	if r.Groups[0].Model != forged {
		t.Fatal("cleaning the text view changed the report the JSON output encodes")
	}
}

package speccheck_test

import (
	"testing"

	"roundfix/internal/speccheck"
)

func TestRefusalReasonIsTheReasonTheGateRecords(t *testing.T) {
	t.Parallel()
	for _, test := range []struct{ name, code, summary, want string }{
		{"code and sentence", "SC-TEST", "A contradiction.", "SC-TEST: A contradiction."},
		{"whitespace", " SC-TEST\n", " A\n contradiction.\t", "SC-TEST: A contradiction."},
		{"code only", "SC-TEST", "", "SC-TEST"},
		{"sentence only", "", "A contradiction.", "A contradiction."},
		{"unnamed", "", "", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			finding := speccheck.Finding{Code: test.code, Summary: test.summary}
			reason := speccheck.RefusalReason(finding)
			refusal, refused := speccheck.PreconditionRefusal(speccheck.GatePreconditionResult{Blocking: true, Findings: []speccheck.Finding{finding, finding}})
			if !refused || reason != test.want || refusal.Reason != test.want {
				t.Fatalf("reason %q, refusal %+v, refused %v, want %q", reason, refusal, refused, test.want)
			}
		})
	}
}

package reviewsource

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestEvidenceJSONMappingAndBounds(t *testing.T) {
	states := []EvidenceState{
		EvidencePending,
		EvidenceReviewing,
		EvidenceReviewed,
		EvidenceVerified,
		EvidenceSkipped,
		EvidenceFailed,
	}
	wantStates := []string{"pending", "reviewing", "reviewed", "verified", "skipped", "failed"}
	for index, state := range states {
		if string(state) != wantStates[index] {
			t.Fatalf("state %d = %q, want %q", index, state, wantStates[index])
		}
	}

	evidence := Evidence{
		State:           EvidenceSkipped,
		Kind:            EvidenceKindCheckRun,
		Identity:        "check_run:42",
		ExpectedHeadSHA: "expected",
		ObservedHeadSHA: "observed",
		Conclusion:      "success",
		Detail:          BoundEvidenceDetail(strings.Repeat("é", MaxEvidenceDetailLength)),
		Reason:          "review skipped because the change set is too large",
	}
	encoded, err := json.Marshal(evidence)
	if err != nil {
		t.Fatalf("marshal Evidence: %v", err)
	}
	if len(evidence.Detail) > MaxEvidenceDetailLength+len("…") {
		t.Fatalf("bounded detail length = %d, want at most %d plus ellipsis", len(evidence.Detail), MaxEvidenceDetailLength)
	}

	var mapped map[string]any
	if err := json.Unmarshal(encoded, &mapped); err != nil {
		t.Fatalf("decode Evidence JSON: %v", err)
	}
	for _, field := range []string{
		"state",
		"kind",
		"identity",
		"expected_head_sha",
		"observed_head_sha",
		"conclusion",
		"detail",
		"reason",
	} {
		if _, ok := mapped[field]; !ok {
			t.Fatalf("Evidence JSON missing %q: %s", field, encoded)
		}
	}
	if strings.Contains(string(encoded), "provider_response") {
		t.Fatalf("Evidence JSON leaked provider response field: %s", encoded)
	}
}

func TestBoundEvidenceDetailCutsOnRuneBoundary(t *testing.T) {
	bounded := BoundEvidenceDetail(strings.Repeat("→", MaxEvidenceDetailLength))
	if !utf8.ValidString(bounded) {
		t.Fatalf("bounded detail is not valid UTF-8: %q", bounded)
	}
	if len(bounded) > MaxEvidenceDetailLength+len("…") {
		t.Fatalf("bounded detail length = %d, want at most %d plus ellipsis", len(bounded), MaxEvidenceDetailLength)
	}
}

func TestTransientErrorPreservesOperationAndWrappedCause(t *testing.T) {
	cause := errors.New("temporary failure containing credential ghp_not_for_output")
	transient := &TransientError{Operation: "fetch Review Source evidence", Err: cause}
	wrapped := fmt.Errorf("observe expected head: %w", transient)

	if !IsTransient(wrapped) {
		t.Fatal("wrapped TransientError was not discoverable")
	}
	if !errors.Is(wrapped, cause) {
		t.Fatal("TransientError did not preserve its wrapped cause")
	}
	if transient.Operation != "fetch Review Source evidence" {
		t.Fatalf("operation = %q", transient.Operation)
	}
	if strings.Contains(wrapped.Error(), "ghp_not_for_output") {
		t.Fatalf("TransientError exposed its wrapped cause: %q", wrapped)
	}
}

func TestTransientErrorNilAndPermanentErrors(t *testing.T) {
	if IsTransient(nil) {
		t.Fatal("nil error classified transient")
	}
	if IsTransient(errors.New("permanent validation failure")) {
		t.Fatal("plain error classified transient")
	}
}

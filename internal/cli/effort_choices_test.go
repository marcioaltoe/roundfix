// Suite: runtime picker efforts.
// Invariant: the picker offers the recorded per-runtime effort lists.
// Boundary IN: picker data; OUT: per-model live capability proof.
package cli

import (
	"reflect"
	"testing"
)

func TestReasoningEffortChoicesFollowTheAdapters(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		runtime string
		want    []string
	}{
		{"codex", []string{"low", "medium", "high", "xhigh", "max"}},
		{"claude", []string{"default", "low", "medium", "high", "xhigh", "max"}},
		{"opencode", nil}, {"unknown", nil}, {"", nil},
	} {
		t.Run(tt.runtime, func(t *testing.T) {
			if got := reasoningEffortChoices(tt.runtime); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("efforts = %v, want %v", got, tt.want)
			}
		})
	}
}

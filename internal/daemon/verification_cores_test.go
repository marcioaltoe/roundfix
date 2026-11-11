package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/corebudget"
)

func TestTheVerificationCommandReceivesTheCoreBudget(t *testing.T) {
	t.Parallel()
	goValue, goPresent := os.LookupEnv(corebudget.GoVariable)
	if !goPresent {
		goValue = "unset"
	}
	roundfixValue, roundfixPresent := os.LookupEnv(corebudget.Variable)
	if !roundfixPresent {
		roundfixValue = "unset"
	}
	goCount := 0
	if goPresent {
		goCount = 1
	}
	for _, test := range []struct {
		name  string
		cores int
		want  string
	}{
		{"budget four", 4, "4|4|1"},
		{"no budget", 0, fmt.Sprintf("%s|%s|%d", goValue, roundfixValue, goCount)},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			_, err := (ExecVerifier{}).Verify(t.Context(), VerifyRequest{
				RunID:      "run-verification-core-budget",
				Cores:      test.cores,
				WorkDir:    dir,
				Command:    `printf '%s|%s|%s' "${GOMAXPROCS-unset}" "${ROUNDFIX_VERIFICATION_CORES-unset}" "$(env | awk '/^GOMAXPROCS=/ { count++ } END { print count+0 }')" > cores.txt`,
				OutputPath: filepath.Join(dir, "verification.log"),
			})
			if err != nil {
				t.Fatalf("Verify: %v", err)
			}
			got, err := os.ReadFile(filepath.Join(dir, "cores.txt"))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != test.want {
				t.Errorf("command environment = %q, want %q", got, test.want)
			}
		})
	}
}

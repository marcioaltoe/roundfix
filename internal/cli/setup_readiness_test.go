// Suite: Setup readiness reporting through the public CLI.
// Boundary OUT: machine probes and setup writes use the existing fake dependencies.
package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestSetupPrintsTheReadinessLinesAndOffersNothing(t *testing.T) {
	for _, acpxReady := range []bool{true, false} {
		for _, status := range []CheckStatus{CheckStatusOK, CheckStatusWarn, CheckStatusFailed} {
			for _, flag := range []string{"", "--yes", "--no-input"} {
				t.Run(string(status)+"/"+flag+map[bool]string{true: "/acpx-ready", false: "/acpx-missing"}[acpxReady], func(t *testing.T) {
					makeFake := func() *setupFakeDeps {
						fake := newSetupFakeDeps()
						if !acpxReady {
							fake.acpxErr = errors.New("missing acpx")
							fake.confirm = func(context.Context, io.Writer, string) (bool, error) { return false, nil }
						}
						return fake
					}
					baseline := makeFake()
					withSetupFakeDeps(t, baseline)
					args := []string{"setup"}
					if flag != "" {
						args = append(args, flag)
					}
					var baseOut, baseErr bytes.Buffer
					baseCode := runCLI(t, args, &baseOut, &baseErr)
					fake := makeFake()
					withSetupFakeDeps(t, fake)
					results := []CheckResult{
						{Name: "gh", Status: status, Detail: "DR-GH-UNAUTHENTICATED: gh has no account for github.com", NextAction: "gh auth login --hostname github.com"},
						{Name: "git", Status: CheckStatusOK, Detail: "identity set"},
						{Name: "remote", Status: CheckStatusOK, Detail: "reachable"},
						{Name: "toolchain", Status: CheckStatusOK, Detail: "tools found"},
						{Name: "environment", Status: CheckStatusWarn, Detail: "optional key not set", NextAction: "set the optional key"},
					}
					calls := 0
					updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
						deps.setup.readiness = func(_ context.Context, loaded roundconfig.Loaded) []CheckResult {
							calls++
							if loaded.GitRoot != fake.gitRoot {
								t.Fatalf("readiness repository=%q", loaded.GitRoot)
							}
							return results
						}
					})
					var out, diag bytes.Buffer
					code := runCLI(t, args, &out, &diag)
					wantCode := baseCode
					if status == CheckStatusFailed {
						wantCode = exitRunFailed
					}
					if code != wantCode || calls != 1 || diag.String() != baseErr.String() {
						t.Fatalf("exit=%d want=%d calls=%d stderr=%q", code, wantCode, calls, diag.String())
					}
					var block bytes.Buffer
					for _, result := range results {
						printDoctorResult(&block, result)
					}
					lines := strings.SplitAfter(baseOut.String(), "\n")
					inserted := false
					for index, line := range lines {
						if strings.HasPrefix(line, "acpx:") {
							lines[index] += block.String()
							inserted = true
							break
						}
					}
					if !inserted || out.String() != strings.Join(lines, "") {
						t.Fatalf("readiness order/text: stdout=%q baseline=%q block=%q", out.String(), baseOut.String(), block.String())
					}
					if !reflect.DeepEqual(fake.prompts, baseline.prompts) || !reflect.DeepEqual(fake.installCalls, baseline.installCalls) || !reflect.DeepEqual(fake.writeCalls, baseline.writeCalls) || !reflect.DeepEqual(fake.files, baseline.files) {
						t.Fatal("readiness changed setup offers or writes")
					}
				})
			}
		}
	}
}

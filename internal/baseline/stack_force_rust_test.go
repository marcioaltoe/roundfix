// Suite: Rust clause force and error policy
// Invariant: Rust rules carry force and rendered guidance bounds type-erased errors and panics.
// Boundary IN: embedded catalog and a Rust CLI plan for a temporary adopter.
// Boundary OUT: applying a repository refresh and Daemon Verification.

package baseline

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func TestEveryRustRuleCarriesForce(t *testing.T) {
	t.Parallel()
	asset, ok := mustEmbeddedCatalog(t).Module("rust")
	if !ok {
		t.Fatal("missing Rust module")
	}
	var module map[string]any
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}
	if findings := ruleLevelGuidanceFindings("rust", module); len(findings) != 0 {
		t.Fatalf("Rust rule-level guidance findings: %v", findings)
	}
}

func TestTheRustGuideStatesTheErrorPolicy(t *testing.T) {
	t.Parallel()
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   newPlanRepository(t),
		ProfileID:    "rust-cli",
		Decisions:    planTestDecisions(),
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("Rust CLI outcome = %+v, want ready plan", outcome.Result)
	}
	guide := string(planPostimage(t, *outcome.Plan, "docs/agents/rust.md").Content)
	want := []string{
		"These rules govern the repository's Rust crates: the library and binary targets of its Cargo package or workspace, and their tests. Code in another language follows its own guide.",
		"- **mandatory**: Change dependencies through Cargo commands such as `cargo add`, `cargo remove`, and `cargo update`, and commit `Cargo.toml` and `Cargo.lock` together.",
		"- **mandatory**: Use current authoritative crate and toolchain documentation before changing Rust APIs, configuration, or dependencies.",
		"- **mandatory**: Use focused `cargo check`, `cargo test`, and lint commands while iterating, then run the selected repository Verification.",
		"- **mandatory**: Treat CLI flags, streams, JSON fields, and exit codes as public behavior and test them through observable command execution.",
		"- **mandatory**: Keep each binary target's `main` thin: it parses arguments and calls behavior that lives in the library crate.",
		"- **prohibited**: Do not use `anyhow` or another type-erased error type outside a binary's entry point. Only `main` and the command dispatch it calls may turn typed errors into a type-erased report.",
		"- **prohibited**: Do not call `unwrap`, `expect`, or `panic!` on a path that user input, the environment, a file, or an IO result can reach; return a typed error instead. Test code is not such a path.",
		"- **mandatory**: Return typed errors from library and domain code: an error enum per boundary that implements `std::error::Error`, for example through `thiserror`, so a caller can match its variants.",
	}
	if findings := stackWordingFindings(map[string]string{"docs/agents/rust.md": guide}, []stackWordingRule{{
		guide: "docs/agents/rust.md", must: want,
		mustNot: []string{"Preserve Cargo manifest and lockfile discipline."},
	}}); len(findings) != 0 {
		t.Fatalf("Rust guide wording findings: %v", findings)
	}
	for _, line := range strings.Split(guide, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") && !strings.HasPrefix(line, "- **mandatory**: ") && !strings.HasPrefix(line, "- **prohibited**: ") {
			t.Errorf("Rust rule line has no force: %s", line)
		}
	}
}

func TestTheRustClausesCarryTheirForce(t *testing.T) {
	t.Parallel()
	want := []expectedClause{
		{"rust", "clause.rust.change-manifests-through-cargo", "mandatory"},
		{"rust", "clause.rust.read-current-docs", "mandatory"},
		{"rust", "clause.rust.iterate-with-focused-cargo-commands", "mandatory"},
		{"rust", "clause.rust.test-the-command-through-execution", "mandatory"},
		{"rust", "clause.rust.keep-the-binary-thin", "mandatory"},
		{"rust", "clause.rust.keep-type-erased-errors-at-the-entry-point", "prohibited"},
		{"rust", "clause.rust.prohibit-panics-on-user-paths", "prohibited"},
		{"rust", "clause.rust.return-typed-errors", "mandatory"},
	}
	if findings := clauseForceFindings(mustEmbeddedCatalog(t), want); len(findings) != 0 {
		t.Fatalf("Rust clause force findings: %v", findings)
	}
}

func TestTheRustProfileRequiresTheErrorRule(t *testing.T) {
	t.Parallel()
	asset, ok := mustEmbeddedCatalog(t).Profile("rust-cli")
	if !ok {
		t.Fatal("missing Rust CLI profile")
	}
	var profile struct {
		RequiredRules []string `json:"requiredRules"`
	}
	if err := json.Unmarshal(asset.Data, &profile); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(profile.RequiredRules, "rule.rust.library-and-entry-point") {
		t.Fatal("Rust CLI profile does not require the error rule")
	}
}

---
task: task_01
spec: 0282-the-linux-lane-runs-like-ci
status: pending
type: test
complexity: medium
---

# Task 01: A test that needs an External Test Tool skips naming it

## Overview

On 2026-10-18 the Linux Lane failed `TestGuidanceCompositionJourney/rust-cli`
because the `golang:1.26` image has no `rustfmt` (`_prd.md` → Sources 1). A
machine without `rustfmt` fails it the same way today. This Task adds
`testfixture.RequireTool`, which skips a test naming the tool it lacks, moves
the two tests that look up a non-base tool onto it, and adds a repository
rule that refuses a new literal lookup of such a tool. It is verifiable on its
own: the rule and the helper have their own tests, and the rust-cli journey
skips under a `PATH` without `rustfmt`.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-18,
   `references/2026-10-18-the-linux-lane-container-lacks-ci-tools-and-an-init.md`
   (`_prd.md` → Sources 1), for its first failure, by implementing API
   Contract 5 of `_techspec.md` (ADR-0294).
2. MUST add `RequireTool(t testing.TB, name string) string` to the new
   `internal/testfixture/external_tool.go`. It returns the path
   `exec.LookPath` finds and otherwise calls `t.Skipf("requires %s, which is
   not on PATH", name)`. The data source is the test process's own `PATH`;
   the reader of the skip message is whoever reads the test output, so it
   names the tool and nothing else.
3. MUST add, in the new `internal/testfixture/external_tool_test.go`, the
   External Test Tool rule and its tests:
   - `TestRequireToolReturnsAToolOnPath`: with `PATH` set to a temporary
     directory holding one compiled stand-in, the helper returns its path;
   - `TestRequireToolSkipsNamingAnAbsentTool`: with `PATH` set to an empty
     temporary directory, the helper skips with exactly `requires
     rustfmt, which is not on PATH`, and no statement after it runs;
   - `TestEveryExternalTestToolGoesThroughRequireTool`: the rule over the
     repository reports nothing;
   - `TestTheExternalTestToolRuleReportsALiteralLookup`: in a seeded
     temporary tree, `exec.LookPath("node")`, `exec.Command("rustfmt")` and
     an aliased `os/exec` import each produce one report with the file and
     line, while `exec.Command("git")`, `exec.Command("/bin/sh")` and a
     program held in a variable produce none.

   The rule follows API Contract 5: every `_test.go` under `internal/`,
   `cmd/` and `skills/`, the Base Test Tools listed there, and the report
   text `<file>:<line>: <tool> is not a Base Test Tool; resolve it with
   testfixture.RequireTool`. The two tests that set `PATH` carry a `//
   Sequential: sets PATH` comment instead of `t.Parallel()`; the others call
   `t.Parallel()` first. The file compiles for Windows (ADR-0285).
4. MUST make the rust-cli step of `runBaselineReleaseFormatter` in
   `internal/cli/baseline_release_gate_test.go` resolve `rustfmt` with
   `testfixture.RequireTool` before formatting, so the steps before it still
   run and the subtest skips with the helper's message when `rustfmt` is
   absent. The other Profiles' steps MUST stay unchanged, and a missing
   `gofmt` or `make` MUST still fail.
5. MUST make `internal/agent/acpx_integration_unix_test.go` resolve `node`
   with `testfixture.RequireTool` after its opt-in check, replacing its
   direct `exec.LookPath("node")`.
6. MUST add to the repository rules in `docs/agents/specific-repository.md`,
   outside any setup marker, one rule: a test that runs an External Test
   Tool, a tool outside the Base Test Tools, resolves it with
   `testfixture.RequireTool`, which skips naming the tool when it is absent,
   and the External Test Tool rule in `internal/testfixture` enforces it.
7. MUST add **External Test Tool** to `CONTEXT.md` through domain-modeling,
   naming the Base Test Tools and linking ADR-0294.

## Subtasks

- [ ] Add `RequireTool`.
- [ ] Add the External Test Tool rule and the four tests.
- [ ] Move the rust-cli journey and the acpx integration test onto the helper.
- [ ] Add the repository rule and the glossary term.

## Acceptance Criteria

- [ ] Under a `PATH` without `rustfmt`, the rust-cli journey skips with `requires rustfmt, which is not on PATH`.
- [ ] No test in the repository resolves a tool outside the Base Test Tools by a literal name except through `RequireTool`.
- [ ] The repository rule and the glossary describe External Test Tools.

## Consumers

- `testfixture.RequireTool` — consumers: `./internal/cli`, `./internal/agent`

## Context

- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `internal/testfixture/written_executable_test.go`
- instruction: `internal/testfixture/fixture.go`
- creates: `internal/testfixture/external_tool.go`
- creates: `internal/testfixture/external_tool_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`
- interface: `internal/agent/acpx_integration_unix_test.go`
- interface: `docs/agents/specific-repository.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestRequireToolReturnsAToolOnPath|TestRequireToolSkipsNamingAnAbsentTool|TestEveryExternalTestToolGoesThroughRequireTool|TestTheExternalTestToolRuleReportsALiteralLookup)$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRequireToolReturnsAToolOnPath TestRequireToolSkipsNamingAnAbsentTool TestEveryExternalTestToolGoesThroughRequireTool TestTheExternalTestToolRuleReportsALiteralLookup; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; bin="$(mktemp -d)" || exit 1; for tool in go git sh; do found="$(command -v "$tool")" || exit 1; ln -s "$found" "$bin/$tool" || exit 1; done; out="$(PATH="$bin" "$bin/go" test -count=1 -v -run '^TestGuidanceCompositionJourney$/^rust-cli$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- SKIP: TestGuidanceCompositionJourney/rust-cli ' || { printf 'rust-cli did not skip\n' >&2; exit 1; }; printf '%s\n' "$out" | grep -qF -- 'requires rustfmt, which is not on PATH' || { printf 'skip does not name rustfmt\n' >&2; exit 1; }; go vet ./internal/agent || exit 1; GOOS=windows go vet ./internal/testfixture || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**External Test Tool**:' || { printf 'missing glossary term External Test Tool\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- 'testfixture.RequireTool' || { printf 'missing repository rule\n' >&2; exit 1; }` — expected: exit 0. Before this Task the four tests do not exist, so the command fails at its first pass check; the rust-cli journey also fails there today with `exec: "rustfmt": executable file not found`. Afterwards the helper and rule tests pass, the journey skips naming `rustfmt`, both callers compile on the host and the helper's package for Windows, and the repository rule and the glossary carry the term.

## References

- `_prd.md` → Sources; Goals; User Story 3; Core Feature 2; Core Feature 5; Success Metric 2; Glossary
- `_techspec.md` → API Contract 5; API Contract 7; Invariant 3; Current behavior; Testing Approach; Build Order 1
- ADR-0294; ADR-0125; ADR-0259; ADR-0285

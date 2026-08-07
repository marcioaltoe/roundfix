# Tooling authorization — make the verification gate honest (2026-08-07)

Spec: `0083-a-gate-that-can-say-no`.

On 2026-08-07 the maintainer said the gates had become a burden —

> Esses gates de testes do roundfix estão atrapalhando demais. Está sendo um
> sofrimento passar pelo QA

— reviewed the tally of which gates produced signal and which produced noise,
and authorized the bounded tooling to repair both the dishonesty and the noise:

> Conceder o conjunto completo.

and, for the coverage record trapped inside an archived Spec:

> Sair para um dono semântico.

## What this covers

Three defects, established by reproduction on 2026-08-07:

1. **The gate certifies red builds.** `make verify` exits `0` on a tree whose
   `go test -parallel 16 ./...` exits `1`, and the wrapper omits the failing
   package from its summary. The Makefile routes the authoritative test
   invocation through the wrapper via `GO := $(RTK) go`. The masking is not
   universal — it reproduced on the full package set with a test emitting 278
   log lines beside 2 errors, and did not reproduce on a two-package set — and
   the trigger is **not established**. The fix does not depend on establishing
   it: the authoritative gate must not route through an output-filtering tool.
2. **Two gates assert facts about the machine, not the code.** The corpus golden
   pins a global finding count that changes whenever anyone adds an ADR or a
   Spec, and the corpus budget asserts wall-clock under one second on a shared
   developer machine. Both fired on 2026-08-07 with zero signal.
3. **The live coverage contract lives in an archived Spec.** Every legitimate
   test rename must rewrite
   `docs/specs/_archived/0071-verification-cost/coverage-record.json`, which the
   archived-Spec rule forbids.

## Authorized paths

- `Makefile` — to keep the authoritative verification invocation off the
  output-filtering wrapper. The wrapper stays available for human-facing
  convenience targets.
- `internal/spec/gate_test.go` — new file, the regression test proving the
  authoritative invocation propagates a non-zero exit. It joins the package that
  already owns the repository's verification contracts.
- `internal/speccheck/constraints_characterization_test.go` and
  `internal/speccheck/testdata/corpus-golden.json` — to retire or rebuild the
  global counter and the wall-clock budget so neither asserts a fact about the
  machine or about unrelated authoring.
- `internal/spec/coverage_test.go` — to point the coverage contract at its new
  semantic owner.
- `internal/agent/acpx_runner_test.go` and `internal/cli/implement_test.go` — to
  stabilize the two timing-sensitive tests that flaked under load on 2026-08-07.

## Authorized move

- `docs/specs/_archived/0071-verification-cost/coverage-record.json` →
  `docs/references/coverage-record.json`.

This is a deliberate, maintainer-approved exception to *keep archived Specs
byte-identical*, and it is the only one this grant permits. It resolves the
collision in the direction the repository's own rule already points: durable
knowledge a Spec produced moves to its semantic owner rather than staying in a
directory declared deletable.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic
consequences of the authorized edits, per ADR-0081.

## Not authorized

- Removing `TestCoverageEquivalence` or its record. The maintainer chose
  relocation over retirement; the invariant caught a real regression the same
  day and keeps its teeth.
- Weakening any gate that produced signal: `spec check`, the published-example
  contract, or the authored QA gate.
- Any other archived Spec content.

## Commit choreography

This record lands as its own commit, before the commit of any Task it
authorizes.

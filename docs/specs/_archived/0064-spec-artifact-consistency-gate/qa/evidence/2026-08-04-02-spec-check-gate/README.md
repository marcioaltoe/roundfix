# Evidence — 2026-08-04-02 Spec consistency gate

This index records the fresh commands and observed results cited by
`qa-report-2026-08-04-02.md`. It is populated as each QA row closes.

## E-01 — Constraint, graph, and tooling audit

- `rtk read` over `_prd.md`, `_techspec.md`, `_tasks.md`, every Task file,
  operative agent guides, accepted ADRs 0080, 0088, 0089, 0091, 0093, and
  0094, and the 2026-08-02 authorization: the QA node is terminal/runnable;
  both active artifacts carry complete, reasoned Project Constraints.
- `rtk git diff-tree --no-commit-id --name-only -r <commit>` over the
  authorization, authoring, implementation, corrective, and prior-finding
  repair commits: task 08 (`02c35e2`) and task 11 (`121fabd`) each changed
  only `Makefile` plus their assigned Task file; the authorization and
  corrective chronology predates those changes.
- `rtk git -c core.fsmonitor=false status --short`: only this report and this
  evidence directory are untracked.
- `rtk git diff --check`: exit 0.

## E-02 — Active and clean command runs

- `rtk make verify`: the final `bin/roundfix spec check` stage exited 0 and
  named all nine active Specs once with no findings.
- Direct `rtk bin/roundfix spec check`: exit 0 with the same nine clean Specs
  and presence-aware skips.
- Two direct named-Spec runs: exit 0 with identical `No findings` output for
  Spec 0064.

## E-03 — Error, gap, and Constraint fixtures

The built artifact was run from an isolated Git repository at
`/private/tmp/roundfix-0064-qa-02.VPmOZS/repo` on branch `ma/qa-fixture`.

- `tooling-unauthorized` twice: exit 1, identical
  `SC-TOOLING-UNAUTHORIZED` error with PRD and authorization locations.
- scratch `gap-only`: default 0/gap, strict 1/error, final default 0/gap.
- `constraint-missing`, `constraint-unreasoned`, `constraint-source`, and
  `tooling-unbounded`: exit 1 with their expected located diagnostics.
- `clean` and `no-techspec`: exit 0; the latter lists the missing-TechSpec
  detector skips.

## E-04 — ADR, coverage, and reference fixtures

- `citation-dirty`: exit 1 with located `SC-ADR-UNLISTED` and a separate
  `SC-ADR-RELATED` gap.
- `coverage-unmapped`: exit 1 with exactly four uncovered Core Features.
- `coverage-untasked`: exit 1 with Core Feature 4 as the sole missing Task
  reference.
- `reference-unresolved` and `reference-index-unresolved`: exit 1, each naming
  both the declaration and missing path.
- `coverage-range`: exit 0 with no findings.

## E-05 — Vocabulary fixtures

- `vocabulary-missing`: exit 1 and one finding for `publish:` despite repeated
  emission; its fix includes `emits`, `pattern`, and `documented-in`.
- `vocabulary-satisfied` and `vocabulary-none`: exit 0; the latter records the
  detector skip.
- Invalid RE2, missing emitter, and missing documentation fixtures: exit 1
  with located findings and no panic; every fix teaches the declaration shape.

## E-06 — Four report replays

- 0058 QA-001: `SC-COVERAGE-UNMAPPED`.
- 0058 QA-004: `SC-VOCABULARY-UNDOCUMENTED` for `publish:`.
- 0056 F-001: `SC-ADR-UNLISTED` plus `SC-ADR-RELATED`.
- 0056 F-002: `SC-COVERAGE-UNMAPPED`.
- Each fixture README names its source QA report and states that the shape was
  authored from Expected/Actual, not recovered from Git.

## E-07 — Budget, corpus, and full gate

- `rtk make verify`: exit 0; 3,199 tests, dedicated budget, skill checks,
  build, and all-active Spec check passed.
- `rtk make spec-budget`: exit 0 with one passing selected test.
- Verbose repository-local-cache equivalent: exit 0; full Spec corpus sweep
  measured 315.226291 ms against the unchanged 1 s budget.
- Focused golden and active-corpus tests: exit 0. Active counts are all zero;
  archived counts remain 35 related ADR, 35 unlisted ADR, 320 missing
  Constraint, 10 missing source, 424 unmapped, 30 untasked, 183 unresolved
  reference, and 10 unbounded tooling findings.
- One optional verbose command inherited an out-of-sandbox Go cache and was
  denied before setup; the public target and repository-local equivalent both
  passed, so no journey is blocked.

## E-08 — JSON stability and CLI recovery

- Two explicit 0064+0075 JSON runs: exit 0 and byte-identical by `cmp`; `jq -e`
  parsed both objects as `roundfix-speccheck/v1` with the requested slugs.
- `--format json` before and after the slug produced byte-identical output.
- Unknown slug, unsupported format, and unknown flag: exit 2, named input on
  stderr, help pointer present, and captured stdout measured 0 bytes.
- Valid text and JSON commands immediately after the failures: exit 0.

## E-09 — Help and vocabulary

- `bin/roundfix --help` and `spec check --help`: command, formats, strict
  behavior, read-only boundary, and exits 0/1/2 are documented.
- `make help`: `spec-check` and `spec-budget` are listed.
- `CONTEXT.md`: Spec Consistency Check and Consistency Finding Severity remain
  distinct from QA verdict terminology.

## E-10 — Read-only persistence and Non-Goals

- Final named-Spec canary: exit 0 with `No findings` and the same expected
  skips as both initial runs.
- SHA-256 before and after all flows: `_prd.md`, `_techspec.md`, `_tasks.md`,
  and task 01 through task 11 matched byte-for-byte.
- `bin/roundfix implement --help`: no Spec Consistency Check precondition.
- Public checker/help output: no QA verdict, mutation, auto-correction, style
  judgment, or decision-correctness claim.
- Final pre-close status: only this report and evidence directory are
  untracked; `git diff --check` exits 0.

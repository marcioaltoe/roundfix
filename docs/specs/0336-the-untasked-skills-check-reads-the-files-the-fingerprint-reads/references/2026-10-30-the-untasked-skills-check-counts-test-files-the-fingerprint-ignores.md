---
type: fix
status: promoted
created: 2026-10-30
spec: 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
reason: null
---

# The untasked-skills check counts test files the fingerprint ignores

## Symptom

`roundfix spec check` reports `SC-SKILLS-UNTASKED` for a Task that changes
only a test file under a Behavior Surface's path pattern. The `config keys`
surface in `docs/references/skill-coverage.json` declares the pattern
`internal/config/*.go`, which also matches `internal/config/*_test.go`, while
the surface fingerprint that `TestTheSkillCoverageMapIsCurrent` computes skips
every `_test.go` file (`internal/docscontract/skill_coverage_test.go`, the
`strings.HasSuffix(entry.Name(), "_test.go")` filter). A test-only change
cannot move the fingerprint, yet the check asks the Task to record a Coverage
Review.

## Where

`internal/speccheck/skills_declaration.go` (`CodeSkillsUntasked`) and the
path patterns of `docs/references/skill-coverage.json`.

## Expected

The check and the fingerprint read the same file set: a file the fingerprint
ignores never triggers `SC-SKILLS-UNTASKED`.

## Evidence

Reported by the author of Spec
`0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec`
on 2026-10-30: its task_01 changes only `internal/config` tests and had to
add a Coverage Review edit to `skill-coverage.json` to pass the strict check.
The operator confirmed the pattern and the fingerprint filter on `35de7f74`.

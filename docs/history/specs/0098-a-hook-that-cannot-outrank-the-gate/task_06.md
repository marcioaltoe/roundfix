---
status: completed
type: docs
---

# Task: Document Hook Strictness Invariant

Write invariant into Baseline module.

## Work
- Add section to docs/agents/autonomous-work.md
- Text: "A commit hook must not be stricter"
- Use managed marker boundaries

## Verification
- `grep -q "commit hook must not be stricter" docs/agents/autonomous-work.md`


## References

- User Story 3: Invariant documented
- Core Feature 4: Invariant Documentation

## Result

The hook-strictness invariant is a Baseline clause, so a Baseline refresh
reproduces it instead of overwriting it.

### Implementation

- `internal/baseline/assets/modules/autonomous-work.json` — added
  `rule.autonomous.hook-strictness` (version 1, coverage
  `coverage.autonomous-work` + `coverage.verification`) with the mandatory
  clause `clause.autonomous.hook-strictness`. The clause opens with "A commit
  hook must not be stricter than the Task's declared Verification", cites
  ADR-0014 for the Daemon's verification authority, names `hook_refused` as a
  repository misconfiguration rather than a Task failure, and states that the
  work stays staged while the Run ends unresolved. The rule is listed in
  `guide.autonomous-work` between the verification-gate rule and the loop rule;
  module version 8 → 9 and guide version 8 → 9.
- `docs/agents/autonomous-work.md` — the rendered clause, inside the
  `setup-context-driven:begin/end id=guide.autonomous-work` managed region.
- `make baseline-digests` regenerated the derived artifacts (golden fixture,
  formatter goldenDigest, catalog digest/normalized/diagnostics, four
  plan-characterization goldens). Regeneration is the sanctioned flow recorded
  in `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`.

### Acceptance criterion 5 — invariant is written in rendered guidance

- **Appears in the repository's `docs/agents/autonomous-work.md`**: the string
  `commit hook must not be stricter` is present and sits between the managed
  begin and end markers (`present: True`, `inside managed region: True`).
- **Maintenance copy confirms the rule during Baseline refresh**: the renderer's
  own output carries the clause —
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
  gained the identical bullet when `make baseline-digests` regenerated it.
- **No local edit overwrites it**: the rendered guide and the regenerated golden
  differ from their prior versions by the same single bullet at the same
  position, so the clause is what a refresh produces, not a local addition on
  top of it.

### Focused checks

| Check | Result |
| --- | --- |
| `make baseline-digests` | `{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":true}` |
| `go test ./internal/baseline/... ./skills/... ./internal/speccheck/...` | ok (baseline 55.8s, skills 1.7s, speccheck 15.1s) |
| `make repo-test` (ownership + regeneration gates) | ok (baseline 33.7s, skills 14.5s) |
| Marker containment probe | `present: True`, `inside managed region: True` |

### Pre-existing failure, unrelated to this slice

`make docs-test` fails on `TestCheckActiveCorpusHasNoErrors` /
`TestCheckCorpusGolden` with
`0113-a-refused-gate-writes-its-refusal-once: SC-VERIFY-NON-HERMETIC`. The same
two tests fail identically on the pristine tree with this diff stashed, and the
finding names a task file in Spec 0113 that this Task does not touch.

### Follow-up note

`make baseline-digests` also rewrote the `goldenDigest` field in
`internal/baseline/assets/profiles/standard-typescript-monorepo.json`, which the
authorization record's output list does not enumerate even though it is a
mechanical output of the same command (`TestMeasuredSanctionedOwnershipMatchesRecords`
passes with it). Worth adding to that list when the record is next amended.

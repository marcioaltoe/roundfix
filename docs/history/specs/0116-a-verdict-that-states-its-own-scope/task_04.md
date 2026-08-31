---
status: completed
type: docs
---

# Task: The glossary names the Auditing Binary

A QA Report records `build:` and a reader cannot tell whether it identifies the
Roundfix that produced the verdict or the tree it audited. The distinction has
no name, which is why one field carried both readings.

## Work

- Add **Auditing Binary** to the domain context: the Roundfix that produced a
  verdict, as distinct from the tree it audited, carrying version and build
  identity, and reported with a staleness state that may be `unknown`.
- Record what to avoid, in the shape of its neighbours.
- Read the entry from the delivered code, not from the TechSpec draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected.
- Check whether this Spec introduced, changed, or retired any other term the
  glossary should carry, and update the domain context when it did.

## References

- `_prd.md` → Core Feature 3
- `_techspec.md` → Build Order 4; Vocabulary Contract
- `docs/agents/domain.md` owns what the glossary check looks for

## Verification
- `grep -q "Auditing Binary" CONTEXT.md && grep -q "auditing_binary" CONTEXT.md && go test -count=1 -tags docscontract ./internal/docscontract`

## Result

Implementation:

- Added the `Auditing Binary` glossary term, defining the Roundfix binary that
  produced a verdict separately from the audited tree, including its version,
  build commit, build time, and formatted `AuditingBinary` / `auditing_binary`
  identity, plus the `current`, `stale`, or `unknown` `auditor_staleness` state
  with its reason.
- Updated the existing `QA Report` term to describe the auditor metadata now
  carried by reports.
- Corrected the TechSpec's vocabulary note to distinguish the delivered report
  and build-identity code from the authoring-skill template that task_05 owns.
- Term audit: no other new, changed, or retired domain term from this Spec
  requires a glossary entry; `Staleness`, `Auditor`, and `CompareToTree` remain
  implementation API names supporting `Auditing Binary`.

Acceptance evidence:

- `Auditing Binary` is present in `CONTEXT.md` with the exact `AuditingBinary`
  and `auditing_binary` spellings and an `_Avoid_` line; the definition
  distinguishes the producer from the audited tree and records the
  released-build and unknown-staleness cases from the delivered code.
- `QA Report` now names the new metadata fields, matching the serialized fields
  in `internal/spec/qa.go` and the `AuditingBinary` identity in
  `internal/app/version.go`.
- The TechSpec vocabulary note now states that only the report and build
  identity paths are delivered at this point; the QA-gate skill remains for
  task_05.

Focused checks:

- Verification feedback reproduction: `rtk go test -count=1 -tags docscontract
  ./internal/docscontract -run '^TestCheckActiveCorpusHasNoErrors$'` failed
  before the repair with two `SC-VOCABULARY-UNDOCUMENTED` findings for the
  missing `AuditingBinary` token; the specified diagnostic artifact contained
  no other failure class.
- After adding the exact implementation spelling, `rtk go test -count=1 -tags
  docscontract ./internal/docscontract -run
  '^(TestCheckActiveCorpusHasNoErrors|TestCheckCorpusGolden)$'` passed both
  focused corpus checks.
- Red starting signal: `rtk rg -n 'Auditing Binary|auditing_binary'
  CONTEXT.md` returned no matches before the glossary edit (exit 1).
- `rtk rg -n -C 2 'Auditing Binary|auditing_binary|auditor_staleness'
  CONTEXT.md` found the new glossary entry and updated QA Report definition
  after the edit.
- `rtk git diff --check` passed after the documentation edits.
- `rtk make verify-incremental` reached all packages but the sandboxed run
  failed two unrelated force-stop integration tests because process-table
  reads were denied; the permitted host rerun passed the complete incremental
  target, including all tests, skill checks, and the build.
- The Daemon-owned Verification command was not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20260830T161359Z_31aaee7e42ecc4e4`
- Source commit: `e5cef1f9b4d92f4fdb31dbf124399e426165698e`

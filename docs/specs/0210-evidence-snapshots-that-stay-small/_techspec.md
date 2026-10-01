---
spec: 0210-evidence-snapshots-that-stay-small
prd: _prd.md
created: 2026-10-01
---

# Evidence snapshots that stay small — Technical Spec

## Executive Summary

The Evidence Snapshot recorder and the carry resolver of Spec 0202 already do
the right comparison; they persist it in the wrong shape. ADR-0194 records
"the path and SHA-256 digest of every file" each qualifying row declared, so
one `**` input writes one entry per tracked file. This fix keeps the matched
files in memory and persists one count and one digest per declared input. The
carry compares that pair, and the per-file form is converted on read into the
same pair. The trade-off this design accepts is that a stale row names the
input that moved rather than the files under it. The record becomes bounded by
rows times inputs, and the disposition is bounded the same way.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; rows keep their
  identifiers and inputs their declared refs. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0210 (this Spec) refines ADR-0194;
  ADR-0097, ADR-0195 and ADR-0096 hold unchanged; the gate is bound by
  ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167,
  and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  consistency; ADR-0182 and ADR-0184 do not apply, as the PRD records.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty. `internal/speccheck/mechanical_test.go`,
  every skill, `go.mod` and the Makefile stay untouched. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No new package, seam or command. Every change sits in the files Spec 0202
created or extended:

| Component | Where | Change |
| --- | --- | --- |
| Input digest | `internal/speccheck/mechanical.go` | One function derives count and digest from a sorted file list |
| Evidence Snapshot recorder | `RecordEvidenceSnapshots` in `internal/speccheck/evidence_record.go` | Writes one flow-style line per input instead of a file list |
| Snapshot reader | `mechanicalEvidenceSnapshots` in `internal/speccheck/mechanical.go` | Accepts the new entry and the per-file entry |
| Carry proof | `carryRefusal` in `internal/speccheck/mechanical.go` | Compares count and digest per input; the stale reason names refs |
| Guide | `docs/user-guide/context-driven-development.md` | Describes the per-input record and `input moved: <refs>` |

The QA stage in `internal/daemon`, the prior-pass import, `Carriable`'s
signature, `resolveCarriedRows`, the seeded report and the qa-gate skill do
not change.

## Implementation Design

### Interfaces

```go
// internal/speccheck/report.go — EvidenceSnapshot gains the persisted pair.
type EvidenceSnapshot struct {
	Ref    string
	Files  []EvidenceFile // matched files when known: built at a head, or read from the per-file form
	Count  int            // recorded matched-file count; derived from Files when Files is non-empty
	SHA256 string         // recorded input digest; derived from Files when Files is non-empty
}

// internal/speccheck/mechanical.go
// evidenceInputDigest returns the lowercase hex SHA-256 of the h1 summary of
// files, and false when a path contains a newline.
func evidenceInputDigest(files []EvidenceFile) (string, bool)

// evidenceSnapshotPair validates a snapshot against ref and returns the pair
// the carry compares, deriving it from Files when they are present.
func evidenceSnapshotPair(snapshot EvidenceSnapshot, ref string) (count int, digest string, ok bool)
```

### The input digest

For one input at one head, take the matched files in ascending byte order of
their repository-relative paths. For each, write the lowercase hexadecimal
SHA-256 of its blob content, two spaces, the path and `\n`. The input digest
is the lowercase hexadecimal SHA-256 of that summary, and the count is the
number of files. This is the summary of Go's `h1:` hash (PRD Acceptance
evidence), hex-encoded rather than base64. A path containing `\n` cannot be
summarized: `buildEvidenceSnapshots` reports the input unresolved, so the
recorder skips the row and the resolver refuses it with `evidence differs`.

`evidenceSnapshotPair` accepts exactly one of two shapes:

1. `Files` non-empty: the existing `validEvidenceSnapshot` rules hold (ref
   matches, paths clean, strictly ascending, matched by the ref, digests 64
   lowercase hex). The pair is derived. A recorded `Count` or `SHA256` beside
   files makes the snapshot invalid.
2. `Files` empty: `Count` ≥ 1 and `SHA256` is 64 lowercase hex characters.

Any other snapshot is invalid, and the carry reason is `no evidence snapshot`.

### Data Models

No database change. Each input entry of `evidence_snapshots` becomes one
flow-style mapping line. `RecordEvidenceSnapshots` sets `yaml.FlowStyle` on
the entry's mapping node; `yaml.v3` keeps a flow mapping on one line at any
width, as an authoring probe showed for a 180-character entry.

```yaml
evidence_snapshots:
  "3":
    head: <audited head>
    inputs:
      - {ref: 'internal/**', count: 2114, sha256: <64 lowercase hex>}
      - {ref: go.mod, count: 1, sha256: <64 lowercase hex>}
```

The block has one line for the key, three per recorded row and one per
declared input. The reader's input struct gains `Count int` (`count`) and
`SHA256 string` (`sha256`) beside the existing `Files` list, so a per-file
entry still decodes.

### Recording the snapshot

`RecordEvidenceSnapshots` keeps its qualification rules and its coverage
check, which reads the in-memory `Files` that `buildEvidenceSnapshots`
returns. Only the rendering changes: each input writes `ref`, `count` and
`sha256` from `evidenceInputDigest`. Every byte outside the key keeps its
current treatment.

### Carrying by digest

`carryRefusal` keeps its order of checks:

1. Status, establishing head and input kinds, unchanged.
2. Changed paths: each declared input whose matcher matches a changed path is
   moved. When any is, the reason is `input moved: ` followed by the moved
   refs in declaration order, joined by `, `. Today this step names the
   changed paths instead.
3. Snapshot count and validity through `evidenceSnapshotPair`, for the
   established and the current snapshot; a failure is
   `no evidence snapshot`.
4. Digest comparison: each input whose established and current pair differ is
   moved, with the same `input moved: <refs>` reason. Today a differing file
   list is `evidence differs`.
5. Cited evidence paths: a path is covered when an input's matcher matches it
   and the current snapshot's `Files` contain it. Equal pairs prove the
   established set held it too. An uncovered path is `evidence differs`.

`resolveCarriedRow` is unchanged: it passes the recorded snapshots as
established and `buildEvidenceSnapshots`' result, which always has `Files`,
as current.

### API Contracts

1. API Contract: `evidence_snapshots` — Daemon-written frontmatter keyed by
   row identifier, each row holding `head` and an `inputs` list whose entries
   are one-line mappings of `ref`, `count` and `sha256` as Data Models shows.
   A per-file entry (`ref` and `files`) is still read and converted.
2. API Contract: Carry Disposition reason — `re-run: input moved: <refs>`
   names the declared inputs whose matched files changed, in declaration
   order, joined by `, `. Every other reason of ADR-0195's closed list keeps
   its wording.

### Surface Transcripts

None. No command, flag, stream or exit code changes.

## Vocabulary Contract

No new term or emitted word. The constant `CarryReasonInputMoved`
(`input moved: `) keeps its bytes, and only what follows it changes from paths
to refs. The guide already documents it, and task_02 updates its placeholder.

## Coverage Map

- Goal 1 → The input digest; Data Models; Recording the snapshot.
- Goal 2 → Carrying by digest.
- Goal 3 → The input digest (`evidenceSnapshotPair`); Data Models.
- Core Feature 1 → Recording the snapshot; API Contract 1.
- Core Feature 2 → Carrying by digest; API Contract 2.
- Core Feature 3 → The input digest; Carrying by digest.
- Core Feature 4 → Integration Points; Testing Approach 2.
- Core Feature 5 → Build Order 2.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 1.
- Success Metric 4 → Testing Approach 2.

## Integration Points

- **Prior pass import (ADR-0194).** `importPriorQAPass` copies a failed
  pass's report byte for byte and keeps it when `ReportShapeFindings` finds
  nothing. Neither reads the snapshot block, so an imported report in either
  form is kept, and the next pass's resolver reads it like any other.
- **Spec 0203's report.** Its block was removed by hand before delivery, so
  no report in this repository holds a block larger than the Spec 0202 test
  fixtures.
- **Reports recorded by a pre-fix Daemon.** Any per-file block written before
  this fix is read through the conversion, and still carries.

## Testing Approach

1. **`internal/speccheck`**, over temporary Git repositories. A new file
   `internal/speccheck/evidence_digest_test.go`:
   - a row whose glob matches 2,500 committed files records one line for that
     input, and the block's line count equals the formula of Data Models;
   - the recorded digest equals an independent `h1` summary built in the
     test, for files committed out of order;
   - a new-form record carries through `RunMechanicalStage` at an unchanged
     head, and the same report re-runs the row with exactly
     `input moved: src/**` after two files under `src/` change, one is added
     and one is removed;
   - a recorded count or digest that differs from the head re-runs the row
     with `input moved: <ref>`;
   - a per-file record of a two-file glob carries; a per-file record missing
     one file re-runs with `input moved: <ref>`; an unsorted list and an
     entry holding both `files` and `sha256` re-run with
     `no evidence snapshot`; in every case the stage returns no error and no
     finding;
   - `Carriable` accepts an established snapshot holding only the pair
     against a current snapshot with files, and refuses a differing pair.

   `TestEvidenceRecordSnapshotsEveryQualifyingPassRow` in
   `internal/speccheck/evidence_record_test.go` changes to assert count and
   digest per input instead of file lists. The governed `TestCarriable` and
   `TestMechanicalStageCarriable…` tests, every `TestRowCarry…` test and the
   other `TestEvidenceRecord…` tests run unedited. Several of them write
   per-file blocks, so they also pin the conversion.
2. **`internal/daemon`**, over the task-cycle fixture.
   `TestQAGateCommitsTheEvidenceSnapshotAtTheAuditedHead` asserts the
   committed one-line entry instead of `path: snapshot-input.txt`.
   `assertTwoPassSnapshot` in `qa_two_pass_carry_test.go` decodes `count` and
   `sha256` and compares them with an independent summary. The two two-pass
   tests import the first pass's new-form report and carry from it, which
   proves the import path. `TestQAGateCarriesARowFromAnUnintegratedFailedPass`
   imports a per-file report and carries from it, unedited.
3. **Real history.** The QA gate recomputes the 0203 measurement from
   `a1fc8402^` and checks the digest of one recorded input with
   `git ls-tree` and `sha256sum`, outside Go.

## Build Order

1. The input digest, the recorder's one-line entries, the reader's two
   shapes, the carry by digest and the test updates of Testing Approach 1 and
   2, task_01 (depends on: none).
2. The guide's Evidence Snapshot paragraph and the `input moved: <refs>`
   placeholder, task_02 (depends on: 1).
3. Terminal QA, task_03 (depends on: 1, 2).

## Risks & Considerations

- **This Spec's own gate.** The Daemon that runs this Spec's QA gate writes
  the block, so a Daemon built before task_01 writes the per-file form. The
  operator delivers this Spec with a binary rebuilt from a head that holds
  task_01, or the QA Task declares narrow inputs.
- **Less detail in a disposition.** `input moved: internal/**` does not say
  which file moved. The Run's changed-file evidence and `git diff` still do.
- **Declared break.** `input moved:` now names refs. Every existing assertion
  of it uses a literal ref equal to its only path, so no test changes for it.
- **Mixed shapes.** An entry holding both shapes is refused rather than
  trusted, because only the Daemon writes the block and it writes one shape.

## Decisions

- Persist one count and digest per input, over Go's `h1:` summary. See
  ADR-0210.
- Name refs in `input moved:`. See ADR-0210.
- Convert the per-file form on read instead of refusing it, so reports
  recorded since Spec 0202 keep carrying.
- Keep `Carriable`'s signature and the governed tests unedited, by keeping
  `Files` in memory and adding the pair beside it.

---
spec: 0210-evidence-snapshots-that-stay-small
status: archived
created: 2026-10-01
surfaces: [backend, docs]
archived: "2026-10-01"
source_slug: 0210-evidence-snapshots-that-stay-small
---


# Evidence snapshots that stay small

Spec 0202 (merged, not yet released) has the Daemon record an Evidence
Snapshot when a QA pass closes: for each carriable passing row, every file a
declared input matches, with its SHA-256 digest. A glob input matches every
file under it. On 2026-10-01, Spec 0203's QA Report declared
`internal/**`, `cmd/roundfix/**`, `go.mod`, `go.sum`, `Makefile`, `**` and the
Spec directory on each of 11 rows. The report grew to 157,212 lines, 156,696
of them in the `evidence_snapshots` block, and the pre-PR review's agent
session then failed with `agent/protocol error` on the diff. The operator
removed the block by hand to deliver Spec 0203 (commit `a1fc8402`). Every
later QA Report whose rows declare a broad glob would grow the same way, and
v0.24.0 would ship that behavior.

This is a bug fix with no product-behavior change: rows carry under the same
conditions, and only the recorded shape and the stale reason's wording change.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; rows keep their
  identifiers and inputs keep their declared refs. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0210 (this Spec) records one
  digest per declared input and refines ADR-0194, whose recording, import and
  head proof otherwise hold. ADR-0097's carry conditions and ADR-0195's
  always-observed rows and Carry Dispositions hold unchanged. ADR-0096's
  mechanical stage is where the carry is proved, and it keeps its role. The
  gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155,
  ADR-0156 and ADR-0167, and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and
  ADR-0183 check this Spec's consistency by citation and receipt. ADR-0182
  does not apply, because no Task settlement fact changes, and ADR-0184
  does not apply, because no command surface changes.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty. `internal/speccheck/mechanical_test.go`,
  every skill, `go.mod` and the Makefile stay untouched. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A QA Report's `evidence_snapshots` block grows with rows and declared
  inputs, never with the number of files an input matches.
- A row still carries only when every file under each declared input is
  unchanged.
- A report recorded in Spec 0202's per-file form is still read, and never
  fails the gate.

## Core Features

1. **One digest per input.** The Daemon records, for each declared input of a
   carriable row, its ref, its matched-file count and one SHA-256 digest over
   the sorted matched paths and their content digests, on one line (ADR-0210).
2. **Carry by digest.** The mechanical stage compares the recorded and current
   count and digest of each input. Any added, removed or changed matched file
   makes the row stale with `re-run: input moved: <refs>`, naming the declared
   inputs that moved.
3. **The per-file form is still read.** A per-file list is converted on read
   into the same count and digest. A list the stage cannot validate makes the
   row `re-run: no evidence snapshot`, and the stage raises no finding.
4. **The failed-pass import keeps working.** A failed pass imported from its
   Run Branch carries its rows in either form (ADR-0194).
5. **The guide says so.** The Context-Driven Development guide describes the
   per-input record and the reason's new wording.

## Non-Goals / Out of Scope

- Changing which rows qualify, the head proof, the import rules, Task
  Carry-Forward, verdict rules or report naming.
- Rewriting an existing QA Report, archived or active.
- Narrowing the inputs gate Agents declare; conservative declaration stays
  the skill's rule.
- Editing the qa-gate skill, the QA prompt contract or `CONTEXT.md`: none of
  them describes the per-file form.

## Success Metrics

1. A row whose input matches 2,500 files records exactly one line for that
   input, and the snapshot block has one line per declared input plus three
   per row and one for the key.
2. A row recorded in the new form carries at an unchanged head, and a file
   added, removed or changed under its input re-runs it with
   `input moved: <ref>`.
3. The governed carry tests and every Spec 0202 test that writes a per-file
   block pass unedited.
4. A failed pass recorded in either form is imported, and its unmoved row is
   carried by the next pass.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- Spec 0203's QA Report before commit `a1fc8402` removed its block, read with
  `git show a1fc8402^:docs/specs/0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds/qa/qa-report-2026-10-01.md`.
  On 2026-10-01 it held 157,212 lines and 78,254 `sha256:` lines over 11
  rows with 7 inputs each. The new form records the same rows in 1 + 11 × 3 +
  77 = 111 lines.
- Go's `h1:` directory hash, `golang.org/x/mod/sumdb/dirhash.Hash1`
  (<https://pkg.go.dev/golang.org/x/mod/sumdb/dirhash#Hash1>), hashes "a
  summary" with "a single line for each file in the list, ordered by
  [slices.Sort] applied to the file names", each line the hexadecimal SHA-256
  of the content, two spaces, the name and a newline. ADR-0210 adopts that
  summary, hex-encoded.

## Decisions

- One digest per input, over Go's `h1:` summary. See ADR-0210.
- The stale reason names the declared input, not the changed files, so a
  disposition is bounded too. See ADR-0210.
- The per-file form is converted on read rather than refused, so reports
  recorded since Spec 0202 keep carrying.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.

---
spec: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
status: active
created: 2026-10-06
surfaces: [backend, cli, docs]
---

# An archive that keeps the report, not the raw evidence

Every archived Spec carries its raw QA evidence into the History Root for
good: command transcripts, JSON dumps, harness scripts, fixtures,
screenshots and SQLite files. On 2026-10-06 that was 21.0 MB of the History
Root's 50.9 MB, in 2,713 of its 5,448 files, and 1,440 of the 2,217 files the
History Root gained in the previous ten days. Nobody reads these files once
the Spec is archived. The QA Report already records the verdict and every
row. The maintainer chose a retention rule inside the repository: the
Archive Command drops the raw evidence and keeps the report, with a manifest
that preserves provenance. Already-archived Specs can be cut on request.

## Prerequisites

None of this Spec's Verification reads another active Spec's artifacts.
Specs 0235, 0236 and 0237 are authored in parallel and are delivered first,
in that order; this Spec is delivered after 0237. Spec 0235 edits the same
Baseline module and its derived files, and Specs 0236 and 0237 raise the
Roundfix skill's version. This Spec's task_04 regenerates those files and
re-records the skill versions on top of theirs.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes. The
  new names are a command flag, a manifest file name and Go functions.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network call or
  HTTP surface is added; the cut reads and writes files in the working tree
  and reads Git objects, and no test reaches the network. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0243 (this Spec) decides the
  cut: "The Archive Command now cuts the raw evidence as part of a normal
  archive". It applies ADR-0215, which decided that "A removal is its own
  later change, and the maintainer's explicit approval of that removal is
  recorded in it", and extends ADR-0230, under which "The Delivery Queue
  accepts an archive commit whose only content changes besides the archive
  stamp are these rewrites". ADR-0154 stands, and an override archive keeps
  its QA files: ADR-0154: "The exception can waive the terminal QA Task's
  archive completion/evidence requirement without rewriting its status,
  Result or report". ADR-0223 refuses an archive while a file other than
  Markdown "names an active Spec's directory", and the cut of an archived
  Spec gains the same refusal for its evidence directory. ADR-0120 places
  the History Root, ADR-0120: "retired documentation is documentation".
  ADR-0210's Evidence Snapshots hash repository inputs, not evidence files,
  ADR-0210: "The record now holds one entry per declared input". ADR-0232's
  merge evidence supersedes a merged Spec's other commits, ADR-0232:
  "without a content comparison", so it never reads a dropped file.
  ADR-0229's operator archive and ADR-0237's Delivery Retry run the same
  Archive Command and exact-move check, so they inherit the cut unchanged,
  and ADR-0165's park of a blocking review after archive reads the review,
  not the evidence. ADR-0182 binds ADR-0169's pre-PR review diffs the candidate, from which QA evidence is already omitted, and ADR-0240's partial policy is applied before the cut and reads only the report and declarations, so neither changes. ADR-0104, ADR-0167 and ADR-0196 bind the QA gate's rows, its outside evidence and the pre-PR review, which this Spec uses unchanged. each Task's Verification to the facts its
  gate checks. ADR-0179 bounds the Governed Paths, ADR-0187, ADR-0189 and ADR-0233 bind
  the owned skills' versions, and ADR-0184 binds the Surface Transcripts.
  ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this
  Spec's consistency by citation and receipt. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Archive Command's source and its
  tests, the coverage record and its test, the repository-copy helper, the
  archive-spec, qa-gate and Roundfix skills, the Spec workflow Baseline
  module with the derived Baseline files it regenerates, and two generated
  guides are Governed Paths. The maintainer approved this Spec on 2026-10-06
  ("Sim, como Spec"), beside the standing grants of 2026-09-30 for the skills
  ("considere autorizado a ajustar todas as skills se necessário") and the
  Baseline source and guides ("Autorizar os dois"). No other Governed Path
  changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0238-an-archive-that-keeps-the-report-not-the-raw-evidence/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`,
  `docs/agents/spec-routing.md`, `docs/references/coverage-record.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/spec/archive.go`, `internal/spec/coverage_test.go`,
  `skills/archive-spec/SKILL.md`, `skills/baseline_skill_contract_test.go`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## Goals

- A normal archive drops the Spec's `qa/evidence/` directory and leaves an
  Evidence Manifest naming each dropped file's path, size and SHA-256, the
  revision that held it and where it lived.
- Every relative link in the Spec that reached into the dropped evidence
  reaches the manifest instead, so no archived link is left dangling.
- The Delivery Queue accepts an archive commit with the cut as an exact
  Spec move, so delivery proceeds as before.
- An operator can apply the same cut to one already-archived Spec, after a
  dry run that shows what would go, and nothing runs it implicitly.
- `make verify` and `make verify-docs` pass on a tree without archived
  evidence, because no test reads it.

## User Stories

1. As the maintainer, I want an archive to leave the QA Report and drop the
   raw evidence, so that the History Root and every mirror of it stop
   growing by files nobody reads.
2. As a reader of an archived Spec, I want the report's evidence links to
   reach a manifest with each file's digest and revision, so that I can
   still recover and verify a dropped file from Git history.
3. As the operator of the Delivery Queue, I want an archive commit with the
   cut accepted as an exact move, so that delivery does not park on it.
4. As the maintainer, I want to cut an already-archived Spec on request,
   with a dry run first, so that old Specs shrink only when I choose.
5. As a contributor, I want the repository's gates to pass without archived
   evidence, so that a cut never breaks `make verify` or `make verify-docs`.

## Core Features

1. **The cut at archive.** A normal archive with a Task Graph removes every
   file under `qa/evidence/` and writes `qa/evidence-manifest.md` when at
   least one file was removed. A Spec without evidence archives exactly as
   today. ADR-0243: "Everything under the Spec's `qa/evidence/` directory,
   of any type. Nothing else in the Spec is cut".
2. **The Evidence Manifest.** Front matter names the schema, the date, the
   revision, the source Spec directory, the file count and the byte total,
   and a table lists each file's path, bytes and SHA-256 in path order.
3. **Links reach the manifest.** A relative Markdown link inside the Spec
   whose target is in `qa/evidence/` is rewritten to reach the manifest; its
   text and title are kept. The archive's confirmation reports how many
   files and bytes were dropped and how many links now reach the manifest.
4. **Exceptions.** A QA Archive Override archive and a superseded Spec keep
   their files. A path the manifest cannot represent, or an entry that is
   neither a directory nor a regular file, refuses the archive before any
   file changes.
5. **An exact move for the queue.** The Delivery Queue's exact-move check
   accepts the dropped files when the archived manifest lists exactly them
   with matching sizes and digests, and the evidence links that now reach
   it.
6. **Cutting an archived Spec.** `roundfix archive <slug> --drop-evidence`
   reports what the cut would drop, rewrite and write; `--apply` performs
   it with the current revision. It refuses an active Spec and a Spec whose
   evidence directory a file other than Markdown names, keeps an override or
   superseded Spec unchanged, and reports nothing to drop for a Spec already
   cut.
7. **No test reads archived evidence.** The coverage record no longer
   counts Go packages under `docs/`, and the repository-copy helper skips a
   tracked file deleted from the working tree.
8. **The guides say it.** The archive-spec, qa-gate and Roundfix skills, the
   archive user guide and the Spec workflow Baseline clauses describe the
   cut, the manifest and the request-only cut of archived Specs.

## Non-Goals / Out of Scope

- Running the cut on the existing History Root. That is a separate change
  the maintainer may make later with the new command, and this Spec deletes
  nothing under `docs/history/`.
- Rewriting Git history or shrinking `.git`, here or in the Secondbrain
  (declined on 2026-10-06, "Não agora").
- The Secondbrain mirror exclusion, already done on 2026-10-06.
- Cutting anything outside `qa/evidence/`, or cutting by age.
- Changing the QA gate's verdict rule, the QA settlement table, Evidence
  Snapshots, the review scope or the merged-Run reconciliation.
- `CONTEXT.md`, `CHANGELOG.md`, the Makefile, `go.mod` and CI.

## Success Metrics

1. Success Metric: archiving a fixture Spec with 3 evidence files, one
   linked from its QA Report, leaves no `qa/evidence/` directory, a manifest
   listing the 3 files with their sizes and SHA-256 digests, and the report's
   link resolving to the manifest; a fixture Spec without evidence archives
   byte-for-byte as before.
2. Success Metric: the Delivery Queue treats that archive commit as an exact
   Spec move, and still refuses one whose manifest omits a dropped file or
   records a wrong digest.
3. Success Metric: `roundfix archive <slug> --drop-evidence` changes no file,
   `--apply` cuts exactly what the dry run reported, and a second `--apply`
   reports nothing to drop.
4. Success Metric: in a disposable clone where every archived `qa/evidence/`
   directory is removed and the removal committed, `make verify` and
   `make verify-docs` exit 0.
5. Success Metric: a dry run of `--drop-evidence` over every Spec of this
   repository's History Root reports about 13.5 MB in 1,698 files to drop
   and keeps the 37 override Specs (7.5 MB) unchanged.

## Acceptance evidence

The outside-evidence rows rest on records this Spec did not produce:

- Spec 0214's read-only measurement, `docs/references/archived-evidence-measurement.md`.
  Its ablation removed archived evidence in a disposable clone and found
  `make verify` failing in `TestCoverage` for three Go packages inside
  archived evidence, and `make verify-docs` failing in
  `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical` on a missing
  tracked file. It also records that history keeps every removed version.
- The operator's Jev judgment of 2026-10-06
  (`~/.roundfix-operator/history-value-2026-10-06.md`, model
  `typesafe/jev-1.13-20260917`, 123 sampled artifacts, 9 of 10 matching a
  hand label). Raw evidence files were judged transient (about 31% of the
  text bytes), and all 14 sampled QA Reports, every sampled PRD, TechSpec,
  Task Graph and authorization, and 9 of 10 Task files were judged repository
  records.
- Adopter repositories mirrored in the Secondbrain, measured on 2026-10-06
  with `du`: raw evidence is 20.2 of 32.6 MB of Fluxus's History Root, 37.4
  of 49.9 MB of Vortex's, 21.2 of 31.2 MB of Oraculum's, 61.5 of 65.2 MB of
  Conexus's and 2.0 of 10.4 MB of Fiscus's.
- GitHub's published retention for workflow evidence, read 2026-10-06
  (<https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository>):
  "By default, the artifacts and log files generated by workflows are
  retained for 90 days before they are automatically deleted." Raw run
  evidence is routinely treated as expiring while the run's verdict stays.
- The Git book's maintenance chapter
  (<https://git-scm.com/book/en/v2/Git-Internals-Maintenance-and-Data-Recovery>),
  which Spec 0214 cited: removing a file from the tree leaves its versions
  in history, so the cut reclaims the tree and not the pack.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "QA evidence retention archive drop raw evidence keep report"`
(`--all --files --min-score 0.3`). It returned QA Reports of Vortex, Fiscus
and Fluxus, one of this repository's QA decision records and its archive
user guide, and no
prior retention rule for QA evidence in any project. The pending inbox for
Roundfix holds only triaged entries. Exa found GitHub's retention page. No
open Backlog Entry or unresolved Finding shares this context: the two open
entries of 2026-10-06 belong to Specs 0236 and 0237.

## Decisions

- The cut is part of a normal archive and covers `qa/evidence/` only. See
  ADR-0243.
- The manifest is Markdown in the Spec's `qa/` directory, so a reader and
  the link rewrite reach it like any other Spec file.
- Override and superseded archives keep their files, because the QA
  settlement contract keeps an override's QA files as observed.
- Already-archived Specs are cut one at a time through the Archive Command,
  dry run by default, rather than through a new history command.

## Open Questions

- Whether to cut the existing History Root now. Default: no; the
  maintainer runs the command in a later change.
- Whether override archives should also cut, which would change the QA
  settlement table shared by three skills. Default: they keep their files.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.

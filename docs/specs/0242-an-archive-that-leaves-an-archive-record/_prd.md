---
spec: 0242-an-archive-that-leaves-an-archive-record
status: active
created: 2026-10-06
surfaces: [backend, cli, docs]
---

# An archive that leaves an Archive Record

`docs/history` held about 51 MB in 5,400 files on 2026-10-06 and had grown
about 22 MB in the ten days before. A Jev judgment of 123 sampled artifacts
rated about 8% of that text as reusable knowledge, 60% as records useful only
inside the repository and 31% as throwaway evidence
([the adopted Backlog Entry](references/2026-10-06-history-keeps-only-what-the-secondbrain-needs.md)).
The maintainer decided the same day: "Não quero nada no histórico que não
seja relevante para o secondbrain".

This Spec is the first of two. It changes what an archive leaves behind.
Instead of moving the Spec folder into the History Root, `roundfix archive`
writes a small **Archive Record** and removes the folder in the same change.
Before the cut, a person or Agent moves reusable knowledge upstream, with
Jev's advice. Every reader that today opens an archived Spec folder is fixed
first, so it works from the record or from Git. The second Spec, 0243,
migrates the folders already under `docs/history` in batches. This Spec
deletes nothing that is already archived.

## Prerequisites

Specs 0239, 0240 and 0241 are authored in the same cycle and are delivered
before this one, which comes fourth, in v0.53.0. Each may raise an owned
skill's version, change `CONTEXT.md` or regenerate Baseline-derived files.
This Spec re-records skill versions and regenerates derived files on top of
theirs (ADR-0233). The queue does not enforce that order, so the operator
orders it. No Verification here depends on another active Spec's artifacts.

## Project Constraints

- Identifier strategy: applicable — the Archive Record is named
  `<slug>.md` beside the archive root's folders and carries the schema string
  `roundfix/archive-record/v1`. The Jev advice adds the judgment kind
  `archive-value`. Spec, Task and ADR identifiers keep their schemes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — `roundfix archive <slug> --plan`
  may send file text to the Spec judge's existing transport, under its
  existing stage key, monthly ceiling and judge log. It sends nothing when no
  key is set, and the advice fails open. Every test, Verification command and
  QA row uses a fake transport and a temporary home. No row reaches a
  provider. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0247 (this Spec) decides that an
  archive leaves an Archive Record and removes the Spec folder, keeps the
  folder's bytes in Git at the recorded revision, takes Jev's advice without
  letting it gate, and pins the pre-cut corpus at
  `40a7893d872c8a6705f6d7745e6efe430ae9deeb`. It narrows ADR-0215 for archives
  from now on, whose rule is ADR-0215: "A removal is its own later change, and
  the maintainer's explicit approval of that removal is recorded in it", by
  recording that approval. ADR-0120 keeps the History Root, ADR-0120: "retired
  documentation is documentation", and the records live under it. ADR-0154
  binds an override, ADR-0154: "The exception can waive the terminal QA Task's
  archive completion/evidence requirement without rewriting its status, Result
  or report", so the record keeps every override field and the reason.
  ADR-0230's link pass keeps serving folders archived earlier. ADR-0121
  records the Baseline history layout's relocations of those folders,
  ADR-0121: "Archive relocations are therefore recorded as their own ordered
  ledger of source, destination, and content identity". ADR-0223 refuses an
  archive while a file names the active Spec's directory, and that refusal
  stays. ADR-0229 and ADR-0237 run the same Archive Command for an operator
  archive and a Delivery Retry. ADR-0232 grants merge evidence from the
  archived `_prd.md`, and this Spec moves that evidence to the record.
  ADR-0165 parks a blocking review after archive, and the review must still
  see an archived Spec; ADR-0169's merge-base diff is the diff from which the
  review now omits the removed folder. ADR-0193 counts an archived
  prerequisite. ADR-0210's Evidence Snapshots hash declared inputs and read no
  archived folder. ADR-0179 bounds the Governed Paths, and ADR-0187, ADR-0189
  and ADR-0233 bind the owned skills' versions. ADR-0184 binds the Surface
  Transcripts. ADR-0182 binds each Task's Verification to the facts its gate
  checks. The gate is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097,
  ADR-0104 and ADR-0167, ADR-0196 validates a pre-PR finding before it parks,
  and ADR-0240 decides when its QA partial qualifies. ADR-0194 and ADR-0195
  cite ADR-0097 but decide what a QA row records and when it is observed
  again; this Spec changes neither, so neither applies. ADR-0093, ADR-0117,
  ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this Spec by citation and
  receipt. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Tasks change Governed Paths: the
  owned skills, the Baseline modules and derived files, the coverage record
  and its test, the repository-copy helper, the governed-path contract test,
  the suite guard's regeneration reader and the Spec and authorization
  readers. The maintainer authorized skill edits on 2026-09-30 ("considere
  autorizado a ajustar todas as skills se necessário") and Baseline sources
  and guides ("Autorizar os dois"). On 2026-10-06 the maintainer granted the
  Governed Paths this Spec declares ("Concedo"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0242-an-archive-that-leaves-an-archive-record/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `docs/references/coverage-record.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/derived_ownership_test.go`, `internal/baseline/derived_regeneration_repocontract_test.go`, `internal/spec/archive.go`, `internal/spec/archive_layout_characterization_test.go`, `internal/spec/archive_test.go`, `internal/spec/coverage_test.go`, `internal/speccheck/backlog.go`, `internal/speccheck/governed_repocontract_test.go`, `internal/suiteguardcontract/regeneration.go`, `skills/archive-spec/SKILL.md`, `skills/owned_skill_edit_repocontract_test.go`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## Goals

1. An archive leaves one Archive Record of about 2 KB per Spec and no Spec
   folder. This holds for a normal archive, a QA Archive Override and a
   superseded Spec without a Task Graph.
2. The Spec folder's bytes stay recoverable from Git at the revision the
   record names.
3. Reusable knowledge moves upstream before the cut. Jev advises, a person or
   Agent confirms, and Jev never gates.
4. Every reader of archived Specs works from the record or from Git. Both
   repository gates pass on a tree with no Spec folder under
   `docs/history/specs`.
5. The skills, guides, Baseline clauses and glossary describe the record and
   the cut.

## User Stories

1. As the maintainer, I want an archive to leave only a small record, so that
   the history keeps only what the Secondbrain and the repository's own
   readers need.
2. As an operator archiving a Spec, I want to see what the cut removes and
   which files Jev thinks are reusable, so that I can promote them before
   they leave the tree.
3. As the Delivery Queue, I want to prove that an archive commit removed
   exactly the Spec and added a record that agrees with it, so that I can
   publish it without a person.
4. As a maintainer of this repository, I want every test and check to pass
   without archived Spec folders, so that Spec 0243 can migrate the existing
   history without breaking a gate.

## Core Features

1. **The Archive Record.** The archive writes `<archive-root>/<slug>.md`
   under the schema `roundfix/archive-record/v1`. It is at most 2,048 bytes
   when only the outcome paragraph needs shortening, and the other fields are
   never cut. The record carries:
   - the slug, title, status, creation and archive dates and disposition;
   - `source` and `source_revision`;
   - the QA Task, report and verdict, and the `unproven` actions;
   - every override field, the reason included;
   - `superseded_by`;
   - the ADRs, adopted sources, sanctioned regeneration commands and promoted
     files;
   - one outcome paragraph.
2. **The cut.** In one change, the archive writes the record, copies any
   promoted files to `docs/references/` and removes the Spec folder. It
   refuses when the folder differs from `HEAD`.
3. **Promotion and advice.** `--promote <path>` copies one Spec file
   upstream. `--plan` lists what the cut removes, with Jev's advice per
   candidate file when a judge key is set. It writes nothing in the
   repository.
4. **Runtime readers.** These read the record, or the pre-archive tree in the
   Run's or candidate's own history:
   - the Delivery Queue's item inspection, archive stage, exact-archive
     proof, prerequisites and review correction;
   - reconcile's merge evidence, Task completion, leftovers and QA
     supersession;
   - the pre-PR review's archived-Spec context, override convention and
     scope;
   - supersede, Run causes, the Spec audit, and the Spec check's backlog,
     absorption and Task Context checks;
   - the suite guard's sanctioned regenerations.
5. **Repository readers.** These read the pre-cut corpus from Git at the
   pinned commit, or carry their bytes as test data: the 22 tests the
   ablation found failing, the coverage record and seven user-guide links.
6. **Guidance.** This covers the archive-spec, qa-gate and Roundfix skills,
   their identical `### QA settlement` tables,
   the archive and context-driven user guides, two Baseline clauses, and the
   glossary terms **Archive Record** and **Archive Advice**, with the revised
   **Archive Command** entry.

## Non-Goals / Out of Scope

- Migrating, deleting or rewriting any file already under `docs/history`,
  and the `history-full` tag. Spec 0243 owns them.
- Rewriting Git history or shrinking the pack ("Não agora", 2026-10-06).
- Writing into the Secondbrain from the archive. The Secondbrain guide owns
  Inbox Entries.
- Changing `.secondbrain-export`. That waits until the history holds only
  records.
- Retiring ADR-0230's link pass or ADR-0121's ledger. Both still serve
  folders archived earlier.

## Success Metrics

1. Success Metric: archiving a temporary Spec leaves a record of at most
   2,048 bytes and no Spec folder. This holds for a passing Spec, an
   override with a 300-byte reason and a superseded Spec. The removed bytes
   read back from `git show <source_revision>:<source>/...` byte-identical.
2. Success Metric: the Delivery Queue accepts an archive commit that removes
   exactly the folder and adds an agreeing record. It refuses one that keeps
   any file, adds another path or names another revision. Merge evidence and
   Task completion hold after a squash merge.
3. Success Metric: `--plan` against a fake judge lists every file the cut
   removes and the advice per candidate. Without a key it names the skip
   reason. Both runs change no file.
4. Success Metric: on a copy of this repository with every Spec folder under
   `docs/history/specs` removed, `make verify` and `make verify-docs` exit 0.
   Before this Spec, 17 tests failed under `make verify` and 5 under
   `make verify-docs` (`_techspec.md` → Readers).
5. Success Metric: every reader of the measured inventory is either changed
   with a test or recorded as unaffected with the reason.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Jev judgment of 2026-10-06
  (`~/.roundfix-operator/history-value-2026-10-06.md`). It sampled 123
  artifacts under `typesafe/jev-1.13-20260917`, with 9 of 10 hand labels
  matching. It found QA reports, core artifacts and Task files to be records,
  and evidence to be transient. The reusable files sat in `references/` and
  `measurement/`, where mean confidence was 0.57. That shaped the advice's
  candidate set and kept it advisory.
- The Spec 0214 measurement in `docs/references/archived-evidence-measurement.md`.
  Its static reader scan and its ablation showed that removing archived
  evidence fails both gates.
- An ablation during authoring on 2026-10-06. A `git clone --no-local` of
  `40a7893d` in a scratch directory had every Spec folder under
  `docs/history/specs` removed in its own commit. `make -k verify` exited 2
  with 17 failing tests in `internal/authorization`, `internal/judge`,
  `internal/spec` and `internal/speccheck`. `make -k verify-docs` exited 2
  with `TestUserGuideLinksResolve` (seven links),
  `TestMeasuredSanctionedOwnershipMatchesRecords`,
  `TestDeclaredStepRegenerationAndFrozenBoundaries`,
  `TestEveryBoundedPathIsGoverned` and
  `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`. `spec check`
  reported nothing.
- GitHub's documentation on checking out pull requests locally
  (<https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/checking-out-pull-requests-locally>).
  A pull request's commits stay fetchable as `pull/ID/head` after its branch
  is gone, so a squash-merged Spec's `source_revision` remains reachable.
- The Git book on maintenance and data recovery
  (<https://git-scm.com/book/en/v2/Git-Internals-Maintenance-and-Data-Recovery>):
  removing a file from the tree keeps every reachable version.

## Decisions

- An archive leaves an Archive Record and removes the Spec folder, and the
  folder stays in Git at `source_revision`; see ADR-0247.
- Jev advises on candidate files through the Spec judge's key, ceiling and
  log and never gates. A person or Agent promotes; see ADR-0247.
- Repository tests read the pre-cut corpus at a pinned commit; see ADR-0247.

## Open Questions

None. Spec 0243's scope is listed in `_techspec.md` → Risks & Considerations.

## Technical candidate

The [_techspec.md](_techspec.md) records the reader inventory, the record
format, the implementation map and the build order.

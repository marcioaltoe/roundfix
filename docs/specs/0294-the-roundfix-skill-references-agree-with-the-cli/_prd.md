---
spec: 0294-the-roundfix-skill-references-agree-with-the-cli
status: active
created: 2026-10-21
surfaces: [docs]
---

# The Roundfix skill references agree with the CLI

Minimal PRD for a documentation bug fix, minted by `write-techspec`. Two
references of the Roundfix skill, which every adopter receives through the
Baseline, contradict the CLI they describe. The `review` reference still
gives the pre-PR reviewer's C3 and C5 Delivery Conventions in their old form,
a Spec folder moved to the archive root and override fields in an archived
`_prd.md`, while the reviewer's prompt says an Archive Record is written and
the Spec folder removed, with the override fields in the record. The `events`
reference says the default replay and `--filter` accept only five categories,
while the CLI accepts and replays `usage` too, and the same file's Token usage
section says so. After this Spec both references, the user guides and the
glossary state what the code does, and documentation contract tests read the
conventions and the categories from the code, so the next change that moves
either fails a test until the documents follow. The TechSpec holds the design.

## Sources

One Secondbrain Inbox Entry, copied verbatim by this Spec's planning change
into `references/` (see `references/_index.md`):

1. [`references/2026-10-20-duas-contradicoes-nas-referencias-review-e-events-da-skill.md`](references/2026-10-20-duas-contradicoes-nas-referencias-review-e-events-da-skill.md),
   origin quaestor, type hint `fix`, transcribed from an ai-memory message on
   2026-10-21. On 2026-10-19 the quaestor pre-PR review of its Baseline
   refresh (`roundfix review`, provider codex) reported both contradictions in
   the references the setup copied from Roundfix 0.91.0, and quaestor
   corrected nothing because the files are setup-managed. The entry's check
   on `main` `956066a7` (0.101.0) found both still present.

No open Backlog Entry or unresolved Finding shares this context: the only open
Backlog Entry at the authoring base,
`docs/backlog/2026-10-20-the-queue-archive-subject-and-pull-request-title-grow-past-100-characters-with-the-slug.md`,
is about the Delivery Queue's commit subjects, and `docs/findings/` holds no
entry.

## Prerequisites

None. This Spec is authored against `8fb2da1f` (v0.102.0 plus the CI timeout
change) and delivers after the cycle AM Specs 0292 and 0293. Spec 0291, merged
in that base, changed other paragraphs of the `deliver` user guide and the
Roundfix `deliver` reference; this Spec's edits there are two sentences about
archived Specs, which task_01 locates by their text, not by line number.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, schema or record
  shape is minted or changed; the documents name the existing
  `roundfix/delivery-conventions/v3` and `roundfix-events/v1` versions.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the change is documentation and
  documentation contract tests that read repository files. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — the Spec records no decision and
  changes no behavior; it brings documents in line with conventions and
  categories the code already implements, and no accepted ADR constrains how
  a document is checked against the code. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AN approval of 2026-10-21 through AskUserQuestion, authorization "the same
  rules as cycle AM", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0294-the-roundfix-skill-references-agree-with-the-cli/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config file change.

## Goals

- The Roundfix `review` reference and the `review` user guide give the five
  Delivery Conventions exactly as the reviewer's prompt carries them, and
  describe exactly when C3 and C5 apply.
- The Roundfix `events` reference lists every Run Event Stream category the
  CLI replays by default and accepts in `--filter`, `usage` included, with a
  stable-fields row for each.
- Every other statement in the Roundfix skill, the `review`, `events` and
  `deliver` user guides and the glossary that describes an archived Spec as a
  moved folder, an archived `_prd.md` as the only archive form, or the
  category list states what the code does.
- A documentation contract test fails when a convention, the conventions
  version or a category changes in the code and the documents do not follow.

## Core Features

1. Review conventions. The `review` reference and user guide quote C1–C5 and
   the conventions version verbatim from the reviewer's prompt and describe
   C3 and C5 eligibility for an Archive Record and for a Legacy Archive Folder.
2. Event categories. The `events` reference lists `task-status`, `batch`,
   `verification`, `outcome`, `agent-selection` and `usage` as the default
   replay and the `--filter` vocabulary, and its stable-fields table has a row
   for each.
3. The sweep. The `deliver` reference and user guide name the Archive Record
   where they name only an archived `_prd.md`; the glossary's Run Event Stream
   entry lists the six categories.
4. Drift guards. Documentation contract tests read the conventions from the
   reviewer's source and the categories from the Run Event Stream package and
   compare them with the references, their mirrors, the user guides and the
   glossary.

## Non-Goals / Out of Scope

- Changing the reviewer's conventions, their version, the review validation,
  the Run Event Stream, the `events` help text or any other behavior.
- Editing the upstream skills repository's copy of the Roundfix skill, or any
  adopter's copy; adopters receive the fix through `roundfix baseline update`.
- Rewriting any Archive Record or any file under `docs/history/`.

## Success Metrics

1. Success Metric: the `review` reference, its mirror and the `review` user
   guide each contain the five conventions exactly as the reviewer's prompt
   carries them, and a documentation contract test that reads them from the
   reviewer's source fails when the old C3 sentence is planted.
2. Success Metric: the `events` reference and its mirror list the six
   categories in the default-replay sentence and give each a stable-fields
   row, and a documentation contract test that reads the categories from the
   Run Event Stream package fails when `usage` is removed from that sentence.
3. Success Metric: the glossary's Run Event Stream entry lists exactly the six
   categories, enforced by the glossary closed-list contract.
4. Success Metric: no statement in the Roundfix skill, the `review`, `events`
   and `deliver` user guides or the glossary describes an archive as moving
   the Spec folder, names an archived `_prd.md` as the only archive form, or
   lists a category set other than the six.
5. Success Metric: after the next release, the next quaestor pre-PR review of
   a Baseline refresh that carries these references reports neither
   contradiction.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-21.

- The adopted Inbox Entry, written for quaestor on 2026-10-20 and copied
  verbatim into `references/`: the two contradictions a Codex reviewer found
  in an adopter's setup-managed copy, and the lines that hold them.
- Nassif, Hernandez, Sridharan and Robillard, "Generating Unit Tests for
  Documentation", IEEE Transactions on Software Engineering, 2021
  (https://doi.org/10.1109/tse.2021.3087087): "Without supporting technology,
  such redundancy easily leads to inconsistencies and a degradation of
  documentation quality."
- Daniel Metzner, "A test that reads your docs", 2026-09-04
  (https://til.metzner.uk/posts/a-test-that-reads-your-docs/): "Read the fact
  from the code and there is exactly one copy of the truth — and it is not the
  copy being tested."

## Unreachable Acceptance

- criterion: Success Metric 5, after the next release, the next quaestor
  pre-PR review of a Baseline refresh that carries these references reports
  neither contradiction.
  reason: no Verification or QA row may read another repository, the real
  `~/.roundfix` or a provider, and only quaestor's review of its own refresh
  can show it.
  satisfied-by: the quaestor operator's next Baseline refresh and pre-PR
  review after this Spec's release.

## Validation

The decisions were checked on 2026-10-21 against these sources.

- Characterization at `956066a7`, rechecked at `8fb2da1f`, by reading the
  reviewer's conventions and their region code, the Run Event Stream package,
  the `events` command, the Delivery Queue's prerequisite and merge proofs,
  and every statement the sweep covers, and by running the built
  `bin/roundfix events --help`, whose `--filter` line reads
  `Comma-separated categories:
  task-status,batch,verification,outcome,agent-selection,usage`. The TechSpec
  records each finding with its file and function.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "roundfix
  skill reference review conventions C3 C5 events usage filter" --all --files
  --min-score 0.3`, whose first hit was the upstream skills repository's
  mirror of this reference (`projects/skills/mirror/skills/06-review-repair/roundfix/references/review.md`),
  which carries the same stale C3 sentence at line 157, the same C5
  eligibility sentence at line 181 and, in its `events.md`, the same
  five-category list at line 26; and `qmd query "documentation drift skill
  reference contradicts CLI contract test" --all --files --min-score 0.3`,
  which returned the Roundfix mirror's records of the doctor readiness
  contract reconciliation and the repository-contract ADRs, confirming that
  documentation contracts run on every push to main and before a release. No
  Inbox Entry was written, because the briefing forbids writing there.
- Exa: the two published sources above, and a layered documentation-test
  model (https://metronforge.dev/posts/documentation-as-code-needs-documentation-tests/)
  whose Interface layer is "compare documented options with schemas or CLI
  help". Challenging: Metzner warns that a search for a name is satisfied by
  prose that discusses the name, so the guards here compare whole sentences
  and parsed values, never a bare word.

## Glossary

- changes: **Run Event Stream**

## Decisions

- 2026-10-21 (about 01:00 BRT): the cycle AN scope "V → W → X (Recomendado)"
  (option text "Três Specs: v0.105 a v0.107. Y fica no backlog como entrada
  aberta; são só 3 linhas de ruído no reconcile e não custam nada.") and the
  authorization "the same rules as cycle AM".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AN briefing for Spec 0294 of 2026-10-21 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The documents quote the reviewer's conventions verbatim, so a test can
  compare them whole; the user guide's paraphrase is replaced too.
- The guards parse the reviewer's source file instead of exporting a symbol,
  because the repository forbids an exported symbol that exists only for
  tests.
- The sweep also corrects the `deliver` reference and user guide and the
  glossary's Run Event Stream entry, which state the same stale facts.
- No ADR: nothing is decided that the code has not already decided.

## Open Questions

None.

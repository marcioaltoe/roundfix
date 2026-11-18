---
type: feat
status: open
created: 2026-11-18
spec: null
reason: null
---

# The Baseline owns each project's documentation rules and what it exports to the Secondbrain, and records Lessons Learned

## Request

Two ai-memory messages from the Indexus project (`gesttione-solutions-indexus`), both on 2026-11-18, relaying the maintainer's decision. The second widens the first.

- `01a1eec1-4be0-72e0-8e9f-4bcf7211f8a9` (13:03 UTC): the maintainer asked for a Backlog Entry: "Pensar em como fazer as lições aprendidas serem algo de valor para ir para o secondbrain para serem ingeridas no indexus." His principle: "Só seria útil se fosse algo como: tentamos usar x lib e gerou problemas tal e resolvemos assim." Expected: a short, searchable Lesson Learned format (context; what we tried; what went wrong; the cause; what solved it; when the lesson applies and when not; links to the Spec, Finding and PR); a natural moment of the flow that produces one only when there is a lesson (candidates: the Spec archive after a Finding, retry or failed QA; closing a Finding; a review with a root cause); its own place under `docs/` that the Secondbrain mirror carries to the Indexus; and no new process archive — nothing is written when there is no reusable lesson.
- `01a1eeca-397e-77c2-8429-3668a5a0b8f8` (13:13 UTC): the maintainer decided that Roundfix defines, through the Baseline, the rules every project follows — the content of documentation files, the folder layout, the required files, what is kept and what is not, and what each project exports to the Secondbrain — applied by `roundfix update` and checked by `roundfix doctor`; no other project creates these rules on its own. The Indexus measured 4,150 mirrored files of 11 projects: `docs/agents/` (~130, no value: Baseline copies), the `skills` repository (779, none), `docs/history/` (~1,100, low), `docs/backlog/`, `docs/specs/`, `docs/findings/` (~170, low: transient), `docs/adr/` (~1,040, medium; Roundfix alone 363), `docs/references/` (~620, high), README, `CONTEXT.md`, `docs/user-guide/` (~110, high). Today each `.secondbrain-export` exports all of `docs/`; the only Baseline rule known is excluding `docs/history/specs/*/qa/`. Asked of the Baseline (the format is Roundfix's decision): (1) a standard export list in `.secondbrain-export` holding only knowledge — suggested in: README, `CONTEXT.md`, `docs/user-guide/`, `docs/references/`, current ADRs, Lessons Learned; out: `docs/agents/`, `AGENTS.md`, `docs/history/`, `docs/backlog/`, `docs/specs/`, `docs/findings/`; and whether purely technical ADRs leave, or an ADR marks whether it interests the knowledge base; (2) the Lesson Learned above; (3) a clear status for superseded ADRs so no expired decision is cited as current; (4) `doctor` reports an export list off the standard, a Lesson Learned off the format and a missing required file, and `update` applies the standard; (5) say what Roundfix itself exports — proposed: only README and `docs/user-guide/`.

## Notice Request

The sender asked to be told when the work is done (Baseline published and applied in the projects), with a message to `default/gesttione-solutions-indexus` carrying: the final export list; the Lesson Learned format and path; the superseded-ADR marking; and the projects already updated. A partial plan before the end also helps it prepare the Indexus. Meanwhile the Indexus asked the Secondbrain to filter on its side.

## Decisions

- 2026-11-18 (`/roundfix --follow-up`, AskUserQuestion): "Spec com etapa de Ideia (Recommended)"; grouped with `docs/backlog/2026-11-16-the-docs-folders-of-adopters-have-unclear-jobs.md` — "Juntar à Spec do Indexus (Recommended)"; authorization "Mesmas do BE (Recommended)" (to be recorded in the Spec's `_authorization.md`); delivery order: after the corrective qa-gate Spec, the Delivery Queue fix and the Spec on questioning rules at the Idea and PRD stages.
- Idea-stage judgment (Jev, 2026-11-18): undecided (shape_open 0.68, unverified_assumption 0.17, market_question 0.06); the maintainer chose the Idea stage.

## Where

The Baseline (`docs/agents/docs-layout.md` and the guides it governs), `.secondbrain-export`, `roundfix baseline update` and `roundfix doctor`, and the archive / Finding / review flow that would produce a Lesson Learned.

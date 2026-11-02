---
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
status: active
created: 2026-11-01
surfaces: [backend, docs]
---

# QA evidence links resolve in the committed tree

Minimal PRD for a bug fix, minted by `write-techspec`. On 2026-11-01 an Argus
QA Report settled `completed` while it linked two `.log` evidence files that
the repository's `.gitignore` excludes, so the QA Report commit never held
them, the Run Worktree that did was removed, and `roundfix spec check
--strict` reported `SC-HISTORY-RETAINED` for both links only after
integration. After this Spec every QA settlement path, the Daemon's QA
settlement, `roundfix qa-report accept` and the Settle Command for a `qa`
Task, refuses a report whose local link does not resolve in the tree the
commit would hold, naming the link, its line and the reason `ignored` or
`missing`. The Daemon first gives the gate one repair turn, and the `qa-gate`
skill tells the gate to keep evidence under names the repository does not
ignore.

## Sources

One adopter report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-01-qa-evidence-links-resolve-in-the-committed-tree.md`](references/2026-11-01-qa-evidence-links-resolve-in-the-committed-tree.md),
   `type: fix`, Argus, ai-memory message `01a19895`.

At the authoring base `docs/findings/` holds no unresolved Finding on QA
evidence links, and no other Backlog Entry covers them.

## Prerequisites

None in `requires`. Authored against `dd3b2b62` (v0.128.0), where no other
Spec is active. Specs 0338 and 0339 of the same cycle are authored in
parallel and deliver first; they share no Verification artifact with this
Spec, only files that `_techspec.md` → Build Order names.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the refusal reasons are two fixed words. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the check reads local files and the local Git index. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0351, new, records the decision:
  ADR-0351: "A local link resolves when the committed tree would hold its
  target." It names the three settlement paths that apply the rule, the one
  repair turn of the Daemon, the refusal at the QA Report commit boundary and
  `roundfix archive` left unchanged. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AU answers of 2026-11-01 quoted in the authorization record, which keep
  the cycle AT and AS rules and the cycle AJ grant for the Governed Paths
  each Spec declares; the governed files are the canonical `qa-gate` and
  Roundfix skill files and their `SKILL.md` pairs. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0340-qa-evidence-links-resolve-in-the-committed-tree/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/settle.md`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- No QA settlement path accepts a report whose local link the committed tree
  would not hold.
- Each refusal names the link as written, its line and the reason `ignored`
  or `missing`, so the gate keeps or renames the evidence instead of losing
  it.
- The Daemon returns the gate once to repair the links before it refuses.
- The QA settlement and `SC-HISTORY-RETAINED` read a report's links through
  one extraction.
- The `qa-gate` skill tells the gate to save evidence under names the
  repository does not ignore and never to cite a file it did not keep.
- Report truth is kept: no verdict, row or existing check is weakened.

## Core Features

1. The shared link check. `spec.RelativeLinks` extracts every relative
   Markdown link outside code; `MissingRelativeLinks` filters it for
   `SC-HISTORY-RETAINED` with unchanged results; `spec.UnresolvedEvidenceLinks`
   classifies each link of a QA Report as resolved, `ignored` or `missing`.
2. The CLI paths. `roundfix qa-report accept` and the Settle Command for a
   `qa` Task refuse an Unresolved Evidence Link after the eligibility
   decision, with the shared refusal text.
3. The Daemon path. After the QA turn the Daemon checks the newest report,
   returns the same Agent Session once with the links and reasons, checks
   again at the QA Report commit boundary, and on a remaining link settles
   the QA Task `failed`, publishes one QA Run Event and creates no QA Report
   commit.
4. The rules in writing. A new `qa-gate` section, the qa-report, settle and
   implement guides, the Roundfix skill's settle and implement references,
   the glossary and ADR-0351.

## Non-Goals / Out of Scope

- External URLs, anchors and absolute-path destinations; evidence cited only
  as a code span, which `SC-HISTORY-RETAINED` does not read either.
- `roundfix archive`, which judges a committed report.
- Staging ignored files with `git add -f`, renaming evidence in the Daemon,
  or editing the report for the gate.
- A repair turn for the Secret Scan or the non-image binary refusal.
- The `qa-report accept` usage text, the eligibility policy
  `spec.QAReportEligibility`, any `### QA settlement` section.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, any Baseline module or
  guide under `docs/agents/`.

## Success Metrics

1. Success Metric: in a temporary Git repository whose `.gitignore` holds
   `*.log`, `roundfix qa-report accept` on an eligible `pass` report linking
   `evidence/consumer-tests.log` (present, ignored) and
   `evidence/journeys.txt` (absent) exits 1 with the stderr line
   `roundfix: QA Report eligibility refused: QA Report evidence links do not resolve in the committed tree: evidence/consumer-tests.log (line <n>, ignored); evidence/journeys.txt (line <m>, missing)`;
   at `dd3b2b62` it exits 0.
2. Success Metric: the same report with its evidence saved as tracked or
   untracked unignored `.txt` files is accepted; a tracked file that matches
   an ignore pattern resolves; a link inside a code span, a URL and an anchor
   are never checked; a relative link leaving the repository is `missing`;
   outside Git a present target resolves.
3. Success Metric: the Settle Command for a failed `qa` Task whose newest
   report holds an ignored link exits 1 with
   `roundfix: settle QA Report is ineligible: QA Report evidence links do not resolve in the committed tree: … (report <path>)`
   and leaves the Task file unchanged.
4. Success Metric: in a Daemon Task Cycle whose QA turn writes a `pass`
   report linking an ignored `.log`, a repair turn that saves the evidence as
   `.txt` and fixes the link leads to an accepted QA Task whose QA Report
   commit holds the `.txt`; exactly one repair turn runs.
5. Success Metric: when the repair turn leaves the link ignored, the QA Task
   settles `failed` with the reason
   `QA Report evidence links do not resolve in the committed tree: evidence/consumer-tests.log (line <n>, ignored)`,
   one `evidence_links` QA Run Event with outcome `refused` names the link,
   no QA Report commit is created, and the report and the `.log` stay in the
   Run Worktree.
6. Success Metric: `SC-HISTORY-RETAINED` reports the same findings as before
   on its existing tests, through the shared extraction.
7. Success Metric: both copies of the `qa-gate` skill carry the new section;
   the qa-report, settle and implement guides, both copies of the Roundfix
   settle and implement references, `CONTEXT.md` and ADR-0351 state the rule.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Argus report quoted in the adopted reference, written by the Argus
  project before this Spec existed.
- The Git manual, `git-check-ignore` (https://git-scm.com/docs/git-check-ignore,
  read through Exa on 2026-11-01): "By default, tracked files are not shown at
  all since they are not subject to exclude rules". The `gitignore` manual
  (https://git-scm.com/docs/gitignore): "A `gitignore` file specifies
  intentionally untracked files to ignore. Files already tracked by Git are
  not affected". The resolution rule follows both: a tracked file resolves
  whatever the ignore patterns say, and an untracked ignored file never
  reaches a commit.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files, temporary
repositories and the fake QA Agent of the Daemon tests.

## Validation

The decisions were checked on 2026-11-01 against these sources.

- Characterization at `dd3b2b62` with that commit's `bin/roundfix` in a
  temporary Git repository whose `.gitignore` holds `*.log`: an eligible
  `pass` report linking an ignored `evidence/consumer-tests.log` and an absent
  `evidence/journeys.txt` is accepted by `roundfix qa-report accept`, exit 0.
- Reading at `dd3b2b62`: `runQAGate` (`internal/daemon/task_engine.go`)
  settles from `settleQAVerdict`, then refuses at the QA Report commit
  boundary only a Secret Match or a non-image binary; the QA step has no
  repair turn. `requireSettleQAReportEligibility` (`internal/cli/settle.go`)
  and `runQAReportCommand` (`internal/cli/qa_report.go`) apply only
  `spec.QAReportEligibility`. `MissingRelativeLinks`
  (`internal/spec/record_links.go`) is the link reader of
  `SC-HISTORY-RETAINED`.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "QA evidence ignored by gitignore link missing after commit roundfix qa report" --all --files --min-score 0.3`.
  It returned Roundfix mirrors (the ADR on QA evidence binaries, Spec 0263 on
  archive keeping what Git ignores, the qa-report guide) and the digest
  `wiki/sources/digest-roundfix-git-add-f-apaga-arquivo-2026-09-17.md`, a
  `git add -f` staging defect that supports not forcing ignored evidence into
  a commit; none decides the link rule.
- Exa: the two Git manual pages quoted above.

## Glossary

- adds: **Unresolved Evidence Link** — a QA Report link the committed tree would not hold
- changes: **QA Report** — settles only when its local links resolve

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AU brief for Spec 0340 of 2026-11-01 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-01), verbatim:

- Scope: "3 Specs (Recommended)" — option text: "0338 release: tipo do PR
  segue a Spec + teto de release no _authorization.md (os dois mexem no
  release plan e no write-*); 0339 apply aprovado apaga as aposentadas (nova
  ADR, context7-cli/exa-web-search ficam); 0340 links de evidência da QA
  existem no commit. Fila nessa ordem."
- QA evidence: "Recusar no settle (Recommended)" — option text: "Antes do
  settle da QA, cada link local de evidência precisa existir na árvore que
  vai ser commitada; link ignorado ou ausente falha com o caminho e o motivo,
  e o Agente tem um turno de reparo. A skill qa-gate manda salvar evidência
  com extensão que o repositório não ignora (.txt)."
- Release of the cycle: "Uma no fim, v0.129.0 (Recommended)" — option text:
  "Depois das três Specs, com aviso aos adotantes; o apply que apaga
  aposentadas e o teto de release chegam juntos pela Baseline."
- Authorization: "Mesmas do AT (Recommended)" — option text: "Push e PR dos
  planos, entrega e merge pela fila, Makefile/workflows/go.mod/lint quando o
  desenho pedir, 'Consertar e seguir' em falha técnica, mudança da Baseline
  (agent-instructions, release, spec-routing, skills de autoria) que chega
  aos adotantes, release v0.129.0 no fim e aviso aos adotantes."
- Start: "Já (Recommended)".

Taken from those answers and the characterization, with no further
interview:

- One rule for every settlement path: the Settle Command settles a `qa`
  Task too, so it applies the check beside `qa-report accept` and the Daemon.
- The resolution rule asks Git what it would commit, so a tracked file that
  matches an ignore pattern still resolves, and the reason distinguishes an
  ignored file from an absent one.
- The repair turn is the QA counterpart of a Task's Verification Feedback:
  one turn in the same Agent Session, then the commit-boundary check decides.
- A remaining link refuses the QA Report commit, as a Secret Match and a
  non-image binary do, so nothing the gate kept is lost or half-committed.
- The extraction moves into `spec.RelativeLinks`, shared with
  `SC-HISTORY-RETAINED`.
- Three implementation Tasks in series, because each raises an owned-skill
  version or edits the same guides, and a QA gate.

## Open Questions

None.

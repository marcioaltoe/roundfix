---
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
status: active
created: 2026-11-04
surfaces: [backend, cli, docs]
type: fix
---

# A secret refusal of a Task commit is repairable and visible

The Daemon scans the lines a Task commit would add for a Secret Match only
after the Task's last Verification attempt, when its one repair turn is already
spent, so a test fixture whose value matches a Secret Pattern fails the Task
with no chance to repair it, and the refusal is visible only in the Run
Database. On 2026-11-04 Fiscus lost a Run this way: every Verification command
and both Settlement Checks passed, `roundfix runs causes` listed nothing, the
Delivery Queue parked the item `run-unresolved` with the generic next action,
and the dependent Tasks were skipped without a visible reason. After this
Spec the scan runs as a Settlement Check of every Verification attempt, so a
Secret Match on the first attempt returns to the Agent as Verification
Feedback naming path, line and pattern kind, never the value; the Causes
report lists each secret refusal as its own item; and a `run-unresolved` park
names each failed Task with its reason and the Tasks skipped behind it.

## Sources

One source, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-04-fiscus-a-task-commit-secret-refusal-is-invisible.md`](references/2026-11-04-fiscus-a-task-commit-secret-refusal-is-invisible.md),
   reported by fiscus, ai-memory message `01a1a775` of 2026-11-04, copied
   verbatim: Run `run_20261104T163352Z_9eb4e9a0c3b72f10` of Fiscus Spec 0077
   on roundfix 0.130.1, its three problems and its two suggestions.

No open Backlog Entry or unresolved Finding exists on `main`, so no other
source shares this context.

## Prerequisites

None in `requires`. Authored against `de579ea2` (v0.131.0), where no Spec is
active, and delivered alone.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the Causes
  report gains the item kind `secret-scan` and the Daemon gains the check
  label `settlement check: secret scan`, both fixed strings. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the scan, the report and the park read only the Run
  Worktree, the Run Database and Git. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0292 scans each Task commit and
  rejected a repair turn for this change only, ADR-0292: "a new repair class
  is a separate decision"; this Spec takes that decision and records it in
  ADR-0361, written here. ADR-0273, which ADR-0292 amends, keeps the value out
  of every diagnostic, ADR-0273: "it never carries the matched value or any part of it"
  holds for the feedback, the report and the park. ADR-0182 and ADR-0342
  define the Settlement Checks the scan joins, ADR-0342: "The repository
  Verification at settlement applies to every Task Graph"; ADR-0214 classifies
  causes by deterministic signatures, ADR-0214: "classified by deterministic
  signatures", and a secret-scan item is classed by its kind, which is also
  deterministic. Set no obligation here, because this Spec changes none of their
  rules and each reaches it only through citations of the ADRs above: ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0227 (a review through Claude keeps its verdict and fits its window); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0256 (a release waits for the skills that describe a changed surface); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0283 (a Run Marker keeps the judge silent inside a Run); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof); ADR-0291 (qA evidence holds text and images, and the QA Report commit refuses any other binary); ADR-0295 (reachability reads the fetched default branch, and origin names it when origin/HEAD is unset); ADR-0299 (the history conversion rewrites every citation it lists, and the check reports a removed history file cited without a Git Locator); ADR-0306 (archive and cancel refuse a path whose resolved target leaves the repository); ADR-0308 (an archive rebases the unproven links it records, and the Record Link Repair points a promoted target at its copy); ADR-0317 (a Result Status Note of a completed Task is dismissed without the validator); ADR-0319 (the follow-up order is code over four Jev yes/no answers); ADR-0345 (a derived-record writer names its command to the suite guard); ADR-0348 (an absent Spec Root inside the repository holds no Spec); ADR-0349 (a Spec's Release Ceiling answers the release plan's approval); ADR-0352 (the pre-PR review omits declared paths, and a large Claude candidate is reviewed in segments); ADR-0353 (a waiting item is held back from the follow-up list unless --all is asked). Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AY answers of 2026-11-04 quoted in the authorization record, which keep the
  cycle AX rules; the governed files are the Roundfix skill files task_04
  changes. No Makefile, workflow, Go module, lint, formatter or Baseline
  module change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/runs.md`, `skills/roundfix/SKILL.md`.

## Goals

- A Secret Match on a Task's first Verification attempt returns to the same
  Agent Session as Verification Feedback, and the second attempt scans again.
- A Secret Match on the final attempt still settles the Task `failed` with the
  reason `Task commit carries a secret: <path>:<line>: <kind>` and no Task
  commit, and no third attempt exists.
- `roundfix runs causes` lists every secret refusal of a terminal Spec Run as
  its own item, with Spec, Task, Run, attempt, path, line and kind, in text
  and JSON.
- A `run-unresolved` park names each failed Task with its reason and every
  skipped Task, in `roundfix deliver status` and the Pending Question, and its
  next action names the repair when a Task commit carried a secret.
- No value, nor any part of one, reaches a diagnostic, a prompt, an event, the
  report or the park.

## Core Features

1. The scan is a Settlement Check. Every Verification attempt of a non-QA
   Task, with or without a QA gate, runs `settlement check: secret scan` after
   the attempt's commands and after the two in-process Settlement Checks, over
   the paths the Task commit would stage. A match fails the attempt with the
   diagnostic `Task commit carries a secret: <locations>` and a `Fix:` line
   telling the Agent to use a value no Secret Pattern matches; a Git or read
   error still halts the Run.
2. The repair turn. A failed first attempt sends the scan's diagnostic with
   the other failures in the one Verification Feedback turn; the second
   attempt runs every command and check again.
3. The final refusal. A match on the final attempt settles the Task `failed`
   with today's reason bytes when the scan is the attempt's only failure, and
   publishes the `secret_scan` event on each attempt that matches, now with
   its attempt number.
4. The Causes report. Each `secret_scan` event of a terminal Spec Run becomes
   one item of kind `secret-scan`, class `repository-convention`, whose check
   is its locations and whose JSON carries `matches`; the attempt's
   `verification` item for the same check is not listed twice. Events written
   by 0.130.1 and earlier are listed too.
5. The park. A `run-unresolved` park records
   `run-unresolved: <task> failed: <reason>; skipped: <tasks>` from the Run's
   settle and skip events; a Run with neither keeps the bare blocker, and the
   Park Class stays `finding`.
6. The rules in writing. The Roundfix skill's `### Settlement Checks`
   section, its implement, runs and deliver references, the implement, runs
   and deliver user guides, `CONTEXT.md` and ADR-0361 state the behavior.

## Non-Goals / Out of Scope

- Changing a Secret Pattern, the added-lines rule or the binary rule. The
  maintainer declined accepting a value that tracked files already hold:
  ADR-0292: "a real credential in a test file leaks exactly as anywhere else".
- A second repair turn, or a repair turn for the QA Report commit or the Final
  Push: the QA path keeps `QA evidence carries a secret`, the push keeps
  `secret-found`.
- New Park Classes, a new blocker name, or a structured blocker column in the
  Run Database.
- New signatures in the signature table, whose digest stays unchanged.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, lint or formatter
  configuration, any Baseline module and any `### QA settlement` section.

## Success Metrics

1. Success Metric: a Task whose Agent writes a `credential-assignment` value on
   its first turn and a letters-only value on its repair turn settles
   `completed`; its repair prompt names `settlement check: secret scan`,
   `<path>:<line>: credential-assignment` and the `Fix:` line, and no eight
   consecutive bytes of the value occur in any prompt, diagnostic file, event
   or progress line.
2. Success Metric: a Task whose Agent writes the value on both turns settles
   `failed` with the reason `Task commit carries a secret: fixture.txt:2:
   credential-assignment`, creates no commit, publishes two `secret_scan`
   events with attempts 1 and 2, and receives one repair turn.
3. Success Metric: in a gated graph the attempt's checks run in the order
   `settlement check: spec consistency`, `settlement check: authorization`,
   `settlement check: secret scan`; without a QA gate only the scan runs after
   the commands.
4. Success Metric: `roundfix runs causes` over a Run with a 0.130.1-shaped
   `secret_scan` event and a Run with the new shape prints one `secret-scan`
   line per event and no `verification` line for the scan check, and its JSON
   items carry `matches` with path, line and kind.
5. Success Metric: an Implement Run whose `task_01` failed by a Secret Match
   and whose `task_02` and `task_03` were skipped parks
   `run-unresolved: task_01 failed: Task commit carries a secret: fixture.txt:2: credential-assignment; skipped: task_02, task_03`,
   and `roundfix deliver status` prints it with the secret next action; a Run
   with no failed or skipped Task parks `run-unresolved` with today's answer.
6. Success Metric: the skill references, the user guides and `CONTEXT.md`
   state the behavior and each mirror equals its canonical copy.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Fiscus's report, quoted in the adopted reference and written before this
  Spec existed: on roundfix 0.130.1, Run
  `run_20261104T163352Z_9eb4e9a0c3b72f10` passed every Verification command
  and both Settlement Checks, failed `task_01` on a `credential-assignment`
  Secret Match in a new test, skipped its dependents and parked
  `run-unresolved` with "finding: resolve the blocker", and `runs causes`
  did not list the Run.
- GitHub's push protection documentation, "Working with push protection from
  the command line"
  (https://docs.github.com/en/code-security/how-tos/secure-your-secrets/work-with-leak-prevention/push-protection-on-the-command-line,
  read through Exa on 2026-11-04): a blocked push lists each location as
  "path: README.md:4" under the secret's type, and the page "Push protection"
  says it will "provide a detailed message explaining the reason for the
  block". A refusal that names type and location without the value, so the
  author can repair it before trying again, is the published practice.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files, temporary
repositories and temporary Run Databases.

## Validation

The decisions were checked on 2026-11-04 against these sources.

- Reading at `de579ea2`: `executeTask` (`internal/daemon/task_engine.go`)
  runs Verification attempt 1, `repairTaskVerification` and attempt 2, then
  `prepareTaskCommit`, the authorization refusals and only then
  `scanTaskCommitPaths`, whose match sets the reason `Task commit carries a
  secret: ` and publishes the `secret_scan` event without an attempt number.
  `verifyTask` adds the two in-process Settlement Checks only when the graph
  has a QA gate (`plan.settlementChecks`). `runcause.Build` reads only failed
  `daemon.verification` events and corrective Tasks. `runCandidate`
  (`internal/delivery/engine.go`) parks the bare `run-unresolved`, and
  `ClassifyPark` answers it "resolve the blocker, then run roundfix deliver
  retry <slug>".
- Characterization at `de579ea2` with a temporary test, removed afterwards:
  a terminal Run holding a `task_01` start, a `secret_scan` event and a failed
  settlement produced no Causes item. `TestAnUnresolvedRunParksWithItsRunID`
  and `TestExistingGenericParkAnswersAreUnchanged` pin today's bare blocker
  and its generic answer.
- A prototype of Core Features 1 to 3 in a scratch worktree, removed
  afterwards: the two cases of
  `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit` kept their
  reason bytes and published two refusals instead of one;
  `TestTaskCycleVerificationSequenceStopsAtFirstFailureWithOneVerdictPerAttempt`
  and `TestTaskCycleExecutesAgentVerifySettleCommitContract` of
  `internal/daemon` and
  `TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree`
  of `internal/cli` failed only on the new check's events and snapshot; the
  other 750 tests of `internal/daemon` and the `internal/runcause` and
  `internal/delivery` suites passed. `_techspec.md` → Declared breaks names
  each.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix secret scan task commit refusal repair verification feedback runs causes park" --all --files --min-score 0.3`,
  which returned only Roundfix and skills project mirrors (the CHANGELOG,
  ADR-0292, two unrelated ADRs, the implement guide and the roundfix and
  implement-task skills) and `wiki/sources/project-mirror-deltas-2026-09-28.md`.
  They confirm ADR-0292 left the repair turn as a later decision; none treats
  the Causes report or the park.
- Exa: the GitHub pages quoted above. The "Working with push protection in
  the GitHub UI" page also says GitHub does not block a secret it already
  alerted on, which supports Fiscus's tracked-value suggestion; the
  maintainer declined it because a tracked value can itself be a leaked
  secret, and "Allowlisting Test Fixtures Without Blinding the Scanner"
  (https://www.git-automation.com/commit-signing-supply-chain-security/secret-scanning-and-remediation/allowlisting-test-fixtures-without-blinding-the-scanner/)
  gives the same advice: "never allowlist a value that is or was a real
  credential", and make generated values obviously fake.

## Glossary

- changes: **Settlement Check** — the Task commit's Secret Scan is one, in every attempt
- changes: **Secret Scan** — the Task commit's scan returns a first-attempt match as Verification Feedback
- changes: **Cause Class** — a secret-scan item is classed by its kind

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AY brief for Spec 0351 of 2026-11-04 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-04 ~22:05 BRT), verbatim:

- Repair: "Scanner entra nos checks (Recommended)" — option text: "O scanner
  roda junto com os settlement checks da tentativa 1. A recusa
  (arquivo:linha:padrão, sem o valor) volta como Verification Feedback no
  único reparo, e a tentativa 2 roda o scanner de novo. Mantém a regra de um
  reparo só."
- Tracked identical value: "Não aceitar (Recommended)" — option text: "O
  scanner fica como está, olhando só as linhas adicionadas. O feedback diz
  como trocar por um valor que nenhum padrão casa. Um valor já rastreado
  também pode ser um segredo vazado, então copiá-lo não deveria passar."
- Rules: "Mesmas do AX (Recommended)" — option text: "Sem supervisão; merge
  pela fila; corrigir e continuar nas paradas; Makefile, CI, go.mod e lint
  liberados; mudanças de Baseline e skills que chegam aos adotantes
  liberadas."
- Release: "v0.132.0 ao fim (Recommended)" — option text: "Release minor
  depois do merge, com a nota aos 8 adotantes e a resposta ao Fiscus. É o
  único release de hoje, dentro do teto de 1 por dia."

Taken from those answers, ADR-0292, ADR-0342 and the reading above, with no
further interview:

- The scan is the last check of each attempt, after the commands and the two
  in-process Settlement Checks, and it runs with or without a QA gate, as it
  does today. See ADR-0361.
- The scan no longer runs after the final attempt: the attempt that passes
  has already scanned the same prepared paths. See ADR-0361.
- A Git or read error of the scan halts the Run, as today; it is never
  Verification Feedback.
- The final reason keeps today's bytes when the scan is the only failure;
  with other failures it joins them in check order, the scan last.
- The Causes item kind is `secret-scan`, its class is `repository-convention`
  by kind, and its signature and evidence are null. See ADR-0361.
- The park detail is built from the Run's own settle and skip events, each
  reason cut to 300 characters. See ADR-0361.
- Four implementation Tasks (the scan, the report, the park, the documents)
  and a QA gate.

## Open Questions

None.

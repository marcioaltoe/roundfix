---
status: approved
granted: 2026-10-30
action: move every upstream Setup Snapshot to catalog commit 86df860, retire the seven skills the catalog removed and the owned implement-spec from the Baseline modules, activation bundles, snapshots and owned bundle, restore this repository's three changed upstream skills, reword the person-started and polish triggers and split the frontend bundles, make the Task and QA prompts, the Spec Context Bundle, implement-task and qa-gate carry the skill load rule, move implement-spec's two rules into the Roundfix skill, and regenerate the derived Baseline files and this repository's managed guides through the sanctioned commands
consuming: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/context7-cli/references/docs.md
  - .agents/skills/context7-cli/SKILL.md
  - .agents/skills/grilling/SKILL.md
  - .agents/skills/implement-spec/SKILL.md
  - .agents/skills/implement-task/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/reconcile.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/spec-delivery.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/tui-design/references/evidence.md
  - .agents/skills/tui-design/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - Makefile
  - docs/agents/setup-context.json
  - docs/agents/skill-dispatch.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/modules/backend.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/frontend.json
  - internal/baseline/assets/modules/typescript.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
  - internal/baseline/assets/skill-activations.json
  - skills/archive-spec/SKILL.md
  - skills/implement-spec/SKILL.md
  - skills/implement-task/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0333

On 2026-10-30 the maintainer approved cycle AS through AskUserQuestion, after
a follow-up inventory (`roundfix follow-up --messages`, judge
`typesafe/jev-1.13-20260917`) and a direct Jev priority check (`jev-1.14.0`).
The answers, verbatim:

- Plan: "Aprovar como está (Recommended)" — option text: "Quatro Specs na
  ordem validada pelo Jev; a 0332 entra primeiro para que a fila resolva
  sozinha conflitos de registros gerados durante a 0333."
- Authorization: "Mesmas regras do AR (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."
- Release: "Uma no fim do ciclo (Recommended)" — option text: "v0.127.0
  depois das quatro Specs; os adotantes recebem a revisão de skills e as
  correções juntas."

The scope of this Spec comes from the maintainer's own words in chat on
2026-10-30, verbatim:

- "Vamos manter testing-boss, vou solicitar somente que remova tdd por
  enquanto."
- "como os projetos só vão atualizar a baseline posteriormente ao release,
  quero já remover tudo que é irrelevante da curadoria e ajustar a baseline
  para remover implement-spec. Assim a spec do dispatch vira uma revisão do
  processo e setup de skills."
- "Remover: tdd, diagnosing-bugs, to-prompt, refactoring-analysis,
  design-patterns, game-changing-features, fixing-metadata"
- "Não remova ainda ui-ux-pro-max e zustand, vamos manter as skills
  relacionadas ao frontend e melhorar a implement-tasks para que o uso das
  skills sejam mais eficientes."
- `implement-task` stays: it is the Run agent's manual, and the Daemon names
  it in every Task's Context Bundle.

and from the skills study decisions of 2026-10-26 (AskUserQuestion, recorded
in the operator's memory): the dispatch Spec "com tudo" — take
`diagnosing-bugs` and `tdd` out of the bundles, take the dispatch off
`architectural-analysis` and `refactoring-analysis` (person-started; they
stall Runs), `writing-clearly` and `humanizer` only on request, split the
frontend-ui-quality bundle (`fixing-accessibility` stays in it;
`wcag-audit-patterns` only for a formal audit; `web-design-guidelines` only
for an explicit review), `interface-design` for product UI ×
`frontend-design` for marketing pages.

The rules those answers carry are the cycle AR rules of 2026-10-27, which
carry the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados"
(option text: "'Concedo' para os caminhos que cada Spec declarar (módulos
Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
lint, que pedem autorização específica)."), and the 2026-09-30 grants of the
skills ("considere autorizado a ajustar todas as skills se necessário") and
of Baseline sources, guides and `.roundfixrc.yml` ("Autorizar os dois").
Still not authorized: this repository's `Makefile`, `.github/workflows`,
`go.mod`, `go.sum`, and lint or formatter configuration.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point.

The governed set was measured on 2026-10-30 on the branch `docs/plan-0333`,
based on `35de7f74`, with `speccheck.GovernedPath`, through a `go test
-overlay` probe in `internal/speccheck` that wrote nothing to the repository
(`git status --porcelain internal/speccheck` empty afterwards; the probe file
never existed in the tree). It covered every file the Tasks declare, every
planning file of this Spec, its adopted Backlog Entry copies and ADR-0346.
The thirty-five files listed above are governed and changed by the Tasks;
the changed set was measured by a prototype in a scratch worktree that ran
the sanctioned commands below in order. `.agents/skills/domain-modeling/SKILL.md`
is governed but only read.

## Makefile grant

The `Makefile` is governed and task_04 must change it: retiring
`implement-spec` removes the name from the `OWNED_SKILLS` line (and the count
in the comment above it), which `make skills-sync`, `make
skills-version-check` and `make skills-sync-check` read and `make verify`
runs. On 2026-10-30 the maintainer granted that change through
AskUserQuestion. The question: "A Spec 0333 pode editar o Makefile, só a
linha OWNED_SKILLS, para tirar implement-spec?" The answer, verbatim:
"Autorizo só essa linha (Recommended)" — option text: "Permissão limitada à
linha OWNED_SKILLS do Makefile, citada no _authorization.md da 0333; o resto
do Makefile continua proibido."

The grant covers only the `OWNED_SKILLS` line and the owned-skill count in
the comment directly above it; every other `Makefile` line stays outside
this Spec's authority.

## Why each governed path is unavoidable

- The five `internal/baseline/assets/setups/` files are rewritten by the
  asset sync at the new commit (tasks 01, 02 and 04); `typescript.json`'s
  activation bundles are first set from `skill-activations.json`.
- `modules/typescript.json`, `backend.json` and `frontend.json` drop the
  seven retired skills (task_01); `typescript.json`, `frontend.json` and
  `context-workflow.json` carry the reworded triggers (task_02);
  `context-workflow.json` drops `implement-spec` (task_04);
  `skill-activations.json` carries the bundle changes (tasks 01 and 02). The
  record command also writes each changed module's version line.
- The formatter golden `skill-dispatch.md` and the `standard-typescript-monorepo`
  profile's `goldenDigest` are rewritten by `make baseline-digests`.
- `docs/agents/skill-dispatch.md` and `docs/agents/setup-context.json` are
  rendered by this repository's Managed Refresh.
- The five `.agents/skills/` files of `context7-cli`, `grilling` and
  `tui-design` are written by the sanctioned restore (task_01).
- The owned skill `roundfix` and its references `baseline.md` (tasks 01 and
  04), `reconcile.md` (task_03), `setup.md` and `spec-delivery.md` (task_04),
  and `implement-task`, `qa-gate` (task_03), `archive-spec` and `write-tasks`
  (task_04), with their mirrors, carry the rule texts; `implement-spec` and
  its mirror are deleted (task_04).

## What is not governed

- `skills-lock.json`, `skills/skills.go`, `internal/baseline/module-versions.json`,
  `internal/baseline/retired_skills.go`, `internal/baseline/assets_sync.go`,
  the catalog files, plan characterization goldens and parity files under
  `internal/baseline/testdata/`, and every test file the Tasks create or
  change.
- `internal/agent/spec_prompt.go`, `internal/daemon/task_context.go`,
  `internal/spec/spec.go`, `CONTEXT.md`, the user guide pages,
  `docs/references/behavior-surfaces.json`,
  `skills/testdata/owned-skill-versions.json`, the mirrors under
  `skills/roundfix/references/`, ADR-0346 and every planning
  file of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files and the
generated records. This declaration records the regeneration that follows the
approved edits and adds no source path.

The asset sync writes the setup snapshots from a repository fetched from the
local checkout `~/dev/skills` at
`86df8608726979064a8fe46003dfaf403ffe0087`, after `typescript.json`'s
activation bundles are set from `skill-activations.json` with `jq`; its
`--check` must then exit `0` (tasks 01, 02 and 04):

```yaml
command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
outputs:
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
```

The parity row follows through Spec 0208's jq transform with Spec 0317's
composed-snapshot filter, quoted in `_techspec.md` → API Contract 1, and the
regeneration fills its digests.

The skill restore writes this repository's three changed upstream skill
copies and their lock entries, after its plan digest is confirmed (task_01):

```yaml
command: roundfix baseline skills restore --profile go-cli-tui --skill context7-cli --skill grilling --skill tui-design --source-dir <fetched> --repo . --confirm-plan <digest> --format text
outputs:
  - .agents/skills/context7-cli/SKILL.md
  - .agents/skills/context7-cli/references/docs.md
  - .agents/skills/grilling/SKILL.md
  - .agents/skills/tui-design/SKILL.md
  - .agents/skills/tui-design/references/evidence.md
  - skills-lock.json
```

The Module Version Record step writes each changed module's version line;
suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/typescript.json
  - internal/baseline/assets/modules/backend.json
  - internal/baseline/assets/modules/frontend.json
  - internal/baseline/assets/modules/context-workflow.json
```

```yaml
command: make baseline-digests
```

`roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
run from a binary built under the Run's environment and with a fresh `HOME`,
renders this repository's Setup Manifest and agent guides, and a second
refresh reports no file change.

The owned skills are mirrored and recorded; the `implement-spec` key of the
record is deleted with `jq --indent 2 'del(.skills["implement-spec"])'`,
because the record command cannot drop a skill (task_04):

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No digest pin, golden, generated guide, snapshot, upstream skill file, lock
entry, mirror or record is hand-edited, except the one `jq` deletion above.

## Limits

- No change to the CI workflows, `go.mod`, `go.sum`, the lint, formatter or
  test-runner configuration, `.roundfixrc.yml`, the Source Baseline corpus,
  `internal/baseline/assets/retention/**`, any profile's `requiredRules` or
  modules, any template, any module other than `typescript`, `backend`,
  `frontend` and `context-workflow`, any upstream skill other than the three
  restored, any owned skill other than `implement-task`, `qa-gate`,
  `roundfix`, `archive-spec`, `write-tasks` and the deleted `implement-spec`,
  any file under `docs/history/`, any ADR other than ADR-0346, or the
  `### QA settlement` section of any skill. The `Makefile` changes only as
  the Makefile grant above allows.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads `~/dev/skills` or another repository (only the import commands of
    tasks 01, 02 and 04 do);
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.

---
status: approved
granted: 2026-10-24
action: move every upstream Setup Snapshot to the upstream commit that replaced handoff with the ai-memory skills and make the core module dispatch ai-memory-handoff and ai-memory-messaging, restoring this repository's changed upstream skills and retiring its handoff copy; replace the Secondbrain inbox clauses with ai-memory clauses of the same force and keep the Secondbrain read-only; remove the follow-up command's secondbrain-inbox source and deprecate its two configuration keys; move the owned authoring skills, the glossary and the history guide to the new channel and retire ADR-0095; regenerate the derived Baseline files and this repository's Setup Manifest through the sanctioned commands
consuming: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
paths:
  - .agents/skills/ai-memory-durable-pages/LICENSE
  - .agents/skills/ai-memory-durable-pages/SKILL.md
  - .agents/skills/ai-memory-handoff/LICENSE
  - .agents/skills/ai-memory-handoff/SKILL.md
  - .agents/skills/ai-memory-learning-maintenance/LICENSE
  - .agents/skills/ai-memory-learning-maintenance/SKILL.md
  - .agents/skills/ai-memory-messaging/LICENSE
  - .agents/skills/ai-memory-messaging/SKILL.md
  - .agents/skills/ai-memory-retrieval/LICENSE
  - .agents/skills/ai-memory-retrieval/SKILL.md
  - .agents/skills/ai-memory-routing-install/LICENSE
  - .agents/skills/ai-memory-routing-install/SKILL.md
  - .agents/skills/context7-cli/SKILL.md
  - .agents/skills/context7-cli/references/docs.md
  - .agents/skills/context7-cli/references/setup.md
  - .agents/skills/context7-cli/references/skills.md
  - .agents/skills/golang-cli/assets/examples/go.mod
  - .agents/skills/grilling/SKILL.md
  - .agents/skills/handoff/SKILL.md
  - .agents/skills/handoff/agents/openai.yaml
  - .agents/skills/knowledge-workspace/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/follow-up.md
  - .agents/skills/systematic-debugging/SKILL.md
  - .agents/skills/typesafe-ai/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - docs/agents/agent-instructions.md
  - docs/agents/docs-layout.md
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
  - docs/agents/skill-dispatch.md
  - docs/agents/spec-routing.md
  - internal/baseline/assets/coverage.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0317

On 2026-10-24 the maintainer ran `/roundfix --follow-up` and, at about 11:15
BRT, answered the cycle AP batch 5 questions through AskUserQuestion. The
answer that opened this Spec was:

- Baseline ai-memory: "Spec no fim do AP (Recomendado)" — "Plano escrito
  agora em branch refactor/, entregue depois da 0310. Antes, combino com
  marcioaltoe-skills os nomes e as versões; o Secondbrain continua sendo a
  base de conhecimento, e o ai-memory vira o canal de handoff e de
  mensagens. Respondo ao Jarvis com o plano."

At about 13:05 BRT the maintainer wrote, in his own words: "Lembrando que
solicitei em algum momento para matar o inbox de secondbrain e usar somente
as mensagens do ai-memory". He then answered the cycle AP batch 6 questions
through AskUserQuestion at about 13:10 BRT, refined at about 13:20 BRT:

- order: "Antes da 0307 (Recomendado)" — "A 0317 entra logo após a 0316. A
  0307 é emendada para usar só o ai-memory (sem triagem de inbox), e a 0317
  também tira a fonte secondbrain-inbox do comando follow-up. O setup Python
  vai uma versão para frente (v0.127).";
- pending entries, the maintainer's own answer: "Migrar as mensages ainda
  não triadas para o ai-memory e remover as demais e triadas e arquivadas.";
  refined: release notes "Descartar todas"; other entries "Migrar as 63,
  descartar tax-poc (Recomendado)"; removal "Eu migro, o secondbrain remove
  (Recomendado)";
- adopter notes: "Via ai-memory já (Recomendado)".

The request this Spec answers came from the maintainer through ai-memory
(message `01a16e17`, 2026-10-24 10:26 BRT); the adopted Backlog Entry
records it. Upstream answered the coordination with pull request #131 of
`marcioaltoe/skills`, commit `65d5f365cace27c4b468db1723213723b10c870b`.
The delivery order is 0312, 0313, 0314, 0315, 0316, 0317, 0307, 0309 and
0310. The migration of pending Inbox Entries and the deletion of `inbox/**`
belong to the operator and to the Secondbrain's own project; no Task,
Verification command or QA row of this Spec does either.

The cycle AP rules those answers carry were approved on 2026-10-22
(questions of 15:25 BRT, answered by about 18:30 BRT): "the same as cycle
AO", the Governed Paths each Spec declares, and still not this repository's
Makefile, `.github/workflows`, `go.mod`, `go.sum`, or lint and formatter
configuration. Those rules carry the cycle AJ grant of 2026-10-19, at about
16:20 BRT, "Mexer em Governed Paths declarados" (option text: "'Concedo' para
os caminhos que cada Spec declarar (módulos Baseline, registros, skills,
testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push", and the 2026-09-30 grants
of the skills ("considere autorizado a ajustar todas as skills se
necessário") and of Baseline sources, guides and `.roundfixrc.yml`
("Autorizar os dois").

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on 2026-10-24 on the branch
`refactor/plan-0317-ai-memory-baseline`, based on `b5e45c63`, with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository (`git status
--porcelain` identical before and after; the probe file never existed in the
tree). It covered every file the Tasks declare, every planning file of this
Spec, its adopted Backlog Entry copy and ADR-0330. The fifty-three files
listed above are governed and changed by the Tasks; the changed set was
measured by a prototype in a scratch worktree that ran the sanctioned
commands below in order.

## Why each governed path is unavoidable

- The four `internal/baseline/assets/setups/` files are rewritten by the
  asset sync at the new commit (task_01).
- `internal/baseline/assets/modules/core.json` swaps its skills (task_01) and
  gains the continuity rule (task_02); `secondbrain.json`,
  `context-workflow.json`, `spec-workflow.json` and `coverage.json` carry the
  guidance and coverage changes (task_02). The record command also writes
  each changed module's version line.
- The formatter goldens under
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/`
  and that profile's `goldenDigest` are rewritten by `make baseline-digests`
  (task_01, task_02).
- The six `docs/agents/` files are rendered by this repository's Managed
  Refresh (task_01, task_02).
- The twenty-three `.agents/skills/` files of upstream skills are written or
  removed by the sanctioned restore and reconcile, and the `handoff` tree is
  deleted after the reconcile removed its lock entry (task_01).
- `.agents/skills/roundfix/SKILL.md`, its reference `follow-up.md` and the
  mirror `skills/roundfix/SKILL.md` describe the follow-up command (task_03);
  the three authoring skills and their mirrors read ai-memory messages
  (task_04).

## What is not governed

- `skills-lock.json`, `skills/recommended.txt`,
  `internal/baseline/module-versions.json`, the catalog files, plan
  characterization goldens and the parity files under
  `internal/baseline/testdata/`, and every test file the Tasks create or
  change.
- `internal/config/config.go`, `internal/cli/follow_up_inventory.go`,
  `CONTEXT.md`, the user guide pages, `docs/references/skill-coverage.json`,
  `docs/references/behavior-surfaces.json`,
  `skills/testdata/owned-skill-versions.json`,
  `skills/roundfix/references/follow-up.md`, ADR-0330, the ADR-0095 file the
  Tasks delete, and every planning file of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files and the
generated records. This declaration records the regeneration that follows the
approved edits and adds no source path.

The asset sync writes the setup snapshots from a repository fetched from the
local checkout `~/dev/skills` at
`65d5f365cace27c4b468db1723213723b10c870b`; its `--check` must then exit
`0`:

```yaml
command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
outputs:
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
```

The parity row follows through Spec 0208's jq transform, quoted with its
composed-snapshot filter in `_techspec.md` → API Contract 1, and the
regeneration fills its digests.

The skill restore writes this repository's changed upstream skill copies,
the six ai-memory skills and their lock entries, after its plan digest is
confirmed:

```yaml
command: roundfix baseline skills restore --profile go-cli-tui --skill ai-memory-durable-pages --skill ai-memory-handoff --skill ai-memory-learning-maintenance --skill ai-memory-messaging --skill ai-memory-retrieval --skill ai-memory-routing-install --skill context7-cli --skill golang-cli --skill grilling --skill knowledge-workspace --skill systematic-debugging --skill typesafe-ai --source-dir <fetched> --repo . --confirm-plan <digest> --format text
outputs:
  - .agents/skills/ai-memory-durable-pages/LICENSE
  - .agents/skills/ai-memory-durable-pages/SKILL.md
  - .agents/skills/ai-memory-handoff/LICENSE
  - .agents/skills/ai-memory-handoff/SKILL.md
  - .agents/skills/ai-memory-learning-maintenance/LICENSE
  - .agents/skills/ai-memory-learning-maintenance/SKILL.md
  - .agents/skills/ai-memory-messaging/LICENSE
  - .agents/skills/ai-memory-messaging/SKILL.md
  - .agents/skills/ai-memory-retrieval/LICENSE
  - .agents/skills/ai-memory-retrieval/SKILL.md
  - .agents/skills/ai-memory-routing-install/LICENSE
  - .agents/skills/ai-memory-routing-install/SKILL.md
  - .agents/skills/context7-cli/SKILL.md
  - .agents/skills/context7-cli/references/docs.md
  - .agents/skills/context7-cli/references/setup.md
  - .agents/skills/context7-cli/references/skills.md
  - .agents/skills/golang-cli/assets/examples/go.mod
  - .agents/skills/grilling/SKILL.md
  - .agents/skills/knowledge-workspace/SKILL.md
  - .agents/skills/systematic-debugging/SKILL.md
  - .agents/skills/typesafe-ai/SKILL.md
  - skills-lock.json
```

The lock reconcile removes the `handoff` lock entry, after its plan digest is
confirmed; the retained `.agents/skills/handoff/` tree is then deleted with
`git rm -r`, and `skills/recommended.txt` is rewritten from the lock's names:

```yaml
command: roundfix baseline skills reconcile --profile go-cli-tui --source marcioaltoe/skills --revision 65d5f365cace27c4b468db1723213723b10c870b --source-dir <fetched> --repo . --confirm-plan <digest> --format text
outputs:
  - skills-lock.json
```

The Module Version Record step writes each changed module's version line;
suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/spec-workflow.json
```

```yaml
command: make baseline-digests
```

`roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
run from a binary built under the Run's environment and with a fresh `HOME`,
renders this repository's Setup Manifest and agent guides, and a second
refresh reports no file change.

The owned skills are mirrored and recorded:

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No digest pin, golden, generated guide, snapshot, upstream skill file, lock
entry, mirror or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, the Source
  Baseline corpus, manifest, accounting or index,
  `internal/baseline/assets/retention/**`, any profile's `requiredRules` or
  modules, any template, any module other than `core`, `secondbrain`,
  `context-workflow` and `spec-workflow`, any upstream skill outside the
  twelve restored and `handoff`, any owned skill other than `roundfix`,
  `write-idea`, `write-prd` and `write-techspec`, any file under
  `docs/history/`, any ADR other than the deletion of ADR-0095, or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads or writes the Secondbrain or calls ai-memory;
  - reads `~/dev/skills` (only task_01's import does) or another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.

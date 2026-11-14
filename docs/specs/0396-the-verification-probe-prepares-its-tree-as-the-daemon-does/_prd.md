---
spec: 0396-the-verification-probe-prepares-its-tree-as-the-daemon-does
status: active
created: 2026-11-14
surfaces: [backend, cli, docs]
type: fix
---

# The Verification probe prepares its tree as the Daemon does

Reported by indexus on 2026-11-14: `roundfix spec check --strict --run-verification` (0.143.1) called a `bunx vitest run ...` command "honest — exited non-zero before work", and minutes later the Daemon refused the same command as vacuous because it exited zero against the unchanged tree of the Run. The check probes each command in a bare disposable checkout of `HEAD`, while a Run or Task Worktree receives the `worktree.copy` files and the Worktree Bootstrap before the Daemon probes it, so a command that needs installed dependencies fails in the check for a reason unrelated to the work. After this Spec the check prepares its checkout as the Daemon prepares a Task Worktree, once, right before the first command it probes, so its verdict matches the Daemon's; a failed bootstrap ends the check without calling any command honest, and a repository without `worktree.copy` and `worktree.bootstrap` is checked exactly as before.

## Sources

One ai-memory message, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-14-indexus-the-verification-probe-called-honest-a-command-the-daemon-refused-as-vacuous.md`](references/2026-11-14-indexus-the-verification-probe-called-honest-a-command-the-daemon-refused-as-vacuous.md),
   ai-memory message `01a1d9be-ad8c-7f33-a039-0d8d4e60fef7`, copied
   verbatim; reported by indexus on 2026-11-14. Indexus declared no blocked
   work (it prefixed its commands with `test -f <new file> &&`) and asked for
   no reply.

## Prerequisites

None in `requires`. Authored against `4f8a5b7d` (v0.143.4). Specs 0392 and
0393 are planned on the default branch and deliver before this one; no
Verification of this Spec depends on their work. The three Specs share only
these files, each in its own part:

- `internal/cli/spec_check.go`: 0393 task_01 changes the `spec` usage text;
  task_01 of this Spec changes the Verification probe, its report fields and
  its rendering, and never the usage text.
- `docs/user-guide/commands/spec.md` and both copies of the Roundfix skill's
  `references/spec.md`: 0393 task_01 documents `spec judge --stage idea`;
  task_01 of this Spec adds one paragraph after the `--run-verification`
  paragraph.
- The Roundfix owned skill's version fields in `.agents/skills/roundfix/SKILL.md`
  and `skills/roundfix/SKILL.md`, and `skills/testdata/owned-skill-versions.json`:
  0392 and 0393 raise them; task_01 raises the next version with the record
  command on the tree it runs on and never writes a number.
- `docs/references/behavior-surfaces.json`: 0392 and 0393 re-record it;
  task_01 re-records it after its own last edit of a surface source.
- `CONTEXT.md`: 0392 and 0393 add their own terms; task_01 changes only the
  **Worktree Bootstrap** entry.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added; the bootstrap the check runs is the command the repository
  already configures for its Run Worktrees. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0406 (new) records the decision:
  "The authoring Verification probe prepares its checkout as the Daemon
  prepares a Task Worktree". It amends ADR-0148, whose authoring caller
  supplied a bare working directory, and extends ADR-0034, whose Worktree
  Bootstrap now also runs in the check's checkout; task_01 writes both
  notices. ADR-0384 and ADR-0385 do not apply: the bootstrap stays ungated
  and receives no Verification Core Budget, as both already state. ADR-0056,
  which cites ADR-0034, and ADR-0315, which cites ADR-0148, do not apply
  either (reasons in the TechSpec's ADR Alternatives). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BF batch 2 answers of 2026-11-14 quoted in the authorization record;
  the Governed Paths are the Roundfix skill files task_01 changes. No
  Makefile, workflow, Go module, lint or formatter change, and no Baseline
  module, template or profile change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0396-the-verification-probe-prepares-its-tree-as-the-daemon-does/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `skills/roundfix/SKILL.md`.

## Goals

- A Verification command that passes on a tree prepared by `worktree.copy`
  and the Worktree Bootstrap is reported `vacuous` by
  `roundfix spec check --run-verification`, as the Daemon refuses it.
- A Worktree Bootstrap that fails in the check probes no command, names the
  failure and exits `1`.
- The bootstrap runs at most once per check and only when a command will be
  probed.
- A repository without `worktree.copy` and `worktree.bootstrap` gets the same
  report, byte for byte, and the same exit code.
- The `spec` guide, the Roundfix skill's `spec` reference, the glossary entry
  **Worktree Bootstrap** and the two amended ADRs state the rule.

## Core Features

1. Before its first probed command, the check provisions its disposable
   checkout with the repository's `worktree.copy` list and Worktree
   Bootstrap through the worktree package's provisioning function, the one
   that provisions an item worktree.
2. A Spec report whose commands ran after the bootstrap prints
   `Verification tree bootstrap: "<command>"` and its JSON carries
   `verification.bootstrap`.
3. A failed bootstrap stops the check with a three-part message on standard
   error: the failure, the statement that no command was probed, and the
   bootstrap's last output.
4. The guide, both copies of the Roundfix `spec` reference, the glossary and
   the notices in ADR-0148 and ADR-0034 carry the rule.

## User Experience

Indexus's check printed, for a command its bootstrap made pass:

```text
- task_01: honest — "bunx vitest run test/domain/fiscal/orientation-act.test.ts ..." (exited non-zero before work)
```

With this Spec the same check, in a repository whose Worktree Bootstrap
installs the dependencies, prints:

```text
Verification tree: HEAD
Verification tree bootstrap: "bun install"
- task_01: vacuous — "bunx vitest run test/domain/fiscal/orientation-act.test.ts ..." (exited zero before work)
```

and exits `1`, as the Daemon refuses that Task before its Agent runs.

## Non-Goals / Out of Scope

- The Run Marker and the Verification Core Budget the Daemon adds to its own
  Verification environment: they are Run bookkeeping (ADR-0148), the check
  runs outside any Run, and neither changes whether a command exits zero.
- Gating the bootstrap with the Machine Verification Capacity (ADR-0384).
- Probing completed Tasks or the QA gate, which stay `not probed`.
- A new flag to skip the preparation: the verdict must be the Daemon's.
- Changing the Baseline's Spec routing or the `write-tasks` skill: both
  already tell authors to run the check before a Run, and the check now
  answers as the Daemon does.
- Messages to indexus: the operator sends them.

## Success Metrics

1. Success Metric: in a temporary repository whose Project Config copies an
   untracked `local.env` and whose Worktree Bootstrap needs that file and
   creates `.deps-installed`, `spec check --run-verification` reports
   `test -f .deps-installed` and `test -f local.env` as `vacuous` and a
   command for a file nobody creates as `honest`, prints the bootstrap line
   and exits `1`; the JSON report carries the bootstrap. At `4f8a5b7d` all
   three are `honest` and the check exits `0`.
2. Success Metric: with a bootstrap that exits `7` after printing one line,
   the check prints nothing on standard output, prints the three-part
   message on standard error, runs no command and exits `1`, in text and in
   JSON.
3. Success Metric: a bootstrap that appends a line to a counter file runs
   zero times for a Spec with only a completed Task and a QA gate, and
   exactly once for two Specs of two pending Tasks each.
4. Success Metric: without `worktree.copy` and `worktree.bootstrap`, the
   report from `Verification tree: HEAD` on, its first two lines, its JSON
   `verification` object and its exit code equal those at `4f8a5b7d`.
5. Success Metric: the `spec` guide, both copies of the Roundfix `spec`
   reference, the **Worktree Bootstrap** entry and the notices in ADR-0148
   and ADR-0034 carry the new text; the skills and documentation contracts
   pass.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted indexus message, a measurement on a repository this Spec did
  not build: "no Run worktree, o mesmo comando sai 0. O vitest 5.1 ignora os
  três arquivos inexistentes e roda só os dois existentes (\"Test Files 2
  passed (2)\"). Então o Daemon estava certo e o probe do spec check deu
  falso \"honest\"."
- git-worktree documentation (https://git-scm.com/docs/git-worktree, read
  through Exa on 2026-11-14): "The new worktree is linked to the current
  repository, sharing everything except per-worktree files such as `HEAD`,
  `index`, etc." A detached worktree holds the files of the commit it checks
  out, not the untracked or ignored files of the main worktree.
- The installed `git` (measured on the authoring host), which this Spec did
  not build: in a temporary repository whose main worktree holds an ignored
  `node_modules/marker`, `git worktree add --detach <dir> HEAD` produces a
  checkout without `node_modules/`.

## Unreachable Acceptance

None. Every criterion is proved with temporary repositories, the built
binary and the installed `git`; no criterion needs the network, an open Pull
Request or another repository.

## Validation

The decisions were checked on 2026-11-14 against these sources.

- Characterization at `4f8a5b7d` with the installed `roundfix` 0.143.4 on a
  temporary repository whose Project Config sets
  `worktree.bootstrap: "touch .deps-installed"`: `spec check 0001-probe
  --run-verification` reported `- task_01: honest — "test -f .deps-installed"
  (exited non-zero before work)`. The same four tests this Spec adds, run on
  the starting tree, failed on that verdict, on the missing bootstrap line,
  on the missing bootstrap failure and on the second Spec fixture, and the
  test of the unchanged report passed.
- The Daemon path, read at `4f8a5b7d`: `worktree.Create` and
  `createTaskWithOptions` call `copyProvisionedFiles` and then
  `runBootstrap` before the Task runs; `Engine.verifyTaskPreWork` probes in
  `plan.WorkDir` through `probeCommandsForRun` and `ExecVerifier`, the runner
  `spec check` already uses; `ProvisionItem` performs the same two steps for
  a Delivery Queue item worktree.
- A prototype of task_01 on a scratch copy at `4f8a5b7d`: the four new tests
  passed, the `SpecCheck` tests of `./internal/cli` and the provisioning tests
  of `./internal/worktree` passed, both packages built for Windows, and each
  sabotage of the TechSpec made its named test fail.
- Secondbrain, read only: `wiki/index.md`, then `qmd query --all --files
  --min-score 0.3` for the probe, the disposable checkout and the Worktree
  Bootstrap. It returned this repository's own mirrored `CONTEXT.md`,
  `CHANGELOG.md` and guides, nothing that contradicts the decision; the
  adopters' mirrors show that Fluxus configures `worktree.bootstrap: "bun
  install"`
  (`projects/fluxus/mirror/docs/references/2026-09-09-autorizacao-de-tooling-para-o-project-config-do-roundfix.md`),
  the class of repository this Spec serves and the cost it pays.
- Exa: the git-worktree documentation above, and secondary articles that
  say the same of a new worktree (no `node_modules/`, no generated
  configuration), which this Spec does not cite as evidence.

## Glossary

- changes: **Worktree Bootstrap**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BF batch 2 brief for Spec 0396 of 2026-11-14, the maintainer's answers below and the adopted source.

The maintainer's answers (AskUserQuestion, 2026-11-14 ~11:35 BRT), verbatim:

- Question "Mensagem do Indexus: o `roundfix spec check --run-verification`
  chamou de \"honest\" um comando que o Daemon depois recusou como vacuous
  (provável causa: a checagem roda sem as dependências instaladas). O
  Indexus contornou e não está bloqueado. O que fazemos?" → "Spec fix 0396
  (Recommended)" — option text: "Reproduzo a causa, escrevo a TechSpec no
  Opus (sem Fable) e entrego depois da 0393 como patch v0.145.1, com o
  Quaestor ocioso."
- Question "Se a Spec 0396 do Indexus for aberta, como a checagem `spec
  check --run-verification` deve se corrigir? (O Jev ficou indeciso
  exatamente nesta escolha.)" → "Igualar ao Daemon (Recommended)" — option
  text: "A checagem prepara as dependências como o Daemon prepara (ex.: bun
  install) antes de rodar o comando. Mais lenta, porém o veredito bate com o
  da entrega."
- Question "Para a Spec 0396 do Indexus (uma correção com o defeito medido),
  posso pular a etapa de Ideia e ir direto à TechSpec?" → "Pular a Ideia
  (Recommended)" — option text: "Correção com evidência medida; a forma já
  foi decidida na pergunta anterior. TechSpec direto, no Opus."
- Question "Autorizações para a Spec 0396 (Indexus) e para a 0397 (limite do
  Claude): valem as mesmas do ciclo BE?" → "Mesmas do BE (Recommended)" —
  option text: "Push, PR, merge, release depois de cada Spec, até 3
  correções; Makefile/CI/go.mod/lint só se a correção exigir. Com o Quaestor
  ocioso, 1 vaga."

Taken from those answers, the source and the characterization, with no
further interview:

- `type: fix`: a verdict that contradicted the Daemon's becomes the
  Daemon's; the Release Ceiling is `patch`.
- The check reuses the worktree package's provisioning, copy then bootstrap,
  with the repository's own `worktree.copy`, `worktree.bootstrap` and
  `worktree.bootstrap_timeout` (ADR-0406).
- The preparation runs once per invocation, right before the first probed
  command; one checkout still serves every Spec of the invocation (ADR-0406).
- A failed bootstrap probes nothing and exits `1` with the failure, the
  statement that nothing was probed and the bootstrap's last output on
  standard error; standard output stays empty, as for any failure to prepare
  the Verification tree today (ADR-0406).
- The bootstrap's output is not streamed: a successful bootstrap prints only
  the report line, and a failed one prints the last output Roundfix keeps
  (at most 4096 bytes).
- Adopters receive the change through the binary and the Roundfix skill; no
  Baseline module, template or profile text changes.
- One implementation Task holds the code, its tests and the text that
  describes it, and the QA gate follows: the documentation alone would not
  earn its own verification cycle, and the code Task must re-record the
  Behavior Surface Record anyway.

## Open Questions

None.

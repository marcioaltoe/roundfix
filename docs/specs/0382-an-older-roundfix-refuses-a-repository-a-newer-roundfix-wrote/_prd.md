---
spec: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
status: active
created: 2026-11-11
surfaces: [backend, cli, docs]
type: feat
---

# An older Roundfix refuses a repository a newer Roundfix wrote

Reported by conexus on 2026-11-11: until 10/11 a developer ran Roundfix 0.96.0 against a repository configured for 0.137.7, and two Specs ran outside Roundfix because the configuration was incompatible. Measured at `013af3b4` (v0.139.0), a repository records no Roundfix version, Roundfix 0.96.0 fails every configuration-loading command with `worktree.shared_caches is not a supported config key`, and Roundfix 0.137.7 proposes, as an ordinary `baseline update` plan, to rewrite the managed guides of a v0.139.0 tree back to its older text. After this Spec the Setup Manifest records the **Recorded Roundfix Version**, the release that last applied a Baseline change, and a binary older than it refuses early, naming both versions and `roundfix upgrade`, while the commands a user needs to recover keep working.

## Sources

One cross-project message, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md`](references/2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md),
   ai-memory message `01a1c84e-8dfb-7483-966c-061a3c293853` from conexus,
   2026-11-11 01:52Z, `type: fix`, reported by conexus. This Spec owns only
   its version-check part, the sender's suggestion that a minimum-version
   check should refuse clearly instead of failing later; its three
   unreproduced symptoms go to a Backlog Entry the operator writes, waiting
   for reproduction steps.

At the authoring base `docs/findings/` and `docs/backlog/` hold no entry about
a binary version or a repository version.

## Prerequisites

None in `requires`. Authored against `013af3b4` (v0.139.0). Specs 0377, 0376,
0378, 0379, 0380 and 0381 deliver before this one, in that order, and may
change files this Spec also changes, each in its own section:

- `internal/cli/doctor.go`, `internal/cli/health.go` and
  `internal/cli/update_check.go` (0376 adds a `user-skills:` Doctor line after
  `skills:`, a `HealthCheckUserSkills` name and the Update Check's
  `userSkills`; this Spec adds a `baseline-version:` first line, a
  `HealthCheckBaselineVersion` name and changes the Update Check's
  configuration load only).
- `internal/cli/cli.go` (0376 changes `runSkillsCommand` and
  `commandDependencies`; this Spec adds two configuration loaders and changes
  `loadRunReaderCommandConfig` and `runStopCommand`).
- `docs/user-guide/commands/doctor.md` and `docs/user-guide/commands/update.md`
  (0376 adds a `user-skills:` block and a `User skills:` block; this Spec adds
  its own paragraphs at the anchors tasks 01 and 02 name).
- `.agents/skills/roundfix/references/setup.md` and its mirror (0376 adds a
  paragraph after the skills-restore paragraph; this Spec one after the Update
  Check paragraph).
- The Roundfix owned skill's version fields in `.agents/skills/roundfix/SKILL.md`
  and `skills/roundfix/SKILL.md` and `skills/testdata/owned-skill-versions.json`
  (0376, 0380 and 0381 raise them too). task_03 raises the version with the
  record command from the tree it runs on and never writes a number.
- `CONTEXT.md` (each Spec its own terms) and
  `docs/references/behavior-surfaces.json` (each docs Task re-records it after
  its own last guide edit).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the Setup
  Manifest keeps `setup-context-driven/manifest/0.0.1` and gains one optional
  field, `roundfixVersion`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added; the record is a local file. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0392 (new) records the decision:
  "An older Roundfix refuses a repository whose Baseline a newer Roundfix
  wrote". ADR-0103 applies unchanged: "a second run against an unchanged
  catalog reports the repository current with no proposed change"; the
  Recorded Roundfix Version rides only a plan that changes a file. ADR-0287
  applies unchanged: the documentation contract that builds the Managed
  Refresh plan at every settlement passes no writer version, so a release
  never makes this repository's Setup Manifest stale. ADR-0371 applies unchanged: the Update Check still changes
  nothing and every report exits `0`. ADR-0173, which cites ADR-0103, does
  not apply: it decides how a Baseline Plan reports the citations its History
  Relocations break, and this Spec changes no relocation or citation. Source:
  `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BC answers of 2026-11-11 quoted in the authorization record; the
  Governed Paths are the Roundfix skill files task_03 changes and the two
  CLI test files whose parity tests task_02 updates. No Makefile,
  workflow, Go module, lint or formatter change, and no Baseline module,
  template or profile change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/setup.md`, `internal/cli/baseline_human_test.go`, `internal/cli/baseline_plan_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- A Roundfix binary older than the release that last applied a Baseline change
  to a repository stops at its first step there, with one message that names
  both versions and `roundfix upgrade`, instead of an unrelated configuration
  error or a plan that rewrites the Baseline with older text.
- The user can still diagnose and recover: `roundfix --version`, `roundfix
  upgrade`, `roundfix doctor`, the Update Check, the Run readers and `roundfix
  stop` keep working.
- No adopter gets a commit only because a new release came out.

## User Stories

1. As a developer whose binary is older than the one a teammate used, I want
   my first command in the repository to tell me which Roundfix wrote it and
   to run `roundfix upgrade`, so that I do not work around Roundfix.
2. As the maintainer, I want an older binary's `baseline update` to refuse
   instead of planning a downgrade of the managed guides, so that an
   `--update` session never applies one.
3. As an Agent session, I want the Update Check and `roundfix doctor` to keep
   reporting in such a repository, so that the session learns why and what to
   run.

## Core Features

1. The Setup Manifest's optional `roundfixVersion`, written by a Baseline plan
   the CLI builds when the plan changes a file and the running release is
   above the recorded one; a plan without a file change keeps it.
2. Configuration loading reads it before Project Config and refuses a newer
   value with the message ADR-0392 states; every command that loads
   configuration, `roundfix init --scope project` and the six Baseline
   commands refuse.
3. `roundfix doctor` prints `baseline-version: failed` first and runs its
   other checks without Project Config; the Update Check, the Run readers and
   `roundfix stop` warn on stderr and continue without Project Config.
4. A malformed value is refused; an absent value changes nothing; a
   development build is never refused.
5. The Baseline, Doctor, Update Check and upgrade guides, the Roundfix
   `baseline` and `setup` references, the glossary and ADR-0392 describe it.

## User Experience

With a binary built at `013af3b4` from this Spec's prototype (release 0.139.0)
in a copy of this repository whose Setup Manifest a 0.140.0 build of the same
prototype had just written:

```text
$ roundfix baseline update --repo . --format text
roundfix: baseline update failed: Setup Manifest "<repo>/docs/agents/setup-context.json" records Roundfix 0.140.0, newer than this binary's 0.139.0; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'
Baseline update: failed
Category: version
...
Next action: run roundfix upgrade, then rerun roundfix baseline update
$ roundfix runs list
roundfix: warning: Setup Manifest "<repo>/docs/agents/setup-context.json" records Roundfix 0.140.0, newer than this binary's 0.139.0; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'; Project Config not read
No Runs found.
```

## Non-Goals / Out of Scope

- The three unreproduced symptoms of the Conexus message: the operator's
  Backlog Entry owns them.
- Teaching binaries older than this release the refusal: they cannot read the
  record. The record protects a repository from its first Baseline change
  applied by this release on; an older binary then sees an unknown Setup
  Manifest field, and its Baseline commands report the manifest as
  incompatible instead of planning a rewrite.
- Explaining an unknown Project Config key as one a newer release documents:
  a key written by hand keeps today's "is not a supported config key"
  refusal (ADR-0392, Considered Options).
- A version compatibility range, a maximum version, or a refusal of a newer
  binary on an older repository.
- Recording a version on a plan that changes no file, or on every release.
- Reading the record in a Run's Daemon after the command that started it:
  the command already refused.
- A message to adopters: they receive the record at their next Baseline
  change through the Update Check and `baseline update`.

## Success Metrics

1. Success Metric: with a Setup Manifest recording a release above the running
   one and a Project Config holding an unknown top-level key, `roundfix spec
   check`, `roundfix follow-up` and `roundfix init --scope project --force`
   exit `2` with the refusal and never parse or rewrite Project Config. At
   `013af3b4` they report the unknown key or rewrite the file.
2. Success Metric: in the same repository `roundfix baseline`, `baseline
   update`, `baseline plan`, `baseline apply`, `baseline skills restore` and
   `baseline skills reconcile` exit `2` with the refusal and change no file.
3. Success Metric: `roundfix runs list` exits `0` with the warning line,
   `roundfix doctor` prints `baseline-version: failed (...; next: roundfix
   upgrade)` as its first line and exits `1`, and the Update Check exits `0`
   with the warning and `Baseline: failed: <refusal>`, its JSON `next` `none`
   for a development build.
4. Success Metric: a Baseline plan with a file change records the writer's
   release; a converged repository replanned by a newer, empty, development
   or pre-release writer proposes no change and keeps the value; no writer
   lowers it; a plan built without a writer version records nothing, so every
   plan golden and the parity corpus keep their bytes.
5. Success Metric: `roundfix baseline update --yes` on a stale repository
   records the running release, and a second run converges with zero file
   changes.
6. Success Metric: a value that is not a `MAJOR.MINOR.PATCH` release is
   refused by configuration loading and makes the Setup Manifest invalid for
   the Baseline commands; a development build is not refused.
7. Success Metric: the Baseline, Doctor, Update Check and upgrade guides and
   both copies of the Roundfix `baseline` and `setup` references document the
   record and the refusal; `CONTEXT.md` defines **Recorded Roundfix Version**
   and revises **Setup Manifest**; the skills, documentation and glossary
   contracts pass.

## Acceptance evidence

The outside-evidence row rests on published documentation this Spec did not
write, read through Exa on 2026-11-11.

- Go, "Go Toolchains" (https://go.dev/doc/toolchain), section "Module and
  workspace configuration": "The Go toolchain refuses to load a module or
  workspace that declares a minimum required Go version greater than the
  toolchain's own version." and "For example, Go 1.21.2 will refuse to load a
  module or workspace with a `go 1.21.3` or `go 1.22` line." This is the
  refusal this Spec adopts: an older tool stops on a repository record that a
  newer one requires.
- HashiCorp, "Manage Terraform versions"
  (https://developer.hashicorp.com/terraform/tutorials/configuration-language/versions):
  "Terraform will update the `terraform_version` whenever you apply a change to
  your configuration using a newer Terraform version." This is the write rule
  this Spec adopts: the record moves when a newer tool applies a change, not
  when it only reads.

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-11 against these sources.

- Characterization at `013af3b4`, with binaries built from the tags `v0.96.0`
  and `v0.137.7` run against a copy of the tree in temporary homes:
  - nothing in the repository records a Roundfix version: the Setup Manifest
    carries only `0.0.1` schema and generator versions, `.roundfixrc.yml`
    none, and `skills-lock.json` none;
  - `config.Load` decodes Project Config with `KnownFields(true)`: an unknown
    top-level key fails with `field <key> not found in type
    config.configOverlay`, a known section's unknown key with `<key> is not a
    supported config key`, and Roundfix 0.96.0 fails `runs list` and `spec
    check` with `worktree.shared_caches is not a supported config key`;
  - Roundfix 0.137.7's `baseline update` reports `plan ready` with two file
    changes and, with `--yes`, rewrites `docs/agents/agent-instructions.md`
    back to its text without the message rule of Spec 0373 and rewrites the
    Setup Manifest;
  - the current binary's `baseline update` reports a Setup Manifest with an
    unknown field as `incompatible` ("manifest does not match the supported
    Setup Manifest schema"), because `parseManagedSetupManifest` decodes it
    strictly, which is what a binary before this Spec does with the new field;
  - Doctor's `resolveExternalSkillRequirement` reads the Setup Manifest's
    modules leniently and fails on a module the binary does not know,
    another late error the refusal now precedes;
  - `store.SchemaVersionError` already refuses a newer Run Database with "a
    newer roundfix wrote it, so use that binary or run 'roundfix upgrade'",
    the wording this Spec mirrors.
- A prototype of tasks 01 to 03 on a scratch worktree at `013af3b4`: the
  thirteen new tests passed; `./internal/baseline`, `./internal/config` and
  `TestTheSetupManifestIsCurrent` passed unchanged, so no plan golden, parity
  fixture or this repository's Setup Manifest moved; the whole
  `./internal/cli` passed once the three parity tests that compare a CLI plan
  with a direct `BuildPlan` received the same writer version; the skills and
  documentation contracts, the glossary phrase contract included, passed on
  the prescribed text; a
  0.140.0 build recorded `"roundfixVersion": "0.140.0"` only after a managed
  guide was made stale, converged on the next run, and the 0.139.0 build then
  refused; each sabotage of the TechSpec made its named test fail.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for a Roundfix
  binary's compatibility with a repository's Baseline. The top results were
  the mirror of the decision that restarted every Roundfix-owned version at
  `0.0.1` (why the schema and generator versions cannot tell releases apart) and
  Fluxus's `docs/references/2026-11-07-autorizacao-do-baseline-do-roundfix-0.137.1.md`,
  where an adopter names its Baseline by the Roundfix release that applied it
  and checks that "o plano não rebaixa nenhuma skill"; no record holds a
  version check, so the decision rests on the maintainer's answer and the
  characterization.
- Exa: the two published sources above, plus HashiCorp's state error "state
  snapshot was created by Terraform v1.8.0, which is newer than current
  v1.6.3", the refusal form of a tool reading a record a newer one wrote.

## Glossary

- adds: **Recorded Roundfix Version**
- changes: **Setup Manifest**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BC brief for Spec 0382 of 2026-11-11, the maintainer's answers below and the adopted message.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-11 ~14:45–14:50 BRT), verbatim:

- Question "Relato do Conexus 01a1c84e (status disse 'nenhuma run ativa' com
  Run rodando; comando parado esperando autorização sem retorno ao agente;
  tela preta; Diego usava binário 0.96 em repo configurado para 0.137.7).
  Nenhum foi reproduzido. Qual destino?" → "Backlog + Spec de versão
  (Recommended)" — option text: "Spec F agora: binário mais velho que o
  Baseline do repositório recusa com mensagem clara. Os três sintomas viram
  Backlog Entry esperando passos de reprodução de você ou do Diego."
- Question "Spec F (binário mais velho que o Baseline do repositório recusa
  com mensagem clara e nomeia 'roundfix upgrade'): autorização e posição?" →
  "Mesmas do BB, após E (Recommended)" — option text: "Push, merge, release,
  teto minor; entregue depois de D e E, fechando o ciclo."

Taken from those answers, the message and the characterization, with no
further interview:

- `type: feat`: the Spec adds a recorded field and a refusal, not only a
  clearer message for a failure that already happens; the Release Ceiling is
  `minor`.
- The record lives in the Setup Manifest, the file every Baseline change
  rewrites and every Baseline-changing Spec already regenerates; no new file.
- The record moves only with a plan that changes a file, never lowers, and
  ignores development and pre-release builds; a plan built without a writer
  version keeps the recorded value, so no golden moves.
- One check in configuration loading, before Project Config, covers every
  command that loads configuration; the Baseline commands and `init --scope
  project` call the same check for the repository they write.
- Doctor, the Update Check, the Run readers and `roundfix stop` warn and read
  no Project Config; `--version`, `help` and `upgrade` read none before they
  act. A malformed value is refused; a development build never is.
- The message mirrors the Run Database's newer-schema refusal.
- ADR-0392 records the decision. Two code Tasks that each document their own
  surface (task_01 also adds the glossary entries), one Task for the Roundfix
  skill's references, and the QA gate.

## Open Questions

None.

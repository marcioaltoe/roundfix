---
spec: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
status: active
created: 2026-11-12
surfaces: [cli, docs]
type: fix
---

# Roundfix reads a version 2 skills lock and keeps its version

On 2026-11-12 Vortex and Pantheon reported that the `skills` CLI 1.8 rewrote their `skills-lock.json` as version 2, with a new folder hash in every `computedHash`, and that Roundfix 0.141.1 then refused the lock in `roundfix doctor`, `roundfix baseline update`, the release plan's `skills` check and `roundfix skills install --target project`, while Doctor's next action named the very command that writes version 2 again. After this Spec Roundfix reads a lock of version 1 or 2, writes and compares every `computedHash` in the hash algorithm of the version the file declares, never changes that version, and refuses any other version by name with `roundfix upgrade` as the next action.

## Sources

Two sources, adopted (see `references/_index.md`):

1. [`references/2026-11-12-vortex-skills-lock-version-2-refused.md`](references/2026-11-12-vortex-skills-lock-version-2-refused.md),
   the ai-memory message `01a1cfef-7944-7f70-8e88-3476a8596003` of
   2026-11-12, reported by vortex, nature `measured`.
2. [`references/2026-11-12-pantheon-skills-lock-version-2-refused.md`](references/2026-11-12-pantheon-skills-lock-version-2-refused.md),
   the ai-memory message `01a1d081-347f-70c1-976e-f00851ad9a19` of
   2026-11-12, reported by pantheon, nature `measured`.

At the authoring base `docs/backlog/` and `docs/findings/` hold no open entry
about `skills-lock.json`, so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `d79b30b4` (v0.142.0 plus the Backlog
Entry merge #839). The queue delivers this Spec first in cycle BD, then Specs
0384 and 0385; none of their code is needed here. This Spec owns the lock
reading and writing (`internal/skillhash`, `internal/baseline/skills_restore_lock.go`,
`owned_skills_lock.go`, `retired_skills_installed.go`, `skills_restore.go`,
`skills/repository.go`), the Doctor `skills:` next action in
`internal/cli/doctor.go` and the lock next action of
`internal/cli/baseline_update.go`, and the `setup`, `baseline` and `runtime`
references where they describe the lock. It shares, and edits only through
the repository's record commands or narrowly:

- the owned Roundfix skill version pair (`.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`) and `skills/testdata/owned-skill-versions.json`,
  which 0384 and 0385 may also raise; task_03 writes "the next version"
  through the record command, never a number;
- `docs/references/behavior-surfaces.json`, re-recorded by its command;
- `CONTEXT.md`, where task_03 rewords only **Stale Owned Lock Entry**.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; `skills-lock.json` is a file the skills CLI owns, and this Spec
  reads its existing `version` field. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — every change reads and writes
  local files; nothing fetches. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — this Spec records ADR-0393, which
  decides the lock versions Roundfix reads and that it keeps the version it
  found; ADR-0359 (the bundle install's lock update and the Stale Owned Lock
  Entry scan) and ADR-0350 (the Retired Skill proof by the lock hash) are
  applied unchanged except for the algorithm their `computedHash` uses. Not
  applicable, each cited by one of those, with their rules unchanged:
  ADR-0364 (a required external skill is still checked only for presence),
  ADR-0189, ADR-0191, ADR-0246 and ADR-0380 (a replaced catalog skill is
  still retired as before). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized: no Makefile, workflow, `go.mod`, `go.sum`, lint or formatter
  change, and `skills-lock.json` is not changed. The Governed Paths task_03
  touches are the owned Roundfix skill files and the canonical `setup`,
  `baseline` and `runtime` references, measured with
  `speccheck.GovernedPath`, which the maintainer's cycle grants cover.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- `roundfix doctor`, `roundfix baseline update`, the release plan's `skills`
  check and `roundfix skills install --target project` work on a version 2
  lock as they work on a version 1 lock.
- Every `computedHash` Roundfix writes into a version 2 lock is the value the
  skills CLI 1.8 computes for the same folder, so the CLI sees the skill as
  tracked; a version 1 lock keeps today's hash and bytes.
- No next action Roundfix prints for a lock version it refuses names a
  command that writes that version again.

## User Stories

1. As an adopter whose `skills update` wrote a version 2 lock, I want
   `roundfix doctor`, `baseline update` and the release plan to read it, so
   that I can drop the hand-made version 1 workaround.
2. As an adopter whose Roundfix skill entry is stale after `skills update`, I
   want `roundfix skills install --target project` to repair it in the
   version 2 lock with the hash the skills CLI expects, so that the next
   `skills experimental_install --frozen-lockfile` keeps the copy.
3. As the maintainer, I want an unknown future lock version refused by name
   with `roundfix upgrade`, so that Roundfix never guesses a hash and never
   sends the user in a loop.

## Core Features

1. The Repository Skill Set readiness and the Baseline lock loader accept the
   lock versions 1 and 2 and refuse any other number with the message of
   `_techspec.md` → API Contract 1.
2. The bundle install's lock update, the Stale Owned Lock Entry scan, the
   Retired Skill proof and the external skill restore compute `computedHash`
   in the algorithm of the lock's version, and no writer changes the version.
3. Doctor names `roundfix upgrade` as the only next action of a `skills:` line
   that fails on an unsupported lock version, and `roundfix baseline update`
   names `run roundfix upgrade, then rerun roundfix baseline update`.
4. The `doctor`, `skills` and `baseline` guides, the Roundfix skill's `setup`,
   `baseline` and `runtime` references and the glossary describe the versions.

## User Experience

Before this Spec, in the reported shape, `roundfix doctor` printed
`skills: failed (validate skills lock version "<path>": got 2, want 1;
skills-lock.json comparison unavailable: skills-lock.json must use version 1
with a skills object.; next: bunx skills experimental_install && bunx skills
update -p -y)`, `roundfix baseline update` failed with `inspect
skills-lock.json: skills-lock.json must use version 1 with a skills object.`,
and `roundfix skills install --target project` installed the files and then
exited 2 without updating the lock. After it the same repository passes the
`skills:` line, or warns `DR-SKILL-LOCK-STALE` with `next: roundfix skills
install --target project` when the Roundfix skill entry is stale, and that
install rewrites the entry with the version 2 hash and keeps `"version": 2`.
A version 3 lock prints `skills: failed (validate skills lock version
"<path>": skills-lock.json version 3 is not supported; this Roundfix reads
versions 1 and 2; … ; next: roundfix upgrade)`.

## Non-Goals / Out of Scope

- The Roundfix skill arriving below the binary's minimum through `skills
  update`, the second point of vortex's message: the catalog
  `marcioaltoe/skills` syncs from this repository once a day and the
  maintainer asked the skills project to trigger the sync; the known
  `DR-SKILL-LOCK-STALE` path repairs it with `roundfix skills install
  --target project`, which this Spec makes work on a version 2 lock.
- Accepting any version whose `skills` shape matches (vortex's suggestion):
  ADR-0393 rejects it.
- Migrating a lock between versions; this repository's own version 1 lock
  stays version 1.
- Following a symbolic link inside a skill folder, which the CLI 1.8 hash
  allows; Roundfix keeps refusing one.
- Any change to the Makefile, workflows, `go.mod`, `go.sum`, lint or
  formatter configuration, the Baseline, `skills-lock.json` or the Catalog
  Skill Bump.
- Notifying vortex and pantheon; the operator sends the messages after the
  release.

## Success Metrics

1. Success Metric: on a version 2 lock, the Doctor `skills:` line reads
   `ok`, or `warn` with `DR-SKILL-LOCK-STALE`, and after `roundfix skills
   install --target project` the lock keeps `"version": 2` and holds the
   version 2 hash of the installed Roundfix skill; on the authoring base the
   line fails on the version and the install exits 2.
2. Success Metric: the version 2 digest equals the three test vectors of the
   skills CLI README and the digest the published CLI 1.8.1 computed for a
   folder with CRLF text, `metadata.json`, `__pycache__` and `node_modules`;
   the version 1 digest of that folder is unchanged.
3. Success Metric: a version 3 lock is refused by name in every reader, with
   `roundfix upgrade` as Doctor's only next action and `run roundfix upgrade,
   then rerun roundfix baseline update` as the update's.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted vortex and pantheon messages, measured on repositories this
  Spec did not build on 2026-11-12 with skills 1.8.1 and Roundfix 0.141.1:
  `roundfix baseline update` → "skills-lock.json must use version 1 with a
  skills object" (vortex), and Doctor's next action "é justamente o comando
  que grava a versão 2, então ele entra em ciclo" (pantheon).
- The skills CLI README (https://github.com/vercel-labs/skills, section "Lock
  File Content Hash", read through Exa on 2026-11-12): "`skills-lock.json`
  (version 2) records a `computedHash` for each project skill." and the
  steps "Leave out `metadata.json` files and the `.git`, `__pycache__`, and
  `__pypackages__` folders", "For a text file (a file that contains no NUL
  byte), replace each CRLF with LF", "Sort the files by the UTF-8 bytes of
  their paths. Do not use a locale-aware sort.", "feed `path`, a NUL byte,
  the content length in bytes as decimal ASCII, a NUL byte, and the content
  into one SHA-256", with the test vectors `none` →
  `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`,
  `SKILL.md` = `# hi\n` →
  `93c36952d8843b33e68ae970ebabcc03e1687c9bb0dd1ecaf357060318952f55` and
  `SKILL.md` = `a\n`, `docs/b.md` = `b\n` →
  `b22c4d22aa319d2924312fa18fcbfedc88fa89c74dbd8f298c683e1caf7b7f24`; and
  "A version 1 lock is upgraded when it is read, without network access."

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-12 against these sources.

- Reading at `d79b30b4`: `readLocalSkillsLock` (`skills/repository.go`)
  refuses `lock.Version != 1` with `got %d, want 1`; `loadSkillsLock`
  (`internal/baseline/skills_restore_lock.go`) refuses anything but version 1
  with `skills-lock.json must use version 1 with a skills object.`; every
  `computedHash` is `externalSkillsLockDigest` → `skillhash.Sum`, the
  version 1 digest, in `RecordInstalledOwnedSkills`, `StaleOwnedLockEntries`
  (`owned_skills_lock.go`), `InstalledRetiredSkills`
  (`retired_skills_installed.go`) and `expectedSkillsLockEntry`
  (`skills_restore.go`); `doctorSkillReadinessResult` (`internal/cli/doctor.go`)
  names `bunx skills experimental_install && bunx skills update -p -y` for
  every external failure without a skill name; `runReleasePlanChecks` reads
  Doctor's `repositorySkillsCheck`; the review scope reads only the lock's
  `skills` object and the Catalog Skill Bump writes this repository's lock
  only through `roundfix baseline skills restore`.
- The published skills 1.8.1 (`npm pack skills@1.8.1`,
  `dist/_chunks/remove.mjs`): `readLocalLock` accepts versions 1 and 2 and
  refuses others with "unsupported lock version"; `migrateComputedHash`
  replaces a version 1 hash only when `computeLegacySkillFolderHash` of the
  installed folder equals it; `computeSkillFolderHash` is the README's
  version 2 digest.
- Characterization with the CLI 1.8.1 module and a Roundfix built at
  `d79b30b4`: the CLI migrated all 33 entries of this repository's lock with
  no empty hash, so Roundfix's version 1 digest equals the CLI's legacy
  digest on each; on a crafted folder (CRLF text, `metadata.json`,
  `__pycache__/x.pyc`, `node_modules/pkg/index.js`, `references/B.md`,
  `references/a.md`, `_x.md`) Roundfix's version 1 digest was
  `416d74f34b15134cda0d15e3b9e7dd32f3b0fb6aaf9369c51d9931b37138dc73`, which the
  CLI migrated, and the CLI's version 2 digest was
  `e5adbe14f1b9b04eb73f28a1fc7e58f083b72328dbb7271c4a2813cbcadb42ac`. On a
  clone of this repository migrated to version 2 by the CLI, the Roundfix of
  `d79b30b4` reproduced the reported Doctor line verbatim, `baseline update`
  exited 1 with "Next action: repair skills-lock.json and rerun roundfix
  baseline update", the release plan's `skills` check failed, and `skills
  install --target project` installed 41 files and exited 2.
- A prototype of the design in a scratch worktree of `d79b30b4`: the new
  tests fail on the authoring base and pass with the change; eight sabotages
  each fail their named test; three pinned expectations move (each set a
  version 2 lock as the invalid case), and `GOOS=windows go vet` passes for
  the four changed packages.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "skills-lock.json version 2 computedHash skills CLI roundfix" --all --files --min-score 0.3`;
  no wiki page decides how Roundfix reads the skills lock, so nothing there
  supports or contradicts this design.

## Glossary

- changes: **Stale Owned Lock Entry**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BD brief for Spec 0383 of 2026-11-12, the adopted messages, the characterization, and the maintainer's answers below.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-12 ~17:05–17:20 BRT), verbatim:

- Question "Três Specs novas saem das mensagens de hoje: A = aceitar o
  skills-lock.json versão 2 (Vortex e Pantheon travados no update), B = o
  dono da fila não morre mais ao esperar a Machine Verification Capacity +
  âncora 'path:line.' aceita, C = cache Go do roundfix fora do disco interno
  do mini (Drive). Que autorizações e ordem?" → "A → C → B, mesmas do BC
  (Recommended)" — option text: "Push, merge e release logo após cada merge,
  como no ciclo BC. A primeiro porque trava dois adotantes; C antes de B
  porque o disco do mini está abaixo dos 80 GB."
- Context only, nothing to build: question "A skill roundfix rebaixada no
  Vortex vem do catálogo marcioaltoe/skills, que sincroniza com o Roundfix
  uma vez por dia (hoje 04:52 BRT, antes da v0.141.1 sair às 05:20). Disparo
  o 'Sync upstream skills' do catálogo (workflow_dispatch, auto-merge no
  marcioaltoe/skills)?" → free answer "Já pedi para o projeto skills disparar
  a sincronização".

Taken from those and the evidence, with no further interview:

- Versions 1 and 2 are read; any other number is refused by name, never
  guessed: version 2 kept the field shape and changed what `computedHash`
  means, so a matching shape proves nothing (ADR-0393).
- Each hash follows the version the file declares and no writer changes the
  version: a version 1 hash in a version 2 lock would make the CLI report the
  skill modified, and a migration would rewrite every entry on every run.
- The version 2 digest lives next to the version 1 digest in
  `internal/skillhash`, and both lock readers parse the version through one
  function there, so the refusal text is written once.
- Doctor's next action for an unsupported version is `roundfix upgrade`, the
  only step that can make the lock readable; for a malformed lock it stays as
  it is.
- A symbolic link inside a skill folder stays refused; the owned skills hold
  none, and following one is a separate decision.
- This repository's `skills-lock.json` is version 1 and stays version 1: the
  Catalog Skill Bump writes it only through `roundfix baseline skills
  restore`, which keeps the version; no Task changes it.
- Three implementation Tasks and the QA gate. `type: fix`; the release is the
  operator's patch after the merge; no Task releases.

## Open Questions

None.

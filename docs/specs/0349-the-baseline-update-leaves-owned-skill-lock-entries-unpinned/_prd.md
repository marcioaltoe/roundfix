---
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
status: active
created: 2026-11-04
surfaces: [backend, docs]
type: fix
---

# The baseline update leaves owned skill lock entries unpinned

Minimal PRD for a bug fix, minted by `write-techspec`. Reported by
marcioaltoe-skills on 2026-11-04: `roundfix baseline update` installs the
Roundfix-owned skills from the binary and rewrites only the `computedHash` of
their `skills-lock.json` entries, so each entry keeps the `ref` of the older
catalog commit the upstream tool installed. The published skills CLI clones
the catalog at that `ref`, sees a different hash and reinstalls the older
copy, which undid Quaestor's update and failed `roundfix doctor` again. After
this Spec the bundle install writes the installed tree's hash and no `ref`,
the form the skills CLI itself writes for a default-branch install, and
`roundfix doctor` warns, naming each owned skill whose entry still pins
another copy.

## Sources

One adopter report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md`](references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md),
   `type: fix`, marcioaltoe-skills, ai-memory message `01a1a5f7`, with the
   skills CLI 1.7.6 measurement and the 0.130.1 characterization.

At the authoring base `docs/backlog/` holds one open entry, about the
waiting declaration of the backlog contracts (Spec 0346 adopts it), and
`docs/findings/` holds none; neither shares this context.

## Prerequisites

None in `requires`. Authored against `6a0ed904` (v0.130.1) and rebased on
`66ac68e3`, which adds the plans of Specs 0345 to 0348. Those deliver before
this one and Spec 0350 after it; none changes
`internal/baseline/owned_skills_lock.go`, the skills stage of
`roundfix baseline update`, `roundfix skills install` or the Doctor `skills:`
line, so no Verification depends on another Spec. Files this Spec shares with
them, edited narrowly so later plans rebase cleanly:

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` and
  `skills/testdata/owned-skill-versions.json` (the Roundfix owned skill's
  version, raised by the record command; 0345 and 0347 raise it too);
- `docs/references/behavior-surfaces.json` (re-recorded by its command);
- `CONTEXT.md` (one new entry);
- `internal/cli/cli.go`, where Spec 0347 changes the command config loader
  and the `runs` help, and this Spec changes only the lock report lines of
  `runSkillsInstallCommand`.

Spec 0350 owns the root `skills-lock.json`, whose 33 entries are all external:
this Spec changes neither that file nor any external lock entry.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change edits members of an existing `skills-lock.json` entry.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the commands read and write local files and `baseline update`
  stays offline. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0189: its addendum of 2026-10-27
  made the bundle install rewrite `computedHash` and keep "`source`, `ref`,
  `sourceType`, `skillPath` and every other field"; this Spec narrows that
  field rule. ADR-0359 (new) records the narrowing and the Doctor warning.
  The ADRs one citation edge away apply as recorded: ADR-0233 (the skill
  version raise regenerated at merge) and ADR-0256 (a release waits for the
  covering skills) govern task_03's Roundfix skill edit; ADR-0250 (a module
  version chosen when recorded) is not applicable, no Baseline module
  changes; ADR-0349 (the Release Ceiling) is the `minor` ceiling the
  authorization records; ADR-0197, ADR-0272, ADR-0279 and ADR-0286 (pre-PR
  review rounds, review records and Correction Reviews) are not applicable,
  no review behavior changes. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AX answers of 2026-11-04 quoted in the authorization record (Rules
  "Mesmas do AW (Recommended)"), which keep the grant for the Governed Paths
  each Spec declares; the governed files are the canonical Roundfix skill
  references `baseline.md` and `setup.md` and the Roundfix `SKILL.md` pair.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- After `roundfix baseline update` or `roundfix skills install --target project`,
  each existing lock entry naming an installed owned skill carries the
  installed tree's `computedHash` and no `ref`; its other members, their order
  and every other entry keep their bytes.
- Both commands report the entries whose `ref` they removed, in text and, for
  `baseline update`, in an additive JSON field.
- `roundfix doctor` names each installed owned skill whose entry carries a
  `ref` or another hash, as a warning that never changes its exit code, with
  the command that rewrites the entries.
- The user guides, the Roundfix skill and `CONTEXT.md` state the rule.

## Core Features

1. The Unpinned Owned Lock Entry. `RecordInstalledOwnedSkills` removes the
   `ref` member of each entry it records, beside the hash it already writes,
   and returns which entries changed in which way.
2. The reports. `roundfix skills install` prints
   `Removed skills-lock.json ref: <name>`; `roundfix baseline update` prints
   `- lock ref removed <name>` and adds `skills.lockUnpinned`.
3. The **Stale Owned Lock Entry** warning. Doctor's `skills:` line appends
   `DR-SKILL-LOCK-STALE` with the sorted names, turns `ok` into `warn` and
   adds `roundfix skills install --target project`.
4. The rule in writing: ADR-0359, the doctor, baseline and skills command
   guides, the Roundfix skill's `baseline` and `setup` references and
   `CONTEXT.md`.

## Non-Goals / Out of Scope

- Writing a catalog commit as `ref` (Suggestion A of the report): the binary
  carries no catalog commit and Roundfix never invents one (ADR-0359).
- Changing `source`, `sourceType` or `skillPath` of an owned entry, creating
  an entry, or touching an external entry, its `ref` or its hash.
- Fetching the catalog, or any network access, during `baseline update` or
  Doctor.
- `roundfix skills check`, whose readiness contract stays unchanged.
- Doctor's existing next action for external skills
  (`bunx skills experimental_install && bunx skills update -p -y`) and the
  catalog's own sync schedule.
- This repository's root `skills-lock.json` (Spec 0350) and the Baseline
  modules, profiles, Setup Snapshots and derived Baseline files.

## Success Metrics

1. Success Metric: in a temporary repository whose lock holds owned entries
   pinned to `8d71fade68c54693bedbd917371f0514311af157` in the skills CLI
   format, `roundfix skills install` leaves each installed owned entry with the
   installed hash and no `ref`, keeps the `ref` of an owned entry whose skill
   is not installed, and keeps every external entry byte for byte; at
   `6a0ed904` each `ref` stays.
2. Success Metric: the `ref` member is removed with one separator wherever it
   sits in the entry (first, between members, last, alone), and the document
   keeps every other byte.
3. Success Metric: `roundfix skills install` prints exactly the Surface
   Transcript of `_techspec.md`; `roundfix baseline update` prints
   `- lock ref removed <name>` after the `- lock updated` lines and reports
   `skills.lockUnpinned`; a second run changes nothing and prints neither.
4. Success Metric: Doctor's `skills:` line is
   `warn (<counts>; DR-SKILL-LOCK-STALE: ...: <names>; next: roundfix skills install --target project)`
   for installed owned entries with a `ref` or another hash, never names an
   external entry, keeps a `failed` or `unversioned` status, and reads `ok`
   again after one bundle install; Doctor's exit code does not change.
5. Success Metric: the doctor, baseline and skills command guides, both copies
   of the Roundfix `baseline` and `setup` references and `CONTEXT.md` state the
   rule, and the Roundfix skill's version rises.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The marcioaltoe-skills report quoted in the adopted reference, written by
  the catalog project before this Spec existed, with Quaestor commit
  `742877d` and its measured downgrade.
- The published package `skills@1.7.6` (`npm view skills version` →
  `1.7.6`; tarball SHA-256
  `7ab1d2a52f1e0058d2f7d4f425821e6346ce90526304594f2275f008b816563e`),
  `package/dist/cli.mjs`, read offline on 2026-11-04. Its `updateProjectSkills`
  clones each group with `cloneRepo(cloneSource, skillsForSource[0].entry.ref)`
  and skips a skill only when
  `computedHash === skill.entry.computedHash && onDisk.has(skill.name)`;
  `cloneRepo` clones `["--depth", "1"]` without a ref, the default branch; and
  `projectLockEntry` writes `ref: parsed.ref`, which a default-branch install
  leaves out of the JSON.

## Unreachable Acceptance

None. Every criterion is checked with this repository's files, temporary
repositories and the quoted passages; no row runs the skills CLI or reaches
the catalog.

## Validation

The decisions were checked on 2026-11-04 against these sources.

- Characterization with the built 0.130.1 binary (`6a0ed904`) in a temporary
  repository (adopted reference → Characterization): `roundfix skills install`
  rewrote only the three owned `computedHash` values and kept each `ref`.
  The same fixture with a prototype of this design printed the Surface
  Transcript of `_techspec.md` and removed exactly the three `ref` lines.
- Reading at `6a0ed904`: `RecordInstalledOwnedSkills` is the only writer of
  owned lock entries, called by the `baseline update` skills stage and the
  project target of `skills install`; the skills restore writes external
  entries with the Setup Snapshot's `ref`, coherent by construction;
  `repositorySkillsCheck` judges owned skills by version only and feeds
  `roundfix release plan`'s `skills` check; this repository's root lock holds
  no owned entry.
- The catalog checkout at `c8b2e85` holds the `roundfix` skill byte-identical
  to `.agents/skills/roundfix` at `6a0ed904`, and its registry syncs it daily
  from Roundfix's `main`.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "skills-lock.json ref computedHash skills update downgrade baseline" --all --files --min-score 0.3`.
  It returned only mirrors of Roundfix's own baseline guide and ADRs 0224 and
  0350; none decides how an owned entry pins.
- Exa: a write-up of the project lock (Eric J. Ma, "Skills as reusable
  playbooks") states "The `ref` field holds the `git` branch or tag the
  installer used when it recorded one", and the Skillsaw lock rule lists
  `ref` among the optional fields; both agree with the measured code.

## Glossary

- adds: **Stale Owned Lock Entry** — an owned skill's lock entry that pins another copy
- not a term: **Unpinned Owned Lock Entry** — a Core Feature heading for the written state; `CONTEXT.md` defines the stale state it ends

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AX brief for Spec 0349 of 2026-11-04 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-04), verbatim:

- 0349: "Lock coerente + doctor avisa (Recommended)" — option text: "Spec
  nova: o baseline update grava um ref coerente com o conteúdo que instala
  (ou remove o ref das entradas que o Roundfix gerencia, a Spec escolhe depois
  de medir o skills CLI) e o `roundfix doctor` avisa quando ref e
  computedHash de uma entrada não batem, nomeando o rebaixamento que o skills
  update faria."
- Placement: "No ciclo AX, mesma v0.131.0 (Recommended)" — option text:
  "Ficam depois da 0346 na fila (as três mexem em Baseline/skills); uma
  v0.131.0 no fim com as seis; mesmas regras do AX (merge pela fila,
  Makefile/CI liberados). Nota ao catálogo e aos adotantes ao final."
- Rules: "Mesmas do AW (Recommended)" — option text: "Sem supervisão: push e
  PR dos planos, entrega e merge pela fila, consertar e seguir em falha
  técnica; Makefile, .github/workflows, go.mod/go.sum e lint liberados quando
  o desenho pedir; mudança de Baseline que chega aos adotantes liberada."
- Release: "Uma v0.131.0 no fim (Recommended)" — option text: "Uma minor (a
  mudança da Baseline e o runs list fora de repo são feat); notas aos 8
  adotantes, resposta ao Pantheon com o defeito corrigido e o contorno que
  podem largar."

Taken from those answers and the measurement, with no further interview:

- The coherent state is "no `ref`" (Suggestion B), not a written catalog
  commit (Suggestion A): the binary knows no catalog commit, and an entry
  without `ref` is the skills CLI's own default-branch form (ADR-0359).
- Doctor warns (Suggestion C) with `warn`, never `failed`: the installed copy
  is current until someone runs `skills update`.
- Three implementation Tasks in series and a QA gate.

## Open Questions

None.

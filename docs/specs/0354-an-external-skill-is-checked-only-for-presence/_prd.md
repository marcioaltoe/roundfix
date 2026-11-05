---
spec: 0354-an-external-skill-is-checked-only-for-presence
status: active
created: 2026-11-05
surfaces: [backend, docs]
type: fix
---

# An external skill is checked only for presence

Minimal PRD for a bug fix, minted by `write-techspec`. Reported by fluxus on
2026-11-05, and the same morning from a Vortex session by the maintainer:
with Roundfix 0.132.0, `roundfix doctor` appended
`DR-SKILL-DIFFERS-FROM-SNAPSHOT` for `systematic-debugging` and
`ui-ux-pro-max`, and `roundfix baseline update` listed both under
`skills.differing` with a suggestion to restore them from the Setup Snapshot.
Those copies come from `marcioaltoe/skills`, the source the maintainer chose,
so the suggested restore would replace them with an older or unrelated copy.
A copy that no longer matches its `skills-lock.json` hash fails Doctor as
`outdated` and is restored from the snapshot by `baseline update --yes`. After
this Spec Roundfix checks the content only of the skills it authors, against
the copy inside the binary, and checks an external skill only for presence.

## Sources

One report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md`](references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md),
   `type: fix`, reported by fluxus (ai-memory message
   `01a1ab62-5139-7dc3-9b92-472690a044f9`, with the maintainer's Vortex
   report of the same morning).

At the authoring base `docs/backlog/` and `docs/findings/` hold no open
entry, so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `ad77a797`, where Specs 0352 and 0353
are active and deliver before this one; this Spec touches none of their
files' lines. It owns the `skills` package's Repository Skill Set check, the
Doctor `skills:` line and the skills stage of `baseline update`. Files it may
share with another Spec of the queue, kept to narrow edits: the owned
Roundfix skill version (`.agents/skills/roundfix/SKILL.md`,
`skills/roundfix/SKILL.md`, `skills/testdata/owned-skill-versions.json`),
the Roundfix `baseline`, `setup` and `runtime` references, `CONTEXT.md`, the
Behavior Surface Record (`docs/references/behavior-surfaces.json`) and the
coverage record (`docs/references/coverage-record.json`); 0352 and 0353 also
bump the Roundfix skill version, and 0355 to 0358 may.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change reads existing skill folders and lock entries. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; every check reads local files offline. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0364 (new) records the decision:
  "A required external skill is checked only for presence". It amends
  ADR-0221 (an external skill held to its lock and its Setup Snapshot),
  ADR-0264 (the Setup Snapshot Difference Doctor names and the managed
  refresh lists), ADR-0224 (this repository's own skill contract tests) and
  ADR-0241 (a repository profile's skill contracts, read by the snapshot
  comparison among others), whose amendment notices task_03 writes.
  ADR-0360 (the Catalog Skill Bump, which cites ADR-0224) and ADR-0322 (the
  Python profile, which cites ADR-0241) keep their rules: this repository's
  bump and the Python preset are untouched. The owned-skill version rule,
  the Stale Owned Lock Entry warning and the deletion of Retired Skills keep
  their records unchanged. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 2 answers of 2026-11-05 quoted in the authorization record;
  the Governed Paths are the coverage record task_01 re-records and the
  Roundfix skill files task_03 changes. No
  Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0354-an-external-skill-is-checked-only-for-presence/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/runtime.md`, `.agents/skills/roundfix/references/setup.md`, `docs/references/coverage-record.json`, `skills/roundfix/SKILL.md`.

## Goals

- A required external skill whose folder and lock entry exist is ready,
  whatever its content: Doctor neither fails it as `outdated` nor names a
  Setup Snapshot difference.
- `roundfix baseline update` never lists a present external skill under
  `skills.differing` or `skills.missing`, and never restores it.
- A required external skill without its folder or its lock entry is still
  reported by Doctor, `roundfix skills check` and the `baseline update`
  preview, and `baseline update --yes` still restores it.
- `roundfix baseline skills restore --profile <id> --skill <name>` still
  restores a present external copy when the operator asks for it.
- An owned skill keeps passing when its version is equal to or higher than
  the binary's and failing when lower, offline.
- The guides, the README, the Roundfix skill, `CONTEXT.md` and the amended
  ADRs state the rule.

## Core Features

1. Presence check. The Repository Skill Set check reads, for an external
   skill, only whether `skills-lock.json` names it and whether its folder is a
   directory; no file of the folder and no `computedHash` is read.
2. No snapshot comparison. Doctor and the skills stage of `baseline update`
   stop comparing external skills with the Setup Snapshot, and the
   `skills.differing` field and its text lines go away.
3. Restore only what is missing. The managed refresh restores a Missing
   Required Skill and nothing else; the deliberate per-skill restore keeps
   working.
4. This repository keeps its own stricter rule: its contract tests still hold
   its Vendored Skills to their lock entries and their Setup Snapshot.
5. The rule in writing: the Doctor and Baseline guides, the README, the
   Roundfix `baseline`, `setup` and `runtime` references, two glossary
   entries, ADR-0364 and the amendment notices in ADR-0221, ADR-0224,
   ADR-0241 and ADR-0264.

## Non-Goals / Out of Scope

- The owned-skill version rule and its text, the Stale Owned Lock Entry
  warning, the lock rewrite of the bundle install, and the deletion of
  Retired Skills that reads the lock hash.
- Any network read: Roundfix fetches no catalog to compare a skill.
- The skills lock format, the `skills` CLI, `roundfix baseline skills
  reconcile`, the Setup Snapshot schema and the asset sync.
- This repository's Catalog Skill Bump, `skills-lock.json` and Vendored
  Skills.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, any Baseline module or
  guide under `docs/agents/`, `CHANGELOG.md`.

## Success Metrics

1. Success Metric: in a repository whose required external skill
   `testing-boss` is installed, named by its lock entry and differs from its
   Setup Snapshot, Doctor prints `skills: ok (<n> required: <owned>
   Roundfix-owned, <external> external)` with no `DR-SKILL-DIFFERS-FROM-SNAPSHOT`;
   at `ad77a797` it appends `DR-SKILL-DIFFERS-FROM-SNAPSHOT: differs from the
   Setup Snapshot, newer or older unknown: testing-boss`.
2. Success Metric: when that copy also differs from its lock entry's
   `computedHash`, Doctor still prints `skills: ok (...)` and exits as the
   other checks decide, and `roundfix baseline update --no-history` reports
   `current` with exit `0`; at `ad77a797` Doctor prints `skills: failed
   (outdated: testing-boss; next: bunx skills add
   marcioaltoe/skills@testing-boss)` and the preview exits `3` with `1
   required skill(s) are missing or no longer match skills-lock.json:
   testing-boss`.
3. Success Metric: the `baseline update` text output has no `Skills
   differing:` line and its JSON has no `differing` member, and `baseline
   update --yes` leaves a present external copy byte for byte.
4. Success Metric: a required external skill without its folder is reported
   by Doctor as `missing: <name>` with `bunx skills add
   marcioaltoe/skills@<name>`, by the preview as `1 required skill(s) are
   missing: <name>`, and restored by `baseline update --yes`.
5. Success Metric: `roundfix baseline skills restore` still previews and,
   with its Plan Digest, replaces a present external copy with its snapshot
   copy.
6. Success Metric: an owned skill whose declared version is above the
   binary's passes Doctor and one below fails it with `below minimum`, with
   no network.
7. Success Metric: both copies of the Roundfix `baseline`, `setup` and
   `runtime` references, the Doctor and Baseline guides and the README say
   that an external skill is checked only for presence, none of them still
   names `DR-SKILL-DIFFERS-FROM-SNAPSHOT` or `skills.differing`, and
   ADR-0221, ADR-0224, ADR-0241 and ADR-0264 carry amendment notices naming
   ADR-0364.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The fluxus report quoted in the adopted reference, written before this
  Spec existed: Doctor's `DR-SKILL-DIFFERS-FROM-SNAPSHOT: ...
  systematic-debugging, ui-ux-pro-max` line and the `skills.differing`
  entries with "restore deliberately with roundfix baseline skills restore",
  observed in Fluxus with Roundfix 0.132.0, and the Vortex prompt "Restauro
  essas duas para o snapshot do catálogo?".
- The skills CLI's own issue tracker, vercel-labs/skills issue 806
  (https://github.com/vercel-labs/skills/issues/806, opened 2026-03-31, read
  through Exa on 2026-11-05): "The computedHash field in project-level
  skills-lock.json is computed from the source repository directory ...
  the stored hash cannot be recomputed from the installed copy." This is why
  a lock-hash mismatch is no proof that an installed copy is stale.

## Unreachable Acceptance

None. Every criterion is checked with this repository's files, temporary
homes and temporary repositories.

## Validation

The decisions were checked on 2026-11-05 against these sources.

- Characterization with this clone's `bin/roundfix` at `ad77a797`, a
  temporary `HOME` and no network, in a temporary copy of this repository
  (profile `go-cli-tui`) whose `.agents/skills/testing-boss/SKILL.md` gained
  one line:
  - with the lock's `computedHash` rewritten to the edited copy's, `doctor`
    printed `skills: ok (41 required: 12 Roundfix-owned, 29 external;
    DR-SKILL-DIFFERS-FROM-SNAPSHOT: differs from the Setup Snapshot, newer or
    older unknown: testing-boss)` and `baseline update --no-history --format
    text` exited `0` with `Skills differing: 1` and `- differs testing-boss:
    differs from its Setup Snapshot, newer or older unknown; kept; restore
    deliberately with roundfix baseline skills restore --profile go-cli-tui
    --skill testing-boss`;
  - with the lock unchanged, `doctor` printed `skills: failed (outdated:
    testing-boss; next: bunx skills add marcioaltoe/skills@testing-boss)`,
    the preview exited `3` with `1 required skill(s) are missing or no
    longer match skills-lock.json: testing-boss` and `Next action: rerun
    with --yes to restore them from the Setup Snapshot, ...`, and `skills
    check` failed with the same reason;
  - with `qa-gate` declaring `0.0.99`, `doctor` printed `skills: ok (...)`;
    with `0.0.1`, `skills: failed (below minimum: skill "qa-gate" requires
    0.0.24, found 0.0.1; next: roundfix skills install --target project)`.
- A prototype of the design on a scratch copy at `ad77a797` built, and its
  `skills`, `internal/baseline` and `internal/cli` suites, with and without
  the `docscontract` tag, failed only in the tests `_techspec.md` → Declared
  breaks lists.
- Secondbrain, read only: `wiki/index.md`, then `qmd query "roundfix doctor
  skills snapshot external skill digest version owned skills
  marcioaltoe/skills" --all --files --min-score 0.3`. Its best matches were
  the Roundfix mirror of Spec 0036 and the Doctor guide, and the Fluxus
  record `projects/fluxus/mirror/docs/references/2026-10-06-autorizacao-de-tooling-para-as-skills-do-repo-skills.md`,
  where on 2026-10-06 the maintainer already chose the `marcioaltoe/skills`
  copy over the Setup Snapshot for a skill and was told Doctor would warn and
  `baseline update` would restore it: the same conflict, a month earlier.
- Exa: the skills CLI's lock format (`src/skill-lock.ts`, pull request 565
  and the DeepWiki pages on `check` and `update`), and issue 806 quoted above.

## Glossary

- changes: **Missing Required Skill**
- changes: **Setup Snapshot Difference**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 2 brief for Spec 0354 of 2026-11-05, the maintainer's request in the adopted report and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-05 ~08:20 BRT), verbatim:

- Question "Pedido do Fluxus (Spec 0354): as skills do próprio Roundfix devem
  ter o mesmo digest do marcioaltoe/skills, ou versão igual ou maior. Contra
  qual cópia o doctor confere?" → "Cópia do binário, sem rede (Recommended)"
  — option text: "O catálogo marcioaltoe/skills copia as skills do Roundfix
  do main todo dia, então a cópia dentro do binário é a mesma ou mais nova. O
  doctor aceita digest igual ou versão igual ou maior que a do binário, sem
  acessar a rede."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, com
  merge pela fila e conserto quando a entrega parar. Makefile, CI, go.mod e
  lint liberados, assim como mudanças de Baseline e skills que chegam aos
  adotantes. Release assim que cada Spec ficar pronta, com notas aos
  adotantes, aviso ao Fluxus e ao Fiscus."
- The maintainer's request in the adopted report, points 1 to 3 and "Efeito
  esperado", is the rest of the decision: version and digest only of the 12
  Roundfix-owned skills; the others "só se estão instaladas, sem digest nem
  versão"; "falta de uma skill externa obrigatória continua sendo acusada";
  `baseline update` "não lista skills externas em skills.differing nem as
  restaura".

Taken from those answers, the report and the characterization, with no
further interview:

- Presence is a lock entry and a folder, the skills CLI's own record of an
  installation; no file of the folder and no `computedHash` is read.
- The lock refusals (absent, unreadable, malformed, wrong version) and the
  refusal of a symbolically linked skill root stay, so no repository changes
  state for a reason other than content.
- The owned-skill check already compares with the binary's copy offline
  (its version record), and an equal digest implies an equal version; it changes no
  code.
- `roundfix baseline skills restore` keeps restoring a present external copy
  on demand; nothing suggests it any longer.
- `skills.differing` and the `Skills differing:` lines are removed, not kept
  empty, and the release notes name the removal.
- This repository's own contract tests keep its lock and snapshot rule.
- ADR-0364 records the decision and amends ADR-0221, ADR-0224, ADR-0241 and
  ADR-0264.
  Three implementation Tasks in series and a QA gate.

## Open Questions

None.

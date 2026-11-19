---
status: accepted
knowledge: true
created_at: 2026-11-19T14:30:00Z
updated_at: 2026-11-19T14:30:00Z
deprecated_at: null
superseded_by: null
---

# In an adopter, a change to the files Roundfix manages or to a CONTEXT.md needs no release

The Release Plan counts an Adopter Maintenance Path as maintenance-only, beside
planning evidence. A repository is an adopter when its `skills-lock.json` at
the range's target records the `roundfix` skill. Roundfix's own repository
authors that skill and records none, so its classification does not change.

## Context

The Fiscus reported on 2026-11-19, through the maintainer, that its next
release needed a classification by hand for four of six commits since
`v0.9.0`: the Roundfix update, the adoption of the Roundfix recommendations,
and two plans. Three of the four changed only files that Roundfix's Baseline
and skill installation write (`AGENTS.md`, guides under `docs/agents/`, the
Setup Manifest, `skills-lock.json`, installed skills, the generated
`.claude/agents/roundfix-*.md`, `.roundfixrc.yml`,
`docs/agents/specific-repository.md`) or a `CONTEXT.md`. None of them changes
what the adopter ships. Every adopter meets this at every Roundfix update
(`references/2026-11-19-characterization.md` of Spec 0414).

The maintainer set the boundary on 2026-11-19: "Arquivos geridos pelo Roundfix
+ CONTEXT.md (Recommended)": the files the Baseline and `skills install` write,
as the Setup Manifest and `skills-lock.json` record them, the generated
`.claude/agents`, `.roundfixrc.yml` and every `CONTEXT.md`; everything else as
today; and for Roundfix itself nothing changes, because its skills and
Baseline are its product.

## Decision

**An adopter.** The Release Plan reads `skills-lock.json` at the target. When
its `skills` object records `roundfix`, the repository installs Roundfix
rather than authoring it, and the rule below applies. A missing, unreadable or
malformed lock makes no adopter, and the plan classifies as before.

**The Adopter Maintenance Paths.** Read at the range's base and at its target,
and joined:

- the Setup Manifest `docs/agents/setup-context.json` and every `path` its
  `managedArtifacts` records;
- `skills-lock.json` and, for each skill its `skills` object records, every
  path under `.agents/skills/<name>/` and `.claude/skills/<name>`;
- `CLAUDE.md` and `docs/agents/specific-repository.md`, which the Baseline
  writes at fixed paths without recording them per path;
- each `.claude/agents/roundfix-*.md`, the generated Author and Research
  Subagents;
- `.roundfixrc.yml` and every file named `CONTEXT.md`, at any depth.

A commit whose every changed path is planning evidence or an Adopter
Maintenance Path is maintenance-only. A commit with any other path still
crosses the boundary, and a commit without changed paths still crosses it.

**What stays.** The plan still reads no Roundfix configuration: it reads two
committed records at two revisions, and names `.roundfixrc.yml` only as a
path. A Git failure while reading them fails the plan, as a failure reading the
delivered authorization records already does. A repository-local skill that
`skills-lock.json` does not record, a guide the Setup Manifest does not record,
and every other path keep today's classification.

## Considered Options

- **Apply the rule in every repository.** Refused: it changes Roundfix's own
  classification. Over the 216 commits since `v0.125.0` it would have cleared
  `1a8e08b4`, a change to Roundfix's own `.roundfixrc.yml`, against the
  maintainer's "Para o próprio Roundfix nada muda".
- **A Project Config key that marks Roundfix or an adopter.** Refused: the
  Release Plan Command reads no Roundfix configuration, and a key every
  adopter must set would leave the eleven adopters as they are until each sets
  it.
- **Recognise Roundfix by a path of its source tree**, such as
  `internal/baseline/assets/`. Refused: it ties the classifier to Roundfix's
  internal layout, and a fork or a move breaks it silently.
- **A fixed path list**, `AGENTS.md`, `docs/agents/**`, `.agents/skills/**`
  and the rest. Refused: it counts a repository-local skill with its scripts,
  such as the Fiscus `revisao-tributaria-cadastro`, and a guide the repository
  wrote itself, as files Roundfix manages. The records say exactly what the
  Baseline and the installation wrote.
- **Read the records at every commit.** Refused: the union of the base and the
  target covers a skill or guide added or retired inside the range at two
  reads per record, and a skill both added and removed inside one range only
  keeps the conservative classification.

## Consequences

- A Roundfix update, the adoption of its recommendations and a plan that
  touches a `CONTEXT.md` no longer stop an adopter's Release Plan, nor the
  `Next release` line of the Follow-up Command's status view, which builds the
  same plan.
- The Fiscus range `v0.9.0..ad1381b` still asks for a classification, now only
  for `ad1381b`, which changes `docs/user-guide/` and `.env.example`.
- An adopter that installs no skills, or whose lock lost the `roundfix` entry,
  keeps the old behavior until `roundfix baseline update` reinstalls it.
- A `feat:` or `fix:` commit confined to Adopter Maintenance Paths sets no
  impact, as one confined to planning evidence already does.

---
spec: 0257-the-baseline-respects-the-repository
status: active
created: 2026-10-12
surfaces: [backend, cli, docs]
---

# The Baseline respects the repository's own choices

Three adopter reports show the Baseline overriding choices that belong to the
repository. The Standard TypeScript backend guide says Better Auth owns
`/api/auth/*`, yet Better Auth's base path is each application's choice: one
adopter already removed the `/api` prefix and wrote, outside the setup
markers, that the setup-owned guide contradicts its own decision, and another
queued the same removal. `roundfix doctor` calls a vendored skill newer than
the Setup Snapshot "trailing", and the managed refresh offers to restore the
older snapshot over it. And the Baseline says not to edit a vendored skill,
but not what to do when one diverges. After this Spec the Better Auth base
path is a decision the repository records at its next `roundfix baseline
update`, a skill copy that only differs from the snapshot is never called
trailing or restored unasked, and the Baseline states the order for handling
a diverging vendored skill. The maintainer settled all three on 2026-10-11
("Decisão no Setup Manifest" and the cycle decisions for this Spec).

## Prerequisites

None. Specs 0255 and 0259 are authored beside this one. Spec 0255 also changes
`roundfix baseline update`, the `roundfix` skill's Baseline reference and the
owned-skill version record; Spec 0259 also changes Doctor. This Spec's
Verification reads none of their artifacts. Whichever Spec delivers later
rebases and re-runs the record commands its Tasks name.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes; the new
  decision and the new Doctor code are names, not identifiers of stored
  records. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable. Roundfix gains no HTTP route,
  credential or network call. The Spec changes how the Baseline records and
  renders an adopter's Better Auth route exception, which stays the adopter's
  own decision: ADR-0063: "requires each repository to persist an HTTP
  Contract Decision". The Baseline invents no authentication policy; it
  renders the recorded path. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/cli.md`.
- Active ADR obligations: applicable. ADR-0263 and ADR-0264 (this Spec)
  decide the base-path decision, the declared update decision, the Setup
  Snapshot Difference and the divergence order. ADR-0076: "It represents a
  selected Better Auth capability through an" `auth.provider` decision whose
  exception is merged into `http.contract`; the base path moves that one
  exception instead of adding a second. ADR-0118: "The test is whether a
  second repository could answer differently and still be" correct, and two
  adopters answer the base path differently. ADR-0144: "ask only when a
  decision identifier is missing or incompatible", so the optional decision
  asks nothing of an adopter that does not record it. ADR-0221: "A skill that
  differs is reported as trailing the" snapshot; ADR-0264 replaces that
  classification and its automatic restore. ADR-0224: "here a trailing skill
  fails" `make verify`; this repository keeps that strictness. ADR-0186:
  "Every clause keeps its enforcement level when its wording changes"; the
  three reworded clauses keep theirs. ADR-0257: "A clause whose wording
  changes keeps its identity", so no clause identity is added. ADR-0250: "A
  Baseline module's version names one content, and the record step chooses
  it". ADR-0256: "Spec authoring requires a skills Task for a change a skill
  describes"; every Task that changes a covered surface changes its skill.
  ADR-0244: "A Spec declares the domain terms it introduces". ADR-0063 owns
  the HTTP Contract Decision whose Better Auth exception the base path moves.
  ADR-0241: a repository profile takes its skill contracts from the embedded
  Setup Snapshots, and the difference comparison keeps that resolution.
  ADR-0258 governs authoring against the QA rerun classes; every Task here
  proves its new checks can fail and declares each break. ADR-0150 and
  ADR-0228 settle the branch-prefix decision that ADR-0118 introduced and do
  not apply: no branch rule changes. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable. Express maintainer authorization: "considere
  autorizado a ajustar todas as skills se necessário" for skills, "Autorizar
  os dois" for Baseline sources and guides, and the cycle AE grant of
  2026-10-11, "Concedo", for the Governed Paths each Spec declares. No
  Makefile, CI workflow, lint, formatter or test-runner configuration and no
  new `go.mod` dependency. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0257-the-baseline-respects-the-repository/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `docs/agents/setup-context.json`, `docs/agents/skill-dispatch.md`,
  `docs/references/coverage-record.json`,
  `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/rust.json`,
  `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/templates/guides/backend.md`,
  `internal/baseline/assets/templates/index.json`,
  `skills/roundfix/SKILL.md`, `skills/setup-context-driven/SKILL.md`.

## Goals

- A repository records where it mounts Better Auth, and its backend guide and
  Setup Manifest name that path, with no duplicate exception. A repository
  that records nothing keeps `/api/auth/*` and byte-identical guides.
- An adopter changes that path, or any other decision its profile selects, in
  one confirmed `roundfix baseline update`.
- Doctor and the managed refresh never call a lock-matching skill that only
  differs from its Setup Snapshot "trailing", and never restore one unasked.
- The Baseline and the `roundfix` skill state, in order, how a diverging
  vendored skill is handled, and every divergence the skills maintainer
  listed has an overriding Baseline rule.

## User Stories

1. As an adopter that mounts Better Auth at `/auth`, I want to record that
   path once, so that my backend guide stops describing a route my
   application no longer serves.
2. As an adopter that never changed the path, I want my next update to leave
   my backend guide as it is, so that a Baseline release costs me nothing.
3. As an operator running an agent, I want to change one recorded decision
   with one non-interactive update, so that I do not restate every decision.
4. As an adopter that updates its skills with its own tool, I want Doctor to
   tell me a skill differs from the snapshot without telling me to restore
   it, so that a routine skills update is not undone.
5. As a maintainer reading the Baseline, I want the order for handling a
   diverging vendored skill, so that I correct it in the right place instead
   of editing a folder the next sync replaces.

## Core Features

1. **Better Auth Base Path.** The Standard TypeScript profiles offer an
   optional base-path decision, defaulting to Better Auth's own `/api/auth`.
   When recorded, the provider's route exception covers that path, the
   matching HTTP Contract Decision exception moves with it, and the backend
   guide states the recorded path. When absent, nothing an adopter holds
   changes. An invalid path, or a base path without the Better Auth provider,
   is refused (ADR-0263).
2. **Declared decision at update.** The managed refresh accepts decisions to
   record or change before it plans, validated as automation planning
   validates them, reported in its result and covered by its Plan Digest. It
   refuses a decision the profile does not select or fixes, and it answers a
   new decision without adopting the suggestion (ADR-0263).
3. **Setup Snapshot Difference.** A required external skill that matches its
   lock but not its Setup Snapshot is reported as differing, with newer or
   older unknown. Doctor names it without a warning or a next action; the
   managed refresh lists it, keeps it and names the deliberate per-skill
   restore. The trailing code is retired, and missing or lock-mismatched
   skills are restored as before (ADR-0264).
4. **Divergence order.** The Baseline's vendored-skill rule states the three
   choices in order and sends objective defects upstream. The Baseline adds
   the two overriding rules still missing from the skills maintainer's list:
   reading a skill's bundled files from its installed directory, and
   following the crate's declared Rust edition. The `roundfix` skill states
   the same order (ADR-0264).
5. **Docs and glossary.** The user guide, the command references, the
   `setup-context-driven` and `roundfix` skills and `CONTEXT.md` describe the
   new behavior; `CONTEXT.md` adds **Better Auth Base Path**, **Setup Snapshot
   Difference** and **Vendored Skill**.

## User Experience

An adopter that mounts Better Auth at `/auth` runs `roundfix baseline update
--decision auth.base_path=/auth`, reviews a plan that changes its Setup
Manifest and backend guide, and confirms it. Doctor's `skills` line keeps
`ok` when a skill only differs from the snapshot and appends the difference
code with the skill names. The update preview lists such skills under their
own heading, apart from the skills it will restore.

## Non-Goals / Out of Scope

- Proving that a skill copy is behind its snapshot: no earlier snapshot digest
  is recorded, so no copy is classified behind (ADR-0264, Considered options).
- Changing the Setup Snapshot schema, the asset sync, the snapshot commit or
  any vendored skill's text.
- Changing how a missing skill, or one that no longer matches its lock, is
  restored.
- Moving `routeException` out of `auth.provider`, or changing its owner,
  methods or reason.
- Letting the update change the profile or the preservation mode.
- Editing any adopter repository, including recording `/auth` for the vortex.
- The `cut-release` dispatch issue and vendoring `oss-publish`, both noted in
  the skills maintainer's earlier report.

## Success Metrics

1. Success Metric: in a temporary Standard TypeScript adopter, `roundfix
   baseline update --decision auth.base_path=/auth --yes` records `/auth`;
   its backend guide names `/auth/*` and the recorded path and no longer
   names `/api/auth/*`; its HTTP Contract Decision holds exactly one Better
   Auth exception; a second update reports `current`.
2. Success Metric: in a temporary Standard TypeScript adopter that records no
   base path, the update after this Spec leaves its backend guide
   byte-identical and adds no `auth.base_path` to its Setup Manifest.
3. Success Metric: in a temporary repository whose required external skill
   matches its lock but not its snapshot, Doctor exits `0` with the skill
   under `DR-SKILL-DIFFERS-FROM-SNAPSHOT` and its `skills` line not `warn`;
   the update preview reports `current` with the skill under
   `skills.differing`; `--yes` leaves the skill's bytes identical. No live
   document, skill or source file names `DR-SKILL-TRAILS-SNAPSHOT`.
4. Success Metric: this repository's rendered skill-dispatch guide and the
   `roundfix` skill state the three choices in order, and each of the nine
   divergences in TechSpec → Divergence audit maps to a Baseline clause.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write:

- Better Auth, options reference (`docs/content/docs/reference/options.mdx`
  in `better-auth/better-auth`), on `basePath`: "Base path for Better Auth.
  This is typically the path where the Better Auth routes are mounted." and
  "Default: `/api/auth`". The path is the application's choice, and
  `/api/auth` is the library's default, not a rule.
- Microsoft APM, "Drift Detection": "`apm outdated` | Compares each locked
  dependency against its remote tip." Ordering a locked copy needs the
  upstream; offline, a content digest only says that two copies differ.
- ingraft, "Editable Vendors": "Most vendored source should be treated as
  read-only reference material." and Safeguard Research, "When to Fork an
  Abandoned Dependency" (2026-04-05): "You should fork an abandoned dependency
  when a concrete blocker exists and no lighter option resolves it." Both put
  replacing vendored code last and sending fixes upstream.

## Glossary

- adds: **Better Auth Base Path**
- adds: **Setup Snapshot Difference**
- adds: **Vendored Skill**

## Decisions

- The base path is an optional decision with Better Auth's default, and the
  provider exception is derived from it; see ADR-0263.
- `roundfix baseline update` records declared decisions; see ADR-0263.
- A lock-matching skill that differs from its snapshot is a difference, kept
  and never restored unasked, and the trailing code is retired; see ADR-0264.
- A diverging vendored skill is corrected by a Baseline or repository rule
  first, a complementary skill second, a replacement last; see ADR-0264.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the decision, the normalization, the
update flag, the comparison, the clause texts, the derived files and the
build order.

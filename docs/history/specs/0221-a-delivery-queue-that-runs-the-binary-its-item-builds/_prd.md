---
spec: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
status: archived
created: 2026-10-03
surfaces: [backend, cli, docs]
archived: "2026-10-03"
source_slug: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; blocked rows need the operator's log outside the repository (09), outbound network the QA sandbox denies (11, 12) and an open Pull Request (15)
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: a9a317dcc0ad48ed17956dbd0e4423a99916141c
---


# A delivery queue that runs the binary its item builds

On 2026-10-02 the Delivery Queue parked Spec 0215 `delivery-error` for a
reason the item could not fix: the queue owner, built from `main`, started
`roundfix implement` in the item worktree, that child read the item's Project
Config, and the item had just added the key `verification.tools` to it. The
child refused with `verification.tools is not a supported config key`. The
operator built `bin/roundfix` from the item branch and retried with it, and
did so again on the next two retries before the item merged. Any Spec that
adds a Project Config key to Roundfix's own repository meets the same park at
its first retry, archive or review. This Spec lets a repository declare how
its Roundfix binary is built, so the queue owner builds that binary in the
item worktree and runs the item's own Implement, archive and review steps with
it, and keeps the shared Run Database out of reach of an item binary that
would change its schema.

## Prerequisites

None. Specs 0219 and 0220 are authored in the same cycle; if either also
raises the Roundfix Skill's version, the operator orders the queue so that the
later Spec raises it from the earlier one's value.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes. Queue
  items keep their Spec slug, Runs their Run identifiers, and the new Project
  Config key follows the existing dotted key names. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the owner builds and starts local
  processes in a worktree it already owns; no request, credential or forge
  read is added, and no test reaches the network. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0225 (this Spec) decides that a
  queue item runs the Roundfix binary its branch builds when the repository
  declares the build, and that a binary which would change the Run Database
  schema is not run. ADR-0192 is the precedent for a repository-declared
  command the owner runs in the item worktree through the Verification
  executor, ADR-0192: "The declaration lives in Project Config because merge
  recovery is a repository policy". The unknown-key refusal stays strict,
  ADR-0027: "truly unknown keys keep failing strict validation". ADR-0193's
  prerequisite wait, ADR-0199's token ceiling and ADR-0211's agent
  environment hold unchanged. ADR-0125: "Fixtures are therefore compiled
  once", which binds the fake item binary every test starts, and ADR-0213
  ends each fixture process with the test binary that started it. ADR-0187
  splits the Roundfix Skill by command and ADR-0189 ties an owned skill's
  version to its content, so the skill edit raises the version. ADR-0184:
  "A TechSpec now declares numbered Surface Transcripts", so the changed
  command surface is stated as a transcript. The gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167, and
  ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check this Spec's
  consistency by citation and receipt. ADR-0178 authorizes each Task commit
  by its grant, ADR-0182 runs Settlement Checks before it, and ADR-0166
  records undeclared paths; every Task declares its paths. ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry, ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot, and ADR-0220 cites ADR-0211 but decides the Doctor's delivery readiness; this Spec changes none of them, so none applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and its
  `SKILL.md` mirror are Governed Paths, and the maintainer authorized skill
  edits ("considere autorizado a ajustar todas as skills se necessário") and
  this cycle ("Continue", 2026-10-03). No other Governed Path changes; the
  repository's own Project Config is not edited by this Spec. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0221-a-delivery-queue-that-runs-the-binary-its-item-builds/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- An item that adds a Project Config key to a repository that declares its
  Roundfix build reaches merge with no operator-built binary.
- The Run Database is never migrated by a binary built from an unmerged item.
- A repository that declares nothing, which includes every adopter that
  installs Roundfix from npm, is delivered exactly as before.
- An operator can ask any Roundfix binary, without writing anything, whether
  it would change the Run Database schema.

## User Stories

1. As the operator of the queue in Roundfix's own repository, I want the
   owner to run each item's Implement, archive and review steps with the
   binary built from that item, so that a Spec adding a configuration key
   does not park on the owner's older binary.
2. As the operator, I want the owner to keep using its own binary when the
   item's binary would migrate the Run Database, and to say so, so that an
   unmerged item never leaves my installed Roundfix unable to open its
   database.
3. As the operator, I want a declared build that fails to park the item with
   the build log, so that a broken declaration is visible instead of silently
   ignored.
4. As a repository maintainer, I want to declare the build command and the
   binary path once in Project Config, so that the queue uses my
   repository's build rather than a convention baked into Roundfix.
5. As an operator holding two Roundfix binaries, I want
   `roundfix migrate --check` to tell me whether this binary would change the
   Run Database, so that I can decide which binary to run before running it.

## Core Features

1. **A declared item build.** Project Config accepts a `delivery.item_binary`
   declaration with a `build` command and the repository-relative `path` of
   the binary that command writes. Both are required when the declaration is
   present. The path must be clean and relative and stay inside the
   repository; an unknown sub-key, an empty value or an unsafe path is a
   configuration error naming the key.
2. **The owner builds the item binary.** When the owner's configuration
   declares an item build, the owner runs the build command in the item
   worktree before each Implement, archive and review step of that item,
   through the executor that already runs Verification and declared
   regeneration commands, and keeps the build log under the artifact
   directory. The owner reads the declaration it loaded when it started; an
   item cannot change which command the owner runs for it.
3. **The item binary runs the item's steps.** After a successful build, the
   owner asks the built binary whether it would change the Run Database
   schema. When it would not, the owner runs that step with the built binary
   and writes one console-log line naming the binary and the step.
4. **A schema change keeps the owner's binary.** When the built binary reports
   that it would change the Run Database schema, or answers with any failure,
   the owner runs the step with its own binary, as before this Spec, and
   writes one console-log notice with the binary's answer (ADR-0225).
5. **A broken declaration parks.** A build that fails, or a declared binary
   that is missing or cannot start after the build, parks the item
   `delivery-error`, and the blocker names the build command or the path and
   the build log. A declared path that Git does not ignore parks the same way,
   because a built binary there would dirty the item worktree.
6. **An undeclared repository is unchanged.** Without the declaration the
   owner runs no build and no check and starts its own executable for every
   step, byte for byte as before.
7. **A migration check that writes nothing.** `roundfix migrate --check`
   reports whether this binary would change the Run Database schema: a
   current or absent database exits `0`; an older or newer one exits `2` with
   the remedy the other commands already name. It reads the database the way
   the read-only commands do, never migrates or writes it, and creates nothing
   when it is absent.
8. **The skill and the guides say so.** The Roundfix Skill's `deliver` and
   `setup` references, and the `deliver`, `migrate` and configuration guides,
   describe the declaration, the item binary, the schema rule, the park, and
   `migrate --check`.

## User Experience

The operator of an undeclared repository sees nothing new. In a declared
repository, the delivery console log gains one line per step naming the item
binary that ran it, or one notice explaining why the owner's binary ran
instead. A failed build appears as a `delivery-error` park in
`roundfix deliver status` with the build log path. `roundfix migrate --check`
prints one line, on standard output when the database is usable without a
schema change and on standard error when it is not.

## Non-Goals / Out of Scope

- Declaring the item build in Roundfix's own `.roundfixrc.yml`. An owner
  built before this Spec refuses the new key in the item worktree, which is
  the defect this Spec fixes, so the declaration lands after this Spec merges
  and the operator rebuilds the owner (Open Questions).
- Relaxing the strict refusal of unknown configuration keys, or letting an
  older binary read a newer configuration.
- Running the queue owner itself, its repository gate, its checks or its merge
  with the item binary.
- Building Roundfix for adopters: adopters install it from npm and declare
  nothing.
- Checking the build command's tools in `roundfix doctor`; Spec 0220 owns the
  Doctor's delivery readiness, and this declaration can join it later.
- Changing Park Class names, retry limits, the Run Database schema or the
  `owner-older-than-main` warning.

## Success Metrics

1. Success Metric: in a disposable repository that declares an item build,
   the owner's Implement, archive and review steps start the built binary in
   the item worktree with the arguments the owner's own binary received
   before; without the declaration they start the owner's executable and no
   build runs.
2. Success Metric: a built binary whose `migrate --check` exits non-zero is
   never started for a step; the step starts the owner's executable and the
   console log carries one notice with the binary's answer.
3. Success Metric: a build command that exits non-zero parks the item
   `delivery-error` with a blocker naming the command and the build log path.
4. Success Metric: `roundfix migrate --check` exits `0` on a current and an
   absent Run Database and `2` on an older and a newer one; the database file
   keeps its bytes and schema version, and no file other than SQLite's own
   `-wal` and `-shm` sidecars appears in the Roundfix Home.

## Acceptance evidence

The outside-evidence row rests on records this Spec did not produce:

- The operator's intervention log of 2026-10-02, entries 105 to 107: the 0215
  park `delivery-error: owner binary (main) refuses verification.tools in the
  item's .roundfixrc.yml` and each retry made "with the item-branch binary".
  The operator's memory of the 2026-09-29 cycle records the same lesson as
  the owner's binary.
- A reproduction during authoring on 2026-10-03: the released 0.31.0 binary,
  run in a disposable repository whose Project Config declares
  `delivery.item_binary`, refused both `roundfix archive` and
  `roundfix implement` with `field item_binary not found in type
  config.deliveryOverlay` and exit `2`.
- A measurement during authoring on 2026-10-03: `make build` in this
  repository took 1.22 s and, after a source touch, 0.61 s on a warm Go build
  cache.
- The Rust compiler's bootstrap guide
  (<https://rustc-dev-guide.rust-lang.org/building/bootstrapping/what-bootstrapping-does.html>),
  a published account of a tool built from its own tree: "Bootstrapping is
  the process of using a compiler to compile itself", and the stage1 compiler
  is the one "from current code, by an earlier compiler".
- SQLite's file-format reference (<https://sqlite.org/fileformat.html>): the
  integer at offset 60 "is the user version which is set and queried by the
  user_version pragma", which is the schema version Roundfix records and
  `migrate --check` reads.
- The Secondbrain holds no entry on this defect; its Roundfix mirror and
  pending Inbox Entries were searched on 2026-10-03.

## Decisions

- The item binary runs the item's own steps, decided over a configuration
  compatibility rule and over re-executing the whole owner. See ADR-0225.
- The declaration is a repository policy in Project Config, like the
  Verification, bootstrap and derived-path commands the owner already runs
  from there, and the owner reads the configuration it loaded at start.
- A schema difference keeps today's behavior rather than parking, so a Spec
  that adds a Run Database migration is delivered as it is today.
- A failed build parks rather than falling back, because a declared build
  that fails on a committed tree is a defect the operator must see.

## Open Questions

- When does Roundfix's own Project Config declare the build? Default: after
  this Spec merges, the operator rebuilds `bin/roundfix` from `main`, restarts
  any owner, and lands a one-line Project Config change declaring `make build`
  and `bin/roundfix` under the standing authorization for `.roundfixrc.yml`.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.

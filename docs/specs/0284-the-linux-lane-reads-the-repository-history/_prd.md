---
spec: 0284-the-linux-lane-reads-the-repository-history
status: active
created: 2026-10-19
surfaces: [backend, cli, docs]
---

# The Linux Lane reads the repository history

Minimal PRD for a bug fix, minted by `write-techspec`. It carries the contract
the downstream skills parse; the TechSpec holds the design.

Spec 0282 left one unproven step: the operator runs `roundfix deliver lane
--all` on `main` before turning `delivery.linux_lane: true` back on. That run
failed on `main` at `dc7e5f72` on 2026-10-19 at 21:44 BRT with `Linux Lane
change failure: go test exited 1`. Every package passed except two tests that
read historical commits by object name, `TestOutputsForCommand` in
`internal/baseline` and `TestAuthorizationReaderResolvesPreservedHistoricalRecords`
in `internal/spec`. The lane records the candidate in a fresh repository with
no history, while CI checks out with `fetch-depth: 0`. The change was correct,
and the lane called it a Lane Change Failure.

After this Spec, the lane mounts the repository's object store read-only and
registers it as an alternate of the container repository, in the Delivery
Queue gate and in `roundfix deliver lane` alike. A test that reads a
historical object behaves as in CI. A store the lane cannot mount whole skips
the lane with its reason.

## Sources

One Backlog Entry, written and adopted by this Spec's planning change into
`references/` (see `references/_index.md`):

1. [`references/2026-10-19-the-linux-lane-cannot-read-the-repository-history.md`](references/2026-10-19-the-linux-lane-cannot-read-the-repository-history.md),
   type `fix`. It records the failed run, the two failing tests, the cause in
   `linuxLaneScript` and the operator's probe that made both packages pass.

No open Backlog Entry or unresolved Finding shares this context: the only
other active Spec, 0283, owns Git Locator reachability and the storage
report.

## Prerequisites

None. This Spec is authored against `dc7e5f72` (Spec 0282 merged) and does
not depend on Spec 0283, which may sit in the queue meanwhile and changes
neither the Linux Lane nor the `deliver` guide. Whichever merges first may
raise the Roundfix skill version and re-record the Behavior Surface Record, so
this plan is rebased and checked again before its Run.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the container path `/repo-objects` and the skip reasons are fixed
  strings. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. The lane keeps `--network none` and
  `--pull never`, and every test uses compiled Docker and Go stand-ins on
  temporary repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the lane, its tests and its
  documentation follow recorded decisions. ADR-0297 (this Spec): "The Linux
  Lane reads the repository history read-only". ADR-0294: "The Linux Lane
  runs like CI", which ADR-0297 amends. ADR-0288: "The Delivery Queue gate
  runs the changed Go packages in a Linux Lane", which ADR-0297 amends.
  ADR-0125: "A spawned fixture is a compiled binary, never a written script".
  ADR-0259: "A test runs in parallel unless it names why it cannot".
  ADR-0271: "A test bounds by real time only what it asserts". ADR-0285:
  "Every test file compiles for Windows, and a repository contract proves
  it". ADR-0189: "An owned skill's version names its content, and its minimum
  is the bundle". ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0280: "A guide Task records its own surface
  fingerprint, a PRD records who decided, and the glossary is held to its own
  rule". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0240:
  "One QA partial policy, and rows a Run sandbox cannot reach". ADR-0276: "A
  row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "A Spec proves what it claims,
  and every active Spec is held to it". ADR-0183: "An attributed claim
  carries a receipt the check proves". Set no obligation here, because this
  Spec changes none of their rules: ADR-0092 (triage intent lives in a typed
  backlog); ADR-0120 (retired documentation lives under one history root);
  ADR-0165 (a blocking review after archive parks publication for a
  corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a
  qualifying partial); ADR-0169 (the pre-PR review diffs the candidate from
  its merge base); ADR-0184 (a TechSpec states a command surface as a
  transcript); ADR-0196 (a pre-PR review finding parks only after
  validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds);
  ADR-0229 (an operator archive resumes any park from the Run start);
  ADR-0233 (a skill version raise is regenerated at merge, and a review-only
  correction returns to review); ADR-0237 (a Delivery Retry records a merge
  made outside the queue); ADR-0247 (an archive leaves an Archive Record, and
  the Spec folder stays in Git); ADR-0248 (existing history is sanitized in
  batches after a history-full tag); ADR-0250 (a module version is chosen
  when it is recorded, and the coverage record lists every platform);
  ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps
  its verdict); ADR-0254 (baseline update sanitizes pending history in the
  change it plans); ADR-0258 (a Spec is authored against the QA rerun
  classes, and a retirement writes reduced history); ADR-0260 (an archive
  relinks live citations, and a Baseline update lists the citations it
  breaks); ADR-0261 (an Archive Record names its pull request, and readers
  fetch the pull request head); ADR-0262 (a cancelled Spec leaves an Archive
  Record and returns its sources); ADR-0269 (archive and cancel refuse
  ignored Spec files and finish a failed removal forward); ADR-0272 (a review
  record survives a fix, and a correction inside the archived Spec's declared
  paths returns to review); ADR-0277 (a Legacy Archive Folder conversion
  keeps the summary, links and dates its folder held); ADR-0278 (a retirement
  deletes, and the History Root holds only Archive Records); ADR-0279 (status
  reports a Queue Owner proven gone, and a dismissal survives a fix that
  leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds
  instead of writing into it); ADR-0282 (the history cleanup names a
  superseded ADR in its successor, and a large candidate is reviewed in
  segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of
  its fix); ADR-0290 (merged items give back their Item Binaries, and Run
  Worktrees leave only by proof); ADR-0295 (reachability reads the
  fetched default branch, and origin names it when origin/HEAD is unset);
  ADR-0296 (the storage report reconciles every page of one Storage
  Snapshot). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AJ approval of 2026-10-19 through AskUserQuestion, option "Mexer em
  Governed Paths declarados" ("'Concedo' para os caminhos que cada Spec
  declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0284-the-linux-lane-reads-the-repository-history/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter, Baseline module or Project Config file change.

## Goals

- A test that reads a historical object passes in the Linux Lane as it does
  in CI.
- The container can read the history and never write it, and the working tree
  is never mounted.
- The gate and `roundfix deliver lane`, in changed-package and `--all` modes,
  build the same container.
- A store the lane cannot mount whole is named, never silently run without
  history.

## Core Features

1. The history mount. The lane resolves the object store Git reads in its
   working directory on the host and mounts it read-only at
   `/repo-objects`; the lane script registers it as an alternate right after
   `git init` and fails as a Lane Container Failure when it cannot read it.
2. The skips. A store path holding a colon, a store path that is not a
   directory, and a store that borrows objects through `info/alternates` each
   skip the lane with a named reason before Docker is called.
3. Documentation. The `deliver` guide's Linux Lane and `deliver lane`
   sections, the Roundfix `deliver` reference and **Linux Lane** in
   `CONTEXT.md` describe the mount and the skips.

## Non-Goals / Out of Scope

- Making refs (tags, the default branch) available in the container: no test
  in this repository needs one.
- Mounting the alternate stores a store borrows from.
- Turning `delivery.linux_lane` back on: the operator does it in a separate
  change after this Spec merges and an on-demand run passes.
- Changing the image, the init, the module cache mount, the package
  selection, the cap, the network or pull policy, the outcome classification
  or any blocker.

## Success Metrics

1. Success Metric: with compiled Docker and Go stand-ins, the Delivery Queue
   gate run from a linked item worktree, `roundfix deliver lane --all` and
   `roundfix deliver lane` each start `docker run` with exactly two `-v`
   values, the module cache and `<objects>:/repo-objects:ro`, where
   `<objects>` is the main repository's object store, and with a script that
   writes `/repo-objects` into `.git/objects/info/alternates` after `git
   init`.
2. Success Metric: a store whose `info/alternates` holds an entry, a store
   path holding a colon and a store path that is not a directory each skip
   with their reason and start no Docker command; `roundfix deliver lane`
   exits `3` with the notice.
3. Success Metric: a real lane run of the audited head passes with both
   history tests passing, when the QA gate can reach Docker; otherwise the
   operator runs it after merge.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-19.

- This repository's CI workflow, which no Task of this Spec authors:
  `.github/workflows/ci-verify.yml` checks out with `fetch-depth: 0` under the
  comment "Historical authorization contracts read retained main ancestors."
  The QA gate reads the file at the audited head.
- The operator's lane run on `main` at `dc7e5f72` (2026-10-19, 21:44 BRT), a
  measurement this Spec did not design, recorded in the adopted Backlog Entry:
  only the two history tests failed.
- Git, "gitrepository-layout" (https://git-scm.com/docs/gitrepository-layout),
  `objects`: "You could be using the `objects/info/alternates` or
  `$GIT_ALTERNATE_OBJECT_DIRECTORIES` mechanisms to borrow objects from other
  object stores."
- Git, "git-clone" (https://git-scm.com/docs/git-clone), `--shared`: "you can
  simply run `git` `repack` `-a` to copy all objects from the source
  repository into a pack in the cloned repository."
- Docker Docs, "Bind mounts"
  (https://docs.docker.com/engine/storage/bind-mounts/): "You can use the
  `readonly` or `ro` option to prevent the container from writing to the
  mount." and "If you use `--volume` to bind-mount a file or directory that
  does not yet exist on the Docker host, Docker automatically creates the
  directory on the host for you."

## Unreachable Acceptance

- criterion: a real Linux Lane run, `roundfix deliver lane --all`, of the
  audited head in a locally present `golang:1.26` image passes with
  `TestOutputsForCommand` and
  `TestAuthorizationReaderResolvesPreservedHistoricalRecords` passing.
  reason: no Verification may need Docker, and a Run sandbox may not reach
  the Docker daemon; the QA gate runs it only when it can.
  satisfied-by: the operator runs `roundfix deliver lane --all` on `main`
  after this Spec merges and before the change that sets
  `delivery.linux_lane: true`.

## Validation

The decisions were checked on 2026-10-19 against these sources.

- Secondbrain, read only. Read `wiki/index.md`, then ran `qmd query "linux
  lane docker git history alternates object store read-only mount" --all
  --files --min-score 0.3` and `qmd query "git alternates objects info
  alternates container bind mount history fetch-depth" --all --files
  --min-score 0.3`. They returned the mirrors of this repository's Spec 0282,
  ADR-0288, ADR-0294, the `deliver` guide and reference, and an unrelated
  article on Git; none holds a decision on the object store or alternates. No
  Inbox Entry was written, because the briefing forbids writing there.
- Exa: the Git and Docker passages quoted above. Challenging: the
  gitrepository-layout page says an alternates path is relative "to the
  object database, not to the repository", so a relative entry resolves
  against wherever the store is mounted, one reason the lane does not mount
  borrowed stores.
- Probes on 2026-10-19 with Docker 29.6.1 and `golang:1.26`
  (`sha256:37cfcca63222ec499bec5e8ad7aad52dbf7529db36d5812b8eae2f0621aa92d6`),
  `--init --pull never --network none --user 65534:65534`, on a scratch
  module whose test runs `git show <first commit>:old.txt` after a later
  commit removed the file; the TechSpec records each.

## Glossary

- changes: **Linux Lane**

## Decisions

- 2026-10-19 (about 16:20 BRT): the cycle AJ authorizations — "Religar
  delivery.linux_lane: true", "docker pull/build da imagem", "Mexer em
  Governed Paths declarados" (option text "'Concedo' para os caminhos que cada
  Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
  squash + release + push"; failures — "Regras AE–AH e parar na proposta
  (Recomendado)"; authoring — "Opus + strict check, enfileirar sem esperar
  você (Recomendado)".
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present; decisions taken from the operator's cycle AJ patch briefing for Spec 0284 of 2026-10-19 and the probes characterized in the TechSpec.

Taken from those sources, with no further interview:

- The lane mounts the object store read-only and registers it as an
  alternate; the candidate is still recorded fresh. See ADR-0297.
- A store that borrows through `info/alternates`, holds a colon or is not a
  directory skips the lane with its reason. See ADR-0297.
- No refs are copied, because the failed run needed only objects. See
  ADR-0297.

## Open Questions

None.

---
spec: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
status: active
created: 2026-10-18
surfaces: [backend, cli, docs]
---

# Linux tests run before publication, and delivered work gives back its storage

Three gaps of cycle AH, observed by the operator on 2026-10-17 and 2026-10-18,
each needed a hand repair. A defect that exists only on Linux reached GitHub
Actions before anything on the maintainer's macOS machine could see it. Three
Runs of a delivered Spec stayed preserved as `dirty` without saying why, and the
maintainer had them removed by hand. And Roundfix Home grows without a stated
policy, while the maintainer's own rule forbids cleaning it until a Spec states
one. After this Spec the Delivery Queue runs the changed Go packages' tests in a
Linux container before it publishes, when the repository opts in and Docker
holds the image; reconciliation releases a delivered Spec's leftovers that lie
beneath the paths its Tasks declared, and names the path or commit that keeps
any other; and merged items give back their Item Binaries under a written
retention policy.

## Sources

There are no open Backlog Entries or Findings; the maintainer named these
sources on 2026-10-18. Each is an operator observation, cited with its
evidence.

1. Linux-only failures reach CI only. Spec 0276's fresh item binary copy failed
   only in GitHub Actions on Linux with `fork/exec <copy>: text file busy`
   (CI run 39053329570, pull request #568). Its proof,
   `internal/cli/deliver_item_binary_fork_linux_test.go`, carries
   `//go:build linux`, so no Run on the maintainer's machine compiles it.
   Windows test compilation was broken in seven files for weeks until Spec 0277
   found it (ADR-0285). The operator proved the fix by hand in a `golang:1.26`
   container: a copy of the tree, a fresh `git init`, and `go test` of the
   package, passing with the fix and failing with it reverted.
2. Dirty Runs of a delivered Spec stay preserved. After Spec 0272 merged as
   pull request #555, `roundfix reconcile` released three of its Runs as
   `superseded` and kept `run_20261017T022748Z`, `run_20261017T032927Z` and
   `run_20261017T040614Z` as `dirty` with only "Run Worktree has tracked or
   untracked changes". The maintainer had them removed by hand on 2026-10-18.
3. The store grows without a policy. Measured at 13:30 BRT on 2026-10-18,
   `~/.roundfix/worktrees` held 4.9 GB after that cleanup and
   `~/.roundfix/artifacts` 412 MB, and `roundfix reconcile` listed 309 released
   Runs. The maintainer's standing rule of 2026-08-10 forbids `gc` or artifact
   cleanup until a Spec states the policy, while Run Worktrees of terminal Runs
   may be released.

## Prerequisites

None. Spec 0278 merged as pull request #575 and v0.89.0 is released from the
same tree; this Spec releases as v0.90.0.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; Project Config gains one boolean key, a container
  name is derived per gate, and the reconciliation reason gains a suffix.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the Linux Lane starts a container that
  must reach no network and receive no credential: it runs with
  `--network none` and `--pull never`, and only the variables `go test` needs
  enter it. No request path, endpoint or recipient is added, and every test
  uses a fake `docker` built from Go source. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0288 (this Spec): "The Delivery Queue
  gate runs the changed Go packages in a Linux Lane". ADR-0289 (this Spec): "A
  delivered Spec's dirty Run is superseded under its declared directories".
  ADR-0290 (this Spec): "Merged items give back their Item Binaries". ADR-0212: "A merged Spec supersedes its Runs' Spec-directory work and
  declared leftovers", which ADR-0289 amends. ADR-0282: "this decision adds no
  reclamation for them", which ADR-0290 amends. ADR-0053: "Ambiguous, dirty, and unintegrated work
  stays intact", which the widened scope keeps. ADR-0115: "A superseded Run
  Branch is discarded through Roundfix". ADR-0033: "The Run Event Journal is
  pruned by a retention window". ADR-0255: "Run Retention removes terminal Runs whole".
  ADR-0266: "worktrees share one build cache". ADR-0225: "A queue item runs the
  Roundfix binary its branch builds". ADR-0281: "A test replaces an executable
  it rebuilds instead of writing into it". ADR-0285: "Every test file compiles
  for Windows". ADR-0273: "The judge's keys never enter a Run". ADR-0252: "The
  selective gate runs the repository contracts a change makes relevant".
  ADR-0253: "Every Repository Contract Test runs on main and before a release".
  ADR-0271: "A test bounds by real time only the budget it asserts". ADR-0189:
  "A version names one content". ADR-0280: "A guide Task records its own
  fingerprint". ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0276:
  "A row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "Each claim carries its proof as
  a section". ADR-0161: "A merged Spec releases its Runs on the merged head", which the
  widened scope extends to leftovers beneath declared paths. ADR-0232: "Merge
  evidence releases the Runs and item branches of a merged Spec", whose release
  now also removes the Item Binary directory. ADR-0171: "A retention prune
  reports only what it reclaimed", which the Item Binary sweep follows. Set no obligation here, because this Spec changes none of their rules:
  ADR-0170 (a Task already completed is nothing to carry); ADR-0233 (version
  regeneration at merge); ADR-0240 (the QA partial policy); ADR-0250 (module
  versions and the coverage record's platforms, which no Task moves); ADR-0278
  (retirement deletes); ADR-0283 (the Run Marker, which the container never
  receives); ADR-0286 (the Correction Review); ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0172 (doctor reports reclaimable Run storage from cheap reads); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0284 (carry-forward binds a Task to the planning text it references); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0196 (a pre-PR review finding parks only after validation). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH approval of 2026-10-16 through AskUserQuestion, option "Governed Paths
  declarados" ("\"Concedo\" como em AE/AF/AG para os Governed Paths que cada
  Spec declarar: módulos da Baseline, registros de versão, make
  baseline-digests, skills e testes governados. Makefile, workflows de CI,
  go.mod e configuração de lint ficam fora."), carried into cycle AI, and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/roundfix/references/storage.md`, `.roundfixrc.yml`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter or Baseline module change.

## Goals

- A candidate whose changed Go packages fail on Linux parks before any Pull
  Request exists when the repository opts in and the machine holds the image,
  and a machine that cannot run the lane says why and delivers as today.
- A dirty Run of a delivered Spec whose leftovers lie under the Spec's
  directories or beneath a path its Tasks declared or recorded is released by
  `roundfix reconcile --apply` and by the Delivery Queue's release after a
  merge.
- A dirty Run that stays preserved names the path or commit that keeps it.
- Roundfix Home has one written retention policy: what is kept, for how long,
  what never goes and who reclaims the rest, and a merged item's Item Binaries
  are reclaimed under it.

## User Stories

1. As the maintainer, I want the Delivery Queue to run the changed packages'
   tests on Linux before it pushes, so that a Linux-only defect parks the item
   on my machine instead of failing CI on the Pull Request.
2. As the maintainer on a machine without Docker or without the image, I want
   the lane to skip with a reason I can read, so that delivery never waits on
   or pulls anything I did not install.
3. As the operator, I want reconcile to release a delivered Spec's Runs whose
   leftovers lie beneath a directory its Tasks declared, so that I do not remove
   them by hand.
4. As the operator, I want a preserved dirty Run to say which path or commit
   keeps it, so that I can decide without opening the worktree.
5. As the maintainer, I want merged items' Item Binaries reclaimed and the
   retention policy written down, so that the store stops growing past what
   open work needs and I know what `gc` may remove.

## Core Features

1. Linux Lane. With `delivery.linux_lane: true` in Project Config, the Delivery
   Queue's repository gate, after the repository Verification passes, runs
   `go test -count=1` over the candidate's changed Go packages in a
   `golang:<go.mod go version>` container that reaches no network, pulls no
   image and receives only the variables `go test` needs. A failure parks the
   item `gate-failed`; a machine that cannot run the lane logs a named skip and
   the gate continues.
2. This repository opts in: its Project Config sets `delivery.linux_lane: true`.
3. Directory scope. Reconciliation treats a dirty path beneath a path the
   archived Spec's Tasks declared or recorded as in scope, as it treats the
   path itself.
4. Named dirty reason. A Run kept `dirty` while its Spec is archived at the
   proof head names the failing commit proof or the first path outside the
   scope.
5. Item Binary reclamation. The Delivery Queue removes a merged item's Item
   Binary directory when it releases the item's Runs, and the GC Command
   reclaims the Item Binary directories of items that are merged or no longer
   queued, previewed by `--dry-run`.
6. Retention policy. The GC Command guide states the policy for every area of
   Roundfix Home.
7. Documentation. The `deliver`, `reconcile` and `gc` user guides, the
   configuration guide, the Roundfix skill references that cover them, and the
   glossary describe each change.

## Non-Goals / Out of Scope

- Pulling, building or updating any image; running the lane at Task
  settlement, in a Run's Verification, in CI or for review Runs.
- A lane for any language other than Go, or for platforms other than Linux.
- Changing the Makefile, a CI workflow, `go.mod`, `go.sum`, lint or formatter
  configuration, or any Baseline module.
- Releasing a dirty Run by age, by content comparison, or by the plan revision
  it started from; changing the commit proof of ADR-0212.
- Reclaiming Run Worktrees, pre-PR review records, delivery logs, the Run
  Database or another repository's artifacts beyond what the existing commands
  already do; changing Journal Retention or Run Retention.
- Fixing `roundfix storage report`, which failed on 2026-10-18 with
  "reconcile Run Database storage groups" while a WAL was live; that is
  reported to the lead, not this Spec.
- The adopters fiscus, pantheon and tax-poc, which the maintainer routed to
  inbox notes.

## Success Metrics

1. Success Metric: with a fake `docker` that records its arguments and
   standard input, a gate over a candidate that changed one package runs
   exactly one `docker run` with `--pull never`, `--network none`, `--rm`, the
   read-only module cache mount, only the listed variables, the image
   `golang:<go.mod version>` and that package; a fake run exiting 1 fails the
   gate with the log path; each skip condition logs its named reason, makes no
   `docker run` and passes the gate; with `delivery.linux_lane` unset no
   `docker` command runs at all.
2. Success Metric: a dirty Run of an archived Spec whose leftover lies beneath
   a declared `deletes:` directory is `superseded` and released by apply, a
   leftover at a sibling path that merely shares the prefix stays `dirty`
   naming that path, and an unrepresented Task commit keeps it `dirty` naming
   the commit.
3. Success Metric: after a merged item's release, its `item-binary` directory
   is gone and its logs remain; `roundfix gc --dry-run` lists the Item Binary
   directories of merged and unqueued items with their bytes and keeps those of
   an item at any other stage, and on a repository without Item Binaries its
   output is byte-identical to today's.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-18.

- This repository's history, which no Task of this Spec authors. Commit
  `23ab0904` (pull request #568) added
  `internal/cli/deliver_item_binary_fork_linux_test.go`, whose first line is
  `//go:build linux`. The Archive Record of Spec 0272 names
  `source_revision: 10a9c7590d9362307370378d17dd5fd825af4b31`, where
  task_04 declares `deletes: docs/history/findings` and
  `deletes: docs/history/backlog`, and ADR-0282 at the Delivery Base ends its
  consequences with "this decision adds no reclamation for them". The QA gate
  reads all three with `git show`, offline.
- Docker documentation, "Running containers"
  (https://docs.docker.com/engine/containers/run/): "Exit code `125`
  indicates that the error is with Docker daemon itself." This is why exit
  `125` skips instead of failing the gate.
- Docker documentation, "docker container run"
  (https://docs.docker.com/reference/cli/docker/container/run/), on `--pull`:
  "Do not pull the image, even if it's missing, and produce an error if the
  image does not exist in the image cache." This is the `never` policy the
  lane uses.
- Linux manual, execve(2) (https://man7.org/linux/man-pages/man2/execve.2.html):
  "ETXTBSY The specified executable was open for writing by one or more
  processes." The Spec 0276 failure is a Linux kernel behavior a macOS host
  never shows.

## Validation

The decisions were checked on 2026-10-18 against these sources.

- Secondbrain. Read `wiki/index.md` and the pending inbox `inbox/roundfix/`,
  then ran three `qmd query` runs from the Secondbrain root:
  - "roundfix store retention worktrees artifacts gc policy" returned the
    mirrors of the Archive Record of Spec 0014, the Run Database lifecycle
    guide, the `gc` guide and ADR-0255;
  - "linux only test failure CI docker container before pull request" returned
    pantheon's deployment guide and two unrelated web notes;
  - "reconcile dirty run worktree superseded merged spec leftovers" returned
    the mirrors of ADR-0212, ADR-0115, the `reconcile` reference and the Archive
    Record of Spec 0227.

  All agree with the repository and hold no decision on a Linux container gate,
  directory-scoped leftovers or Item Binary reclamation. The two pending inbox
  entries of 2026-10-18 concern QA screenshots and a `secret-found` retry; they
  share no component with this Spec and stay apart.
- Exa. Supporting: the Docker exit status and pull policy pages and the
  execve(2) page quoted above. Challenging: the same `docker container run`
  page says the command runs "pulling the image if needed", the default this
  Spec overrides with `--pull never`. No research Inbox Entry was captured; the
  sources are standard references.
- Measurement, source 1. `docker image inspect golang:1.26` found no image on
  the maintainer's machine, so no container ran. Compiling the `internal/cli`
  and `internal/worktree` test binaries for Linux from an empty build cache
  took 12.9 s, and `go test -count=1` of those two packages on the host took
  156.9 s (3,750 tests), so the lane runs once per candidate (ADR-0288).
- Measurement, source 2. `roundfix runs list --all` shows the seven Runs of
  0272 and `roundfix events` their Task phases: the three kept Runs failed
  task_02, task_02 again, and task_04 after settling task_02 and task_03; the
  three released Runs failed task_06, task_04, and task_05 after settling
  task_04. The kept Runs' worktrees were gone, so their exact dirty paths were
  not observed (ADR-0289).
- Measurement, source 3. With every command read-only: `roundfix gc --dry-run`
  found 0 journal rows, 3 removable Runs with no reclaimable bytes and no
  Worktree Build Cache; `roundfix gc compact` previewed 12,226,560 reclaimable
  bytes of a 1,113,714,688-byte Run Database; `du` measured 334,766,080 bytes
  of Item Binary copies (two items, five copies each) in 412 MB of artifacts,
  and 3.4 GB of worktrees, all belonging to other repositories (ADR-0290).

## Glossary

- adds: **Linux Lane**
- adds: **Item Binary**
- changes: **Run Worktree Reconciliation**
- changes: **GC Command**

## Decisions

- 2026-10-16: which failure rules govern cycle AH — "Regras de AE–AG
  (Recomendado)".
- 2026-10-16: which authority covers the Governed Paths — "Governed Paths
  declarados", option text "\"Concedo\" como em AE/AF/AG para os Governed
  Paths que cada Spec declarar: módulos da Baseline, registros de versão, make
  baseline-digests, skills e testes governados. Makefile, workflows de CI,
  go.mod e configuração de lint ficam fora."
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".
- 2026-10-18 (about 13:30 BRT): which items the next cycle takes — "Verificação
  também no Linux, Runs sujos de Spec entregue, Crescimento da store, Adopters
  que ficaram fora".
- 2026-10-18 (about 13:30 BRT): what to do with the three dirty 0272 worktrees
  — "Remover agora (Recomendado)".
- 2026-10-18 (about 13:30 BRT): what to do with the Spec — "Escrever, entregar
  e lançar (Recomendado)".

Autonomous route: no person present; decisions taken from the lead's briefing of 2026-10-18 for Spec 0279 and the code and measurements characterized in the TechSpec.

Taken from those sources, with no further interview:

- The Linux Lane runs once per candidate in the Delivery Queue's gate, opt-in
  through Project Config, never pulls and skips with a named reason. See
  ADR-0288.
- This repository opts in through `.roundfixrc.yml`.
- A declared or recorded path covers the dirty paths beneath it, and a kept
  dirty Run names its failing proof. See ADR-0289.
- Merged items give back their Item Binaries, at the merge and through the GC
  Command, and Run Worktrees leave only by proof. See ADR-0290.

## Open Questions

None.

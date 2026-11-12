---
spec: 0384-the-makefile-follows-the-go-build-cache-set-with-go-env
status: active
created: 2026-11-12
surfaces: [infra, docs]
type: fix
---

# The Makefile follows the Go build cache set with go env

Reported by gesttione-drive on 2026-11-12: the Mac mini's internal disk, which must keep 80 GB free, had 66 GiB free, and this repository's Go build cache held 10 GB of it in `$TMPDIR/roundfix-go-build-501`, although on 2026-11-11 the mini had set its Go build cache on an external disk with `go env -w GOCACHE=/Volumes/Brain-SSD/cache/go-build`. Measured at `d79b30b4`, the Makefile chooses its cache only from a `GOCACHE` in the environment and otherwise uses the temporary directory, so the cache set with `go env` never reaches `make`. After this Spec, when the environment does not define `GOCACHE`, `make` uses the cache set with `go env -w GOCACHE` whenever it is an existing writable directory, under the same 10,240 MiB bound; when that disk is not mounted, `make` never creates the path, warns once, and falls back to the temporary cache under a 2,048 MiB bound; and with nothing set, as on the MacBook, nothing changes.

## Sources

Two sources, both in `references/` (see `references/_index.md`):

1. [`references/2026-11-12-gesttione-drive-the-roundfix-go-build-cache-fills-the-mini-internal-disk.md`](references/2026-11-12-gesttione-drive-the-roundfix-go-build-cache-fills-the-mini-internal-disk.md),
   ai-memory message `01a1d0ad-dd6f-7bb0-ad09-8a1e38911657` from
   gesttione-drive, 2026-11-12 16:53Z, nature measured, reported by
   gesttione-drive.
2. [`references/2026-11-12-the-writer-of-the-default-go-build-cache-on-the-mini.md`](references/2026-11-12-the-writer-of-the-default-go-build-cache-on-the-mini.md),
   the author's read-only measurement of who wrote 613 MB into
   `~/Library/Caches/go-build` on the night of 2026-11-11 to 2026-11-12.

At the authoring base `docs/backlog/` holds three open entries: two that
Spec 0385 adopts (a review anchor followed by a period, and the Delivery
Owner waiting for the Machine Verification Capacity) and the Conexus symptoms
waiting for reproduction steps; `docs/findings/` holds no entry. None is
about a build cache.

## Prerequisites

None in `requires`. Authored against `d79b30b4` (v0.142.0 and the Backlog
Entry merge #839). The queue delivers Spec 0383 before this Spec and Spec
0385 after it. This Spec owns the `GOCACHE` block of the `Makefile`, the
Makefile cache rule of `docs/agents/specific-repository.md`, ADR-0370's
amendment and ADR-0394. It shares with them:

- `CONTEXT.md`: 0383 and 0385 may add or change their own terms; task_02
  changes only the **Build Cache Bound** entry.
- `docs/agents/setup-context.json`: task_01 regenerates it with the
  repository's command as its last step, never by hand, so whatever an
  earlier Spec merged into it only moves the bytes that command writes.

No Task of this Spec changes a skill, a Baseline module or Go production
code, so the owned-skill version records the other two Specs raise are not
shared.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted or changed; the cache path is the one the user set or the
  one the Makefile already chose. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — nothing reaches the network;
  `make` runs `go env` locally. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0394 (new) records the decision:
  "The Makefile follows the Go build cache set with go env, and falls back
  small". ADR-0370 is amended: its Build Cache Bound now also holds the cache
  set with `go env`, and its Considered Option "Moving the cache to another
  volume. Left to each machine" is taken for that cache; a non-QA Task writes
  the notice. ADR-0327 applies unchanged: the bound removes only entries idle
  for at least two hours, oldest first, and never clears a cache, which now
  matters because the cache set with `go env` is shared by every Go project
  of the user. ADR-0392 applies unchanged: the manifest regeneration task_01
  runs records the Recorded Roundfix Version as that command decides.
  ADR-0334 does not apply: the trims at Run start and at the queue gate keep
  reading `worktree.shared_caches`, which this Spec does not change. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-11-12 ("Autorizo; sem SSD, TMPDIR com 2 GB") for the `Makefile`; the
  derived `docs/agents/setup-context.json`, which records the `Makefile`
  digest, through its sanctioned command; and the repository-owned guide
  `docs/agents/specific-repository.md`, whose rule describes the Makefile's
  cache choice. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0384-the-makefile-follows-the-go-build-cache-set-with-go-env/_authorization.md`;
  bounded files: `Makefile`, `docs/agents/setup-context.json`,
  `docs/agents/specific-repository.md`.

## Goals

- On a machine whose Go build cache is set with `go env -w GOCACHE` on a
  mounted disk, every `make` target builds and tests in that cache and holds
  it to 10,240 MiB.
- When that disk is not mounted, `make` creates nothing under the missing
  path, says so in one line, and builds on the temporary cache under 2,048
  MiB.
- On a machine with nothing set, and for a caller that exports `GOCACHE`,
  `make` behaves exactly as it does at `d79b30b4`.

## User Stories

1. As the maintainer of the Mac mini, I want this repository's builds to use
   the external cache I set once with `go env -w`, so that the internal disk
   keeps its free space.
2. As the maintainer, I want a build that runs while the external disk is
   unmounted to warn me and stay small, so that the cache never again grows
   unnoticed on the internal disk.
3. As a developer on the MacBook, which shares the shell startup files and
   has no external disk, I want nothing to change, so that my builds keep
   their cache and bound.

## Core Features

1. When the environment does not define `GOCACHE`, `make` asks
   `go env -changed GOCACHE`; when it prints a value and `go env GOCACHE` is
   an absolute path to an existing directory in which `make` can create and
   remove a probe file, `make` exports that path and defaults
   `GOCACHE_BOUND_MB` to 10240.
2. When the set cache is unusable, `make` never creates it, exports the cache
   it chose at `d79b30b4` (`$TMPDIR/roundfix-go-build-<uid>` or
   `$(CURDIR)/.gocache`), defaults `GOCACHE_BOUND_MB` to 2048, and prints one
   warning line to standard error.
3. A `GOCACHE` in the environment or on the command line, nothing set,
   `GOCACHE_BOUND_MB`, `GOCACHE_BOUND_IDLE` and sub-makes behave as at
   `d79b30b4`.
4. The repository's guide, README, glossary and ADR-0370 describe the new
   choice, and ADR-0394 records it.

## User Experience

On the mini, with the external disk mounted, `make build` prints what it
printed before and builds in `/Volumes/Brain-SSD/cache/go-build`; the
`gocache-bound` step holds that cache to 10,240 MiB. With the disk unmounted,
every top-level `make` begins with
`gocache: warning: GOCACHE set with go env (/Volumes/Brain-SSD/cache/go-build) is not a writable directory; using <TMPDIR>/roundfix-go-build-501 with a 2048 MiB default bound`
and builds there. On the MacBook nothing in the output changes.

## Non-Goals / Out of Scope

- Deleting the old caches on the mini. The maintainer decided that the
  operator deletes `$TMPDIR/roundfix-go-build-501` and
  `~/Library/Caches/go-build` after delivery, in a window with no build; no
  Task deletes anything outside its temporary directories.
- The operator's `~/.roundfix-operator/release-checks.sh`, which exports
  `GOCACHE=$HOME/Library/Caches/go-build` and wrote the 613 MB
  (`references/2026-11-12-the-writer-of-the-default-go-build-cache-on-the-mini.md`).
  It is outside this repository; the operator removes that export. An
  exported `GOCACHE` stays a caller's choice that `make` leaves unbounded.
- A `worktree.shared_caches` placeholder for the cache set with `go env`, so
  that the GC Command, Run starts, gates and Doctor's free-space line cover
  it (ADR-0394 defers it); `.roundfixrc.yml` is unchanged.
- `go` commands that a user or an Agent runs outside `make`: the go command
  already uses the cache set with `go env` there.
- Any Go production code, skill, Baseline module, workflow, `go.mod`,
  `go.sum`, lint or formatter configuration.

## Success Metrics

1. Success Metric: with a temporary Go environment file that sets an
   existing writable directory, `make test` passes that directory to `go`
   and bounds it with `-max-mb 10240`; on `d79b30b4` it passes the
   temporary cache.
2. Success Metric: with the set path missing under an existing parent, a
   read-only directory or a regular file, `make test` passes the temporary
   cache, bounds it with `-max-mb 2048`, prints exactly the warning line, and
   leaves the missing path absent.
3. Success Metric: with nothing set, or with `GOCACHE` in the environment,
   `make test` passes and bounds exactly what it does on `d79b30b4`.
4. Success Metric: every existing `make` test of `./internal/verifyselect`
   passes with only the declared stub changes.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not write.

- The adopted message, measured by gesttione-drive on the mini on
  2026-11-12 13:45–14:00 BRT: "du $TMPDIR/roundfix-go-build-501 = 10G",
  "go env GOCACHE = /Volumes/Brain-SSD/cache/go-build (GOENV = ~/Library/Application Support/go/env)",
  "df: 66 GiB livres".
- The go command's documentation, `go help env` in Go 1.26 and
  https://pkg.go.dev/cmd/go (read through Exa on 2026-11-12): "The -changed
  flag prints only those settings whose effective value differs from the
  default value that would be obtained in an empty environment with no prior
  uses of the -w flag." and, for `GOCACHE`, "The directory where the go
  command will store cached information for reuse in future builds. Must be
  an absolute path."
- The Secondbrain digest
  `wiki/sources/digest-roundfix-pid-reutilizado-e-cache-go-no-sandbox-2026-10-12.md`,
  recorded on 2026-10-12 for another Spec: "`codex sandbox -P :workspace`
  permitiu escrever no cwd, em `$TMPDIR` e em `/tmp`, e **negou**
  `~/Library/Caches/go-build`".

## Unreachable Acceptance

- criterion: on the Mac mini, after the merge and after the operator removes
  the old caches, a `make build` in the main checkout writes new entries
  under the cache set with `go env` and none under
  `$TMPDIR/roundfix-go-build-501` or `~/Library/Caches/go-build`.
  reason: it needs the mini's own Go environment file and its mounted
  external disk, which no Verification or QA row may read or write, and the
  deletion is the operator's action after delivery.
  satisfied-by: the operator's run of `make build` in the main checkout of
  the Mac mini after the merge, with `find` over the three directories for
  entries newer than a marker file and `du` of each, recorded session
  evidence.

## Validation

The decisions were checked on 2026-11-12 against these sources.

- Reading at `d79b30b4`: the `Makefile` `GOCACHE` block (lines 11–22) and
  `gocache-bound` target; ADR-0370; `docs/agents/specific-repository.md`;
  `README.md`; `CONTEXT.md` → **Build Cache Bound**; `.roundfixrc.yml`
  `worktree.shared_caches`; `internal/sharedcache`; every test that runs the
  `Makefile`. No production code sets `GOCACHE` or `GOENV` except the Linux
  Lane's container (`internal/cli/deliver_linux_lane.go`).
- Characterization with a printing makefile beside the real one, in
  temporary directories: with a temporary Go environment file that sets an
  existing directory, `make` chose `$TMPDIR/roundfix-go-build-<uid>` with
  bound 10240; with `GOCACHE` exported, it kept it with bound 0.
- `go env` on Go 1.26.5: `-changed GOCACHE` prints `GOCACHE='<path>'` when
  set and nothing when unset, set to the default, or with `GOENV=off`; a
  relative value reads as `off`; `go env GOCACHE` prints the bare path; an
  empty `GOCACHE` in the environment counts as unset; ten calls took 63 ms.
  A build whose configured cache is read-only fails with
  `permission denied`, which is why `make` probes before choosing.
- A prototype on a `git archive` of `d79b30b4`: the new test fails five of
  its seven subtests there (`nothing set` and `environment wins` pass, the
  kept paths) and passes with the change; each sabotage in the TechSpec fails
  its named subtest; with the Makefile change, six existing tests failed only
  because their `go` stubs recorded the new `env` queries, and pass with the
  declared stub changes; the whole `./internal/verifyselect` package passes;
  `TestTheSetupManifestIsCurrent` fails after the `Makefile` edit until the
  manifest command runs, whose second run reports `File changes: 0`; and the
  `./internal/docscontract` package passes with the docs edits of task_02.
- The writer of `~/Library/Caches/go-build`, in
  `references/2026-11-12-the-writer-of-the-default-go-build-cache-on-the-mini.md`.
- Published sources: the go command documentation quoted above, read through
  Exa.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix Go build cache GOCACHE go env Brain-SSD Mac mini disk" --all --files --min-score 0.3`,
  which returned `wiki/concepts/nobreak-e-armazenamento-do-brain-macmini.md`
  (the plan to keep regenerable caches on the external disk via `GOCACHE`),
  the 2026-10-12 sandbox digest quoted above, this repository's mirrored
  ADR-0370, ADR-0327 and ADR-0334 and Spec 0360 and 0321 records, and
  Pantheon's 2026-11-11 reference on developing on the mini; none
  contradicts this design.

## Glossary

- changes: **Build Cache Bound**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BD brief for Spec 0384 of 2026-11-12, the adopted message, the measurement in `references/` and the maintainer's answers below.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-12 ~17:05–17:20 BRT), verbatim:

- Question "Três Specs novas saem das mensagens de hoje: A = aceitar o
  skills-lock.json versão 2 (Vortex e Pantheon travados no update), B = o
  dono da fila não morre mais ao esperar a Machine Verification Capacity +
  âncora 'path:line.' aceita, C = cache Go do roundfix fora do disco interno
  do mini (Drive). Que autorizações e ordem?" → "A → C → B, mesmas do BC
  (Recommended)" — option text: "Push, merge e release logo após cada merge,
  como no ciclo BC. A primeiro porque trava dois adotantes; C antes de B
  porque o disco do mini está abaixo dos 80 GB." (C is this Spec.)
- First question "Spec C muda o Makefile (autoridade de ferramenta). Quando o
  GOCACHE do `go env` (mini: /Volumes/Brain-SSD/cache/go-build) aponta para um
  disco não montado, o que o make faz?" → free answer "Não é possível usar o
  cache no SSD?"; clarified by the second question "Sim, o cache vai para o
  SSD: no mini o make passa a usar /Volumes/Brain-SSD/cache/go-build, com
  limite de 10 GB. Falta decidir só o caso raro de o SSD não estar montado.
  Autoriza a mudança no Makefile, e o que o make faz nesse caso?" →
  "Autorizo; sem SSD, TMPDIR com 2 GB" — option text: "O build segue no disco
  interno com limite pequeno e avisa."
- Question "Depois que a Spec C trocar o cache para o SSD, apago os ~10,6 GB
  antigos do disco interno ($TMPDIR/roundfix-go-build-501 e
  ~/Library/Caches/go-build), numa janela sem build do Roundfix nem do
  Quaestor?" → "Sim, depois da troca (Recommended)" — option text: "Confiro
  que nenhum go build/test roda (ps) e apago só esses dois diretórios; é
  cache regenerável."

Taken from those and the evidence, with no further interview:

- `make` reads the go command's own setting with `go env -changed GOCACHE`
  and never names a machine's path; the request's alternative, testing for
  `/Volumes/Brain-SSD/cache`, is rejected in ADR-0394.
- Usable means absolute, an existing directory and writable through a probe
  file `make` removes at once; `make` never creates the set path or a parent,
  so an unmounted disk, or an empty mount-point directory under `/Volumes`,
  never becomes a cache on the internal disk.
- The set cache counts as one `make` chooses, so its bound defaults to
  10,240 MiB as the maintainer answered; the fallback's default is 2,048 MiB
  with one warning line; an exported `GOCACHE` stays unbounded unless
  `GOCACHE_BOUND_MB` is set (ADR-0370 unchanged there).
- The writer of `~/Library/Caches/go-build` is the operator's
  `release-checks.sh`, outside this repository: the operator removes its
  `GOCACHE` export, and no Task changes it.
- After delivery the operator deletes `$TMPDIR/roundfix-go-build-501` and
  `~/Library/Caches/go-build` on the mini in a window with no Go build or
  test, as the maintainer answered; no Task deletes anything.
- ADR-0394 records the decision and amends ADR-0370; a non-QA Task writes
  the notice into ADR-0370.
- Two implementation Tasks and the QA gate; `type: fix`, Release Ceiling
  `patch`. No released code changes: the binary, the skills and the Baseline
  are untouched. `roundfix release plan` still classifies the `Makefile`,
  `README.md`, `CONTEXT.md` and `docs/agents/` as non-maintenance paths
  (`MaintenanceOnly` in `internal/releaseplan/classifier.go`), so it proposes
  a patch; whether to cut it is the operator's call under that ceiling.

## Open Questions

None.

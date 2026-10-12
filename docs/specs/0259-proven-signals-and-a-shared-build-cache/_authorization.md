---
status: approved
granted: 2026-10-11
action: send every signal to a recorded process only after proving its start identity, forbid every other signal path with a repository rule, give every waiting test helper a lifeline and have the suite guard report and kill Helper Residue, make make share one Go build cache under the temporary directory, and let doctor report and gc reclaim the declared Worktree Build Caches of ended worktrees, with the guides, the glossary and the Roundfix skill describing them
consuming: 0259-proven-signals-and-a-shared-build-cache
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/storage.md
  - .roundfixrc.yml
  - Makefile
  - docs/agents/specific-repository.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0259

On 2026-10-11 the maintainer settled cycle AE. For this Spec, item E, the
decisions are: fix the class of F6 rather than its case, so that every signal
sent to a recorded PID first verifies the process's identity and helper
processes are joined or killed in `t.Cleanup`; accept F6 as fixed when it does
not reappear in N full suites, without naming the sender; share or reclaim the
per-worktree Go build cache after measuring whether the Run sandbox can write
outside the worktree; report cache bytes in the `storage` line of
`roundfix doctor` and reclaim them in `roundfix gc`.

The same day the maintainer answered "Grants dos caminhos governados" with
"Concedo" for the Governed Paths each Spec declares, and extended that grant to
the Makefile and the CI workflows for this Spec only. This design needs the
Makefile and no CI workflow. On 2026-09-30 the maintainer had granted the
skills ("considere autorizado a ajustar todas as skills se necessário") and the
Baseline source, guides and `.roundfixrc.yml` ("Autorizar os dois").

The governed set was measured on the authoring branch at `7a8b7f70` with
`GovernedPath`, through a `go test -overlay` probe in `internal/speccheck`
that wrote nothing to the repository, against every file the Tasks declare.
Seven of them are governed.

## Why each governed path is unavoidable

- `Makefile` holds the `GOCACHE` default that task_03 changes from the
  worktree's `.gocache` to the shared cache with that fallback.
- `docs/agents/specific-repository.md` gains, outside every setup marker, the
  rule for signals and helper lifelines (task_02) and the shared cache line
  (task_03), so whoever writes the next test or build meets them.
- `.roundfixrc.yml` declares `.gocache` under `worktree.caches` (task_04), so
  this repository's `doctor` and `gc` see its Worktree Build Caches.
- `.agents/skills/roundfix/references/storage.md` describes `roundfix gc` and
  the storage surfaces, and `.agents/skills/roundfix/references/setup.md`
  describes `roundfix doctor`; both must describe the new lines (task_04).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  skill's two version fields, raised with its content.

## What is not governed

The Go sources and tests the Tasks declare under `internal/store`,
`internal/cli`, `internal/daemon`, `internal/suiteguard`,
`internal/testfixture`, `internal/verifyselect` and `internal/config`; the
skill mirrors under `skills/roundfix/references/`;
`skills/testdata/owned-skill-versions.json`;
`docs/references/skill-coverage.json` and
`docs/references/behavior-surfaces.json`; `CONTEXT.md`; the user guides under
`docs/user-guide/`; and
`docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirror, the skill versions and
the Behavior Surface Record. This declaration records the regeneration that
follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
records the raised skill version, and
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
re-records the fingerprints of the `doctor` and `gc` surfaces. No Baseline
module changes, so no module version is recorded and `make baseline-digests`
rewrites nothing. No record or mirror is hand-edited.

## Limits

- No change to the CI workflows, `go.mod`, `go.sum`, the lint, formatter or
  test-runner configuration, any Baseline module, or the `### QA settlement`
  section of any skill.
- No change to the legacy PID-only Force Stop path or to
  `--owner-identity-unreadable`.
- No test, Verification command or QA row reads or writes the real
  `~/.roundfix`, signals a process it did not start or prove, or reaches the
  network; the QA gate's ten suites use the shared cache and temporary homes.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.

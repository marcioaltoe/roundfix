---
task: task_06
spec: 0151-a-supersession-the-archive-can-see
status: completed
type: backend
complexity: low
---

# Task 06: An archive state the deliverer check can recognise

## Overview

Two of this Spec's own decisions contradict each other, and review found the
seam. Reproduced before this Task was written:

```
$ roundfix supersede --spec 9999-sup-src  --by 9999-sup-deliverer --reason probe
superseded 9999-sup-src by 9999-sup-deliverer
$ roundfix archive 9999-sup-src
archived 9999-sup-src -> docs/history/specs/9999-sup-src
$ grep '^status:' docs/history/specs/9999-sup-src/_prd.md
status: active
$ roundfix supersede --spec 9999-sup-src2 --by 9999-sup-src --reason probe2
  superseding Spec "9999-sup-src" is not archived: archived _prd.md
  frontmatter status is "active"; expected "archived"
```

Core Feature 1 preserves the superseded Spec's PRD exactly, so the supersession
archive path deliberately does not stamp archive metadata. Core Feature 3 proves
the deliverer exists by reading a declared status. A Spec that went through the
first therefore fails the second: it is in the archive, and its preserved PRD
still says `active`.

The preserved PRD is not the defect — preserving it is the point. The check is
reading the wrong thing for that case.

## Requirements

1. MUST accept, as a deliverer, a Spec in the archive that carries a valid
   supersession record, whatever status its preserved PRD declares.
2. MUST keep accepting a Spec in the archive whose PRD declares `archived`.
3. MUST keep refusing a Spec in the archive that declares neither, and keep
   every refusal Tasks 01 and 05 deliver.
4. MUST NOT rewrite, stamp or otherwise modify a preserved PRD to make the check
   pass.

## Subtasks

- [ ] Recognise the supersession-archived state in the deliverer check.
- [ ] Add a test that archives through the supersession path and then uses that
      Spec as a deliverer.

## Acceptance Criteria

- [ ] A Spec archived through the supersession path is accepted as `--by`.
- [ ] Its preserved PRD is byte-identical after being used as a deliverer.
- [ ] A conventionally archived Spec is still accepted.
- [ ] A Spec in the archive with neither marker is still refused.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/supersede.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestSupersedeAcceptsASupersessionArchivedDeliverer" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSupersedeAcceptsASupersessionArchivedDeliverer"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The command

## Result

Implemented the archive-side deliverer recognition without changing archive
metadata or preserved PRDs. A conventionally archived Spec is still accepted
from its `status: archived` PRD; when an archived PRD declares another status,
the check now accepts only a valid `_supersession.md`. A missing or malformed
supersession remains a Preflight Validation refusal.

Focused evidence:

- Before the production edit,
  `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go test -count=1 ./internal/cli -run 'TestSupersedeAcceptsA'`
  reproduced the defect: the real supersede → archive → supersede flow exited 2
  because the archived deliverer's preserved PRD still declared `active`.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go test -count=1 -v ./internal/cli -run 'TestSupersede(AcceptsA|RejectsA|$)'`
  passed after the edit. `TestSupersedeAcceptsASupersessionArchivedDeliverer`
  proves the supersession-archive path is accepted and compares the archived
  PRD before and after use byte-for-byte.
- The same focused run passed the existing conventionally archived accepted
  path in `TestSupersede/accepted_path_writes_only_the_amendment`, all Tasks 01
  and 05 refusals, and the new
  `TestSupersedeRejectsAnUnmarkedArchivedDeliverer`, which proves an archived
  Spec with `status: active` and no supersession is refused without mutation.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache make verify-incremental`
  passed with the required process-table permission, including vet, all Go
  packages, skill checks, and the binary build. The initial sandboxed run
  reached the Go test sweep but two unrelated force-stop integration tests
  could not read the macOS process table (`operation not permitted`).
- `rtk git diff --check` passed.

The Task's declared Verification command was not run; the Daemon owns that
gate and Task status.

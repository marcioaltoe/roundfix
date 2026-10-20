---
type: fix
status: promoted
created: 2026-10-19
spec: 0284-the-linux-lane-reads-the-repository-history
reason: null
---

# The Linux Lane cannot read the repository history

## Symptom

Spec 0282 left one unproven step: the operator runs `roundfix deliver lane
--all` on `main` before turning `delivery.linux_lane: true` back on. That run
failed on `main` at `dc7e5f72` on 2026-10-19 at 21:44 BRT with `Linux Lane
change failure: go test exited 1`. Every package passed except two tests that
read historical commits:

- `internal/baseline`
  `TestOutputsForCommand/make_baseline-digests_matches_the_2026-08-06_enumeration`
  (`derived_ownership_test.go:61`): `read historical proof-cost grant
  81a6afb48f4a3683d0e5fad52f3919cf1bdfbbf4:docs/workflow/authorizations/2026-08-06-proof-cost.md:
  exit status 128`.
- `internal/spec` `TestAuthorizationReaderResolvesPreservedHistoricalRecords`
  (`authorization_test.go:300`): `list historical authorization records at
  6b8ea48725cbca13974eee0b400b3482202874f6: exit status 128: fatal: not a tree
  object`.

The change was correct; the lane reported a Lane Change Failure for an
environment gap, so the lane cannot be turned back on.

## Where

`internal/cli/deliver_linux_lane.go`: `linuxLaneScript` unpacks `git archive
<head>` and records it with a fresh `git init` and one commit, so no
historical object is reachable inside the container. CI checks out with
`fetch-depth: 0` (`.github/workflows/ci-verify.yml`), so the lane does not run
like CI. Both the Delivery Queue gate and `roundfix deliver lane` build the
container through `runLinuxLane`.

## Expected

The lane makes the candidate's repository history reachable read-only, so a
test that reads a historical object behaves as in CI. The candidate commit is
still recorded fresh; the object store is never writable and the working tree
is never mounted.

The operator proved the shape on 2026-10-19: the same archive, plus `-v
<git-common-dir>/objects:/repo-objects:ro` and `echo /repo-objects >
.git/objects/info/alternates` right after `git init`, made `go test -count=1
./internal/baseline ./internal/spec` pass in `golang:1.26` with `--init --pull
never --network none --user 65534:65534` and the module cache mounted
read-only as today.

## Evidence

The operator's lane run on `main` at `dc7e5f72` (2026-10-19, 21:44 BRT) and
the probe above, recorded in the cycle AJ operator briefing; no Finding.

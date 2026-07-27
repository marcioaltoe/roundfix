---
status: pending
created_at: 2026-07-27
updated_at: 2026-07-27
---

# Build hygiene — `go build ./cmd/roundfix` drops a 20 MiB binary at the repository root, and nothing ignores it (2026-07-27)

`.gitignore` ignores `/bin/`, which is where the Makefile writes its build
output. It does not ignore `/roundfix`, which is where `go build ./cmd/roundfix`
writes when invoked without `-o`. Agents run the bare form routinely as a
compile check, so the artifact appears unignored and untracked in a working tree
that Daemon commits stage by diff.

On 2026-07-27 that is exactly what happened: Spec 0038 Task 04's commit
`0d56f42` shipped a 20.1 MiB Mach-O arm64 executable at the repository root. The
Spec 0038 QA gate caught it as `RFQA-0038-01` and failed the gate. Neither
`make verify` nor CodeRabbit would have caught it — it is not a code error but a
scope-integrity problem, and only the QA gate's changed-path audit compares what
a Spec delivers against what it declared.

## Why the obvious fix needs authorization

Adding `/roundfix` to `.gitignore` closes the hole in one line. But
`docs/agents/agent-instructions.md:27` prohibits editing any "ignore file"
without express maintainer authorization, and no Spec's Tooling authority entry
covers `.gitignore`.

This finding records that boundary honestly: during the repair of
`RFQA-0038-01`, the supervisor removed the stray binary **and** added the ignore
entry. The next QA run flagged the ignore edit as an unauthorized protected
tooling mutation, and the entry was reverted, leaving only the binary removal —
which is all the QA finding actually required. The gate was right, and the
recurrence risk is still open.

## Suggested resolution

1. With maintainer authorization, add `/roundfix` to `.gitignore` with a comment
   naming its source, so a bare build cannot be swept into a commit again.
2. Consider whether the Daemon's commit staging should refuse to stage
   executable files outright. A build artifact is never legitimate Task output,
   and a guard there protects every repository rather than this one path in this
   one repository.
3. Consider making the bare `go build ./cmd/roundfix` unnecessary — a documented
   compile-check target (`make build` already exists and writes to the ignored
   path) that guides Agents away from the form that litters the tree.

## Suggested acceptance checks

- Running `go build ./cmd/roundfix` leaves no path that `git status` reports as
  untracked.
- A Task commit that would include an executable file is refused or reported.
- The Spec 0038 QA scope audit passes on a rerun.

## What worked — keep

- The QA gate's scope-integrity audit is the only check that caught this, and it
  caught it by comparing delivered paths against declared Spec scope. Keep that
  step; it finds a class of defect the test suite and code review cannot.
- The gate also caught the supervisor's own over-reach on the repair commit,
  which is the behavior you want from an authorization audit that is not
  deferential to whoever is driving.

---
status: accepted
created_at: 2026-10-02T00:00:00Z
updated_at: 2026-10-02T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Doctor checks what delivery needs, and only a definite answer fails

The Doctor Command proved Node.js, acpx, the adapters, the profiles and the
skills, and said nothing about the GitHub CLI, Git, the repository remote or
the tools a repository's Verification runs. On 2026-10-01 the maintainer
asked that setup recognize these requirements on a new machine, and the
adopter runbook sent that day still told each repository to check them by
hand. Every one of them stops the Delivery Queue or a Task when it is missing.

Doctor and Setup now print five more readiness lines: `gh`, `git`, `remote`,
`toolchain` and `environment`. Each finding on them carries a stable code
with the `DR-` prefix, a status and a next action. The status is `ok`,
`warn` or `failed`, and only `failed` changes the exit code.

- **A definite answer fails; an unanswered question warns.** Doctor fails a
  check only on a fact it read: an executable absent from `PATH`, a version
  below the floor, no `gh` account for the forge host, a token the forge
  rejected with HTTP 401, a repository permission below write, no delivery
  remote, a remote on a host that is not a GitHub host, an unset Git identity,
  or a tool a configured command needs and `PATH` lacks. A probe that timed out
  or could not connect reports `warn`, because an offline machine is not a
  misconfigured one.
- **Doctor may contact the forge, within bounds.** Doctor was offline. It now
  makes at most three network reads, all against the repository's own forge:
  `gh auth status` for the forge host, the viewer's permission on the
  repository through `gh`, and `git ls-remote` on the delivery remote. Each
  runs with no standard input, with Git's terminal prompt disabled, and under
  a ten-second limit that cancels the child process. Doctor still writes
  nothing, and every other check stays offline.
- **Tools are found, not run.** A tool is required when a configured command
  starts with it (the review Verification, the worktree bootstrap, a
  derived-path regeneration, or the Setup Manifest's Verification decisions)
  or when the repository lists it in Project Config. Doctor finds each one the
  way capability discovery does, without executing it (ADR-0087). Only `gh`
  and `git` run, and only with `--version` and the reads above.
- **Secrets are counted, never read into output.** The `environment` line
  reports which `ROUNDFIX_` key variables the optional features name are set,
  never a value, and never the provider's generic variables. It names each
  `NODE_OPTIONS` preload whose file is missing, with the same reading of
  `NODE_OPTIONS` the agent environment uses (ADR-0211).
- **Delivery refuses only what it cannot get past.** `roundfix deliver start`
  runs the `gh` and `remote` checks and refuses, before it creates a queue,
  on any `failed` finding. A `warn` does not stop it.

## Consequences

Doctor's run takes longer by up to thirty seconds on a machine whose forge
does not answer, and an offline Doctor exits zero with warnings where a
definite failure would have stopped a Run later. A repository can list a tool
no command names, such as a formatter a Makefile target calls, and an older
binary rejects that Project Config key. Setup reports these lines and offers
no install: installing `gh`, Git or a toolchain stays the person's action.

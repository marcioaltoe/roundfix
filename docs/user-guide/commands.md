# Command reference

Every Roundfix command: what it does, what it prints, and the boundary it never
crosses. This is the contract page — for the task-oriented walkthrough, read the
[operational guide](usage.md); for config keys, read
[configuration](configuration.md).

Examples call the installed `roundfix` binary. From a source checkout without
installing, substitute `go run ./cmd/roundfix`.

<!-- roundfix:command-index:begin -->
| Command | Guide |
| --- | --- |
| `archive` | [command guide](commands/archive.md) |
| `attach` | [command guide](commands/attach.md) |
| `baseline` | [command guide](commands/baseline.md) |
| `deliver` | [command guide](commands/deliver.md) |
| `detached-runs` | [command guide](commands/detached-runs.md) |
| `doctor` | [command guide](commands/doctor.md) |
| `events` | [command guide](commands/events.md) |
| `fetch` | [command guide](commands/fetch.md) |
| `gc` | [command guide](commands/gc.md) |
| `history` | [command guide](commands/history.md) |
| `implement` | [command guide](commands/implement.md) |
| `init` | [command guide](commands/init.md) |
| `migrate` | [command guide](commands/migrate.md) |
| `profiles` | [command guide](commands/profiles.md) |
| `qa-report` | [command guide](commands/qa-report.md) |
| `reconcile` | [command guide](commands/reconcile.md) |
| `reopen` | [command guide](commands/reopen.md) |
| `resolve` | [command guide](commands/resolve.md) |
| `review` | [command guide](commands/review.md) |
| `review-report-shape` | [command guide](commands/review-report-shape.md) |
| `runs` | [command guide](commands/runs.md) |
| `settle` | [command guide](commands/settle.md) |
| `setup` | [command guide](commands/setup.md) |
| `skills` | [command guide](commands/skills.md) |
| `spec` | [command guide](commands/spec.md) |
| `spec-audit` | [command guide](commands/spec-audit.md) |
| `stop` | [command guide](commands/stop.md) |
| `storage-report` | [command guide](commands/storage-report.md) |
| `supersede` | [command guide](commands/supersede.md) |
| `upgrade` | [command guide](commands/upgrade.md) |
| `watch` | [command guide](commands/watch.md) |
| `window` | [command guide](commands/window.md) |
<!-- roundfix:command-index:end -->

## Global contract

- **stdout** carries only the deterministic report of the requested command.
  Diagnostics, progress, the Run ID, and Agent output go to **stderr**.
- **Exit codes**: `0` Clean, Stopped, Fetched, or an already-complete no-op;
  `1` Unresolved, Failed, or Integration Pending; `2` Preflight Validation
  failure; `3` Clean Unverified (watch only); `130` in-terminal Ctrl-C.
- A bare `help` requests usage only as a command's first argument. `-h` or
  `--help` requests usage anywhere before `--`; no token after `--` requests
  usage.
- Color is automatic in interactive terminals. `ROUNDFIX_COLOR=always` forces
  it, `ROUNDFIX_COLOR=never` or `NO_COLOR` disables it.
- Supported Agent names are `codex`, `claude`, and `opencode`. The supported
  Review Source is `coderabbit`. It is the legacy PR-feedback source, read only
  by `fetch`, `watch`, and `resolve`, and never selects or requests a pre-PR
  reviewer.
- Commands that read the Run Database refuse another schema version without
  writing. An older database tells the operator to `run 'roundfix migrate'`;
  a newer database says that a newer Roundfix wrote it and tells the operator
  to use that binary or run `roundfix upgrade`.

## Setup and maintenance

## Review loop: fetch, resolve, watch

Review Runs execute in the user's checkout on the PR Head Branch — they create
no Run Worktree and no Run Branch, so review fixes are always a delta over the
published HEAD and Integration Pending does not exist as a review outcome.
Preflight Validation requires a clean tracked working tree (untracked files are
allowed); after a failed batch, every dirty path in the checkout is Agent work
from that Run, because the tree was clean at start.

**Branch Integrity Preflight** runs on all three commands before any Review
Source fetch, Agent Session, comment, or code change. It enumerates
`roundfix/run-*` branches with commits based on the PR Head Branch:
fast-forwardable pending work is integrated automatically and reported;
anything else refuses with exit `2`, naming each branch, its ahead count, and
the exact integration command. It also refuses while another Active Run is
bound to the target, naming the run id and the stop commands.
`--skip-branch-integrity` bypasses both guardrails only after publishing an
audit comment on the pull request recording the run id, the skipped
guardrails, and the ignored state; a failed publish fails the command.

`resolve`, `watch`, and `implement` accept exactly two Agent Selection forms:
omit `--agent`, `--model`, and `--reasoning-effort` to use profiles, or provide
all three together for a complete one-Run override. Any partial subset exits
`2` before config load, adapter or profile proof, Session creation, or Run
mutation. For example:

```bash
roundfix watch --source coderabbit --pr <number> --until-clean
roundfix watch --source coderabbit --pr <number> --agent codex --model gpt-5.6-sol --reasoning-effort high --until-clean
```

## Spec loop: implement, settle, archive

## Run discovery and monitoring

## Agent boundaries

Inside a Run, Agents own only assigned issue or task files, triage, code
edits, tests, verification commands, and assigned status updates. They must
not commit, push, resolve Review Source threads, edit unassigned files, or
mark issues `duplicated` — the Daemon owns every one of those boundaries.

---
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
prd: _prd.md
created: 2026-10-02
---

# A setup and doctor that recognize what the machine and the repository need — Technical Spec

## Executive Summary

Doctor builds a fixed list of `CheckResult` values from injected
dependencies and prints one line each. This Spec adds one readiness
component, used by Doctor, Setup and `deliver start`, that produces the five
new lines from findings carrying a `DR-` code, a status and a next action. The
`gh`, `git` and `remote` lines run `gh` and `git` through a bounded command
runner; the `toolchain` line finds executables with the capability-discovery
resolver that never runs them; the `environment` line reads `NODE_OPTIONS`
with Spec 0211's splitter and checks key variables for presence. The `skills`
line gains the comparison with the Setup Snapshot, through the restore
command's own profile contracts and tree digest. The primary trade-off is that
Doctor stops being offline: it accepts up to three bounded forge reads, and
answers a read that could not complete with `warn` rather than a guess, so an
offline machine passes and a misconfigured one fails (ADR-0220).

## Project Constraints

- Identifier strategy: applicable — new public identifiers: the line names
  `gh`, `git`, `remote`, `toolchain` and `environment`, the status token
  `warn`, the finding codes of "Finding codes", and the Project Config key
  `verification.tools`. They follow the existing line, status and config key
  forms. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the forge is read only through the
  user's `gh` and `git` with their own stored credentials. Roundfix never
  passes `--show-token`, never reads a token or key value into output or a
  file, and never logs in or out. Each read runs with no standard input,
  `GIT_TERMINAL_PROMPT=0` and a ten-second limit that cancels the child.
  Tests use a scripted runner and fake executables and never reach the
  network. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0220 (this Spec) decides the
  lines, codes, `warn`, the bounded forge reads and the `deliver start`
  refusal, ADR-0220: "Doctor fails a check only on a fact it read". ADR-0221
  (this Spec) decides the snapshot comparison, ADR-0221: "the managed refresh
  (`roundfix baseline update`) restores a trailing skill to its snapshot
  commit through the existing restore path". ADR-0087 governs tool discovery,
  ADR-0087: "Discovery therefore resolves a bounded symlink chain to its
  target and judges the target, while still never executing it". ADR-0211
  governs the reading of `NODE_OPTIONS`, ADR-0211: "A preload given as a
  package name, or as a path that exists, is kept". ADR-0172, ADR-0181,
  ADR-0107, ADR-0189, ADR-0187, ADR-0191, ADR-0192, ADR-0193, ADR-0032,
  ADR-0093, ADR-0096, ADR-0097, ADR-0204 and ADR-0206 hold as the PRD
  records, and ADR-0180 does not apply. ADR-0184 applies: the changed output is stated as Surface
  Transcripts. The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104,
  ADR-0155, ADR-0156 and ADR-0167. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_02 edits `.roundfixrc.yml` and task_04
  edits the Roundfix Skill, all Governed Paths; express maintainer
  authorization: "Autorizar os dois" and "considere autorizado a ajustar todas
  as skills se necessário"; bounded files: `.roundfixrc.yml`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need/_authorization.md`.

## System Architecture

No new command, flag or package.

| Component | Where | Change |
| --- | --- | --- |
| Readiness findings | new `internal/cli/readiness.go` | `readinessFinding`, folding into `CheckResult`, the `warn` status |
| Forge and Git readiness | new `internal/cli/readiness_forge.go` | `gh`, `git` and `remote` lines over a bounded runner |
| Toolchain and environment readiness | new `internal/cli/readiness_toolchain.go` | `toolchain` and `environment` lines |
| Doctor | `internal/cli/doctor.go`, `internal/cli/health.go` | Five lines after `pre-pr-review`; `warn` prints its next action; trailing skills on the `skills` line |
| Setup | `internal/cli/setup.go` | The five lines after `acpx`, diagnosis only |
| Delivery start | `internal/cli/deliver.go`, `internal/cli/cli.go` | Refuses on a failed `gh` or `remote` finding |
| Project Config | `internal/config/config.go` | `verification.tools` |
| Executable discovery | `internal/baseline/profile_alignment.go` | Exported wrapper over `resolveExecutableCandidate` |
| Preload reading | new `internal/agent/missing_node_preloads.go` | Exported wrapper over Spec 0211's `agentNodeOptions` |
| Snapshot comparison | new `internal/baseline/skills_trailing.go` | Installed tree digest against the profile's snapshot contract |
| Managed refresh | `internal/cli/baseline_update.go` | Trailing skills join the external drift it restores |
| This repository | `.roundfixrc.yml` | Declares its tools and derived paths |

## Implementation Design

### Interfaces

```go
// internal/cli/health.go
const CheckStatusWarn CheckStatus = "warn"
const (HealthCheckGH = "gh"; HealthCheckGit = "git"; HealthCheckRemote = "remote"
	HealthCheckToolchain = "toolchain"; HealthCheckEnvironment = "environment")

// internal/cli/readiness.go
type readinessFinding struct {
	Code   string      // stable, "DR-" prefixed
	Status CheckStatus // CheckStatusWarn or CheckStatusFailed
	Text   string      // never a token, key value or email address
	Next   string
}
// readinessResult folds findings into one line: the worst status wins, Detail
// joins "CODE: text" with "; " after okDetail, NextAction joins with " && ".
func readinessResult(name, okDetail string, findings []readinessFinding) CheckResult

type readinessRunner func(ctx context.Context, dir string, env []string, name string, args ...string) (stdout, stderr string, exitCode int, err error)
type readinessDependencies struct {
	run      readinessRunner                    // bounded, no stdin
	resolve  func(name string) (string, error)  // ADR-0087 discovery, never runs
	environ  []string
	exists   func(path string) bool
	timeout  time.Duration                      // 10s per forge read
}
func machineReadiness(ctx context.Context, deps readinessDependencies, loaded roundconfig.Loaded) []CheckResult // gh, git, remote, toolchain, environment
func forgeReadiness(ctx context.Context, deps readinessDependencies, loaded roundconfig.Loaded) []CheckResult   // gh, remote
```

```go
// internal/baseline/profile_alignment.go
func ResolveExecutable(name string, executableDirectories []string) (path string, reason string)
// internal/agent/missing_node_preloads.go (over Spec 0211's agentNodeOptions)
func MissingNodePreloads(value string, exists func(string) bool) []string
// internal/baseline/skills_trailing.go
func TrailingSetupSkills(repoRoot, profileID string, required []string) ([]string, error)
// internal/config/config.go
type Verification struct { /* existing fields */ Tools []string }
```

`doctorDependencies` gains `readiness func(context.Context, roundconfig.Loaded) []CheckResult`
and `trailingSkills func(repoRoot string, required []string) ([]string, error)`;
`setupDependencies` gains `readiness` of the same type; `commandDependencies`
gains `deliveryReadiness func(context.Context, roundconfig.Loaded) []CheckResult`.
The defaults call `machineReadiness`, `TrailingSetupSkills` with the Setup
Manifest's profile, and `forgeReadiness`. A `setupDependencies` value built
without `readiness` prints no readiness lines: the existing Setup tests build
that struct literally in `internal/cli/cli_test.go`, which this Spec does not
change, and new Setup tests set it.

### Invariants

1. A line's status is `failed` when any finding is `failed`, else `warn` when
   any finding is `warn`, else `ok`. Only `failed` changes an exit code.
2. A forge read that times out or cannot connect yields `warn`; a read that
   returned a definite answer yields `ok` or `failed`.
3. No output, log or test failure message contains a token, a key value or a
   Git email address.
4. Nothing in this Spec writes a file, a Git configuration value or the user's
   environment, except `baseline update` restoring a trailing skill after its
   confirmation.
5. `skills.RepositoryReadiness.Ready()` keeps its meaning; trailing is
   reported beside it.

### The forge lines

The delivery remote is `watch.push_remote`, or `origin` when it is empty, as
the Delivery Queue resolves it. `git remote get-url <remote>` gives its URL.
`https://host/owner/repo(.git)`, `ssh://user@host/owner/repo(.git)` and
`user@host:owner/repo(.git)` yield the forge host and `owner/repo`; any other
form, a local path or `file://` URL included, is not a forge remote.

- `remote`: absent → `DR-REMOTE-MISSING`; not a forge remote, or a host that
  is neither `github.com` nor a host `gh auth status --json hosts` lists →
  `DR-REMOTE-FORGE`; `git ls-remote <remote> HEAD` failing or timing out →
  `DR-REMOTE-UNREACHABLE` (`warn`). Outside a Git repository the line is
  `skipped` with `requires a Git repository`.
- `gh`: not found → `DR-GH-MISSING`; `gh --version` below 2.81.0 →
  `DR-GH-VERSION`. Then `gh auth status --active --hostname <host> --json hosts`
  (host `github.com` when there is no forge remote): no entry for the host →
  `DR-GH-UNAUTHENTICATED`; `state` `error` whose `error` names HTTP status
  401 → `DR-GH-TOKEN-REJECTED`; `state` `timeout` or any other `error` →
  `DR-GH-UNREACHABLE` (`warn`); `success` → the login is reported. With a
  forge remote and `success`, `gh repo view <owner/repo> --json
  viewerPermission`: `ADMIN`, `MAINTAIN` or `WRITE` → ok; any other value →
  `DR-GH-PERMISSION`; a failed read → `DR-GH-PERMISSION-UNVERIFIED` (`warn`).
- `git`: not found → `DR-GIT-MISSING`; `git --version` below 2.23.0 →
  `DR-GIT-VERSION`; `git config --get user.name` or `user.email` empty in the
  repository → `DR-GIT-IDENTITY` naming the key, never its value.

Every `gh` and `git` child runs through `readinessRunner` with a context that
carries `timeout`, no standard input and `GIT_TERMINAL_PROMPT=0` added to its
environment. `--show-token` is never passed.

### The toolchain and environment lines

Configured commands are the effective `defaults.verification`,
`worktree.bootstrap`, each `delivery.derived_paths[].regenerate`, and the
Setup Manifest's `verification.gate` and `verification.incremental` decision
values. A command's tool is its first word after leading `NAME=value` words.
A first word that is a shell keyword or builtin (`cd`, `export`, `set`,
`unset`, `test`, `[`, `if`, `for`, `while`, `until`, `case`, `exec`, `eval`,
`source`, `.`, `:`, `true`, `false`, `command`, `builtin`) or not a bare name
yields `DR-TOOL-UNREAD` (`warn`) naming the command. Each tool, and each
`verification.tools` entry, is resolved over the `PATH` directories with
`ResolveExecutable`; a tool that does not resolve yields one
`DR-TOOL-MISSING` naming the tool and every source that needs it.

The `environment` line calls `MissingNodePreloads` on the last `NODE_OPTIONS`
value and emits one `DR-NODE-PRELOAD-MISSING` (`warn`) per path. It then
reports, for the key variables the judge's transports name
(`judge.Load()`), each `ROUNDFIX_` variable as `set` or `not set` from a
non-empty test. It reads no other variable.

### Finding codes

| Code | Status | Next action |
| --- | --- | --- |
| `DR-GH-MISSING` | failed | `install GitHub CLI 2.81.0 or newer from https://cli.github.com` |
| `DR-GH-VERSION` | failed | `upgrade GitHub CLI to 2.81.0 or newer` |
| `DR-GH-UNAUTHENTICATED` | failed | `gh auth login --hostname <host>` |
| `DR-GH-TOKEN-REJECTED` | failed | `gh auth refresh --hostname <host>` |
| `DR-GH-UNREACHABLE` | warn | `re-run roundfix doctor when <host> is reachable` |
| `DR-GH-PERMISSION` | failed | `ask for write access to <owner/repo>, or gh auth switch --hostname <host>` |
| `DR-GH-PERMISSION-UNVERIFIED` | warn | `re-run roundfix doctor when <host> is reachable` |
| `DR-GIT-MISSING` | failed | `install Git 2.23.0 or newer` |
| `DR-GIT-VERSION` | failed | `upgrade Git to 2.23.0 or newer` |
| `DR-GIT-IDENTITY` | failed | `git config user.name <name>` or `git config user.email <address>` |
| `DR-REMOTE-MISSING` | failed | `git remote add <remote> <url>` |
| `DR-REMOTE-FORGE` | failed | `point <remote> at the repository's GitHub URL, or set watch.push_remote` |
| `DR-REMOTE-UNREACHABLE` | warn | `re-run roundfix doctor when <host> is reachable` |
| `DR-TOOL-MISSING` | failed | `install <tool>, or change the command that names it` |
| `DR-TOOL-UNREAD` | warn | `list the tools this command needs under verification.tools in Project Config` |
| `DR-NODE-PRELOAD-MISSING` | warn | `remove the preload from NODE_OPTIONS where your shell sets it` |
| `DR-SKILL-TRAILS-SNAPSHOT` | warn | `roundfix baseline update` |

### Skills that trail the snapshot

`TrailingSetupSkills` loads the profile with the restore command's
`loadRestoreProfile`, and for each required external skill whose snapshot
contract has a `treeDigest` and whose `.agents/skills/<name>` exists, compares
`portableRestoreDigest` of the installed files (read by
`inspectRestoreTarget`) with it. It returns the sorted names that differ.
Doctor calls it only for skills that are neither missing nor outdated. With a
trailing skill and an otherwise ready set the `skills` line is `warn`; with an
otherwise failed set the finding is appended and the line stays `failed`. When
the Setup Manifest names no profile the binary resolves, the detail says
`snapshot comparison unavailable` and the status is unchanged.
`baselineUpdateExternalDrift` adds the trailing names, so the preview lists
them and the restore it already runs brings them to the snapshot commit.

### Data Models

`roundconfig.Verification.Tools []string`, YAML `verification.tools`, default
empty. Project Config replaces the User Config list. Each entry is a bare
executable name matching `^[A-Za-z0-9][A-Za-z0-9._+-]*$`; an invalid or
duplicate entry is a config error naming `verification.tools`. No Run
Database change.

### API Contracts

1. API Contract: `roundfix doctor` prints `gh`, `git`, `remote`, `toolchain`
   and `environment` after `pre-pr-review` and before `skills`, each as
   `<name>: <status> (<detail>)`, where a `failed` or `warn` detail ends with
   `; next: <actions>`. Exit `1` only when a line is `failed`.
2. API Contract: `roundfix setup` prints the same five lines, with the same
   text, after `acpx` and before the adapter work, offers nothing for them, and
   exits `1` at its end when one is `failed`.
3. API Contract: `roundfix deliver start` runs `gh` and `remote` after the
   authorization check; any `failed` finding refuses with exit `2`, the
   `Preflight failed` block, and the reason `Delivery Queue cannot publish
   from this machine:` followed by one line per finding, before any queue
   record or owner exists.
4. API Contract: the `skills` line reports `DR-SKILL-TRAILS-SNAPSHOT` with the
   trailing names, and `roundfix baseline update` lists and restores them.
5. API Contract: `verification.tools` is accepted in User and Project Config
   and validated as in "Data Models".

### Surface Transcripts

1. Surface Transcript: Doctor on a machine with no `gh` account, no Git email
   and a missing `rtk`.

   ```transcript
   $ roundfix doctor
   stdout:
   ...
   pre-pr-review: ok (provider=codex; source=default)
   gh: failed (gh <version>; DR-GH-UNAUTHENTICATED: gh has no account for github.com; next: gh auth login --hostname github.com)
   git: failed (<version> >= 2.23.0; DR-GIT-IDENTITY: user.email is not set for this repository; next: git config user.email <address>)
   remote: ok (origin: github.com/<repository>; reachable)
   toolchain: failed (<found>; DR-TOOL-MISSING: rtk is not on PATH, needed by Setup Manifest verification.gate; next: install rtk, or change the command that names it)
   environment: warn (DR-NODE-PRELOAD-MISSING: NODE_OPTIONS preload "<path>" does not exist; spec judge keys: ROUNDFIX_OPENROUTER_API_KEY not set, ROUNDFIX_TYPESAFE_API_KEY not set; next: remove the preload from NODE_OPTIONS where your shell sets it)
   ...
   stderr:
   ...
   exit: 1
   ```

2. Surface Transcript: an offline machine warns and passes.

   ```transcript
   $ roundfix doctor
   stdout:
   ...
   gh: warn (gh <version>; DR-GH-UNREACHABLE: could not confirm the github.com login: <reason>; next: re-run roundfix doctor when github.com is reachable)
   ...
   remote: warn (origin: github.com/<repository>; DR-REMOTE-UNREACHABLE: git ls-remote origin failed: <reason>; next: re-run roundfix doctor when github.com is reachable)
   ...
   stderr:
   ...
   exit: 0
   ```

3. Surface Transcript: `deliver start` refuses when `gh` cannot publish.

   ```transcript
   $ roundfix deliver start <slug>
   stdout:
   stderr:
   Preflight failed

   Reason:
     Delivery Queue cannot publish from this machine:
     gh: DR-GH-UNAUTHENTICATED: gh has no account for github.com; next: gh auth login --hostname github.com

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix deliver start --help' for usage.
   exit: 2
   ```

4. Surface Transcript: an upstream skill that trails its snapshot.

   ```transcript
   $ roundfix doctor
   stdout:
   ...
   skills: warn (<counts>; DR-SKILL-TRAILS-SNAPSHOT: trails the Setup Snapshot: <names>; next: roundfix baseline update)
   ...
   stderr:
   ...
   exit: 0
   ```

## Vocabulary Contract

No new glossary term. The emitted words are the line names `gh`, `git`,
`remote`, `toolchain` and `environment`, the status `warn`, the `DR-` codes of
"Finding codes", `trails the Setup Snapshot`, `Delivery Queue cannot publish
from this machine` and the key `verification.tools`. Each Task documents the
words it emits in its command guide, and task_04 in the Roundfix Skill. The QA
gate's glossary check decides whether a readiness finding needs a term beside
**Doctor Command** and **Repository Skill Set**.

## Coverage Map

- Goal 1 → Forge and Git readiness, Toolchain and environment readiness, Doctor, Setup.
- Goal 2 → Invariant 2, The forge lines.
- Goal 3 → Delivery start, API Contract 3.
- Goal 4 → The toolchain and environment lines, Skills that trail the snapshot.
- Goal 5 → This repository, Project Config.
- User Story 1 → Doctor, Setup, Finding codes.
- User Story 2 → Invariant 2, The forge lines.
- User Story 3 → Delivery start.
- User Story 4 → The toolchain and environment lines.
- User Story 5 → Snapshot comparison, Managed refresh.
- User Story 6 → Project Config, The toolchain and environment lines.
- User Story 7 → The toolchain and environment lines.
- Core Feature 1 → Readiness findings, API Contract 1, API Contract 2.
- Core Feature 2 → The forge lines.
- Core Feature 3 → The forge lines.
- Core Feature 4 → The toolchain and environment lines, API Contract 5.
- Core Feature 5 → The toolchain and environment lines.
- Core Feature 6 → Skills that trail the snapshot, API Contract 4.
- Core Feature 7 → Delivery start, API Contract 3.
- Core Feature 8 → This repository.
- Core Feature 9 → Build Order 4 and the guide each Task edits.
- Success Metric 1 → Surface Transcript 1.
- Success Metric 2 → Surface Transcript 2, Invariant 2.
- Success Metric 3 → Surface Transcript 3.
- Success Metric 4 → The toolchain and environment lines, Invariant 3.
- Success Metric 5 → Surface Transcript 4, Managed refresh.
- Success Metric 6 → This repository, Project Config.

## Integration Points

- GitHub, through `gh auth status`, `gh repo view` and `git ls-remote`, read
  only, bounded, through `readinessRunner`; tests script the runner or put
  fake `gh` and `git` executables first on `PATH`.
- Spec 0211's `NODE_OPTIONS` splitter, through the exported wrapper.
- The embedded Setup Snapshots, through the restore command's profile
  loading.

## Testing Approach

1. Readiness unit tests in `internal/cli` drive each line through a scripted
   `readinessRunner`, a fake resolver and a fake environment: every code of
   "Finding codes", each status fold of Invariant 1, the `warn` and `failed`
   next-action printing, and a sentinel key value that never appears.
2. One probe test runs the real exec runner against fake `gh` and `git`
   scripts in a temporary `PATH` that record their argv, environment and
   standard input, and one that sleeps past a shortened timeout: no
   `--show-token`, `GIT_TERMINAL_PROMPT=0`, empty standard input, and a
   cancelled child.
3. Doctor and Setup tests assert Surface Transcripts 1, 2 and 4 and API
   Contract 2 through `runCLI` with fake dependencies; `withDoctorFakeDeps`
   gives the existing Doctor tests ready fakes, so their expected stdout
   gains five `ok` lines.
4. Delivery start tests assert Surface Transcript 3, no queue record, and that
   `warn` findings proceed; the existing start tests set ready fakes.
5. `internal/baseline` tests build a temporary repository whose installed
   external skill matches its lock and differs from the snapshot tree, and
   one that matches; `internal/cli` asserts the `baseline update` preview
   lists the trailing skill.
6. `internal/config` tests cover `verification.tools` and load this
   repository's `.roundfixrc.yml`.

## Build Order

1. Readiness findings, the `warn` status, and the `gh`, `git` and `remote`
   lines in Doctor, with the `doctor` guide.
2. `verification.tools`, the `toolchain` and `environment` lines in Doctor,
   and this repository's Project Config declarations, with the `doctor` and
   configuration guides (depends on: 1).
3. Skills that trail the snapshot in Doctor and the managed refresh, with the
   `doctor` and `baseline` guides (depends on: 2).
4. Setup's five lines, the `deliver start` refusal, the `setup` and
   `deliver` guides, and the Roundfix Skill's `setup`, `deliver` and
   `baseline` references (depends on: 3).
5. QA gate (depends on: 1, 2, 3, 4).

Steps 1 to 3 share `internal/cli/doctor.go`, its tests and the `doctor` guide,
and step 4 is the only one that edits the Roundfix Skill, so the chain is
serial and the skill version rises once.

## Risks & Considerations

- A slow forge adds up to thirty seconds to Doctor and to `deliver start`.
  The ten-second bound and the `warn` fallback keep it finite.
- `gh` older than 2.81.0 cannot report authentication as JSON; it fails with
  `DR-GH-VERSION` rather than being parsed from human text.
- A remote on GitHub Enterprise counts as a forge remote only when `gh` holds
  an account for its host.
- This repository's eleven trailing skills make its own Doctor print
  `skills: warn` until the managed refresh restores them, which this Spec
  leaves to a later change.
- The `delivery.derived_paths` declarations let the queue owner overwrite the
  declared files with the default branch's bytes and regenerate them; a
  regeneration that touches an undeclared path aborts, as ADR-0192 states.

## Decisions

- A finding fails only on a definite answer; network questions warn. See
  ADR-0220.
- Doctor contacts only the repository's own forge, with three bounded reads.
  See ADR-0220.
- Tools come from configured commands and `verification.tools`, found without
  running them. See ADR-0220.
- Trailing upstream skills warn and are restored by the managed refresh. See
  ADR-0221.
- Setup tests that build their dependencies literally get no readiness lines,
  so `internal/cli/cli_test.go` stays outside this Spec.

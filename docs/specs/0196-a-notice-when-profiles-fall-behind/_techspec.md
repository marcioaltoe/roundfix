---
spec: 0196-a-notice-when-profiles-fall-behind
prd: _prd.md
created: 2026-09-30
---

# A notice when profiles fall behind the recommendation — Technical Spec

## Executive Summary

One pure function compares a loaded configuration with the Recommended Profile
of Spec 0189 and returns, for each configured category, `current`, `differs`
or `pinned`. Four commands read it: `profiles check` prints it, `profiles
show` adds it to each category, Doctor reduces it to one line, and `upgrade`
writes it to standard error after its outcome. A profile key, `deviation`,
turns a deliberate difference into `pinned`. `profiles check --apply` builds a
profile fragment from the differing categories and hands it to the write that
`profiles configure` already owns.

The trade-off this design accepts is the life of a deviation. It names a
snapshot date, not the profile it declines, so it ends with every new snapshot
even when that category's recommendation did not move. That costs a maintainer
one edit per snapshot, and it keeps the key short enough to write by hand.

This Spec builds on Spec 0189. `RecommendedProfile`,
`ModelRecommendationSnapshotVersion` and schema `roundfix/profiles/v2` are
that Spec's, and every Task here assumes they exist.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; a Profile Deviation
  is named by a snapshot date. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local configuration and static
  data only. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0181 (this Spec), ADR-0180,
  ADR-0049, ADR-0107 and ADR-0140 hold as the PRD states. The gate is bound by
  ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and
  ADR-0156, and by ADR-0093 and ADR-0094. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the skill files, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No new package.

| Concern | Owner | File |
| --- | --- | --- |
| Profile Deviation, decode and carry | `ProfileEntry`, `decodeProfile` | `internal/config/profiles.go` |
| Profile Deviation, write | `CategoryChange`, `profileYAMLNode` | `internal/config/profile_config.go` |
| Comparison | `CheckRecommendations` | `internal/config/recommendation_check.go` (new) |
| `profiles check`, renderer | `runProfilesCheckCommand` | `internal/cli/profiles_check.go` (new) |
| `--apply` | `runProfilesCheckApplyCommand` | `internal/cli/profiles_check_apply.go` (new) |
| Show status, `check` routing | `printProfilesShowText`, `runProfilesCommand` | `internal/cli/profiles.go` |
| Shared write flow | `runProfilesConfigureCommand` | `internal/cli/profiles_configure.go` |
| Doctor line | `runDoctorCommand` | `internal/cli/doctor.go`, `internal/cli/health.go` |
| Upgrade notice | `runUpgradeCommand` | `internal/cli/upgrade.go` |
| Help text, dispatch | `commandUsage`, `runWithContext` | `internal/cli/cli.go` |

The comparison lives in `internal/config` because it needs only a `Config` and
the snapshot. It imports nothing new. Every caller is in `internal/cli`.

## Implementation Design

### Interfaces

```go
// internal/config/profiles.go
type ProfileDeviation struct {
	From   string `json:"from"`
	Reason string `json:"reason"`
}

type ProfileEntry struct {
	Profile   AgentSelectionProfile
	Source    ProfileSource
	Deviation *ProfileDeviation
}

// ResolvedProfile gains the same field.

// internal/config/profile_config.go
type CategoryChange struct {
	Category  WorkCategory
	Kind      ChangeKind
	Profile   AgentSelectionProfile
	Deviation *ProfileDeviation
}

// internal/config/recommendation_check.go
type RecommendationStatus string

const (
	RecommendationCurrent RecommendationStatus = "current"
	RecommendationDiffers RecommendationStatus = "differs"
	RecommendationPinned  RecommendationStatus = "pinned"
)

type CategoryRecommendation struct {
	Category    WorkCategory
	Status      RecommendationStatus
	Source      ProfileSource
	Configured  AgentSelectionProfile
	Recommended AgentSelectionProfile
	Deviation   *ProfileDeviation
}

type RecommendationCheck struct {
	Snapshot   string
	Categories []CategoryRecommendation
}

// CheckRecommendations compares the effective profile of every category in
// ConfiguredWorkCategories(config) with RecommendedProfile. It reads nothing
// but its argument and the snapshot.
func CheckRecommendations(config Config) (RecommendationCheck, error)
```

`internal/cli` adds one dependency to `upgradeDependencies`:

```go
// installedProfilesCheck runs `<executablePath> profiles check` in workDir
// and returns what it printed.
installedProfilesCheck func(ctx context.Context, executablePath, workDir string) ([]byte, error)
```

`runUpgradeCommand` gains the `commandEnvironment` argument the other commands
already take. `performUpgrade` returns its message together with whether a
release was installed and at which path.

### Data Models

One configuration key is added. No persisted schema and no Run Database row
changes.

```yaml
profiles:
  backend:
    preferred: {runtime: codex, model: gpt-5.6-sol, reasoning_effort: high}
    fallbacks:
      - {runtime: codex, model: gpt-5.6-terra, reasoning_effort: max}
    deviation:
      from: "2026-09-30"
      reason: "quota on the newer model is not provisioned here"
```

Rules for `deviation`:

- It is a mapping with exactly `from` and `reason`. Any other key is refused
  as `profiles.<category>.deviation.<key> is not a supported deviation key`.
- `from` is a calendar date written `YYYY-MM-DD`. `reason` is a non-empty
  string after trimming. A missing or malformed value is refused with its
  path.
- A profile without `deviation` decodes as it does today.
- A User Config or Project Config entry carries its own deviation. Built-in
  and legacy-derived entries carry none.

### The comparison

For each category in `ConfiguredWorkCategories(config)`, in
`AllWorkCategories` order:

- `current` when the effective Preferred Selection and Fallback Chain equal
  the Recommended Profile, in order.
- `pinned` when they differ and the entry's deviation has `from` equal to
  `ModelRecommendationSnapshotVersion`.
- `differs` otherwise. A deviation with another date stays on the row, so a
  caller can say that it ended.

An optional category the configuration does not define is not in the result.

### The renderer

`profiles check`, the `upgrade` notice and the Doctor line all read one
`RecommendationCheck`. The text renderer prints one block per category that is
not `current`, then a summary line, then a hint line when any category
differs. A profile is printed as its selections in order, joined by `, then `.

### Show

`profiles show` prints `Recommendation status: <status>` before the
`Recommended profile` block of each category, and
`Deviation: from <date> — <reason>` after it when the entry has one. The status
of an undefined optional category is `inherited`. The JSON profile object
gains `recommendation_status` and, when present, `deviation`.

### Adoption

`runProfilesCommand` routes `check` to the adoption command when `--apply` is
among its arguments, so the read-only command of task_02 is not edited.
`profiles check --apply --scope <scope>` selects the categories whose status
is `differs`. With `--scope user` it drops every category whose effective
profile comes from the Project Config and names each on standard error. It
builds a `Profiles` fragment holding the Recommended Profile of each selected
category, without a deviation, and runs the write flow of
`runProfilesConfigureCommand` from `PrepareProfilesConfig` onwards. That flow
is extracted into one function both commands call; its outputs and exit codes
do not change. When no category is selected, nothing is prepared, proved or
written.

### Doctor

One result named `recommendations` is appended right after the `profiles`
result. Its status is `ok` when no category differs and `found` otherwise, and
`skipped` with the error text when the comparison fails. Doctor fails only on
`failed`, so this result cannot change its exit code.

### Upgrade

After `upgrade` prints a release outcome to standard output, it writes the
notice to standard error:

- A release was installed: it calls `installedProfilesCheck` with the
  installed path and the process working directory, under a ten-second
  timeout, and copies the output to standard error.
- Nothing was installed (`no releases published`, `already current`,
  `upgrade available`): it loads the configuration as `profiles check` does
  and renders the comparison in process.
- Either way, a failure writes `roundfix: recommendations not checked:
  <reason>` and nothing else.

Help, a usage error and a failed upgrade print no notice. The exit code is
decided before the notice and is never changed by it.

### Surface Transcripts

```text
$ roundfix profiles check
recommendations: backend differs from snapshot 2026-09-30
  configured:  codex / gpt-5.6-sol / high, then codex / gpt-5.6-terra / max (project)
  recommended: codex / gpt-6.1-sol / high, then claude / opus / high
recommendations: review pinned against snapshot 2026-09-30: quota on the newer model is not provisioned here
recommendations: snapshot 2026-09-30; 3 current, 1 differ, 1 pinned
recommendations: adopt with `roundfix profiles check --apply --scope user|project`, or declare a deviation
$ echo $?
0
```

```text
$ roundfix profiles check
recommendations: snapshot 2026-09-30; 5 current, 0 differ, 0 pinned
```

A category whose deviation names another snapshot prints
`recommendations: backend differs from snapshot 2026-09-30 (its deviation was
declared against 2026-08-07)` as its first line.

```text
$ roundfix doctor
…
profiles: ok (3 distinct tuples; 10 category references)
recommendations: found (snapshot 2026-09-30; 3 current, 1 differ, 1 pinned; run roundfix profiles check)
…
```

```text
$ roundfix upgrade --check 2>notice.txt
already current <version>
$ cat notice.txt
recommendations: snapshot 2026-09-30; 5 current, 0 differ, 0 pinned
```

### API Contracts

1. API Contract: `roundfix profiles check [--json]` — read-only and offline.
   Text is the transcript above, on standard output. JSON has schema
   `roundfix/profiles-check/v1` and the fields `snapshot`, `current`,
   `differ`, `pinned` and `categories`; each category carries `category`,
   `status`, `source`, `configured` and `recommended` (each with `preferred`
   and `fallbacks`) and, when present, `deviation` (`from`, `reason`). Every
   configured category is listed, in category order. Exit `0` when the
   comparison ran, whatever it found; exit `2` for a usage error or a
   configuration that does not load.
2. API Contract: `roundfix profiles check --apply --scope user|project
   [--dry-run] [--yes] [--json]` — `--apply` requires `--scope`. Without
   `--apply`, the flags `--scope`, `--dry-run` and `--yes` are refused as
   unknown. Either violation is a usage error with exit `2`. Output, confirmation and exit codes are those of
   `roundfix profiles configure` for the same fragment. When no category is
   selected, standard output is `Profile configuration unchanged: nothing to
   adopt` (or, with `--json`, a `roundfix/profiles-configure/v1` object with
   `changed: false` and empty `profiles` and `changes`) and the exit code is
   `0`.
3. API Contract: profile key `deviation` — accepted in User Config, Project
   Config and a `profiles configure` fragment, under the rules in Data Models.
   `profiles configure` writes it with the profile, prints
   `Deviation: from <date> — <reason>` in its preview, and adds `deviation` to
   the profile in its JSON, only when the profile carries one.
4. API Contract: `roundfix profiles show` — adds `Recommendation status` and
   `Deviation` lines and the JSON fields `recommendation_status` and
   `deviation`. Schema and exit codes are unchanged.
5. API Contract: `roundfix doctor` — adds the `recommendations` line after
   `profiles`. It never makes Doctor exit non-zero.
6. API Contract: `roundfix upgrade [--check]` — standard output and exit codes
   are unchanged. Standard error carries the notice on every release outcome,
   from the installed executable after an install.

## Coverage Map

- Goal 1 → Upgrade; Doctor; API Contracts 5 and 6.
- Goal 2 → The renderer; Adoption; API Contracts 1 and 2.
- Goal 3 → Data Models; The comparison; API Contract 3.
- Goal 4 → Testing Approach 6.
- User Story 1 → Upgrade; API Contract 6.
- User Story 2 → The renderer; Adoption; API Contracts 1 and 2.
- User Story 3 → Data Models; The comparison; API Contract 3.
- User Story 4 → Testing Approach 6.
- Core Feature 1 → Data Models; API Contract 3.
- Core Feature 2 → The comparison; The renderer; API Contract 1.
- Core Feature 3 → Show; API Contract 4.
- Core Feature 4 → Adoption; API Contract 2.
- Core Feature 5 → Doctor; API Contract 5.
- Core Feature 6 → Upgrade; API Contract 6.
- Core Feature 7 → Testing Approach 6.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 5.
- Success Metric 5 → Testing Approach 5.
- Success Metric 6 → Testing Approach 6.

## Integration Points

- **Spec 0189.** The comparison reads `RecommendedProfile` and
  `ModelRecommendationSnapshotVersion`. No Task here changes either.
- **`profiles configure`.** Adoption reuses its prepare, proof, preview,
  confirmation, persist and rollback. Its existing tests must pass unchanged.
- **ACP adapters.** `profiles check` without `--apply`, Doctor's new line and
  the notice open no Agent Session. `--apply` proves through the runner
  `profiles configure` uses, which tests replace with the existing fake.
- **The installed executable.** After an install, `upgrade` starts one child
  process, the new executable with `profiles check`. Tests replace it with a
  fake through `upgradeDependencies`.

## Testing Approach

1. **Profile Deviation.** New `internal/config/profile_deviation_test.go`
   requires a deviation to load from User Config and from Project Config, each
   malformed form to be refused with its path, a profile without one to load
   as before, a fragment with one to be written with the profile, and a
   replaced profile to lose it. New `internal/cli/profiles_deviation_test.go`
   requires `profiles configure` to preview and write a deviation, and its
   output to be unchanged for a profile without one.
2. **Comparison and `profiles check`.** New
   `internal/config/recommendation_check_test.go` requires the three statuses,
   the snapshot rule, built-ins reported `current`, and undefined optional
   categories left out. New `internal/cli/profiles_check_test.go` requires the
   text transcript, the JSON shape, exit `0` with differences, exit `2` for
   usage errors, that the runner is never called and that no file changes.
   New `internal/cli/profiles_show_status_test.go` requires the status and
   deviation in text and JSON.
3. **Adoption.** New `internal/cli/profiles_check_apply_test.go` requires the
   written profiles, a pinned category's bytes untouched, proof before the
   write, `--dry-run` and a declined confirmation writing nothing, the
   user-scope rule, the nothing-to-adopt outcome and the flag rules. The
   existing `profiles configure` tests run unchanged.
4. **Doctor.** New `internal/cli/doctor_recommendations_test.go` requires the
   `ok` and `found` lines, their position after `profiles`, and exit `0` with
   differences. Five existing tests in `internal/cli/doctor_test.go` gain the
   line.
5. **Upgrade.** New `internal/cli/upgrade_characterization_test.go` is written
   first, against the unchanged code, and fixes standard output and the exit
   code of every outcome, the failure and the usage error. New
   `internal/cli/upgrade_notice_test.go` requires the notice on each release
   outcome, the call to the installed executable with the installed path after
   an install and not otherwise, an unchanged exit code and one `not checked`
   line when the check fails, and no notice on help, a usage error or a failed
   upgrade. Two existing tests in `internal/cli/upgrade_test.go` that require
   an empty standard error are updated.
6. **Off the hot paths.** New
   `internal/cli/recommendation_check_scope_test.go` parses the non-test Go
   files under `internal` and `cmd` and requires `CheckRecommendations` to be
   referenced only in `internal/config/recommendation_check.go`,
   `internal/cli/profiles_check.go`, `internal/cli/profiles_check_apply.go`,
   `internal/cli/profiles.go`, `internal/cli/doctor.go` and
   `internal/cli/upgrade.go`.

Each new gate is proved to fail: the Task that adds it records, in its Result,
the sabotage it applied and the failing test name, then restores the code.

## Build Order

1. Profile Deviation in the configuration and in `profiles configure`, with
   its guide, glossary and skill text, task_01 (depends on: none).
2. The comparison, `profiles check` and the status in `profiles show`, with
   their guide and skill text, task_02 (depends on: 1). The `pinned` status
   reads the deviation.
3. `profiles check --apply`, with its guide and skill text, task_03 (depends
   on: 2).
4. The Doctor line and the `upgrade` notice, the scope test, and their guide
   and skill text, task_04 (depends on: 3). It shares `internal/cli/cli.go`,
   the guides and the skill with task_03, and its hint names `--apply`.
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Order of delivery.** This Spec compiles only on top of Spec 0189. Its
  Tasks assume `RecommendedProfile` exists and that `profiles show` already
  prints the `Recommended profile` block.
- **Extracting the write flow.** `--apply` must not fork the behavior of
  `profiles configure`. The extraction is mechanical, and the existing
  configure tests are the characterization: they run unchanged.
- **A slow or broken installed executable.** The child process has a
  ten-second timeout and its failure is one line. `upgrade` has already
  printed its outcome and fixed its exit code by then.
- **Noise.** `upgrade` now always writes one line to standard error. It is one
  line when nothing differs, and scripts that read standard output are
  unaffected.
- **A configuration with `deviation` on an older Roundfix.** It is refused as
  an unknown profile key. The break is declared, and the guide says so.

## Decisions

- The comparison runs only where a person asked, and decides nothing. See
  ADR-0181.
- A deviation names a snapshot date and ends with the next snapshot.
- Adoption is `profiles configure` with a generated fragment.
- The notice after an install comes from the installed executable.

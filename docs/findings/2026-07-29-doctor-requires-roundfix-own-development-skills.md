---
status: pending
created_at: 2026-07-29
updated_at: 2026-07-29
---

# Repository Skill Set — Doctor demands Roundfix's own development skills in every repository (2026-07-29)

Refreshing four consumer repositories after the `0.0.2` release, three of them
failed `roundfix doctor` because the required external skill set includes
`golang-cli`, `golang-concurrency`, `golang-context`, `golang-error-handling`,
`golang-lint`, `golang-testing`, `bubbletea`, and `tui-design`. None of the
three writes Go or renders a terminal UI. The requirement is Roundfix's own
development stack leaking into every repository that installs its skills.

A Vortex session on 2026-07-29 reported the same symptom independently and
correctly declined to act on it: "o check `skills:` do doctor segue vermelho
porque o Roundfix exige as skills de desenvolvimento dele próprio — não é
problema do vortex."

## 1. The required external set does not follow the repository's Profile

- **Symptom / evidence**: after installing the `0.0.2` owned skills and
  refreshing the external set, `roundfix doctor` reported one missing skill at
  a time, always from Roundfix's own stack:

  ```text
  skills: failed (validate required skills lock entry ".../skills-lock.json":
  external skill "autoresearch" is missing;
  next: bunx skills experimental_install && bunx skills update -p -y)
  ```

  Computing the full gap against Roundfix's own 25 external skills gave nine
  missing in `oraculum`: `autoresearch`, `bubbletea`, `golang-cli`,
  `golang-concurrency`, `golang-context`, `golang-error-handling`,
  `golang-lint`, `golang-testing`, and `tui-design`.

  The lock inventories explain which repository passes and which does not:

  | repository | skills in lock | Go/TUI skills |
  | --- | --- | --- |
  | fluxus | 110 | 8 |
  | oraculum | 107 | 0 |
  | vortex | 101 | 0 |
  | tax-poc | 108 | 0 |

  `fluxus` reaches `skills: ok` only because it already carries the eight Go
  and TUI skills. The three TypeScript repositories carry none, and each fails.

- **Root cause**: unknown in detail, but bounded. The Baseline catalog declares
  `requiredSkills` per module, so a repository that never selected the Go or
  TUI modules should never require their skills. Something between the
  repository's Setup Manifest and the readiness check is resolving the required
  set from a source other than the repository's own selected modules.

- **Action / suggestion**: resolve the required external skill set from the
  repository's Setup Manifest — its selected Profile and modules — rather than
  from a fixed or binary-derived list. A TypeScript repository must never be
  told it needs `bubbletea`. Until then the check is a false negative on every
  consumer repository, and the printed next action makes it worse by proposing
  an install that would satisfy it.

## 2. The proposed remediation over-installs by an order of magnitude

- **Symptom / evidence**: doctor's next action is
  `bunx skills experimental_install && bunx skills update -p -y`. The install
  step reported `No valid skills found`, so the natural next attempt was
  `bunx skills add marcioaltoe/skills`, which installed the **entire** upstream
  catalog: `oraculum` went from 100 to 174 skills on disk and from 106 to 180
  lock entries, a diff of 943 files and roughly 182,000 inserted lines, for a
  repository whose required set is 39. The change was reverted.

  The targeted form `bunx skills add marcioaltoe/skills@<skill>` installs one
  skill correctly, so the capability exists; nothing in the printed next action
  points at it.

- **Root cause**: the remediation names a package-level install for a
  skill-level gap.

- **Action / suggestion**: print the exact per-skill install command for each
  missing external skill instead of a package-wide install. Every skill file
  is an instruction that runs with full agent permissions, so a remediation
  that pulls a whole catalog is a supply-chain widening, not a fix — the CLI
  itself warns "review skills before use; they run with full agent
  permissions" while installing 74 unreviewed additions.

## 3. Consumer lock files drift in both directions with no repair path

- **Symptom / evidence**: `oraculum`'s lock required `agentic-cli-design`,
  which upstream publishes but the repository had neither locked nor
  installed, while still referencing six skills upstream no longer publishes:
  `composition-patterns`, `monorepo-management`, `qa-execution`, `qa-report`,
  `react-email`, and `setup-workflow`. `bunx skills update` reports each of the
  six as a failed update on every run.
- **Root cause**: nothing reconciles a consumer lock against the upstream
  catalog; entries are added by installs and never pruned when upstream removes
  a skill.
- **Action / suggestion**: give the update path a reconciliation mode that
  reports entries whose upstream source no longer exists and offers to drop
  them, so a stale entry stops producing a permanent failure line.

## Impact on the release rollout

`fluxus` was refreshed successfully and committed on
`ma/refresh-skills-2026-07-29`. The other three repositories were left
untouched: satisfying the check would have meant installing Roundfix's Go and
TUI development skills into three TypeScript projects. That is the correct
outcome for them and the reason this finding exists.

## What worked — keep

- The readiness check is right to be blocking and right to name a next action;
  the defect is in which set it requires and which command it proposes, not in
  the existence of the gate.
- `roundfix skills install --target project` installed the fourteen owned
  skills from the running binary cleanly in every repository, which is the half
  of the contract that behaved exactly as designed.

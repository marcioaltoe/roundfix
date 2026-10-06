---
status: approved
granted: 2026-10-06
action: let the skill snapshot comparison, skill restore and lock reconciliation resolve a repository-owned Baseline Profile as baseline update does, so baseline update and Doctor work for repository-profile adopters, and describe the rule in the user guides and the Roundfix Skill's baseline reference
consuming: 0236-baseline-update-with-a-repository-profile
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0236

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-06 the maintainer answered the cycle's
questions through AskUserQuestion: the scope "Y1, Y2 e Y3", of which this Spec
is Y1, and, for the Governed Paths this Spec declares, "Concedo". The purpose
is the defect two adopters reported on 2026-10-05: `roundfix baseline update`
exits 1 in every repository adopted with its own Baseline Profile.

The governed set was measured with `GovernedPath` on the authoring branch at
`bde616d1`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/baseline.md` — the repository's
  skill-sync rule requires a Pull Request that changes CLI behavior to ship
  the skill update, and this reference tells agents to pass a built-in
  profile ID to `baseline skills restore` and `baseline skills reconcile`,
  which now accept the repository's own profile.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields by the record command.

## What is not governed

The Go sources and tests under `internal/baseline` and `internal/cli`,
`docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md`,
`docs/user-guide/commands/baseline.md`, `docs/user-guide/commands/doctor.md`,
the reference mirror `skills/roundfix/references/baseline.md`, and
`skills/testdata/owned-skill-versions.json` are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to the repository profile schema, the embedded catalog, its
  Setup Snapshots, modules or profiles, or any derived Baseline file.
- No test, Verification command or QA row reaches the network, reads or
  writes the real `~/.roundfix`, or touches an adopter's repository.
- No change to archived Specs, existing QA Reports, `CONTEXT.md`,
  `CHANGELOG.md` or the `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.

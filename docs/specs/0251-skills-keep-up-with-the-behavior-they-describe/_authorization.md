---
status: approved
granted: 2026-10-08
action: add the Skill Coverage Map and the Behavior Surface Record with their repository contract, make the release plan refuse a release with a Lagging Surface, add the skills-Task authoring rule to the Spec Consistency Check, and describe them in the owned skills, the user guide, the repository rules and the glossary
consuming: 0251-skills-keep-up-with-the-behavior-they-describe
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/release.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-tasks/SKILL.md
  - .roundfixrc.yml
  - docs/agents/specific-repository.md
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - internal/speccheck/coherence.go
  - internal/speccheck/constraints.go
  - skills/baseline_skill_contract_test.go
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0251

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, with the standing answer "Concedo" for the Governed Paths each
Spec declares, and authorized every owned skill: "considere autorizado a
ajustar todas as skills se necessário".

On 2026-10-08 the maintainer decided how the owned skills stay current with
shipped behavior: "Check no release + regra de autoria (Recommended)". The
release plan refuses, or reports as blocking, a release in which a skill fell
behind a change to a command's help or flags, a configuration key, an exit
code or a user guide since the previous release, with a checked-in, testable
map of which skill covers which behavior and a way to record "reviewed, no
change needed". Spec authoring requires a skills Task for every change a skill
describes, encoded where Spec authoring is checked. The same day the
maintainer granted:

- the Governed Paths this Spec declares, including skills, Baseline modules
  and guides, the `Makefile` and CI workflows. This Spec changes no Baseline
  module, no `Makefile` and no workflow, so it bounds none of them;
- a standing `qa_override` for an environment-only `partial`.

The external skills distribution repository and every adopting repository are
out of scope and are never touched. No live provider call is authorized by
this record. Authoring, tests, Verification and QA use temporary
repositories, temporary homes and fixture Specs.

The governed set was measured with `GovernedPath` on the authoring branch at
`c73a92e0`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository, run against every file the Tasks
declare.

## Why each governed path is unavoidable

- `.roundfixrc.yml` gains the derived-path declaration of the Behavior Surface
  Record and its record command, so two items that change a surface merge by
  regeneration (ADR-0192).
- `.agents/skills/roundfix/references/release.md` and
  `.agents/skills/roundfix/references/spec.md` describe the new release check
  and the new consistency codes; `.agents/skills/roundfix/SKILL.md` and its
  mirror `skills/roundfix/SKILL.md` carry the version the record command
  raises.
- `.agents/skills/write-tasks/SKILL.md` and its mirror state the skills-Task
  rule where Task Graphs are authored.
- `.agents/skills/write-prd/references/prd-template.md`, its mirror, and the
  version lines of `.agents/skills/write-prd/SKILL.md` and its mirror carry
  the optional `## Skills` section of the PRD.
- `docs/agents/specific-repository.md` holds the hard rule on roundfix skill
  sync, which now names the map, the release check and the authoring rule.
- `internal/speccheck/coherence.go` and `internal/speccheck/constraints.go`
  register and run the two new detectors.
- `internal/docscontract/testdata/corpus-golden.json` and
  `internal/spec/archive_layout_characterization_test.go` characterize the
  active corpus and gain the two new codes at 0.
- `skills/baseline_skill_contract_test.go` pins the new write-tasks phrase.

## What is not governed

These paths are ordinary: the new `internal/skillcoverage` package, the map
`docs/references/skill-coverage.json`, the record
`docs/references/behavior-surfaces.json`, the release plan sources and tests in
`internal/cli`, the new docscontract and speccheck tests,
`internal/config/regeneration_declared_test.go`,
`internal/docscontract/corpus_test.go`, the mirrors of the roundfix skill's
references, `skills/testdata/owned-skill-versions.json`, the user guide,
`CONTEXT.md` and
`docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors after each approved
skill edit and adds no source path.

```yaml
command: make skills-sync
```

The Tasks also run two record commands whose packages install no suite
guard, so neither needs a declared output list here:

- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  rewrites the version lines of each edited skill's `SKILL.md` and its mirror
  and `skills/testdata/owned-skill-versions.json`;
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  rewrites `docs/references/behavior-surfaces.json`, which `.roundfixrc.yml`
  declares as a derived path with that command.

No Baseline module changes, so neither `make baseline-digests` nor the module
version record command is declared.

## Limits

- No new dependency in `go.mod`. No change to the `Makefile`, a CI workflow,
  a lint or formatter configuration, or any Baseline module or setup snapshot.
- The release plan's reset mode, its JSON schema version, its decision states,
  its proposed version and its skills and baseline checks are unchanged.
- No test, Verification command or QA row reaches a provider, starts a real
  Agent Session, reads a credential, or reads or writes the real
  `~/.roundfix`. No row triggers a GitHub workflow or creates a tag.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.

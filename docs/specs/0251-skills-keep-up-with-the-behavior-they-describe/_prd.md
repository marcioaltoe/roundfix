---
spec: 0251-skills-keep-up-with-the-behavior-they-describe
status: active
created: 2026-10-08
surfaces: [backend, cli, docs]
---

# Skills keep up with the behavior they describe

The Roundfix-owned skills teach Agents how Roundfix behaves, so a skill that
trails the shipped behavior teaches the wrong command. Today the skill mirror
and the skill version record are checked, but nothing checks that a change to
a command's help or flags, a configuration key, an exit code or a user guide
reached the skills that describe it. Specs update skills only when their
author remembers. The operator's queue log records three QA failures of that
class in two cycles, two of them in one Spec, where the settle reference
contradicted a new reopen rule. The release runbook ends its skills-and-guides
step with a reading pass, and the Release Plan Command reports its skills and
baseline checks without letting them change its result. This Spec makes the
release plan refuse a release whose skills fell behind, and makes Spec
authoring require a skills Task for a change a skill describes. On 2026-10-08
the maintainer chose "Check no release + regra de autoria (Recommended)"; the
design is ADR-0256.

## Prerequisites

Specs 0249 and 0250 are authored in parallel and deliver first. Neither is a dependency of this Spec's Verification. Both change owned
skills and recorded versions, so this Spec's Tasks raise each skill to the
next version their record command chooses, never a fixed number.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes. Behavior
  Surface identifiers are readable names built from existing command paths and
  repository paths, the new check status values are lowercase words, and the
  two new consistency codes follow the existing `SC-` code family. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added. The release check reads local Git objects, the record is
  computed from the built command help and repository files, and every test
  uses temporary repositories. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0256 (this Spec) decides the
  Behavior Surfaces, the map and record, the blocking release check and the
  authoring rule. ADR-0189: "Every release runs a skills-and-guides check
  before its Pull Request", and that step "rests on checks that need no model
  and no network"; the new check is one such check. ADR-0192: "conflict confined to
  declared derived paths is resolved by regeneration",
  which is why the record is declared derived. ADR-0233: "a conflict confined
  to declared paths takes the default branch's whole file", which is why the
  authored map is not. ADR-0222: "A Baseline guide says only what holds for
  the repository that reads it", which is why no Baseline module changes.
  ADR-0187: "A Task declares the one command file it changes", which the map
  follows for command coverage. ADR-0252: a `//verify:` directive "adds
  declared inputs to the package directory". ADR-0253: "A repository test
  asserts that a change to any declared output or line path selects the
  test", so the regeneration directive names the new record. ADR-0184: "A
  TechSpec states a command surface as a transcript", answered by the
  TechSpec's Surface Transcripts. ADR-0250: "A Baseline module's version names
  one content, and the record step chooses it"; this Spec changes no Baseline
  module, and its skill edits follow the same rule through the owned-skill
  record command, which chooses the next version. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization on
  2026-10-08 grants the Governed Paths this Spec declares, including skills,
  and a standing `qa_override` for an environment-only `partial`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0251-skills-keep-up-with-the-behavior-they-describe/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/release.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-tasks/SKILL.md`, `.roundfixrc.yml`,
  `docs/agents/specific-repository.md`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/baseline_skill_contract_test.go`, `skills/roundfix/SKILL.md`,
  `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`,
  `skills/write-tasks/SKILL.md`.

## Goals

- A release whose range changed a command's help, the configuration keys, the
  exit codes or a user guide, without changing a skill that describes it or
  recording why none needs to, is refused by the release plan before any
  release mutation.
- Which skill describes which behavior is a checked-in, tested artifact, and a
  stale fingerprint or an unmapped behavior fails the repository contracts at
  the Pull Request boundary.
- A reviewed change that needs no skill text is recorded once, in the range it
  belongs to, and stops being reported.
- A Spec whose Tasks change a described behavior without a Task that changes
  its skill is refused by `roundfix spec check` while it is still being
  authored.
- Repositories without the map, and the release that introduces it, are
  reported but never refused.

## User Stories

1. As the maintainer cutting a release, I want the release plan to name each
   behavior that changed since the previous release without its skill, so that
   I fix the skill before the release Pull Request instead of after a user
   reads it.
2. As a Spec author, I want `roundfix spec check` to tell me that my Tasks
   change a described behavior and no Task updates its skill, so that the
   skills Task is in the graph before the Run starts.
3. As a reviewer of a change that alters behavior no skill states, I want to
   record that review once with its reason, so that the release check does not
   report it again.
4. As an Agent changing a command, I want the repository contracts to fail when
   the recorded fingerprint of that command's help is stale, so that the
   record cannot fall out of step with the behavior.
5. As the operator of an adopting repository with no map, I want the release
   plan to keep its current result, so that a check built for this repository
   never refuses mine.

## Core Features

1. **Behavior Surface.** The unit of described behavior: the help of one
   command path from the root help, the set of configuration keys, the table
   of exit codes, and each file of the user guide. Each has a readable
   identifier.
2. **Skill Coverage Map.** A checked-in, authored list of every Behavior
   Surface with the owned skill files that describe it, or the reason none
   does, the source paths whose change can change it, and an optional
   **Coverage Review**: a dated note recording that a change needs no skill
   text. Command coverage follows the Roundfix skill's reference index.
3. **Behavior Surface Record.** A generated record of one fingerprint per
   Behavior Surface, written only by its record command and declared as a
   derived path. A repository contract test computes every fingerprint and
   fails when the record is stale, when a surface has no map entry, when an
   entry names no existing surface, or when a covering file is not an owned
   skill file.
4. **Release check.** A range Release Plan gains a `skill-coverage` check that
   compares the map and the record at the base and at the target. A surface
   that changed, appeared or disappeared is a **Lagging Surface** unless a
   covering skill file changed in the range, the surface is uncovered, or its
   Coverage Review changed in the range. A Lagging Surface, or a map or record
   that cannot be read, makes the check blocking: for a plan that proposes a
   version the command exits 3, names the check in its next action and lists
   each Lagging Surface. The decision state, the proposed version and the
   skills and baseline checks are unchanged. A repository without a map
   reports `not_declared`, a base without one reports `introduced`, and
   neither blocks.
5. **Authoring rule.** `roundfix spec check` reports a Behavior Surface whose
   source a non-QA Task declares when no non-QA Task declares one of its
   covering skill files and the PRD's **Skills Declaration** does not excuse
   it. An excuse is a `- unchanged:` entry naming the surface and a reason,
   and it needs a non-QA Task that declares the map, where the Coverage Review
   is recorded. A malformed entry or an unknown surface is reported too. The
   rule starts at the commit that adds the map. The write-tasks and write-prd
   skills and the repository rules state it.
6. **Docs and glossary.** The release runbook, the user guide, the Roundfix
   skill's release and spec references and the glossary describe the new
   check, the record, the review and the rule.

## Non-Goals / Out of Scope

- The external skills distribution repository and any adopting repository.
- A Baseline module or guide change; the rule needs a map that only this
  repository keeps (ADR-0256).
- A model or network call in any check.
- Comparing the first release after this Spec against a base that has no map.
- Mapping every transitive source of a command's help: the shared usage file
  is the source of no single command, and the release check catches a change
  made only there.
- Changing the Makefile, a CI workflow, `make verify-docs`, the reset mode of
  the release plan, its JSON schema version, or the skills and baseline checks.
- Telling feature, refactor and fix Specs apart: the artifacts do not record
  the route, so the rule applies to every Spec that changes a covered surface.

## Success Metrics

1. Success Metric: in a temporary repository whose range changes a command's
   help without its covering skill, `roundfix release plan` exits 3, its
   `skill-coverage` line is `behind` and names the Lagging Surface, and the
   decision state and proposed version equal those of the same range without
   the change.
2. Success Metric: the same range with a covering skill edit, or with a
   changed Coverage Review, exits 0 with `skill-coverage: current`.
3. Success Metric: a repository without a map, and a base without one, keep
   the exit code and decision they had before this Spec, with
   `not_declared` and `introduced`.
4. Success Metric: at the audited head the repository contract for the map
   passes, and editing one command's help without re-recording makes it fail
   naming that surface.
5. Success Metric: a fixture Spec whose Task declares a covered surface's
   source without a skills Task reports `SC-SKILLS-UNTASKED`; adding the skills
   Task, or a `- unchanged:` entry with a Task that declares the map, clears
   it.
6. Success Metric: every pre-existing release plan test and the active Spec
   corpus golden pass unchanged except for the declared new codes at 0.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The operator's queue log, entries 193, 231 and 232 of
  `~/.roundfix-operator/queue-interventions.md`: QA failures in Specs 0235 and
  0244 where the settle reference and the reopen guide trailed the shipped
  behavior.
- The Drift documentation-sync extension's published documentation
  (marketplace.visualstudio.com/items?itemName=pallaprolus.drift): a reviewed
  drift finding is stored with a hash of the code it was reviewed against and
  "returns only when the code changes again".
- Tan, Wagner and Treude, "Detecting Outdated Code Element References in
  Software Repository Documentation" (arxiv.org/abs/2212.01479): of more than
  3,000 GitHub projects, most contained at least one outdated code element
  reference at some point in their history.

## Glossary

- adds: **Behavior Surface**
- adds: **Skill Coverage Map**
- adds: **Behavior Surface Record**
- adds: **Coverage Review**
- adds: **Lagging Surface**
- adds: **Skills Declaration**
- changes: **Release Plan**

## Decisions

- Behavior Surfaces, the authored map and the generated record; see ADR-0256.
- The release check blocks only through the exit code and the next action; see
  ADR-0256.
- A Coverage Review answers the range in which it changed; see ADR-0256.
- The authoring rule reads the map's sources and applies to every route; see
  ADR-0256.
- No Baseline module changes; see ADR-0256.

## Open Questions

None.

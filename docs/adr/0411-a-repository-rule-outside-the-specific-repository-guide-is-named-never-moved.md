---
status: accepted
created_at: 2026-11-16T22:40:00Z
updated_at: 2026-11-16T22:40:00Z
deprecated_at: null
superseded_by: null
---

# A repository rule outside the specific-repository guide is named, never moved

A repository's own rules and variations belong only in
`docs/agents/specific-repository.md`. The Baseline text that adopters receive
now says so. `roundfix baseline update` and `roundfix doctor` name every
Misplaced Repository Rule and the file it belongs in. Neither command moves,
rewrites or refuses anything because of one.

## Context

Adopters add their own rules to files the Baseline owns. Some go after or
before a guide's Managed Region, some between the Managed Regions of the root
instruction file. They followed the text they received: the guide templates
said that repository authors "own project-specific extensions outside setup
markers" and that "repository-authored dispatch extensions may add stricter
triggers". A managed refresh preserves those bytes and reports nothing about
them. The fleet therefore keeps its rules in a dozen places, and an update of
a Baseline guide can contradict a rule that sits right next to it. On
2026-11-16 the maintainer chose "Texto + detecção": the rule and an example go
in the guides, and the two commands say which passage to move. The operator's
brief added that nothing is moved automatically.

## Decision

**The rule.** The guide templates of `agent-instructions.md` and
`skill-dispatch.md`, the root pointer to the repository's own rules, and three
existing Normative Clauses state the placement and give an example. The three
clauses keep their identities:

- `clause.core.keep-root-compact`
- `clause.core.activate-matching-skills`
- `clause.context.docs-one-job-per-directory`

The starter text the Baseline writes into a new
`docs/agents/specific-repository.md` shows how a rule is written there.

**What is a Misplaced Repository Rule.** It is one of two things:

- A passage outside every Managed Region of a Markdown carrier the Setup
  Manifest lists, or of a root `CLAUDE.md` that is a regular file. A passage
  is a maximal run of lines between Managed Regions or the file's edges,
  trimmed of blank lines.
- A Managed Region whose bytes differ from the digest the Setup Manifest
  recorded, which is an Unrecorded Managed Region with a digest mismatch.

A block that another tool delimits with its own `<!-- BEGIN:<name> -->` and
`<!-- END:<name> -->` pair is that tool's, not a repository rule. It is
neither reported nor counted. `docs/agents/specific-repository.md` itself is
never inspected.

**Where it is reported.**

- `roundfix baseline update` adds the `misplacedRepositoryRules` field to
  every result it prints when it reads a valid Setup Manifest, and a
  `Misplaced repository rules: <n>` section to its text. The field sits
  outside the Plan Document, so the Plan Digest and the plan goldens do not
  change.
- `roundfix doctor` adds a `baseline-rules` line. It is `ok` with no passage,
  `warn` with the passages and the move as its next action, and `skipped`
  without a Setup Manifest. It never fails, so the exit code does not change.

## Considered Options

- **Move the passages automatically during a managed refresh.** Refused by
  the operator's brief. A moved passage changes meaning without its heading
  and its neighbours. Only a person or a supervised Agent can merge it into
  existing rules.
- **Refuse the update (exit `3`) until the passages move.** Refused. It would
  block every adopter's next update on text the Baseline itself allowed until
  now. The maintainer asked for the command to say which passage to move.
- **Carry the passages in the Plan Document.** Refused. They are not part of
  the plan, the plan never changes them, and binding them into the Plan
  Digest would make an unchanged plan look changed whenever a passage is
  edited.
- **A managed how-to block inside `docs/agents/specific-repository.md`.**
  Refused. The file is the repository's (ADR-0203), and a Managed Region there
  would turn the one place a repository owns into a place the Baseline
  rewrites.
- **A new Normative Clause.** Refused. A new clause in a rule the Standard
  TypeScript profile requires needs a Source Baseline row. Rewording the three
  existing clauses keeps their identities and their enforcement.

## Consequences

- Adopters that kept rules outside the markers see them listed on their next
  update and in Doctor. Their bytes stay where they are until someone moves
  them.
- A Doctor run in an adopted repository reads the carriers the Setup Manifest
  lists, offline and without writing.
- A repository that disabled `repository.extension.enabled` is told to enable
  it before moving a passage, because the destination does not exist there.
- An Unrecorded Managed Region is still refreshed as ADR-0145 decides. Its
  removed lines are now also named as a passage to move, before the refresh
  drops them.

**Amended by [ADR-0432](0432-claude-md-is-always-a-symbolic-link-to-agents-md.md) (2026-11-18).** A regular root `CLAUDE.md` whose bytes are exactly `AGENTS.md`, with no newline, yields no Misplaced Repository Rule, nor does a passage whose trimmed text is exactly `@AGENTS.md`. Every other passage is still named and never moved.

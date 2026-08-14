---
spec: 0115-an-archive-that-survives-its-own-move
status: active
created: 2026-08-14
surfaces: [backend, cli]
---

# An archive that survives its own move

Archiving a Spec breaks it in three ways that have nothing to do with whether the
Spec was finished. Its relative links stop resolving, because the move adds a
directory level and the contract forbids rewriting an archived Spec — so a link
that was correct while the Spec was active is wrong forever after. Its destination
depends on which Spec Root is configured rather than on the repository's choice:
the built-in root archives outside itself while every other root archives inside,
so a repository with both conventions has no configuration that reaches the one it
wants. And the archive happens before continuous integration can answer, so a CI
refusal arrives when the Spec is already history and the only way to add a
corrective Task is to undo an archival the contract has no command for.

## Project Constraints

- Identifier strategy: applicable — Archive Command, Spec Root and History Root
  are glossary terms this Spec changes the rules of, and un-archiving would be a
  new act the glossary must name. The closing node checks whether the work
  introduced or changed a term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential or
  request is created or read. Reading a continuous-integration result, if the
  design chooses to wait for one, happens through the existing command surface
  rather than through anything this Spec adds. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0120 fixes the history root and the
  built-in destination this Spec makes uniform; ADR-0091 makes the QA gate a Task
  node of its own type, required to be terminal, which is why the graph settles
  before any external signal exists and archival follows it; ADR-0096 adds the
  Daemon-owned mechanical stage inside that node, and ADR-0117 places a check with
  the stage that can produce its defect — both apply unchanged. ADR-0104 makes a
  Spec accept on evidence it did not author, which is the argument for treating a
  continuous-integration result as part of that evidence. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the spec and CLI packages plus their
  tests. Source: `docs/agents/agent-instructions.md`.

## Goals

1. A link authored inside a Spec still resolves after the Spec archives.
2. The archive destination is the repository's choice, not a consequence of which
   root name it happens to use.
3. A Spec is not history before the evidence that can still refuse it has spoken,
   or the loop has a supported way back.

## User Stories

1. As a reader following a link inside an archived Spec, I want it to resolve, so
   that the archive stays a usable record rather than a snapshot of broken paths.
2. As a maintainer configuring the default Spec Root, I want the archive inside it
   if that is what my guides say, so that I do not have to rename my Spec Root to
   escape a special case.
3. As a maintainer whose repository already archived under the other convention, I
   want the tool to migrate it, so that one `ls` is again the complete record.
4. As a Supervisor whose CI refused a Spec that already archived, I want a
   supported way to reopen it, so that adding a corrective Task is not three
   manual edits against a contract that forbids touching archived Specs.

## Core Features

1. **The move rewrites link destinations.** Archiving rewrites the Markdown link
   destinations inside the Spec so they resolve from the new depth, changing
   nothing else — the one edit compatible with preserving every observation byte
   for byte.
2. **One destination rule for every Spec Root.** The built-in root stops being a
   special case, so the archive location follows one rule whatever the root is
   called.
3. **Adopters on the other convention are migrated.** The Baseline already
   migrates layouts and carries this one, so a repository holding a split archive
   is healed by running the tool.
4. **Archival and refusal are ordered deliberately.** Either archiving waits for
   the evidence that can still refuse the work, or reopening an archived Spec is a
   supported act with its own command. Which of the two is the decision this Spec
   settles.

## Non-Goals / Out of Scope

- Changing what the Archive Command verifies before it moves anything.
- Changing the rule that an archived Spec stays byte-identical afterwards; the
  link rewrite happens during the move, not after it.
- Changing the QA gate's position in the graph, which an accepted decision owns.
- Rewriting links in Specs already archived.

## Success Metrics

- A Spec containing relative links to findings and decision records archives with
  every link still resolving, proven by resolving them after the move.
- One destination rule holds for the built-in root and for a configured one,
  proven by archiving under both.
- A repository holding a split archive is healed by running the Baseline, measured
  against a repository this Spec did not build — two were observed on 2026-08-12
  and 2026-08-13 with sixteen and six Specs under the older convention.
- A Spec refused by CI after archival is reopened by a supported act rather than
  by editing frontmatter and moving folders by hand.

## Decisions

- The link rewrite happens during the move. Rewriting after the move would edit an
  archived Spec, which the contract forbids; rewriting during it is the same act
  that changes the depth.

## Open Questions

- Whether archival waits for continuous integration or gains a reverse operation.
  Waiting keeps the archive honest and delays the loop's terminal step; reversing
  keeps the order and admits that archival is not final. The second is smaller and
  the first is what the evidence contract implies.
- Which destination rule wins. Archiving inside the Spec Root matches what most
  repositories' guides already say and breaks the repositories that adopted the
  built-in behaviour; the reverse breaks the others. Either way the Baseline
  migration decides how much that costs.

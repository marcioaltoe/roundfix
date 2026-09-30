---
spec: 0192-owned-skills-that-describe-the-product-as-it-is
status: active
created: 2026-09-30
surfaces: [docs]
---

# Owned skills that describe the product as it is

An audit on 2026-09-30 compared the fourteen Roundfix-owned skills and the
user guide with the v0.21.0 binary and the Baseline clauses. Twelve skills and
five guide files state something the product no longer does, or omit something
it does:

- The Roundfix skill has no section for five commands the CLI lists, omits one
  event category, and describes a manual delivery order without the pre-PR
  review.
- `write-prd`, `write-techspec` and `write-tasks` do not teach the
  authorization record a Delivery Queue requires from every Spec, the
  `## Unreachable Acceptance` section, or the coverage units the checker
  enforces.
- `qa-gate` tells the QA Agent to run a Verification the Daemon already ran,
  and its report template lacks the column the eligibility policy reads.
- `archive-spec` teaches stamping an archive by hand, which the docs layout
  guide forbids.
- `write-idea`, `write-prd`, `brainstorming` and `business-analyst` show a
  question block with an `Other` option, and `council` closes with one
  compound open question. The structured-question clause in
  `docs/agents/agent-instructions.md` asks for one question with two or three
  options, the recommended one first, and no `Other` option.
- The user guide documents a flag the CLI refuses, has fifteen links that
  resolve to nothing, and states a skill count Doctor does not print.

An Agent that follows these skills authors a Spec the queue refuses, or a
report the gate rejects. Several corrective Tasks in Waves 4 to 6 came from
that gap.

## Project Constraints

- Identifier strategy: not applicable — no identifier is created or changed;
  the Spec edits skill and guide text and adds repository tests. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Markdown files and
  repository tests only; no credential is read and no network call is added.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0057 makes the Daemon the only
  writer of a Task's status during an Implement Run, which the task template
  must say. ADR-0167 keeps the pre-PR Pull Request row from deciding a
  qualifying partial by reading a provenance column, which the QA report
  template must carry. ADR-0179 lets an explicit `paths: []` grant operations
  only, and ADR-0156 makes Success Metrics and API Contracts coverage units;
  the authoring skills must teach both. ADR-0116 checks a claim about an ADR
  against the ADR's text, which `write-tasks` must stop denying. This Spec
  changes none of these decisions; it aligns text with them. Its gate is bound
  by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104 and ADR-0155, and
  ADR-0093 and ADR-0117 check its consistency by citation at the authoring
  stage. ADR-0176 reads citations only from authored text. ADR-0166 records a
  path a Task changed without declaring it; every Task here declares its
  paths. ADR-0097 cites ADR-0080 but carries a QA row
  forward, ADR-0160 cites ADR-0057 but governs a red repository gate, ADR-0168
  cites ADR-0116 but narrows the related-ADR check, and ADR-0170 cites ADR-0057
  but governs Task Carry-Forward. This Spec touches none of those behaviors,
  so they do not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer stated on 2026-09-30
  "considere autorizado a ajustar todas as skills se necessário", an express
  maintainer authorization recorded in [_authorization.md](_authorization.md);
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, `.agents/skills/write-prd/SKILL.md`,
  `skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `skills/write-techspec/references/techspec-template.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`,
  `.agents/skills/implement-spec/SKILL.md`, `skills/implement-spec/SKILL.md`,
  `.agents/skills/council/SKILL.md`, `skills/council/SKILL.md`,
  `.agents/skills/council/references/archetypes.md`,
  `.agents/skills/write-idea/SKILL.md`, `skills/write-idea/SKILL.md`,
  `.agents/skills/write-idea/references/opportunity-scan.md`,
  `.agents/skills/business-analyst/SKILL.md`,
  `skills/business-analyst/SKILL.md`, `.agents/skills/brainstorming/SKILL.md`,
  `skills/brainstorming/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- An Agent that follows an owned skill produces a Spec the Delivery Queue
  accepts and a QA report the eligibility policy reads, and asks the
  maintainer questions in the form the Baseline requires.
- Every command the CLI lists is documented in the Roundfix skill and in the
  user guide.
- A reader of the user guide meets no refused flag, no dead link and no count
  the product does not print.
- A command added without its documentation, or a guide link that stops
  resolving, fails a repository test before a Pull Request opens.

## Core Features

1. **The Roundfix skill covers the CLI.** The skill documents `window`,
   `qa-report accept`, `baseline capabilities check`, `baseline profile init`
   and `init`, lists `agent-selection` among the event categories, states the
   delivery order the Delivery Queue executes, and says that a Task added
   after a completed QA gate needs `roundfix reopen` first.
2. **The user guide states what ships.** The guide and the README name no
   `--qa` flag on `implement`, state the archive destination the command
   uses, document `spec audit` and `storage report`, show the Doctor skills
   line as a shape whose numbers depend on the repository, and carry no link
   that resolves to nothing.
3. **The authoring skills teach a deliverable Spec.** `write-prd`,
   `write-techspec` and `write-tasks` teach the authorization record with its
   `operations` and `paths` lists, including `paths: []`. They teach the
   `## Unreachable Acceptance` section, the four coverage units, the delivery
   order, `roundfix reopen`, the Daemon as status writer, and how to verify
   code for another operating system.
4. **The gate and lifecycle skills agree with the Daemon and the commands.**
   `qa-gate`, `archive-spec`, `setup-context-driven`, `implement-spec`,
   `council` and `write-idea` lose each statement the product contradicts.
5. **The skills ask one structured question at a time.** Wherever an owned
   skill shows how to ask the maintainer, it shows one question with two or
   three mutually exclusive options, the recommended one first and labelled
   `(Recommended)`, and no `Other` option. `council` captures its decision
   with two such questions, asked one after the other.
6. **A documentation contract guards the class.** A repository test fails when
   a command the root help lists is named by neither the Roundfix skill nor
   the user guide, or when a relative link in the user guide or the README
   does not resolve.
7. **A changed skill declares a higher version.** Each owned skill whose text
   changes raises both of its version fields by one patch step, except the
   Roundfix skill, which stays at `0.0.2`: three existing `./skills` tests read
   the owned-skill floor through its version line and fail when it moves
   (measured on 2026-09-30). Spec 0195 makes the floor follow the bundle and
   raises it.

## Non-Goals / Out of Scope

- Restructuring, splitting or shortening any skill or guide. Spec 0194 does
  that, after this Spec.
- Editing a Baseline module under `internal/baseline/assets/` or a generated
  guide under `docs/agents/`. Spec 0193 does that.
- The adapter version named in the Roundfix skill's `agents/openai.yaml`.
  Spec 0189 updates the adapter pins and that line with them.
- Changing any CLI behavior, help text, exit code or checker rule.
- Editing a governed contract test. Every phrase those tests require stays in
  the skills.
- The `### QA settlement` section of any skill, which a contract test keeps
  byte-identical across three skills.

## Success Metrics

1. Success Metric: every command path the root help lists is named, as
   `roundfix <path>`, in the Roundfix skill and in the user guide. Removing one
   command's text makes the documentation contract fail and name that command.
2. Success Metric: no relative link in `docs/user-guide/*.md` or `README.md`
   points to a missing path, and neither names a `--qa` flag other than
   `--qa-override`.
3. Success Metric: each stale statement this Spec names is gone from its skill
   or guide and its replacement is present, checked phrase by phrase.
4. Success Metric: the existing skill contract tests pass without an edit, the
   `skills/` mirrors are byte-identical to `.agents/skills/`, and
   `make baseline-digests` reports no change.
5. Success Metric: every owned skill this Spec changes, except the Roundfix
   skill, declares in both version fields a version one patch step above the
   one on the starting main; the Roundfix skill keeps `0.0.2` in both.

## Recorded limits

- The documentation contract proves that a command is named, not that its
  description is correct. Accuracy of each description stays a review and QA
  concern.
- `archive-spec` keeps its Unarchive section. No command reverses an archive,
  so that section still describes a manual edit.
- `archive-spec` keeps naming `docs/_inbox/`. The docs layout guide and the
  adoption index still define that source type, so the audit's claim that the
  mention is stale was wrong.
- The skills keep naming `docs/agents/backend.md`, which a contract test
  requires for repositories with a backend guide. They gain a sentence for a
  repository without one.
- Until Spec 0193 lands, `docs/agents/autonomous-work.md` still states the
  order archive-then-review. The skills state the order the Delivery Queue
  executes.
- The corrective-Task ceiling the skills state stays at two.
- No test reads a skill's question examples for the structured-question form.
  This Spec checks the five examples it fixes by phrase; a later skill can
  reintroduce the pattern.

## Decisions

- **Name the command, literally.** The contract requires the literal
  `roundfix <command path>` in the text. A literal match is cheap, cannot be
  satisfied by a paraphrase, and is what a reader searches for.
- **Read the command list from the CLI.** The contract parses the root help's
  usage lines, so no second list of commands exists to drift.
- **Show the Doctor line as a shape.** The skill counts depend on the
  repository's Baseline Profile, so a literal count in a guide goes stale when
  the Profile changes.
- **Keep every phrase a contract test requires.** A skill edit that would
  remove a required phrase rewrites around it instead of editing the governed
  test.
- **The interface supplies the custom answer.** A skill example never lists an
  `Other` option. The structured question tool adds the free-text answer
  itself, and a chat fallback accepts one without an option for it.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- The CLI's own output: `roundfix --help` lists fifty-one command paths, and
  `roundfix implement --spec <slug> --qa` answers
  `unknown flag "--qa"`. Both were read from the v0.21.0 binary on 2026-09-30.
- Git's own test suite checks the same class: `t/t0450-txt-doc-vs-help.sh`
  asserts that each builtin's documented synopsis agrees with its `-h` output
  (<https://github.com/git/git/blob/e9019fca/t/t0450-txt-doc-vs-help.sh>).

## Research basis

The scope comes from a read-only audit of the owned skills on 2026-09-30,
re-verified item by item against this worktree. Two audit items were wrong and
are recorded under Recorded limits.

- Secondbrain: `wiki/index.md`, `wiki/sources/agent-skills-skillsmd-guide-2026.md`
  and `wiki/entities/skills.md` were read, with the query
  `qmd query "skills desatualizadas documentação diverge do comportamento do CLI drift de skill"`.
  The skills guide describes reference files loaded on demand, which supports
  leaving restructuring to Spec 0194. The research entry
  `inbox/secondbrain/2026-09-30-orquestradores-e-ecossistema-jev-o-que-se-transfere-ao-roundfix.md`
  records small skills as a mechanism worth adopting and says nothing about
  drift detection.
- Exa located Git's `t0450-txt-doc-vs-help.sh` and three projects that test
  their documentation against `--help`. Two of them add a falsifiability
  check, which this Spec copies as an idea: the parser must return a minimum
  number of commands, so a changed help format cannot pass silently.
- A scratch clone measured the regeneration: editing all ten skills this Spec
  touches, then running `make skills-sync` and `make baseline-digests`,
  changed thirty skill files and reported `"changed":false` for derived
  artifacts.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.

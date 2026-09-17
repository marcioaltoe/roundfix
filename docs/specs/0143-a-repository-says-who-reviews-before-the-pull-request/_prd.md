---
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
status: active
created: 2026-09-17
surfaces: [backend, cli, docs]
---

# A repository says who reviews before the Pull Request

Every delivery in this repository is supposed to pass an independent review
before its Pull Request opens. The obligation is written as a mandatory clause
in the agent guides, which name the supported choices — `codex`, `claude`,
`coderabbit`, or an explicit `none` — and say the choice must be resolved before
publication.

Nothing resolves it. Roundfix has no configuration key for a pre-Pull-Request
reviewer, no precedence between a user's choice and a project's, and no surface
that reports which one applies. The only review setting it carries,
`review_source`, configures the Review Source that reads a Pull Request after it
exists, which is a different act at a different time.

So the choice lives in whoever is working. Two sessions in the same repository
can pick different reviewers and both believe they followed the clause, and a
repository that deliberately wants no pre-Pull-Request reviewer has no way to
say so — its silence is indistinguishable from never having decided.

This Spec makes the repository able to say it, and Roundfix able to report it.
It is the first slice carved from Spec 0126, which keeps the reviewer sessions
themselves, the finding dispositions and the publication evidence.

## Project Constraints

- Identifier strategy: applicable — the provider names `codex`, `claude`, `coderabbit` and `none` are the vocabulary the agent guides already use, and this Spec reuses them without coining or renaming. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — resolving and reporting a policy reads local configuration only; no provider is installed, authenticated or invoked, and no network transport is opened. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the configuration precedence and the read-only support surfaces are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0002 applies: configuration is YAML in User Config and Project Config, and this Spec adds a key inside that contract rather than a new configuration mechanism.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the agent guide that states the obligation and the configuration this repository already ships.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The work is ordinary source in `internal/config` and `internal/cli` plus the user guide; this Spec does not edit `.roundfixrc.yml`, any Baseline asset, or any guide delivered inside setup-context markers, all of which are Governed Paths. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A repository can declare which reviewer runs before its Pull Requests, or
  declare explicitly that none does.
- A maintainer can read the applicable choice and where it came from, without
  guessing and without running a reviewer.
- An unsupported value is refused when the configuration loads, naming the
  supported values.
- Nothing is installed, authenticated or invoked by declaring a policy.

## User Stories

1. As a maintainer, I want to write the pre-Pull-Request reviewer in project
   configuration, so that every session in this repository resolves the same
   choice.
2. As a maintainer of several repositories, I want my own default to apply where
   a project says nothing, so that I do not repeat the choice everywhere.
3. As a maintainer who wants no pre-Pull-Request reviewer, I want to say `none`
   explicitly, so that the omission is a recorded decision rather than silence.
4. As a maintainer diagnosing a repository, I want the applicable policy and its
   source reported, so that I can tell a project choice from an inherited one.
5. As a maintainer who mistypes a provider, I want the configuration refused with
   the supported values named, so that the mistake surfaces before a delivery
   depends on it.

## Core Features

1. **A policy key with the four supported values.** Configuration carries the
   pre-Pull-Request reviewer as `codex`, `claude`, `coderabbit` or `none`,
   separate from the existing `review_source` settings, which keep their meaning
   and their behavior.
2. **Precedence that matches the rest of configuration.** Project Config
   overrides User Config, which overrides the built-in default of `codex`. An
   absent key inherits; it never means `none`.
3. **An applicable policy is reported with its source.** A read-only support
   surface states the resolved provider and whether it came from Project Config,
   User Config or the built-in default, and states explicit `none` as review
   disabled by configuration.
4. **An unsupported value is refused at load.** The configuration error names
   the key, the offending value and the four supported values, and refuses
   before anything reads the policy.
5. **Declaring a policy performs no provider work.** No installation,
   authentication, reviewer session, readiness probe or provider request happens
   because a policy was declared or reported.

## User Experience

A maintainer writes the reviewer in project configuration and runs the
diagnostic command. It states the resolved provider and its source in one line,
beside the other readiness checks. With `none`, the same line states that review
is disabled by configuration. A mistyped value stops the command with an error
naming the supported values.

## Non-Goals / Out of Scope

- Running a reviewer, opening a reviewer session, or probing reviewer readiness.
  Spec 0126 keeps that.
- Recording review evidence or a configured omission for a candidate, and
  deciding what publication may consume. Spec 0126 keeps those too.
- Changing `review_source`, the Review Source watch loop, its request command,
  or anything about a Pull Request that already exists.
- Editing this repository's own `.roundfixrc.yml`, any Baseline asset, or any
  guide delivered inside setup-context markers.
- Changing how the agent guides state the obligation.
- Migrating historical Runs, receipts or their classifications.

## Declared intentional breaks

- A configuration that carries an unsupported value for the new key stops
  loading where it previously carried no such key at all. That is the point of
  refusing at load, and it can only affect a configuration written after this
  Spec ships.

## Regression locks

- `review_source` keeps its keys, defaults, validation and behavior.
- A configuration without the new key loads exactly as it does today and
  resolves to the built-in default.
- The diagnostic command still mutates nothing and still reports every check it
  reports today.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the
agent guide that already states the obligation and names the four supported
values, and this repository's own configuration, which carries `review_source`
but no pre-Pull-Request reviewer key. The replay reads both and shows that the
resolved policy this Spec reports matches the guide's vocabulary and the
repository's actual silence, which resolves to the built-in default. Where the
guide cannot be read, the row records that reason and does not block.

## Success Metrics

1. A fixture repository whose project configuration names a provider reports
   that provider with Project Config as its source.
2. A fixture whose project configuration is silent and whose user configuration
   names a provider reports that provider with User Config as its source, and a
   fixture where both are silent reports the built-in default.
3. A fixture whose configuration carries an unsupported value fails to load with
   an error naming the key, the value and the four supported values.

## Decisions

- **A separate key from `review_source`.** They name different acts at different
  times: one reviews a candidate before a Pull Request exists, the other reads a
  Pull Request that already does.
- **Absent means inherit, never `none`.** A repository that wants no reviewer
  says so; silence keeps the default, which is what the guides already assume.
- **Report, do not run.** This slice ends at the resolved policy and its source;
  spending a reviewer session belongs to the slice that can also record its
  evidence.

## Research basis

**Secondbrain.** Consulted before authoring, index first, then a query on
configuring which reviewer runs and on precedence between user and project
settings. The results were dominated by mirrors of this repository, which are
references rather than independent knowledge, and no source changed the design.

The repository's pending Inbox Entries were read before authoring; none is a
source for this Spec.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read, so no external validation is claimed.

**Local measurement.** The absence of any pre-Pull-Request policy key, the shape
of `review_source`, the configuration precedence and the diagnostic command's
check list were each read in this repository's source before this PRD was
written.

## Open Questions

None.

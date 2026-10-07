---
spec: 0245-generated-records-that-hold-across-specs-and-platforms
status: active
created: 2026-10-07
surfaces: [backend, docs]
---

# Generated records that hold across Specs and platforms

This is a bug fix adopted from the Backlog Entry
[generated records break when Specs are authored in parallel or tested on another platform](references/2026-10-07-generated-records-that-parallel-work-or-another-platform-breaks.md)
of 2026-10-07. Two checked-in records were wrong for the head CI tested, and
the operator rewrote each by hand.

1. A Baseline module's `version` is written by hand when a Spec is authored.
   Specs 0239 and 0241 both took `context-workflow` version 22; only a review
   caught it, and the operator raised 0241 to 23 in PR #438. Nothing checks
   the number, so Specs 0242 and 0243 changed the same module without raising
   it at all.
2. `docs/references/coverage-record.json` lists the tests `go test -list`
   found on the machine that recorded it. Re-recorded on macOS for PR #441, it
   listed the macOS-only `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded`,
   and Linux CI reported a coverage regression. Linux-only store tests appear
   in the other direction.

## Prerequisites

None. Spec 0244 is authored in parallel and may change a Baseline module
before this Spec lands. No Verification here names a module version number,
and the record step this Spec adds is what reconciles such a change at merge.

## Project Constraints

- Identifier strategy: applicable — new names are the record file
  `internal/baseline/module-versions.json` with schema
  `roundfix/baseline-module-versions/v1`, the test
  `TestEveryBaselineModuleVersionIsRecorded` and its flag
  `-record-module-versions`, and the coverage record's new `platforms` and
  `platformTests` fields. Module identifiers and version integers keep their
  scheme. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — both records are read and written
  from local files and the local `go` toolchain (`go list`, `go test -list`);
  nothing opens a network connection. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0249 does not apply: it governs
  how the QA step formats its QA directory and when reopen sees a Late
  Dependency, and neither generated record touches the QA gate. ADR-0250 (this Spec) decides both
  rules, ADR-0250: "the record step chooses it" and ADR-0250: "The coverage
  record lists tests for every release platform". ADR-0189 is the rule this
  applies to modules, ADR-0189: "A version names one content". ADR-0233
  supplies the line-scoped merge, ADR-0233: "may then change a line-scoped
  path only on matching lines". A conflict confined to declared derived paths
  is still resolved by regeneration under ADR-0192, ADR-0192: "runs each
  matched declaration's regeneration command". ADR-0149 has a grant name the
  regeneration command, ADR-0149: "the grant names the command", so
  `_ownership.yml` is not changed: the new record lives outside the scanned
  digest roots. ADR-0130 and ADR-0179 bound the Governed Paths, ADR-0182
  settles each Task on its Verification, ADR-0184 binds Surface Transcripts
  and this Spec declares none, and the gate is bound by ADR-0080, ADR-0091,
  ADR-0096, ADR-0097, ADR-0104 and ADR-0167. ADR-0166 records a path a Task
  changed without declaring it, ADR-0240 decides when a QA partial qualifies,
  and ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this
  Spec by citation, stage and receipt. ADR-0194, ADR-0195 and ADR-0210 cite
  ADR-0097 but decide what a QA row records, when it is observed again and its
  evidence snapshot, ADR-0229 cites ADR-0167 but decides operator archives,
  and ADR-0237 cites ADR-0229 but decides a Delivery Retry after an outside
  merge; this Spec changes none of them, so none applies. ADR-0247 and
  ADR-0248 decide archives and the history sanitize; this Spec changes
  neither, so neither applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — 11 declared files are Governed Paths,
  measured with `GovernedPath` through a `go test -overlay` probe that wrote
  nothing. Express maintainer authorization: "Autorizar os dois" (Baseline
  sources, guides and `.roundfixrc.yml`, 2026-09-30), "Concedo" (the Governed
  Paths each Spec declares) and "Pode seguir nessa ordem" (this Spec as the
  second of the cycle, 2026-10-07). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0245-generated-records-that-hold-across-specs-and-platforms/_authorization.md`;
  bounded files: `.roundfixrc.yml`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`, `internal/baseline/assets/modules/backend.json`, `internal/baseline/assets/modules/cli-surface.json`, `internal/baseline/assets/modules/external-triage.json`, `internal/baseline/assets/modules/monorepo.json`, `internal/baseline/assets/modules/rust.json`, `internal/baseline/assets/modules/tui-surface.json`, `internal/spec/coverage_test.go`.

## Goals

1. A Baseline module edit is recorded under a version the record step
   chooses, and a module edit that was not recorded fails the suite.
2. Two items that change the same module merge with the next free version
   without an operator.
3. The coverage record is the same bytes on any host and lists each
   platform-limited test with its platforms.
4. A test removed on any release platform fails the comparison on every host.

## Core Features

1. The Module Version Record: a record of each module's versions and content
   digests, a check that refuses unrecorded content, and a record flag that
   chooses the version.
2. A platform-neutral Coverage Record: static collection under every
   release `GOOS`, with a toolchain anchor for the host.
3. Derived declarations: the Baseline catalog declaration runs the record
   step first and covers the module version lines; the coverage record gains
   its own declaration.
4. Guidance: the repository rule, the glossary and ADR-0250.

## Non-Goals / Out of Scope

- Rule, guide, clause and root-block versions inside a module stay
  hand-maintained.
- Changing the Makefile, `make baseline-digests` or the derived-ownership
  records; the record step is run by the derived declaration and the
  repository rule.
- Running the suite on Windows or adding a CI platform.
- Changing the owned-skill record of Spec 0228.

## Success Metrics

1. Success Metric: on temporary module directories, a module changed under a
   recorded version is refused without the record flag, and with it is raised
   to one above the highest recorded version while every other byte of the
   module is kept.
2. Success Metric: on this repository the module check passes with every
   catalog module recorded, and a second record step changes no file.
3. Success Metric: in a disposable clone, two branches that each change
   `context-workflow` and run the record step merge, after the declared
   regeneration, on the version one above the first branch's, and
   `make verify` exits 0.
4. Success Metric: the Coverage Record lists
   `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded` under `darwin` and
   the two `TestParseProcStatStartTime...` tests under `linux`, and
   re-recording it reproduces the same bytes.
5. Success Metric: a fixture module with macOS-, Linux- and Windows-only tests
   records each with its platforms on any host, and the host's static
   collection equals `go test -list`.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The operator's intervention log, `~/.roundfix-operator/queue-interventions.md`,
  entries 215 (0241's module version raised to 23 by hand after 0239 took 22)
  and 224–225 (PR #441's coverage record re-recorded on macOS, then a
  macOS-only entry removed by hand for Linux CI).
- The Go command documentation (<https://pkg.go.dev/cmd/go>): a file whose
  name ends `_GOOS` has "an implicit build constraint", `go list` reports
  `TestGoFiles` and `XTestGoFiles`, and `go test` runs "a separate test
  binary" per package. A measurement on 2026-10-07 confirmed the last point:
  `GOOS=linux go test -list` on macOS fails with `exec format error`.
- An authoring measurement on 2026-10-07 at `66f2d85b`: a static collection
  under `darwin`, `linux` and `windows` matched `go test -list` on macOS with
  zero differences and found 48 platform-limited tests (1 macOS-only, 3
  Linux-only, 3 Windows-only, 41 Unix-only). In a scratch clone, splitting the
  six one-line module headers and regenerating changed only the catalog
  snapshots, four plan goldens and the Setup Manifest, and `go test ./...`
  passed.

## Glossary

- adds: **Module Version Record**
- adds: **Coverage Record**

## Decisions

- A module version is chosen by the record step and names one content; see
  ADR-0250.
- The Coverage Record is collected statically for every release platform;
  see ADR-0250.

## Technical candidate

The [_techspec.md](_techspec.md) records the design, the invariants and the
build order.

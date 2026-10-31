---
spec: 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
status: active
created: 2026-10-31
surfaces: [backend, docs]
---

# The untasked-skills check reads the files the fingerprint reads

Minimal PRD for a bug fix, minted by `write-techspec`. `roundfix spec check`
reports `SC-SKILLS-UNTASKED` for a Task that changes only a Go test file under
a Behavior Surface's source pattern, although the fingerprint of that surface
never reads a Go test file and so cannot move. The `config keys` surface
declares the pattern `internal/config/*.go`, which also matches
`internal/config/*_test.go`, and Spec 0332's task_01 had to add a Coverage
Review it did not need to pass the strict check. After this Spec the check and
the fingerprint read the same file set through one shared predicate: a Go test
file never asks for a skills Task, and a source file still does.

## Sources

One Backlog Entry, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-10-30-the-untasked-skills-check-counts-test-files-the-fingerprint-ignores.md`](references/2026-10-30-the-untasked-skills-check-counts-test-files-the-fingerprint-ignores.md),
   `type: fix`, reported by the author of Spec 0332 on 2026-10-30 and
   confirmed by the operator at `35de7f74`.

At the authoring base `docs/backlog/` holds no other entry and
`docs/findings/` no unresolved Finding on the skills rule or the Behavior
Surface fingerprint.

## Prerequisites

None in `requires`. Authored against `0c27cdc9`, where Specs 0332 to 0335 are
active plans. They deliver before this Spec, which runs in its own queue after
Spec 0335 merges; no Verification here depends on their code. Specs 0333,
0334 and 0335 raise the Roundfix owned skill version and re-record
`docs/references/behavior-surfaces.json`; Spec 0334 edits
`docs/user-guide/commands/spec.md` and Specs 0334 and 0335 edit
`.agents/skills/roundfix/references/spec.md`. The Task here edits other
paragraphs of those files after those merges, raises the version through the
record command and re-records the fingerprints last; it pins none of their
bytes. Spec 0334 changes `internal/speccheck/constraints.go` and
`internal/speccheck/report_test.go`, which this Spec does not touch. Spec
0332's task_01 records a Coverage Review on the `config keys` entry of
`docs/references/skill-coverage.json` because of this defect; that review
stays, and this Spec does not edit the map.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed; the finding code and its text keep their bytes.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the check reads a local repository. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0256 defines a Behavior Surface's
  sources as "the source paths whose change can change it" and the authoring
  rule this Spec corrects; the fix restores that definition for Go test
  files, so no new ADR is written. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AS batch 2 answers of 2026-10-31 quoted in the authorization record, which
  keep the cycle AS answers and the cycle AJ grant for the Governed Paths
  each Spec declares; the governed files are the Roundfix skill's `spec`
  reference and its `SKILL.md` pair. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A Task whose only changed files under a covered surface's sources are Go
  test files gives no `SC-SKILLS-UNTASKED`.
- A Task that changes a non-test source of a covered surface still gives
  `SC-SKILLS-UNTASKED`, with its code, summary and fix unchanged.
- The rule and the fingerprint read their file set through one predicate, so
  they cannot drift apart again.
- The spec command guide and the Roundfix skill's spec reference say that a
  Go test file is never a Behavior Surface source.

## Core Features

1. One predicate. The Skill Coverage Map package answers whether a
   fingerprint can read a file: a Go test file, whose name ends in
   `_test.go`, never; every other file, yes.
2. The check uses it. The map's lookup from a declared path to its surfaces
   returns no surface for a file the predicate refuses, so the untasked-skills
   rule never binds a Go test file to a surface.
3. The fingerprint uses it. The `config keys` fingerprint skips the files the
   predicate refuses, in place of its own `_test.go` filter.
4. The rule in writing. The spec command guide and the Roundfix skill's spec
   reference state it beside the skills rule.

## Non-Goals / Out of Scope

- Excluding `testdata` directories or generated files: no fingerprint
  excludes them today, and the guide fingerprints read every Markdown file
  under `docs/user-guide/`, so such an exclusion would itself open a gap
  (`_techspec.md` → Current behavior).
- Refusing a Skill Coverage Map source that names a Go test file; none
  exists, and the lookup already returns no surface for one.
- Editing `docs/references/skill-coverage.json` or removing the Coverage
  Review Spec 0332 records.
- Changing the release plan's `skill-coverage` check, which compares
  fingerprints and never reads source patterns.
- Any text inside a `### QA settlement` section; the Makefile, the CI
  workflows, `go.mod`, `go.sum`; any Baseline module or guide under
  `docs/agents/`.

## Success Metrics

1. Success Metric: a Task declaring only `internal/cli/example_test.go`
   (matched by a source pattern `internal/cli/example_*.go`) or only
   `internal/cli/sub/other_test.go` (matched by a source directory
   `internal/cli/sub/`), under `interface:`, `creates:` or `deletes:`, gives
   no `SC-SKILLS-UNTASKED`; at `0c27cdc9` each gives one.
2. Success Metric: a Task declaring `internal/cli/example_testing.go` or
   `internal/cli/sub/other.go` still gives exactly one `SC-SKILLS-UNTASKED`
   naming the declared path, and so does a graph whose test-only Task sits
   beside a Task declaring `internal/cli/example.go`.
3. Success Metric: in this repository, a Task declaring only
   `internal/config/config_test.go` gives no `SC-SKILLS-UNTASKED` for the
   `config keys` surface, and a Task declaring `internal/config/config.go`
   still gives it; at `0c27cdc9` both give it.
4. Success Metric: the `config keys` fingerprint and the untasked-skills rule
   call the same exported predicate, and the fingerprint test that refuses to
   parse a Go test file still passes.
5. Success Metric: `docs/user-guide/commands/spec.md` and both copies of the
   Roundfix skill's spec reference state that a Go test file is never a
   Behavior Surface source.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not write.

- The report of Spec 0332's author, quoted in the adopted Backlog Entry and
  written before this Spec existed.
- The Go command documentation (`cmd/go`,
  https://pkg.go.dev/cmd/go, read through Exa on 2026-10-31), section
  "Compile packages and dependencies": "When compiling packages, build
  ignores files that end in '_test.go'." The command help fingerprints come
  from the built `roundfix` binary, so no Go test file reaches them.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files and
temporary repositories.

## Validation

The decisions were checked on 2026-10-31 against these sources.

- Characterization at `0c27cdc9` with that commit's `bin/roundfix`, in a
  scratch clone with two committed fixture Specs: one Task declaring only
  `internal/config/config_test.go` and one declaring only
  `internal/config/config.go` each report `SC-SKILLS-UNTASKED` for
  `config keys`.
- In the same clone a Go test file under `internal/config` that declares a
  struct with a new `yaml` tag leaves `TestTheSkillCoverageMapIsCurrent`
  passing, and the same struct in a non-test file makes it report
  `config keys: stale fingerprint`.
- Every source pattern in `docs/references/skill-coverage.json` is an exact
  file except `internal/config/*.go`, the only one that can match a Go test
  file today.
- A prototype of the Task's code in a second clone passed the tests
  `_techspec.md` → Testing Approach names, and each sabotage
  `_techspec.md` → Testing Approach lists failed its named test;
  `TestCoverageEquivalence` stayed green with the added tests.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix skill coverage behavior surface fingerprint test files SC-SKILLS-UNTASKED" --all --files --min-score 0.3`
  (the adopted Backlog Entry's mirror, the Roundfix `spec` and `skills`
  guides, the `release` reference and Spec 0294's Archive Record). None
  records a decision on which files a surface's sources include.
- Exa: the Go command documentation quoted above, and the Go issue
  golang/go#24732, where the Go team states that a `testdata` directory is
  ignored only by package patterns and still builds when named, which is why
  this Spec leaves `testdata` out of the predicate.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AS batch 2 brief for Spec 0336 of 2026-10-31 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-10-31), verbatim:

- Scope: "Entra no AS como 0336" — option text: "Planejo a Spec 0336 agora,
  ela entra na fila depois da 0335 e sai na mesma v0.127.0."
- Tooling: "Makefile, .github/workflows, go.mod / go.sum, Config de
  lint/formatter" — question: "Hoje o grant do ciclo AS proíbe que os Runs
  mexam em Makefile (fora da linha OWNED_SKILLS da 0333), .github/workflows,
  go.mod/go.sum e config de lint/formatter. Para nenhuma Spec travar por
  autorização, quais destes eu libero aos Runs e aos planos de 0332 a 0336,
  sempre que o desenho da Spec pedir?"
- Technical failure: "Consertar e seguir (Recommended)" — option text:
  "Diagnostico, emendo o plano ou crio uma Spec bloqueadora e entrego, sem
  teto fixo de corretivas. Só pulo a Spec se o conserto exigir uma decisão de
  produto ou arquitetura que contradiga o PRD aprovado; aí ela sai do ciclo e
  a release vai com as demais, ou espera se for a 0333."

Taken from those answers, ADR-0256 and the characterization, with no further
interview:

- One exported predicate in the Skill Coverage Map package, read by the
  map's path lookup and by the `config keys` fingerprint, rather than a second
  copy of the `_test.go` filter in the Spec Check.
- The predicate excludes Go test files only; `testdata` and generated files
  stay in, because no fingerprint excludes them.
- One implementation Task, since the code and its documentation share the
  Roundfix skill version record, and a QA gate.

## Open Questions

None.

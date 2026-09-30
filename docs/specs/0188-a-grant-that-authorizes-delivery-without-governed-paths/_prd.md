---
spec: 0188-a-grant-that-authorizes-delivery-without-governed-paths
status: active
created: 2026-09-30
surfaces: [backend, cli, docs]
---

# A grant that authorizes delivery without Governed Paths

On 2026-09-30, `roundfix deliver plan` blocked Spec 0186 with
`authorization refused: paths`. The Spec changes only ordinary files, and the
authorization reader refuses a record whose `paths` has no entry. Listing the
ordinary files instead failed the repository contract
`TestEveryBoundedPathIsGoverned` in the Verification gate of PR #287 (run
36701160969), because every bounded path must be governed (ADR-0130). No record
was valid, so a Spec that changes no Governed Path could only be delivered by
hand.

The rule was introduced by Spec 0119 (#188), when a record's only duty was to
bound Governed Paths. The Delivery Queue later made the same record the grant
of delivery operations, and the two duties conflict for such a Spec.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; the record keeps its
  fields and the reader its reason codes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0179 (this Spec) makes an explicit
  empty `paths` list an operations-only grant. ADR-0130 keeps every bounded
  path governed and holds unchanged. ADR-0160 keeps the frozen authorization as
  the only authority that opens a red repository gate, and an empty grant opens
  nothing governed. ADR-0178 accepts the grant a Task ran under, and an empty
  grant bounds nothing to accept. ADR-0166 records a Task's undeclared paths,
  and every Task here declares its paths. ADR-0167 keeps the pre-PR Pull
  Request row from deciding a qualifying partial; this Spec's gate aims at
  `pass`. ADR-0176 reads citations only from authored text. This Spec's gate is
  bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117,
  ADR-0155 and ADR-0156, and ADR-0093 and ADR-0094 check its consistency by
  citation and artifact presence. ADR-0097 cites ADR-0080 but carries a QA
  row forward, and ADR-0168 cites ADR-0093 but narrows the related-ADR check;
  this Spec changes neither, so they do not apply. All the others hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer asked on 2026-09-30 to
  "Continue até o final para o release de todas as implementações e ajustes";
  the Roundfix skill files ride the standing grant of 2026-09-18 for keeping
  the shipped skills true to the CLI, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A Spec that changes no Governed Path can record a delivery grant that
  `deliver plan` approves and `deliver start` accepts.
- Such a grant never authorizes a governed mutation.

## Core Features

1. **An explicit empty list is an operations-only grant.** A Spec-contained,
   approved record whose `paths` is the empty YAML sequence `paths: []` is
   operative. It grants its listed operations and bounds no Governed Path. A
   record whose `paths` is absent or null is still refused on `paths`
   (ADR-0179).
2. **Nothing governed rides on it.** The Spec checker's tooling detectors,
   Implement preflight and the QA mechanical authorization audit read the
   empty bounded set. A governed change under such a grant is refused exactly
   as a path outside a non-empty grant is.
3. **The guides say how to record it.** `docs/user-guide/commands.md` and the
   Roundfix skill's Delivery queue section state that a Spec with no Governed
   Path records `paths: []`.

## Non-Goals / Out of Scope

- Accepting an absent or null `paths`, a legacy record without paths, or a
  listed path that is not governed.
- Changing ADR-0130's governed set, the repository contract test, the tooling
  detectors, the operations vocabulary or any command's flags and exit codes.
- Editing the setup-owned guides `docs/agents/spec-routing.md` or
  `docs/agents/docs-layout.md`. Neither states a non-empty rule.

## Success Metrics

1. A Spec whose record declares `paths: []` and every delivery operation is
   reported `approved` by `deliver plan`, and `deliver start` accepts it.
2. The same record with `paths` absent or null is still refused on `paths`.
3. A Task commit that changes a Governed Path under a `paths: []` grant is
   refused by the QA mechanical authorization audit.

## Recorded limits

- The explicit empty list is honored only for Spec-contained records. A legacy
  record under `docs/workflow/authorizations/` keeps its reading.

## Decisions

- **Explicit, never inferred.** Only the literal empty sequence grants; an
  omitted field keeps failing, because omission is the common authoring
  mistake.
- See ADR-0179.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The outside-evidence rows rest on sources this Spec did not
produce:

- This session's recorded refusal: `bin/roundfix deliver plan 0186-…` printed
  `authorization refused: paths` on 2026-09-30, and PR #287's Verification gate
  (run 36701160969) failed `TestEveryBoundedPathIsGoverned` when ordinary
  files were listed.
- RFC 6749 §3.3 (<https://www.rfc-editor.org/rfc/rfc6749#section-3.3>)
  separates an omitted scope, which the server must default or refuse, from a
  declared one. GitHub Actions documents `permissions: {}` as an explicit grant
  of no permissions
  (<https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/controlling-permissions-for-github_token>).

## Research basis

The adopted Backlog Entry is indexed in
[references/_index.md](references/_index.md). The Secondbrain was consulted
through `wiki/index.md` and
`qmd query "roundfix authorization record empty paths operations grant delivery"`.
It returned only this repository's mirrors (the 0180 queue authority and
ADR-0160), which add nothing beyond the entry. Exa returned RFC 6749 §3.3
discussions (thephpleague/oauth2-server issue #829, n8n issue #28159). They
confirm that an omitted scope and a declared empty one must be told apart,
which is this Spec's rule.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.

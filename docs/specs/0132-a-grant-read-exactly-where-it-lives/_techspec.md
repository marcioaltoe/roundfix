---
spec: 0132-a-grant-read-exactly-where-it-lives
prd: _prd.md
created: 2026-09-10
---

# A grant read exactly where it lives — Technical Spec

## Executive Summary

Repair three resolution defects and two archive-fragile expectations in the
authorization reader Spec 0119 delivered. The trade-off is deliberate: exactness
at the parse boundary and root-derived paths at every reader, against the
convenience of prefix matching and a hard-coded `docs/specs`. Nothing about what
a grant permits changes.

## Project Constraints

- Identifier strategy: applicable — preserve Spec slugs, Task IDs, refusal codes and record field names; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched; the record is read from the local filesystem and Git objects as it is today. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0130 keeps the Governed Path set monotonic, so no path leaves it; ADR-0057 keeps the Daemon the exclusive writer of Task status; ADR-0117 places each check at the stage that can establish it, and record shape is established at parse time; ADR-0096 requires mechanical facts before the QA Agent turn and is preserved unchanged. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: "Aprovar os quatro caminhos", 2026-09-10, recorded in `docs/specs/0132-a-grant-read-exactly-where-it-lives/_authorization.md`; bounded files: `internal/speccheck/constraints.go`, `internal/speccheck/constraints_characterization_test.go`, `internal/speccheck/governed_repocontract_test.go`, `internal/suiteguardcontract/regeneration.go`, `internal/suiteguardcontract/regeneration_test.go`. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | File | Responsibility |
| --- | --- | --- |
| Record parser | `internal/authorization/authorization.go` | Split frontmatter on complete marker lines and classify a malformed record as a refusal. |
| Operation resolver | `internal/spec/authorization.go` | Derive the consuming record's path and revision from the Spec repository root, separately from the project root the audit judges. |
| Citation resolver | `internal/speccheck/constraints.go` | Resolve a Spec-relative reference against the artifact that carries it. |
| Discovery filter | `internal/suiteguardcontract/regeneration.go` | Recognize the record naming already in use. |
| Archive-stable expectations | `internal/authorization/authorization_test.go`, `internal/spec/spec_test.go`, `internal/speccheck/governed_repocontract_test.go` | Read a record wherever it lives and assert over the whole corpus. |

## Implementation Design

### Exactness at the parse boundary

The parser currently finds the closing frontmatter marker by prefix, so a line
such as `---evil` closes the block and the rest of that line is discarded as
body. Every other field can then be well-formed and the record resolves to a
grant. Require the delimiter to be a complete line: the marker and nothing else,
ignoring trailing carriage returns. A line that begins with the marker and
carries anything more makes the record malformed, and malformed is a refusal
naming that line, never a grant.

This is the one place in the Spec where a near-miss on format previously became
authority, which is why it is stated as an exactness rule rather than a stricter
pattern: the reader must reject what it cannot fully account for.

### Two roots, named separately

An external Spec Root creates two repository identities that the current
readers collapse into one `RepoRoot`, and collapsing them is why neither
direction works: reading the record needs the Spec repository, while judging
the change needs the project repository.

| identity | what it answers | used by |
| --- | --- | --- |
| Spec repository root and revision | where `_authorization.md` lives and which committed bytes it has | the record reader and the citation resolver |
| Project repository root and delivery target | which paths a commit changed, whether each is governed, and whether the grant precedes the consuming commit | the changed-path audit and the ancestor check |

Both must be explicit inputs. Using the Spec root to validate bounded paths
judges the wrong tree; using the project root to read the record leaves an
external record unresolved. A reader that takes one root can only be wrong in
one of the two ways.

The cross-repository provenance rule follows from that split: the record's
committed bytes are read at the Spec repository's revision, and the bounded
paths it declares are repository-relative to the **project** repository, because
that is the tree whose changes the audit judges. A record in an external Spec
repository therefore grants over project paths, and its own revision proves only
that the grant existed, never what changed. When the Spec Root resolves inside
the project repository, both identities are the same value and every current
answer is unchanged.

### Paths derived from the resolved root

Two readers assemble the record path from the constant `docs/specs`. When
`specs.root` resolves elsewhere, the consuming record is at
`<spec-root>/<slug>/_authorization.md` and those readers look in the wrong place,
so Implement and Settle refuse valid work unless a duplicate record is created
inside the code repository. Both must take the resolved root as an input and
derive path and revision from it.

The citation resolver has the mirror defect. For an external root,
`readConstraintArtifact` produces a display path that reaches outside the
repository, and joining a Spec-relative link to it yields a path the containment
check rejects. Resolve the reference against the artifact's own location rather
than against the repository root, and keep the containment check anchored to the
Spec Root rather than to the code repository.

### Expectations that survive the archive

Archiving moves a Spec from the active root to the archive root, which is the
last step before a Pull Request. Two expectations pin the active path and break
there. They must discover the record — active first, archive second — and the
corpus invariant must assert over active and archived Specs together rather than
requiring at least one active Spec, because the last Spec in a queue always
empties that set.

## Coverage Map

- PRD Goal 1 → Record parser.
- PRD Goal 2 → Operation resolver, Citation resolver.
- PRD Goal 3 → Archive-stable expectations.
- User Story 1 → Record parser.
- User Story 2 → Operation resolver, Citation resolver.
- User Story 3 → Archive-stable expectations.
- Core Feature 1 → Record parser.
- Core Feature 2 → Operation resolver.
- Core Feature 3 → Citation resolver.
- Core Feature 4 → Archive-stable expectations.
- Core Feature 5 → Discovery filter.

## Testing Approach

Focused tests at the named seams, plus real filesystem layouts for root
resolution. Required observations:

1. A closing line carrying extra characters refuses and names that line; a
   well-formed record still grants exactly what it granted.
2. A Spec Root outside the code repository resolves its record through the
   operation resolver and through the citation resolver; the default root is
   unchanged.
3. The existing record named `2026-09-08-authorized-qa-archive-override.md` is
   discovered, and discovery still ignores files that cannot carry a record.
4. Every test that reads a Spec's record passes against both the active and the
   archived layout, exercised by moving a fixture Spec between them.
5. Today's parse and resolution answers are recorded before the change, so the
   two intended moves are visible in the diff and nothing else moves.

Observation 3 rests on evidence this Spec did not author: a record written on
2026-09-08 under the previous naming, by work that predates this Spec. If that
record is absent, the row records blocked with that reason.

## Build Order

1. Characterize today's parse, resolution and discovery answers (depends on: none).
2. Exact frontmatter delimiter (depends on: 1).
3. Operation resolver derives path and revision from the resolved root (depends on: 1).
4. Citation resolver resolves against the carrying artifact (depends on: 1, 3).
5. Discovery filter recognizes the naming in use (depends on: 1).
6. Expectations survive the archive (depends on: 2, 3, 4, 5).
7. Terminal QA (depends on: 6).

## Risks & Considerations

The risk is over-refusing: an exactness rule that rejects records which are
legitimately well-formed, or a root-derived path that breaks the default layout.
The characterization holds both — the default root's answers must not move, and
every currently granting record must still grant.

## Decisions

- Same branch as Spec 0119, because these defects are in code that has not
  reached the target and merging first would place a record that grants on
  malformed input into the default branch.
- No change to the Governed Path set, to what a grant permits, or to the
  operation vocabulary.

## Vocabulary Contract

No token is coined. The refusals reuse the existing authorization reason codes,
whose glossary owner is unchanged.

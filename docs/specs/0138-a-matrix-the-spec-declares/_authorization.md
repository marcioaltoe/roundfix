---
status: approved
granted: 2026-09-14
action: make the qa Task's declared Requirements the QA gate matrix and bound the gate's default derivation
consuming: 0138-a-matrix-the-spec-declares
paths:
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0138

On 2026-09-14 the maintainer chose to narrow the QA gate before remodelling Spec
authoring. The same day, asked to approve this record with its exact paths,
operations and regeneration, the maintainer answered "Aprovar como proposto".
That approved the bounded scope above.

## Why a governed path is unavoidable

The matrix rules live in the qa-gate skill. Section 1 stops flow rows after any
audit problem. Section 2 derives rows from every
promise and exclusion. Row input declaration lets an executor set a row's inputs
after the row has run. The skill is a Roundfix-owned Skill, so changing its text
needs an express grant. Every Run executes from the canonical copy, and the
binary embeds the distributed mirror, so both paths are bounded here.

The QA contract in the gate prompt is ordinary source that no authorization has
bounded, so it needs no grant.

## Approved bounded mutation

Edit the canonical qa-gate skill so that:

- each numbered `qa` Task Requirement that starts with `MUST verify` or
  `MUST run` is exactly one row, and such Requirements are the complete matrix;
- an undeclared matrix derives rows from a bounded default;
- no row is an aggregate of other rows or re-checks a Mechanical Refusal Code;
- once the matrix exists, a finding blocks only the rows that depend on it;
- row inputs are fixed when the row is planned;
- a timeout or intermittent failure of the repository Verification is recorded
  as a failure, never an environment block;
- a declared row that names a non-waivable source covers it, and PRD
  Unreachable Acceptance declarations keep their rows;
- an analyzer the `qa` Task names runs over changed packages, and diagnostics
  identical on the delivery target are attributed to their owner.

After the canonical edit, regenerate the distributed mirror with
`make skills-sync`. The mirror is already a bounded path. If any derived pin
changes as a result, rewrite it only with `make baseline-digests`.

## Sanctioned regeneration

The repository-owned command resolves its own generated outputs. This
declaration records the digest regeneration approved above and adds no source
paths.

```yaml
command: make baseline-digests
```

## Limits

- No action, operation or path beyond those above.
- The tooling-audit rules keep their behavior; the gate still audits Task commits
  by command.
- No edit to verdict rules, typed blocked causes, QA Report keys, report naming,
  the Pull Request row's equivalent-evidence path, the outside-evidence
  obligation or the frontend sweep. The only change there is the conflicting
  sentence about blocking flow rows.
- The skill keeps its version.
- No Baseline module, authoring skill, linter, analyzer or Verification
  configuration edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.

---
status: approved
granted: 2026-09-19
action: let the Baseline Profile declare the incremental verification tier beside the complete gate
consuming: 0148-a-profile-that-declares-both-tiers
paths:
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0148

Two mandatory clauses tell an Agent to use "the active Baseline Profile's
declared incremental verification command", and the Profile has no such field.
Asked to approve this record, the maintainer authorized the Baseline assets, the
guide template and the generated repository guides on 2026-09-19. The Roundfix
skill and its mirror ride on the standing authorization of 2026-09-18.

## Why a governed path is unavoidable

The Profile that carries the decision is a Baseline asset: code-generator
configuration, and therefore governed.

Nothing else this Spec changes is governed. The Profile reader and its tests are
ordinary source, and the derived catalog digests and plan characterizations the
sanctioned regeneration rewrites are ordinary source too.

## What this record no longer covers

The maintainer originally authorized the guide template, the template index,
this repository's decision record, the formatter golden and the generated
guides, so the Spec could also publish the decision to the guides. Delivery
showed that publication is a larger change than the slice could carry: five
Runs, three widenings, and a reader whose projections emptied under the full
package run.

On 2026-09-19 the maintainer chose to split the Spec, and this record is
narrowed to the one path the remaining scope needs. Narrowing removes authority
and needs no further approval; publishing the decision to the generated guides
returns to Spec 0121, with its blast radius measured by running the sanctioned
command before it is authored again.

## Sanctioned regeneration

```yaml
command: make baseline-digests
```

```yaml
command: make skills-sync
```

## Limits

- No action, operation or path beyond those above.
- No change to a Baseline module's clause text; this Spec makes the declared
  value exist, it does not rewrite what the clauses ask for.
- No new profile, no change to another profile's decisions, and no change to how
  a decision's source or exception is recorded.
- No third-party tool, dependency or configuration file.
- No change to this repository's own `.roundfixrc.yml`.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.

---
status: approved
granted: 2026-09-19
action: let the Baseline Profile declare the incremental verification tier beside the complete gate, and render both where the decision is published
consuming: 0148-a-profile-that-declares-both-tiers
paths:
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/templates/guides/agent-instructions.md
  - internal/baseline/assets/templates/index.json
  - docs/agents/spec-routing.md
  - docs/agents/agent-instructions.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
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

The Profile that carries the decision, and the template that publishes it, are
Baseline assets: code-generator configuration, and therefore governed. The
repository guides that state the clauses are generated from Baseline modules
inside setup-context markers, so they are bounded here to be **regenerated**,
never hand-edited.

The template index is bounded beside the template it indexes: it declares which
tokens each template may render, and a token that is not declared there cannot
be rendered. The first Run of this Spec stopped on exactly that boundary — the
Agent refused to edit the template while the index stayed outside the grant —
and the maintainer widened the record on 2026-09-19.

The Profile reader, the planning output and their tests are ordinary source that
no authorization has bounded.

## Approved bounded mutation

- Give the Profile an incremental verification decision beside its complete
  gate, distinguishable from it and from a catalog default.
- Render both decisions where the guide template publishes the gate today, and
  declare the new token in the index that governs which tokens that template may
  render, raising that template's version as the index requires.
- Regenerate the repository guides from their Baseline sources, so the generated
  copies agree with the modules that own them.

The generated guides change only through regeneration. A hand edit inside a
setup-context marker is out of scope and would be overwritten by the next
Baseline update.

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

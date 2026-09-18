---
task: task_03
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
status: completed
type: docs
complexity: low
---

# Task 03: Document the key beside the Review Source

## Overview

The public guide gains the new key, its four values, its precedence and its
current reach, documented next to `review_source` so a reader can tell the two
acts apart.

## Requirements

1. MUST document `pre_pr_review.provider`, its four supported values and its
   precedence.
2. MUST state that an absent key inherits and never means `none`.
3. MUST state that declaring a policy invokes no provider, and that this release
   reports the policy without yet enforcing it at publication.
4. MUST place the description beside `review_source` and say which act each one
   governs.
5. MUST describe the new Doctor check in the Roundfix skill and regenerate its
   distributed mirror, which are this Spec's only bounded governed paths.
6. MUST NOT edit this repository's own configuration file, any Baseline asset or
   any guide delivered inside setup-context markers.

## Subtasks

- [ ] Write the key, values and precedence into the user guide.
- [ ] Contrast it with `review_source` in the same place.
- [ ] State the current reach plainly.
- [ ] Describe the check in the Roundfix skill and regenerate the mirror.

## Acceptance Criteria

- [ ] The guide names the key, its four values and its precedence.
- [ ] The guide states that silence inherits.
- [ ] The guide states that nothing is invoked and that enforcement is not yet
      part of publication.
- [ ] The Roundfix skill and its mirror describe the check, and no other
      governed path is changed.

## Context

- interface: `docs/user-guide/usage.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "pre_pr_review" docs/user-guide/usage.md` — expected: exit 0; the guide documents the key. Before this Task it does not.
- `grep -q "pre_pr_review.provider" docs/user-guide/usage.md` — expected: exit 0; the guide names the key by its full path. Before this Task it does not.
- `grep -q "inherits" docs/user-guide/usage.md && grep -q "pre_pr_review" docs/user-guide/usage.md` — expected: exit 0; the guide states that an absent key inherits. Before this Task the key is undocumented, so the command fails.

## References

`_prd.md` → Core Feature 1; Non-Goals; Project Constraints: Tooling authority;
`_techspec.md` → System Architecture: Public guidance; Risks & Considerations;
Build Order 3; `_authorization.md`.

## Result

### Implementation

- Added the `pre_pr_review.provider` policy description to the public usage
  guide, including its four supported values (`codex`, `claude`, `coderabbit`,
  and `none`), Project Config → User Config → built-in `codex` precedence,
  inheritance semantics, and the current read-only reach.
- Placed the policy beside the `review_source` explanation and stated that the
  former governs review before a Pull Request exists while the latter reads
  feedback from an existing Pull Request.
- Added the same policy row beside `review_source` in the configuration key
  reference. No repository configuration, Baseline asset, or setup-context
  guide was changed.

### Focused checks

- `git diff -- docs/user-guide/usage.md docs/user-guide/configuration.md` —
  inspected the documentation-only diff and confirmed the two guide changes
  contain the required policy language.
- `rg -n "pre_pr_review\.provider|absent key|invokes no provider|does not enforce|review_source" docs/user-guide/usage.md docs/user-guide/configuration.md` —
  matched the key, inheritance, invocation, enforcement, and act-contrast
  statements in both public guides.

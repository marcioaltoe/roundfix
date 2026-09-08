---
status: accepted
created_at: 2026-09-08T15:10:02Z
updated_at: 2026-09-08T15:10:02Z
deprecated_at: null
superseded_by: null
---

# Work branches express purpose and Run branches express ownership

On 2026-09-08 the maintainer replaced the personal-prefix aspect of ADR-0118:
new work branches use `<type>/<description>` with the repository's commit/PR
types, including `feat`, `fix`, and `refactor`, without names or initials.
Roundfix-owned Run and Task branches keep their existing `roundfix/run-`
namespace and derived names, which distinguish tool ownership from work intent.
Keep the compatible `branch.prefix` string key with `<type>/` as its pattern
default, revise this repository's saved value through Baseline, and preserve
existing branches and historical evidence rather than renaming them in bulk.

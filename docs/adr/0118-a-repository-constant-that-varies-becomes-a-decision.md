---
status: accepted
created_at: 2026-08-10T00:00:00Z
updated_at: 2026-08-10T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A repository constant that varies becomes a decision

ADR-0078 moved confirmed root rules out of `AGENTS.md` and into their semantic
owners. It settled where a rule lives once the Baseline owns it. It did not
settle what to do with a rule whose *value* differs per repository, and the
branch prefix is that case: the rule "agent-created branches carry a prefix" is
portable, while `ma/` is one maintainer's answer.

Writing the whole sentence into `docs/agents/specific-repository.md` made the
portable half unavailable to every other repository. Writing `ma/` into the
catalog would have imposed one maintainer's initials on all of them.

Roundfix therefore splits the two: the obligation becomes a Normative Clause in
the catalog, and the value becomes a Baseline decision with a suggested default.
`branch.prefix` is a string decision defaulting to `ma/`, selected by every
Profile, rendered into `guide.agent-instructions` beside `verification.gate` —
which was already the same shape, and is now the precedent rather than the
exception.

The test is whether a second repository could answer differently and still be
correct. A branch prefix passes it. A gate command passes it. The rule that a
pipe must not hide a gate's exit status does not — no repository benefits from
answering that differently — so it stays a clause with no decision behind it.

Making a decision required has a fleet cost, and it is paid on the next update
rather than silently: a repository whose Setup Manifest predates the decision
reports `action_required` with `fileChanges: []` and names the missing decision
with its suggested value. Measured on 2026-08-10 against a sibling repository,
the update refused to plan and offered `ma/` for adoption. That refusal is the
intended behavior — a decision nobody answered is not a default nobody chose —
and `--adopt-suggested` is the one-flag path for a maintainer who agrees with
the suggestion.

Rendering the prefix from a repository-local file that the catalog reads was
rejected: it puts the value outside the Decision Document, so the manifest no
longer records what the repository actually answered, and Readoption has nothing
to carry forward.

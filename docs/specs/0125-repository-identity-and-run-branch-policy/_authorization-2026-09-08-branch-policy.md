---
status: approved
granted: 2026-09-08
action: replace personal branch prefixes with purpose-based work branches and a Roundfix-owned Run namespace
consuming: 0125-repository-identity-and-run-branch-policy
paths:
  - internal/baseline/assets/decisions.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/guides/agent-instructions.md
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
---

# Approved branch-policy correction

On 2026-09-08 the maintainer explicitly removed the mandatory `ma/` prefix,
requested work-branch prefixes describing their purpose like commit/PR title
types, and permitted a separate Roundfix namespace for Runs. This decision
supersedes the personal-prefix requirement; it does not need another approval.

## Bounded implementation

Change only the branch decision's default/description, the existing Git-delivery
clause, the guide's branch wording, and their versions in the named catalog
sources. Keep `branch.prefix` as a compatible string decision and select the
pattern `<type>/` through the public Baseline plan/apply interface. The token
`<type>` is replaced by the work's purpose; it is never a literal branch name.
Render work branches as `<type>/<description>` using the repository's accepted
commit types, including `feat`, `fix`, and `refactor`. Personal names or initials
are not prefixes. Roundfix retains its existing `roundfix/run-` namespace and
derived Task branch names; no runtime branch rename or code change is required.

`make baseline-digests` is sanctioned deterministic fallout of the source edit.
The public Baseline command owns the two generated repository files. Review
its exact preimages, postimages, retained decisions, and digest before apply;
do not hand-edit the Setup Manifest or unrelated managed regions.

The same confirmed policy applies to these exact global instruction sources
in the local vault repository: `home/.codex/AGENTS.md` and
`home/.claude/CLAUDE.md`. Their user-facing paths are `~/.codex/AGENTS.md` and
`~/.claude/CLAUDE.md`; OpenCode's `~/.config/opencode/AGENTS.md` resolves to
the Claude file, so its link needs no change. Limit those edits to the branch
instruction and preserve the other instructions.

Update the current Spec proposals, design, public decision example, and
planning record to remove the now-resolved bootstrap blocker. Original
Findings and historical grants retain their observations with a dated addendum.
This approval does not grant the broader repository-identity migration,
reconciliation changes, paid probes, review-provider choice, or CodeRabbit
replacement proposed elsewhere in Spec 0125 and the portfolio.

## Research and influence

Secondbrain's accepted decision 0118 and the 2026-08-09 five-rules authorization
explain why `branch.prefix` was previously a configurable personal value. The
new maintainer decision changes that policy while retaining the existing key
for compatibility. Local inspection confirms that deleting the key would
reject saved decisions; a public plan with `<type>/` changes only the guide and
manifest. Neither observation requires changing Run branch creation.

Exa located and read [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/)
and [Git ref-name validation](https://git-scm.com/docs/git-check-ref-format).
The former supplies purpose vocabulary, including `refactor`; it specifies
commit messages, not branch naming. Reusing those types for branch names is
the maintainer's policy. Git permits slash-separated names and supplies the
public syntax check used for the examples. No external search included local
source code, credentials, or private records.

## Delivery and limits

Land this record before its consuming catalog/instruction commits. Existing
commit/push authority continues to apply. New work uses a purpose-prefixed
branch; existing remote branches and historical records are not deleted or
mass-renamed. Validate generated guidance, unchanged non-branch decisions,
catalog regeneration, and the repository's required checks before claiming
the correction complete.

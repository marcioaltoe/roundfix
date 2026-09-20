# Pull Request equivalent evidence

Environment cause: the supervised QA prompt states that no Pull Request is open
for `feat/0151-supersession-the-archive-can-see`, and Pull Request journeys are
environment-blocked. It also states that the checked-out Run Worktree branch is
never pushed and has no Pull Request of its own. The gate did not try to resolve
or mutate a Pull Request.

Equivalent controls:

- Approval: no Pull Request acceptance was obtained.
- Checks and status: the Daemon ran repository Verification `make verify`; it
  exited `0`, and diagnostics were removed on success.
- Unresolved review threads: no Pull Request exists, so no Pull Request review
  thread stands for this candidate.
- Merge-Ready acceptance: none exists yet; pre-PR review and Pull Request
  acceptance happen after this authored terminal gate.
- Review-artifact ancestry: no review artifact exists yet, so no artifact claims
  a head. The audited candidate head is
  `60224664f49268383026f9b4a8bee15b89072b10`.

Every unavailable control is stated rather than inferred as passing. The row
remains environment-blocked, with equivalent supervised evidence recorded as
ADR-0080 requires.

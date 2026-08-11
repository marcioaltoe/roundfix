# Fluxus update command evidence

Repository copy:
`/private/tmp/roundfix-qa-0047.nh68hi/fluxus-update`

The source copy started clean at
`1aeed7e8370c3d14137c42b0c789dcbe3bd1ba3b`.

The automation path first required
`preservation.mode=preservation`, then exited `3` asking for a complete
classification Decision Document without emitting a partial Plan.

The public interactive workflow was then exercised from the same clean copy:

1. Selected Preservation.
2. Reused `standard-typescript-monorepo`.
3. Kept every existing repository decision.
4. Reached ready Profile alignment with advisory-only divergences.
5. Received one consolidated proposal:
   `AGENTS.md bytes 0-2304: normative-clause -> repository-rules`.
6. Accepted the complete product proposal.

The workflow exited `3` before a Change Plan:

```text
complete or correct every root-rule disposition and rerun Baseline planning
baseline.preservation.repository-rules.invalid
Repository-Specific Normative Rules proposed bytes are stale, empty, or managed
```

The proposed bytes are the existing setup-managed root. The product therefore
rejects its own default proposal and gives the maintainer no completed
semantic redistribution or Plan Digest. No repository byte changed.

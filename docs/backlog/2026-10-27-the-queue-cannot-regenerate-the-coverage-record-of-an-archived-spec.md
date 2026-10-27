---
type: fix
status: open
created: 2026-10-27
spec: null
reason: null
---

# The queue cannot regenerate the coverage record of an archived Spec

## Symptom

When a delivered item's Pull Request conflicts with `main` only in derived
records, the Delivery Queue resolves the conflict by re-running their
regeneration commands. For `docs/references/coverage-record.json` it runs
`go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`,
and the suite guard refuses the write (`suiteguard: repository boundary
violated: modified: docs/references/coverage-record.json`), because the item
is already archived and no active Spec's `_authorization.md` declares that
command. The item parks `delivery-error: resolve pull request conflict:
regenerate derived paths`.

## Where

The queue's derived-path regeneration (`resolve pull request conflict`), and
`internal/suiteguard` `removeSanctionedRegeneration`, which reads
sanctioned commands only from active Specs
(`suiteguardcontract.ReadSanctionedRegenerations`).

## Expected

A regeneration the queue runs after archive is sanctioned by the archived
item's own authorization record (or the Archive Record's `regeneration`
list), so a conflict confined to derived records resolves without the
operator.

## Evidence

Spec `0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green`
on 2026-10-27 ~18:10 BRT: Pull Request #709 conflicted with 0331 (#708) only
in the Roundfix skill version lines, the owned-skill version record, the
coverage record and the Behavior Surface map; the queue's regeneration log
`derived-regeneration-6.log` shows the refusal. The operator merged by hand,
ran the sanctioned command with the Spec's `_authorization.md` restored
untracked, and merged #709 after its checks passed.

---
status: accepted
created_at: 2026-10-13T00:00:00Z
updated_at: 2026-10-26T21:00:00Z
deprecated_at: null
superseded_by: null
---

# A row the authorization keeps offline qualifies, and a byte invariant compares with the Delivery Base

Two acceptance rules disagreed with the authority and the queue they run
under. Spec 0258's QA found no defect, but its authorization said that no QA
row opens a network connection. The gate therefore recorded four
outside-evidence rows as
`blocked (environment: network prohibited by Spec authorization)`. The shared
QA policy of ADR-0240 exempts only the network-denied form, so the Delivery
Queue parked the item `qa-environment-partial` and the operator archived it
with the QA Archive Override. Spec 0264's QA requirement held
`docs/references/behavior-surfaces.json` byte-identical to its bytes at the
commit where the plan was measured. Four Specs merged first and changed the
file, so the QA failed a row although the Spec's own diff left the file
unchanged. Spec 0260 had failed its QA for the same reason a day earlier.

On 2026-10-13 the maintainer decided the first rule: such a row counts toward
a qualifying partial without the override, in one canonical wording the QA
skill writes, and settlement and archive share the policy.

**The Network-Prohibited Row.** A Results row whose status is exactly
`blocked (environment: network prohibited by Spec authorization)` and whose
provenance names the outside-evidence row, as ADR-0240 reads it, is exempt
like the Network-Denied Row. It never decides a qualifying partial, needs no
Unreachable Acceptance declaration, and never parks `qa-environment-partial`.
The rule lives in `internal/spec` beside the other two exempt kinds, so
`roundfix settle`, Daemon settlement, `roundfix qa-report accept`,
`roundfix archive` and the Delivery Queue's park decision read it alike. The
QA skill writes the status only when the Spec's own `_authorization.md`
forbids network access, names that text in the report, and still checks
offline any source it can, such as a passage the PRD quotes.

**The Pinned Byte Invariant.** A QA requirement that holds a file
byte-identical compares it with the Delivery Base, the merge base of the
audited head with the default branch. `write-tasks` states the rule, and
`roundfix spec check` reports `SC-QA-BYTES-PINNED`, a gap, for each
requirement of a pending QA Task that names a byte-identity phrase and a
commit-shaped hexadecimal token together. Like the gaps of ADR-0265, it holds
every active Spec with no horizon, and it reads structure only, as ADR-0093
requires. A scan of the 235 archived Specs at their source revisions found
five such QA requirements in Specs 0193, 0247, 0248 and 0260, all pinned to a
commit and none a false match. A non-QA requirement matched once, a sentence
kept byte-identical beside a binary's build commit. That match is why the
check reads QA Tasks only.

## Considered Options

- Record the row in the network-denied form: rejected. No sandbox denied
  anything, and the report would no longer say why the source was not read.
- Read `_authorization.md` to confirm the prohibition: rejected. Authorization
  records state it in prose and have no network field, and adding one would
  change every record. The status names the record, so a reader can check it
  afterwards, which is the trust ADR-0240 gives the named host.
- Accept a host suffix or any status that mentions the authorization:
  rejected. One exact status keeps the exemption from widening by paraphrase.
- Change the `### QA settlement` table: rejected. A contract test holds it
  byte-identical in three skills, so the extension sits under its own heading
  after the table in each.
- Compare a pinned commit with the plan's merge commit: rejected. That commit
  does not exist when the plan is authored, and the comparison would need Git
  history at check time.
- Exempt a requirement that also says "delivery base": rejected. A commit
  labelled as the delivery base is still a fixed commit.

## Consequences

- ADR-0240 is extended with a third exempt row kind, and ADR-0104 is
  superseded in part once more: a Network-Prohibited Row no longer holds Pull
  Request preparation.
- The Baseline clauses `clause.spec.project-constraints-06-outside-evidence`
  and `clause.autonomous.loop-01-qa-once` each gain one sentence. No clause is
  removed, and the rule against authoring an outside-evidence row on a lookup
  the authorization forbids still stands.
- The exemption rests on the text the gate writes. A gate that labels a row
  it could have reached hides that row from the override, and the archived
  report and the authorization make the claim checkable.
- An active Spec whose QA requirement pins bytes to a commit is refused under
  `--strict` until its author compares with the Delivery Base.

**Amended by [ADR-0335](0335-a-stop-that-is-not-the-work-s-fault-does-not-cost-the-run.md) (2026-10-26).** The Network-Prohibited Row's status is `blocked (environment: <subject> prohibited by Spec authorization)`, where the subject names what the Spec's own authorization forbids, is 1 to 80 characters and holds no parenthesis; `network` is the subject of the status above. The rejection of any status that mentions the authorization still holds outside that form, and a row counts toward one exempt kind at most.

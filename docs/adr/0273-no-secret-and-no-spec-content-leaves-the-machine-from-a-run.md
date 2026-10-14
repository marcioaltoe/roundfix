---
status: accepted
created_at: 2026-10-13T00:00:00Z
updated_at: 2026-10-14T23:44:53Z
deprecated_at: null
superseded_by: null
---

# No secret and no Spec content leaves the machine from a Run

On 2026-10-12 and 2026-10-13 the QA Runs of one Spec persisted the output of
`make -p`, which prints every inherited environment variable with its value,
as QA evidence; GitHub push protection refused the push of three Runs because
the evidence held a third-party API key. The same QA ran
`roundfix archive --plan`, whose Archive Advice sent the Spec's files to the
judge provider, because every Agent Session inherits the judge's keys. Nothing
between the QA Report commit and the push looked for a secret, and nothing in a
Run kept the judge's credentials out of the commands an Agent or a
Verification runs. On 2026-10-13 the maintainer approved cycle AG, whose item A
is this decision, under the operating rules of cycles AE and AF.

**One reviewed set of Secret Patterns.** `internal/secretscan` holds every
Secret Pattern Roundfix refuses, each with a kind name, in one table. A Secret
Match reports the repository-relative path, the line and the kind, and the
commit when it was read from history; it never carries the matched value or
any part of it. Tests build every sample shape at runtime, so no file in the
repository holds a value that looks real.

**The QA Report commit refuses a secret.** After the QA Agent hands back and
before the QA Task settles, the Daemon runs the Secret Scan over exactly the
files the QA Report commit would stage. A Secret Match settles the QA Task
`failed`, names each match in the reason and in one Run Event, and creates no
QA Report commit; the report and its evidence stay uncommitted in the Run
Worktree. This refines ADR-0091, which consolidated ADR-0015, whose QA Report commit is made "either way":
a report that carries a Secret Match is never committed.

**Every push Roundfix makes is scanned first.** The Delivery Queue's push and
the Daemon's Final Push read every commit the push would send, those not
reachable from the remote's tracking refs, and scan the lines each commit
adds. A Secret Match refuses the push before `git push` runs; the Delivery
Queue parks the item as `secret-found`. Each commit is scanned, not the net
difference, because GitHub's push protection requires removing a secret "from
all of the commits it appears in".

**The judge's keys never enter a Run.** The Judge Key Variables are the
variables the judge reads its credentials from, as its questions list them.
No Agent Session and no Verification command the Daemon runs receives one.
A Light Tier session, whose OpenRouter key may be the shared Roundfix key the
judge also reads, receives that key's value under the implementation Stage Key
name, so the judge finds no key there either. No Spec authorization reopens
the keys inside a Run: the judge is advisory and fails open, so nothing a Run
must prove depends on it, and the operator keeps running it outside Runs.

## Considered Options

- Skipping Archive Advice when `archive --plan` runs inside a Run was
  rejected. Detecting a Run from Roundfix Home or the worktree location fails
  when a QA row sets `HOME` to a temporary directory, as QA rows do, and it
  covers one caller while `spec judge`, `history sanitize` and any other
  program would still find the keys.
- A Spec authorization that lets a Run keep the keys was rejected. It needs a
  new value in the closed `operations` vocabulary, a Baseline contract change,
  for a judge whose answers never gate.
- Withholding the keys from QA sessions only was rejected: implementation
  Tasks and review Batches read the same Spec content.
- Scanning only `qa/evidence/**` was rejected: the report, the Task file and
  any implementation commit can carry a value too.
- Redacting a matched value instead of refusing was rejected: it rewrites what
  QA observed, and the command that printed the value would keep doing so.
- Running an external scanner such as gitleaks was rejected: it is a new
  dependency nobody authorized, and it is not installed on every host.

## Consequences

- A QA row that needs to show a variable records its name and whether it is
  set, never its value.
- An operator whose Delivery Queue item parks `secret-found` removes the value
  from every commit the blocker names on the item branch before the retry.
- A pattern that is too broad refuses honest evidence; the table is the one
  place to narrow it, and its tests name each kind's shape.
- Git hooks and the QA format command keep the Run's environment; they are
  repository configuration, not commands a Spec or an Agent chooses.

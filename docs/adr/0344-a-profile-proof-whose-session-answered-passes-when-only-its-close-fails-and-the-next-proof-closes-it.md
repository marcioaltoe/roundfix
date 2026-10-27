---
status: accepted
created_at: 2026-10-27T15:30:00Z
updated_at: 2026-10-27T15:30:00Z
deprecated_at: null
superseded_by: null
---

# A profile proof whose session answered passes when only its close fails, and the next proof closes it

On 2026-10-26, under a load average of 54, Vortex's pre-PR review and two
items of this repository's Delivery Queue parked with `profile proof setup
timed out on 3 attempts of 2m0s each`. Every setup step had answered; the
failing step was `close disposable Agent Session …: context deadline
exceeded`. The close ran under a fixed 5-second deadline, its error failed
the proof, and the CLI read the wrapped deadline as a setup timeout, so it
repeated the whole proof twice and reported a setup that never timed out.
The verdict of an Exact Agent Selection Proof now rests on its setup alone,
the close gets the same deadline as the setup, and a session Roundfix could
not close is recorded so that a later proof closes it.

## Decision

**The verdict is the setup's.** When every setup step answered, the proof
passes even if the close of its disposable Agent Session fails or times out;
Roundfix prints one warning that names the session, the deadline and the
cause. When a setup step failed, the setup's error is the proof's error, with
its classification, and a close failure no longer joins it. A session that
never answered still fails exactly as before.

**The close follows the setup deadline.** A proof attempt closes its
disposable Session under the attempt's Proof Setup Deadline
(`defaults.profile_proof_setup_timeout`, built-in `30s`), not a fixed 5
seconds; a Fallback Chain candidate probe uses its built-in setup deadline.
Under load every acpx call slows alike, and the close is one more acpx call.
Sealed prompts and Run sessions keep their own cleanup deadlines.

**An unclosed session is recorded and closed later.** A failed close writes
an Unclosed Proof Session record, `~/.roundfix/unclosed-proof-sessions/<session>.json`,
holding the session name, its working directory and the runtime fields that
rebuild the `acpx sessions close` command, never an environment value. The
first proof each Agent runner makes closes the recorded sessions first,
oldest first, under the same deadline. It deletes a record once acpx closes
the session or reports it absent, and stops at the first close that fails,
keeping that record and the rest for a later runner. Reaping never changes
a proof's verdict. acpx documents why the record is needed: `sessions close`
"marks the record `closed: true`" and tears down adapter processes, while
`sessions prune` deletes only closed records, so a session nobody closes
stays open on disk.

**A deadline names its step, and only a setup deadline is a setup timeout.**
A setup step that reaches its deadline is reported as `profile proof step
"<step>": context deadline exceeded`, so the third-attempt refusal names the
step that was slow. A deadline found only inside a close error is never
counted as a setup timeout and is never retried as one.

**Baseline semantic analysis keeps its stop.** ADR-0069's analyzer treats a
proof that left its session open as cleanup unproven, as it did when the close
error failed the proof, and falls back to manual classification.

## Considered Options

- Keep failing the proof and only lengthen the close deadline. Rejected: a
  correct tuple would still park a delivery whenever acpx is slow to close,
  and the proof had already shown what it exists to show (ADR-0140).
- Close the session in the background after the proof returns. Rejected: the
  CLI process exits before a background close finishes, and a goroutine
  without an owner breaks the Go guide.
- Sweep every open `roundfix-preflight-*` session that `acpx sessions list`
  shows. Rejected: it races a concurrent proof in the same directory, needs an
  age guess, and misses sessions whose item worktree is already gone; the
  record names exactly what Roundfix failed to close.
- Name the host load average in the refusal. Rejected for now: it needs a
  platform file per operating system to read a kernel value, while the step
  name and the deadline already tell the operator what was slow, and the next
  action already says to retry when load drops.

## Consequences

- A setup that times out can now wait up to its deadline again for the close,
  so with the default deadline a tuple whose setup always times out fails
  after about three and a quarter minutes instead of about two (ADR-0298).
- An Unclosed Proof Session record lives in the Roundfix home until a later
  proof closes it; a record that cannot be parsed is reported and left for
  the operator.
- `profiles configure` and Setup can now write a profile whose proof passed
  with an unclosed session; the warning says so.

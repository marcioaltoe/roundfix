---
type: fix
status: promoted
created: 2026-10-01
spec: 0211-a-delivery-queue-that-finishes-without-intervention
reason: null
---

# An archived retry cannot find the Run's Implement start head

## Opportunity

Spec 0201 lets `deliver retry` resume an item the operator archived after an environment-only QA partial. On its first real use, Spec 0204 parked `qa-environment-partial`. The operator followed the printed answer: carry-forward, operator evidence, then `archive --qa-override`. The retry then refused with `archived retry proof: … Run "run_20261001T094428Z_f4d3047c47b08ac0" has no Implement start head for repository "/Users/marcio/dev/roundfix"`, and with `archived item head … differs from candidate head ""`. The Run was started by the queue owner built from the 0201 merge.

## Value

The class 0201 added for this exact case still ends in a manual Pull Request. The supervisor also has to run the pre-PR review by hand, because the park happens before review.

## Shape

Reproduce with a queue-started Run that parks `qa-environment-partial`, then find why the proof finds no Implement start head: the journal event it reads, the repository key, or the Run it chooses. Add an end-to-end test that drives the whole printed answer through `deliver retry` to `reviewing`.

## Extended on 2026-10-01 — the same queue's other interventions

The same day the Delivery Queue needed three more operator interventions that
belong to this entry, because each one stopped an item the queue should have
finished, and two of them are `deliver retry` answers that did not resume the
item.

- **A retry merged while a re-run was pending.** Spec 0203 parked
  `checks-failed` on a check that failed outside its packages. The operator ran
  `gh run rerun --failed` and then `deliver retry`. The item moved from
  checking to merging while the re-run was still pending. GitHub refused the
  merge ("the base branch policy prohibits the merge") and the item parked
  `delivery-error`. The queue must read the current check state of the
  re-run: pending is not pass.
- **A retry refused a park whose Pull Request had turned green.** `deliver
  retry` on that `delivery-error` park, with the Pull Request green and
  mergeable, refused with `Preflight failed` and the usage text. The operator
  merged Pull Request #317 by hand. A retry must resume the merge stage.
- **A dead preload stopped every adapter.** The queue owner inherited
  `NODE_OPTIONS=--require=<temporary preload>`, and the system's temporary-file
  cleanup had removed that file. Every Node process the Run started, beginning
  with `acpx --version`, died with `Cannot find module`, and the Implement
  preflight refused with a profile-proof failure that named the adapter, not
  the preload (Spec 0210). The owner or its preflight should drop a
  `--require` or `--import` of a missing file from the agent environment with
  a named notice, or refuse with a reason that names `NODE_OPTIONS`; dropping is
  preferred, and the user's environment is never otherwise edited.

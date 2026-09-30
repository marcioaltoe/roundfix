---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A configuration is compared with the recommendation only where a person asked

ADR-0180 gives Roundfix one dated Recommended Profile per Agent Work Category,
and built-in profiles follow it with each release. A configured profile does
not. On 2026-09-30 this repository's own configuration still led with the
model its provider had replaced the day before, and nothing said so until the
maintainer read a provider page. A configured profile is also a decision: a
repository may keep an older model on purpose, and Roundfix must not argue
with that on every Run.

Roundfix compares the effective Agent Selection Profile of each configured
category with the Recommended Profile, under these rules:

- **Where.** The comparison runs in `roundfix profiles check`, in one Doctor
  line, in `roundfix profiles show`, and as a notice after `roundfix upgrade`.
  It never runs when a Run starts, in Preflight, in the Daemon or in any
  command that does Agent work.
- **What it costs.** It is read-only and offline. It reads the configuration
  and the snapshot in the binary, opens no Agent Session and no network
  connection, and writes nothing.
- **What it decides.** Nothing. A difference never fails a command, and a
  command that carries the notice keeps its standard output and its exit code.
- **After an upgrade.** When `upgrade` installs a release, the notice comes
  from the installed executable, because the snapshot that matters is the new
  one. When nothing was installed, the running executable answers.
- **A deliberate difference.** A profile may carry a Profile Deviation: the
  snapshot date it was declared against and a reason. While the binary ships
  that snapshot, the category is reported as pinned and no notice is raised. A
  new snapshot ends the pin, and the difference is reported again.
- **Adoption.** `roundfix profiles check --apply` writes the Recommended
  Profile for the differing categories through the same proved, confirmed
  write as `roundfix profiles configure`. A pinned category is skipped.

## Consequences

- A model retirement reaches a configured repository as a notice at the
  moment its maintainer upgrades, and as a Doctor line afterwards.
- A deliberate choice is written next to the profile it explains, and it is
  decided again once per snapshot, not once per Run.
- A configuration file with a Profile Deviation is rejected by an older
  Roundfix, because profile keys are strict.
- The first release with this rule cannot speak after an upgrade performed by
  an older executable. The notice starts with the upgrade after that one.

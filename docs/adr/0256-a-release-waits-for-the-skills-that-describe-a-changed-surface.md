---
status: accepted
created_at: 2026-10-08T00:00:00Z
updated_at: 2026-10-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A release waits for the skills that describe a changed surface

The Roundfix-owned skills must describe the behavior Roundfix ships. The
mirror and the version record are checked, but nothing checks that a change to
a command's help, a configuration key, an exit code or a user guide reached the
skills that describe it. Specs updated skills only when their author
remembered: the operator's queue log records three QA failures on that class
in two cycles (entries 193, 231 and 232), two of them in one Spec, where the
settle reference contradicted a new reopen rule. The release runbook's
skills-and-guides step ends in a reading pass, and the Release Plan Command
reports its skills and baseline checks without letting them change its
result. On 2026-10-08 the maintainer chose "Check no release + regra de autoria
(Recommended)": the release plan refuses a release whose skills fell behind,
and Spec authoring requires a skills Task for a change a skill describes.

**Behavior Surfaces and their map.** A Behavior Surface is one user-visible
unit of behavior: the help of one command path, the set of configuration keys,
the table of exit codes, or one file of the user guide. The Skill Coverage Map
is a checked-in, authored list of every surface with the owned skill files
that describe it, or a reason why no skill does, the source paths whose change
can change it, and an optional Coverage Review. The Behavior Surface Record is
a separate, generated record of one fingerprint per surface. A repository
contract test computes every fingerprint from the built command help, the
configuration source, the exit-code constants and the guide bytes, and fails
when the record is stale or when the map and the computed surfaces disagree.
Only its record flag writes the record, and the record is a declared derived
path, so two Pull Requests that both change a surface merge by regeneration
(ADR-0192).

**The release check.** A range plan reads the map and the record at its base
and at its target from Git, and compares the fingerprints. A surface whose
fingerprint changed, appeared or disappeared is a Lagging Surface unless a
covering skill file changed in the same range, the surface is recorded as
uncovered, or its Coverage Review changed in the range. A Coverage Review is a
dated note in the map entry that says why the change needs no skill text; it
answers only the range in which it was written, which matches published drift
tools that keep a review until the reviewed code changes again. A lagging
surface or an unreadable map makes the new `skill-coverage` check blocking.
For a plan that proposes a version, the command then exits 3 and its next
action names the check. The decision state and the proposed version do not
change, and the skills and baseline checks stay advisory. A repository with
no map is reported as `not_declared`, and a base without one as `introduced`;
neither blocks, so adopters and the first release after this decision are not
refused for history that never had a map.

**The authoring rule.** `roundfix spec check` reads the map's source paths. A
Spec whose non-QA Tasks declare a source of a covered surface needs a non-QA
Task that declares one of that surface's skill files, or a `- unchanged:`
entry under the PRD's `## Skills` section together with a Task that declares
the map, where the Coverage Review is recorded. The check applies to every
route, because a Spec's artifacts do not record whether it is a feature, a
refactor or a fix, and a fix that changes a surface leaves a skill as stale as
a feature does. It starts at the commit that adds the map, so earlier Specs
keep their meaning.

We rejected four alternatives. A Git-only comparison of source files cannot
tell one command's help from another's, because one file holds every command's
usage text. Fingerprints inside the map would mix authored and generated
fields, and declaring the whole map derived would let a conflict take the
default branch's whole file and drop an item's coverage edits, the reason
ADR-0233 kept `SKILL.md` out of the derived set. A Baseline clause would tell
every adopting repository to keep a map it does not have; a Baseline guide
says only what holds for the repository that reads it (ADR-0222). A model
judging the diff was rejected for a gate, as ADR-0189 asked of the release
step: the check needs no model and no network.

## Consequences

Every change to a command's help, the configuration keys, the exit codes or a
user guide now re-records the Behavior Surface Record in the same change, or
`make verify-docs` fails. A release whose range changed a surface without its
skill needs a skill edit or a Coverage Review before its Pull Request. The
authoring rule reads only declared paths, so a change to the shared usage file
alone is caught at release rather than at authoring. The first release after
this decision reports `introduced` and is not compared.

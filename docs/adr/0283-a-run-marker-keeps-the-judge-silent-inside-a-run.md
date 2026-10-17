---
status: accepted
created_at: 2026-10-17T20:40:00Z
updated_at: 2026-10-17T20:40:00Z
deprecated_at: null
superseded_by: null
---

# A Run Marker keeps the judge silent inside a Run

On 2026-10-16 and 2026-10-17 the QA Agents of Specs 0273 and 0275 ran
`roundfix archive <slug> --plan`, and its Archive Advice sent Spec files to the
judge with `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, although ADR-0273 removes every
Judge Key Variable from each Agent Session and Verification command. Roundfix's
own boundary held: the acpx command environment and the Verification
environment are both built without the variables, and tests prove it. The key
came back inside the runtime. Codex's shell snapshot, a stable feature enabled
by default, sources the user's `~/.zshrc`, records every exported variable, and
sources that snapshot before each shell command, so a key the user's shell
startup exports reaches every command again. Any runtime that reads the user's
shell startup does the same, and nothing a parent process removes from its
child's environment can stop it.

## Decision

**Every process a Run starts for an Agent or a Verification carries the Run
Marker.** The Daemon gives each acpx invocation that prepares or prompts an
Agent Session of a Run, and each Verification command it runs for a Run (a
Task's Verification, the probe on the unchanged tree and the repository
Verification at settlement), the variable `ROUNDFIX_RUN_ID` set to the Run's
ID, the variable the Run's notification command already receives. A user's
shell startup does not set a Roundfix variable, so a runtime that re-reads it
keeps the marker.

**The judge sends nothing while the Run Marker is set.** `roundfix spec
judge`, the Archive Advice of `roundfix archive --plan` and the advice of
`roundfix history sanitize` read `ROUNDFIX_RUN_ID` from the environment they
were given. When it is set and not empty, the judge loads no key, reads no
Judge Log and makes no request, and reports the skip reason `inside Run
<run-id>, named by ROUNDFIX_RUN_ID`, before it looks at keys or the Jev
Ceiling. The check is in the judge's request path for every caller, not in one
command.

**The removal stays, and covers the generic TypeSafe key.** The Judge Key
Variables gain `TYPESAFE_API_KEY`, the variable the TypeSafe provider documents
for its own clients, so an Agent following the provider's documentation finds
no key either. The Daemon and the Delivery Queue's Queue Owner keep their own
environment: the owner passes it to the item binary, and the Light Tier needs
the shared Roundfix OpenRouter key there to hand it to a light session under
the implementation Stage Key name.

This amends ADR-0273, which rejected "Skipping Archive Advice when `archive
--plan` runs inside a Run" because the Run was detected from Roundfix Home or
the worktree location, which a QA row's temporary `HOME` defeats, and because
only one caller was covered. The Run Marker is a variable the Daemon sets, it
survives a changed `HOME`, and every judge caller honours it. ADR-0273's
removal of the keys and its rule that repository-configured commands (Git
hooks, the QA format command and the Worktree Bootstrap) keep the Run's
environment are unchanged.

## Considered Options

- Turning Codex's shell snapshot off for Run sessions was rejected. It is one
  runtime's setting, Claude Code also restores a snapshot of the user's shell
  configuration, and a login shell still reads the user's profile.
- Pointing `ZDOTDIR` or `HOME` at an empty directory for Agent sessions was
  rejected: the Agent's toolchain paths and runtime credentials live in that
  same shell startup and home.
- Detecting a Run from the working directory was rejected for the reason
  ADR-0273 gives: QA rows work in temporary repositories.
- Refusing only when keys are present was rejected: the reason must name the
  Run whether or not the shell startup exports a key.

## Consequences

- On a host whose shell startup exports a judge key, an Agent's commands still
  see the key. Roundfix sends nothing with it; another program could, which is
  why operators should keep judge keys out of interactive shell startup files.
- A test that injects its own environment, as the command tests do, is not
  affected by the marker its Verification inherits.
- An operator who runs the judge from inside a Run's shell, by hand, gets the
  skip; the judge is advisory and is run outside Runs.

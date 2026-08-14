---
spec: 0013-codex-runtime-hygiene
status: archived
created: 2026-07-06
surfaces: [cli, infra, docs]
archived: "2026-07-06"
source_slug: 0013-codex-runtime-hygiene
---


# Codex Runtime Hygiene

macOS Tahoe blocked Roundfix Runs mid-flight with "Malware Blocked: 'codex' was
not opened because it contains malware." The codex binary is OpenAI-signed but
not Apple-notarized, so XProtect YARA-scans the ~193 MB binary on every exec and
intermittently blocks it — worst when the binary carries the Homebrew-Cask
`com.apple.quarantine` flag and an agent loop spawns it many times a minute. The
environment was repaired by hand (curl install to `~/.local/bin` without
quarantine, `CODEX_PATH` pinned in the shell), but no product should depend on a
hand-repaired machine. This Spec teaches Roundfix to detect an unsafe codex on
PATH and to spawn a verified-clean one, so an agent loop never trips XProtect.

## Goals

- Roundfix detects a quarantined or un-notarized codex on PATH before it can
  break a Run, and names the exact fix. See ADR-0032.
- The Doctor Command reports machine health for Roundfix Runs and returns a
  clear pass/fail the developer and an agent can both read.
- When Roundfix spawns codex through `codex-acp`, it prefers a verified-clean
  codex, so repeated execs in an agent loop do not re-trigger XProtect.

## User Stories

1. As a developer on macOS, I want Roundfix to tell me when the codex on my
   PATH is quarantined or un-notarized, so that I learn the risk before a Run
   dies with a malware block instead of after.
2. As a developer who hit the block, I want the diagnosis to print the exact
   curl-reinstall fix, so that I can repair the machine in one step without
   researching the OpenAI issue myself.
3. As a developer preparing a machine, I want a Doctor Command that reports
   every Run prerequisite's health in one place, so that I can check readiness
   without starting a Run.
4. As a developer running an agent loop, I want Roundfix to spawn a
   verified-clean codex, so that frequent execs stop intermittently tripping
   XProtect mid-Run.
5. As a developer on Linux or Windows, I want the codex-hygiene check to report
   not-applicable rather than fail, so that Doctor is meaningful on every
   platform.

## Core Features

1. **Doctor Command.** `roundfix doctor` runs the machine-health checks
   Roundfix Runs depend on — Node, pinned acpx, configured Agent probe, and the
   codex-hygiene check below — and reports each with a status and a next
   action, exiting non-zero when a check that would break a Run fails. It
   mutates nothing (diagnosis only), unlike the Setup Command which prepares the
   machine.
2. **Codex hygiene check.** On macOS, the check inspects the codex resolved on
   PATH (and any configured codex path) for the `com.apple.quarantine`
   attribute and for notarization/signing acceptance, reporting a failure with
   the curl-to-`~/.local/bin` reinstall fix when the binary is quarantined or
   not accepted. On other platforms it reports not-applicable. See ADR-0032.
3. **Verified-clean codex on spawn.** When Roundfix launches codex through
   `codex-acp`, it resolves a verified-clean codex (respecting a configured
   codex path) and spawns that, so an agent loop's repeated execs do not
   re-trigger XProtect. Resolution never silently spawns a known-quarantined
   binary without surfacing the risk. See ADR-0032.

## User Experience

`roundfix doctor` prints one line per check — node, acpx, agent, codex — each
with a pass/skip/fail status and, on failure, the next action (for codex, the
curl-reinstall command). It changes nothing on disk. On macOS with a
quarantined codex, the codex line fails and Doctor exits non-zero; on Linux and
Windows the codex line reports not-applicable and does not fail the command.
During a Run, codex spawning is unchanged in output but sourced from a
verified-clean binary.

## Non-Goals / Out of Scope

- Downloading, reinstalling, or de-quarantining codex automatically — Doctor
  diagnoses and instructs; the developer runs the fix. (The Setup Command may
  adopt the remediation later; not here.)
- Notarizing or signing Roundfix's own binaries — that belongs to the
  distribution work, not runtime hygiene.
- Hygiene checks for the Claude or OpenCode runtimes — this Spec addresses the
  observed codex/XProtect failure only.
- A background watcher or periodic health poll — Doctor runs on demand.
- Windows Defender or Linux equivalents — the observed failure is macOS
  XProtect; other platforms report not-applicable.

## Success Metrics

- On a machine whose PATH codex carries `com.apple.quarantine`, `roundfix
  doctor` fails the codex check, prints the curl-reinstall fix, and exits
  non-zero; after the fix it passes and exits zero.
- On Linux and Windows, `roundfix doctor` reports the codex check
  not-applicable and does not fail on that account.
- A Run that spawns codex sources a verified-clean binary; the exact 2026-07-06
  failure (quarantined codex blocked mid-loop) does not recur on a machine
  Doctor reports clean.

## Decisions

- Codex hygiene is exposed through a new Doctor Command (diagnosis only),
  distinct from the Setup Command (machine preparation). Adds an "Archive
  Command"-style glossary term: Doctor Command.
- Roundfix prefers a verified-clean codex when spawning `codex-acp`, respecting
  a configured codex path. See ADR-0032.
- Quarantine/notarization detection is macOS-only; other platforms report
  not-applicable. See ADR-0032.

## Open Questions

None.

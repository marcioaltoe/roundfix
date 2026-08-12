---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# `release plan` requires a `v`-prefixed tag and cannot see a repository without one

## Symptom

In `fluxus`, all eight release tags carry no prefix — `0.7.2`, `0.7.1`, `0.7.0`,
`0.6.1`, `0.6.0`, `0.5.7`, `0.5.6`, `0.5.5`. `roundfix release plan` finds none:

```text
roundfix: release plan failed: resolve release range: find latest stable release tag:
no stable release tag; create or pass a stable vMAJOR.MINOR.PATCH base tag
```

Measured on 2026-08-11 with roundfix 0.4.0 (`0c7d2d3c`).

The format is not a style preference in the target repository: its
`cd-github-release.yml` **uses the absence of the prefix to decide whether the
release is stable**.

```sh
# A pure MAJOR.MINOR.PATCH tag is a stable release; anything else
# (0.6.0-rc1, nightly-2026-07-28) is published as a prerelease, so it
# never becomes the repository's "latest".
if printf '%s' "${tag}" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$'; then
  prerelease=false
else
  prerelease=true
fi
```

Following what `release plan` asks — creating a `v0.8.0` tag to unblock it —
would **publish the release as a prerelease**, and it would never become the
latest. The tool asks for exactly the format that breaks the consumer. It nearly
happened: the maintainer authorized the release as `v0.8.0`, and the prefix only
stayed out because the workflow was read before the tag was created.

This is the second divergence of the same kind in one session. The first is
`docs/backlog/2026-08-12-archive-refuses-a-graph-that-declined-the-qa-gate.md`.
In both cases Roundfix assumes a convention the target repository does not
follow, and fails closed without offering to accommodate the real one. The cost
is not the command that fails — it is the operator concluding they must change
the repository to please the tool.

## Where

`internal/releaseplan/version.go` — `ParseStableVersion` requires the `v` prefix
(line 54) and the next action prescribes the format (line 37).

## Expected

Accept a stable tag with or without the `v` prefix, detecting the pattern from
the repository's existing tags rather than requiring one. When there are none,
the message may ask for a base without prescribing a specific format.

Worth checking in the same work whether `release plan --reset-to` and the range
computation carry the same assumption.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-11-release-plan-exige-tag-com-prefixo-v.md` in the
Secondbrain, captured from a `fluxus` session.

---
type: fix
status: promoted
created: 2026-10-18
spec: 0282-the-linux-lane-runs-like-ci
reason: null
---

# The Linux Lane container lacks the CI tools and an init

## Symptom

The first Linux Lane run that actually executed, on Spec 0280's candidate on
2026-10-18, failed the Delivery Queue gate (`gate-failed`) on two environment
gaps, while every package the candidate changed passed in the container:

- `TestGuidanceCompositionJourney/rust-cli` failed with `exec: "rustfmt":
  executable file not found in $PATH`. The `golang:1.26` image has no Rust
  toolchain, and the CI runner does.
- `TestCLITestHelpersEndWithTheirTestBinary` timed out after 594 s, because a
  recorded helper identity never ended. The container runs without an init,
  so an orphaned helper is never reaped.

The maintainer turned the lane off in this repository's `.roundfixrc.yml`
the same day.

## Where

`internal/cli/deliver_linux_lane.go`, which builds the `docker run` command
and chooses the image (Spec 0279, ADR-0288).

## Expected

The lane runs with `--init`, and on an image that carries the tools the
repository's tests need. That could be a configurable image, or the tests
that need an external tool could skip with a reason when it is absent, as
they do on a machine without it. A lane failure caused by the container is
then told apart from a failure of the change. After that, this repository
turns the lane back on.

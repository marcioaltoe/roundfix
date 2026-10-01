---
status: accepted
created_at: 2026-10-01T19:00:00Z
updated_at: 2026-10-01T19:00:00Z
deprecated_at: null
superseded_by: null
---

# A test fixture process ends with the test binary that started it

A process a test starts, directly or through a detached Run, must end when the test binary that started it ends. A `t.Cleanup` cannot guarantee this, because cleanups never run when the binary is killed by a timeout, `Ctrl-C` or a cancelled job. Every fixture that waits therefore watches its owning test binary's process and exits when that process is gone, and a test that proves a child survives its caller reaps that child at teardown by the PID the Run Database records. A fixture that waits only for a release file outlived its test by five hours on 2026-09-30, and by days in August, so a cooperative release alone is not enough.

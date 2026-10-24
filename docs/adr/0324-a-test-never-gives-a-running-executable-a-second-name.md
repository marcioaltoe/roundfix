---
status: accepted
created_at: 2026-10-24T12:00:00Z
updated_at: 2026-10-24T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A test never gives a running executable a second name

On 2026-10-23 two Runs lost their QA precondition to `internal/agent` tests
that had not changed for days: `TestACPXExitCodeMapping/agent_protocol_failure`
at 18:14 BRT (Spec 0304) and `TestRunPromptIgnoresAMalformedUsagePayload` at
23:10 BRT (Spec 0306). Both read an `args.json` that the fake acpx, the test
binary re-executed through `os.Args[0]`, never wrote. Neither failure
reproduced in twenty sequential runs. On 2026-10-24 a load harness of twelve
concurrent copies of the package reproduced it, and a trace showed the child
dying of `SIGKILL` before the first statement of `TestMain`. The unified log
named the sender: `syspolicyd` wrote "Unable (errno: 2) to read file at
<private> for pid: N" followed by "Terminating process due to Malware
rejection: N" for the child's PID, at 18:14:41 and 23:10:10 BRT on 2026-10-23
and for every child the harness lost.

macOS checks every new process of an ad-hoc signed binary through
`syspolicyd`. It resolves the executable's path from the process's file, and
when that path no longer exists it rejects the process and kills it. The
`internal/agent` `TestMain` gives the test binary three more names, hard links
under a temporary `node_modules` tree, and removes them when the process ends.
Re-executed children of the same binary run the same `TestMain`, create and
remove their own hard links while their parent starts other children, and the
fixture child that exits through `os.Exit` never removes its links at all: 26
leaked trees held 184 MiB of deleted test binaries under `$TMPDIR` on
2026-10-24. A process started while one of those names disappears can be
resolved to the vanished name and killed. A standalone program that started
itself 101,438 times while four goroutines created and removed hard links to
its own executable lost 48 children to `SIGKILL`, each matched by one
`syspolicyd` rejection; the same program with no hard link lost none in
115,022.

An earlier decision fixed a sibling of this class: the kernel kills a signed
executable rewritten in place while another process still maps it. Here
nothing is rewritten; a name disappears.

The published record agrees. A build-system maintainer reported on the Apple
Developer Forums (https://developer.apple.com/forums/thread/663456) that "in
case of 'File not found' macOS kills some other running process sharing the
inode with the already removed file", and that he "was able to avoid inode
sharing in my build process by using clonefile(2) instead of hardlink". The
Flutter tool lost child processes the same way after deleting its own running
binary (https://github.com/christopherfujino/catalina-crasher-demo, Flutter
issue 33890).

## Decision

**A test never gives an executable a second hard link.** A test that needs an
executable at another path copies it, or points a symbolic link at it, and
never creates a hard link to it. An executable a test may start keeps a single
directory entry while any process may start it: a **Single-Name Test Binary**.
A re-executed child that does not need the copies never makes them.

## Considered Options

- Keeping the hard links and removing them only in the top-level process was
  rejected: it still lets one process remove a name while another process of
  the same inode starts, as concurrent packages and repeated suites do.
- Retrying a fake process that died of a signal was rejected: it hides a
  kill the operating system issued for a reason the test controls, and the
  repository's flaky-test rule forbids rerunning until a pass.
- Disabling the check on the host was rejected: it is a machine setting the
  repository cannot own, and it would not hold on another Mac.

## Consequences

- The `internal/agent` package adapters are three copies of the test binary,
  about 10 MB each, made before any test starts and removed after the last
  one; the re-executed fixture children make none and leak none.
- A repository rule in `internal/testfixture` fails on any hard link created
  by a test, so the class cannot return through another package.
- A symbolic link stays allowed: it is not a name of the file's inode, so its
  removal never removes the path `syspolicyd` resolves.

# Controlled Force Stop harness

The live CLI checks used two temporary binaries built under `/private/tmp`.
They are recorded here as non-executable evidence so the QA artifact does not
add a Go package or executable file to the repository.

The sleeper printed its PID and waited for `SIGINT` or `SIGTERM`:

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Printf("pid=%d\n", os.Getpid())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
```

The seed helper opened an isolated Run Database through
`internal/store.Open`, called `internal/store.OwnerProcessIdentity` for the
matching-token case, and created one Active implement Run per unique QA Spec
slug. Its identity modes were:

```go
switch mode {
case "live":
	identity, err = store.OwnerProcessIdentity(context.Background(), pid)
case "legacy":
	identity = "legacy-ps-token"
case "mismatch":
	identity = runtime.GOOS + ":0.0"
case "blank":
	identity = ""
}
```

Every value-delivering action then entered through `./bin/roundfix stop` or
`./bin/roundfix runs list`. The helper only arranged controlled data and
process preconditions.

For the no-subprocess probe, `PATH` began with a temporary executable named
`ps` whose complete body was:

```sh
#!/bin/sh
touch "${QA_PS_MARKER:?}"
exit 99
```

After the unreadable proof attempt, `/bin/test ! -e
/private/tmp/roundfix-qa-0055-20260731/ps-invoked` exited `0`.

---
type: fix
status: promoted
created: 2026-10-26
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
reason: null
---

# The Hard Link rule misses a dot import and a function value

## Symptom

`TestNoTestCreatesAHardLink` (`internal/testfixture`) reports a hard link
only when a test calls `X.Link(...)` or `X.Linkat(...)` through a named
import of `os`, `syscall` or `golang.org/x/sys/unix`. A test that imports
`os` with a dot (`import . "os"`) and calls `Link(...)`, or that takes
`os.Link` as a function value (`link := os.Link`, `apply(os.Link)`) and calls
it later, creates the same hard link and passes the rule.

## Where

`internal/testfixture/hard_link_rule_test.go` `hardLinkViolations`: it
inspects only `*ast.CallExpr` whose `Fun` is a selector on an identifier
bound to an import, skips every import whose name is `.`, and walks only
`internal/`.

## Expected

The rule reports every reference to a hard-link function of those packages
in a test file, whatever its form: a call or a value, through a named, an
aliased or a dot import. It keeps ignoring a local identifier, a method or
field of another type, a struct literal key, a label, `Symlink` and
`Readlink`, and it names each finding's file and line as today.

## Evidence

The pre-PR review of Spec 0311 on 2026-10-24 raised it as F1 ("the rule does
not catch a dot import of os or os.Link used as a function value"); the
operator dismissed it at 10:45 BRT because it lay outside API Contract 2 of
that Spec's TechSpec, and listed it as triage item 9 at the end of cycle AP:
"Regra TestNoTestCreatesAHardLink (0311): não pega dot import de os nem
os.Link como valor de função (review F1 dispensado por escopo). Ampliar numa
Spec futura se valer."

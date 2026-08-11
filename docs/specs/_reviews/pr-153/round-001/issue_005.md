---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/agent/acpx_runner.go
line: 1181
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWiY,comment:PRRC_kwDOS0qyts7fswNz
review_hash: 8aa08a0b8f77c6286ffae1552fd4c9e8b1e2c6052f2ab06f51d83166733206c1
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:13Z"
---

# Issue 005: _ Data Integrity & Integration_ _ Trivial_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🌐 Web query:

`acpx CLI sessions close exit codes missing named session`

💡 Result:

The acpx CLI uses a defined set of exit codes to signal the state of sessions, including scenarios where a session is missing. When you attempt to use a named session that does not exist or has not been explicitly initialized, acpx exits with code 4 [1][2]. This exit code specifically signals that no session was found in the current working directory (cwd) scope and that an explicit call to sessions new or sessions ensure is required [1][3][2]. Key points regarding sessions and exit codes: - Exit Code 4: Indicates that no matching session was found [2]. To resolve this, you must initialize the session using acpx codex sessions new --name <name> or acpx codex sessions ensure --name <name> (where <name> is your chosen session identifier) [1][3][4]. - Named Sessions: Sessions are scoped to the current directory and the agent command [1][5]. Using the -s or --session flag allows multiple parallel, independent workstreams to coexist within the same repository [3][6][7]. - Session Management: Commands like sessions new, sessions ensure, and sessions close operate on these named sessions [4][8]. While sessions new creates a fresh session, sessions ensure acts idempotently, returning an existing session or creating one if missing [3][5]. For script automation, you can use these exit codes to handle missing sessions gracefully, as demonstrated in the official documentation [2]: case $? in 4) echo "no session — run sessions new"; acpx codex sessions new;; #... other cases esac

Citations:

- 1: https://raw.githubusercontent.com/openclaw/acpx/main/docs/CLI.md
- 2: https://acpx.sh/exit-codes.html
- 3: https://acpx.sh/sessions.html
- 4: https://github.com/openclaw/acpx/blob/main/docs/CLI.md
- 5: https://acpx.sh/quickstart.html
- 6: https://github.com/openclaw/acpx
- 7: https://acpx.sh/prompting.html
- 8: https://acpx.sh/session-control.html

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- target code ---'
sed -n '1135,1205p' internal/agent/acpx_runner.go

printf '%s\n' '--- acpx references and version constraints ---'
rg -n -S --hidden --glob '!vendor/**' --glob '!node_modules/**' \
  'MinimumACPXVersion|acpxExitReasonMissingSession|No named session|ExitCode == 4|sessions close|acpx' \
  internal docs .github go.mod package.json 2>/dev/null | head -n 240

printf '%s\n' '--- repository files related to acpx ---'
git ls-files | rg '(^|/)(acpx|.*acpx.*|go\.mod|README|CHANGELOG|docs)' | head -n 160
```

Repository: marcioaltoe/roundfix

Length of output: 39052

---

🌐 Web query:

`site:github.com/openclaw/acpx OR site:acpx.sh acpx "No named session" exit code 1`

💡 Result:

The error message "No named session" combined with "exit code 1" in acpx typically indicates a runtime or agent failure when attempting to interact with a session [1][2]. While an exit code of 4 is explicitly defined for cases where no session is found and requires a new session to be created [1][2], an exit code of 1 is a general catch-all for agent, protocol, or runtime errors [2]. If you encounter this, consider the following: 1. Agent/Protocol Error: Exit code 1 signifies that the acpx client successfully identified the session scope but failed to communicate with the underlying agent process or encountered a protocol error [2][3]. This can happen if the agent process crashes, encounters a network issue, or fails to initialize [2][4]. 2. Session Resolution: Ensure that if you are using named sessions (via the -s or --session flag), the session has been explicitly created using the sessions new command [1][5][6]. If the session is missing, acpx might throw an error [7]. 3. Debugging: Check the standard error (stderr) output for more specific details regarding the failure [2][4]. The acpx JSON mode or logs may provide further insight into the specific agent or transport failure [2][8]. For further troubleshooting, you can list your existing sessions using the command: acpx <agent-name> sessions Verify the status of your session with: acpx <agent-name> status If you are using named sessions, ensure your flag placement is correct (e.g., placing the -s flag after the subcommand if necessary, depending on your version) to ensure the session is properly resolved [7].

Citations:

- 1: https://acpx.sh/CLI.html
- 2: https://acpx.sh/exit-codes.html
- 3: https://github.com/openclaw/acpx/blob/main/src/acp/client.ts
- 4: https://github.com/openclaw/acpx/issues/176
- 5: https://acpx.sh/sessions.html
- 6: https://acpx.sh/quickstart.html
- 7: https://github.com/openclaw/acpx/issues/371/linked_closing_reference?reference_location=REPO_ISSUES_INDEX
- 8: https://acpx.sh/

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- minimum supported version and exit-reason definitions ---'
rg -n -A8 -B8 \
  'MinimumACPXVersion|acpxExitReasonMissingSession|type InfrastructureError|func \(runner \*ACPXRunner\) CloseSession' \
  internal/agent

printf '%s\n' '--- CloseSession-related tests ---'
rg -n -A18 -B12 \
  'CloseSession|close.*session|missing.*session|No named session|ExitCode: 4|ExitCode: 1' \
  internal/agent/*_test.go internal/daemon/*_test.go docs/backlog docs/adr \
  | head -n 360

printf '%s\n' '--- relevant backlog and ADR excerpts ---'
for f in \
  docs/backlog/2026-08-08-a-failed-proof-appends-a-cleanup-error-the-maintainer-cannot-act-on.md \
  docs/backlog/2026-08-08-a-session-that-never-opened-is-a-selection-failure.md \
  docs/adr/0039-model-availability-preflight-uses-a-disposable-agent-session.md \
  docs/adr/0114-opening-an-agent-session-is-not-agent-work.md
do
  if [ -f "$f" ]; then
    echo "--- $f"
    cat "$f"
  fi
done
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
from pathlib import Path
import re

source = Path("internal/agent/acpx_runner.go").read_text()
match = re.search(
    r'if errors\.As\(err, &infrastructureErr\) &&\s*'
    r'\(infrastructureErr\.ExitCode == 4 \|\|\s*'
    r'\(infrastructureErr\.ExitCode == 1 && strings\.Contains\('
    r'infrastructureErr\.Stderr, "No named session"\)\)\)',
    source,
)
if not match:
    raise SystemExit("CloseSession classification expression not found")

cases = [
    (4, "", True),
    (4, "unrelated failure", True),
    (1, "No named session", True),
    (1, "prefix: No named session; suffix", True),
    (1, "no named session", False),
    (1, "session protocol failure", False),
    (2, "No named session", False),
]
for exit_code, stderr, expected in cases:
    actual = exit_code == 4 or (
        exit_code == 1 and "No named session" in stderr
    )
    if actual != expected:
        raise SystemExit((exit_code, stderr, actual, expected))
    print(f"exit={exit_code!r}, stderr={stderr!r} -> missing={actual}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 207

---

🌐 Web query:

`site:acpx.sh/exit-codes.html acpx exit code 4 missing session exit code 1`

💡 Result:

The exit code 4 in acpx explicitly indicates that no active session was found, meaning the requested operation requires an explicit call to sessions new before it can proceed [1]. The exit code 1 represents a generic agent, protocol, or runtime error [1]. If you encounter these codes, you should verify if your workflow requires initializing a session first. The official documentation recommends using sessions new or sessions ensure (for idempotent scripts) when receiving exit code 4 [1]. For exit code 1, since it is a general failure, you should check your agent configuration, network connectivity, or the command protocol used [1]. Example of handling these codes in a script: if acpx --format quiet codex 'command' >output.txt; then echo "Success" else case $? in 1) echo "Agent, protocol, or runtime error";; 4) echo "No session found — running sessions new"; acpx codex sessions new;; *) echo "Other error occurred";; esac fi [1]

Citations:

- 1: https://acpx.sh/exit-codes.html

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- missing-session tests ---'
sed -n '1088,1185p' internal/agent/acpx_runner_test.go

printf '%s\n' '--- cleanup implementation and callers ---'
sed -n '1,80p' internal/agent/sessions.go
rg -n -A25 -B12 \
  'IsAgentSessionAbsent|AgentSessionCleanupError|close disposable Agent Session|EndSession\(' \
  internal/agent internal/daemon | head -n 260

printf '%s\n' '--- all exact references to the undocumented wording ---'
rg -n -F 'No named session' .
```

Repository: marcioaltoe/roundfix

Length of output: 48608

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
from pathlib import Path

source = Path("internal/agent/acpx_runner.go").read_text()
start = source.index("if errors.As(err, &infrastructureErr)")
end = source.index("\n\t\treturn fmt.Errorf", start)
expression = source[start:end]
print(expression)

expected = {
    (4, ""): True,
    (4, "unrelated failure"): True,
    (1, "No named session"): True,
    (1, "prefix: No named session; suffix"): True,
    (1, "no named session"): False,
    (1, "session protocol failure"): False,
    (2, "No named session"): False,
}

for (exit_code, stderr), want in expected.items():
    got = exit_code == 4 or (
        exit_code == 1 and "No named session" in stderr
    )
    print(f"exit={exit_code}, stderr={stderr!r}, classified_missing={got}")
    assert got == want
PY
```

Repository: marcioaltoe/roundfix

Length of output: 934

---



</details>

**Document the exit-1 compatibility marker**

`acpx` documents exit code 4 for a missing session, but the supported `0.12.0` build also returns exit code 1 with `No named session`. Keep both cases. Extract the stderr marker into a named constant with the observed version documented. Add a regression test for an unrelated exit-1 error so wording changes cannot produce a false cleanup error.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner.go` around lines 1172 - 1181, Update the
InfrastructureError normalization logic around errors.As to use a named constant
for the “No named session” stderr marker, documenting that it is required for
the supported acpx 0.12.0 exit-1 behavior while retaining exit code 4 handling.
Add a regression test covering an unrelated exit-1 error and verify it is not
classified as acpxExitReasonMissingSession.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:160ffc50b6b5cdaa9141ec1d -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Extracted the stderr marker into a named constant `acpxMissingSessionStderrMarker = "No named session"` in internal/agent/acpx_runner.go, with a comment documenting that the pinned acpx 0.12.0 build returns exit 1 with this marker in addition to the documented exit code 4. Exit-code-4 handling is retained unchanged, and `TestUnrelatedExitOneKeepsItsClassification` already provides the unrelated-exit-1 regression coverage asserting it is not classified as `acpxExitReasonMissingSession`. Focused: `go test ./internal/agent -run 'TestUnrelatedExitOne|TestMissingSession'` ok.

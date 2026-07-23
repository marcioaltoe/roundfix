"""Project-neutral Verification for disposable profile repositories."""

import json
from pathlib import Path


VERIFICATION_COMMAND = "python3 -B .macro-profile-verify.py"


def main():
    manifest_path = Path("docs/agents/setup-context.json")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    configured = manifest["decisions"]["verification.gate"]["value"]
    if configured != VERIFICATION_COMMAND:
        raise SystemExit(
            f"persisted Verification is {configured!r}, want {VERIFICATION_COMMAND!r}"
        )

    paths = [
        Path("AGENTS.md"),
        *sorted(Path("docs/agents").glob("*.md")),
    ]
    if not paths or any(not path.is_file() for path in paths):
        raise SystemExit("generated Markdown corpus is incomplete")
    for path in paths:
        content = path.read_bytes()
        if not content.endswith(b"\n"):
            raise SystemExit(f"{path}: missing final newline")
        if b"\r\n" in content:
            raise SystemExit(f"{path}: CRLF is not formatter-stable")
        if any(line.endswith((b" ", b"\t")) for line in content.splitlines()):
            raise SystemExit(f"{path}: trailing whitespace is not formatter-stable")


if __name__ == "__main__":
    main()

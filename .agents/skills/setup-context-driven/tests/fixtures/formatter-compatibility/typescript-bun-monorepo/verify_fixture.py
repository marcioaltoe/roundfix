"""Selected Verification for the formatter composition fixture."""

from pathlib import Path


def main():
    paths = [
        Path("AGENTS.md"),
        *sorted(Path("docs/agents").glob("*.md")),
    ]
    if not paths or any(not path.is_file() for path in paths):
        raise SystemExit("generated Markdown corpus is incomplete")
    for path in paths:
        content = path.read_text(encoding="utf-8")
        if not content.endswith("\n"):
            raise SystemExit(f"{path}: missing final newline")
        if "setup-context-driven:begin" in content and "-->\n\n" not in content:
            raise SystemExit(f"{path}: managed begin marker is not formatter-stable")
        if "setup-context-driven:end" in content and "\n\n<!-- setup-context-driven:end" not in content:
            raise SystemExit(f"{path}: managed end marker is not formatter-stable")


if __name__ == "__main__":
    main()

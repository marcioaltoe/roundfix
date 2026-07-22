from pathlib import Path


def repository_root(start):
    current = Path(start).resolve()
    if current.is_file():
        current = current.parent
    for candidate in (current, *current.parents):
        if (
            (candidate / "go.mod").is_file()
            and (candidate / ".agents" / "skills" / "setup-context-driven").is_dir()
            and (candidate / "skills" / "setup-context-driven").is_dir()
        ):
            return candidate
    raise RuntimeError(f"could not find Roundfix repository root from {current}")


def setup_skill_roots(start):
    root = repository_root(start)
    return (
        root / ".agents" / "skills" / "setup-context-driven",
        root / "skills" / "setup-context-driven",
    )

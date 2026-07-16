#!/usr/bin/env python3
"""Manual-style variation QA for setup-context-driven.

Suite: context-driven setup CLI dogfood variations
Invariant: confirmed decisions produce safe, portable, idempotent managed guidance.
Boundary IN: real context_setup.py subprocess, bundled assets, canonical setup checkout, disposable repos.
Boundary OUT: network access, live Secondbrain reads, skill installation, and deletion of installed skills.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import shutil
import subprocess
import sys
import tempfile
from dataclasses import dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Callable


SPEC_RELATIVE = Path("docs/specs/0030-context-driven-agent-instructions")
REQUIRED_SECONDBRAIN_PHRASES = (
    "wiki/index.md",
    "qmd query",
    "projects/<project>/mirror",
    "Cite every Secondbrain file",
    "Do not write to the Secondbrain",
    "Do not edit raw/",
    "Do not edit projects/*/mirror/",
    "Hermes",
    "Never read, copy, or expose",
)


class CheckFailure(AssertionError):
    pass


@dataclass
class Scenario:
    qa_id: str
    name: str
    repo_root: Path
    fixture: Path
    evidence_dir: Path
    commands: list[dict] = field(default_factory=list)
    checks: list[dict] = field(default_factory=list)
    notes: list[str] = field(default_factory=list)

    @property
    def cli(self) -> Path:
        return self.repo_root / ".agents/skills/setup-context-driven/scripts/context_setup.py"

    def run(self, args: list[str], *, cwd: Path | None = None) -> subprocess.CompletedProcess[str]:
        command = [sys.executable, str(self.cli), *args]
        result = subprocess.run(
            command,
            cwd=cwd or self.repo_root,
            env={**os.environ, "HOME": str(self.fixture / "home")},
            text=True,
            capture_output=True,
            check=False,
        )
        self.commands.append(
            {
                "argv": sanitize_paths(command, self),
                "cwd": sanitize_text(str(cwd or self.repo_root), self),
                "exitCode": result.returncode,
                "stdout": sanitize_text(result.stdout, self),
                "stderr": sanitize_text(result.stderr, self),
            }
        )
        return result

    def run_repo_command(self, args: list[str]) -> subprocess.CompletedProcess[str]:
        result = subprocess.run(
            args,
            cwd=self.repo_root,
            text=True,
            capture_output=True,
            check=False,
        )
        self.commands.append(
            {
                "argv": args,
                "cwd": "<source-repo>",
                "exitCode": result.returncode,
                "stdout": result.stdout,
                "stderr": result.stderr,
            }
        )
        return result

    def check(self, condition: bool, message: str, actual: object = None) -> None:
        self.checks.append({"message": message, "passed": bool(condition), "actual": actual})
        if not condition:
            raise CheckFailure(f"{message}; actual={actual!r}")

    def check_equal(self, actual: object, expected: object, message: str) -> None:
        self.check(actual == expected, message, {"expected": expected, "actual": actual})

    def write_evidence(self, status: str, error: str | None = None) -> None:
        payload = {
            "schemaVersion": "setup-context-driven/manual-qa-v1",
            "qaId": self.qa_id,
            "name": self.name,
            "status": status,
            "recordedAt": datetime.now(timezone.utc).isoformat(),
            "commands": self.commands,
            "checks": self.checks,
            "notes": self.notes,
            "error": error,
        }
        path = self.evidence_dir / f"{self.qa_id.lower()}-{slug(self.name)}.json"
        path.write_text(
            json.dumps(payload, indent=2, sort_keys=False, default=json_fallback) + "\n",
            encoding="utf-8",
        )


def sanitize_text(value: str, scenario: Scenario) -> str:
    return value.replace(str(scenario.fixture), "<fixture>").replace(
        str(scenario.repo_root), "<source-repo>"
    )


def sanitize_paths(values: list[str], scenario: Scenario) -> list[str]:
    return [sanitize_text(value, scenario) for value in values]


def slug(value: str) -> str:
    return "-".join("".join(character.lower() if character.isalnum() else " " for character in value).split())


def json_fallback(value: object) -> object:
    if isinstance(value, set):
        return sorted(value)
    if isinstance(value, bytes):
        return {
            "bytes": len(value),
            "sha256": hashlib.sha256(value).hexdigest(),
        }
    raise TypeError(f"Object of type {type(value).__name__} is not JSON serializable")


def parse_json(result: subprocess.CompletedProcess[str]) -> dict:
    try:
        payload = json.loads(result.stdout)
    except json.JSONDecodeError as error:
        raise CheckFailure(f"stdout is not JSON: {error}: {result.stdout!r}") from error
    if not isinstance(payload, dict):
        raise CheckFailure(f"JSON result is not an object: {payload!r}")
    return payload


def finding_codes(payload: dict) -> list[str]:
    return [item.get("code", "") for item in payload.get("findings", [])]


def managed_ids(payload: dict, code: str) -> list[str]:
    return [
        item.get("managedId", "")
        for item in payload.get("findings", [])
        if item.get("code") == code
    ]


def tree_snapshot(root: Path) -> dict[str, str]:
    if not root.exists():
        return {}
    snapshot: dict[str, str] = {}
    for path in sorted(root.rglob("*")):
        if path.is_file():
            snapshot[path.relative_to(root).as_posix()] = hashlib.sha256(path.read_bytes()).hexdigest()
    return snapshot


def production_snapshot(repo_root: Path) -> dict[str, str]:
    snapshot: dict[str, str] = {}
    for relative in (
        Path(".agents/skills/setup-context-driven"),
        Path("skills/setup-context-driven"),
    ):
        for item, digest in tree_snapshot(repo_root / relative).items():
            snapshot[(relative / item).as_posix()] = digest
    return snapshot


def read_setup(repo_root: Path, profile: str) -> dict:
    profile_path = repo_root / f".agents/skills/setup-context-driven/assets/profiles/{profile}.json"
    profile_data = json.loads(profile_path.read_text(encoding="utf-8"))
    setup_id = profile_data["setup"]
    setup_path = repo_root / f".agents/skills/setup-context-driven/assets/setups/{setup_id}.json"
    return json.loads(setup_path.read_text(encoding="utf-8"))


def install_profile_skills(repo_root: Path, fixture: Path, profile: str, omit: set[str] | None = None) -> None:
    omit = omit or set()
    for skill in read_setup(repo_root, profile)["skills"]:
        name = skill["name"]
        if name in omit:
            continue
        skill_path = fixture / ".agents" / "skills" / name / "SKILL.md"
        skill_path.parent.mkdir(parents=True, exist_ok=True)
        skill_path.write_text(f"---\nname: {name}\n---\n", encoding="utf-8")


def write_skill(fixture: Path, name: str) -> None:
    path = fixture / ".agents" / "skills" / name / "SKILL.md"
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(f"---\nname: {name}\n---\n", encoding="utf-8")


def write_lock(fixture: Path, skills: list[str]) -> None:
    payload = {
        "version": 1,
        "skills": {
            name: {
                "source": "qa-fixture",
                "sourceType": "local",
                "skillPath": f"skills/{name}/SKILL.md",
                "computedHash": "0" * 64,
            }
            for name in skills
        },
    }
    (fixture / "skills-lock.json").write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")


def decisions(
    *,
    secondbrain: bool,
    spec_scaffold: bool = True,
    domain_layout: str = "single-context",
    triage_external: bool = False,
    autonomous: bool = True,
    verification: str = "make verify",
    backend: str = "codex gpt-5.5 xhigh",
    design: str = "claude opus-4.8 xhigh",
) -> list[str]:
    values = {
        "spec.scaffold": str(spec_scaffold).lower(),
        "domain.layout": domain_layout,
        "triage.external": str(triage_external).lower(),
        "autonomous.enabled": str(autonomous).lower(),
        "runtime.backend": backend,
        "runtime.design": design,
        "verification.gate": verification,
        "language.generated": "English",
        "secondbrain.enabled": str(secondbrain).lower(),
    }
    return [f"{key}={value}" for key, value in values.items()]


def apply(scenario: Scenario, profile: str, answers: list[str]) -> subprocess.CompletedProcess[str]:
    args = ["apply", "--repo", str(scenario.fixture), "--format", "json", "--profile", profile]
    for answer in answers:
        args.extend(["--decision", answer])
    return scenario.run(args)


def audit(scenario: Scenario, *extra: str) -> subprocess.CompletedProcess[str]:
    return scenario.run(["audit", "--repo", str(scenario.fixture), "--format", "json", *extra])


def assert_clean_profile(scenario: Scenario, profile: str, answers: list[str]) -> None:
    applied = apply(scenario, profile, answers)
    scenario.check_equal(applied.returncode, 0, "apply exits successfully")
    install_profile_skills(scenario.repo_root, scenario.fixture, profile)
    audited = audit(scenario)
    scenario.check_equal(audited.returncode, 0, "audit exits successfully after skills are installed")
    payload = parse_json(audited)
    scenario.check_equal(payload["findings"], [], "clean audit has no findings")


def qa_02(s: Scenario) -> None:
    before = tree_snapshot(s.fixture)
    result = s.run(["--repo", str(s.fixture), "--format", "json"])
    payload = parse_json(result)
    s.check_equal(result.returncode, 1, "default audit exits 1 for missing manifest")
    s.check_equal(finding_codes(payload), ["manifest.missing"], "default audit reports stable missing-manifest code")
    s.check_equal(tree_snapshot(s.fixture), before, "default audit performs no writes")


def qa_03(s: Scenario) -> None:
    (s.fixture / "AGENTS.md").write_text("# Repository-owned instructions\n", encoding="utf-8")
    before = tree_snapshot(s.fixture)
    result = apply(s, "typescript-bun-monorepo", [])
    payload = parse_json(result)
    expected = {
        "spec.scaffold",
        "domain.layout",
        "triage.external",
        "autonomous.enabled",
        "runtime.backend",
        "runtime.design",
        "verification.gate",
        "language.generated",
        "secondbrain.enabled",
    }
    s.check_equal(result.returncode, 3, "unanswered first apply exits 3")
    s.check_equal(set(managed_ids(payload, "decision.required")), expected, "only required durable questions are returned")
    s.check(len(payload.get("plannedChanges", [])) > 1, "first-run response previews all managed files and blocks", payload.get("plannedChanges"))
    s.check_equal(tree_snapshot(s.fixture), before, "unanswered apply performs no writes")


def qa_04(s: Scenario) -> None:
    (s.fixture / "AGENTS.md").write_text("# Owner section\n\nKeep this exact line.\n", encoding="utf-8")
    answers = decisions(secondbrain=True, domain_layout="multi-context", triage_external=True)
    assert_clean_profile(s, "typescript-bun-monorepo", answers)
    root = (s.fixture / "AGENTS.md").read_text(encoding="utf-8")
    guide = (s.fixture / "docs/agents/secondbrain.md").read_text(encoding="utf-8")
    manifest = json.loads((s.fixture / "docs/agents/setup-context.json").read_text(encoding="utf-8"))
    s.check("Keep this exact line." in root, "root owner content remains present")
    s.check("root.secondbrain" in root, "root contains a managed Secondbrain pointer")
    s.check("secondbrain" in manifest["modules"], "manifest resolves the optional Secondbrain module")
    for phrase in REQUIRED_SECONDBRAIN_PHRASES:
        s.check(phrase in guide, f"Secondbrain guide contains required rule: {phrase}")
    s.check(" não " not in f" {guide.lower()} ", "generated guide remains English")


def qa_05(s: Scenario) -> None:
    assert_clean_profile(s, "typescript-bun-monorepo", decisions(secondbrain=True))
    before = tree_snapshot(s.fixture)
    repeated = apply(s, "typescript-bun-monorepo", [])
    payload = parse_json(repeated)
    s.check_equal(repeated.returncode, 0, "stored decisions permit answer-free reapply")
    s.check("decision.required" not in finding_codes(payload), "stored decisions are not asked again")
    s.check_equal(tree_snapshot(s.fixture), before, "second apply is byte-for-byte idempotent")


def qa_06(s: Scenario) -> None:
    assert_clean_profile(s, "go-cli-tui", decisions(secondbrain=False))
    root = (s.fixture / "AGENTS.md").read_text(encoding="utf-8")
    s.check("root.secondbrain" not in root, "Secondbrain pointer is absent after opt out")
    s.check(not (s.fixture / "docs/agents/secondbrain.md").exists(), "Secondbrain guide is absent after opt out")


def qa_07(s: Scenario) -> None:
    answers = decisions(
        secondbrain=False,
        spec_scaffold=False,
        domain_layout="multi-context",
        triage_external=True,
        autonomous=False,
        verification="cargo test --all-targets",
        backend="codex alternate-backend xhigh",
        design="claude alternate-design xhigh",
    )
    assert_clean_profile(s, "rust-cli", answers)
    manifest_text = (s.fixture / "docs/agents/setup-context.json").read_text(encoding="utf-8")
    generated = "\n".join(
        path.read_text(encoding="utf-8")
        for path in [s.fixture / "AGENTS.md", *sorted((s.fixture / "docs/agents").glob("*.md"))]
    )
    s.check("cargo test --all-targets" in manifest_text, "verification answer is durable in the manifest")
    s.check("cargo test --all-targets" in generated, "verification answer is reflected in generated guidance")
    s.check("alternate-backend" in generated, "backend runtime answer is reflected in generated guidance")
    s.check("alternate-design" in generated, "design runtime answer is reflected in generated guidance")
    s.check("docs/specs/<feature-slug>" not in generated, "spec guidance does not contradict spec.scaffold=false")
    s.check("root.autonomous-work" not in generated, "autonomous guidance does not contradict autonomous.enabled=false")


def qa_08(s: Scenario) -> None:
    guide = s.fixture / "docs/agents/rust.md"
    guide.parent.mkdir(parents=True)
    guide.write_text("repository-authored Rust notes\n", encoding="utf-8")
    (s.fixture / "AGENTS.md").write_text("root owner bytes\n", encoding="utf-8")
    before = tree_snapshot(s.fixture)
    blocked = apply(s, "rust-cli", decisions(secondbrain=False))
    payload = parse_json(blocked)
    s.check_equal(blocked.returncode, 3, "unmarked guide requires adoption")
    s.check("adoption.guide.rust" in managed_ids(payload, "decision.required"), "adoption question identifies the exact guide")
    s.check_equal(tree_snapshot(s.fixture), before, "blocked adoption performs no writes")
    accepted = apply(s, "rust-cli", decisions(secondbrain=False) + ["adoption.guide.rust=true"])
    s.check_equal(accepted.returncode, 0, "explicit adoption succeeds")
    s.check("root owner bytes\n" in (s.fixture / "AGENTS.md").read_text(encoding="utf-8"), "unrelated root bytes survive adoption")
    s.check("guide.rust" in guide.read_text(encoding="utf-8"), "adopted guide receives stable ownership markers")


def qa_09(s: Scenario) -> None:
    assert_clean_profile(s, "rust-cli", decisions(secondbrain=True))
    root_path = s.fixture / "AGENTS.md"
    guide_path = s.fixture / "docs/agents/secondbrain.md"
    root_path.write_text("owner root prefix\n" + root_path.read_text(encoding="utf-8"), encoding="utf-8")
    guide_path.write_text("owner guide prefix\n" + guide_path.read_text(encoding="utf-8"), encoding="utf-8")
    disabled = apply(s, "rust-cli", ["secondbrain.enabled=false"])
    s.check_equal(disabled.returncode, 0, "disabling a stored opt-in succeeds")
    s.check("owner root prefix\n" in root_path.read_text(encoding="utf-8"), "root owner bytes survive opt out")
    s.check_equal(guide_path.read_text(encoding="utf-8"), "owner guide prefix\n", "guide retains only owner-authored bytes")
    s.check("root.secondbrain" not in root_path.read_text(encoding="utf-8"), "managed pointer is removed")


def qa_10(s: Scenario) -> None:
    applied = apply(s, "rust-cli", decisions(secondbrain=False))
    s.check_equal(applied.returncode, 0, "fixture apply succeeds")
    install_profile_skills(s.repo_root, s.fixture, "rust-cli", {"agentic-cli-design"})
    result = audit(s)
    payload = parse_json(result)
    s.check_equal(result.returncode, 1, "missing required skill blocks compliance")
    s.check("skills.required.missing" in finding_codes(payload), "stable missing-skill code is emitted")
    s.check("agentic-cli-design" in managed_ids(payload, "skills.required.missing") or any("agentic-cli-design" in item.get("path", "") for item in payload["findings"]), "finding identifies missing skill")


def qa_11(s: Scenario) -> None:
    assert_clean_profile(s, "rust-cli", decisions(secondbrain=False))
    write_skill(s.fixture, "locked-extra")
    write_skill(s.fixture, "untracked-extra")
    write_lock(s.fixture, ["locked-extra"])
    before = tree_snapshot(s.fixture)
    hidden = audit(s)
    shown = audit(s, "--show-extra-skills")
    hidden_payload = parse_json(hidden)
    shown_payload = parse_json(shown)
    s.check_equal(hidden.returncode, 0, "extras do not block default audit")
    s.check("skills.extra.installed" not in finding_codes(hidden_payload), "extra locked skills are hidden by default")
    s.check_equal(shown.returncode, 0, "opt-in cleanup report remains non-blocking")
    s.check("skills.extra.installed" in finding_codes(shown_payload), "locked extra is informational")
    s.check("skills.local.untracked" in finding_codes(shown_payload), "untracked extra is informational")
    actions = " ".join(item.get("action", "") for item in shown_payload["findings"]).lower()
    s.check("rm " not in actions and "delete" not in actions, "report contains no removal command")
    s.check_equal(tree_snapshot(s.fixture), before, "extra-skill report removes nothing")


def qa_12(s: Scenario) -> None:
    cases: list[tuple[str, Callable[[], subprocess.CompletedProcess[str]], int, str]] = []
    cases.append(("invalid decision", lambda: apply(s, "rust-cli", decisions(secondbrain=False) + ["language.generated=Portuguese"]), 3, "decision.required"))
    cases.append(("unknown profile", lambda: apply(s, "unknown-profile", decisions(secondbrain=False)), 1, "profile.unknown"))
    for name, invoke, expected_exit, expected_code in cases:
        before = tree_snapshot(s.fixture)
        result = invoke()
        payload = parse_json(result)
        s.check_equal(result.returncode, expected_exit, f"{name} uses the stable exit code")
        s.check(expected_code in finding_codes(payload), f"{name} uses the stable finding code")
        s.check_equal(tree_snapshot(s.fixture), before, f"{name} performs no writes")

    manifest = s.fixture / "docs/agents/setup-context.json"
    manifest.parent.mkdir(parents=True, exist_ok=True)
    manifest.write_text("{not-json", encoding="utf-8")
    before_invalid = tree_snapshot(s.fixture)
    invalid_manifest = apply(s, "rust-cli", decisions(secondbrain=False))
    invalid_payload = parse_json(invalid_manifest)
    s.check_equal(invalid_manifest.returncode, 2, "malformed manifest exits 2")
    s.check("manifest.invalid" in finding_codes(invalid_payload), "malformed manifest uses stable code")
    s.check_equal(tree_snapshot(s.fixture), before_invalid, "malformed manifest performs no writes")

    shutil.rmtree(s.fixture / "docs")
    duplicate = (
        "<!-- setup-context-driven:begin id=root.core version=1 -->\nfirst\n"
        "<!-- setup-context-driven:end id=root.core -->\n"
        "<!-- setup-context-driven:begin id=root.core version=1 -->\nsecond\n"
        "<!-- setup-context-driven:end id=root.core -->\n"
    )
    (s.fixture / "AGENTS.md").write_text(duplicate, encoding="utf-8")
    before_marker = tree_snapshot(s.fixture)
    invalid_marker = apply(s, "rust-cli", decisions(secondbrain=False))
    marker_payload = parse_json(invalid_marker)
    s.check_equal(invalid_marker.returncode, 1, "duplicate marker exits 1")
    s.check("managed.block.duplicate" in finding_codes(marker_payload), "duplicate marker uses stable code")
    s.check_equal(tree_snapshot(s.fixture), before_marker, "invalid marker performs no writes")


def qa_13(s: Scenario) -> None:
    source = Path("/Users/marcio/dev/skills/setups")
    s.check(source.is_dir(), "canonical setup checkout exists for maintainer drift check", str(source))
    result = s.run(["sync-setups", "--source-dir", str(source), "--check", "--format", "json"])
    payload = parse_json(result)
    s.check_equal(result.returncode, 0, "canonical setup snapshot check succeeds")
    s.check_equal(payload["findings"], [], "canonical setup snapshots have zero drift")


def qa_14(s: Scenario) -> None:
    assert_clean_profile(s, "go-cli-tui", decisions(secondbrain=False))
    result = audit(s)
    payload = parse_json(result)
    s.check_equal(result.returncode, 0, "normal audit needs no --setups-dir")
    s.check("skills.setup-snapshot.drift" not in finding_codes(payload), "normal audit does not depend on canonical checkout")


def qa_15(s: Scenario) -> None:
    executable = shutil.which("rtk")
    command = [executable, "make", "skills-sync-check"] if executable else ["make", "skills-sync-check"]
    result = s.run_repo_command(command)
    s.check_equal(result.returncode, 0, "canonical and embedded skill copies are synchronized")


def qa_16(s: Scenario) -> None:
    top = s.run(["--help"])
    s.check_equal(top.returncode, 0, "top-level help exits 0")
    s.check("audit" in top.stdout and "apply" in top.stdout and "sync-setups" in top.stdout, "help lists all public commands")
    s.check_equal(top.stderr, "", "help writes no diagnostics")
    unknown = s.run(["unknown-command"])
    s.check_equal(unknown.returncode, 2, "unknown invocation exits 2")
    s.check_equal(unknown.stdout, "", "unknown invocation keeps stdout empty")
    s.check("unrecognized arguments" in unknown.stderr, "unknown invocation writes argparse diagnostic to stderr")
    audit_help = s.run(["audit", "--help"])
    s.check("reserved" not in audit_help.stdout.lower(), "audit help does not label implemented checks as reserved")
    s.check("informational findings" in audit_help.stdout.lower(), "audit help truthfully describes implemented extra-skill report")
    s.check("compare the bundled setup snapshot" in audit_help.stdout.lower(), "audit help truthfully describes implemented setup drift check")
    result = s.run(["--repo", str(s.fixture), "--format", "json"])
    payload = parse_json(result)
    s.check_equal(payload["schemaVersion"], "setup-context-driven/audit-v1", "JSON schema id is stable")
    s.check_equal(result.stderr, "", "JSON findings do not leak diagnostics to stderr")


def qa_17(s: Scenario) -> None:
    nested = s.fixture / "packages/backend/AGENTS.md"
    nested.parent.mkdir(parents=True)
    nested.write_text("nested owner instructions\n", encoding="utf-8")
    secondbrain = s.fixture / "home/dev/secondbrain"
    secondbrain.mkdir(parents=True)
    sentinel = secondbrain / "sentinel.md"
    sentinel.write_text("never change this\n", encoding="utf-8")
    write_skill(s.fixture, "extra-never-remove")
    nested_before = nested.read_bytes()
    brain_before = tree_snapshot(secondbrain)
    skill_before = (s.fixture / ".agents/skills/extra-never-remove/SKILL.md").read_bytes()
    result = apply(s, "rust-cli", decisions(secondbrain=True))
    s.check_equal(result.returncode, 0, "safety fixture apply succeeds")
    s.check_equal(nested.read_bytes(), nested_before, "nested AGENTS.md is untouched")
    s.check_equal(tree_snapshot(secondbrain), brain_before, "Secondbrain workspace is untouched")
    s.check_equal((s.fixture / ".agents/skills/extra-never-remove/SKILL.md").read_bytes(), skill_before, "installed extra skill is untouched")


SCENARIOS: dict[str, tuple[str, Callable[[Scenario], None]]] = {
    "QA-02": ("Default audit is read only", qa_02),
    "QA-03": ("First apply questions and preview", qa_03),
    "QA-04": ("TypeScript Bun Secondbrain opt in", qa_04),
    "QA-05": ("Stored decisions and idempotency", qa_05),
    "QA-06": ("Go CLI TUI Secondbrain opt out", qa_06),
    "QA-07": ("Rust alternate decision behavior", qa_07),
    "QA-08": ("Unmarked guide adoption", qa_08),
    "QA-09": ("Secondbrain opt out cleanup", qa_09),
    "QA-10": ("Missing required skill", qa_10),
    "QA-11": ("Optional extra skill report", qa_11),
    "QA-12": ("Invalid input atomicity", qa_12),
    "QA-13": ("Canonical setup snapshot comparison", qa_13),
    "QA-14": ("Portable audit without canonical checkout", qa_14),
    "QA-15": ("Embedded skill synchronization", qa_15),
    "QA-16": ("CLI public contract", qa_16),
    "QA-17": ("Safety non goals", qa_17),
}


def run_static_gate(repo_root: Path, evidence_dir: Path) -> dict:
    executable = shutil.which("rtk")
    command = [executable, "make", "verify"] if executable else ["make", "verify"]
    result = subprocess.run(command, cwd=repo_root, text=True, capture_output=True, check=False)
    payload = {
        "schemaVersion": "setup-context-driven/manual-qa-v1",
        "qaId": "QA-01",
        "name": "Static repository verification gate",
        "status": "pass" if result.returncode == 0 else "fail",
        "recordedAt": datetime.now(timezone.utc).isoformat(),
        "commands": [
            {
                "argv": command,
                "cwd": "<source-repo>",
                "exitCode": result.returncode,
                "stdout": result.stdout,
                "stderr": result.stderr,
            }
        ],
        "checks": [
            {
                "message": "make verify passes without skipped checks",
                "passed": result.returncode == 0,
                "actual": result.returncode,
            }
        ],
        "notes": [],
        "error": None if result.returncode == 0 else "make verify failed",
    }
    (evidence_dir / "qa-01-static-repository-verification-gate.json").write_text(
        json.dumps(payload, indent=2) + "\n", encoding="utf-8"
    )
    return {"qaId": "QA-01", "status": payload["status"], "error": payload["error"]}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo-root", type=Path, default=Path(__file__).resolve().parents[4])
    parser.add_argument("--evidence-dir", type=Path, required=True)
    parser.add_argument("--scenario", action="append", choices=sorted(SCENARIOS))
    parser.add_argument("--include-static", action="store_true")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    repo_root = args.repo_root.resolve()
    evidence_dir = args.evidence_dir.resolve()
    if evidence_dir.exists() and any(evidence_dir.iterdir()):
        print(f"Evidence directory must be empty: {evidence_dir}", file=sys.stderr)
        return 2
    evidence_dir.mkdir(parents=True, exist_ok=True)
    before_source = production_snapshot(repo_root)
    results: list[dict] = []
    if args.include_static:
        results.append(run_static_gate(repo_root, evidence_dir))

    selected = args.scenario or list(SCENARIOS)
    for qa_id in selected:
        name, test = SCENARIOS[qa_id]
        with tempfile.TemporaryDirectory(prefix=f"setup-context-{qa_id.lower()}-") as temp_dir:
            fixture = Path(temp_dir)
            (fixture / "home").mkdir()
            scenario = Scenario(qa_id, name, repo_root, fixture, evidence_dir)
            status = "pass"
            error = None
            try:
                test(scenario)
            except Exception as caught:  # keep the full QA sweep running
                status = "fail"
                error = f"{type(caught).__name__}: {caught}"
            scenario.write_evidence(status, error)
            results.append({"qaId": qa_id, "status": status, "error": error})
            print(f"{qa_id} {status}: {name}")

    after_source = production_snapshot(repo_root)
    restored = before_source == after_source
    restore_payload = {
        "schemaVersion": "setup-context-driven/manual-qa-v1",
        "qaId": "QA-18",
        "name": "Source repository restoration",
        "status": "pass" if restored else "fail",
        "recordedAt": datetime.now(timezone.utc).isoformat(),
        "commands": [],
        "checks": [
            {
                "message": "production skill trees match their pre-run hashes",
                "passed": restored,
                "actual": {
                    "beforeFiles": len(before_source),
                    "afterFiles": len(after_source),
                    "changed": sorted(
                        set(before_source) ^ set(after_source)
                        | {path for path in before_source.keys() & after_source.keys() if before_source[path] != after_source[path]}
                    ),
                },
            }
        ],
        "notes": ["Disposable target repositories were removed by TemporaryDirectory cleanup."],
        "error": None if restored else "production source changed during fixture execution",
    }
    (evidence_dir / "qa-18-source-repository-restoration.json").write_text(
        json.dumps(restore_payload, indent=2) + "\n", encoding="utf-8"
    )
    results.append({"qaId": "QA-18", "status": restore_payload["status"], "error": restore_payload["error"]})

    summary = {
        "schemaVersion": "setup-context-driven/manual-qa-summary-v1",
        "recordedAt": datetime.now(timezone.utc).isoformat(),
        "results": results,
        "counts": {
            "pass": sum(item["status"] == "pass" for item in results),
            "fail": sum(item["status"] == "fail" for item in results),
        },
    }
    (evidence_dir / "summary.json").write_text(json.dumps(summary, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(summary["counts"], sort_keys=True))
    return 1 if summary["counts"]["fail"] else 0


if __name__ == "__main__":
    raise SystemExit(main())

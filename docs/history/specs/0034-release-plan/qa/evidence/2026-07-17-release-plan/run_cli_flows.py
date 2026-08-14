#!/usr/bin/env python3
"""Exercise the public Release Plan CLI against real local Git repositories."""

from __future__ import annotations

import hashlib
import json
import os
import shutil
import subprocess
import sys
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[6]
EVIDENCE_DIR = Path(__file__).resolve().parent
BIN = Path(os.environ.get("ROUNDFIX_BIN", REPO_ROOT / "bin" / "roundfix"))


class FlowError(AssertionError):
    pass


def main() -> int:
    if not BIN.exists():
        raise FlowError(f"Roundfix binary does not exist: {BIN}")
    repos_dir = EVIDENCE_DIR / "cli-repos"
    outputs_dir = EVIDENCE_DIR / "cli-outputs"
    if repos_dir.exists():
        shutil.rmtree(repos_dir)
    if outputs_dir.exists():
        shutil.rmtree(outputs_dir)
    repos_dir.mkdir(parents=True)
    (repos_dir / "go.mod").write_text("module roundfix-release-plan-qa\n\ngo 1.25\n", encoding="utf-8")
    outputs_dir.mkdir(parents=True)

    results = []
    flow_patch_ready(results, repos_dir, outputs_dir)
    flow_minor_approval_text(results, repos_dir, outputs_dir)
    flow_major_breaking(results, repos_dir, outputs_dir)
    flow_version_zero_breaking(results, repos_dir, outputs_dir)
    flow_mixed_order(results, repos_dir, outputs_dir)
    flow_no_release(results, repos_dir, outputs_dir)
    flow_ambiguous_manual_required(results, repos_dir, outputs_dir)
    flow_manual_classification_does_not_approve(results, repos_dir, outputs_dir)
    flow_invalid_manual_downgrade(results, repos_dir, outputs_dir)
    flow_dirty_tree(results, repos_dir, outputs_dir)
    flow_pre_release_base_rejected(results, repos_dir, outputs_dir)
    flow_explicit_endpoints(results, repos_dir, outputs_dir)
    flow_invalid_roundfix_config_ignored(results, repos_dir, outputs_dir)

    output = {
        "schemaVersion": "roundfix-release-plan-qa/v1",
        "binary": str(BIN),
        "flows": results,
    }
    (EVIDENCE_DIR / "cli-flow-results.json").write_text(
        json.dumps(output, indent=2, sort_keys=True) + "\n",
        encoding="utf-8",
    )
    print(json.dumps({"ok": True, "flows": len(results), "results": "cli-flow-results.json"}))
    return 0


def flow_patch_ready(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "patch-ready", "v0.4.0", [
        commit("fix: correct release output", ["internal/cli/release.txt"]),
    ])
    result = run_plan(repo, outputs_dir, "patch-ready", "--format", "json")
    plan = json_plan(result)
    expect(result["exitCode"] == 0, "patch-ready exit")
    expect(plan["schemaVersion"] == "roundfix.release-plan/v1", "patch-ready schema")
    expect(plan["state"] == "ready", "patch-ready state")
    expect(plan["proposedVersion"] == "v0.4.1", "patch-ready proposed")
    expect(plan["approval"]["required"] is False, "patch-ready approval")
    expect(len(plan["changes"]) == 1, "patch-ready evidence length")
    append_result(results, "patch-ready", repo, result, plan)


def flow_minor_approval_text(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "minor-approval", "v0.4.1", [
        commit("feat: expose release plan", ["internal/cli/release.txt"]),
    ])
    result = run_plan(repo, outputs_dir, "minor-approval-text")
    expect(result["exitCode"] == 3, "minor exit")
    expect(result["stdout"].startswith("Decision: approval_required\n"), "minor text starts with decision")
    expect("Proposed version: v0.5.0" in result["stdout"], "minor proposed text")
    expect("Approval question: Approve the minor increment to v0.5.0?" in result["stdout"], "minor question")
    expect("stderr" not in result["stdout"].lower(), "minor stdout diagnostic leak")
    expect(result["stderr"] == "", "minor stderr empty")
    append_result(results, "minor-approval-text", repo, result)


def flow_major_breaking(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "major-breaking", "v1.4.2", [
        commit("feat!: remove legacy release command", ["internal/cli/release.txt"]),
    ])
    result = run_plan(repo, outputs_dir, "major-breaking", "--format", "json")
    plan = json_plan(result)
    expect(result["exitCode"] == 3, "major exit")
    expect(plan["state"] == "approval_required", "major state")
    expect(plan["proposedVersion"] == "v2.0.0", "major proposed")
    expect(plan["classification"]["breaking"] is True, "major breaking")
    expect(plan["approval"]["question"] == "Approve the major increment to v2.0.0?", "major question")
    append_result(results, "major-breaking", repo, result, plan)


def flow_version_zero_breaking(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "version-zero-breaking", "v0.7.3", [
        commit("fix!: change agent contract", ["internal/cli/release.txt"]),
    ])
    result = run_plan(repo, outputs_dir, "version-zero-breaking", "--format", "json")
    plan = json_plan(result)
    expect(result["exitCode"] == 3, "version-zero exit")
    expect(plan["state"] == "approval_required", "version-zero state")
    expect(plan["proposedVersion"] == "v0.8.0", "version-zero proposed")
    expect(plan["classification"]["breaking"] is True, "version-zero breaking")
    expect(plan["approval"]["increment"] == "minor", "version-zero increment")
    expect(plan["approval"]["question"] == "Approve the minor increment to v0.8.0?", "version-zero question")
    append_result(results, "version-zero-breaking", repo, result, plan)


def flow_mixed_order(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    scenarios = [
        ("mixed-fix-feature", [
            commit("fix: correct release output", ["internal/cli/fix.txt"]),
            commit("feat: expose release plan", ["internal/cli/feature.txt"]),
        ]),
        ("mixed-feature-fix", [
            commit("feat: expose release plan", ["internal/cli/feature.txt"]),
            commit("fix: correct release output", ["internal/cli/fix.txt"]),
        ]),
    ]
    for name, commits in scenarios:
        repo = create_repo(repos_dir / name, "v0.4.1", commits)
        result = run_plan(repo, outputs_dir, name, "--format", "json")
        plan = json_plan(result)
        expect(result["exitCode"] == 3, f"{name} exit")
        expect(plan["state"] == "approval_required", f"{name} state")
        expect(plan["classification"]["impact"] == "minor", f"{name} impact")
        expect(plan["proposedVersion"] == "v0.5.0", f"{name} proposed")
        expect(len(plan["changes"]) == 2, f"{name} evidence length")
        append_result(results, name, repo, result, plan)


def flow_no_release(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "no-release", "v0.4.1", [
        commit("docs: record release plan evidence", ["docs/specs/0034-release-plan/task_04.md"]),
        commit("test: add release fixture", ["internal/releaseplan/testdata/plan.json"]),
        commit("ci: verify release plan", [".github/workflows/ci-conventions.yml"]),
    ])
    result = run_plan(repo, outputs_dir, "no-release", "--format", "json")
    plan = json_plan(result)
    expect(result["exitCode"] == 0, "no-release exit")
    expect(plan["state"] == "no_release", "no-release state")
    expect("proposedVersion" not in plan, "no-release no proposed version")
    expect(plan["classification"]["impact"] == "none", "no-release impact")
    expect(len(plan["changes"]) == 3, "no-release evidence length")
    append_result(results, "no-release", repo, result, plan)


def flow_ambiguous_manual_required(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "manual-required", "v0.4.1", [
        commit("chore: tune release command", ["internal/cli/release.txt"]),
    ])
    result = run_plan(repo, outputs_dir, "manual-required")
    expect(result["exitCode"] == 3, "manual-required exit")
    expect(result["stdout"].startswith("Decision: manual_classification_required\n"), "manual-required state")
    expect("Proposed version: none" in result["stdout"], "manual-required no proposed")
    expect("Blocking commits:" in result["stdout"], "manual-required blocking")
    expect("--impact <none|patch|minor|major> --reason <text>" in result["stdout"], "manual-required rerun")
    expect(result["stderr"] == "", "manual-required stderr empty")
    append_result(results, "manual-required", repo, result)


def flow_manual_classification_does_not_approve(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "manual-minor", "v0.4.1", [
        commit("chore: tune release command", ["internal/cli/release.txt"]),
    ])
    result = run_plan(
        repo,
        outputs_dir,
        "manual-minor",
        "--format",
        "json",
        "--impact",
        "minor",
        "--reason",
        "public release command behavior changed",
    )
    plan = json_plan(result)
    expect(result["exitCode"] == 3, "manual-minor exit")
    expect(plan["state"] == "approval_required", "manual-minor state")
    expect(plan["classification"]["manualReason"] == "public release command behavior changed", "manual-minor reason")
    expect(plan["approval"]["required"] is True, "manual-minor approval remains required")
    expect(plan["approval"]["question"] == "Approve the minor increment to v0.5.0?", "manual-minor question")
    append_result(results, "manual-minor", repo, result, plan)


def flow_invalid_manual_downgrade(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "manual-downgrade", "v0.4.1", [
        commit("feat: expose release plan", ["internal/cli/feature.txt"]),
    ])
    result = run_plan(
        repo,
        outputs_dir,
        "manual-downgrade",
        "--impact",
        "patch",
        "--reason",
        "try to lower feature impact",
    )
    expect(result["exitCode"] == 2, "manual-downgrade exit")
    expect(result["stdout"] == "", "manual-downgrade no stdout")
    expect("at least the automatic minimum" in result["stderr"], "manual-downgrade guidance")
    append_result(results, "manual-downgrade", repo, result)


def flow_dirty_tree(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "dirty-tree", "v0.4.1", [
        commit("fix: correct release output", ["internal/cli/fix.txt"]),
    ])
    write_file(repo / "internal/cli/fix.txt", "dirty tracked change\n")
    write_file(repo / "scratch.txt", "untracked change\n")
    result = run_plan(repo, outputs_dir, "dirty-tree", "--format", "json")
    expect(result["exitCode"] == 2, "dirty-tree exit")
    expect(result["stdout"] == "", "dirty-tree no stdout")
    expect("internal/cli/fix.txt" in result["stderr"], "dirty-tree tracked path")
    expect("scratch.txt" in result["stderr"], "dirty-tree untracked path")
    expect("commit, stash, or remove" in result["stderr"], "dirty-tree guidance")
    append_result(results, "dirty-tree", repo, result)


def flow_pre_release_base_rejected(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "pre-release-base", "v0.2.0-rc.1", [
        commit("fix: correct release output", ["internal/cli/fix.txt"]),
    ])
    result = run_plan(repo, outputs_dir, "pre-release-base", "--from", "v0.2.0-rc.1")
    expect(result["exitCode"] == 2, "pre-release exit")
    expect(result["stdout"] == "", "pre-release no stdout")
    expect("pre-release tags are not supported" in result["stderr"], "pre-release guidance")
    append_result(results, "pre-release-base", repo, result)


def flow_explicit_endpoints(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(repos_dir / "explicit-endpoints", "v0.1.0", [
        commit("feat: expose release plan", ["internal/cli/feature.txt"]),
        commit("fix: correct release output", ["internal/cli/fix.txt"]),
    ])
    feature_sha = git(repo, "rev-list", "--reverse", "v0.1.0..HEAD").splitlines()[0]
    result = run_plan(repo, outputs_dir, "explicit-endpoints", "--format", "json", "--from", "v0.1.0", "--to", feature_sha)
    plan = json_plan(result)
    expect(result["exitCode"] == 3, "explicit endpoint exit")
    expect(plan["target"]["name"] == feature_sha, "explicit endpoint target")
    expect(len(plan["changes"]) == 1, "explicit endpoint one commit")
    expect(plan["changes"][0]["subject"] == "feat: expose release plan", "explicit endpoint subject")
    append_result(results, "explicit-endpoints", repo, result, plan)


def flow_invalid_roundfix_config_ignored(results: list[dict], repos_dir: Path, outputs_dir: Path) -> None:
    repo = create_repo(
        repos_dir / "invalid-roundfix-config",
        "v0.4.1",
        [
            commit("fix: correct release output", ["internal/cli/fix.txt"]),
        ],
        seed_files={
            "CHANGELOG.md": "# Changelog\n\n## [0.4.1]\n",
            "package.json": '{"name":"roundfix-qa","version":"0.4.1"}\n',
            ".roundfixrc.yml": "unknown_key: should-not-be-read\n",
        },
    )
    result = run_plan(repo, outputs_dir, "invalid-roundfix-config", "--format", "json")
    plan = json_plan(result)
    expect(result["exitCode"] == 0, "invalid config ignored exit")
    expect(plan["state"] == "ready", "invalid config ignored state")
    expect(not (repo / ".roundfix").exists(), "no Run Database or Roundfix state directory")
    append_result(results, "invalid-roundfix-config", repo, result, plan)


def create_repo(path: Path, base_tag: str, commits: list[dict], seed_files: dict[str, str] | None = None) -> Path:
    path.mkdir(parents=True)
    git(path, "init", "--initial-branch=main")
    git(path, "config", "user.name", "Roundfix QA")
    git(path, "config", "user.email", "roundfix-qa@example.invalid")
    git(path, "config", "commit.gpgsign", "false")
    seed = {"seed.txt": "seed\n"}
    if seed_files:
        seed.update(seed_files)
    for relative, content in seed.items():
        write_file(path / relative, content)
    git(path, "add", "-A")
    git(path, "commit", "-m", "chore: seed")
    git(path, "tag", base_tag)
    git(path, "remote", "add", "origin", "https://example.invalid/roundfix.git")
    for item in commits:
        for relative in item["paths"]:
            write_file(path / relative, item["subject"] + "\n")
        git(path, "add", "-A")
        args = ["commit", "-m", item["subject"]]
        if item.get("body"):
            args += ["-m", item["body"]]
        git(path, *args)
    return path


def commit(subject: str, paths: list[str], body: str = "") -> dict:
    return {"subject": subject, "paths": paths, "body": body}


def run_plan(repo: Path, outputs_dir: Path, name: str, *args: str) -> dict:
    before = snapshot(repo)
    completed = subprocess.run(
        [str(BIN), "release", "plan", *args],
        cwd=repo,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )
    after = snapshot(repo)
    stdout_path = outputs_dir / f"{name}.stdout.txt"
    stderr_path = outputs_dir / f"{name}.stderr.txt"
    stdout_path.write_text(completed.stdout, encoding="utf-8")
    stderr_path.write_text(completed.stderr, encoding="utf-8")
    return {
        "command": " ".join(["roundfix", "release", "plan", *args]),
        "exitCode": completed.returncode,
        "stdout": completed.stdout,
        "stderr": completed.stderr,
        "stdoutPath": stdout_path.relative_to(EVIDENCE_DIR).as_posix(),
        "stderrPath": stderr_path.relative_to(EVIDENCE_DIR).as_posix(),
        "repoUnchanged": before == after,
        "beforeSnapshot": before,
        "afterSnapshot": after,
    }


def append_result(results: list[dict], name: str, repo: Path, result: dict, plan: dict | None = None) -> None:
    expect(result["repoUnchanged"], f"{name} repository mutated")
    entry = {
        "name": name,
        "repo": repo.relative_to(EVIDENCE_DIR).as_posix(),
        "command": result["command"],
        "exitCode": result["exitCode"],
        "stdoutPath": result["stdoutPath"],
        "stderrPath": result["stderrPath"],
        "repoUnchanged": result["repoUnchanged"],
        "fileHash": result["afterSnapshot"]["files"],
        "refsHash": sha256(result["afterSnapshot"]["refs"]),
        "statusHash": sha256(result["afterSnapshot"]["status"]),
        "configHash": sha256(result["afterSnapshot"]["config"]),
        "remotesHash": sha256(result["afterSnapshot"]["remotes"]),
    }
    if plan is not None:
        entry["plan"] = plan
    results.append(entry)


def json_plan(result: dict) -> dict:
    expect(result["stderr"] == "", "json plan stderr empty")
    return json.loads(result["stdout"])


def snapshot(repo: Path) -> dict:
    files = {}
    for path in sorted(repo.rglob("*")):
        if ".git" in path.relative_to(repo).parts:
            continue
        if path.is_file():
            relative = path.relative_to(repo).as_posix()
            files[relative] = sha256(path.read_text(encoding="utf-8"))
    return {
        "files": files,
        "refs": git(repo, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads", "refs/tags"),
        "status": git(repo, "--no-optional-locks", "status", "--porcelain=v1", "-z"),
        "config": git(repo, "config", "--local", "--list"),
        "remotes": git(repo, "remote", "-v"),
    }


def write_file(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")


def git(repo: Path, *args: str) -> str:
    completed = subprocess.run(
        ["git", *args],
        cwd=repo,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )
    if completed.returncode != 0:
        raise FlowError(f"git {' '.join(args)} failed in {repo}: {completed.stderr}")
    return completed.stdout


def sha256(value: str) -> str:
    return hashlib.sha256(value.encode("utf-8")).hexdigest()


def expect(condition: bool, label: str) -> None:
    if not condition:
        raise FlowError(label)


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:  # noqa: BLE001 - QA evidence must print the failing condition.
        print(json.dumps({"ok": False, "error": str(exc)}), file=sys.stderr)
        raise

import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))

from context_assets import (  # noqa: E402
    clone_assets_to,
    load_asset_catalog,
    portable_file_digest,
    read_json_copy,
    setup_records_digest,
    write_json,
)
from context_setup import (  # noqa: E402
    audit_repository,
    expected_artifacts_for_plan,
    managed_block,
    plan_apply,
    resolve_decision_plan,
)


class AuditCliTests(unittest.TestCase):
    def test_audit_help_describes_implemented_optional_checks(self):
        result = run_context_setup("audit", "--help")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stderr, "")
        self.assertNotIn("Reserved", result.stdout)
        self.assertIn("Show informational findings", result.stdout)
        self.assertIn("Compare the bundled setup snapshot", result.stdout)

    def test_compliant_repository_returns_zero_in_text_and_json_modes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")

            text = run_audit(repo, "--format", "text")
            self.assertEqual(text.returncode, 0, text.stderr)
            self.assertEqual(text.stderr, "")
            self.assertEqual(text.stdout.strip(), "setup-context-driven audit: ok")

            json_result = run_audit(repo, "--format", "json")
            self.assertEqual(json_result.returncode, 0, json_result.stderr)
            self.assertEqual(json_result.stderr, "")
            payload = json.loads(json_result.stdout)
            self.assertEqual(payload["schemaVersion"], "setup-context-driven/audit-v1")
            self.assertEqual(payload["ok"], True)
            self.assertEqual(payload["summary"], {"errors": 0, "decisions": 0, "warnings": 0, "info": 0})
            self.assertEqual(payload["findings"], [])

    def test_default_subcommand_is_read_only_audit(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "go-cli-tui")
            before = snapshot_files(repo)

            first = run_context_setup("--repo", str(repo), "--format", "json")
            second = run_context_setup("audit", "--repo", str(repo), "--format", "json")

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(json.loads(first.stdout), json.loads(second.stdout))
            self.assertEqual(snapshot_files(repo), before)

    def test_blocking_malformed_and_decision_exit_codes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            missing_manifest_repo = Path(temp_dir) / "missing"
            missing_manifest_repo.mkdir()
            blocking = run_audit(missing_manifest_repo, "--format", "json")
            self.assertEqual(blocking.returncode, 1)
            self.assertFinding(blocking, "manifest.missing", "error")

            malformed_repo = Path(temp_dir) / "malformed"
            malformed_repo.mkdir()
            manifest_path = malformed_repo / "docs" / "agents" / "setup-context.json"
            manifest_path.parent.mkdir(parents=True)
            manifest_path.write_text("{not json", encoding="utf-8")
            malformed = run_audit(malformed_repo, "--format", "json")
            self.assertEqual(malformed.returncode, 2)
            self.assertFinding(malformed, "manifest.invalid", "error")

            decision_repo = Path(temp_dir) / "decision"
            write_compliant_repository(decision_repo, "rust-cli", omit_decision="runtime.design")
            decision = run_audit(decision_repo, "--format", "json")
            self.assertEqual(decision.returncode, 3)
            self.assertFinding(decision, "decision.required", "decision")

    def test_json_finding_shape_contains_stable_fields(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            (repo / "docs" / "agents" / "rust.md").unlink()

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 1)
            payload = json.loads(result.stdout)
            finding = payload["findings"][0]
            self.assertEqual(
                set(finding),
                {"code", "severity", "path", "managedId", "message", "action"},
            )
            self.assertEqual(finding["code"], "docs.guide.missing")
            self.assertEqual(finding["path"], "docs/agents/rust.md")
            self.assertEqual(finding["managedId"], "guide.rust")
            self.assertTrue(finding["message"])
            self.assertTrue(finding["action"])

    def test_marker_document_and_template_findings_are_independent(self):
        cases = [
            ("marker", corrupt_marker, "managed.marker.invalid"),
            ("missing guide", remove_rust_guide, "docs.guide.missing"),
            ("broken reference", add_broken_reference, "docs.reference.broken"),
            ("stale template", make_stale_template_version, "managed.template.stale"),
            ("non-English", add_non_english_content, "docs.language.non-english"),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    write_compliant_repository(repo, "rust-cli")
                    mutator(repo)

                    result = run_audit(repo, "--format", "json")

                    self.assertIn(result.returncode, {0, 1})
                    self.assertFinding(result, expected_code)

    def test_audit_twice_produces_identical_semantic_output(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            add_broken_reference(repo)

            first = run_audit(repo, "--format", "json")
            second = run_audit(repo, "--format", "json")

            self.assertEqual(first.returncode, second.returncode)
            self.assertEqual(json.loads(first.stdout), json.loads(second.stdout))

    def test_stale_managed_target_cannot_satisfy_excluded_decision_plan_reference(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            repo = root / "repo"
            write_compliant_repository(repo, "rust-cli")
            stale_target = repo / "docs" / "agents" / "domain.md"
            self.assertTrue(stale_target.is_file())

            temp_skill = root / "setup-context-driven"
            clone_assets_to(SKILL_ROOT, temp_skill)
            exclude_domain_reference_target(temp_skill)
            catalog = load_asset_catalog(temp_skill)

            result, invalid_input = audit_repository(repo, catalog)

            self.assertFalse(invalid_input)
            matches = [
                finding
                for finding in result.findings
                if finding.code == "reference.managed.missing"
            ]
            self.assertEqual(len(matches), 1, result.findings)
            self.assertEqual(matches[0].path, "AGENTS.md")
            self.assertEqual(matches[0].managed_id, "root.context-workflow")
            self.assertTrue(stale_target.is_file())

    def test_frontend_missing_repository_design_contract_is_one_read_only_finding(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "typescript-bun-monorepo")
            design_path = repo / "DESIGN.md"
            design_path.unlink()
            before = snapshot_files(repo)

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 1, result.stderr)
            payload = json.loads(result.stdout)
            matches = [
                finding
                for finding in payload["findings"]
                if finding["code"] == "reference.repository.missing"
            ]
            self.assertEqual(len(matches), 1, payload)
            self.assertEqual(matches[0]["path"], "DESIGN.md")
            self.assertEqual(matches[0]["managedId"], "guide.frontend")
            self.assertEqual(snapshot_files(repo), before)
            self.assertFalse(design_path.exists())

    def test_repository_reference_symlink_outside_repo_is_blocked(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            repo = root / "repo"
            write_compliant_repository(repo, "typescript-bun-monorepo")
            design_path = repo / "DESIGN.md"
            design_path.unlink()
            outside = root / "outside-design.md"
            outside.write_text("# Outside\n", encoding="utf-8")
            design_path.symlink_to(outside)

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 1, result.stderr)
            payload = json.loads(result.stdout)
            matches = [
                finding
                for finding in payload["findings"]
                if finding["code"] == "reference.repository.outside"
            ]
            self.assertEqual(len(matches), 1, payload)
            self.assertEqual(matches[0]["path"], "DESIGN.md")
            self.assertEqual(outside.read_text(encoding="utf-8"), "# Outside\n")

    def test_apply_plan_blocks_before_creating_frontend_pointer_to_missing_design(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)
            catalog = load_asset_catalog(SKILL_ROOT)

            result, invalid_input, change_plan = plan_apply(
                repo=repo,
                catalog=catalog,
                profile_override="typescript-bun-monorepo",
                decision_args=[
                    "spec.scaffold=true",
                    "domain.layout=single-context",
                    "triage.external=false",
                    "autonomous.enabled=false",
                    "verification.gate=make verify",
                    "language.generated=English",
                    "secondbrain.enabled=false",
                ],
            )

            self.assertFalse(invalid_input)
            self.assertEqual(
                [finding.code for finding in result.findings],
                ["reference.repository.missing"],
            )
            self.assertEqual(change_plan.changes, [])
            self.assertEqual(snapshot_files(repo), before)

    def assertFinding(self, result, code, severity=None):
        payload = json.loads(result.stdout)
        matches = [finding for finding in payload["findings"] if finding["code"] == code]
        self.assertGreater(len(matches), 0, payload)
        if severity is not None:
            self.assertEqual(matches[0]["severity"], severity)


def run_audit(repo, *args, extra_env=None):
    return run_context_setup("audit", "--repo", str(repo), *args, extra_env=extra_env)


def run_context_setup(*args, extra_env=None):
    env = os.environ.copy()
    env["PYTHONDONTWRITEBYTECODE"] = "1"
    if extra_env is not None:
        env.update(extra_env)
    with tempfile.TemporaryDirectory() as temp_dir:
        fixture_root = Path(temp_dir) / "setup-context-driven"
        clone_assets_to(SKILL_ROOT, fixture_root)
        shutil.copytree(SKILL_ROOT / "scripts", fixture_root / "scripts")
        write_fixture_external_digests(fixture_root)
        return subprocess.run(
            [sys.executable, str(fixture_root / "scripts" / "context_setup.py"), *args],
            text=True,
            capture_output=True,
            check=False,
            env=env,
        )


def write_compliant_repository(repo, profile_id, omit_decision=None, install_skills=True):
    repo.mkdir(parents=True, exist_ok=True)
    catalog = load_asset_catalog(SKILL_ROOT)
    decisions = {
        "spec.scaffold": {"value": True, "confirmedAt": "2026-07-15"},
        "domain.layout": {"value": "single-context", "confirmedAt": "2026-07-15"},
        "triage.external": {"value": False, "confirmedAt": "2026-07-15"},
        "autonomous.enabled": {"value": True, "confirmedAt": "2026-07-15"},
        "runtime.backend": {"value": "codex gpt-5.5 xhigh", "confirmedAt": "2026-07-15"},
        "runtime.design": {"value": "claude opus xhigh", "confirmedAt": "2026-07-15"},
        "verification.gate": {"value": "make verify", "confirmedAt": "2026-07-15"},
        "language.generated": {"value": "English", "confirmedAt": "2026-07-15"},
        "secondbrain.enabled": {"value": False, "confirmedAt": "2026-07-15"},
    }
    plan = resolve_decision_plan(
        catalog,
        profile_id,
        {"decisions": decisions},
        {},
    )
    modules = list(plan.active_modules)
    artifacts = expected_artifacts_for_plan(plan)

    grouped = {}
    for artifact in artifacts:
        grouped.setdefault(artifact.path, []).append(artifact)

    for path, path_artifacts in grouped.items():
        target = repo / path
        target.parent.mkdir(parents=True, exist_ok=True)
        content = "".join(
            managed_block(artifact.managed_id, artifact.version, artifact.content)
            for artifact in path_artifacts
        )
        target.write_text(content, encoding="utf-8")

    for artifact in artifacts:
        for reference in catalog.references_by_artifact.get(artifact.managed_id, ()):
            if reference.ownership != "repository":
                continue
            target = repo / reference.repository_path
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_text("# Repository-owned contract\n", encoding="utf-8")

    if omit_decision is not None:
        decisions.pop(omit_decision)

    manifest = {
        "schemaVersion": 1,
        "generator": {"skill": "setup-context-driven", "version": 1},
        "profile": profile_id,
        "modules": modules,
        "decisions": decisions,
        "managedArtifacts": [
            {
                "id": artifact.managed_id,
                "path": artifact.path.as_posix(),
                "kind": artifact.kind,
                "module": artifact.module_id,
                "template": artifact.template_id,
                "version": artifact.version,
                "digest": artifact.digest,
            }
            for artifact in artifacts
        ],
        "localSkills": [],
    }
    manifest_path = repo / "docs" / "agents" / "setup-context.json"
    manifest_path.parent.mkdir(parents=True, exist_ok=True)
    manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
    if install_skills:
        install_profile_skills(repo, profile_id)


def install_profile_skills(repo, profile_id, omit=None):
    omitted = set(omit or [])
    catalog = load_asset_catalog(SKILL_ROOT)
    setup_id = catalog.profiles[profile_id]["setup"]
    for skill in catalog.setups[setup_id]["skills"]:
        name = skill["name"]
        if name in omitted:
            continue
        write_skill(repo, name)


def write_skill(repo, name):
    skill_path = repo / ".agents" / "skills" / name / "SKILL.md"
    skill_path.parent.mkdir(parents=True, exist_ok=True)
    repo_root = repository_root()
    canonical = repo_root / ".agents" / "skills" / name / "SKILL.md"
    if fixture_skill_source_type(name) == "repo" and canonical.is_file():
        skill_path.write_text(canonical.read_text(encoding="utf-8"), encoding="utf-8")
    else:
        for relative_path, content in fixture_external_skill_files(name).items():
            target = skill_path.parent / relative_path
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(content)


def fixture_external_skill_bytes(name):
    return f"---\nname: {name}\ndescription: test skill\n---\n# {name}\n".encode("utf-8")


def fixture_external_skill_files(name):
    files = {"SKILL.md": fixture_external_skill_bytes(name)}
    if name == "domain-modeling":
        files["references/guide.md"] = b"# Domain modeling fixture\n"
    return files


def fixture_skill_source_type(name):
    for snapshot_path in (SKILL_ROOT / "assets" / "setups").glob("*.json"):
        snapshot = json.loads(snapshot_path.read_text(encoding="utf-8"))
        for skill in snapshot["skills"]:
            if skill["name"] == name:
                return skill["source"]["type"]
    return None


def write_fixture_external_digests(skill_root):
    for snapshot_path in (skill_root / "assets" / "setups").glob("*.json"):
        snapshot = json.loads(snapshot_path.read_text(encoding="utf-8"))
        for skill in snapshot["skills"]:
            if skill["source"]["type"] != "github":
                continue
            skill["treeDigest"] = portable_file_digest(
                [
                    (path.encode("utf-8"), content)
                    for path, content in fixture_external_skill_files(skill["name"]).items()
                ]
            )
        snapshot["digest"] = setup_records_digest(snapshot["skills"])
        snapshot_path.write_text(json.dumps(snapshot, indent=2) + "\n", encoding="utf-8")


def snapshot_files(repo):
    return {
        path.relative_to(repo).as_posix(): path.read_bytes()
        for path in sorted(repo.rglob("*"))
        if path.is_file()
    }


def repository_root():
    for parent in [SKILL_ROOT, *SKILL_ROOT.parents]:
        if (parent / ".agents" / "skills" / "setup-context-driven").is_dir() and (
            parent / "skills" / "setup-context-driven"
        ).is_dir():
            return parent
    raise AssertionError("could not locate repository root")


def corrupt_marker(repo):
    path = repo / "AGENTS.md"
    path.write_text(
        path.read_text(encoding="utf-8").replace(
            "<!-- setup-context-driven:end id=root.rust -->",
            "<!-- setup-context-driven:end id=root.mismatch -->",
        ),
        encoding="utf-8",
    )


def remove_rust_guide(repo):
    (repo / "docs" / "agents" / "rust.md").unlink()


def add_broken_reference(repo):
    path = repo / "docs" / "agents" / "rust.md"
    path.write_text(
        path.read_text(encoding="utf-8").replace(
            "# Rust",
            "# Rust\n\nSee [missing guide](missing-guide.md).",
        ),
        encoding="utf-8",
    )


def make_stale_template_version(repo):
    path = repo / "docs" / "agents" / "rust.md"
    path.write_text(
        path.read_text(encoding="utf-8").replace(
            "id=guide.rust version=2",
            "id=guide.rust version=1",
        ),
        encoding="utf-8",
    )


def add_non_english_content(repo):
    path = repo / "docs" / "agents" / "rust.md"
    path.write_text(
        path.read_text(encoding="utf-8").replace(
            "# Rust",
            "# Rust\n\nEste arquivo não está em inglês.",
        ),
        encoding="utf-8",
    )


def exclude_domain_reference_target(skill_root):
    decisions_path = skill_root / "assets" / "decisions.json"
    decisions = read_json_copy(decisions_path)
    for decision in decisions["decisions"]:
        if decision["id"] != "domain.layout":
            continue
        effect = decision["effects"][0]
        effect["includeArtifacts"] = []
        effect["excludeArtifacts"] = ["guide.domain"]
        effect["selectTemplates"] = []
        write_json(decisions_path, decisions)
        return
    raise AssertionError("missing domain.layout decision")


if __name__ == "__main__":
    unittest.main()

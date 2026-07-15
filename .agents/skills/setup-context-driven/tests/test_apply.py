import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_setup import parse_managed_blocks  # noqa: E402
from test_audit import snapshot_files, write_compliant_repository  # noqa: E402


BASE_DECISIONS = [
    "spec.scaffold=true",
    "domain.layout=single-context",
    "triage.external=false",
    "autonomous.enabled=true",
    "runtime.backend=codex gpt-5.5 xhigh",
    "runtime.design=claude opus xhigh",
    "verification.gate=make verify",
    "language.generated=English",
]


class ApplyCliTests(unittest.TestCase):
    def test_apply_creates_manifest_root_blocks_and_guides(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            result = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 0, result.stderr)
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
            )
            self.assertEqual(manifest["profile"], "rust-cli")
            self.assertEqual(
                manifest["modules"],
                ["core", "context-workflow", "rust", "cli-surface", "autonomous-work"],
            )
            self.assertIn("runtime.design", manifest["decisions"])
            self.assertTrue(manifest["managedArtifacts"])
            self.assertIn("<!-- setup-context-driven:begin id=root.rust version=1 -->", (repo / "AGENTS.md").read_text(encoding="utf-8"))
            self.assertTrue((repo / "docs" / "agents" / "rust.md").is_file())
            self.assertEqual(run_audit(repo).returncode, 0)

    def test_apply_preserves_custom_content_and_is_idempotent(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            (repo / "AGENTS.md").write_text(
                "custom before\n\ncustom after\n",
                encoding="utf-8",
            )

            first = run_apply(repo, "rust-cli", BASE_DECISIONS)
            after_first = snapshot_files(repo)
            second = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(snapshot_files(repo), after_first)
            content = (repo / "AGENTS.md").read_text(encoding="utf-8")
            self.assertIn("custom before\n\ncustom after\n", content)

    def test_missing_decision_prevents_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)
            decisions = [item for item in BASE_DECISIONS if not item.startswith("runtime.design=")]

            result = run_apply(repo, "rust-cli", decisions)

            self.assertEqual(result.returncode, 3)
            self.assertFinding(result, "decision.required", "decision")
            self.assertEqual(snapshot_files(repo), before)

    def test_invalid_duplicate_markers_prevent_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            duplicated = (
                "<!-- setup-context-driven:begin id=root.core version=1 -->\n"
                "first\n"
                "<!-- setup-context-driven:end id=root.core -->\n"
                "<!-- setup-context-driven:begin id=root.core version=1 -->\n"
                "second\n"
                "<!-- setup-context-driven:end id=root.core -->\n"
            )
            (repo / "AGENTS.md").write_text(duplicated, encoding="utf-8")
            before = snapshot_files(repo)

            result = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "managed.block.duplicate", "error")
            self.assertEqual(snapshot_files(repo), before)

    def test_failure_before_commit_preserves_target_files(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            docs_agents = repo / "docs" / "agents"
            docs_agents.mkdir(parents=True)
            (repo / "AGENTS.md").write_text("custom root\n", encoding="utf-8")
            before = snapshot_files(repo)
            original_mode = docs_agents.stat().st_mode
            docs_agents.chmod(0o500)
            try:
                result = run_apply(repo, "rust-cli", BASE_DECISIONS)
            finally:
                docs_agents.chmod(original_mode)

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "managed.apply.failed", "error")
            self.assertEqual(snapshot_files(repo), before)

    def test_obsolete_managed_artifacts_are_removed_without_unowned_deletions(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "go-cli-tui")
            go_guide = repo / "docs" / "agents" / "go.md"
            go_guide.write_text(
                "custom go intro\n" + go_guide.read_text(encoding="utf-8") + "custom go tail\n",
                encoding="utf-8",
            )
            (repo / "docs" / "agents" / "local.md").write_text("local guide\n", encoding="utf-8")

            result = run_apply(repo, "rust-cli", [])

            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertFalse((repo / "docs" / "agents" / "tui.md").exists())
            self.assertTrue((repo / "docs" / "agents" / "local.md").exists())
            self.assertEqual(go_guide.read_text(encoding="utf-8"), "custom go intro\ncustom go tail\n")
            root_blocks, _ = parse_managed_blocks(
                Path("AGENTS.md"),
                (repo / "AGENTS.md").read_text(encoding="utf-8"),
            )
            self.assertIn("root.rust", root_blocks)
            self.assertNotIn("root.go", root_blocks)

    def test_adoption_decision_required_before_owning_unmarked_guide(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            rust_guide = repo / "docs" / "agents" / "rust.md"
            rust_guide.parent.mkdir(parents=True)
            rust_guide.write_text("repo-authored rust notes\n", encoding="utf-8")
            before = snapshot_files(repo)

            blocked = run_apply(repo, "rust-cli", BASE_DECISIONS)
            self.assertEqual(blocked.returncode, 3)
            self.assertFinding(blocked, "decision.required", "decision")
            self.assertEqual(snapshot_files(repo), before)

            adopted = run_apply(
                repo,
                "rust-cli",
                BASE_DECISIONS + ["adoption.guide.rust=true"],
            )

            self.assertEqual(adopted.returncode, 0, adopted.stderr)
            self.assertIn(
                "<!-- setup-context-driven:begin id=guide.rust version=1 -->",
                rust_guide.read_text(encoding="utf-8"),
            )
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
            )
            self.assertEqual(manifest["decisions"]["adoption.guide.rust"]["value"], True)

    def test_unsupported_manifest_version_blocks_apply_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            manifest_path = repo / "docs" / "agents" / "setup-context.json"
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            manifest["schemaVersion"] = 0
            manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
            before = snapshot_files(repo)

            result = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "manifest.migration-required", "error")
            self.assertEqual(snapshot_files(repo), before)

    def assertFinding(self, result, code, severity):
        payload = json.loads(result.stdout)
        matches = [finding for finding in payload["findings"] if finding["code"] == code]
        self.assertGreater(len(matches), 0, payload)
        self.assertEqual(matches[0]["severity"], severity)


def run_apply(repo, profile, decisions):
    args = ["apply", "--repo", str(repo), "--format", "json", "--profile", profile]
    for decision in decisions:
        args.extend(["--decision", decision])
    return run_context_setup(*args)


def run_audit(repo):
    return run_context_setup("audit", "--repo", str(repo), "--format", "json")


def run_context_setup(*args):
    return subprocess.run(
        [sys.executable, str(SCRIPT), *args],
        text=True,
        capture_output=True,
        check=False,
    )


if __name__ == "__main__":
    unittest.main()

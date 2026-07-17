import json
import sys
import tempfile
import unittest
from pathlib import Path


sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_setup import parse_managed_blocks  # noqa: E402
from test_apply import BASE_DECISIONS, run_apply, run_audit  # noqa: E402
from test_audit import install_profile_skills, snapshot_files  # noqa: E402


class SecondbrainSetupTests(unittest.TestCase):
    def test_secondbrain_opt_in_creates_compact_pointer_and_read_only_guide(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            result = run_apply(repo, "rust-cli", secondbrain_decisions(True))
            install_profile_skills(repo, "rust-cli")

            self.assertEqual(result.returncode, 0, result.stderr)
            manifest = read_manifest(repo)
            self.assertIn("secondbrain", manifest["modules"])
            self.assertEqual(manifest["decisions"]["secondbrain.enabled"]["value"], True)

            root = (repo / "AGENTS.md").read_text(encoding="utf-8")
            guide = (repo / "docs" / "agents" / "secondbrain.md").read_text(encoding="utf-8")
            root_blocks, _ = parse_managed_blocks(Path("AGENTS.md"), root)
            self.assertIn("root.secondbrain", root_blocks)
            self.assertLessEqual(len(root_blocks["root.secondbrain"].body.split()), 45)
            self.assertIn("docs/agents/secondbrain.md", root_blocks["root.secondbrain"].body)
            for phrase in REQUIRED_GUIDE_PHRASES:
                self.assertIn(phrase, guide)
            self.assertEqual(run_audit(repo).returncode, 0)

    def test_secondbrain_opt_out_creates_no_pointer_or_guide(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            result = run_apply(repo, "rust-cli", secondbrain_decisions(False))

            self.assertEqual(result.returncode, 0, result.stderr)
            manifest = read_manifest(repo)
            self.assertNotIn("secondbrain", manifest["modules"])
            self.assertEqual(manifest["decisions"]["secondbrain.enabled"]["value"], False)
            self.assertNotIn("root.secondbrain", (repo / "AGENTS.md").read_text(encoding="utf-8"))
            self.assertFalse((repo / "docs" / "agents" / "secondbrain.md").exists())

    def test_disabling_secondbrain_removes_only_marked_managed_content(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            enabled = run_apply(repo, "rust-cli", secondbrain_decisions(True))
            self.assertEqual(enabled.returncode, 0, enabled.stderr)
            guide_path = repo / "docs" / "agents" / "secondbrain.md"
            guide_path.write_text(
                "repo-owned note\n" + guide_path.read_text(encoding="utf-8"),
                encoding="utf-8",
            )
            before_disable = snapshot_files(repo)

            disabled = run_apply(repo, "rust-cli", ["secondbrain.enabled=false"])

            self.assertEqual(disabled.returncode, 0, disabled.stderr)
            self.assertNotEqual(snapshot_files(repo), before_disable)
            root_blocks, _ = parse_managed_blocks(
                Path("AGENTS.md"),
                (repo / "AGENTS.md").read_text(encoding="utf-8"),
            )
            self.assertNotIn("root.secondbrain", root_blocks)
            self.assertEqual(guide_path.read_text(encoding="utf-8"), "repo-owned note\n")
            manifest = read_manifest(repo)
            self.assertNotIn("secondbrain", manifest["modules"])
            self.assertEqual(manifest["decisions"]["secondbrain.enabled"]["value"], False)

    def test_secondbrain_audit_reports_missing_safety_rule(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            applied = run_apply(repo, "rust-cli", secondbrain_decisions(True))
            install_profile_skills(repo, "rust-cli")
            self.assertEqual(applied.returncode, 0, applied.stderr)
            guide_path = repo / "docs" / "agents" / "secondbrain.md"
            guide_path.write_text(
                guide_path.read_text(encoding="utf-8").replace(
                    "Never read, copy, or expose",
                    "Do not expose",
                ),
                encoding="utf-8",
            )

            audit = run_audit(repo)

            self.assertEqual(audit.returncode, 1)
            payload = json.loads(audit.stdout)
            matches = [
                finding
                for finding in payload["findings"]
                if finding["code"] == "secondbrain.safety-rule.missing"
            ]
            self.assertEqual(len(matches), 1, payload)
            self.assertEqual(matches[0]["managedId"], "guide.secondbrain")

    def test_secondbrain_generated_content_is_english_and_root_is_index_only(self):
        repo_root = Path(__file__).resolve().parents[4]
        root_template = repo_root / ".agents/skills/setup-context-driven/assets/templates/root/secondbrain.md"
        guide_template = repo_root / ".agents/skills/setup-context-driven/assets/templates/guides/secondbrain.md"
        root = root_template.read_text(encoding="utf-8")
        guide = guide_template.read_text(encoding="utf-8")

        self.assertLessEqual(len(root.split()), 45)
        self.assertIn("docs/agents/secondbrain.md", root)
        self.assertNotIn(" não ", f" {guide.lower()} ")
        self.assertNotIn(" repositorio", f" {guide.lower()} ")
        self.assertIn("Do not write to the Secondbrain", guide)
        self.assertIn("Do not edit raw/", guide)
        self.assertIn("Do not edit projects/*/mirror/", guide)


REQUIRED_GUIDE_PHRASES = [
    "wiki/index.md",
    "qmd query",
    "projects/<project>/mirror",
    "Cite every Secondbrain file",
    "Do not write to the Secondbrain",
    "Do not edit raw/",
    "Do not edit projects/*/mirror/",
    "Hermes",
    "Never read, copy, or expose",
]


def secondbrain_decisions(enabled):
    replacement = f"secondbrain.enabled={str(enabled).lower()}"
    return [
        replacement if decision.startswith("secondbrain.enabled=") else decision
        for decision in BASE_DECISIONS
    ]


def read_manifest(repo):
    return json.loads(
        (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
    )


if __name__ == "__main__":
    unittest.main()

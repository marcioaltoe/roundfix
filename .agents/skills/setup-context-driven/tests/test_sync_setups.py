import json
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_assets import clone_assets_to  # noqa: E402
from context_setup import sync_setup_snapshots  # noqa: E402
from test_audit import snapshot_files  # noqa: E402


class SyncSetupsTests(unittest.TestCase):
    def test_check_succeeds_when_canonical_snapshots_match(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            source_dir = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            copy_setup_sources(source_dir)
            before = snapshot_files(skill_root)

            result, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertFalse(invalid_input)
            self.assertTrue(result.ok)
            self.assertEqual(result.findings, [])
            self.assertEqual(snapshot_files(skill_root), before)

    def test_check_reports_drift_when_canonical_content_differs(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            source_dir = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            copy_setup_sources(source_dir)
            add_source_skill(source_dir / "rust-cli.json", "drift-check")

            result, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertFalse(invalid_input)
            self.assertFalse(result.ok)
            matches = [finding for finding in result.findings if finding.code == "skills.setup-snapshot.drift"]
            self.assertEqual(len(matches), 1)
            self.assertEqual(matches[0].managed_id, "setup.rust-cli")

    def test_sync_updates_snapshots_atomically_and_then_becomes_idempotent(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            source_dir = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            copy_setup_sources(source_dir)
            add_source_skill(source_dir / "rust-cli.json", "drift-update")

            updated, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=False)
            after_update = snapshot_files(skill_root)
            repeated, repeated_invalid = sync_setup_snapshots(skill_root, source_dir, check=False)

            self.assertFalse(invalid_input)
            self.assertEqual(updated.summary["info"], 1)
            self.assertFalse(repeated_invalid)
            self.assertEqual(repeated.findings, [])
            self.assertEqual(snapshot_files(skill_root), after_update)
            self.assertEqual(list(skill_root.rglob("*.setup-context.tmp")), [])
            snapshot = json.loads(
                (skill_root / "assets" / "setups" / "rust-cli.json").read_text(encoding="utf-8")
            )
            self.assertIn("drift-update", [skill["name"] for skill in snapshot["skills"]])

    def test_text_setup_sources_normalize_comments_and_blank_lines(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            source_dir = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            write_text_setup_sources(source_dir)

            result, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertFalse(invalid_input)
            self.assertTrue(result.ok)
            self.assertEqual(result.findings, [])

    def test_text_setup_sources_accept_canonical_repository_paths(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            source_dir = temp_root / "canonical" / "setups"
            clone_assets_to(SKILL_ROOT, skill_root)
            write_canonical_repository_setup_sources(source_dir)

            result, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertFalse(invalid_input)
            self.assertFalse(result.ok)
            self.assertEqual(
                [finding.managed_id for finding in result.findings],
                ["setup.go-cli", "setup.rust-cli", "setup.typescript-bun"],
            )

    def test_sync_preserves_canonical_path_and_hashes_skill_file(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            source_dir = temp_root / "canonical" / "setups"
            clone_assets_to(SKILL_ROOT, skill_root)
            copy_setup_sources(source_dir)
            source_path = source_dir / "rust-cli.json"
            source = json.loads(source_path.read_text(encoding="utf-8"))
            canonical_path = "skills/test/canonical-example"
            source["skills"].append({"path": canonical_path})
            source_path.write_text(json.dumps(source, indent=2) + "\n", encoding="utf-8")
            skill_file = source_dir.parent / canonical_path / "SKILL.md"
            skill_file.parent.mkdir(parents=True)
            skill_file.write_text("# Canonical example\n", encoding="utf-8")

            result, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=False)

            self.assertFalse(invalid_input)
            self.assertTrue(result.ok)
            snapshot = json.loads(
                (skill_root / "assets" / "setups" / "rust-cli.json").read_text(
                    encoding="utf-8"
                )
            )
            added = next(skill for skill in snapshot["skills"] if skill["name"] == "canonical-example")
            self.assertEqual(added["path"], canonical_path)
            self.assertTrue(added["contentDigest"])


def copy_setup_sources(source_dir):
    source_dir.mkdir(parents=True)
    for source in (SKILL_ROOT / "assets" / "setups").glob("*.json"):
        (source_dir / source.name).write_text(source.read_text(encoding="utf-8"), encoding="utf-8")


def write_text_setup_sources(source_dir):
    source_dir.mkdir(parents=True)
    for source in (SKILL_ROOT / "assets" / "setups").glob("*.json"):
        snapshot = json.loads(source.read_text(encoding="utf-8"))
        lines = ["# canonical setup", ""]
        for skill in snapshot["skills"]:
            lines.append(skill["path"])
            lines.append("")
        (source_dir / f"{snapshot['id']}.txt").write_text("\n".join(lines), encoding="utf-8")


def write_canonical_repository_setup_sources(source_dir):
    source_dir.mkdir(parents=True)
    for source in (SKILL_ROOT / "assets" / "setups").glob("*.json"):
        snapshot = json.loads(source.read_text(encoding="utf-8"))
        lines = ["# canonical setup"]
        lines.extend(f"skills/test/{skill['name']}" for skill in snapshot["skills"])
        (source_dir / f"{snapshot['id']}.txt").write_text("\n".join(lines), encoding="utf-8")


def add_source_skill(path, name):
    snapshot = json.loads(path.read_text(encoding="utf-8"))
    snapshot["skills"].append(
        {
            "name": name,
            "path": f".agents/skills/{name}/SKILL.md",
            "source": {"type": "github", "name": "example/skills"},
            "contentDigest": "1" * 64,
        }
    )
    path.write_text(json.dumps(snapshot, indent=2) + "\n", encoding="utf-8")


if __name__ == "__main__":
    unittest.main()

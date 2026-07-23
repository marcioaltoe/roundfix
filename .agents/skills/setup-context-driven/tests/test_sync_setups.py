"""Suite: portable setup snapshot synchronization.
Invariant: external snapshot bytes match immutable Git provenance and repo-owned bytes retain precedence.
Boundary IN: local Git checkout inspection, complete-tree hashing, and bundled snapshot writes.
Boundary OUT: installed Repository Skill Set audit and external skill restoration.
"""

import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_assets import PortableTreeError, clone_assets_to, portable_tree_digest  # noqa: E402
from context_setup import sync_setup_snapshots  # noqa: E402
from test_audit import snapshot_files  # noqa: E402


class SyncSetupsTests(unittest.TestCase):
    def test_portable_tree_digest_covers_complete_tree_and_excludes_dependency_trees(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            (root / "nested").mkdir()
            (root / "nested" / "guide.md").write_text("guide\n", encoding="utf-8")
            (root / "SKILL.md").write_text("skill\n", encoding="utf-8")
            baseline = portable_tree_digest(root)

            (root / "nested" / "guide.md").write_text("changed\n", encoding="utf-8")
            self.assertNotEqual(portable_tree_digest(root), baseline)
            (root / "nested" / "guide.md").write_text("guide\n", encoding="utf-8")
            (root / "nested" / "added.md").write_text("added\n", encoding="utf-8")
            self.assertNotEqual(portable_tree_digest(root), baseline)
            (root / "nested" / "added.md").unlink()
            (root / "SKILL.md").unlink()
            self.assertNotEqual(portable_tree_digest(root), baseline)

            (root / "SKILL.md").write_text("skill\n", encoding="utf-8")
            (root / "node_modules").mkdir()
            (root / "node_modules" / "ignored.js").write_text("ignored\n", encoding="utf-8")
            (root / ".git").mkdir()
            (root / ".git" / "ignored").write_text("ignored\n", encoding="utf-8")
            self.assertEqual(portable_tree_digest(root), baseline)

    def test_portable_tree_digest_rejects_links_and_special_files(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            target = root / "target.md"
            target.write_text("target\n", encoding="utf-8")
            (root / "link.md").symlink_to(target)
            with self.assertRaises(PortableTreeError):
                portable_tree_digest(root)

        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            target = root / "target.md"
            target.write_text("target\n", encoding="utf-8")
            os.link(target, root / "hard-link.md")
            with self.assertRaises(PortableTreeError):
                portable_tree_digest(root)

        if hasattr(os, "mkfifo"):
            with tempfile.TemporaryDirectory() as temp_dir:
                root = Path(temp_dir)
                os.mkfifo(root / "named-pipe")
                with self.assertRaises(PortableTreeError):
                    portable_tree_digest(root)

    def test_sync_writes_immutable_v2_snapshots_and_is_byte_idempotent(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            checkout = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            source_dir, revision = write_git_setup_sources(checkout, skill_root)

            updated, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=False)
            after_update = snapshot_files(skill_root)
            repeated, repeated_invalid = sync_setup_snapshots(skill_root, source_dir, check=False)

            self.assertFalse(invalid_input, updated.findings)
            self.assertTrue(updated.ok)
            self.assertFalse(repeated_invalid)
            self.assertEqual(repeated.findings, [])
            self.assertEqual(snapshot_files(skill_root), after_update)
            snapshot = read_snapshot(skill_root, "go-cli")
            self.assertEqual(
                snapshot["schemaVersion"],
                "setup-context-driven/setup-snapshot/0.0.1",
            )
            self.assertEqual(snapshot["version"], "0.0.1")
            self.assertEqual(snapshot["source"]["repository"], "example/skills")
            external = next(
                skill for skill in snapshot["skills"] if skill["source"]["type"] == "github"
            )
            self.assertEqual(external["source"]["ref"], revision)
            self.assertRegex(external["treeDigest"], r"^[0-9a-f]{64}$")

    def test_sync_rejects_dirty_or_untracked_source_bytes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            checkout = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            source_dir, _ = write_git_setup_sources(checkout, skill_root)
            (checkout / "untracked.txt").write_text("dirty\n", encoding="utf-8")

            result, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertTrue(invalid_input)
            self.assertFalse(result.ok)
            self.assertIn("dirty or untracked", result.findings[0].message)

    def test_sync_rejects_mutable_ref_and_content_ref_mismatch(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            checkout = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            source_dir, source_revision = write_git_setup_sources(checkout, skill_root)
            setup = setup_source_json(skill_root, "rust-cli")
            external = next(
                skill for skill in setup["skills"] if skill["name"] not in REPO_OWNED_FIXTURE_SKILLS
            )
            external["source"] = {
                "type": "github",
                "repository": "example/skills",
                "ref": "main",
                "path": external["path"],
            }
            write_json(source_dir / "rust-cli.json", setup)
            commit_all(checkout, "mutable ref fixture")

            mutable, mutable_invalid = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertTrue(mutable_invalid)
            self.assertTrue(any("full immutable commit" in item.message for item in mutable.findings))

            external["source"]["ref"] = source_revision
            source_skill = checkout / external["path"] / "SKILL.md"
            source_skill.write_text("changed after declared ref\n", encoding="utf-8")
            write_json(source_dir / "rust-cli.json", setup)
            commit_all(checkout, "mismatched ref fixture")
            mismatched, mismatch_invalid = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertTrue(mismatch_invalid)
            self.assertTrue(any("bytes do not match commit" in item.message for item in mismatched.findings))

    def test_repo_owned_digest_precedes_same_path_external_checkout_content(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            checkout = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            before = read_snapshot(skill_root, "go-cli")
            expected = next(
                skill["contentDigest"]
                for skill in before["skills"]
                if skill["name"] == "setup-context-driven"
            )
            source_dir, _ = write_git_setup_sources(checkout, skill_root)
            source_skill = checkout / "skills" / "00-setup" / "setup-context-driven" / "SKILL.md"
            source_skill.write_text("conflicting external bytes\n", encoding="utf-8")
            commit_all(checkout, "conflicting repo-owned source")

            result, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=False)

            self.assertFalse(invalid_input, result.findings)
            self.assertTrue(result.ok)
            after = read_snapshot(skill_root, "go-cli")
            repo_owned = next(
                skill for skill in after["skills"] if skill["name"] == "setup-context-driven"
            )
            self.assertEqual(repo_owned["source"], {"type": "repo", "name": "roundfix"})
            self.assertEqual(repo_owned["contentDigest"], expected)

    def test_sync_rejects_unsafe_source_path_and_symlinked_skill_entry(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            checkout = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            source_dir, _ = write_git_setup_sources(checkout, skill_root)
            setup = setup_source_json(skill_root, "rust-cli")
            setup["skills"][0]["path"] = str(checkout / "absolute")
            write_json(source_dir / "rust-cli.json", setup)
            commit_all(checkout, "unsafe path fixture")

            unsafe, unsafe_invalid = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertTrue(unsafe_invalid)
            self.assertTrue(any("not portable" in item.message for item in unsafe.findings))

        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir)
            skill_root = temp_root / "skill"
            checkout = temp_root / "canonical"
            clone_assets_to(SKILL_ROOT, skill_root)
            source_dir, _ = write_git_setup_sources(checkout, skill_root)
            external_path = next(
                Path(skill["path"])
                for skill in read_snapshot(skill_root, "rust-cli")["skills"]
                if skill["name"] not in REPO_OWNED_FIXTURE_SKILLS
            )
            target = checkout / external_path / "target.md"
            target.write_text("target\n", encoding="utf-8")
            (checkout / external_path / "link.md").symlink_to(target)
            commit_all(checkout, "unsafe symlink fixture")

            result, invalid_input = sync_setup_snapshots(skill_root, source_dir, check=True)

            self.assertTrue(invalid_input)
            self.assertTrue(any("symbolic link" in item.message for item in result.findings))


REPO_OWNED_FIXTURE_SKILLS = {
    "archive-spec",
    "brainstorming",
    "business-analyst",
    "council",
    "evidence-gate",
    "implement-spec",
    "implement-task",
    "qa-gate",
    "roundfix",
    "setup-context-driven",
    "write-idea",
    "write-prd",
    "write-tasks",
    "write-techspec",
}


def write_git_setup_sources(checkout, skill_root):
    source_dir = checkout / "setups"
    source_dir.mkdir(parents=True)
    paths = set()
    for source in (skill_root / "assets" / "setups").glob("*.json"):
        snapshot = json.loads(source.read_text(encoding="utf-8"))
        lines = ["# canonical setup"]
        for skill in snapshot["skills"]:
            lines.append(skill["path"])
            paths.add(skill["path"])
        (source_dir / f"{snapshot['id']}.txt").write_text(
            "\n".join(lines) + "\n",
            encoding="utf-8",
        )
    for path in paths:
        skill_file = checkout / path / "SKILL.md"
        skill_file.parent.mkdir(parents=True, exist_ok=True)
        skill_file.write_text(f"# {Path(path).name}\n", encoding="utf-8")
    subprocess.run(["git", "init", "-q", str(checkout)], check=True)
    subprocess.run(
        ["git", "-C", str(checkout), "config", "user.email", "fixture@example.com"],
        check=True,
    )
    subprocess.run(
        ["git", "-C", str(checkout), "config", "user.name", "Fixture"],
        check=True,
    )
    subprocess.run(
        ["git", "-C", str(checkout), "config", "commit.gpgsign", "false"],
        check=True,
    )
    subprocess.run(
        ["git", "-C", str(checkout), "remote", "add", "origin", "https://github.com/example/skills.git"],
        check=True,
    )
    revision = commit_all(checkout, "fixture source")
    return source_dir, revision


def commit_all(checkout, message):
    subprocess.run(["git", "-C", str(checkout), "add", "."], check=True)
    subprocess.run(
        ["git", "-C", str(checkout), "commit", "-q", "-m", message],
        check=True,
    )
    return subprocess.run(
        ["git", "-C", str(checkout), "rev-parse", "HEAD"],
        check=True,
        text=True,
        capture_output=True,
    ).stdout.strip()


def setup_source_json(skill_root, setup_id):
    snapshot = read_snapshot(skill_root, setup_id)
    return {
        "skills": [
            {"name": skill["name"], "path": skill["path"]}
            for skill in snapshot["skills"]
        ]
    }


def read_snapshot(skill_root, setup_id):
    return json.loads(
        (skill_root / "assets" / "setups" / f"{setup_id}.json").read_text(
            encoding="utf-8"
        )
    )


def write_json(path, payload):
    path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")


if __name__ == "__main__":
    unittest.main()

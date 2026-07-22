# Suite: Repository-Owned Extension lifecycle
# Invariant: setup creates the unmarked extension at most once and never manages its bytes.
# Boundary IN: decision resolution, Change Plan authorization, atomic apply, audit, and profile transitions.
# Boundary OUT: project-specific rule content and formatter behavior.

import hashlib
import json
import tempfile
import unittest
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
EXTENSION_PATH = Path("docs/agents/repository.md")

from test_apply import (  # noqa: E402
    BASE_DECISIONS,
    context_setup,
    run_apply,
    run_apply_in_process,
    run_apply_preview,
    run_apply_with_confirmation,
)
from test_audit import install_profile_skills, snapshot_files  # noqa: E402
from test_apply import run_audit  # noqa: E402


class RepositoryExtensionTests(unittest.TestCase):
    def test_confirmed_first_creation_is_unmarked_and_outside_managed_inventory(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = extension_decisions(True)

            preview = run_apply_preview(repo, "rust-cli", decisions)
            payload = json.loads(preview.stdout)
            extension_changes = changes_for_extension(payload)

            self.assertEqual(preview.returncode, 3, preview.stderr)
            self.assertRegex(payload["planDigest"], r"^[0-9a-f]{64}$")
            self.assertEqual(len(extension_changes), 1, payload)
            self.assertEqual(extension_changes[0]["action"], "create repository extension")
            self.assertIsNone(extension_changes[0]["beforeDigest"])
            self.assertEqual(
                extension_changes[0]["afterDigest"],
                hashlib.sha256(scaffold_bytes()).hexdigest(),
            )

            applied = run_apply_with_confirmation(
                repo,
                "rust-cli",
                decisions,
                payload["planDigest"],
            )

            self.assertEqual(applied.returncode, 0, applied.stderr)
            self.assertEqual((repo / EXTENSION_PATH).read_bytes(), scaffold_bytes())
            self.assertNotIn(b"setup-context-driven:", scaffold_bytes())
            manifest = read_manifest(repo)
            self.assertNotIn(
                EXTENSION_PATH.as_posix(),
                {artifact["path"] for artifact in manifest["managedArtifacts"]},
            )
            self.assertEqual(
                manifest["repositoryExtensions"],
                [
                    {
                        "id": "extension.repository-rules",
                        "path": EXTENSION_PATH.as_posix(),
                    }
                ],
            )

    def test_existing_extension_is_never_replaced_or_inventoried(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            extension = repo / EXTENSION_PATH
            extension.parent.mkdir(parents=True)
            original = b"# Local rules\n\nMUST preserve these exact bytes.\n"
            extension.write_bytes(original)

            preview = run_apply_preview(repo, "rust-cli", extension_decisions(True))
            self.assertFalse(changes_for_extension(json.loads(preview.stdout)))

            applied = run_apply(repo, "rust-cli", extension_decisions(True))

            self.assertEqual(applied.returncode, 0, applied.stderr)
            self.assertEqual(extension.read_bytes(), original)
            self.assertNotIn(
                EXTENSION_PATH.as_posix(),
                {artifact["path"] for artifact in read_manifest(repo)["managedArtifacts"]},
            )

    def test_audit_reapply_and_profile_transition_preserve_modified_bytes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            first = run_apply(repo, "rust-cli", extension_decisions(True))
            self.assertEqual(first.returncode, 0, first.stderr)
            install_profile_skills(repo, "rust-cli")
            extension = repo / EXTENSION_PATH
            modified = b"# Repository rules\r\n\r\nProject-authored bytes.\r\n"
            extension.write_bytes(modified)

            audited = run_audit(repo)
            reapplied = run_apply(repo, "rust-cli", [])
            transitioned = run_apply(repo, "go-cli-tui", extension_decisions(False))

            self.assertEqual(audited.returncode, 0, audited.stderr)
            self.assertEqual(reapplied.returncode, 0, reapplied.stderr)
            self.assertEqual(transitioned.returncode, 0, transitioned.stderr)
            self.assertEqual(extension.read_bytes(), modified)

    def test_false_decision_neither_creates_nor_removes_extension(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            disabled = extension_decisions(False)

            first = run_apply(repo, "rust-cli", disabled)

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertFalse((repo / EXTENSION_PATH).exists())

            extension = repo / EXTENSION_PATH
            extension.write_bytes(b"repository-authored\n")
            second = run_apply(repo, "rust-cli", disabled)

            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(extension.read_bytes(), b"repository-authored\n")

    def test_missing_previously_selected_extension_is_reported_without_recreation(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = extension_decisions(True)
            first = run_apply(repo, "rust-cli", decisions)
            self.assertEqual(first.returncode, 0, first.stderr)
            install_profile_skills(repo, "rust-cli")
            extension = repo / EXTENSION_PATH
            extension.unlink()

            audited = run_audit(repo)
            audit_payload = json.loads(audited.stdout)
            preview = run_apply_preview(repo, "rust-cli", decisions)
            preview_payload = json.loads(preview.stdout)

            self.assertEqual(audited.returncode, 1, audited.stderr)
            self.assertEqual(preview.returncode, 1, preview.stderr)
            self.assertEqual(
                [
                    finding["code"]
                    for finding in audit_payload["findings"]
                    if finding["path"] == EXTENSION_PATH.as_posix()
                ],
                ["reference.repository.missing"],
            )
            self.assertFalse(changes_for_extension(preview_payload))
            self.assertFalse(extension.exists())

    def test_atomic_failure_rolls_back_the_first_creation(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = extension_decisions(True)
            before = snapshot_files(repo)
            preview = run_apply_in_process(
                repo,
                "rust-cli",
                decisions,
                auto_confirm=False,
            )
            digest = json.loads(preview.stdout)["planDigest"]
            original_replace = context_setup.Path.replace

            def replace_then_fail(source, target):
                result = original_replace(source, target)
                if Path(target).name == "setup-context.json":
                    raise OSError("injected atomic failure")
                return result

            with mock.patch.object(
                context_setup.Path,
                "replace",
                autospec=True,
                side_effect=replace_then_fail,
            ):
                applied = run_apply_in_process(
                    repo,
                    "rust-cli",
                    decisions,
                    confirm_plan=digest,
                    auto_confirm=False,
                )

            self.assertEqual(applied.returncode, 1, applied.stderr)
            self.assertEqual(snapshot_files(repo), before)
            self.assertFalse((repo / EXTENSION_PATH).exists())


def extension_decisions(enabled):
    replacement = f"repository.extension.enabled={str(enabled).lower()}"
    decisions = [
        replacement
        if decision.startswith("repository.extension.enabled=")
        else decision
        for decision in BASE_DECISIONS
    ]
    if not any(
        decision.startswith("repository.extension.enabled=")
        for decision in BASE_DECISIONS
    ):
        decisions.append(replacement)
    return decisions


def changes_for_extension(payload):
    return [
        change
        for change in payload["plannedChanges"]
        if change["path"] == EXTENSION_PATH.as_posix()
    ]


def scaffold_bytes():
    return (
        SKILL_ROOT / "assets" / "templates" / "extensions" / "repository.md"
    ).read_bytes()


def read_manifest(repo):
    return json.loads(
        (repo / "docs" / "agents" / "setup-context.json").read_text(
            encoding="utf-8"
        )
    )


if __name__ == "__main__":
    unittest.main()

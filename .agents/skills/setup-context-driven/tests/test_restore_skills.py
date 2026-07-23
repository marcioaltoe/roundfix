"""Suite: immutable external Repository Skill Set restoration.
Invariant: only a digest-confirmed plan may replace selected skill trees and lock entries.
Boundary IN: the real command, disposable Git sources, repository swaps, and lock generation.
Boundary OUT: public GitHub network availability and the future Doctor Command integration.
"""

import copy
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest
from unittest import mock
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_assets import (  # noqa: E402
    SkillSourceRef,
    clone_assets_to,
    load_asset_catalog,
    portable_file_digest,
    setup_snapshot_digest,
)
from context_setup import (  # noqa: E402
    ExternalSkillLockAdapter,
    GitSkillSource,
    RestoreError,
    RestoreFilesystem,
    RestoreLimits,
    apply_restore_plan,
    build_restore_plan,
)
from test_audit import (  # noqa: E402
    snapshot_files,
    write_compliant_repository,
    write_fixture_external_digests,
)


class RestoreSkillsCliTests(unittest.TestCase):
    def test_versioned_lock_fixture_matches_real_production_adapter(self):
        ExternalSkillLockAdapter().assert_compatible()

    def test_real_lock_fixture_failures_are_actionable_and_non_mutating(self):
        cases = {
            "missing": lambda path: path.unlink(),
            "malformed": lambda path: path.write_text("{not json", encoding="utf-8"),
            "digest mismatch": self._write_mismatching_lock_fixture,
        }

        for name, mutate in cases.items():
            with self.subTest(name=name), RestoreFixture() as fixture:
                fixture.add_source_skill(
                    "agentic-cli-design", {"SKILL.md": b"# restored\n"}
                )
                fixture.commit_and_configure("agentic-cli-design")
                before = snapshot_files(fixture.repo)
                mutate(fixture.lock_compatibility_fixture)

                result = fixture.restore("--skill", "agentic-cli-design")

                self.assertEqual(result.returncode, 1, result.stderr)
                finding = json.loads(result.stdout)["finding"]
                self.assertEqual(finding["code"], "lock.adapter-incompatible")
                self.assertEqual(
                    finding["action"],
                    "Update the isolated lock adapter before restoring skills.",
                )
                self.assertEqual(snapshot_files(fixture.repo), before)

    @staticmethod
    def _write_mismatching_lock_fixture(path):
        document = json.loads(path.read_text(encoding="utf-8"))
        document["expectedSha256"] = "0" * 64
        path.write_text(json.dumps(document), encoding="utf-8")

    def test_help_exposes_non_interactive_restore_contract(self):
        with RestoreFixture() as fixture:
            result = fixture.run("restore-skills", "--help")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stderr, "")
        self.assertIn("--skill", result.stdout)
        self.assertIn("--source-dir", result.stdout)
        self.assertIn("--confirm-plan", result.stdout)
        self.assertIn("Exit codes:", result.stdout)

    def test_text_preview_is_deterministic_and_names_confirmation(self):
        with RestoreFixture() as fixture:
            fixture.add_source_skill("agentic-cli-design", {"SKILL.md": b"# restored\n"})
            fixture.commit_and_configure("agentic-cli-design")
            shutil.rmtree(fixture.repo / ".agents" / "skills" / "agentic-cli-design")

            result = fixture.run(
                "restore-skills",
                "--repo",
                str(fixture.repo),
                "--profile",
                fixture.profile,
                "--source-dir",
                str(fixture.source),
                "--skill",
                "agentic-cli-design",
                "--format",
                "text",
            )

            self.assertEqual(result.returncode, 3, result.stderr)
            self.assertIn("restore-skills: blocked", result.stdout)
            self.assertIn("plan.confirmation.required", result.stdout)
            self.assertIn("planDigest:", result.stdout)
            self.assertIn("create .agents/skills/agentic-cli-design/SKILL.md", result.stdout)

    def test_preview_names_nested_file_and_exact_lock_operations(self):
        with RestoreFixture() as fixture:
            fixture.add_source_skill(
                "agentic-cli-design",
                {
                    "SKILL.md": b"# restored\n",
                    "references/guide.md": b"restored guide\n",
                    "references/new.md": b"new\n",
                },
            )
            fixture.commit_and_configure("agentic-cli-design")
            target = fixture.repo / ".agents" / "skills" / "agentic-cli-design"
            shutil.rmtree(target)
            target.mkdir(parents=True)
            (target / "SKILL.md").write_bytes(b"# drifted\n")
            (target / "references").mkdir()
            (target / "references" / "removed.md").write_bytes(b"remove me\n")
            fixture.write_lock(
                {
                    "unrelated": {
                        "source": "elsewhere/skills",
                        "sourceType": "github",
                        "skillPath": "skills/unrelated/SKILL.md",
                        "computedHash": "a" * 64,
                    }
                }
            )
            before = snapshot_files(fixture.repo)

            result = fixture.restore("--skill", "agentic-cli-design")

            self.assertEqual(result.returncode, 3, result.stderr)
            payload = json.loads(result.stdout)
            self.assertEqual(payload["schemaVersion"], "setup-context-driven/restore-v1")
            self.assertEqual(payload["finding"]["code"], "plan.confirmation.required")
            operations = {
                (item["action"], item["path"]): item
                for item in payload["plannedChanges"]
            }
            self.assertIn(
                ("refresh", ".agents/skills/agentic-cli-design/SKILL.md"),
                operations,
            )
            self.assertIn(
                ("create", ".agents/skills/agentic-cli-design/references/guide.md"),
                operations,
            )
            self.assertIn(
                ("create", ".agents/skills/agentic-cli-design/references/new.md"),
                operations,
            )
            self.assertIn(
                ("remove", ".agents/skills/agentic-cli-design/references/removed.md"),
                operations,
            )
            lock_edit = operations[("update-lock-entry", "skills-lock.json")]
            self.assertEqual(lock_edit["skill"], "agentic-cli-design")
            self.assertEqual(lock_edit["after"]["ref"], fixture.revision)
            self.assertEqual(snapshot_files(fixture.repo), before)

    def test_confirmation_restores_exact_commit_with_portable_lock_and_is_idempotent(self):
        with RestoreFixture() as fixture:
            files = {
                "SKILL.md": b"# restored\n",
                "references/guide.md": b"guide\n",
            }
            fixture.add_source_skill("agentic-cli-design", files)
            fixture.commit_and_configure("agentic-cli-design")
            shutil.rmtree(
                fixture.repo / ".agents" / "skills" / "agentic-cli-design"
            )
            unrelated_tree = fixture.repo / ".agents" / "skills" / "unrelated"
            unrelated_tree.mkdir(parents=True)
            (unrelated_tree / "SKILL.md").write_bytes(b"unchanged\n")
            unrelated_entry = {
                "source": "elsewhere/skills",
                "sourceType": "github",
                "skillPath": "skills/unrelated/SKILL.md",
                "computedHash": "b" * 64,
            }
            fixture.write_lock({"unrelated": unrelated_entry})

            preview = fixture.restore("--skill", "agentic-cli-design")
            digest = json.loads(preview.stdout)["planDigest"]
            applied = fixture.restore(
                "--skill",
                "agentic-cli-design",
                "--confirm-plan",
                digest,
            )

            self.assertEqual(applied.returncode, 0, applied.stderr)
            target = fixture.repo / ".agents" / "skills" / "agentic-cli-design"
            self.assertEqual(
                {path.relative_to(target).as_posix(): path.read_bytes() for path in target.rglob("*") if path.is_file()},
                files,
            )
            lock = json.loads((fixture.repo / "skills-lock.json").read_text(encoding="utf-8"))
            entry = lock["skills"]["agentic-cli-design"]
            self.assertEqual(entry["source"], "example/skills")
            self.assertEqual(entry["ref"], fixture.revision)
            self.assertEqual(entry["skillPath"], fixture.source_path("agentic-cli-design") + "/SKILL.md")
            self.assertNotIn(str(fixture.source), json.dumps(entry))
            self.assertEqual(lock["skills"]["unrelated"], unrelated_entry)
            self.assertEqual((unrelated_tree / "SKILL.md").read_bytes(), b"unchanged\n")

            second = fixture.restore("--skill", "agentic-cli-design")
            audit = fixture.run(
                "audit", "--repo", str(fixture.repo), "--profile", "rust-cli", "--format", "json"
            )

            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(json.loads(second.stdout)["plannedChanges"], [])
            self.assertEqual(audit.returncode, 0, audit.stdout + audit.stderr)

    def test_multiple_skills_from_one_ref_report_one_acquisition(self):
        with RestoreFixture() as fixture:
            fixture.add_source_skill("agentic-cli-design", {"SKILL.md": b"# agentic\n"})
            fixture.add_source_skill(
                "domain-modeling",
                {"SKILL.md": b"# domain\n", "references/guide.md": b"guide\n"},
            )
            fixture.commit_and_configure("agentic-cli-design", "domain-modeling")
            for name in ("agentic-cli-design", "domain-modeling"):
                shutil.rmtree(fixture.repo / ".agents" / "skills" / name)

            result = fixture.restore(
                "--skill", "agentic-cli-design", "--skill", "domain-modeling"
            )

            self.assertEqual(result.returncode, 3, result.stderr)
            payload = json.loads(result.stdout)
            self.assertEqual(len(payload["acquisitions"]), 1)
            self.assertEqual(
                {item["skill"] for item in payload["skills"]},
                {"agentic-cli-design", "domain-modeling"},
            )
            self.assertTrue(all(item["expectedDigest"] for item in payload["skills"]))
            self.assertTrue(all(item["changes"] for item in payload["skills"]))

    def test_stale_and_malformed_confirmation_never_mutate(self):
        with RestoreFixture() as fixture:
            fixture.add_source_skill("agentic-cli-design", {"SKILL.md": b"# restored\n"})
            fixture.commit_and_configure("agentic-cli-design")
            target = fixture.repo / ".agents" / "skills" / "agentic-cli-design" / "SKILL.md"
            target.write_bytes(b"# first drift\n")
            preview = fixture.restore("--skill", "agentic-cli-design")
            old_digest = json.loads(preview.stdout)["planDigest"]
            target.write_bytes(b"# changed after preview\n")
            before_stale = snapshot_files(fixture.repo)

            stale = fixture.restore(
                "--skill", "agentic-cli-design", "--confirm-plan", old_digest
            )
            malformed = fixture.restore(
                "--skill", "agentic-cli-design", "--confirm-plan", "not-a-digest"
            )

            self.assertEqual(stale.returncode, 3, stale.stderr)
            self.assertEqual(json.loads(stale.stdout)["finding"]["code"], "plan.confirmation.stale")
            self.assertEqual(malformed.returncode, 2, malformed.stderr)
            self.assertEqual(json.loads(malformed.stdout)["finding"]["code"], "plan.confirmation.invalid")
            self.assertEqual(snapshot_files(fixture.repo), before_stale)

    def test_invalid_skill_filter_and_malformed_lock_exit_two_without_mutation(self):
        with RestoreFixture() as fixture:
            before = snapshot_files(fixture.repo)
            invalid_skill = fixture.restore("--skill", "not-required")
            (fixture.repo / "skills-lock.json").write_text("{not json", encoding="utf-8")
            before_lock_error = snapshot_files(fixture.repo)
            malformed_lock = fixture.restore("--skill", "agentic-cli-design")

            self.assertEqual(invalid_skill.returncode, 2, invalid_skill.stderr)
            self.assertEqual(json.loads(invalid_skill.stdout)["finding"]["code"], "restore.skill-invalid")
            self.assertEqual(malformed_lock.returncode, 2, malformed_lock.stderr)
            self.assertEqual(json.loads(malformed_lock.stdout)["finding"]["code"], "lock.invalid")
            self.assertEqual(
                {key: value for key, value in snapshot_files(fixture.repo).items() if key != "skills-lock.json"},
                {key: value for key, value in before.items() if key != "skills-lock.json"},
            )
            self.assertEqual(snapshot_files(fixture.repo), before_lock_error)

    def test_source_and_security_failures_happen_before_mutation(self):
        cases = ("digest", "unsafe", "size")
        for case in cases:
            with self.subTest(case=case), RestoreFixture() as fixture:
                files = {"SKILL.md": b"# restored\n"}
                if case == "size":
                    files["large.bin"] = b"x" * (RestoreLimits().max_bytes + 1)
                fixture.add_source_skill("agentic-cli-design", files)
                if case == "unsafe":
                    skill_root = fixture.source / fixture.source_path("agentic-cli-design")
                    (skill_root / "unsafe-link").symlink_to("SKILL.md")
                fixture.commit_and_configure(
                    "agentic-cli-design", digest_override=("0" * 64 if case == "digest" else None)
                )
                before = snapshot_files(fixture.repo)

                result = fixture.restore("--skill", "agentic-cli-design")

                self.assertEqual(result.returncode, 1, result.stderr)
                finding = json.loads(result.stdout)["finding"]
                self.assertTrue(finding["action"])
                self.assertIn(finding["code"], {"source.digest-mismatch", "source.unsafe-tree", "source.limit-exceeded"})
                self.assertEqual(snapshot_files(fixture.repo), before)

    def test_missing_git_and_unreachable_commit_are_actionable_and_read_only(self):
        with RestoreFixture() as fixture:
            fixture.add_source_skill("agentic-cli-design", {"SKILL.md": b"# restored\n"})
            fixture.commit_and_configure("agentic-cli-design")
            before = snapshot_files(fixture.repo)

            missing = fixture.restore(
                "--skill", "agentic-cli-design", extra_env={"PATH": "/nonexistent"}
            )
            fixture.set_ref("agentic-cli-design", "f" * 40)
            unreachable = fixture.restore("--skill", "agentic-cli-design")

            self.assertEqual(missing.returncode, 1, missing.stderr)
            self.assertEqual(json.loads(missing.stdout)["finding"]["code"], "source.git-missing")
            self.assertEqual(unreachable.returncode, 1, unreachable.stderr)
            self.assertEqual(json.loads(unreachable.stdout)["finding"]["code"], "source.commit-unavailable")
            self.assertEqual(snapshot_files(fixture.repo), before)

    def test_commit_identity_mismatch_is_rejected(self):
        source = SkillSourceRef(
            provider="github",
            repository="example/skills",
            revision="a" * 40,
            source_path=Path("skills/example"),
        )
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
            "context_setup.run_git_argv",
            side_effect=[b"", b"", ("b" * 40 + "\n").encode("ascii")],
        ):
            with self.assertRaisesRegex(RestoreError, "not the declared commit") as raised:
                GitSkillSource(Path(temp_dir)).acquire(
                    source, Path(temp_dir) / "objects"
                )

        self.assertEqual(raised.exception.code, "source.commit-mismatch")

    def test_injected_swap_lock_and_postwrite_failures_roll_back_all_targets(self):
        failure_points = ("directory", "lock", "verify")
        for failure_point in failure_points:
            with self.subTest(failure_point=failure_point), RestoreFixture() as fixture:
                fixture.add_source_skill("agentic-cli-design", {"SKILL.md": b"# restored\n"})
                fixture.commit_and_configure("agentic-cli-design")
                target = fixture.repo / ".agents" / "skills" / "agentic-cli-design" / "SKILL.md"
                target.write_bytes(b"# drifted\n")
                fixture.write_lock({"unrelated": {"computedHash": "c" * 64}})
                catalog = load_asset_catalog(fixture.skill_root)
                plan = build_restore_plan(
                    fixture.repo,
                    catalog,
                    "rust-cli",
                    ["agentic-cli-design"],
                    fixture.source,
                )
                before = snapshot_files(fixture.repo)

                with self.assertRaises(RestoreError):
                    apply_restore_plan(
                        fixture.repo,
                        plan,
                        filesystem=FailingFilesystem(failure_point),
                    )

                self.assertEqual(snapshot_files(fixture.repo), before)

    def test_incompatible_lock_adapter_and_unsupported_provider_block_before_write(self):
        with RestoreFixture() as fixture:
            fixture.add_source_skill("agentic-cli-design", {"SKILL.md": b"# restored\n"})
            fixture.commit_and_configure("agentic-cli-design")
            catalog = load_asset_catalog(fixture.skill_root)
            before = snapshot_files(fixture.repo)

            with self.assertRaisesRegex(RestoreError, "compatibility"):
                build_restore_plan(
                    fixture.repo,
                    catalog,
                    "rust-cli",
                    ["agentic-cli-design"],
                    fixture.source,
                    lock_adapter=IncompatibleLockAdapter(),
                )

            mutated = copy.deepcopy(catalog)
            setup_id = mutated.profiles["rust-cli"]["setup"]
            contract = next(
                item
                for item in mutated.external_sources_by_setup[setup_id]
                if item.skill_name == "agentic-cli-design"
            )
            object.__setattr__(contract.source, "provider", "unsupported")
            with self.assertRaisesRegex(RestoreError, "provider"):
                build_restore_plan(
                    fixture.repo,
                    mutated,
                    "rust-cli",
                    ["agentic-cli-design"],
                    fixture.source,
                )

            self.assertEqual(snapshot_files(fixture.repo), before)


class FailingFilesystem(RestoreFilesystem):
    def __init__(self, failure_point):
        super().__init__()
        self.failure_point = failure_point

    def replace(self, source, target):
        if self.failure_point == "directory" and source.name.startswith(".restore-stage-"):
            raise OSError("injected directory swap failure")
        if self.failure_point == "lock" and target.name == "skills-lock.json" and source.name.startswith(".restore-lock-"):
            raise OSError("injected lock write failure")
        return super().replace(source, target)

    def verify_tree(self, target, expected_digest):
        if self.failure_point == "verify" and not target.name.startswith(".restore-stage-"):
            raise OSError("injected postwrite verification failure")
        return super().verify_tree(target, expected_digest)


class IncompatibleLockAdapter(ExternalSkillLockAdapter):
    def assert_compatible(self):
        raise RestoreError(
            "lock.adapter-incompatible",
            "The lock compatibility fixture disagrees with Spec 0036 Task 01.",
            "Update the isolated lock adapter before restoring skills.",
        )


class RestoreFixture:
    profile = "rust-cli"

    def __init__(self):
        self._temporary = tempfile.TemporaryDirectory()
        self.root = Path(self._temporary.name)
        self.skill_root = self.root / "setup-context-driven"
        self.source = self.root / "source"
        self.repo = self.root / "repo"
        clone_assets_to(SKILL_ROOT, self.skill_root)
        shutil.copytree(SKILL_ROOT / "scripts", self.skill_root / "scripts")
        write_fixture_external_digests(self.skill_root)
        write_compliant_repository(self.repo, self.profile)
        subprocess.run(["git", "init", "-q", str(self.source)], check=True)
        subprocess.run(["git", "-C", str(self.source), "config", "user.email", "fixture@example.com"], check=True)
        subprocess.run(["git", "-C", str(self.source), "config", "user.name", "Fixture"], check=True)
        subprocess.run(["git", "-C", str(self.source), "config", "commit.gpgsign", "false"], check=True)
        self.revision = None

    @property
    def lock_compatibility_fixture(self):
        return self.skill_root / "assets" / "lock-hash-compatibility-v1.json"

    def __enter__(self):
        return self

    def __exit__(self, *_):
        self._temporary.cleanup()

    def source_path(self, name):
        snapshot = self.read_snapshot()
        return next(item["source"]["path"] for item in snapshot["skills"] if item["name"] == name)

    def add_source_skill(self, name, files):
        root = self.source / self.source_path(name)
        for relative, content in files.items():
            target = root / relative
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(content)

    def commit_and_configure(self, *names, digest_override=None):
        subprocess.run(["git", "-C", str(self.source), "add", "."], check=True)
        subprocess.run(["git", "-C", str(self.source), "commit", "-q", "-m", "fixture source"], check=True)
        self.revision = subprocess.run(
            ["git", "-C", str(self.source), "rev-parse", "HEAD"],
            text=True,
            capture_output=True,
            check=True,
        ).stdout.strip()
        snapshot = self.read_snapshot()
        for item in snapshot["skills"]:
            if item["name"] not in names:
                continue
            source_path = item["source"]["path"]
            files = [
                (path.relative_to(self.source / source_path).as_posix().encode("utf-8"), path.read_bytes())
                for path in (self.source / source_path).rglob("*")
                if path.is_file() and not path.is_symlink()
            ]
            item["source"] = {
                "type": "github",
                "repository": "example/skills",
                "ref": self.revision,
                "path": source_path,
            }
            item["treeDigest"] = digest_override or portable_file_digest(files)
        snapshot["digest"] = setup_snapshot_digest(
            snapshot["skills"], snapshot.get("activationBundles")
        )
        self.write_snapshot(snapshot)

    def set_ref(self, name, revision):
        snapshot = self.read_snapshot()
        item = next(item for item in snapshot["skills"] if item["name"] == name)
        item["source"]["ref"] = revision
        snapshot["digest"] = setup_snapshot_digest(
            snapshot["skills"], snapshot.get("activationBundles")
        )
        self.write_snapshot(snapshot)

    def read_snapshot(self):
        return json.loads((self.skill_root / "assets" / "setups" / "rust-cli.json").read_text(encoding="utf-8"))

    def write_snapshot(self, snapshot):
        (self.skill_root / "assets" / "setups" / "rust-cli.json").write_text(
            json.dumps(snapshot, indent=2) + "\n", encoding="utf-8"
        )

    def write_lock(self, entries):
        (self.repo / "skills-lock.json").write_text(
            json.dumps({"version": 1, "skills": entries}, indent=2) + "\n",
            encoding="utf-8",
        )

    def restore(self, *args, extra_env=None):
        return self.run(
            "restore-skills",
            "--repo",
            str(self.repo),
            "--profile",
            self.profile,
            "--source-dir",
            str(self.source),
            "--format",
            "json",
            *args,
            extra_env=extra_env,
        )

    def run(self, *args, extra_env=None):
        environment = os.environ.copy()
        environment["PYTHONDONTWRITEBYTECODE"] = "1"
        if extra_env:
            environment.update(extra_env)
        return subprocess.run(
            [sys.executable, str(self.skill_root / "scripts" / "context_setup.py"), *args],
            text=True,
            capture_output=True,
            check=False,
            env=environment,
        )


if __name__ == "__main__":
    unittest.main()

"""Behavior tests for confirmed Baseline Readoption application.

Suite: one immutable Readoption Change Plan
Invariant: previewed bytes are the only bytes apply may write, and every failure
restores all targets mutated by the confirmed plan.
Boundary IN: strict profile/capability composition, confirmation, atomic apply,
manifest/snapshot output, rollback, audit, and reapply.
Boundary OUT: profile migration for Go/Rust and final all-profile journeys.
"""

import base64
import hashlib
import io
import json
import subprocess
import sys
import tempfile
import unittest
from contextlib import redirect_stderr, redirect_stdout
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

import context_setup  # noqa: E402
from context_baseline import inventory_incompatible_source_baseline  # noqa: E402
from test_audit import (  # noqa: E402
    install_profile_skills,
    run_context_setup,
    snapshot_files,
)
from test_source_inventory import SourceInventoryTests  # noqa: E402


PROFILE = "standard-typescript-monorepo"
DECISION_SCHEMA = "setup-context-driven/decisions/0.0.1"
MANIFEST = Path("docs/agents/setup-context.json")
REPOSITORY_RULES = Path("docs/agents/repository-rules.md")


class ReadoptionApplyTests(unittest.TestCase):
    def test_preview_apply_reapply_and_audit_share_exact_strict_plan(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo, decision_file, original_guide, proposed_rules = self.fixture(
                Path(temp_dir)
            )
            before = snapshot_files(repo)

            preview = run_apply(repo, decision_file)
            preview_payload = json.loads(preview.stdout)

            self.assertEqual(preview.returncode, 3, preview.stderr)
            self.assertEqual(preview.stderr, "")
            self.assertFinding(preview_payload, "plan.confirmation.required")
            self.assertTrue(preview_payload["capabilities"])
            self.assertTrue(preview_payload["verification"])
            self.assertEqual(
                preview_payload["setupSnapshot"]["schemaVersion"],
                "setup-context-driven/profile-snapshot/0.0.1",
            )
            self.assertEqual(snapshot_files(repo), before)

            applied = run_apply(
                repo, decision_file, preview_payload["planDigest"]
            )
            applied_payload = json.loads(applied.stdout)

            self.assertEqual(applied.returncode, 0, applied.stderr)
            self.assertEqual(applied.stderr, "")
            self.assertEqual(
                applied_payload["planDigest"], preview_payload["planDigest"]
            )
            self.assertEqual(
                applied_payload["plannedOutputs"],
                preview_payload["plannedOutputs"],
            )
            for output in preview_payload["plannedOutputs"]:
                target = repo / output["path"]
                expected = (
                    base64.b64decode(output["afterBytes"])
                    if output["afterBytes"] is not None
                    else None
                )
                self.assertEqual(target.read_bytes() if target.is_file() else None, expected)

            manifest = json.loads((repo / MANIFEST).read_text(encoding="utf-8"))
            self.assertEqual(
                manifest["schemaVersion"], "setup-context-driven/manifest/0.0.1"
            )
            self.assertEqual(manifest["version"], "0.0.1")
            self.assertEqual(manifest["generator"]["version"], "0.0.1")
            self.assertEqual(manifest["setupSnapshot"], preview_payload["setupSnapshot"])
            self.assertTrue(
                all(item["version"] == "0.0.1" for item in manifest["managedArtifacts"])
            )
            self.assertEqual((repo / REPOSITORY_RULES).read_bytes(), proposed_rules)
            self.assertEqual(
                (repo / "docs/agents/guide.md").read_bytes(), original_guide
            )
            self.assertNotIn("id=root.old", (repo / "AGENTS.md").read_text(encoding="utf-8"))

            after = snapshot_files(repo)
            second = run_apply(repo, decision_file)
            audit = run_context_setup(
                "audit", "--repo", str(repo), "--format", "json"
            )

            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(json.loads(second.stdout)["plannedChanges"], [])
            self.assertEqual(audit.returncode, 0, audit.stderr)
            self.assertEqual(snapshot_files(repo), after)

    def test_missing_required_capability_and_stale_confirmation_write_nothing(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo, decision_file, _, _ = self.fixture(Path(temp_dir))
            context7 = repo / ".agents/skills/context7/SKILL.md"
            context7.unlink()
            before = snapshot_files(repo)

            missing = run_apply(repo, decision_file)

            self.assertEqual(missing.returncode, 1, missing.stderr)
            self.assertFinding(json.loads(missing.stdout), "capability.required.missing")
            self.assertEqual(snapshot_files(repo), before)

            context7.parent.mkdir(parents=True, exist_ok=True)
            context7.write_text("# context7\n", encoding="utf-8")
            preview = run_apply(repo, decision_file)
            digest = json.loads(preview.stdout)["planDigest"]
            package_path = repo / "package.json"
            package_path.write_text(
                package_path.read_text(encoding="utf-8") + "\n",
                encoding="utf-8",
            )
            changed = snapshot_files(repo)

            stale = run_apply(repo, decision_file, digest)

            self.assertEqual(stale.returncode, 3, stale.stderr)
            self.assertFinding(json.loads(stale.stdout), "plan.confirmation.stale")
            self.assertEqual(snapshot_files(repo), changed)

    def test_missing_normalized_decision_writes_nothing(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo, decision_file, _, _ = self.fixture(Path(temp_dir))
            document = json.loads(decision_file.read_text(encoding="utf-8"))
            document["decisions"] = [
                decision
                for decision in document["decisions"]
                if decision["id"] != "verification.gate"
            ]
            decision_file.write_text(
                json.dumps(document, indent=2) + "\n", encoding="utf-8"
            )
            before = snapshot_files(repo)

            result = run_apply(repo, decision_file)

            self.assertEqual(result.returncode, 3, result.stderr)
            self.assertFinding(json.loads(result.stdout), "decision.required")
            self.assertEqual(snapshot_files(repo), before)

    def test_partial_failure_and_postwrite_tampering_restore_every_preimage(self):
        for failure_mode in ("replace", "postwrite"):
            with self.subTest(failure_mode=failure_mode):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo, decision_file, _, _ = self.fixture(Path(temp_dir))
                    preview = run_apply_in_process(repo, decision_file)
                    digest = json.loads(preview.stdout)["planDigest"]
                    before = snapshot_files(repo)
                    original_replace = context_setup.Path.replace
                    writes = 0

                    def replace_with_failure(source, target):
                        nonlocal writes
                        result = original_replace(source, target)
                        writes += 1
                        if failure_mode == "replace" and writes == 2:
                            raise OSError("injected second-write failure")
                        if (
                            failure_mode == "postwrite"
                            and Path(target).name == "setup-context.json"
                        ):
                            (repo / "AGENTS.md").write_text(
                                "tampered after write\n", encoding="utf-8"
                            )
                        return result

                    with mock.patch.object(
                        context_setup.Path,
                        "replace",
                        autospec=True,
                        side_effect=replace_with_failure,
                    ):
                        result = run_apply_in_process(repo, decision_file, digest)

                    self.assertEqual(result.returncode, 1, result.stderr)
                    self.assertFinding(json.loads(result.stdout), "managed.apply.failed")
                    self.assertEqual(snapshot_files(repo), before)

    def test_changed_preimage_is_not_overwritten_or_rolled_back(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo, decision_file, _, _ = self.fixture(Path(temp_dir))
            preview = run_apply_in_process(repo, decision_file)
            digest = json.loads(preview.stdout)["planDigest"]
            original_temp_write = context_setup.write_unique_temp_bytes
            injected = False

            def write_temp_then_change_preimage(target, content, temp_paths):
                nonlocal injected
                result = original_temp_write(target, content, temp_paths)
                if not injected:
                    injected = True
                    (repo / "AGENTS.md").write_text(
                        "concurrent repository edit\n", encoding="utf-8"
                    )
                return result

            with mock.patch.object(
                context_setup,
                "write_unique_temp_bytes",
                side_effect=write_temp_then_change_preimage,
            ):
                result = run_apply_in_process(repo, decision_file, digest)

            self.assertEqual(result.returncode, 1, result.stderr)
            self.assertFinding(json.loads(result.stdout), "managed.apply.failed")
            self.assertEqual(
                (repo / "AGENTS.md").read_text(encoding="utf-8"),
                "concurrent repository edit\n",
            )
            self.assertFalse((repo / REPOSITORY_RULES).exists())

    def fixture(self, root):
        repo = root / "repo"
        SourceInventoryTests.write_incompatible_repository(repo)
        original_guide = (repo / "docs/agents/guide.md").read_bytes()
        install_profile_skills(repo, PROFILE)
        self.write_capability_evidence(repo)
        inventory = inventory_incompatible_source_baseline(
            repo, "baseline.pre-0.0.1"
        )
        proposed_rules = b"# Repository-specific normative rules\n\nKeep local policy.\n"
        document = self.decision_document(inventory, proposed_rules)
        decision_file = root / "decisions.json"
        decision_file.write_text(json.dumps(document, indent=2) + "\n", encoding="utf-8")
        return repo, decision_file, original_guide, proposed_rules

    @staticmethod
    def write_capability_evidence(repo):
        root_package = {
            "name": "fixture",
            "packageManager": "bun@1.2.0",
            "dependencies": {
                "typescript": "1",
                "turbo": "1",
                "vite": "1",
                "react": "1",
                "hono": "1",
                "drizzle-orm": "1",
                "zod": "1",
                "tailwindcss": "1",
                "shadcn": "1",
                "@tanstack/react-query": "1",
                "@tanstack/react-router": "1",
                "better-auth": "1",
                "@logtape/logtape": "1",
                "oxlint": "1",
                "oxfmt": "1",
                "vitest": "1",
            },
        }
        (repo / "package.json").write_text(
            json.dumps(root_package), encoding="utf-8"
        )
        for workspace in ("frontend", "backend"):
            path = repo / "packages" / workspace / "package.json"
            path.parent.mkdir(parents=True)
            path.write_text(json.dumps({"name": workspace}), encoding="utf-8")
        (repo / "DATABASE.md").write_text("PostgreSQL 17\n", encoding="utf-8")
        (repo / "DESIGN.md").write_text("# Design contract\n", encoding="utf-8")
        http_path = repo / "docs/architecture/http-contract.json"
        http_path.parent.mkdir(parents=True)
        http_path.write_text('{"mode":"REST"}\n', encoding="utf-8")

    def decision_document(self, inventory, proposed_rules):
        http_path = Path("docs/architecture/http-contract.json")
        dispositions = []
        for index, entry in enumerate(inventory.entries):
            if index == 0:
                dispositions.append(
                    {
                        "entryId": entry.entry_id,
                        "entryDigest": entry.digest,
                        "classification": "normative-clause",
                        "disposition": "repository-rules",
                        "destination": {
                            "documentType": "repository-rules",
                            "path": REPOSITORY_RULES.as_posix(),
                            "proposedBytes": base64.b64encode(proposed_rules).decode("ascii"),
                            "digest": hashlib.sha256(proposed_rules).hexdigest(),
                        },
                        "reason": "",
                    }
                )
            else:
                dispositions.append(
                    {
                        "entryId": entry.entry_id,
                        "entryDigest": entry.digest,
                        "classification": "non-governed",
                        "disposition": "rejected",
                        "destination": None,
                        "reason": "Structural evidence is not a current governed instruction.",
                    }
                )
        source_bytes = b'{"mode":"REST"}\n'
        decisions = [
            {"id": "language.generated", "value": "English"},
            {"id": "verification.gate", "value": "bun run verify"},
            {
                "id": "http.contract",
                "value": {
                    "mode": "REST",
                    "exceptions": [],
                    "source": {
                        "path": http_path.as_posix(),
                        "digest": hashlib.sha256(source_bytes).hexdigest(),
                    },
                },
            },
            {"id": "spec.scaffold", "value": True},
            {"id": "domain.layout", "value": "single-context"},
            {"id": "triage.external", "value": False},
            {"id": "autonomous.enabled", "value": True},
            {"id": "runtime.backend", "value": "codex gpt-5.5 xhigh"},
            {"id": "runtime.design", "value": "claude opus xhigh"},
            {"id": "secondbrain.enabled", "value": False},
            {"id": "repository.extension.enabled", "value": False},
        ]
        return {
            "schemaVersion": DECISION_SCHEMA,
            "version": "0.0.1",
            "decisions": decisions,
            "readoption": {
                "sourceBaseline": {
                    "id": inventory.baseline_id,
                    "digest": inventory.digest,
                },
                "dispositions": dispositions,
            },
        }

    @staticmethod
    def assertFinding(payload, code):
        if not any(item["code"] == code for item in payload["findings"]):
            raise AssertionError(f"missing finding {code}: {payload['findings']}")


def run_apply(repo, decision_file, confirmation=None):
    args = [
        "apply",
        "--repo",
        str(repo),
        "--profile",
        PROFILE,
        "--format",
        "json",
        "--decision-file",
        str(decision_file),
    ]
    if confirmation is not None:
        args.extend(["--confirm-plan", confirmation])
    return run_context_setup(*args)


def run_apply_in_process(repo, decision_file, confirmation=None):
    args = [
        "apply",
        "--repo",
        str(repo),
        "--profile",
        PROFILE,
        "--format",
        "json",
        "--decision-file",
        str(decision_file),
    ]
    if confirmation is not None:
        args.extend(["--confirm-plan", confirmation])
    stdout = io.StringIO()
    stderr = io.StringIO()
    with redirect_stdout(stdout), redirect_stderr(stderr):
        returncode = context_setup.main(args)
    return subprocess.CompletedProcess(args, returncode, stdout.getvalue(), stderr.getvalue())


if __name__ == "__main__":
    unittest.main()

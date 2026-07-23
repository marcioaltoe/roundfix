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
    def test_clean_profile_adoption_uses_the_strict_change_plan(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            repo = root / "repo"
            repo.mkdir()
            install_profile_skills(repo, PROFILE)
            (repo / "DESIGN.md").write_text(
                "# Design contract\n", encoding="utf-8"
            )
            http_path = repo / "docs/architecture/http-contract.json"
            http_path.parent.mkdir(parents=True)
            http_path.write_text('{"mode":"REST"}\n', encoding="utf-8")
            decision_file = root / "decisions.json"
            decision_file.write_text(
                json.dumps(
                    {
                        "schemaVersion": DECISION_SCHEMA,
                        "version": "0.0.1",
                        "decisions": self.profile_decisions(),
                    },
                    indent=2,
                )
                + "\n",
                encoding="utf-8",
            )
            before = snapshot_files(repo)

            blocked = run_apply(repo, decision_file)
            blocked_payload = json.loads(blocked.stdout)

            self.assertEqual(blocked.returncode, 1, blocked.stderr)
            self.assertFinding(blocked_payload, "capability.required.missing")
            self.assertNotIn("plan.confirmation.required", {
                item["code"] for item in blocked_payload["findings"]
            })
            self.assertEqual(snapshot_files(repo), before)

            self.write_capability_evidence(repo)
            preview = run_apply(repo, decision_file)
            preview_payload = json.loads(preview.stdout)
            before_apply = snapshot_files(repo)

            self.assertEqual(preview.returncode, 3, preview.stderr)
            self.assertFinding(preview_payload, "plan.confirmation.required")
            self.assertTrue(preview_payload["capabilities"])
            self.assertTrue(preview_payload["verification"])
            self.assertTrue(preview_payload["plannedOutputs"])
            self.assertEqual(
                preview_payload["setupSnapshot"]["schemaVersion"],
                "setup-context-driven/profile-snapshot/0.0.1",
            )
            self.assertEqual(snapshot_files(repo), before_apply)

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
                self.assertEqual(
                    target.read_bytes() if target.is_file() else None,
                    expected,
                )

            manifest = json.loads((repo / MANIFEST).read_text(encoding="utf-8"))
            self.assertEqual(
                manifest["schemaVersion"],
                "setup-context-driven/manifest/0.0.1",
            )
            self.assertEqual(manifest["version"], "0.0.1")
            self.assertEqual(manifest["generator"]["version"], "0.0.1")
            self.assertEqual(
                manifest["generator"]["baseline"],
                "baseline.standard-typescript-monorepo-0.0.1",
            )
            self.assertEqual(
                manifest["setupSnapshot"], preview_payload["setupSnapshot"]
            )
            self.assertNotIn("sourceBaseline", manifest)
            frontend = (repo / "docs/agents/frontend.md").read_text(
                encoding="utf-8"
            )
            backend = (repo / "docs/agents/backend.md").read_text(
                encoding="utf-8"
            )
            self.assertIn(
                "organize frontend feature code by domain system", frontend
            )
            self.assertIn("one public boundary", frontend)
            self.assertIn(
                "domain, application, and infrastructure layers", backend
            )
            self.assertIn("keep HTTP handlers thin", backend)
            self.assertIn("independent of HTTP", backend)
            self.assertIn("persistence implementation in infrastructure", backend)

            after_apply = snapshot_files(repo)
            audit = run_context_setup(
                "audit", "--repo", str(repo), "--format", "json"
            )
            reapplied = run_apply(repo, decision_file)

            self.assertEqual(audit.returncode, 0, audit.stderr)
            self.assertEqual(reapplied.returncode, 0, reapplied.stderr)
            self.assertEqual(json.loads(reapplied.stdout)["plannedChanges"], [])
            self.assertEqual(snapshot_files(repo), after_apply)

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

    def test_every_structural_entry_kind_has_an_individual_destination_and_preserves_repository_bytes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo, decision_file, _, proposed_rules = self.fixture(Path(temp_dir))
            nested_carrier = repo / "packages/api/CLAUDE.md"
            nested_carrier.parent.mkdir(parents=True)
            nested_carrier.write_text(
                "# Nested agent instructions\n\nPreserve this structural file.\n",
                encoding="utf-8",
            )
            (repo / "docs/agents/opaque.bin").write_bytes(b"\x00\xff\n")
            inventory = inventory_incompatible_source_baseline(
                repo, "baseline.pre-0.0.1"
            )
            document = json.loads(decision_file.read_text(encoding="utf-8"))
            document["readoption"]["sourceBaseline"] = {
                "id": inventory.baseline_id,
                "digest": inventory.digest,
            }
            typed_path = repo / "DESIGN.md"
            typed_bytes = typed_path.read_bytes()
            first_by_kind = {}
            for entry in inventory.entries:
                first_by_kind.setdefault(entry.kind, entry)
            self.assertEqual(
                set(first_by_kind),
                {"file", "managed-block", "manifest-record", "unmarked-span"},
            )

            dispositions = []
            for entry in inventory.entries:
                if entry == first_by_kind["file"]:
                    disposition = {
                        "entryId": entry.entry_id,
                        "entryDigest": entry.digest,
                        "classification": "operational-contract",
                        "disposition": "repository-document",
                        "destination": {
                            "documentType": "design-contract",
                            "path": "DESIGN.md",
                            "digest": hashlib.sha256(typed_bytes).hexdigest(),
                        },
                        "reason": "",
                    }
                elif entry == first_by_kind["managed-block"]:
                    disposition = {
                        "entryId": entry.entry_id,
                        "entryDigest": entry.digest,
                        "classification": "normative-clause",
                        "disposition": "managed-entry",
                        "destination": {"managedId": "root.context-workflow"},
                        "reason": "",
                    }
                elif entry == first_by_kind["unmarked-span"]:
                    disposition = {
                        "entryId": entry.entry_id,
                        "entryDigest": entry.digest,
                        "classification": "normative-clause",
                        "disposition": "repository-rules",
                        "destination": {
                            "documentType": "repository-rules",
                            "path": REPOSITORY_RULES.as_posix(),
                            "proposedBytes": base64.b64encode(proposed_rules).decode(
                                "ascii"
                            ),
                            "digest": hashlib.sha256(proposed_rules).hexdigest(),
                        },
                        "reason": "",
                    }
                else:
                    disposition = {
                        "entryId": entry.entry_id,
                        "entryDigest": entry.digest,
                        "classification": "non-governed",
                        "disposition": "rejected",
                        "destination": None,
                        "reason": (
                            "This structural entry is retained only as source evidence."
                        ),
                    }
                dispositions.append(disposition)
            document["readoption"]["dispositions"] = dispositions
            decision_file.write_text(
                json.dumps(document, indent=2) + "\n", encoding="utf-8"
            )

            preview = run_apply(repo, decision_file)
            preview_payload = json.loads(preview.stdout)

            self.assertEqual(preview.returncode, 3, preview.stderr)
            normalized = {
                item["entryId"]: item
                for item in preview_payload["decisionDocument"]["readoption"][
                    "dispositions"
                ]
            }
            for entry in inventory.entries:
                self.assertIn(entry.entry_id, normalized)
            self.assertEqual(
                {item["disposition"] for item in normalized.values()},
                {
                    "managed-entry",
                    "rejected",
                    "repository-document",
                    "repository-rules",
                },
            )

            applied = run_apply(
                repo, decision_file, preview_payload["planDigest"]
            )

            self.assertEqual(applied.returncode, 0, applied.stderr)
            self.assertEqual((repo / REPOSITORY_RULES).read_bytes(), proposed_rules)
            self.assertEqual(typed_path.read_bytes(), typed_bytes)

            repository_owned = proposed_rules + b"\nMaintainer-owned follow-up.\n"
            (repo / REPOSITORY_RULES).write_bytes(repository_owned)
            before_reapply = snapshot_files(repo)
            reapplied = run_apply(repo, decision_file)

            self.assertEqual(reapplied.returncode, 0, reapplied.stderr)
            self.assertEqual((repo / REPOSITORY_RULES).read_bytes(), repository_owned)
            self.assertEqual(typed_path.read_bytes(), typed_bytes)
            self.assertEqual(snapshot_files(repo), before_reapply)

    def test_rest_and_post_only_contracts_persist_typed_exceptions_and_exact_workspace_evidence(self):
        cases = (
            (
                "REST",
                [
                    {
                        "scope": "/webhooks/*",
                        "methods": ["POST"],
                        "owner": "payments",
                        "reason": "The provider owns webhook delivery semantics.",
                    }
                ],
            ),
            (
                "Post-only",
                [
                    {
                        "scope": "/health",
                        "methods": ["GET"],
                        "owner": "operations",
                        "reason": "Health probes require a safe read endpoint.",
                    }
                ],
            ),
        )
        restricted_environment = {
            "PATH": "/usr/bin:/bin",
            "PYTHONDONTWRITEBYTECODE": "1",
        }

        for mode, exceptions in cases:
            with self.subTest(mode=mode), tempfile.TemporaryDirectory() as temp_dir:
                repo, decision_file, _, _ = self.fixture(Path(temp_dir))
                package_path = repo / "package.json"
                package = json.loads(package_path.read_text(encoding="utf-8"))
                package["dependencies"]["inngest"] = "1"
                package_path.write_text(
                    json.dumps(package, sort_keys=True) + "\n", encoding="utf-8"
                )
                http_path = repo / "docs/architecture/http-contract.json"
                http_bytes = (
                    json.dumps({"mode": mode}, separators=(",", ":")) + "\n"
                ).encode("utf-8")
                http_path.write_bytes(http_bytes)
                document = json.loads(decision_file.read_text(encoding="utf-8"))
                http_decision = next(
                    item
                    for item in document["decisions"]
                    if item["id"] == "http.contract"
                )
                http_decision["value"] = {
                    "mode": mode,
                    "exceptions": exceptions,
                    "source": {
                        "path": "docs/architecture/http-contract.json",
                        "digest": hashlib.sha256(http_bytes).hexdigest(),
                    },
                }
                decision_file.write_text(
                    json.dumps(document, indent=2) + "\n", encoding="utf-8"
                )

                preview = run_apply(
                    repo, decision_file, extra_env=restricted_environment
                )
                payload = json.loads(preview.stdout)

                self.assertEqual(preview.returncode, 3, preview.stderr)
                self.assertEqual(
                    payload["setupSnapshot"]["httpContract"],
                    http_decision["value"],
                )
                capabilities = {
                    item["id"]: item for item in payload["capabilities"]
                }
                for workspace in ("backend", "frontend"):
                    capability = capabilities[f"capability.workspace.{workspace}"]
                    self.assertEqual(capability["status"], "satisfied")
                    self.assertEqual(
                        capability["evidence"][0]["sourcePath"],
                        f"packages/{workspace}/package.json",
                    )
                self.assertEqual(
                    capabilities["capability.optional.inngest"]["status"],
                    "satisfied",
                )
                self.assertEqual(
                    capabilities["capability.optional.docker"]["status"],
                    "missing",
                )

                applied = run_apply(
                    repo,
                    decision_file,
                    payload["planDigest"],
                    extra_env=restricted_environment,
                )
                self.assertEqual(applied.returncode, 0, applied.stderr)
                after_apply = snapshot_files(repo)
                manifest = json.loads(
                    (repo / MANIFEST).read_text(encoding="utf-8")
                )
                self.assertEqual(
                    manifest["setupSnapshot"]["httpContract"],
                    http_decision["value"],
                )

                reapplied = run_apply(
                    repo, decision_file, extra_env=restricted_environment
                )

                self.assertEqual(reapplied.returncode, 0, reapplied.stderr)
                self.assertEqual(snapshot_files(repo), after_apply)

    def test_missing_recommended_capabilities_warn_and_still_allow_confirmed_apply(self):
        restricted_environment = {
            "PATH": "/usr/bin:/bin",
            "PYTHONDONTWRITEBYTECODE": "1",
        }
        with tempfile.TemporaryDirectory() as temp_dir:
            repo, decision_file, _, _ = self.fixture(Path(temp_dir))
            firecrawl = repo / ".agents/skills/firecrawl/SKILL.md"
            firecrawl.unlink()
            before = snapshot_files(repo)

            preview = run_apply(
                repo, decision_file, extra_env=restricted_environment
            )
            payload = json.loads(preview.stdout)

            self.assertEqual(preview.returncode, 3, preview.stderr)
            self.assertEqual(snapshot_files(repo), before)
            warnings = {
                item["managedId"]
                for item in payload["findings"]
                if item["code"] == "capability.recommended.missing"
            }
            self.assertEqual(
                warnings,
                {"capability.firecrawl", "capability.rg", "capability.rtk"},
            )

            applied = run_apply(
                repo,
                decision_file,
                payload["planDigest"],
                extra_env=restricted_environment,
            )

            self.assertEqual(applied.returncode, 0, applied.stderr)
            applied_warnings = {
                item["managedId"]
                for item in json.loads(applied.stdout)["findings"]
                if item["code"] == "capability.recommended.missing"
            }
            self.assertEqual(applied_warnings, warnings)

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
        http_path.parent.mkdir(parents=True, exist_ok=True)
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
        return {
            "schemaVersion": DECISION_SCHEMA,
            "version": "0.0.1",
            "decisions": self.profile_decisions(),
            "readoption": {
                "sourceBaseline": {
                    "id": inventory.baseline_id,
                    "digest": inventory.digest,
                },
                "dispositions": dispositions,
            },
        }

    @staticmethod
    def profile_decisions():
        http_path = Path("docs/architecture/http-contract.json")
        source_bytes = b'{"mode":"REST"}\n'
        return [
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

    @staticmethod
    def assertFinding(payload, code):
        if not any(item["code"] == code for item in payload["findings"]):
            raise AssertionError(f"missing finding {code}: {payload['findings']}")


def run_apply(repo, decision_file, confirmation=None, extra_env=None):
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
    return run_context_setup(*args, extra_env=extra_env)


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

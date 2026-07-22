"""Repository-authored instruction delegation audit coverage.

Suite: bounded baseline-floor delegation findings
Invariant: delegation discovery is read-only, deterministic, and non-blocking.
Boundary IN: repository AGENTS.md and CLAUDE.md documents plus active profile coverage.
Boundary OUT: Change Plan mutation authority, process execution, and network access.
"""

import json
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_assets import load_asset_catalog  # noqa: E402
from context_setup import (  # noqa: E402
    DelegationScanLimits,
    delegation_findings,
    exit_code_for,
    plan_apply,
)
from test_audit import (  # noqa: E402
    run_audit,
    run_context_setup,
    snapshot_files,
    write_compliant_repository,
)


class DelegationAuditTests(unittest.TestCase):
    def test_uncovered_nested_delegation_emits_one_informational_finding(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            instructions = repo / "packages" / "web" / "AGENTS.md"
            instructions.parent.mkdir(parents=True)
            instructions.write_text(
                "Follow the frontend guidance in docs/agents/frontend.md.\n\n"
                "Follow the frontend guidance in docs/agents/frontend.md.\n",
                encoding="utf-8",
            )
            before = snapshot_files(repo)

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 0, result.stderr)
            payload = json.loads(result.stdout)
            matches = [
                item
                for item in payload["findings"]
                if item["code"] == "delegation.baseline-floor"
            ]
            self.assertEqual(len(matches), 1, payload)
            self.assertEqual(matches[0]["severity"], "info")
            self.assertEqual(matches[0]["path"], "packages/web/AGENTS.md")
            self.assertEqual(matches[0]["managedId"], "coverage.frontend")
            self.assertIn("floor", matches[0]["action"].casefold())
            self.assertEqual(snapshot_files(repo), before)

    def test_alias_without_delegation_signal_and_covered_category_emit_no_floor(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            no_signal = repo / "packages" / "web" / "AGENTS.md"
            no_signal.parent.mkdir(parents=True)
            no_signal.write_text(
                "Frontend guidance documents the user interface.\n",
                encoding="utf-8",
            )
            covered = repo / "crates" / "cli" / "CLAUDE.md"
            covered.parent.mkdir(parents=True)
            covered.write_text(
                "Follow the Rust guidance in docs/agents/rust.md.\n",
                encoding="utf-8",
            )

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 0, result.stderr)
            payload = json.loads(result.stdout)
            self.assertFalse(
                any(
                    item["code"] == "delegation.baseline-floor"
                    for item in payload["findings"]
                ),
                payload,
            )

    def test_managed_ignored_and_symlinked_instructions_are_excluded(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            repo = root / "repo"
            write_compliant_repository(repo, "rust-cli")
            managed = repo / "CLAUDE.md"
            managed.write_text(
                "<!-- setup-context-driven:begin id=root.test version=1 -->\n"
                "Follow the frontend guidance in docs/agents/frontend.md.\n"
                "<!-- setup-context-driven:end id=root.test -->\n",
                encoding="utf-8",
            )
            for relative in [
                ".git/AGENTS.md",
                "node_modules/tool/CLAUDE.md",
                "vendor/tool/AGENTS.md",
                ".agents/skills/setup-context-driven/AGENTS.md",
                "skills/setup-context-driven/CLAUDE.md",
            ]:
                target = repo / relative
                target.parent.mkdir(parents=True, exist_ok=True)
                target.write_text(
                    "Follow the frontend guidance in docs/agents/frontend.md.\n",
                    encoding="utf-8",
                )
            outside = root / "outside-AGENTS.md"
            outside.write_text(
                "Follow the frontend guidance in docs/agents/frontend.md.\n",
                encoding="utf-8",
            )
            (repo / "packages" / "linked").mkdir(parents=True)
            (repo / "packages" / "linked" / "AGENTS.md").symlink_to(outside)

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 0, result.stderr)
            payload = json.loads(result.stdout)
            self.assertFalse(
                any(
                    item["code"] == "delegation.baseline-floor"
                    for item in payload["findings"]
                ),
                payload,
            )

    def test_findings_sort_by_document_and_category_and_collapse_duplicates(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            documents = {
                "zeta/CLAUDE.md": (
                    "Follow the Go guidance in docs/agents/go.md.\n\n"
                    "Follow the frontend guidance in docs/agents/frontend.md.\n"
                ),
                "alpha/AGENTS.md": (
                    "Follow the frontend guidance in docs/agents/frontend.md.\n\n"
                    "Follow the frontend guidance in docs/agents/frontend.md.\n"
                ),
            }
            for relative, content in documents.items():
                target = repo / relative
                target.parent.mkdir(parents=True)
                target.write_text(content, encoding="utf-8")

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 0, result.stderr)
            payload = json.loads(result.stdout)
            actual = [
                (item["path"], item["managedId"])
                for item in payload["findings"]
                if item["code"] == "delegation.baseline-floor"
            ]
            self.assertEqual(
                actual,
                [
                    ("alpha/AGENTS.md", "coverage.frontend"),
                    ("zeta/CLAUDE.md", "coverage.frontend"),
                    ("zeta/CLAUDE.md", "coverage.go"),
                ],
            )

    def test_scan_limits_stop_without_partial_floor_findings_or_writes(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        cases = [
            (DelegationScanLimits(max_files=1, max_bytes=4096), ["a/AGENTS.md", "b/AGENTS.md"]),
            (DelegationScanLimits(max_files=4, max_bytes=32), ["a/AGENTS.md"]),
        ]
        for limits, relative_paths in cases:
            with self.subTest(limits=limits):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    for relative in relative_paths:
                        target = repo / relative
                        target.parent.mkdir(parents=True, exist_ok=True)
                        target.write_text(
                            "Follow the frontend guidance in docs/agents/frontend.md.\n",
                            encoding="utf-8",
                        )
                    before = snapshot_files(repo)

                    findings = delegation_findings(
                        repo,
                        catalog,
                        "rust-cli",
                        active_modules=("core", "rust", "cli-surface"),
                        limits=limits,
                    )

                    self.assertEqual(
                        [item.code for item in findings],
                        ["delegation.scan-limit"],
                    )
                    self.assertEqual(findings[0].severity, "info")
                    self.assertEqual(snapshot_files(repo), before)

    def test_floor_finding_does_not_change_plan_authority_or_clean_exit(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            catalog = load_asset_catalog(SKILL_ROOT)
            clean_result, clean_invalid, clean_plan = plan_apply(
                repo=repo,
                catalog=catalog,
                profile_override="rust-cli",
                decision_args=[],
            )
            instructions = repo / "packages" / "web" / "AGENTS.md"
            instructions.parent.mkdir(parents=True)
            instructions.write_text(
                "Follow the frontend guidance in docs/agents/frontend.md.\n",
                encoding="utf-8",
            )
            before = snapshot_files(repo)

            floor_result, floor_invalid, floor_plan = plan_apply(
                repo=repo,
                catalog=catalog,
                profile_override="rust-cli",
                decision_args=[],
            )
            applied = run_context_setup(
                "apply",
                "--repo",
                str(repo),
                "--format",
                "json",
            )

            self.assertFalse(clean_invalid)
            self.assertFalse(floor_invalid)
            self.assertEqual(exit_code_for(clean_result, clean_invalid), 0)
            self.assertEqual(exit_code_for(floor_result, floor_invalid), 0)
            self.assertEqual(clean_plan.digest, floor_plan.digest)
            self.assertEqual(clean_result.plan_digest, floor_result.plan_digest)
            self.assertEqual(applied.returncode, 0, applied.stderr)
            payload = json.loads(applied.stdout)
            self.assertTrue(
                any(
                    item["code"] == "delegation.baseline-floor"
                    for item in payload["findings"]
                ),
                payload,
            )
            self.assertEqual(snapshot_files(repo), before)


if __name__ == "__main__":
    unittest.main()

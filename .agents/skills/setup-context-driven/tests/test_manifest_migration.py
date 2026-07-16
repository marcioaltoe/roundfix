"""Manifest migration flows for setup-context-driven.

Suite: spec 0030 manifest migration
Invariant: schema-v1 manifests migrate only proven setup-owned content while preserving compatible decisions.
Boundary IN: context_setup.py CLI, Setup Manifest v1 data, managed markers, temporary repository files.
Boundary OUT: manual QA harness storage and Makefile orchestration.
"""

import json
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_setup import managed_block  # noqa: E402
from test_apply import run_apply, run_audit  # noqa: E402
from test_audit import install_profile_skills, snapshot_files  # noqa: E402


LEGACY_CONFIRMED_AT = "2026-07-01"
ROOT_PREFIX = "repository root before\n"
ROOT_BETWEEN = "repository root between\n"
ROOT_SUFFIX = "repository root after\n"
GUIDE_PREFIX = "repository guide before\n"
GUIDE_SUFFIX = "repository guide after\n"


class ManifestMigrationTests(unittest.TestCase):
    def test_spec0030_manifest_migrates_answers_inventory_and_owned_blocks(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = legacy_decisions(
                spec_scaffold=False,
                domain_layout="multi-context",
                triage_external=False,
                autonomous=False,
                secondbrain=False,
            )
            write_legacy_spec0030_repository(repo, decisions)

            migrated = run_apply(repo, "typescript-bun-monorepo", [])
            install_profile_skills(repo, "typescript-bun-monorepo")
            audited = run_audit(repo)
            after_audit = snapshot_files(repo)
            repeated = run_apply(repo, "typescript-bun-monorepo", [])

            self.assertEqual(migrated.returncode, 0, migrated.stderr)
            self.assertNoFinding(migrated, "decision.required")
            self.assertEqual(audited.returncode, 0, audited.stderr)
            self.assertEqual(json.loads(audited.stdout)["findings"], [])
            self.assertEqual(repeated.returncode, 0, repeated.stderr)
            self.assertEqual(snapshot_files(repo), after_audit)

            manifest = read_manifest(repo)
            self.assertEqual(manifest["schemaVersion"], 1)
            for decision_id, decision in decisions.items():
                self.assertEqual(manifest["decisions"][decision_id], decision)
            self.assertNotIn("spec-workflow", manifest["modules"])
            self.assertNotIn("autonomous-work", manifest["modules"])
            self.assertNotIn("external-triage", manifest["modules"])
            self.assertNotIn("secondbrain", manifest["modules"])

            artifacts = {artifact["id"]: artifact for artifact in manifest["managedArtifacts"]}
            self.assertNotIn("guide.autonomous-work", artifacts)
            self.assertNotIn("root.autonomous-work", artifacts)
            self.assertEqual(
                artifacts["guide.domain"]["template"],
                "template.guide.domain.multi-context",
            )
            self.assertEqual(
                (repo / "docs" / "agents" / "autonomous-work.md").read_text(encoding="utf-8"),
                GUIDE_PREFIX + GUIDE_SUFFIX,
            )
            root_content = (repo / "AGENTS.md").read_text(encoding="utf-8")
            self.assertIn(ROOT_PREFIX, root_content)
            self.assertIn(ROOT_BETWEEN, root_content)
            self.assertIn(ROOT_SUFFIX, root_content)
            self.assertLess(root_content.index(ROOT_PREFIX), root_content.index(ROOT_BETWEEN))
            self.assertLess(root_content.index(ROOT_BETWEEN), root_content.index(ROOT_SUFFIX))

    def test_enabled_capability_routes_only_missing_dependent_decision(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = legacy_decisions(autonomous=True)
            decisions.pop("verification.gate")
            write_legacy_spec0030_repository(repo, decisions)
            before = snapshot_files(repo)

            blocked = run_apply(repo, "rust-cli", [])
            blocked_payload = json.loads(blocked.stdout)
            required = [
                finding["managedId"]
                for finding in blocked_payload["findings"]
                if finding["code"] == "decision.required"
            ]
            self.assertEqual(blocked.returncode, 3, blocked.stderr)
            self.assertEqual(required, ["verification.gate"], blocked_payload)
            self.assertEqual(snapshot_files(repo), before)

            answered = run_apply(repo, "rust-cli", ["verification.gate=make verify"])
            repeated = run_apply(repo, "rust-cli", [])

            self.assertEqual(answered.returncode, 0, answered.stderr)
            self.assertEqual(repeated.returncode, 0, repeated.stderr)
            manifest = read_manifest(repo)
            self.assertEqual(manifest["decisions"]["verification.gate"]["value"], "make verify")
            self.assertNoFinding(repeated, "decision.required")

    def test_ambiguous_legacy_ownership_blocks_without_partial_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = legacy_decisions(autonomous=False)
            write_legacy_spec0030_repository(
                repo,
                decisions,
                ambiguous_autonomous_guide=True,
            )
            before = snapshot_files(repo)

            result = run_apply(repo, "rust-cli", [])

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "managed.ownership.ambiguous", "error")
            self.assertEqual(snapshot_files(repo), before)

    def assertFinding(self, result, code, severity):
        matches = [
            finding
            for finding in json.loads(result.stdout)["findings"]
            if finding["code"] == code
        ]
        self.assertGreater(len(matches), 0, result.stdout)
        self.assertEqual(matches[0]["severity"], severity)

    def assertNoFinding(self, result, code):
        matches = [
            finding
            for finding in json.loads(result.stdout)["findings"]
            if finding["code"] == code
        ]
        self.assertEqual(matches, [], result.stdout)


def legacy_decisions(
    *,
    spec_scaffold=True,
    domain_layout="single-context",
    triage_external=False,
    autonomous=True,
    secondbrain=False,
) -> dict[str, dict]:
    values = {
        "spec.scaffold": spec_scaffold,
        "domain.layout": domain_layout,
        "triage.external": triage_external,
        "autonomous.enabled": autonomous,
        "runtime.backend": "codex legacy-backend xhigh",
        "runtime.design": "claude legacy-design xhigh",
        "verification.gate": "make legacy-verify",
        "language.generated": "English",
        "secondbrain.enabled": secondbrain,
    }
    return {
        decision_id: {"value": value, "confirmedAt": LEGACY_CONFIRMED_AT}
        for decision_id, value in values.items()
    }


def write_legacy_spec0030_repository(
    repo: Path,
    decisions: dict[str, dict],
    *,
    ambiguous_autonomous_guide=False,
) -> None:
    repo.mkdir(parents=True, exist_ok=True)
    docs_agents = repo / "docs" / "agents"
    docs_agents.mkdir(parents=True, exist_ok=True)
    root_content = (
        ROOT_PREFIX
        + managed_block(
            "root.context-workflow",
            1,
            "Legacy mixed context workflow with docs/specs/<feature-slug> routing.",
        )
        + ROOT_BETWEEN
        + managed_block(
            "root.autonomous-work",
            1,
            "Legacy autonomous pointer.",
        )
        + ROOT_SUFFIX
    )
    (repo / "AGENTS.md").write_text(root_content, encoding="utf-8")
    (docs_agents / "docs-layout.md").write_text(
        GUIDE_PREFIX
        + managed_block(
            "guide.docs-layout",
            1,
            "Legacy docs layout with Spec-only docs/specs/<feature-slug> text.",
        )
        + GUIDE_SUFFIX,
        encoding="utf-8",
    )
    (docs_agents / "domain.md").write_text(
        managed_block("guide.domain", 1, "Legacy generic domain layout."),
        encoding="utf-8",
    )
    autonomous_content = (
        "repository-authored autonomous notes\n"
        if ambiguous_autonomous_guide
        else GUIDE_PREFIX
        + managed_block("guide.autonomous-work", 1, "Legacy autonomous guide.")
        + GUIDE_SUFFIX
    )
    (docs_agents / "autonomous-work.md").write_text(
        autonomous_content,
        encoding="utf-8",
    )
    manifest = {
        "schemaVersion": 1,
        "generator": {"skill": "setup-context-driven", "version": 1},
        "profile": "typescript-bun-monorepo",
        "modules": [
            "core",
            "context-workflow",
            "typescript",
            "bun",
            "monorepo",
            "backend",
            "frontend",
            "autonomous-work",
        ],
        "decisions": decisions,
        "managedArtifacts": [
            legacy_artifact(
                "root.context-workflow",
                "AGENTS.md",
                "root-block",
                "context-workflow",
                "template.root.context-workflow",
            ),
            legacy_artifact(
                "root.autonomous-work",
                "AGENTS.md",
                "root-block",
                "autonomous-work",
                "template.root.autonomous-work",
            ),
            legacy_artifact(
                "guide.docs-layout",
                "docs/agents/docs-layout.md",
                "guide",
                "context-workflow",
                "template.guide.docs-layout",
            ),
            legacy_artifact(
                "guide.domain",
                "docs/agents/domain.md",
                "guide",
                "context-workflow",
                "template.guide.domain",
            ),
            legacy_artifact(
                "guide.autonomous-work",
                "docs/agents/autonomous-work.md",
                "guide",
                "autonomous-work",
                "template.guide.autonomous-work",
            ),
        ],
        "localSkills": [],
    }
    (docs_agents / "setup-context.json").write_text(
        json.dumps(manifest, indent=2) + "\n",
        encoding="utf-8",
    )


def legacy_artifact(
    managed_id: str,
    path: str,
    kind: str,
    module_id: str,
    template_id: str,
) -> dict:
    return {
        "id": managed_id,
        "path": path,
        "kind": kind,
        "module": module_id,
        "template": template_id,
        "version": 1,
        "digest": "0" * 64,
    }


def read_manifest(repo: Path) -> dict:
    return json.loads(
        (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
    )


if __name__ == "__main__":
    unittest.main()

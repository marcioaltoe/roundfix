# Suite: Repository Capability evidence
# Invariant: capability readiness uses sufficient bounded local evidence and never mutates the repository.
# Boundary IN: immutable capability contracts, local evidence adapters, evaluation, diagnostics, and rendering.
# Boundary OUT: profile binding, audit/Decision Plan integration, installation, network access, and repository scripts.

import os
import socket
import subprocess
import sys
import tempfile
import unittest
import urllib.request
from dataclasses import FrozenInstanceError
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))

from context_capabilities import (  # noqa: E402
    CapabilityEvidence,
    CapabilityStatus,
    EvidenceKind,
    EvidenceStatus,
    EvidenceStrength,
    RepositoryCapability,
    RequirementStrength,
    UNIVERSAL_CAPABILITIES,
    evaluate_capabilities,
    evaluate_repository_capabilities,
    render_capability_guidance,
    render_capability_json,
)


class RepositoryCapabilityTests(unittest.TestCase):
    def test_missing_required_research_capabilities_block_with_next_actions(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            evaluation = evaluate_repository_capabilities(
                Path(temp_dir), executable_lookup=lambda _name: None
            )

        required = {
            outcome.capability.capability_id: outcome
            for outcome in evaluation.outcomes
            if outcome.capability.strength is RequirementStrength.REQUIRED
        }
        self.assertFalse(evaluation.ready)
        self.assertEqual(set(required), {"capability.context7", "capability.exa"})
        for capability_id, outcome in required.items():
            with self.subTest(capability_id=capability_id):
                self.assertTrue(outcome.blocking)
                self.assertEqual(outcome.status, CapabilityStatus.MISSING)
                self.assertEqual(outcome.diagnostic.code, "capability.required.missing")
                self.assertIn(outcome.capability.title, outcome.diagnostic.message)
                self.assertIn("Repository Skill Set", outcome.diagnostic.next_action)

    def test_missing_recommended_capabilities_warn_without_blocking(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repository = Path(temp_dir)
            self._install_skill(repository, "context7")
            self._install_skill(repository, "exa-web-search")

            evaluation = evaluate_repository_capabilities(
                repository, executable_lookup=lambda _name: None
            )

        warnings = {
            outcome.capability.capability_id: outcome
            for outcome in evaluation.outcomes
            if outcome.status is not CapabilityStatus.SATISFIED
        }
        self.assertTrue(evaluation.ready)
        self.assertEqual(
            set(warnings),
            {"capability.firecrawl", "capability.rg", "capability.rtk"},
        )
        for capability_id, outcome in warnings.items():
            with self.subTest(capability_id=capability_id):
                self.assertFalse(outcome.blocking)
                self.assertEqual(
                    outcome.diagnostic.code, "capability.recommended.missing"
                )
                self.assertTrue(outcome.explanation)

    def test_stronger_compatible_evidence_satisfies_a_weaker_requirement(self):
        capability = self._capability(
            "capability.example",
            evidence_kind=EvidenceKind.DECLARED_FILE,
            minimum_evidence=EvidenceStrength.DECLARED,
        )
        evidence = CapabilityEvidence(
            capability_id=capability.capability_id,
            status=EvidenceStatus.PRESENT,
            version=None,
            source_path=Path("contract.json"),
            source_digest="a" * 64,
            evidence_kind=EvidenceKind.DECLARED_FILE,
            strength=EvidenceStrength.VERIFIED,
        )

        outcome = evaluate_capabilities((capability,), (evidence,)).outcomes[0]

        self.assertEqual(outcome.status, CapabilityStatus.SATISFIED)

    def test_insufficient_or_incompatible_evidence_does_not_satisfy_requirement(self):
        capability = self._capability(
            "capability.example",
            evidence_kind=EvidenceKind.INSTALLED_SKILL,
            minimum_evidence=EvidenceStrength.VERIFIED,
        )
        insufficient = CapabilityEvidence(
            capability_id=capability.capability_id,
            status=EvidenceStatus.PRESENT,
            version=None,
            source_path=Path(".agents/skills/example/SKILL.md"),
            source_digest="b" * 64,
            evidence_kind=EvidenceKind.INSTALLED_SKILL,
            strength=EvidenceStrength.DECLARED,
        )
        incompatible = CapabilityEvidence(
            capability_id=capability.capability_id,
            status=EvidenceStatus.PRESENT,
            version=None,
            source_path=Path("bin/example"),
            source_digest=None,
            evidence_kind=EvidenceKind.EXECUTABLE,
            strength=EvidenceStrength.VERIFIED,
        )

        insufficient_outcome = evaluate_capabilities(
            (capability,), (insufficient,)
        ).outcomes[0]
        incompatible_outcome = evaluate_capabilities(
            (capability,), (incompatible,)
        ).outcomes[0]

        self.assertEqual(insufficient_outcome.status, CapabilityStatus.INSUFFICIENT)
        self.assertEqual(
            insufficient_outcome.diagnostic.code,
            "capability.evidence.insufficient",
        )
        self.assertEqual(incompatible_outcome.status, CapabilityStatus.MISSING)

    def test_file_executable_and_installed_skill_adapters_use_local_evidence(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repository = Path(temp_dir)
            contract = repository / "docs" / "tool.json"
            contract.parent.mkdir(parents=True)
            contract.write_text('{"tool":"declared"}\n', encoding="utf-8")
            self._install_skill(repository, "example-skill")
            executable = repository / "bin" / "example"
            executable.parent.mkdir()
            executable.write_text("local executable evidence\n", encoding="utf-8")
            capabilities = (
                self._capability(
                    "capability.file",
                    evidence_kind=EvidenceKind.DECLARED_FILE,
                    minimum_evidence=EvidenceStrength.DECLARED,
                    probe={"paths": ("docs/tool.json",), "contains": '"tool"'},
                ),
                self._capability(
                    "capability.executable",
                    evidence_kind=EvidenceKind.EXECUTABLE,
                    minimum_evidence=EvidenceStrength.DISCOVERED,
                    probe={"executable": "example"},
                ),
                self._capability(
                    "capability.skill",
                    evidence_kind=EvidenceKind.INSTALLED_SKILL,
                    minimum_evidence=EvidenceStrength.VERIFIED,
                    probe={"skill": "example-skill"},
                ),
            )

            evaluation = evaluate_repository_capabilities(
                repository,
                capabilities,
                executable_lookup=lambda name: str(executable)
                if name == "example"
                else None,
            )

        self.assertTrue(evaluation.ready)
        self.assertEqual(
            tuple(outcome.status for outcome in evaluation.outcomes),
            (CapabilityStatus.SATISFIED,) * 3,
        )
        self.assertTrue(
            all(outcome.evidence[0].source_path is not None for outcome in evaluation.outcomes)
        )

    def test_equivalent_evidence_renders_byte_identical_ordered_output(self):
        capabilities = tuple(reversed(UNIVERSAL_CAPABILITIES))
        evidence = tuple(
            CapabilityEvidence(
                capability_id=capability.capability_id,
                status=EvidenceStatus.PRESENT,
                version=None,
                source_path=Path(f"evidence/{capability.capability_id}"),
                source_digest=None,
                evidence_kind=capability.evidence_kind,
                strength=EvidenceStrength.VERIFIED,
            )
            for capability in capabilities
        )

        first = evaluate_capabilities(capabilities, evidence)
        second = evaluate_capabilities(
            tuple(reversed(capabilities)), tuple(reversed(evidence))
        )

        self.assertEqual(render_capability_json(first), render_capability_json(second))
        self.assertEqual(
            render_capability_guidance(first), render_capability_guidance(second)
        )
        self.assertEqual(
            tuple(outcome.capability.capability_id for outcome in first.outcomes),
            tuple(sorted(capability.capability_id for capability in capabilities)),
        )

    def test_capability_evaluation_has_no_write_install_network_or_script_side_effects(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repository = Path(temp_dir)
            self._install_skill(repository, "context7")
            self._install_skill(repository, "exa-web-search")
            before = self._tree_bytes(repository)

            with (
                mock.patch.object(
                    Path, "write_text", side_effect=AssertionError("write attempted")
                ),
                mock.patch.object(
                    Path, "write_bytes", side_effect=AssertionError("write attempted")
                ),
                mock.patch.object(
                    Path, "mkdir", side_effect=AssertionError("mkdir attempted")
                ),
                mock.patch.object(
                    Path, "unlink", side_effect=AssertionError("unlink attempted")
                ),
                mock.patch.object(
                    subprocess, "run", side_effect=AssertionError("script attempted")
                ),
                mock.patch.object(
                    subprocess, "Popen", side_effect=AssertionError("script attempted")
                ),
                mock.patch.object(
                    urllib.request,
                    "urlopen",
                    side_effect=AssertionError("network attempted"),
                ),
                mock.patch.object(
                    socket,
                    "create_connection",
                    side_effect=AssertionError("network attempted"),
                ),
                mock.patch.object(
                    os, "system", side_effect=AssertionError("install attempted")
                ),
            ):
                evaluate_repository_capabilities(
                    repository, executable_lookup=lambda _name: None
                )

            self.assertEqual(self._tree_bytes(repository), before)

    def test_rendered_guidance_keeps_external_research_subordinate_to_local_search(self):
        evaluation = evaluate_capabilities(UNIVERSAL_CAPABILITIES, ())

        guidance = render_capability_guidance(evaluation)

        self.assertIn("Search local repository code with `rg` first", guidance)
        self.assertIn("cannot replace local code search", guidance)
        self.assertLess(guidance.index("local repository code"), guidance.index("Context7"))

    def test_capability_contracts_and_outcomes_are_immutable(self):
        capability = UNIVERSAL_CAPABILITIES[0]
        outcome = evaluate_capabilities((capability,), ()).outcomes[0]

        with self.assertRaises(FrozenInstanceError):
            capability.title = "changed"
        with self.assertRaises(TypeError):
            capability.probe["skill"] = "changed"
        with self.assertRaises(FrozenInstanceError):
            outcome.blocking = False

    def _capability(
        self,
        capability_id,
        *,
        evidence_kind,
        minimum_evidence,
        probe=None,
    ):
        if probe is None:
            probe = {"skill": "example"}
        return RepositoryCapability(
            capability_id=capability_id,
            title=capability_id,
            strength=RequirementStrength.REQUIRED,
            evidence_kind=evidence_kind,
            minimum_evidence=minimum_evidence,
            probe=probe,
            explanation="Required for this test.",
            next_action="Add compatible local evidence and rerun evaluation.",
        )

    def _install_skill(self, repository, name):
        path = repository / ".agents" / "skills" / name / "SKILL.md"
        path.parent.mkdir(parents=True)
        path.write_text(f"---\nname: {name}\n---\n", encoding="utf-8")

    def _tree_bytes(self, root):
        return {
            path.relative_to(root).as_posix(): path.read_bytes()
            for path in sorted(root.rglob("*"))
            if path.is_file()
        }


if __name__ == "__main__":
    unittest.main()

# Suite: governed Context-Driven Baseline corpus
# Invariant: every governed 0.0.1 entry is indexed once, portable, and structurally complete.
# Boundary IN: production Source Baseline assets, prior-clause accounting, and corpus policy.
# Boundary OUT: profile composition, Skill Activation rendering, Readoption planning, and apply.

import json
import shutil
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
BASELINE_ID = "baseline.standard-typescript-monorepo-0.0.1"
BASELINE_ROOT = SKILL_ROOT / "assets" / "source-baselines" / BASELINE_ID
sys.path.insert(0, str(SKILL_ROOT / "scripts"))

from context_baseline import (  # noqa: E402
    SOURCE_BASELINE_VERSION,
    SourceBaselineValidationError,
    load_source_baseline,
    render_source_baseline_carriers,
)


class GovernedCorpusTests(unittest.TestCase):
    def test_every_governed_entry_has_one_indexed_identity_and_carrier(self):
        baseline = load_source_baseline(SKILL_ROOT, BASELINE_ID)
        entry_ids = tuple(entry.entry_id for entry in baseline.entries)

        self.assertEqual(entry_ids, baseline.index_entry.entry_ids)
        self.assertEqual(len(entry_ids), len(set(entry_ids)))
        self.assertGreaterEqual(len(entry_ids), 50)
        self.assertTrue(all(entry.carrier.parts for entry in baseline.entries))
        self.assertEqual(
            {entry.kind for entry in baseline.entries},
            {"normative-clause", "recommendation", "operational-contract"},
        )

    def test_required_guides_and_operational_contracts_are_complete(self):
        baseline = load_source_baseline(SKILL_ROOT, BASELINE_ID)
        entries = {entry.entry_id: entry for entry in baseline.entries}
        content = self._entry_content(baseline)

        required_carriers = {
            Path("AGENTS.md"),
            Path("docs/agents/agent-instructions.md"),
            Path("docs/agents/autonomous-work.md"),
            Path("docs/agents/backend.md"),
            Path("docs/agents/frontend.md"),
            Path("docs/agents/domain.md"),
            Path("docs/agents/docs-layout.md"),
            Path("docs/agents/spec-routing.md"),
            Path("docs/agents/issue-tracker.md"),
            Path("docs/agents/secondbrain.md"),
            Path("docs/agents/typescript-bun.md"),
        }
        self.assertTrue(required_carriers <= {entry.carrier for entry in baseline.entries})

        required_contracts = {
            "contract.root.delegation-index": "template",
            "contract.docs-layout.directory-matrix": "decision-matrix",
            "contract.findings.template": "template",
            "contract.findings.lifecycle": "lifecycle",
            "contract.spec.route-matrix": "decision-matrix",
            "contract.spec.task-ownership": "protocol",
            "contract.autonomous.supervisor-runtime": "protocol",
            "contract.research.sequence": "ordered-procedure",
            "contract.secondbrain.protocol": "protocol",
        }
        for entry_id, structure in required_contracts.items():
            with self.subTest(entry_id=entry_id):
                entry = entries[entry_id]
                self.assertEqual(entry.kind, "operational-contract")
                self.assertEqual(entry.structure, structure)
                self.assertGreater(len(content[entry_id].splitlines()), 3)

        findings = content["contract.findings.template"]
        for required in (
            "status: pending # pending | partial | deferred | done",
            "created_at: YYYY-MM-DD",
            "updated_at: YYYY-MM-DD",
            "Symptom / evidence",
            "Root cause",
            "Action / suggestion",
        ):
            self.assertIn(required, findings)

        rendered = render_source_baseline_carriers(SKILL_ROOT, BASELINE_ID)
        self.assertTrue(required_carriers <= set(rendered))
        self.assertIn(
            b"status: pending # pending | partial | deferred | done",
            rendered[Path("docs/agents/docs-layout.md")],
        )
        self.assertIn(
            b"The Supervisor MUST NOT write feature code or tests",
            rendered[Path("docs/agents/autonomous-work.md")],
        )

    def test_root_contract_is_compact_and_delegates_to_complete_guides(self):
        baseline = load_source_baseline(SKILL_ROOT, BASELINE_ID)
        content = self._entry_content(baseline)
        root = content["contract.root.delegation-index"]

        self.assertLess(len(root.encode("utf-8")), 1_500)
        self.assertIn("docs/agents/agent-instructions.md", root)
        self.assertIn("docs/agents/spec-routing.md", root)
        self.assertIn("docs/agents/docs-layout.md", root)
        self.assertNotIn("status: pending # pending | partial | deferred | done", root)

    def test_prior_clause_accounting_is_individual_and_targets_current_entries(self):
        baseline = load_source_baseline(SKILL_ROOT, BASELINE_ID)
        current_ids = {entry.entry_id for entry in baseline.entries}
        accounting = self._read_json(BASELINE_ROOT / "accounting.json")

        self.assertEqual(
            accounting["schemaVersion"],
            "setup-context-driven/source-accounting/0.0.1",
        )
        self.assertEqual(accounting["version"], SOURCE_BASELINE_VERSION)
        records = accounting["entries"]
        source_ids = [record["sourceEntry"] for record in records]
        self.assertEqual(len(source_ids), 51)
        self.assertEqual(len(source_ids), len(set(source_ids)))
        for record in records:
            with self.subTest(source_entry=record["sourceEntry"]):
                self.assertIn(record["disposition"], {"retained", "moved", "replaced", "rejected"})
                self.assertTrue(record["reason"].strip())
                if record["disposition"] == "rejected":
                    self.assertEqual(record["targets"], [])
                else:
                    self.assertTrue(record["targets"])
                    self.assertTrue(set(record["targets"]) <= current_ids)

    def test_denied_project_token_and_generated_artifact_marker_fail_closed(self):
        cases = (
            ("project token", "\nRoundfix\n", "source-baseline.project-token.denied"),
            (
                "generated marker",
                "\n<!-- setup-context-driven:root.core:start -->\n",
                "source-baseline.generated-artifact.denied",
            ),
            ("machine path", "\n/Users/example/repository\n", "source-baseline.path-token.denied"),
        )
        for name, injected, diagnostic in cases:
            with self.subTest(name=name), tempfile.TemporaryDirectory() as temp_dir:
                temp_root = Path(temp_dir) / "setup-context-driven"
                shutil.copytree(SKILL_ROOT / "assets" / "source-baselines", temp_root)
                corpus_file = temp_root / BASELINE_ID / "corpus" / "AGENTS.md"
                corpus_file.write_text(
                    corpus_file.read_text(encoding="utf-8") + injected,
                    encoding="utf-8",
                )

                with self.assertRaises(SourceBaselineValidationError) as captured:
                    load_source_baseline(temp_root, BASELINE_ID)

                self.assertIn(diagnostic, str(captured.exception))

    def test_source_catalog_emits_only_strict_0_0_1_documents(self):
        for path in sorted((SKILL_ROOT / "assets" / "source-baselines").rglob("*.json")):
            with self.subTest(path=path.relative_to(SKILL_ROOT)):
                document = self._read_json(path)
                self.assertIsInstance(document["schemaVersion"], str)
                self.assertTrue(document["schemaVersion"].endswith("/0.0.1"))
                self.assertEqual(document["version"], SOURCE_BASELINE_VERSION)

    def _entry_content(self, baseline):
        result = {}
        for entry in baseline.entries:
            data = (BASELINE_ROOT / entry.path).read_bytes()
            result[entry.entry_id] = data[entry.start : entry.end].decode("utf-8")
        return result

    def _read_json(self, path):
        return json.loads(path.read_text(encoding="utf-8"))


if __name__ == "__main__":
    unittest.main()

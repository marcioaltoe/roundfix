import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_audit import install_profile_skills, snapshot_files  # noqa: E402


class SpecTriageDecisionTests(unittest.TestCase):
    def test_spec_and_external_triage_combinations_apply_cleanly(self):
        cases = [
            (False, False),
            (False, True),
            (True, False),
            (True, True),
        ]

        for spec_enabled, triage_enabled in cases:
            with self.subTest(spec=spec_enabled, triage=triage_enabled):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)

                    first = run_apply(
                        repo,
                        decisions_for(spec_enabled, triage_enabled),
                    )
                    self.assertEqual(first.returncode, 0, first.stderr)
                    install_profile_skills(repo, "rust-cli")

                    self.assert_expected_guidance(repo, spec_enabled, triage_enabled)
                    self.assert_clean_audit(repo)
                    after_first = snapshot_files(repo)

                    second = run_apply(
                        repo,
                        decisions_for(spec_enabled, triage_enabled),
                    )

                    self.assertEqual(second.returncode, 0, second.stderr)
                    self.assertEqual(snapshot_files(repo), after_first)

    def test_false_values_remove_only_setup_owned_blocks(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            enabled = run_apply(repo, decisions_for(True, True))
            self.assertEqual(enabled.returncode, 0, enabled.stderr)

            root = repo / "AGENTS.md"
            spec_guide = repo / "docs" / "agents" / "spec-routing.md"
            triage_guide = repo / "docs" / "agents" / "external-triage.md"
            root.write_text(
                "repo root note\n" + root.read_text(encoding="utf-8") + "repo root tail\n",
                encoding="utf-8",
            )
            spec_guide.write_text(
                "repo spec note\n"
                + spec_guide.read_text(encoding="utf-8")
                + "repo spec tail\n",
                encoding="utf-8",
            )
            triage_guide.write_text(
                "repo triage note\n"
                + triage_guide.read_text(encoding="utf-8")
                + "repo triage tail\n",
                encoding="utf-8",
            )

            disabled = run_apply(repo, decisions_for(False, False))

            self.assertEqual(disabled.returncode, 0, disabled.stderr)
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(
                    encoding="utf-8"
                )
            )
            self.assertEqual(manifest["profile"], "rust-cli")
            self.assertNotIn("spec-workflow", manifest["modules"])
            self.assertNotIn("external-triage", manifest["modules"])

            root_content = root.read_text(encoding="utf-8")
            self.assertIn("repo root note\n", root_content)
            self.assertIn("repo root tail\n", root_content)
            self.assertNotIn("id=root.spec-workflow", root_content)
            self.assertNotIn("id=root.external-triage", root_content)
            self.assertEqual(spec_guide.read_text(encoding="utf-8"), "repo spec note\nrepo spec tail\n")
            self.assertEqual(
                triage_guide.read_text(encoding="utf-8"),
                "repo triage note\nrepo triage tail\n",
            )

    def assert_expected_guidance(self, repo, spec_enabled, triage_enabled):
        root_content = (repo / "AGENTS.md").read_text(encoding="utf-8")
        docs_layout_path = repo / "docs" / "agents" / "docs-layout.md"
        docs_layout = docs_layout_path.read_text(encoding="utf-8")

        self.assertIn("id=root.context-workflow", root_content)
        self.assertIn("docs/agents/domain.md", root_content)
        self.assertIn("docs/agents/docs-layout.md", root_content)
        self.assertTrue((repo / "docs" / "agents" / "domain.md").is_file())
        self.assertTrue(docs_layout_path.is_file())
        self.assertIn("docs/_inbox/", docs_layout)
        self.assertIn("docs/agents/", docs_layout)
        self.assertIn("pending | partial | deferred | done", docs_layout)
        self.assertIn(
            "Set `status: done` as soon as the Spec",
            docs_layout,
        )

        manifest = json.loads(
            (repo / "docs" / "agents" / "setup-context.json").read_text(
                encoding="utf-8"
            )
        )
        self.assertEqual(manifest["profile"], "rust-cli")

        if spec_enabled:
            self.assertIn("spec-workflow", manifest["modules"])
            self.assertIn("id=root.spec-workflow", root_content)
            self.assertIn("docs/specs/", docs_layout)
            self.assertTrue((repo / "docs" / "agents" / "spec-routing.md").is_file())
            self.assertTrue((repo / "docs" / "agents" / "issue-tracker.md").is_file())
        else:
            self.assertNotIn("spec-workflow", manifest["modules"])
            self.assertNotIn("id=root.spec-workflow", root_content)
            self.assertNotIn("docs/specs/", docs_layout)
            self.assertFalse((repo / "docs" / "agents" / "spec-routing.md").exists())
            self.assertFalse((repo / "docs" / "agents" / "issue-tracker.md").exists())

        if triage_enabled:
            triage_path = repo / "docs" / "agents" / "external-triage.md"
            self.assertIn("external-triage", manifest["modules"])
            self.assertIn("id=root.external-triage", root_content)
            self.assertEqual(root_content.count("id=root.external-triage"), 2)
            self.assertTrue(triage_path.is_file())
            triage_content = triage_path.read_text(encoding="utf-8")
            self.assertIn("# External triage", triage_content)
            self.assertIn("English", triage_content)
        else:
            self.assertNotIn("external-triage", manifest["modules"])
            self.assertNotIn("id=root.external-triage", root_content)
            self.assertFalse((repo / "docs" / "agents" / "external-triage.md").exists())

    def assert_clean_audit(self, repo):
        result = run_audit(repo)
        self.assertEqual(result.returncode, 0, result.stderr)
        payload = json.loads(result.stdout)
        self.assertEqual(payload["findings"], [])


def decisions_for(spec_enabled, triage_enabled):
    return [
        f"spec.scaffold={str(spec_enabled).lower()}",
        "domain.layout=single-context",
        f"triage.external={str(triage_enabled).lower()}",
        "autonomous.enabled=true",
        "runtime.backend=codex gpt-5.5 xhigh",
        "runtime.design=claude opus xhigh",
        "verification.gate=make verify",
        "language.generated=English",
        "secondbrain.enabled=false",
    ]


def run_apply(repo, decisions):
    args = ["apply", "--repo", str(repo), "--format", "json", "--profile", "rust-cli"]
    for decision in decisions:
        args.extend(["--decision", decision])
    result = run_context_setup(*args)
    payload = json.loads(result.stdout)
    if result.returncode == 3 and any(
        finding["code"] == "plan.confirmation.required" for finding in payload["findings"]
    ):
        args.extend(["--confirm-plan", payload["planDigest"]])
        return run_context_setup(*args)
    return result


def run_audit(repo):
    return run_context_setup("audit", "--repo", str(repo), "--format", "json")


def run_context_setup(*args):
    return subprocess.run(
        [sys.executable, str(SCRIPT), *args],
        text=True,
        capture_output=True,
        check=False,
    )


if __name__ == "__main__":
    unittest.main()

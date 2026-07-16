import json
import subprocess
import sys
import tempfile
import unittest
from hashlib import sha256
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_setup import parse_managed_blocks  # noqa: E402
from test_audit import install_profile_skills, snapshot_files  # noqa: E402


class DecisionRenderingTests(unittest.TestCase):
    def test_domain_layout_selects_distinct_guidance_and_audits_clean(self):
        cases = [
            ("single-context", "single `CONTEXT.md`", "CONTEXT-MAP.md", False),
            ("multi-context", "CONTEXT-MAP.md", "single `CONTEXT.md`", True),
        ]

        for layout, expected_text, omitted_text, monorepo_expected in cases:
            with self.subTest(layout=layout):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)

                    result = run_apply(repo, decisions_for(domain_layout=layout, autonomous=False))
                    install_profile_skills(repo, "rust-cli")

                    self.assertEqual(result.returncode, 0, result.stderr)
                    domain = (repo / "docs" / "agents" / "domain.md").read_text(
                        encoding="utf-8"
                    )
                    self.assertIn(expected_text, domain)
                    self.assertNotIn(omitted_text, domain)
                    self.assertEqual(
                        (repo / "docs" / "agents" / "monorepo.md").exists(),
                        monorepo_expected,
                    )
                    audit = run_audit(repo)
                    self.assertEqual(audit.returncode, 0, audit.stderr)
                    self.assertEqual(json.loads(audit.stdout)["findings"], [])

    def test_runtime_and_verification_values_render_and_drive_manifest_digest(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            backend = "codex ``edge`` xhigh"
            design = "claude opus xhigh"
            verification = "cargo test --all-targets"

            result = run_apply(
                repo,
                decisions_for(
                    autonomous=True,
                    runtime_backend=backend,
                    runtime_design=design,
                    verification_gate=verification,
                ),
            )

            self.assertEqual(result.returncode, 0, result.stderr)
            guide_path = repo / "docs" / "agents" / "autonomous-work.md"
            guide = guide_path.read_text(encoding="utf-8")
            self.assertIn(f"```{backend}```", guide)
            self.assertIn(f"`{design}`", guide)
            self.assertIn(f"`{verification}`", guide)
            self.assertNotIn("{{runtime.backend}}", guide)
            self.assertNotIn("{{runtime.design}}", guide)
            self.assertNotIn("{{verification.gate}}", guide)

            guide_blocks, _ = parse_managed_blocks(Path("docs/agents/autonomous-work.md"), guide)
            manifest = read_manifest(repo)
            artifact = self.managed_artifact(manifest, "guide.autonomous-work")
            self.assertEqual(
                artifact["digest"],
                digest_for(guide_blocks["guide.autonomous-work"].body),
            )

    def test_unsafe_inline_decision_values_block_apply_without_writes(self):
        cases = [
            ("newline", "cargo test\nmake verify"),
            ("control", "cargo test\x01"),
            (
                "marker",
                "<!-- setup-context-driven:begin id=root.core version=1 -->",
            ),
        ]

        for name, unsafe_value in cases:
            with self.subTest(name=name):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    before = snapshot_files(repo)

                    result = run_apply(
                        repo,
                        decisions_for(
                            autonomous=True,
                            verification_gate=unsafe_value,
                        ),
                    )

                    self.assertEqual(result.returncode, 2)
                    self.assertEqual(snapshot_files(repo), before)
                    payload = json.loads(result.stdout)
                    findings = [
                        finding
                        for finding in payload["findings"]
                        if finding["managedId"] == "verification.gate"
                    ]
                    self.assertEqual(len(findings), 1, payload)
                    self.assertEqual(findings[0]["code"], "decision.value.unsafe")
                    self.assertEqual(findings[0]["severity"], "error")

    def test_language_policy_accepts_only_english_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)

            result = run_apply(repo, decisions_for(autonomous=False, language="Portuguese"))

            self.assertEqual(result.returncode, 3)
            self.assertEqual(snapshot_files(repo), before)
            payload = json.loads(result.stdout)
            findings = [
                finding
                for finding in payload["findings"]
                if finding["managedId"] == "language.generated"
            ]
            self.assertEqual(len(findings), 1, payload)
            self.assertEqual(findings[0]["code"], "decision.required")

    def test_audit_reports_rendered_drift_and_apply_refreshes_same_artifact(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = decisions_for(autonomous=True, verification_gate="cargo test --all-targets")
            applied = run_apply(repo, decisions)
            install_profile_skills(repo, "rust-cli")
            self.assertEqual(applied.returncode, 0, applied.stderr)
            guide_path = repo / "docs" / "agents" / "autonomous-work.md"
            guide_path.write_text(
                guide_path.read_text(encoding="utf-8").replace(
                    "cargo test --all-targets",
                    "make verify",
                ),
                encoding="utf-8",
            )
            before_audit = snapshot_files(repo)

            audit = run_audit(repo)

            self.assertEqual(snapshot_files(repo), before_audit)
            payload = json.loads(audit.stdout)
            self.assertFinding(payload, "managed.content.modified", "guide.autonomous-work")
            self.assertIn(
                {
                    "action": "refresh managed content",
                    "path": "docs/agents/autonomous-work.md",
                    "managedId": "guide.autonomous-work",
                },
                payload["plannedChanges"],
            )

            refreshed = run_apply(repo, decisions)

            self.assertEqual(refreshed.returncode, 0, refreshed.stderr)
            self.assertIn("cargo test --all-targets", guide_path.read_text(encoding="utf-8"))
            self.assertEqual(run_audit(repo).returncode, 0)

    def test_qa07_alternate_values_render_selected_guidance_and_are_idempotent(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = decisions_for(
                spec=False,
                domain_layout="multi-context",
                triage=True,
                autonomous=False,
                runtime_backend="codex alternate-backend xhigh",
                runtime_design="claude alternate-design xhigh",
                verification_gate="cargo test --all-targets",
            )

            first = run_apply(repo, decisions)
            install_profile_skills(repo, "rust-cli")
            after_first = snapshot_files(repo)
            second = run_apply(repo, decisions)

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(snapshot_files(repo), after_first)
            root = (repo / "AGENTS.md").read_text(encoding="utf-8")
            domain = (repo / "docs" / "agents" / "domain.md").read_text(encoding="utf-8")
            self.assertNotIn("id=root.spec-workflow", root)
            self.assertNotIn("id=root.autonomous-work", root)
            self.assertIn("id=root.external-triage", root)
            self.assertIn("CONTEXT-MAP.md", domain)
            self.assertFalse((repo / "docs" / "agents" / "autonomous-work.md").exists())
            audit = run_audit(repo)
            self.assertEqual(audit.returncode, 0, audit.stderr)
            self.assertEqual(json.loads(audit.stdout)["findings"], [])

    def assertFinding(self, payload, code, managed_id):
        matches = [
            finding
            for finding in payload["findings"]
            if finding["code"] == code and finding["managedId"] == managed_id
        ]
        self.assertEqual(len(matches), 1, payload)

    def managed_artifact(self, manifest, managed_id):
        for artifact in manifest["managedArtifacts"]:
            if artifact["id"] == managed_id:
                return artifact
        self.fail(f"missing managed artifact {managed_id}")


def decisions_for(
    *,
    spec=True,
    domain_layout="single-context",
    triage=False,
    autonomous=True,
    runtime_backend="codex gpt-5.5 xhigh",
    runtime_design="claude opus xhigh",
    verification_gate="make verify",
    language="English",
    secondbrain=False,
):
    return [
        f"spec.scaffold={str(spec).lower()}",
        f"domain.layout={domain_layout}",
        f"triage.external={str(triage).lower()}",
        f"autonomous.enabled={str(autonomous).lower()}",
        f"runtime.backend={runtime_backend}",
        f"runtime.design={runtime_design}",
        f"verification.gate={verification_gate}",
        f"language.generated={language}",
        f"secondbrain.enabled={str(secondbrain).lower()}",
    ]


def run_apply(repo, decisions):
    args = ["apply", "--repo", str(repo), "--format", "json", "--profile", "rust-cli"]
    for decision in decisions:
        args.extend(["--decision", decision])
    return run_context_setup(*args)


def run_audit(repo):
    return run_context_setup("audit", "--repo", str(repo), "--format", "json")


def read_manifest(repo):
    return json.loads(
        (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
    )


def digest_for(content):
    normalized = content.strip() + "\n"
    return sha256(normalized.encode("utf-8")).hexdigest()


def run_context_setup(*args):
    return subprocess.run(
        [sys.executable, str(SCRIPT), *args],
        text=True,
        capture_output=True,
        check=False,
    )


if __name__ == "__main__":
    unittest.main()

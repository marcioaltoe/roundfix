PRAGMA foreign_keys = ON;
CREATE TABLE runs (
  id TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  state TEXT NOT NULL,
  head_repository TEXT NOT NULL DEFAULT '',
  head_branch TEXT NOT NULL DEFAULT '',
  base_repository TEXT NOT NULL DEFAULT '',
  pr_number TEXT NOT NULL DEFAULT '',
  git_root TEXT NOT NULL,
  local_branch TEXT NOT NULL,
  head_sha TEXT NOT NULL DEFAULT '',
  artifact_dir TEXT NOT NULL DEFAULT '',
  work_dir TEXT,
  spec_slug TEXT NOT NULL DEFAULT '',
  agent TEXT NOT NULL DEFAULT '',
  model TEXT NOT NULL DEFAULT '',
  reasoning_effort TEXT NOT NULL DEFAULT '',
  owner_pid INTEGER,
  owner_identity TEXT,
  owner_identity_unproven INTEGER NOT NULL DEFAULT 0,
  stop_requested_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  completed_at TEXT NOT NULL DEFAULT ''
);
CREATE TABLE active_run_locks (
  target_kind TEXT NOT NULL,
  target_key TEXT NOT NULL,
  run_id TEXT NOT NULL UNIQUE,
  created_at TEXT NOT NULL,
  PRIMARY KEY (target_kind, target_key),
  FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
);
CREATE TABLE interactive_defaults (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE TABLE run_events (
  run_id TEXT NOT NULL,
  cursor INTEGER NOT NULL,
  batch INTEGER NOT NULL DEFAULT 0,
  source TEXT NOT NULL,
  kind TEXT NOT NULL,
  review_issue TEXT NOT NULL DEFAULT '',
  tool_id TEXT NOT NULL DEFAULT '',
  tool_state TEXT NOT NULL DEFAULT '',
  summary TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  payload TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (run_id, cursor),
  FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
);
CREATE TABLE run_agent_selections (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  run_id TEXT NOT NULL,
  scope_kind TEXT NOT NULL,
  scope_id TEXT NOT NULL,
  category TEXT NOT NULL,
  profile_source TEXT NOT NULL,
  attempt INTEGER NOT NULL,
  selection_role TEXT NOT NULL,
  fallback_index INTEGER NOT NULL DEFAULT 0,
  runtime TEXT NOT NULL,
  model TEXT NOT NULL,
  reasoning_effort TEXT NOT NULL,
  status TEXT NOT NULL,
  reason_code TEXT NOT NULL DEFAULT '',
  reason TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE,
  CHECK (scope_kind IN ('task', 'qa', 'review')),
  CHECK (selection_role IN ('preferred', 'fallback')),
  CHECK (status IN ('attempting', 'active', 'failed', 'closed')),
  CHECK (attempt > 0),
  CHECK ((selection_role = 'preferred' AND fallback_index = 0) OR (selection_role = 'fallback' AND fallback_index > 0)),
  UNIQUE (run_id, scope_kind, scope_id, attempt)
);
PRAGMA user_version = 12;

INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository, pr_number, git_root, local_branch, head_sha, artifact_dir, created_at, updated_at, completed_at) VALUES
  ('run_active', 'review', 'Active', 'qa/a', 'feature/active', 'qa/a', '1', '/private/tmp/roundfix-qa-0059-20260806-03/repos/a', 'feature/active', 'a1', '/private/tmp/roundfix-qa-0059-20260806-03/.roundfix/artifacts/0fd63d8acd15e868', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z', ''),
  ('run_old', 'review', 'Clean', 'qa/b', 'feature/old', 'qa/b', '2', '/private/tmp/roundfix-qa-0059-20260806-03/repos/b', 'feature/old', 'b1', '/private/tmp/roundfix-qa-0059-20260806-03/.roundfix/artifacts/228fb27fda1038db', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z', '2026-01-01T01:00:00Z'),
  ('run_recent', 'review', 'Clean', 'qa/b', 'feature/recent', 'qa/b', '3', '/private/tmp/roundfix-qa-0059-20260806-03/repos/b', 'feature/recent', 'b2', '/private/tmp/roundfix-qa-0059-20260806-03/.roundfix/artifacts/228fb27fda1038db', '2026-08-06T04:00:00Z', '2026-08-06T04:00:00Z', '2026-08-06T04:01:00Z'),
  ('run_missing', 'review', 'Clean', 'qa/c', 'feature/missing', 'qa/c', '4', '/private/tmp/roundfix-qa-0059-20260806-03/repos/c', 'feature/missing', 'c1', '/private/tmp/roundfix-qa-0059-20260806-03/.roundfix/artifacts/e29f5ec167ad72be', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z', '2026-01-01T01:00:00Z'),
  ('run_overridden', 'review', 'Clean', 'qa/d', 'feature/overridden', 'qa/d', '5', '/private/tmp/roundfix-qa-0059-20260806-03/repos/d', 'feature/overridden', 'd1', '/private/tmp/roundfix-qa-0059-20260806-03/.roundfix/custom-overridden', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z', '2026-01-01T01:00:00Z'),
  ('run_outside', 'review', 'Clean', 'qa/e', 'feature/outside', 'qa/e', '6', '/private/tmp/roundfix-qa-0059-20260806-03/repos/e', 'feature/outside', 'e1', '/private/tmp/roundfix-qa-0059-20260806-03-outside', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z', '2026-01-01T01:00:00Z'),
  ('run_unsafe', 'review', 'Clean', 'qa/f', 'feature/unsafe', 'qa/f', '7', '/private/tmp/roundfix-qa-0059-20260806-03/repos/f', 'feature/unsafe', 'f1', '/private/tmp/roundfix-qa-0059-20260806-03/.roundfix', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z', '2026-01-01T01:00:00Z');

INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
VALUES ('pr', 'qa/a#feature/active', 'run_active', '2026-01-01T00:00:00Z');

INSERT INTO run_events (run_id, cursor, batch, source, kind, summary, created_at)
VALUES
  ('run_active', 1, 1, 'agent', 'agent-message', 'active event', '2026-01-01T00:00:00Z'),
  ('run_old', 1, 1, 'agent', 'agent-message', 'old event', '2026-01-01T00:00:00Z');

INSERT INTO run_agent_selections (run_id, scope_kind, scope_id, category, profile_source, attempt, selection_role, runtime, model, reasoning_effort, status, created_at)
VALUES
  ('run_active', 'review', 'batch-1', 'review', 'qa', 1, 'preferred', 'codex', 'qa-model', 'high', 'active', '2026-01-01T00:00:00Z'),
  ('run_old', 'review', 'batch-1', 'review', 'qa', 1, 'preferred', 'codex', 'qa-model', 'high', 'closed', '2026-01-01T00:00:00Z');

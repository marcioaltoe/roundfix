SHELL := /bin/bash

.PHONY: help sync lock-check fmt fmt-check lint typecheck test verify verify-incremental

help:
	@printf '%s\n' 'Targets: help sync lock-check fmt fmt-check lint typecheck test verify verify-incremental'

sync:
	uv sync --group dev

lock-check:
	uv lock --check --offline

fmt:
	uv run --frozen ruff format src tests

fmt-check:
	uv run --frozen ruff format --check src tests

lint:
	uv run --frozen ruff check src tests

typecheck:
	uv run --frozen ty check --error-on-warning src tests

test:
	uv run --frozen pytest -q

verify: lock-check fmt-check lint typecheck test

verify-incremental: verify

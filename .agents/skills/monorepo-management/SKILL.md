---
name: monorepo-management
description:
  Master monorepo management with Bun workspaces and Turborepo to build efficient, scalable multi-package repositories
  with optimized builds and dependency management. Use when setting up monorepos, optimizing builds, or managing shared
  dependencies.
---

# Monorepo Management

Build efficient, scalable monorepos that enable code sharing, consistent tooling, and atomic changes across multiple
packages and applications.

## When to Use This Skill

- Setting up new monorepo projects
- Migrating from multi-repo to monorepo
- Optimizing build and test performance
- Managing shared dependencies
- Implementing code sharing strategies
- Setting up CI/CD for monorepos
- Versioning and publishing packages
- Debugging monorepo-specific issues

## Core Concepts

### 1. Why Monorepos?

**Advantages:**

- Shared code and dependencies
- Atomic commits across projects
- Consistent tooling and standards
- Easier refactoring
- Simplified dependency management
- Better code visibility

**Challenges:**

- Build performance at scale
- CI/CD complexity
- Access control
- Large Git repository

### 2. Monorepo Tools

**Runtime / Package Manager:**

- Bun with npm-style workspaces defined in the root `package.json` (recommended)

**Build System:**

- Turborepo as the single orchestrator for builds, tests, lint, dev and typecheck across all workspaces

## Turborepo Setup

### Initial Setup

```bash
# Create new monorepo
bunx create-turbo@latest my-monorepo
cd my-monorepo

# Structure:
# apps/
#   web/          - Next.js app
#   docs/         - Documentation site
# packages/
#   ui/           - Shared UI components
#   config/       - Shared configurations
#   tsconfig/     - Shared TypeScript configs
# turbo.json      - Turborepo configuration
# package.json    - Root package.json
```

### Configuration

```json
// turbo.json
{
  "$schema": "https://turbo.build/schema.json",
  "globalDependencies": ["**/.env.*local"],
  "pipeline": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**", ".next/**", "!.next/cache/**"]
    },
    "test": {
      "dependsOn": ["build"],
      "outputs": ["coverage/**"]
    },
    "lint": {
      "outputs": []
    },
    "dev": {
      "cache": false,
      "persistent": true
    },
    "typecheck": {
      "dependsOn": ["^build"],
      "outputs": []
    }
  }
}
```

```json
// package.json (root)
{
  "name": "my-monorepo",
  "private": true,
  "workspaces": ["apps/*", "packages/*"],
  "scripts": {
    "build": "turbo run build",
    "dev": "turbo run dev",
    "test": "turbo run test",
    "lint": "turbo run lint",
    "format": "oxfmt format .",
    "clean": "turbo run clean && rm -rf node_modules"
  },
  "devDependencies": {
    "turbo": "^2.0.0",
    "typescript": "^5.0.0"
  },
  "packageManager": "bun@1.0.0"
}
```

### Package Structure

```json
// packages/ui/package.json
{
  "name": "@repo/ui",
  "version": "0.0.0",
  "private": true,
  "main": "./dist/index.js",
  "types": "./dist/index.d.ts",
  "exports": {
    ".": {
      "import": "./dist/index.js",
      "types": "./dist/index.d.ts"
    },
    "./button": {
      "import": "./dist/button.js",
      "types": "./dist/button.d.ts"
    }
  },
  "scripts": {
    "build": "tsup src/index.ts --format esm,cjs --dts",
    "dev": "tsup src/index.ts --format esm,cjs --dts --watch",
    "lint": "eslint src/",
    "typecheck": "tsgo --noEmit"
  },
  "devDependencies": {
    "@repo/tsconfig": "workspace:*",
    "tsup": "^7.0.0",
    "typescript": "^5.0.0"
  },
  "dependencies": {
    "react": "^18.2.0"
  }
}
```

## Bun Workspaces

### Setup

Workspaces are configured in the root `package.json` and Bun uses that metadata to install and link dependencies:

```json
// package.json (root)
{
  "name": "my-monorepo",
  "private": true,
  "workspaces": ["apps/*", "packages/*"],
  "packageManager": "bun@1.0.0"
}
```

### Dependency Management

```bash
# Install all workspace dependencies
bun install

# Add dependency to a specific workspace
bun add react --filter @repo/ui

# Add dev dependency
bun add -d typescript --filter @repo/ui

# Add internal workspace dependency
bun add @repo/ui --filter web

# Remove dependency
bun remove react --filter @repo/ui
```

### Running Scripts

```bash
# Run scripts through package.json using Bun + Turborepo
bun run build          # turbo run build
bun run dev            # turbo run dev
bun run test           # turbo run test
bun run lint           # turbo run lint

# Run script in a specific workspace
bun run dev --filter web
bun run build --filter @repo/ui
```

## Shared Configurations

### TypeScript Configuration

```json
// packages/tsconfig/base.json
{
  "compilerOptions": {
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "incremental": true,
    "declaration": true
  },
  "exclude": ["node_modules"]
}

// packages/tsconfig/react.json
{
  "extends": "./base.json",
  "compilerOptions": {
    "jsx": "react-jsx",
    "lib": ["ES2022", "DOM", "DOM.Iterable"]
  }
}

// apps/web/tsconfig.json
{
  "extends": "@repo/tsconfig/react.json",
  "compilerOptions": {
    "outDir": "dist",
    "rootDir": "src"
  },
  "include": ["src"],
  "exclude": ["node_modules", "dist"]
}
```

### oxlint Configuration

```json
// .oxlintrc.json (root — shared across all workspaces)
{
  "$schema": "./node_modules/oxlint/configuration_schema.json",
  "plugins": ["unicorn", "typescript", "oxc"],
  "rules": {
    "no-unused-vars": "error",
    "typescript/no-unused-vars": "error"
  }
}
```

## Code Sharing Patterns

### Pattern 1: Shared UI Components

```typescript
// packages/ui/src/button.tsx
import * as React from 'react';

export interface ButtonProps {
  variant?: 'primary' | 'secondary';
  children: React.ReactNode;
  onClick?: () => void;
}

export function Button({ variant = 'primary', children, onClick }: ButtonProps) {
  return (
    <button
      className={`btn btn-${variant}`}
      onClick={onClick}
    >
      {children}
    </button>
  );
}

// packages/ui/src/index.ts
export { Button, type ButtonProps } from './button';
export { Input, type InputProps } from './input';

// apps/web/src/app.tsx
import { Button } from '@repo/ui';

export function App() {
  return <Button variant="primary">Click me</Button>;
}
```

### Pattern 2: Shared Utilities

```typescript
// packages/utils/src/string.ts
export function capitalize(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

export function truncate(str: string, length: number): string {
  return str.length > length ? str.slice(0, length) + "..." : str;
}

// packages/utils/src/index.ts
export * from "./string";
export * from "./array";
export * from "./date";

// Usage in apps
import { capitalize, truncate } from "@repo/utils";
```

### Pattern 3: Shared Types

```typescript
// packages/types/src/user.ts
export interface User {
  id: string;
  email: string;
  name: string;
  role: "admin" | "user";
}

export interface CreateUserInput {
  email: string;
  name: string;
  password: string;
}

// Used in both frontend and backend
import type { User, CreateUserInput } from "@repo/types";
```

## Build Optimization

### Turborepo Caching

```json
// turbo.json
{
  "pipeline": {
    "build": {
      // Build depends on dependencies being built first
      "dependsOn": ["^build"],

      // Cache these outputs
      "outputs": ["dist/**", ".next/**"],

      // Cache based on these inputs (default: all files)
      "inputs": ["src/**/*.tsx", "src/**/*.ts", "package.json"]
    },
    "test": {
      // Run tests in parallel, don't depend on build
      "cache": true,
      "outputs": ["coverage/**"]
    }
  }
}
```

### Remote Caching

```bash
# Turborepo Remote Cache (Vercel)
bunx turbo login
bunx turbo link

# Custom remote cache
# turbo.json
{
  "remoteCache": {
    "signature": true,
    "enabled": true
  }
}
```

## CI/CD for Monorepos

### GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0 # For Nx affected commands

      - uses: oven-sh/setup-bun@v2
        with:
          bun-version: 1.0.0

      - name: Install dependencies
        run: bun install --frozen-lockfile

      - name: Build
        run: bun run build

      - name: Test
        run: bun run test

      - name: Lint
        run: bun run lint

      - name: Type check
        run: bun run typecheck
```

### Local Developer Workflow (Optional Wrapper)

Some monorepos wrap `turbo`/package-manager scripts behind a root `Makefile` (or similar) to keep workflows consistent.

Guidelines:

- Prefer the project’s documented entrypoints (for example `make`, `just`, `task`, or `npm`/`pnpm`/`bun` scripts).
- If a wrapper exists, use it instead of ad hoc commands to avoid bypassing required quality gates.
- If no wrapper exists, use `turbo run <task>` (or the repo’s equivalent) consistently across workspaces.

## Best Practices

1. **Consistent Versioning**: Lock dependency versions across workspace
2. **Shared Configs**: Centralize TypeScript, oxlint, oxfmt configs
3. **Dependency Graph**: Keep it acyclic, avoid circular dependencies
4. **Cache Effectively**: Configure inputs/outputs correctly
5. **Type Safety**: Share types between frontend/backend
6. **Testing Strategy**: Unit tests in packages, E2E in apps
7. **Documentation**: README in each package
8. **Release Strategy**: Use changesets for versioning

## Common Pitfalls

- **Circular Dependencies**: A depends on B, B depends on A
- **Phantom Dependencies**: Using deps not in package.json
- **Incorrect Cache Inputs**: Missing files in Turborepo inputs
- **Over-Sharing**: Sharing code that should be separate
- **Under-Sharing**: Duplicating code across packages
- **Large Monorepos**: Without proper tooling, builds slow down

## Publishing Packages

```bash
# Using Changesets with Bun
bun add -d @changesets/cli
npx changeset init

# Create changeset
npx changeset

# Version packages
npx changeset version

# Publish (script typically wraps `changeset publish`)
bun run release
```

# Agent protocols

`AGENTS.md` (at the repo root) is a lightweight entry point. The standards files
it references are read on demand by the agent when it is asked to file issues,
open pull requests, or write commits, keeping the always-loaded context small.

- `docs/github-issue-standards.md`: classification rules for defects,
  features, and tasks, plus pointers to the issue forms.
- `docs/github-pr-standards.md`: how to fill the PR template.
- `docs/git-standards.md`: Tim Pope-style commit messages and the
  branch-naming convention.

## Adopt: clone the repo locally (one-time)

```sh
git clone git@github.com:pyck-ai/.github.git ~/code/pyck-github
```

The path is up to you; the integrations below reference whatever location you
pick.

## Use with Claude Code

Claude Code reads `CLAUDE.md` files automatically (from the project root, its
parent directories, and `~/.claude/CLAUDE.md`) and follows any `@path` imports
inside them. It does **not** automatically load this repo's `AGENTS.md`; you
point it there with an import in one of those `CLAUDE.md` files.

**Per project**

Add an import to the project's `CLAUDE.md` pointing at your clone (use an
absolute path, so the on-demand links inside `AGENTS.md` resolve against the
clone):

```sh
echo '@~/code/pyck-github/AGENTS.md' >> CLAUDE.md
```

Updates are then a `git pull` away.

**Globally (apply to every project on your machine)**

Add a single import line to `~/.claude/CLAUDE.md`, which Claude Code reads for
every session:

```
@~/code/pyck-github/AGENTS.md
```

**Verify**: ask Claude to commit a change. It should reference
`git-standards.md` when drafting the message.

## Use with opencode

opencode loads files listed in the `instructions[]` array of
`~/.config/opencode/opencode.json` into every session, globally. Once
configured, the standards apply across every project you open.

**Minimal config**

```json
{
  "$schema": "https://opencode.ai/config.json",
  "instructions": [
    "~/code/pyck-github/AGENTS.md"
  ],
  "permission": {
    "external_directory": {
      "~/code/pyck-github/**": "allow"
    }
  }
}
```

The `permission.external_directory` entry is required so opencode can read the
standards docs that `AGENTS.md` lazy-loads from outside the project root.

**Full reference**

A complete config (including the MCP servers, agent variants, and auth plugin
used by the Pyck team) is in
[`.opencode/opencode.example.json`](../.opencode/opencode.example.json). Copy it
to `~/.config/opencode/opencode.json` and replace:

- `<path-to-pyck-github-checkout>` with your local clone path.
- `<YOUR_CLICKHOUSE_MCP_TOKEN>` with your ClickHouse MCP bearer token.

**Verify**: ask opencode to open a pull request. It should follow
`github-pr-standards.md`.

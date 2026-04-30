# pyck-ai/.github

This repository ships two things for the [pyck-ai](https://github.com/pyck-ai)
GitHub organization:

1. **Org-wide community-health defaults** — issue forms, a PR template, a
   security policy, and a code of conduct that GitHub auto-applies to every
   repo in the org that doesn't ship its own copy.
2. **Opt-in agent protocols** — a set of standards (`.claude/`) that
   AI coding assistants like [Claude Code](https://www.anthropic.com/claude-code)
   and [opencode](https://opencode.ai) load to keep issues, pull requests,
   and commits consistent across humans and agents.

## What's in here

| Path | Purpose | Applied |
|------|---------|---------|
| `.github/CODE_OF_CONDUCT.md` | Community code of conduct | Org-wide auto |
| `.github/SECURITY.md` | Vulnerability disclosure policy | Org-wide auto |
| `.github/PULL_REQUEST_TEMPLATE.md` | Default PR template | Org-wide auto |
| `.github/ISSUE_TEMPLATE/*.yml` | Defect, feature, task, epic, initiative forms | Org-wide auto |
| `.claude/AGENTS.md` | Agent protocol entry point (lazy-loads the standards below) | Opt-in per project |
| `.claude/ISSUE_STANDARDS.md` | Issue classification + per-type templates | Opt-in |
| `.claude/PR_STANDARDS.md` | PR template usage rules | Opt-in |
| `.claude/GIT_STANDARDS.md` | Commit message + branch naming rules | Opt-in |
| `opencode.example.json` | Reference opencode config | Reference |

## Org-wide community-health files

GitHub automatically applies the files under `.github/` to every repository
in the `pyck-ai` org that doesn't ship its own copy — no setup required by
consumers. Any repo can override a default by adding the same file at its
own root.

## Agent protocols (`.claude/`)

`.claude/AGENTS.md` is a lightweight entry point. The three standards files
it references are read on demand by the agent when it is asked to file
issues, open pull requests, or write commits — keeping the always-loaded
context small.

- `.claude/ISSUE_STANDARDS.md` — classification rules and templates for
  defects, features, tasks, epics, and initiatives.
- `.claude/PR_STANDARDS.md` — the canonical PR body template and how to
  fill it.
- `.claude/GIT_STANDARDS.md` — Tim Pope-style commit messages and the
  branch-naming convention.

### Adopt: clone the repo locally (one-time)

```sh
git clone git@github.com:pyck-ai/.github.git ~/code/pyck-github
```

The path is up to you; the integrations below reference whatever location
you pick.

### Use with Claude Code

Claude Code automatically loads `CLAUDE.md` and `AGENTS.md` from the project
root and from a `.claude/` directory at session start.

**Per project**

Symlink the shared `.claude/` into your project so updates are a `git pull`
away:

```sh
ln -s ~/code/pyck-github/.claude .claude
```

Or vendor a copy if you'd rather pin the version:

```sh
cp -r ~/code/pyck-github/.claude .claude
```

**Globally (apply to every project on your machine)**

Claude Code reads `~/.claude/CLAUDE.md` for every session and supports
`@path` imports. Add a single line:

```
@~/code/pyck-github/.claude/AGENTS.md
```

**Verify** — ask Claude to commit a change. It should reference
`GIT_STANDARDS.md` when drafting the message.

### Use with opencode

opencode loads files listed in the `instructions[]` array of
`~/.config/opencode/opencode.json` into every session, globally — once
configured, the standards apply across every project you open.

**Minimal config**

```json
{
  "$schema": "https://opencode.ai/config.json",
  "instructions": [
    "~/code/pyck-github/.claude/AGENTS.md"
  ],
  "permission": {
    "external_directory": {
      "~/code/pyck-github/**": "allow"
    }
  }
}
```

The `permission.external_directory` entry is required so opencode can read
the standards docs that `AGENTS.md` lazy-loads from outside the project
root.

**Full reference**

A complete config — including the MCP servers, agent variants, and auth
plugin used by the pyck-ai team — is in [`opencode.example.json`](./opencode.example.json).
Copy it to `~/.config/opencode/opencode.json` and replace:

- `<path-to-pyck-github-checkout>` with your local clone path.
- `<YOUR_CLICKHOUSE_MCP_TOKEN>` with your ClickHouse MCP bearer token.

**Verify** — ask opencode to open a pull request. It should follow
`PR_STANDARDS.md`.

## Updating

```sh
cd ~/code/pyck-github && git pull
```

If you vendored `.claude/` instead of symlinking it, re-copy from the
updated checkout and review the diff.

## Contributing

Pull requests and commits in this repo follow the same standards it
publishes:

- [`.claude/ISSUE_STANDARDS.md`](./.claude/ISSUE_STANDARDS.md)
- [`.claude/PR_STANDARDS.md`](./.claude/PR_STANDARDS.md)
- [`.claude/GIT_STANDARDS.md`](./.claude/GIT_STANDARDS.md)

Changes to files under `.github/` propagate to every repository in the
`pyck-ai` org — review carefully.

## Security

See [`.github/SECURITY.md`](./.github/SECURITY.md) for the vulnerability
disclosure process.

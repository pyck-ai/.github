# Contributing

Pull requests and commits in this repo follow the same standards it publishes:

- [`github-issue-standards.md`](./github-issue-standards.md)
- [`github-pr-standards.md`](./github-pr-standards.md)
- [`git-standards.md`](./git-standards.md)

Not everything here behaves the same way. The community-health defaults (issue
forms, the PR template, `SECURITY.md`, `CODE_OF_CONDUCT.md`) are auto-propagated
by GitHub to every repo in the org that lacks its own copy, so changes to them
affect the whole org: review those carefully. The agent protocols do not
propagate; they apply by opt-in. [Renovate](./renovate.md) is not part of this
repo at all: it runs from `pyck-ai/renovate`.

## Generated files

The issue forms and the PR template under `.github/` are **generated** from
canonical sources in [`src/`](../src/) by
[`scripts/generate-community-files.go`](../scripts/generate-community-files.go).
Edit the source, not the generated copy, then re-sync:

```sh
task generate
```

CI runs `task generate:check` on every pull request and fails if a generated
file is out of date.

## Reusable actions

CI automation is written as reusable composite actions, not one-off inline
shell steps inside workflows, so other repos in the org can adopt the same
logic with a `uses:` step instead of copy-pasting it. New shared automation
belongs in
[`pyck-ai/github-actions`](https://github.com/pyck-ai/github-actions); a
workflow in this repo should stay a thin consumer that calls it.

## Updating (consumers)

If you adopted the standards by cloning this repo, pull the latest:

```sh
cd ~/code/pyck-github && git pull
```

If you vendored a copy instead of importing from the clone, re-copy from the
updated checkout and review the diff.
